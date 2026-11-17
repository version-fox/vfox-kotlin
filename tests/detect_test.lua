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
local plugin = test.load({ os = "linux", arch = "amd64", env = { HOME = "/home/test" } })
local detect = require("detect_kotlin")

-- A minimal fake filesystem: files maps paths to contents, dirs maps a parent
-- directory to its immediate child names.
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
            if pattern == "*a" then
                return content
            end
            return nil
        end,
        close = function()
            return true
        end,
    }
end

io.popen = function(cmd)
    local dir = cmd:match("'([^']+)'") or cmd:match('"([^"]+)"')
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

local function add_child(parent, name)
    local key = norm(parent)
    dirs[key] = dirs[key] or {}
    for _, child in ipairs(dirs[key]) do
        if child == name then
            return
        end
    end
    table.insert(dirs[key], name)
end

local function add_compiler(home, version)
    files[norm(home .. "/build.txt")] = version .. "-release-358"
    files[norm(home .. "/bin/kotlinc")] = "#!/bin/sh\n"
    local parent, name = home:match("^(.*)/([^/]+)$")
    if parent then
        add_child(parent, name)
    end
end

local sdkman = "/home/test/.sdkman/candidates/kotlin/2.0.0"
local system = "/usr/share/kotlinc"
add_compiler(sdkman, "2.0.0")
add_compiler(system, "1.9.24")

-- build.txt parsing: release, EAP and malformed values.
assert(detect.version_from_build("1.9.0-release-358") == "1.9.0")
assert(detect.version_from_build("2.1.0-Beta2") == "2.1.0")
assert(detect.version_from_build("2.0.0") == "2.0.0")
assert(detect.version_from_build("garbage") == nil)
assert(detect.version_from_build(nil) == nil)
-- A directory name is the fallback when build.txt is missing.
assert(detect.version_for("/home/test/.sdkman/candidates/kotlin/1.9.24") == "1.9.24")

-- A directory is a compiler home only with both a launcher and the jar/build.txt.
assert(detect.kotlin_home(sdkman))
assert(not detect.kotlin_home("/home/test/.sdkman/candidates/kotlin"))
assert(not detect.kotlin_home("/usr/share"))

local found = plugin:Detect({})
local by_path = {}
for _, entry in ipairs(found) do
    by_path[entry.path] = entry
end
assert(by_path[sdkman] ~= nil, "SDKMAN compiler was not detected")
assert(by_path[sdkman].version == "2.0.0")
assert(by_path[sdkman].note == "installed by SDKMAN")
assert(by_path[system] ~= nil, "system compiler was not detected")
assert(by_path[system].version == "1.9.24")
assert(by_path[system].note == "system Kotlin compiler")
assert(found[1].version == "1.9.24", "results must be sorted by version")
assert(found[2].version == "2.0.0")

print("PASS: Detect finds and sorts compilers on Linux")
