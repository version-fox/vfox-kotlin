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

// Package common holds target-independent helpers for the Kotlin E2E suites.
package common

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// Plugin is metadata.lua's PLUGIN.name.
	Plugin = "kotlin"

	// DefaultVersion is a release JetBrains publishes a .sha256 for, so the
	// default run exercises PreInstall's checksum path.
	DefaultVersion = "2.0.0"
	// LegacyVersion predates the .sha256 assets; PreInstall installs it
	// without a checksum. Keeping it in the suite covers both branches.
	LegacyVersion = "1.8.22"

	// VfoxVersion and VfoxTag pin the release installed in the containers.
	// 1.1.0 is the release that added the Detect hook used by --detect.
	VfoxVersion = "1.1.0"
	VfoxTag     = "v" + VfoxVersion

	// Checksums of the pinned vfox archives (sha256sum of the release assets).
	VfoxLinuxAMD64SHA256   = "3b33c3c468b27bc17aecf0d6d3af24fe49fe204972984904a136a47b402a0ee4"
	VfoxLinuxARM64SHA256   = "e702e057dfe588200727dfb3cf786d181b13f015cbcc9e13679a3edcfd07e430"
	VfoxWindowsAMD64SHA256 = "1bbec8fffddbb62de0d93030c3f6203c69680662bc65257ca4a5e611cb607d0a"
)

// Version returns the Kotlin release to test, honouring $KOTLIN_VERSION.
func Version() string {
	if v := strings.TrimSpace(os.Getenv("KOTLIN_VERSION")); v != "" {
		return v
	}
	return DefaultVersion
}

// HasChecksum mirrors PreInstall's version threshold: releases from 1.9.0 on
// have a .sha256 asset on GitHub and are verified by vfox.
func HasChecksum(version string) bool {
	return CompareVersions(version, "1.9.0") >= 0
}

// HasNative reports whether JetBrains publishes a Kotlin/Native archive for
// the E2E target. Both suites run on amd64, where native builds exist from
// 1.9.0 on; Linux arm64 has no prebuilt native compiler.
func HasNative(version, arch string) bool {
	return arch == "amd64" && CompareVersions(version, "1.9.0") >= 0
}

// CompareVersions compares dot-separated numeric versions. Missing components
// count as zero, as in lib/util.lua.
func CompareVersions(a, b string) int {
	parse := func(v string) []int {
		var parts []int
		for _, s := range strings.Split(v, ".") {
			n, err := strconv.Atoi(s)
			if err != nil {
				n = 0
			}
			parts = append(parts, n)
		}
		return parts
	}
	x, y := parse(a), parse(b)
	for i := 0; i < len(x) || i < len(y); i++ {
		xi, yi := 0, 0
		if i < len(x) {
			xi = x[i]
		}
		if i < len(y) {
			yi = y[i]
		}
		if xi < yi {
			return -1
		}
		if xi > yi {
			return 1
		}
	}
	return 0
}

// MaxJobs caps parallel containers; $E2E_MAX_JOBS overrides the default.
func MaxJobs(def int) int {
	if v := strings.TrimSpace(os.Getenv("E2E_MAX_JOBS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			return n
		}
	}
	return def
}

// DefaultMaxJobs is how many cases run at once when E2E_MAX_JOBS is unset.
const DefaultMaxJobs = 3

// Case is one independent run of a suite, usually one Kotlin release in its
// own container. Run is the body of the subtest called Name.
type Case struct {
	Name string
	Run  func(t *testing.T)
}

// RunCases runs each case as a subtest and keeps at most MaxJobs of them
// running at once. A case holds its slot until its subtest finishes, so the
// containers the cases start never exceed E2E_MAX_JOBS. Cases start in the
// order given, and a failing case does not stop the others.
func RunCases(t *testing.T, cases []Case) {
	t.Helper()
	slots := make(chan struct{}, MaxJobs(DefaultMaxJobs))
	var wg sync.WaitGroup
	for _, c := range cases {
		slots <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			t.Run(c.Name, c.Run)
		}()
	}
	wg.Wait()
}

// RetryAttempts is the number of tries for network-touching steps.
func RetryAttempts() int {
	if v := strings.TrimSpace(os.Getenv("E2E_RETRY_ATTEMPTS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			return n
		}
	}
	return 3
}

// NormalizeArch maps Go and Windows architecture names to amd64/arm64.
func NormalizeArch(s string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "x86_64", "x64", "amd64", "amd64 ":
		return "amd64", nil
	case "aarch64", "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("unknown architecture %q", s)
	}
}

// DetectArch reads $PROCESSOR_ARCHITECTURE on Windows hosts, else GOARCH.
func DetectArch() (string, error) {
	if arch, ok := os.LookupEnv("PROCESSOR_ARCHITECTURE"); ok && strings.TrimSpace(arch) != "" {
		return NormalizeArch(arch)
	}
	return NormalizeArch(runtime.GOARCH)
}

// Sanitize turns a label into a Docker-name-safe string.
func Sanitize(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
			b.WriteByte(c)
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}

// ContainerName gives each combo a unique, greppable container name.
func ContainerName(slug string) string {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("vfox-kotlin-e2e-%s-%d", slug, time.Now().UnixNano())
	}
	return "vfox-kotlin-e2e-" + slug + "-" + hex.EncodeToString(b[:])
}

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// StripANSI removes the colour codes vfox writes to its progress and status
// lines so assertions can match plain text.
func StripANSI(s string) string { return ansiEscape.ReplaceAllString(s, "") }

var installSuccess = regexp.MustCompile(`Install kotlin@(\S+) success!`)

// InstalledVersion extracts the name vfox reports after an install.
func InstalledVersion(output string) string {
	if m := installSuccess.FindStringSubmatch(StripANSI(output)); m != nil {
		return m[1]
	}
	return ""
}

// CheckInstalled asserts the install succeeded under the requested name.
func CheckInstalled(t *testing.T, output, version string) string {
	t.Helper()
	installed := InstalledVersion(output)
	require.NotEmptyf(t, installed, "vfox install kotlin@%s did not report success\n--- output ---\n%s", version, output)
	require.Equalf(t, version, installed, "kotlin@%s installed as %q", version, installed)
	return installed
}

// Retry runs run until it exits 0 or the attempt budget is exhausted, backing
// off linearly between tries. Network calls in the suites all go through it.
func Retry(ctx context.Context, t *testing.T, label string, run func() (int, string)) string {
	t.Helper()
	max := RetryAttempts()
	if max < 1 {
		max = 1
	}
	var out string
	var code, attempts int
	for {
		attempts++
		code, out = run()
		if code == 0 {
			return out
		}
		if attempts >= max {
			break
		}
		t.Logf("retrying %s (attempt %d, exit %d)\n--- output ---\n%s", label, attempts, code, out)
		select {
		case <-ctx.Done():
			t.Fatalf("%s: %v", label, ctx.Err())
		case <-time.After(time.Duration(10*attempts) * time.Second):
		}
	}
	t.Fatalf("%s failed after %d attempt(s) (exit %d)\n--- output ---\n%s", label, attempts, code, out)
	return ""
}
