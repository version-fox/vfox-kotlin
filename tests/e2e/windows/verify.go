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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

func verifySearch(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox, version string) {
	t.Helper()
	out := common.StripANSI(execRetry(ctx, t, ctr, "vfox search kotlin", "vfox search kotlin", box.vars...))
	require.Containsf(t, out, "- "+version, "search must list %s\n--- output ---\n%s", version, out)
	require.Contains(t, out, "(installed)", "search must mark the installed release\n--- output ---\n"+out)
	t.Logf("PASS: vfox search lists kotlin@%s as installed", version)
}

func useKotlin(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox, activation, version string) {
	t.Helper()
	execOK(ctx, t, ctr, activated(activation, fmt.Sprintf("vfox use --global kotlin@%s", version)), box.vars...)
	t.Logf("PASS: kotlin@%s selected globally", version)
}

// verifyToolchain checks EnvKeys and the compiler on Windows. Native stderr
// is not redirected with 2>&1: pwsh wraps it in an ErrorRecord and the
// prelude would abort the script.
func verifyToolchain(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox, activation, version string) {
	t.Helper()

	script := activatedEnv(activation,
		`$cmd = Get-Command kotlinc; `+
			`if (-not $cmd) { throw 'kotlinc not on PATH' }; `+
			`$cmd.Source; `+
			`Write-Output ('PATH=' + $env:PATH); `+
			`kotlin -version; `+
			`kotlinc -version; `+
			`vfox current kotlin`)
	out := common.StripANSI(execOK(ctx, t, ctr, script, box.vars...))

	binDir := strings.TrimSpace(firstLineAfter(out, "kotlinc.bat"))
	require.NotEmpty(t, binDir, "kotlinc must resolve to a path\n--- output ---\n"+out)
	require.Containsf(t, out, "PATH=", "EnvKeys must add the SDK bin to PATH\n--- output ---\n%s", out)
	require.Containsf(t, strings.ReplaceAll(binDir, "/", `\`), `\bin\`,
		"kotlinc must resolve from a bin directory\n--- output ---\n%s", out)
	require.Containsf(t, out, "Kotlin version "+version,
		"kotlin -version must report %s\n--- output ---\n%s", version, out)
	require.Contains(t, out, version, "kotlinc/current must name the version\n--- output ---\n"+out)
	t.Logf("PASS: kotlin %s runs on Windows", version)
}

// firstLineAfter returns the first non-empty line at or after needle, used to
// pick the resolved kotlinc path out of the combined output.
func firstLineAfter(out, needle string) string {
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(needle)) {
			for j := i; j < len(lines); j++ {
				if s := strings.TrimSpace(lines[j]); s != "" {
					return s
				}
			}
		}
	}
	return ""
}

// verifyNative checks the Kotlin/Native addition on Windows: its .bat
// launcher must resolve from the added bin directory and report the release.
func verifyNative(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox, activation, version string) {
	t.Helper()

	script := activatedEnv(activation,
		`$cmd = Get-Command kotlinc-native; `+
			`if (-not $cmd) { throw 'kotlinc-native not on PATH' }; `+
			`$cmd.Source; `+
			`kotlinc-native -version`)
	out := common.StripANSI(execOK(ctx, t, ctr, script, box.vars...))
	require.Containsf(t, strings.ReplaceAll(out, "/", `\`), `\bin\kotlinc-native.bat`,
		"kotlinc-native must resolve from a bin directory\n--- output ---\n%s", out)
	require.Containsf(t, out, version, "kotlinc-native -version must report %s\n--- output ---\n%s", version, out)
	t.Logf("PASS: kotlinc-native %s is on PATH on Windows", version)
}

func checkMissingVersionRejected(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox) {
	t.Helper()
	code, out := execRun(ctx, t, ctr, "vfox install kotlin@0.0.0-e2e-missing", box.vars...)
	require.NotZerof(t, code, "installing a nonexistent release must fail\n--- output ---\n%s", out)
	t.Logf("PASS: a nonexistent release is rejected")
}
