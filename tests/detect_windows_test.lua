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
local plugin = test.load({
    os = "windows",
    arch = "amd64",
    env = {
        USERPROFILE = "C:\\Users\\John Doe",
        ProgramFiles = "C:\\Program Files",
    },
})
local detect = require("detect_kotlin")

assert(RUNTIME.osType == "windows")
assert(detect.join("C:\\Kotlin", "bin") == "C:\\Kotlin\\bin")

local files = {}
local dirs = {}

local function norm(path)
    return (path or ""):gsub("\\", "/")
end

io.open = function(path, mode)
    if (mode or "r") ~= "r" then
        return nil
    end
    local content = files[norm(path)]
    if content == nil then
        return nil
    end
    return {
        read = function(_, pattern)
            return pattern == "*a" and content or nil
        end,
        close = function()
            return true
        end,
    }
end

io.popen = function(cmd)
    local dir = cmd:match('"([^"]+)"') or cmd:match("'([^']+)'")
    local children = (dir and dirs[norm(dir)]) or {}
    local index = 0
    return {
        lines = function()
            return function()
                index = index + 1
                return children[index]
            end
        end,
        close = function()
            return true
        end,
    }
end

local function add_compiler(home, version)
    files[norm(home .. "\\build.txt")] = version .. "-release-1"
    files[norm(home .. "\\bin\\kotlinc.bat")] = "@echo off\n"
    local parent, name = home:match("^(.*)\\([^\\]+)$")
    if parent then
        local key = norm(parent)
        dirs[key] = dirs[key] or {}
        table.insert(dirs[key], name)
    end
end

-- A home with spaces in its path is a normal case on Windows.
local sdkman = "C:\\Users\\John Doe\\.sdkman\\candidates\\kotlin\\2.0.0"
local system = "C:\\Program Files\\Kotlin\\1.9.24"
add_compiler(sdkman, "2.0.0")
add_compiler(system, "1.9.24")

assert(detect.kotlin_home(sdkman), "a .bat launcher must be accepted on Windows")
local found = plugin:Detect({})
local versions = {}
for _, entry in ipairs(found) do
    versions[entry.version] = entry.path
end
assert(versions["2.0.0"] == sdkman, "SDKMAN compiler under a spaced path was not found")
assert(versions["1.9.24"] == system, "system compiler in Program Files was not found")

print("PASS: Detect finds compilers on Windows, including spaced paths")
