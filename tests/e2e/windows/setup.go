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

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// slotBox carries the per-combo work directory and the env overrides that
// apply it. Each exec starts a fresh process, so the assignments ride along
// as tcexec.WithEnv entries rather than being set once.
type slotBox struct {
	workDir   string
	vfoxExe   string
	pluginZip string
	vars      []string
}

const curlRetry = "--retry 3 --retry-delay 5 --retry-all-errors"

// setupSlotEnv creates an isolated work directory for the combo and points
// VFOX_HOME at a path containing a space (C:\Users\John Doe\.vfox), the case
// vfox has supported since 1.0.13.
func setupSlotEnv(ctx context.Context, t *testing.T, ctr tc.Container, slot string) slotBox {
	t.Helper()
	root := strings.TrimSpace(execOK(ctx, t, ctr, "$env:USERPROFILE"))
	path := strings.TrimSpace(execOK(ctx, t, ctr, "$env:PATH"))
	box := slotBox{
		workDir:   root + `\vfox-e2e-runs\` + slot + `\work`,
		vfoxExe:   root + `\vfox-e2e-runs\` + slot + `\work\vfox.exe`,
		pluginZip: root + `\vfox-e2e-runs\` + slot + `\work\kotlin.zip`,
	}
	box.vars = []string{
		"USERPROFILE=" + root + `\vfox-e2e-runs\` + slot,
		"VFOX_HOME=" + root + `\John Doe\.vfox`,
		"TEMP=" + box.workDir,
		"TMP=" + box.workDir,
		"PATH=" + box.workDir + ";" + path,
	}
	execOK(ctx, t, ctr, fmt.Sprintf("New-Item -ItemType Directory -Force -Path %s | Out-Null",
		pwshQuote(box.workDir)), box.vars...)
	return box
}

// setupVfox installs the pinned vfox release and adds the local plugin.
func setupVfox(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox, arch string) {
	t.Helper()
	if arch != "amd64" {
		t.Fatalf("the Windows E2E image currently supports amd64 only (got %s)", arch)
	}
	zipPath := box.workDir + `\vfox.zip`
	url := fmt.Sprintf("https://github.com/version-fox/vfox/releases/download/%s/vfox_%s_windows_x86_64.zip",
		common.VfoxTag, common.VfoxVersion)
	sha := common.VfoxWindowsAMD64SHA256

	execRetry(ctx, t, ctr, "download vfox", fmt.Sprintf(
		"curl.exe %s -fsSL -o %s %s; "+
			"$actual = (Get-FileHash -Algorithm SHA256 -Path %s).Hash.ToLowerInvariant(); "+
			"if ($actual -ne %s) { throw ('vfox checksum mismatch: ' + $actual) }; "+
			"tar -xf %s -C %s; Remove-Item %s; "+
			"Copy-Item -Path (Get-ChildItem -Path %s -Recurse -Filter 'vfox.exe' | Select-Object -First 1) -Destination %s; "+
			"& %s --version",
		curlRetry, pwshQuote(zipPath), pwshQuote(url), pwshQuote(zipPath),
		pwshQuote(sha), pwshQuote(zipPath), pwshQuote(box.workDir), pwshQuote(zipPath),
		pwshQuote(box.workDir), pwshQuote(box.vfoxExe), pwshQuote(box.vfoxExe)), box.vars...)

	execOK(ctx, t, ctr, fmt.Sprintf(
		"tar -a -cf %s -C 'C:\\e2e' metadata.lua hooks lib Injection.lua; "+
			"& %s add kotlin --source %s; & %s info kotlin",
		pwshQuote(box.pluginZip), pwshQuote(box.vfoxExe), pwshQuote(box.pluginZip),
		pwshQuote(box.vfoxExe)), box.vars...)
	t.Logf("PASS: vfox installed and local kotlin plugin added")
}

// fetchActivation captures `vfox activate pwsh` once per container.
func fetchActivation(ctx context.Context, t *testing.T, ctr tc.Container, box slotBox) string {
	t.Helper()
	out := execOK(ctx, t, ctr, "& "+pwshQuote(box.vfoxExe)+" activate pwsh", box.vars...)
	if strings.TrimSpace(out) == "" {
		t.Fatal("vfox activate pwsh produced no output")
	}
	return out
}
