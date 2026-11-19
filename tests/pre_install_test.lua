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
    -- Only these two releases ship a Kotlin/Native archive in this fixture.
    if url:find("kotlin-native-", 1, true) then
        if not (url:find("2.0.0", 1, true) or url:find("1.9.0", 1, true)) then
            return { status_code = 404, body = "Not Found" }
        end
        if url:sub(-#".sha256") == ".sha256" then
            return { status_code = 200, body = sha256 }
        end
        return { status_code = 200, body = "" }
    end
    if url:sub(-#".sha256") == ".sha256" then
        -- 1.9.0-Beta is the one pre-release in this fixture with no checksum asset.
        if url:find("1.9.0-Beta", 1, true) then
            return { status_code = 404, body = "Not Found" }
        end
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

local function nativeUrlFor(version)
    local prefix = "kotlin-native-"
    if version == "2.0.0" then
        prefix = "kotlin-native-prebuilt-"
    end
    return string.format(
        "https://github.com/JetBrains/kotlin/releases/download/v%s/%slinux-x86_64-%s.tar.gz",
        version,
        prefix,
        version
    )
end

local function countCompilerChecksums()
    local count = 0
    for _, url in ipairs(requests) do
        if url:find("kotlin-compiler-", 1, true) and url:sub(-#".sha256") == ".sha256" then
            count = count + 1
        end
    end
    return count
end

local cases = 0
local function check(version, wantNative)
    requests = {}
    local result = plugin:PreInstall({ version = version })
    assert(result.version == version, version .. ": unexpected installed version " .. tostring(result.version))
    assert(result.url == urlFor(version), version .. ": unexpected URL " .. tostring(result.url))
    if version == "1.8.22" then
        assert(result.sha256 == nil or result.sha256 == "", version .. ": no checksum should be returned before 1.9.0")
        assert(countCompilerChecksums() == 0, "older releases need no compiler checksum request")
    else
        assert(result.sha256 == sha256, version .. ": expected the .sha256 body")
        assert(countCompilerChecksums() == 1, version .. ": the compiler checksum file must be fetched")
    end
    if wantNative then
        assert(result.addition ~= nil and #result.addition == 1, version .. ": expected a native addition")
        local native = result.addition[1]
        assert(native.name == "native", version .. ": unexpected addition name " .. tostring(native.name))
        assert(native.version == version, version .. ": unexpected addition version " .. tostring(native.version))
        assert(native.url == nativeUrlFor(version), version .. ": unexpected native URL " .. tostring(native.url))
        assert(native.sha256 == sha256, version .. ": expected the native .sha256 body")
    else
        assert(result.addition == nil or #result.addition == 0, version .. ": no native addition should be present")
    end
    cases = cases + 1
end

check("2.0.0", true)
check("1.9.0", true)
check("1.8.22", false)

-- Pre-releases at the boundaries JetBrains changed. 1.9.0-Beta has no .sha256,
-- so it must install without asking for one; 1.9.0-RC has one and is verified.
-- 2.0.0-Beta2 still uses the pre-2.0 native archive name, Beta3 the new one.
local function checkPreRelease(version, wantChecksum, wantNativeUrl)
    requests = {}
    local result = plugin:PreInstall({ version = version })
    assert(result.version == version, version .. ": unexpected installed version " .. tostring(result.version))
    assert(result.url == urlFor(version), version .. ": unexpected URL " .. tostring(result.url))
    if wantChecksum then
        assert(result.sha256 == sha256, version .. ": expected the .sha256 body")
        assert(countCompilerChecksums() == 1, version .. ": the compiler checksum file must be fetched")
    else
        assert(result.sha256 == nil, version .. ": no checksum should be returned")
        assert(countCompilerChecksums() == 0, version .. ": no compiler checksum request")
    end
    if wantNativeUrl ~= nil then
        assert(result.addition ~= nil and #result.addition == 1, version .. ": expected a native addition")
        local got = result.addition[1].url
        assert(got == wantNativeUrl, version .. ": unexpected native URL " .. tostring(got))
    end
    cases = cases + 1
end

checkPreRelease("1.9.0-Beta", false, nil)
checkPreRelease(
    "1.9.0-RC",
    true,
    "https://github.com/JetBrains/kotlin/releases/download/v1.9.0-RC/kotlin-native-linux-x86_64-1.9.0-RC.tar.gz"
)
checkPreRelease(
    "2.0.0-Beta2",
    true,
    "https://github.com/JetBrains/kotlin/releases/download/v2.0.0-Beta2/kotlin-native-linux-x86_64-2.0.0-Beta2.tar.gz"
)
checkPreRelease(
    "2.0.0-Beta3",
    true,
    "https://github.com/JetBrains/kotlin/releases/download/v2.0.0-Beta3/kotlin-native-prebuilt-linux-x86_64-2.0.0-Beta3.tar.gz"
)
checkPreRelease(
    "2.0.0-RC2",
    true,
    "https://github.com/JetBrains/kotlin/releases/download/v2.0.0-RC2/kotlin-native-prebuilt-linux-x86_64-2.0.0-RC2.tar.gz"
)

-- `latest` resolves through Available to the first listed release.
requests = {}
local latest = plugin:PreInstall({ version = "latest" })
assert(latest.version == "2.0.0", "latest must resolve to 2.0.0, got " .. tostring(latest.version))
assert(latest.url == urlFor("2.0.0"), "latest must use the resolved release's URL")
assert(latest.sha256 == sha256, "latest must use the resolved release's checksum")
assert(latest.addition ~= nil and latest.addition[1].url == nativeUrlFor("2.0.0"), "latest must bundle Kotlin/Native")
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
