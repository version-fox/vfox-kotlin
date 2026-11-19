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

local util = {}

util.DownloadURL = "https://github.com/JetBrains/kotlin/releases/download/v%s/kotlin-compiler-%s.zip"

-- Kotlin/Native is published as a separate archive next to the compiler.
-- JetBrains only builds it for these host platforms; anything else (Linux
-- arm64 for example) falls back to the compiler-only install.
local NATIVE_PLATFORMS = {
    linux = {
        amd64 = "linux-x86_64",
    },
    darwin = {
        amd64 = "macos-x86_64",
        arm64 = "macos-aarch64",
    },
    windows = {
        amd64 = "windows-x86_64",
    },
}

--- Return the Kotlin/Native platform token and archive extension for a
--- runtime's OS and architecture, or nil when no prebuilt compiler exists.
--- @param osType string RUNTIME.osType (linux, darwin, windows)
--- @param archType string RUNTIME.archType (amd64, arm64)
--- @return string|nil platform
--- @return string|nil extension
function util:NativePlatform(osType, archType)
    local byArch = NATIVE_PLATFORMS[osType]
    if byArch == nil then
        return nil, nil
    end
    local platform = byArch[archType]
    if platform == nil then
        return nil, nil
    end
    local extension = "tar.gz"
    if osType == "windows" then
        extension = "zip"
    end
    return platform, extension
end

--- URL of the Kotlin/Native prebuilt archive for a release. JetBrains renamed
--- the asset from kotlin-native-<platform>-<version> to
--- kotlin-native-prebuilt-<platform>-<version> in 2.0.0-Beta3, so Beta1 and
--- Beta2 keep the old name. The canonical name is the one that carries a
--- .sha256 file.
--- @param version string Kotlin release
--- @param osType string RUNTIME.osType
--- @param archType string RUNTIME.archType
--- @return string|nil url
function util:NativeDownloadURL(version, osType, archType)
    local platform, extension = self:NativePlatform(osType, archType)
    if platform == nil then
        return nil
    end
    local prefix = "kotlin-native-"
    if self:at_least(version, "2.0.0-Beta3") then
        prefix = "kotlin-native-prebuilt-"
    end
    return string.format(
        "https://github.com/JetBrains/kotlin/releases/download/v%s/%s%s-%s.%s",
        version,
        prefix,
        platform,
        version,
        extension
    )
end

--- Whether a Kotlin version is the minimum or later, pre-releases included:
--- 2.0.0-RC2 is at least 2.0.0-Beta3, and 1.9.0-Beta is not at least 1.9.0-RC.
--- @param version string Kotlin release
--- @param minimum string Lowest release that counts
--- @return boolean
function util:at_least(version, minimum)
    return version == minimum or self:compare_versions({ version = version }, { version = minimum })
end

--- Split a Kotlin version into its numeric parts and its pre-release tag:
--- "2.0.0-RC2" gives {2, 0, 0} and "RC2", and "2.0.0" gives {2, 0, 0} and nil.
--- @param version string Kotlin release
--- @return table parts
--- @return string|nil pre
local function splitVersion(version)
    local core, pre = string.match(version, "^([^-]+)%-(.+)$")
    if core == nil then
        core = version
    end
    local parts = {}
    for part in string.gmatch(core, "[^.]+") do
        table.insert(parts, tonumber(part))
    end
    return parts, pre
end

--- Order two pre-release tags such as Beta4 and RC2 by their letters, then by
--- the number after them: Beta4 < Beta5 < RC < RC2 < RC10.
--- @param a string
--- @param b string
--- @return integer -1, 0 or 1
local function comparePreRelease(a, b)
    local aLetters, aNumber = string.match(a, "^(%D*)(%d*)")
    local bLetters, bNumber = string.match(b, "^(%D*)(%d*)")
    if aLetters ~= bLetters then
        return aLetters < bLetters and -1 or 1
    end
    local aN, bN = tonumber(aNumber) or 0, tonumber(bNumber) or 0
    if aN == bN then
        return 0
    end
    return aN < bN and -1 or 1
end

--- Whether v1 comes after v2. The numeric parts decide first. A pre-release
--- comes before its release, so 2.0.0-RC2 is after 2.0.0-Beta5 and before 2.0.0.
--- @param v1o table {version = string}
--- @param v2o table {version = string}
--- @return boolean
function util:compare_versions(v1o, v2o)
    local v1_parts, v1_pre = splitVersion(v1o.version)
    local v2_parts, v2_pre = splitVersion(v2o.version)

    for i = 1, math.max(#v1_parts, #v2_parts) do
        local v1_part = v1_parts[i] or 0
        local v2_part = v2_parts[i] or 0
        if v1_part > v2_part then
            return true
        elseif v1_part < v2_part then
            return false
        end
    end

    if v1_pre == nil or v2_pre == nil then
        return v1_pre == nil and v2_pre ~= nil
    end
    return comparePreRelease(v1_pre, v2_pre) > 0
end

return util
