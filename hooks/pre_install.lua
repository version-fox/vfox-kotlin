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

--- Returns some pre-installed information, such as version number, download address, local files, etc.
--- If checksum is provided, vfox will automatically check it for you.
--- @param ctx table Context information; ctx.version is the requested version
--- @return table Version information
function PLUGIN:PreInstall(ctx)
    local version = ctx.version

    if version == "latest" then
        local lists = self:Available({})
        version = lists[1].version
    end
    local url = util.DownloadURL:format(version, version)
    local resp, err = http.head({
        url = url,
    })
    if err ~= nil or resp.status_code ~= 200 then
        error("Current version information not detected.")
    end
    if util:compare_versions({ version = version }, { version = "1.9.0" }) or version == "1.9.0" then
        resp, err = http.get({
            url = url .. ".sha256",
        })
        if err ~= nil or resp.status_code ~= 200 then
            error("Current version checksum not detected.")
        end
        return {
            version = version,
            url = url,
            sha256 = resp.body,
        }
    else
        return {
            version = version,
            url = url,
        }
    end
end
