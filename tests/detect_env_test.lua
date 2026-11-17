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
    os = "linux",
    arch = "amd64",
    env = {
        HOME = "/home/u",
        KOTLIN_HOME = "/opt/kotlin-language",
    },
})

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

local function add_compiler(home, version)
    files[norm(home .. "/build.txt")] = version .. "-release-1"
    files[norm(home .. "/bin/kotlinc")] = "#!/bin/sh\n"
    local parent, name = home:match("^(.*)/([^/]+)$")
    if parent then
        local key = norm(parent)
        dirs[key] = dirs[key] or {}
        table.insert(dirs[key], name)
    end
end

-- Each version-manager root and KOTLIN_HOME contributes a compiler.
add_compiler("/home/u/.sdkman/candidates/kotlin/2.0.0", "2.0.0")
add_compiler("/home/u/.asdf/installs/kotlin/1.9.24", "1.9.24")
add_compiler("/home/u/.local/share/mise/installs/kotlin/2.1.0", "2.1.0")
add_compiler("/opt/kotlin-language", "1.8.22")

local found = plugin:Detect({})
local notes = {}
local versions = {}
for _, entry in ipairs(found) do
    notes[entry.path] = entry.note
    versions[entry.version] = true
end

assert(notes["/opt/kotlin-language"] == "pointed to by KOTLIN_HOME")
assert(notes["/home/u/.sdkman/candidates/kotlin/2.0.0"] == "installed by SDKMAN")
assert(notes["/home/u/.asdf/installs/kotlin/1.9.24"] == "installed by asdf")
assert(notes["/home/u/.local/share/mise/installs/kotlin/2.1.0"] == "installed by mise")
assert(versions["2.1.0"] and versions["1.8.22"], "all four compilers must be found")
assert(#found == 4, "expected 4 compilers, got " .. #found)

print("PASS: Detect honours KOTLIN_HOME and version-manager roots")
