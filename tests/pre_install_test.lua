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

local releases = [[
<html><body><table><tbody>
<tr><td><b>2.0.0</b></td><td>release</td></tr>
<tr><td><b>1.9.0</b></td><td>release</td></tr>
</tbody></table></body></html>
]]

local requests = {}
local sha256 = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

local function answer(request)
    table.insert(requests, request.url)
    local url = request.url
    if url:find("kotlinlang.org", 1, true) then
        return { status_code = 200, body = releases }
    end
    if url:sub(-#".sha256") == ".sha256" then
        if url:find("2.0.0", 1, true) or url:find("1.9.0", 1, true) then
            return { status_code = 200, body = sha256 }
        end
        return { status_code = 404, body = "Not Found" }
    end
    -- The HEAD of the compiler archive: everything but the deliberately
    -- missing release exists, so an absent .sha256 can be tested on its own.
    if url:find("0.0.0-e2e-missing", 1, true) then
        return { status_code = 404, body = "Not Found" }
    end
    return { status_code = 200, body = "" }
end

local plugin = test.load({ os = "linux", arch = "amd64", http = answer })

local function urlFor(version)
    return string.format(
        "https://github.com/JetBrains/kotlin/releases/download/v%s/kotlin-compiler-%s.zip",
        version,
        version
    )
end

local cases = 0
local function check(version)
    requests = {}
    local result = plugin:PreInstall({ version = version })
    assert(result.version == version, version .. ": unexpected installed version " .. tostring(result.version))
    assert(result.url == urlFor(version), version .. ": unexpected URL " .. tostring(result.url))
    if version == "1.8.22" then
        assert(result.sha256 == nil or result.sha256 == "", version .. ": no checksum should be returned before 1.9.0")
        assert(#requests == 1, "older releases need only the HEAD request")
    else
        assert(result.sha256 == sha256, version .. ": expected the .sha256 body")
        assert(
            #requests == 2 and requests[2] == urlFor(version) .. ".sha256",
            version .. ": the checksum file must be fetched"
        )
    end
    cases = cases + 1
end

check("2.0.0")
check("1.9.0")
check("1.8.22")

-- `latest` resolves through Available to the first listed release.
requests = {}
local latest = plugin:PreInstall({ version = "latest" })
assert(latest.version == "2.0.0", "latest must resolve to 2.0.0, got " .. tostring(latest.version))
assert(latest.url == urlFor("2.0.0"), "latest must use the resolved release's URL")
assert(latest.sha256 == sha256, "latest must use the resolved release's checksum")
cases = cases + 1

-- A version JetBrains never published fails in PreInstall's HEAD check.
local ok, err = pcall(plugin.PreInstall, plugin, { version = "0.0.0-e2e-missing" })
assert(not ok, "a missing release must raise an error")
assert(tostring(err):find("not detected", 1, true), "unexpected error: " .. tostring(err))
cases = cases + 1

-- A release with a checksum asset that does not answer fails too.
local ok2, err2 = pcall(plugin.PreInstall, plugin, { version = "1.9.1" })
assert(not ok2, "a missing checksum must raise an error")
assert(tostring(err2):find("checksum", 1, true), "unexpected error: " .. tostring(err2))
cases = cases + 1

print("PASS: " .. cases .. " PreInstall cases")
