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

local page = [[
<html><body><table><tbody>
<tr><td><b>2.1.0</b></td><td>release</td></tr>
<tr><td><b>2.0.0</b></td><td>release</td></tr>
<tr><td><b>1.9.24</b></td><td>release</td></tr>
</tbody></table></body></html>
]]

local requests = {}
local failNext = false
local networkDown = false
local plugin = test.load({
    http = function(request)
        table.insert(requests, request.url)
        if request.url:find("kotlinlang.org", 1, true) then
            if networkDown then
                return nil, "dial tcp: lookup kotlinlang.org: no such host"
            end
            if failNext then
                return { status_code = 503, body = "unavailable" }
            end
            return { status_code = 200, body = page }
        end
        return { status_code = 404, body = "Not Found" }
    end,
})

local versions = {}
for _, item in ipairs(plugin:Available({})) do
    table.insert(versions, item.version)
    assert(item.note == "", "Available must not set a note")
end
assert(#versions == 3, "expected 3 versions, got " .. #versions)
assert(versions[1] == "2.1.0", "the first listed version must be 2.1.0")
assert(versions[2] == "2.0.0", "the second listed version must be 2.0.0")
assert(#requests == 1 and requests[1]:find("releases.html", 1, true), "Available must request the release details page")

-- A failed request raises an error that names the request and its status. It must not look like an empty list.
failNext = true
local ok, err = pcall(plugin.Available, plugin, {})
assert(not ok, "a failed request must raise an error")
assert(tostring(err):find("get Kotlin releases failed: HTTP 503", 1, true), "unexpected error: " .. tostring(err))
failNext = false

-- A request that never gets an answer raises the connection error.
networkDown = true
ok, err = pcall(plugin.Available, plugin, {})
assert(not ok, "a connection failure must raise an error")
assert(tostring(err):find("get Kotlin releases failed: ", 1, true), "unexpected error: " .. tostring(err))
assert(tostring(err):find("no such host", 1, true), "the connection error must be included: " .. tostring(err))
networkDown = false

print("PASS: Available lists versions and raises when the releases page cannot be read")
