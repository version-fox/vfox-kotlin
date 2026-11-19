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
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// installKotlin installs a release and asserts the plugin's two checksum
// paths: releases from 1.9.0-RC on carry a .sha256 that vfox verifies, and
// the others, such as 1.9.0-Beta and 1.8.22, install with vfox's "checksum not
// provided" warning.
func installKotlin(ctx context.Context, t *testing.T, ctr tc.Container, version string) string {
	t.Helper()
	command := fmt.Sprintf("vfox install kotlin@%s", version)
	out := execRetry(ctx, t, ctr, command, command)
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
