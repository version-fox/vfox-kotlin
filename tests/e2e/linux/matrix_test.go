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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShQuote(t *testing.T) {
	assert.Equal(t, "'plain'", shQuote("plain"))
	assert.Equal(t, `'it'\''s'`, shQuote("it's"))
}

func TestActivated(t *testing.T) {
	assert.Equal(t, `eval "$(vfox activate bash)"; vfox --version`, activated("vfox --version"))
}
