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
    -- A pre-release sorts below its release and above the earlier pre-releases.
    { "2.0.0-RC2", "2.0.0", false },
    { "2.0.0", "2.0.0-RC2", true },
    { "2.0.0-RC2", "2.0.0-Beta5", true },
    { "2.0.0-Beta4", "2.0.0-Beta3", true },
    { "2.0.0-Beta2", "2.0.0-Beta3", false },
    { "2.0.0-RC10", "2.0.0-RC2", true },
    { "2.0.0-Beta3", "2.0.0-Beta3", false },
    { "1.9.0-Beta", "1.9.0-RC", false },
    { "1.9.0-RC", "1.9.0", false },
    { "2.0.1-Beta1", "2.0.0", true },
}
for _, case in ipairs(cases) do
    local got = util:compare_versions({ version = case[1] }, { version = case[2] })
    assert(
        got == case[3],
        string.format("compare_versions(%s, %s) = %s, want %s", case[1], case[2], tostring(got), tostring(case[3]))
    )
end

print("PASS: " .. #cases .. " version comparison cases")

-- at_least takes the lowest release that counts, pre-releases included. The
-- Kotlin/Native names change at 2.0.0-Beta3 and checksums start at 1.9.0-RC.
local minimums = {
    { "2.0.0-Beta3", "2.0.0-Beta3", true },
    { "2.0.0-RC2", "2.0.0-Beta3", true },
    { "2.0.0-Beta2", "2.0.0-Beta3", false },
    { "1.9.0", "1.9.0-RC", true },
    { "1.9.0-RC", "1.9.0-RC", true },
    { "1.9.0-Beta", "1.9.0-RC", false },
    { "1.8.22", "1.9.0-RC", false },
}
for _, case in ipairs(minimums) do
    local got = util:at_least(case[1], case[2])
    assert(
        got == case[3],
        string.format("at_least(%s, %s) = %s, want %s", case[1], case[2], tostring(got), tostring(case[3]))
    )
end
print("PASS: " .. #minimums .. " minimum version cases")

-- Kotlin/Native is a separate asset. The platform token and the archive name
-- must match JetBrains' release naming, which changed in 2.0.
local nativePlatforms = {
    { "linux", "amd64", "linux-x86_64", "tar.gz" },
    { "darwin", "amd64", "macos-x86_64", "tar.gz" },
    { "darwin", "arm64", "macos-aarch64", "tar.gz" },
    { "windows", "amd64", "windows-x86_64", "zip" },
}
for _, case in ipairs(nativePlatforms) do
    local platform, extension = util:NativePlatform(case[1], case[2])
    assert(platform == case[3], case[1] .. "/" .. case[2] .. ": unexpected platform " .. tostring(platform))
    assert(extension == case[4], case[1] .. "/" .. case[2] .. ": unexpected extension " .. tostring(extension))
end
-- No prebuilt native compiler exists for Linux arm64.
assert(util:NativePlatform("linux", "arm64") == nil, "Linux arm64 must have no native platform")

local nativeUrls = {
    {
        "2.0.21",
        "linux",
        "amd64",
        "https://github.com/JetBrains/kotlin/releases/download/v2.0.21/kotlin-native-prebuilt-linux-x86_64-2.0.21.tar.gz",
    },
    {
        "2.0.21",
        "windows",
        "amd64",
        "https://github.com/JetBrains/kotlin/releases/download/v2.0.21/kotlin-native-prebuilt-windows-x86_64-2.0.21.zip",
    },
    {
        "1.9.24",
        "linux",
        "amd64",
        "https://github.com/JetBrains/kotlin/releases/download/v1.9.24/kotlin-native-linux-x86_64-1.9.24.tar.gz",
    },
    {
        "1.9.24",
        "darwin",
        "arm64",
        "https://github.com/JetBrains/kotlin/releases/download/v1.9.24/kotlin-native-macos-aarch64-1.9.24.tar.gz",
    },
    -- The archive name changed in 2.0.0-Beta3, so Beta2 keeps the pre-2.0 name.
    {
        "1.9.0-RC",
        "linux",
        "amd64",
        "https://github.com/JetBrains/kotlin/releases/download/v1.9.0-RC/kotlin-native-linux-x86_64-1.9.0-RC.tar.gz",
    },
    {
        "2.0.0-Beta2",
        "linux",
        "amd64",
        "https://github.com/JetBrains/kotlin/releases/download/v2.0.0-Beta2/kotlin-native-linux-x86_64-2.0.0-Beta2.tar.gz",
    },
    {
        "2.0.0-Beta3",
        "linux",
        "amd64",
        "https://github.com/JetBrains/kotlin/releases/download/v2.0.0-Beta3/kotlin-native-prebuilt-linux-x86_64-2.0.0-Beta3.tar.gz",
    },
    {
        "2.0.0-RC2",
        "windows",
        "amd64",
        "https://github.com/JetBrains/kotlin/releases/download/v2.0.0-RC2/kotlin-native-prebuilt-windows-x86_64-2.0.0-RC2.zip",
    },
}
for _, case in ipairs(nativeUrls) do
    local url = util:NativeDownloadURL(case[1], case[2], case[3])
    assert(url == case[4], case[1] .. "/" .. case[2] .. ": unexpected native URL " .. tostring(url))
end
assert(util:NativeDownloadURL("2.0.21", "linux", "arm64") == nil, "unsupported platforms must have no native URL")

print("PASS: " .. #nativeUrls .. " Kotlin/Native URL cases")
