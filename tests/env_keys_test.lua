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

-- An installed Kotlin/Native addition contributes its own bin directory, so
-- kotlinc-native is on PATH next to kotlinc.
local withNative = plugin:EnvKeys({
    path = "/sdks/kotlin/2.0.0",
    sdkInfo = {
        kotlin = { name = "kotlin", version = "2.0.0", path = "/sdks/kotlin/2.0.0" },
        ["add-native"] = { name = "add-native", version = "2.0.0", path = "/sdks/kotlin/add-native-2.0.0" },
    },
})
assert(#withNative == 2, "EnvKeys must return a second PATH entry for Kotlin/Native")
assert(withNative[2].key == "PATH", "the native entry must also extend PATH")
assert(withNative[2].value == "/sdks/kotlin/add-native-2.0.0/bin", "EnvKeys must add the native bin directory to PATH")

-- Linked installs and releases without a native build have no such addition.
local withoutNative = plugin:EnvKeys({
    path = "/sdks/kotlin/2.0.0",
    sdkInfo = { kotlin = { name = "kotlin", version = "2.0.0", path = "/sdks/kotlin/2.0.0" } },
})
assert(#withoutNative == 1, "EnvKeys must not add a PATH entry without a native SDK")

print("PASS: EnvKeys adds the SDK bin (and native bin) to PATH")
