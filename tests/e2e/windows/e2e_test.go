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

// TestE2E builds the Windows target image once and drives one throwaway
// container through setup -> install -> verify. It needs a Windows Docker
// daemon in Windows-container mode; CI runs it on windows-latest.
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

	ctr, err := provider.RunContainer(ctx, tc.ContainerRequest{
		Image:         image,
		ImagePlatform: platform,
		Name:          common.ContainerName("windows-" + common.Sanitize(common.Version())),
		Env:           map[string]string{"KOTLIN_VERSION": common.Version()},
	})
	if err != nil {
		t.Fatalf("start container: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })

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

	installKotlin(ctx, t, ctr, box, common.LegacyVersion)
}
