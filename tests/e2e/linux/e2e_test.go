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

// Package linux is the Linux container E2E suite, driven from Go through
// testcontainers-go. See docs/e2e.md.
package linux

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/version-fox/vfox-kotlin/tests/e2e/common"
)

// TestE2E builds the target image once and runs one case per Kotlin release,
// each in a container of its own. The release under test goes through setup
// -> install -> verify. The pre-1.9.0 release only checks its install, which
// takes the no-checksum branch. E2E_MAX_JOBS caps how many containers run at
// once. Scope it while iterating with KOTLIN_VERSION=1.8.22 go test -v -run
// TestE2E ./linux.
func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skip container E2E in short mode")
	}

	arch, err := common.DetectArch()
	if err != nil {
		t.Fatal(err)
	}
	image := "vfox-kotlin-e2e:linux-" + arch
	repoRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	provider, err := tc.NewDockerProvider()
	if err != nil {
		t.Fatalf("new docker provider: %v", err)
	}
	t.Cleanup(func() { _ = provider.Close() })

	tag, err := provider.BuildImage(ctx, &tc.ContainerRequest{
		FromDockerfile: tc.FromDockerfile{
			Context:        repoRoot,
			Dockerfile:     "tests/e2e/linux/Dockerfile",
			Repo:           "vfox-kotlin-e2e",
			Tag:            "linux-" + arch,
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
		ctr, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Image: image,
				Name:  common.ContainerName("linux-" + common.Sanitize(version)),
				Env:   map[string]string{"KOTLIN_VERSION": version},
			},
			Started: true,
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
			ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
			defer cancel()

			t.Logf("=== linux/%s, kotlin@%s ===", arch, common.Version())
			setupVfox(ctx, t, ctr, arch)
			installKotlin(ctx, t, ctr, common.Version())
			verifySearch(ctx, t, ctr, common.Version())
			checkMissingVersionRejected(ctx, t, ctr)
			useKotlin(ctx, t, ctr, common.Version())
			verifyToolchain(ctx, t, ctr, common.Version())
			if common.HasNative(common.Version(), arch) {
				verifyNative(ctx, t, ctr, common.Version())
			}
			checkDetect(ctx, t, ctr)
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
				ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
				defer cancel()

				setupVfox(ctx, t, ctr, arch)
				installKotlin(ctx, t, ctr, common.LegacyVersion)
			},
		})
	}
	common.RunCases(t, cases)
}
