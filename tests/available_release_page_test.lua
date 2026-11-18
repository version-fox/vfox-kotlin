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

-- An excerpt of the release page with its id and class attributes removed. The first table lists
-- release lines such as "2.4", which are not releases. The second table lists the releases.
local page = [[
<html><body>
<table><thead><tr><th><p>Kotlin release line</p></th><th><p>Release date</p></th><th><p>End of support</p></th><th><p>Latest release</p></th></tr></thead><tbody><tr><td><p><b>2.4</b></p></td><td><p>June 3, 2026</p></td><td><p>December 3, 2027</p></td><td><p>2.4.21</p></td></tr></tbody></table>
<table><thead><tr><th><p>Build info</p></th><th><p>Build highlights</p></th></tr></thead><tbody>
<tr><td><p><b>2.4.21</b></p><p>Released: <b>October 8, 2026</b></p></td><td><p>A bug fix release for Kotlin 2.4.21.</p></td></tr>
<tr><td><p><b>2.4.20</b></p><p>Released: <b>September 7, 2026</b></p></td><td><p>A tooling release.</p></td></tr>
<tr><td><p><b>2.4.0</b></p><p>Released: <b>June 3, 2026</b></p></td><td><p>A language release.</p></td></tr>
<tr><td><p><b>2.3.21</b></p><p>Released: <b>April 23, 2026</b></p></td><td><p>A bug fix release for Kotlin 2.3.20.</p></td></tr>
</tbody></table>
</body></html>
]]

local archiveUrl = "https://github.com/JetBrains/kotlin/releases/download/v2.4.21/kotlin-compiler-2.4.21.zip"
local checksum = "0123456789abcdef"

local plugin = test.load({
    http = function(request)
        if request.url:find("releases.html", 1, true) then
            return { status_code = 200, body = page }
        end
        if request.method == "HEAD" and request.url == archiveUrl then
            return { status_code = 200, body = "" }
        end
        if request.url == archiveUrl .. ".sha256" then
            return { status_code = 200, body = checksum }
        end
        return { status_code = 404, body = "Not Found" }
    end,
})

local versions = {}
for _, item in ipairs(plugin:Available({})) do
    table.insert(versions, item.version)
end
local expected = { "2.4.21", "2.4.20", "2.4.0", "2.3.21" }
assert(#versions == #expected, "expected " .. #expected .. " releases, got " .. #versions)
for i, version in ipairs(expected) do
    assert(versions[i] == version, "expected " .. version .. " at position " .. i .. ", got " .. versions[i])
end

-- `latest` takes the first listed version. It must be the newest release, not the release line 2.4.
local archive = plugin:PreInstall({ version = "latest" })
assert(archive.version == "2.4.21", "latest must resolve to 2.4.21, got " .. archive.version)
assert(archive.url == archiveUrl, "latest must use the 2.4.21 archive, got " .. archive.url)
assert(archive.sha256 == checksum, "latest must use the 2.4.21 checksum")
