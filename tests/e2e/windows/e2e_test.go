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

// Package windows is the Windows container E2E suite, driven from Go through
// testcontainers-go. See docs/e2e.md.
package windows

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// TestE2E builds the Windows target image once and runs one case per Kotlin
// release, each in a container of its own, as the Linux suite does. The
// release under test goes through the whole flow; the pre-1.9.0 release only
// checks its install. It needs a Windows Docker daemon in Windows-container
// mode; CI runs it on windows-latest. E2E_MAX_JOBS caps how many containers
// run at once.
func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skip container E2E in short mode")
	}
	// Ryuk has no Windows image, so the reaper can never start here.
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	arch, err := common.DetectArch()
	if err != nil {
		t.Fatal(err)
	}
	image := "vfox-kotlin-e2e:windows-" + arch
	// Pinned to the Server 2025 build used by the windows-latest runners;
	// process isolation needs the container's build to match the host.
	baseImage := "mcr.microsoft.com/windows/servercore:ltsc2025-KB5122871-amd64"
	platform := "windows/" + arch
	repoRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	provider, err := tc.NewDockerProvider(tc.WithDefaultBridgeNetwork("nat"))
	if err != nil {
		t.Fatalf("new docker provider: %v", err)
	}

	tag, err := provider.BuildImage(ctx, &tc.ContainerRequest{
		FromDockerfile: tc.FromDockerfile{
			Context:        repoRoot,
			Dockerfile:     "tests/e2e/windows/Dockerfile",
			Repo:           "vfox-kotlin-e2e",
			Tag:            "windows-" + arch,
			BuildArgs:      map[string]*string{"BASE_IMAGE": &baseImage},
			KeepImage:      true,
			BuildLogWriter: os.Stdout,
		},
	})
	if err != nil {
		t.Fatalf("build image: %v", err)
	}
	if tag != image {
		t.Fatalf("built tag %q, want %q", tag, image)
	}

	// start gives a case its own throwaway container, removed when the case ends.
	start := func(t *testing.T, version string) tc.Container {
		t.Helper()
		ctr, err := provider.RunContainer(ctx, tc.ContainerRequest{
			Image:         image,
			ImagePlatform: platform,
			Name:          common.ContainerName("windows-" + common.Sanitize(version)),
			Env:           map[string]string{"KOTLIN_VERSION": version},
		})
		if err != nil {
			t.Fatalf("start container: %v", err)
		}
		t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })
		return ctr
	}

	cases := []common.Case{{
		Name: common.Version(),
		Run: func(t *testing.T) {
			ctr := start(t, common.Version())
			ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
			defer cancel()

			t.Logf("=== windows/%s, kotlin@%s ===", arch, common.Version())
			box := setupSlotEnv(ctx, t, ctr, strings.ToLower(common.Sanitize(common.Version())))
			setupVfox(ctx, t, ctr, box, arch)
			activation := fetchActivation(ctx, t, ctr, box)
			installKotlin(ctx, t, ctr, box, common.Version())
			verifySearch(ctx, t, ctr, box, common.Version())
			checkMissingVersionRejected(ctx, t, ctr, box)
			useKotlin(ctx, t, ctr, box, activation, common.Version())
			verifyToolchain(ctx, t, ctr, box, activation, common.Version())
			if common.HasNative(common.Version(), arch) {
				verifyNative(ctx, t, ctr, box, activation, common.Version())
			}
			checkDetect(ctx, t, ctr, box, activation)
		},
	}}
	// The pre-1.9.0 release has no checksum asset, so PreInstall takes the
	// other branch. It runs in a container of its own unless it is already the
	// release under test.
	if common.Version() != common.LegacyVersion {
		cases = append(cases, common.Case{
			Name: common.LegacyVersion,
			Run: func(t *testing.T) {
				ctr := start(t, common.LegacyVersion)
				ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
				defer cancel()

				box := setupSlotEnv(ctx, t, ctr, strings.ToLower(common.Sanitize(common.LegacyVersion)))
				setupVfox(ctx, t, ctr, box, arch)
				installKotlin(ctx, t, ctr, box, common.LegacyVersion)
			},
		})
	}
	common.RunCases(t, cases)
}
