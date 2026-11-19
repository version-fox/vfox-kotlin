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

local util = require("util")
local http = require("http")

--- Build the additional package holding Kotlin/Native, when JetBrains
--- publishes a prebuilt compiler for this host. vfox downloads and extracts
--- it into `add-native-<version>` beside the main SDK, and EnvKeys puts its
--- bin directory on PATH so `kotlinc-native` is available too. A missing
--- asset is not fatal: the compiler-only install still succeeds.
--- @param version string Kotlin release
--- @return table|nil addition
local function nativeAddition(version)
    local url = util:NativeDownloadURL(version, RUNTIME.osType, RUNTIME.archType)
    if url == nil then
        return nil
    end
    local resp, err = http.head({
        url = url,
    })
    if err ~= nil or resp.status_code ~= 200 then
        return nil
    end
    local addition = {
        name = "native",
        version = version,
        url = url,
    }
    -- Older releases have no native .sha256; verify when one is published.
    local checksum, checksumErr = http.get({
        url = url .. ".sha256",
    })
    if checksumErr == nil and checksum.status_code == 200 and checksum.body ~= "" then
        addition.sha256 = checksum.body
    end
    return addition
end

--- Returns some pre-installed information, such as version number, download address, local files, etc.
--- If checksum is provided, vfox will automatically check it for you.
--- @param ctx table Context information; ctx.version is the requested version
--- @return table Version information
function PLUGIN:PreInstall(ctx)
    local version = ctx.version

    if version == "latest" then
        local lists = self:Available({})
        if #lists == 0 then
            error("No Kotlin release is listed, so latest cannot be resolved.")
        end
        version = lists[1].version
    end
    local url = util.DownloadURL:format(version, version)
    local resp, err = http.head({
        url = url,
    })
    if err ~= nil or resp.status_code ~= 200 then
        error("Current version information not detected.")
    end
    local result = {
        version = version,
        url = url,
    }
    -- JetBrains publishes a .sha256 from 1.9.0-RC on; 1.9.0-Beta has none.
    if util:at_least(version, "1.9.0-RC") then
        resp, err = http.get({
            url = url .. ".sha256",
        })
        if err ~= nil or resp.status_code ~= 200 then
            error("Current version checksum not detected.")
        end
        result.sha256 = resp.body
    end
    local native = nativeAddition(version)
    if native ~= nil then
        result.addition = { native }
    end
    return result
end
