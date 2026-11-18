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

-- `latest` resolves through Available. When the releases page cannot be read, or lists no release,
-- PreInstall must stop with a message that says so, not index a list that is not there.
local mode = "unavailable"
local plugin = test.load({
    http = function(request)
        if mode == "unavailable" then
            return { status_code = 503, body = "unavailable" }
        end
        -- A page whose release table has no rows.
        return { status_code = 200, body = "<html><body><table><tbody></tbody></table></body></html>" }
    end,
})

local ok, err = pcall(plugin.PreInstall, plugin, { version = "latest" })
assert(not ok, "latest must fail when the releases page is unavailable")
assert(tostring(err):find("get Kotlin releases failed: HTTP 503", 1, true), "unexpected error: " .. tostring(err))
assert(not tostring(err):find("attempt to index", 1, true), "latest must not index a missing list: " .. tostring(err))

mode = "empty"
ok, err = pcall(plugin.PreInstall, plugin, { version = "latest" })
assert(not ok, "latest must fail when no release is listed")
assert(tostring(err):find("No Kotlin release is listed", 1, true), "unexpected error: " .. tostring(err))
assert(not tostring(err):find("attempt to index", 1, true), "latest must not index an empty list: " .. tostring(err))

print("PASS: latest reports a releases page that cannot be read or lists no release")
