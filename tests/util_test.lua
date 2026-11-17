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
test.load()

local util = require("util")

assert(
    util.DownloadURL == "https://github.com/JetBrains/kotlin/releases/download/v%s/kotlin-compiler-%s.zip",
    "unexpected download URL template"
)

local cases = {
    { "2.0.0", "1.9.0", true },
    { "1.9.0", "1.9.0", false },
    { "1.8.22", "1.9.0", false },
    { "1.9.0", "1.10.0", false },
    { "1.9", "1.9.0", false },
    { "2.0.0", "2.0.0", false },
}
for _, case in ipairs(cases) do
    local got = util:compare_versions({ version = case[1] }, { version = case[2] })
    assert(
        got == case[3],
        string.format("compare_versions(%s, %s) = %s, want %s", case[1], case[2], tostring(got), tostring(case[3]))
    )
end

print("PASS: " .. #cases .. " version comparison cases")
