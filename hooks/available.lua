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

local html = require("html")
local http = require("http")
--- Return all available versions provided by this plugin.
--- Raises an error naming the failed request when the releases page cannot be read.
--- @param ctx table Empty table used as context, for future extension
--- @return table Descriptions of available versions and accompanying tool descriptions
function PLUGIN:Available(ctx)
    local resp, err = http.get({
        url = "https://kotlinlang.org/docs/releases.html#release-details",
    })
    if err ~= nil then
        error("get Kotlin releases failed: " .. tostring(err))
    end
    if resp.status_code ~= 200 then
        error("get Kotlin releases failed: HTTP " .. tostring(resp.status_code))
    end
    local result = {}
    html.parse(resp.body):find("tbody tr"):each(function(i, sel)
        local version = sel:find("td b"):first():text()
        -- The first table of the page lists release lines such as "2.4"; only X.Y.Z entries are releases.
        if version:match("^%d+%.%d+%.%d+$") then
            table.insert(result, {
                version = version,
                note = "",
            })
        end
    end)
    return result
end
