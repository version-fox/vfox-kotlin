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

// checksumFrom is the first release that PreInstall returns a checksum for:
// JetBrains publishes a .sha256 asset from 1.9.0-RC on, and 1.9.0-Beta has none.
const checksumFrom = "1.9.0-RC"

// nativeFrom is the first release the suite checks kotlinc-native for. It is a
// pre-release, not 1.9.0: a pre-release sorts before its release, so 1.9.0 alone
// would leave out the 1.9.0 pre-releases, which have a native archive too.
const nativeFrom = "1.9.0-Beta"

// HasChecksum reports whether vfox verifies a checksum for version. It mirrors
// the threshold in hooks/pre_install.lua.
func HasChecksum(version string) bool {
	return atLeast(version, checksumFrom)
}

// HasNative reports whether the suite checks kotlinc-native for version on the
// E2E target, which is amd64 in both suites. Some older releases ship a native
// archive too; the suite does not check them.
func HasNative(version, arch string) bool {
	return arch == "amd64" && atLeast(version, nativeFrom)
}

// atLeast reports whether version is minimum or later, pre-releases included.
// It mirrors util:at_least in lib/util.lua.
func atLeast(version, minimum string) bool {
	return version == minimum || CompareVersions(version, minimum) > 0
}

// CompareVersions orders two Kotlin releases as lib/util.lua does: the numeric
// parts first, with missing parts counting as zero, then the pre-release tag.
// A release sorts after its pre-releases, so 1.9.0-RC comes before 1.9.0. It
// returns -1, 0 or 1.
func CompareVersions(a, b string) int {
	aParts, aPre, aHasPre := splitVersion(a)
	bParts, bPre, bHasPre := splitVersion(b)
	for i := range max(len(aParts), len(bParts)) {
		if c := cmpInt(partAt(aParts, i), partAt(bParts, i)); c != 0 {
			return c
		}
	}
	switch {
	case !aHasPre && !bHasPre:
		return 0
	case !aHasPre:
		return 1
	case !bHasPre:
		return -1
	default:
		return comparePreRelease(aPre, bPre)
	}
}

// splitVersion splits a release such as 2.0.0-RC2 into its numeric parts and
// its pre-release tag, as lib/util.lua does.
func splitVersion(version string) (parts []int, pre string, hasPre bool) {
	core := version
	if i := strings.Index(version, "-"); i > 0 && i < len(version)-1 {
		core, pre, hasPre = version[:i], version[i+1:], true
	}
	for _, s := range strings.Split(core, ".") {
		n, err := strconv.Atoi(s)
		if err != nil {
			n = 0
		}
		parts = append(parts, n)
	}
	return parts, pre, hasPre
}

// comparePreRelease orders pre-release tags such as Beta4 and RC2 by their
// letters, then by the number after them: Beta4 < Beta5 < RC < RC2 < RC10.
func comparePreRelease(a, b string) int {
	aLetters, aNumber := splitTag(a)
	bLetters, bNumber := splitTag(b)
	if aLetters != bLetters {
		if aLetters < bLetters {
			return -1
		}
		return 1
	}
	return cmpInt(aNumber, bNumber)
}

// splitTag splits a pre-release tag into the letters before its first digit
// and the number that follows them, which is zero when there is none.
func splitTag(tag string) (letters string, number int) {
	i := 0
	for i < len(tag) && (tag[i] < '0' || tag[i] > '9') {
		i++
	}
	j := i
	for j < len(tag) && tag[j] >= '0' && tag[j] <= '9' {
		j++
	}
	number, _ = strconv.Atoi(tag[i:j]) // an empty number gives zero
	return tag[:i], number
}

// partAt returns the i-th numeric part, or zero past the end of parts.
func partAt(parts []int, i int) int {
	if i < len(parts) {
		return parts[i]
	}
	return 0
}

// cmpInt returns -1, 0 or 1 as a is less than, equal to or greater than b.
func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
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
