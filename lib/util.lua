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

--- URL of the Kotlin/Native prebuilt archive for a release. Kotlin 2.0
--- renamed the asset from kotlin-native-<platform>-<version> to
--- kotlin-native-prebuilt-<platform>-<version>; the canonical name is the
--- one that carries a .sha256 file.
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
    if self:compare_versions({ version = version }, { version = "2.0.0" }) or version == "2.0.0" then
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

function util:compare_versions(v1o, v2o)
    local v1 = v1o.version
    local v2 = v2o.version
    local v1_parts = {}
    for part in string.gmatch(v1, "[^.]+") do
        table.insert(v1_parts, tonumber(part))
    end

    local v2_parts = {}
    for part in string.gmatch(v2, "[^.]+") do
        table.insert(v2_parts, tonumber(part))
    end

    for i = 1, math.max(#v1_parts, #v2_parts) do
        local v1_part = v1_parts[i] or 0
        local v2_part = v2_parts[i] or 0
        if v1_part > v2_part then
            return true
        elseif v1_part < v2_part then
            return false
        end
    end

    return false
end

return util
