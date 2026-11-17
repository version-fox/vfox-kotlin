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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, DefaultVersion, Version())
	t.Setenv("KOTLIN_VERSION", "1.9.24")
	assert.Equal(t, "1.9.24", Version())
	t.Setenv("KOTLIN_VERSION", "   ")
	assert.Equal(t, DefaultVersion, Version())
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.9.0", "1.9.0", 0},
		{"1.9.0", "1.8.22", 1},
		{"1.8.22", "1.9.0", -1},
		{"2.0.0", "1.99.99", 1},
		{"1.9", "1.9.0", 0},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, CompareVersions(c.a, c.b), "%s vs %s", c.a, c.b)
	}
}

func TestHasChecksum(t *testing.T) {
	assert.False(t, HasChecksum(LegacyVersion))
	assert.True(t, HasChecksum("1.9.0"))
	assert.True(t, HasChecksum(DefaultVersion))
}

func TestNormalizeArch(t *testing.T) {
	for _, s := range []string{"amd64", "x86_64", "AMD64", "x64"} {
		got, err := NormalizeArch(s)
		require.NoError(t, err)
		assert.Equal(t, "amd64", got)
	}
	for _, s := range []string{"arm64", "aarch64"} {
		got, err := NormalizeArch(s)
		require.NoError(t, err)
		assert.Equal(t, "arm64", got)
	}
	_, err := NormalizeArch("mips")
	assert.Error(t, err)
}

func TestInstalledVersion(t *testing.T) {
	out := "\x1b[32mInstall kotlin@2.0.0 success!\x1b[0m\nother text"
	assert.Equal(t, "2.0.0", InstalledVersion(out))
	assert.Equal(t, "", InstalledVersion("nothing here"))
}

func TestSanitizeAndContainerName(t *testing.T) {
	assert.Equal(t, "kotlin_2_0_0", Sanitize("kotlin/2.0.0"))
	first := ContainerName("v2_0_0")
	second := ContainerName("v2_0_0")
	assert.Contains(t, first, "vfox-kotlin-e2e-v2_0_0-")
	assert.NotEqual(t, first, second)
}

func TestMaxJobsAndAttempts(t *testing.T) {
	assert.Equal(t, 3, MaxJobs(3))
	t.Setenv("E2E_MAX_JOBS", "1")
	assert.Equal(t, 1, MaxJobs(3))
	t.Setenv("E2E_MAX_JOBS", "bogus")
	assert.Equal(t, 3, MaxJobs(3))

	assert.Equal(t, 3, RetryAttempts())
	t.Setenv("E2E_RETRY_ATTEMPTS", "5")
	assert.Equal(t, 5, RetryAttempts())
}
