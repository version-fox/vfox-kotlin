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
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// checkDetect places a compiler outside vfox's directories, points
// KOTLIN_HOME at it and verifies `vfox link --detect kotlin` links it.
func checkDetect(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox, activation string) {
	t.Helper()
	const target = `C:\KotlinDetect`

	script := activatedEnv(activation,
		`$sdkHome = Split-Path (Split-Path (Get-Command kotlinc).Source); `+
			`New-Item -ItemType Directory -Force -Path '`+target+`' | Out-Null; `+
			`Copy-Item -Path (Join-Path $sdkHome '*') -Destination '`+target+`' -Recurse -Force; `+
			`Set-Content -Path '`+target+`\build.txt' -Value '1.7.0-release-e2e' -NoNewline; `+
			`Write-Output $sdkHome`)
	execOK(ctx, t, ctr, script, box.vars...)

	vars := append(append([]string{}, box.vars...), "KOTLIN_HOME="+target)
	code, out := execRun(ctx, t, ctr, "vfox link --detect -y kotlin", vars...)
	require.Zerof(t, code, "vfox link --detect failed\n--- output ---\n%s", out)
	plain := common.StripANSI(out)
	require.Containsf(t, plain, "kotlin@1.7.0", "Detect must find the compiler under KOTLIN_HOME\n--- output ---\n%s", plain)
	require.Containsf(t, plain, "Linked", "Detect must link the compiler\n--- output ---\n%s", plain)

	list := common.StripANSI(execOK(ctx, t, ctr, activatedEnv(activation, "vfox list kotlin"), box.vars...))
	require.Containsf(t, list, "1.7.0 (linked", "linked compiler must appear in vfox list\n--- output ---\n%s", list)
	t.Logf("PASS: Detect found and linked a system Kotlin compiler")
}
