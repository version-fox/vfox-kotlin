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
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

func installKotlin(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox, version string) string {
	t.Helper()
	script := fmt.Sprintf("vfox install kotlin@%s", version)
	out := execRetry(ctx, t, ctr, script, script, box.vars...)
	installed := common.CheckInstalled(t, out, version)
	plain := common.StripANSI(out)
	if common.HasChecksum(version) {
		require.Regexp(t, `Verifying checksum [0-9a-fA-F]{64}\.\.\.`, plain,
			"vfox must verify the .sha256 asset for kotlin@%s\n--- output ---\n%s", version, plain)
	} else {
		require.Contains(t, plain, "Checksum is not provided",
			"older releases install without a checksum\n--- output ---\n%s", plain)
	}
	t.Logf("PASS: kotlin@%s installed as %s", version, installed)
	return installed
}
