// Copyright 2026 Han Li and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package windows

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// pwshPrelude makes pwsh fail on cmdlet and native-command errors. By default
// pwsh is non-terminating and would exit 0 even when a step failed.
const pwshPrelude = "$ErrorActionPreference = 'Stop'; $PSNativeCommandUseErrorActionPreference = $true;"

func execRun(ctx context.Context, t *testing.T, ctr tc.Container, script string, extraEnv ...string) (int, string) {
	t.Helper()
	opts := []tcexec.ProcessOption{tcexec.Multiplexed()}
	if len(extraEnv) > 0 {
		opts = append(opts, tcexec.WithEnv(extraEnv))
	}
	// Bound every exec: a pwsh process that waits forever must fail in
	// minutes, not consume the whole job timeout.
	execCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	code, r, err := ctr.Exec(execCtx, []string{"pwsh", "-NoProfile", "-NonInteractive", "-Command", pwshPrelude + " " + script}, opts...)
	if err != nil {
		t.Fatalf("exec: %v\n--- script ---\n%s", err, script)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("exec read: %v\n--- script ---\n%s", err, script)
	}
	return code, string(out)
}

func execOK(ctx context.Context, t *testing.T, ctr tc.Container, script string, extraEnv ...string) string {
	t.Helper()
	code, out := execRun(ctx, t, ctr, script, extraEnv...)
	if code != 0 {
		t.Fatalf("exec exited %d\n--- script ---\n%s\n--- output ---\n%s", code, script, out)
	}
	return out
}

func execRetry(ctx context.Context, t *testing.T, ctr tc.Container, label, script string, extraEnv ...string) string {
	t.Helper()
	return common.Retry(ctx, t, label, func() (int, string) {
		return execRun(ctx, t, ctr, script, extraEnv...)
	})
}

// pwshQuote renders a single-quoted PowerShell string literal.
func pwshQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// activated prepends the captured `vfox activate pwsh` output. Every exec is
// a fresh process, so activation cannot persist between calls.
func activated(activation, script string) string {
	return activation + "\r\n" + script
}

// activatedEnv is activation plus an explicit `vfox env -s pwsh` export.
// vfox's pwsh activation refreshes the toolchain env through a prompt hook,
// which never fires in non-interactive `pwsh -Command` execs, so each call
// that consumes EnvKeys must pull the export itself.
func activatedEnv(activation, script string) string {
	return activation + "\r\n$__vfoxExport = vfox env -s pwsh | Out-String; " +
		"if ($__vfoxExport.Trim()) { Invoke-Expression -Command $__vfoxExport }\r\n" + script
}
