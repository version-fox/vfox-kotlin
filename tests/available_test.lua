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
local plugin = test.load({
    http = function(request)
        table.insert(requests, request.url)
        if request.url:find("kotlinlang.org", 1, true) then
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

-- A failed or unparsable page yields an empty list, not an error.
failNext = true
assert(#plugin:Available({}) == 0, "a failed request must produce no versions")
failNext = false

print("PASS: Available lists versions from the release page")
