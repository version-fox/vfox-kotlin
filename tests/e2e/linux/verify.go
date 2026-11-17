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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// verifySearch checks the Available hook through vfox search.
func verifySearch(ctx context.Context, t *testing.T, ctr tc.Container, version string) {
	t.Helper()
	out := common.StripANSI(execRetry(ctx, t, ctr, "vfox search kotlin", "vfox search kotlin"))
	require.Containsf(t, out, "- "+version, "search must list %s\n--- output ---\n%s", version, out)
	require.Contains(t, out, "(installed)", "search must mark the installed release\n--- output ---\n"+out)
	t.Logf("PASS: vfox search lists kotlin@%s as installed", version)
}

// useKotlin selects the installed compiler globally; activation is required
// for vfox use to know which shell to write to.
func useKotlin(ctx context.Context, t *testing.T, ctr tc.Container, version string) {
	t.Helper()
	execOK(ctx, t, ctr, activated(fmt.Sprintf("vfox use --global kotlin@%s", version)))
	t.Logf("PASS: kotlin@%s selected globally", version)
}

// verifyToolchain checks EnvKeys and the compiler end to end: the SDK's bin
// must be on PATH and kotlinc/kotlin must report the installed version.
func verifyToolchain(ctx context.Context, t *testing.T, ctr tc.Container, version string) {
	t.Helper()

	binDir := strings.TrimSpace(execOK(ctx, t, ctr, activated(`dirname "$(command -v kotlinc)"`)))
	require.NotEmpty(t, binDir, "kotlinc must be on PATH after vfox use")
	require.Truef(t, strings.HasSuffix(binDir, "/bin"), "unexpected kotlinc directory %q", binDir)

	path := execOK(ctx, t, ctr, activated(`printf '%s' "$PATH"`))
	require.Containsf(t, path, binDir, "EnvKeys must add %s to PATH\n--- PATH ---\n%s", binDir, path)

	// vfox env reports the cache SDK directory while command -v resolves
	// through the sdks symlink; both point at the same compiler.
	envJSON := execOK(ctx, t, ctr, activated("vfox env -s bash --json"))
	require.Containsf(t, envJSON, "kotlin-"+version,
		"vfox env must export the installed SDK\n--- output ---\n%s", envJSON)
	require.Containsf(t, envJSON, "/bin", "vfox env must export a bin directory\n--- output ---\n%s", envJSON)

	banner := common.StripANSI(execOK(ctx, t, ctr, activated("kotlin -version 2>&1")))
	require.Containsf(t, banner, "Kotlin version "+version,
		"kotlin -version must report %s\n--- output ---\n%s", version, banner)
	compilerBanner := common.StripANSI(execOK(ctx, t, ctr, activated("kotlinc -version 2>&1")))
	require.Containsf(t, compilerBanner, version,
		"kotlinc -version must mention %s\n--- output ---\n%s", version, compilerBanner)

	current := common.StripANSI(execOK(ctx, t, ctr, activated("vfox current kotlin")))
	require.Contains(t, current, version, "vfox current must name the selection\n--- output ---\n"+current)
	t.Logf("PASS: kotlin %s runs from %s", version, binDir)
}

// checkMissingVersionRejected asserts PreInstall's HEAD check rejects a
// release JetBrains never published.
func checkMissingVersionRejected(ctx context.Context, t *testing.T, ctr tc.Container) {
	t.Helper()
	code, out := execRun(ctx, t, ctr, "vfox install kotlin@0.0.0-e2e-missing")
	require.NotZerof(t, code, "installing a nonexistent release must fail\n--- output ---\n%s", out)
	t.Logf("PASS: a nonexistent release is rejected")
}
