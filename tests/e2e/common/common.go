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
	VfoxVersion = "1.0.13"
	VfoxTag     = "v" + VfoxVersion

	// Checksums of the pinned vfox archives (sha256sum of the release assets).
	VfoxLinuxAMD64SHA256   = "7ad0a4c82c8173f09393e09abecc3e27468253bab5fe53ffbf58f2b12b3fc768"
	VfoxLinuxARM64SHA256   = "7cec88e6e9be551cb138e061d37521286d34096fdc0c4fdb65d89cef30c58583"
	VfoxWindowsAMD64SHA256 = "dab9f12aa1139eb4a055afccbfa2b970b389e1fd3f14c94e9e6793ae5fc66e0a"
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
