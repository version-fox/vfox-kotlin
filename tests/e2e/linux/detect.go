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

package linux

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// checkDetect verifies `vfox link --detect kotlin` against a compiler placed
// outside vfox's own directories. It copies the installed compiler to /opt,
// gives it a version no other install has, points KOTLIN_HOME at it and links
// everything the plugin finds.
func checkDetect(ctx context.Context, t *testing.T, ctr tc.Container) {
	t.Helper()
	binDir := strings.TrimSpace(execOK(ctx, t, ctr, activated(`dirname "$(command -v kotlinc)"`)))
	sdkHome := strings.TrimSuffix(binDir, "/bin")
	require.NotEmpty(t, sdkHome, "kotlinc must be installed before checking Detect")

	execOK(ctx, t, ctr, "mkdir -p /opt/kotlin-detect && cp -a "+shQuote(sdkHome)+"/. /opt/kotlin-detect/ && "+
		"printf '1.7.0-release-e2e' > /opt/kotlin-detect/build.txt && chmod +x /opt/kotlin-detect/bin/kotlinc")

	code, out := execRun(ctx, t, ctr, "vfox link --detect -y kotlin", "KOTLIN_HOME=/opt/kotlin-detect")
	require.Zerof(t, code, "vfox link --detect failed\n--- output ---\n%s", out)
	plain := common.StripANSI(out)
	require.Containsf(t, plain, "kotlin@1.7.0", "Detect must find the compiler under KOTLIN_HOME\n--- output ---\n%s", plain)
	require.Containsf(t, plain, "Linked", "Detect must link the compiler\n--- output ---\n%s", plain)

	list := common.StripANSI(execOK(ctx, t, ctr, activated("vfox list kotlin")))
	require.Containsf(t, list, "1.7.0 (linked", "linked compiler must appear in vfox list\n--- output ---\n%s", list)
	t.Logf("PASS: Detect found and linked a system Kotlin compiler")
}
