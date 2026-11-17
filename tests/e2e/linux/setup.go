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

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// setupVfox installs the pinned vfox release from its GitHub archive, then
// adds the local plugin from a zip assembled from this checkout.
func setupVfox(ctx context.Context, t *testing.T, ctr tc.Container, arch string) {
	t.Helper()
	work := strings.TrimSpace(execOK(ctx, t, ctr, "mktemp --directory"))

	archSuffix, checksum := map[string]string{
		"amd64": "x86_64", "arm64": "aarch64",
	}[arch], ""
	if arch == "amd64" {
		checksum = common.VfoxLinuxAMD64SHA256
	} else if arch == "arm64" {
		checksum = common.VfoxLinuxARM64SHA256
	}
	if archSuffix == "" {
		t.Fatalf("unsupported arch %q", arch)
	}
	asset := fmt.Sprintf("vfox_%s_linux_%s", common.VfoxVersion, archSuffix)
	url := fmt.Sprintf("https://github.com/version-fox/vfox/releases/download/%s/%s.tar.gz",
		common.VfoxTag, asset)

	execRetry(ctx, t, ctr, "download vfox", fmt.Sprintf(
		"set -o pipefail; curl --fail --silent --show-error --location --retry 3 --retry-delay 3 %s -o %s/vfox.tar.gz && "+
			"echo %s > %s/vfox.sha256 && sed -i 's|$|  %s/vfox.tar.gz|' %s/vfox.sha256 && "+
			"sha256sum --check %s/vfox.sha256 && "+
			"tar --extract --gzip --file %s/vfox.tar.gz --directory %s --strip-components=1 && "+
			"install -m 0755 %s/vfox /usr/local/bin/vfox && vfox --version",
		shQuote(url), work, checksum, work, work, work, work, work, work, work))

	execOK(ctx, t, ctr, fmt.Sprintf(
		"(cd /e2e && zip -qr %s/kotlin.zip metadata.lua hooks lib Injection.lua) && "+
			"vfox add kotlin --source %s/kotlin.zip && vfox --version && vfox info kotlin",
		work, work))
	t.Logf("PASS: vfox installed and local kotlin plugin added")
}
