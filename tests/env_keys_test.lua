-- Copyright 2026 Han Li and contributors
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--   http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

local test = require("vfox.test")
local plugin = test.load()

local result = plugin:EnvKeys({ path = "/sdks/kotlin/2.0.0" })
assert(#result == 1, "EnvKeys must return exactly one variable")
assert(result[1].key == "PATH", "EnvKeys must set PATH")
assert(result[1].value == "/sdks/kotlin/2.0.0/bin", "EnvKeys must add the SDK's bin directory")

-- A path with spaces must be passed through unchanged; vfox quotes it.
local spaced = plugin:EnvKeys({ path = "/sdks/My Kotlin/2.0.0" })
assert(spaced[1].value == "/sdks/My Kotlin/2.0.0/bin", "EnvKeys must preserve spaces in the path")

print("PASS: EnvKeys adds the SDK bin to PATH")
