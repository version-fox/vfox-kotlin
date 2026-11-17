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

--- Finds Kotlin compilers that are already on the machine, for
--- `vfox link --detect kotlin`. A compiler home holds `bin/kotlinc` (or
--- `kotlinc.bat` on Windows) and a `build.txt` with its version, such as
--- `1.9.0-release-358`. Each result is named the way `vfox install` names it.
local detect = {}

local function trim(value)
    return (value:gsub("^%s+", ""):gsub("%s+$", ""))
end

local function lower(value)
    return (value or ""):lower()
end

--- Returns the separator for the current platform: `\` on Windows, `/` else.
local function sep()
    if RUNTIME and RUNTIME.osType == "windows" then
        return "\\"
    end
    return "/"
end

--- Joins two path segments with the platform separator.
function detect.join(left, right)
    local s = sep()
    if left:sub(-1) == "/" or left:sub(-1) == "\\" then
        return left .. right
    end
    return left .. s .. right
end

local function home_dir()
    return os.getenv("HOME") or os.getenv("USERPROFILE") or ""
end

local function shell_quote(value)
    if RUNTIME and RUNTIME.osType == "windows" then
        return '"' .. value:gsub('"', '""') .. '"'
    end
    return "'" .. value:gsub("'", "'\\''") .. "'"
end

--- Lists the immediate children of dir as full paths, or {}.
function detect.list_children(dir)
    if dir == "" then
        return {}
    end
    local cmd
    if RUNTIME and RUNTIME.osType == "windows" then
        cmd = "dir /b /ad " .. shell_quote(dir) .. " 2>nul"
    else
        cmd = "ls -1 " .. shell_quote(dir) .. " 2>/dev/null"
    end
    local handle = io.popen(cmd, "r")
    if not handle then
        return {}
    end
    local result = {}
    local function add(line)
        line = trim(line)
        if line ~= "" and line ~= "." and line ~= ".." then
            table.insert(result, detect.join(dir, line))
        end
    end
    -- Some vfox runtimes do not give popen handles a `lines` method; read
    -- line by line there (notably on Windows).
    if handle.lines then
        for line in handle:lines() do
            add(line)
        end
    else
        while true do
            local line = handle:read("*l")
            if not line then
                break
            end
            add(line)
        end
    end
    if handle.close then
        handle:close()
    end
    return result
end

local function read_file(path)
    local file = io.open(path, "r")
    if not file then
        return nil
    end
    local content = file:read("*a")
    file:close()
    return content
end

--- Extracts the version from a Kotlin `build.txt` such as
--- `1.9.0-release-358` or `2.1.0-Beta2`. Returns nil if there is none.
function detect.version_from_build(text)
    if type(text) ~= "string" then
        return nil
    end
    local version = trim(text):match("^(%d+%.[%d%.]+)")
    if not version then
        return nil
    end
    version = version:gsub("%.$", "")
    if not version:match("^%d+%.%d+") then
        return nil
    end
    return version
end

--- The version of a compiler home: its `build.txt` first, then a version in
--- the directory name (as SDKMAN, asdf and mise name their installs).
function detect.version_for(home)
    local from_build = detect.version_from_build(read_file(detect.join(home, "build.txt")) or "")
    if from_build then
        return from_build
    end
    local name = home:match(".*[/\\]([^/\\]+)$")
    if name then
        return detect.version_from_build(name .. "-")
    end
    return nil
end

--- Whether path is a Kotlin compiler home: a `bin` launcher and the compiler
--- jar or `build.txt` are both present.
function detect.kotlin_home(path)
    if path == "" then
        return false
    end
    local launcher = detect.join(path, "bin")
    local on_windows = RUNTIME and RUNTIME.osType == "windows"
    local has_launcher = read_file(detect.join(launcher, on_windows and "kotlinc.bat" or "kotlinc")) ~= nil
        or read_file(detect.join(launcher, "kotlinc.bat")) ~= nil
        or read_file(detect.join(launcher, "kotlinc.cmd")) ~= nil
    if not has_launcher then
        return false
    end
    return read_file(detect.join(path, "build.txt")) ~= nil
        or read_file(detect.join(path, detect.join("lib", "kotlin-compiler.jar"))) ~= nil
end

--- Candidate roots whose children may be Kotlin compiler homes:
--- {path, note, deep}. Deep roots are scanned a few levels down for IDE
--- bundles, which keep the compiler under `plugins/kotlin/kotlinc`.
function detect.candidate_roots()
    local roots = {}
    local function add(path, note, deep)
        if type(path) == "string" and path ~= "" then
            table.insert(roots, { path = path, note = note, deep = deep })
        end
    end
    local function child(base, ...)
        if base == "" then
            return ""
        end
        local result = base
        for _, part in ipairs({ ... }) do
            result = detect.join(result, part)
        end
        return result
    end

    local home = home_dir()
    local s = sep()

    -- Env roots first, so the note explains why the compiler was included
    -- even when the same path also sits under a later system root.
    local kotlin_home = os.getenv("KOTLIN_HOME")
    if kotlin_home and kotlin_home ~= "" then
        add(kotlin_home, "pointed to by KOTLIN_HOME")
    end

    if RUNTIME and RUNTIME.osType == "windows" then
        local userprofile = os.getenv("USERPROFILE") or home
        local programfiles = os.getenv("ProgramFiles") or ("C:" .. s .. "Program Files")
        local programfiles_x86 = os.getenv("ProgramFiles(x86)")
        add(child(userprofile, ".sdkman", "candidates", "kotlin"), "installed by SDKMAN")
        add(child(userprofile, ".asdf", "installs", "kotlin"), "installed by asdf")
        add(child(userprofile, ".local", "share", "mise", "installs", "kotlin"), "installed by mise")
        add(child(userprofile, "scoop", "apps"), "installed by Scoop", 2)
        add(child(programfiles, "Kotlin"), "system Kotlin compiler")
        if programfiles_x86 and programfiles_x86 ~= "" then
            add(child(programfiles_x86, "Kotlin"), "system Kotlin compiler")
        end
        add(child(programfiles, "JetBrains"), "bundled with IntelliJ IDEA", 4)
        add(child(userprofile, "AppData", "Local", "Programs"), "installed by a user installer", 2)
    else
        add(child(home, ".sdkman", "candidates", "kotlin"), "installed by SDKMAN")
        add(child(home, ".asdf", "installs", "kotlin"), "installed by asdf")
        add(child(home, ".local", "share", "mise", "installs", "kotlin"), "installed by mise")
        if RUNTIME and RUNTIME.osType == "darwin" then
            add(
                child(home, "Library", "Application Support", "JetBrains", "Toolbox", "apps"),
                "bundled with IntelliJ IDEA",
                4
            )
            add("/Applications", "bundled with IntelliJ IDEA", 3)
        else
            add(child(home, ".local", "share", "JetBrains", "Toolbox", "apps"), "bundled with IntelliJ IDEA", 4)
        end
        add("/opt/homebrew/opt", "installed by Homebrew")
        add("/usr/local/opt", "installed by Homebrew")
        add("/home/linuxbrew/.linuxbrew/opt", "installed by Linuxbrew")
        add(child(home, ".linuxbrew", "opt"), "installed by Linuxbrew")
        add("/usr/share/kotlinc", "system Kotlin compiler")
        add("/opt/kotlinc", "system Kotlin compiler")
        add("/usr/local/kotlinc", "system Kotlin compiler")
        add("/opt", "system Kotlin compiler")
        add("/usr/local", "system Kotlin compiler")
    end

    return roots
end

--- True for paths vfox manages, which must never be detected.
local function is_vfox_managed(path)
    local p = lower(path or "")
    return p:find("/.vfox/", 1, true) ~= nil or p:find("\\.vfox\\", 1, true) ~= nil
end

local function add_home(path, note, seen, out)
    if seen[path] then
        return
    end
    seen[path] = true
    local version = detect.version_for(path)
    if version then
        table.insert(out, { version = version, path = path, note = note })
    end
end

local function find_shallow(root, note, seen, out)
    if detect.kotlin_home(root) then
        add_home(root, note, seen, out)
        return
    end
    for _, dir in ipairs(detect.list_children(root)) do
        local name = dir:match("[/\\]([^/\\]+)$") or ""
        if not is_vfox_managed(dir) and name:sub(1, 1) ~= "." then
            if detect.kotlin_home(dir) then
                add_home(dir, note, seen, out)
            end
        end
    end
end

local function find_deep(root, note, depth, seen, out)
    if depth <= 0 then
        return
    end
    if detect.kotlin_home(root) then
        add_home(root, note, seen, out)
        return
    end
    for _, dir in ipairs(detect.list_children(root)) do
        local name = dir:match("[/\\]([^/\\]+)$") or ""
        -- Skip the noisy trees; an IDE bundle is never inside them.
        if
            not is_vfox_managed(dir)
            and name:sub(1, 1) ~= "."
            and name ~= "node_modules"
            and name ~= "lib"
            and name ~= "license"
        then
            find_deep(dir, note, depth - 1, seen, out)
        end
    end
end

local function version_key(value)
    local parts = {}
    for number in value:gmatch("%d+") do
        table.insert(parts, tonumber(number))
    end
    return parts
end

local function less_version(left, right)
    local a, b = version_key(left), version_key(right)
    for i = 1, math.max(#a, #b) do
        local x, y = a[i] or 0, b[i] or 0
        if x ~= y then
            return x < y
        end
    end
    return false
end

--- Finds every Kotlin compiler on the machine. Results are sorted by version
--- and path so `vfox link --detect` offers them deterministically.
function detect.find_all()
    local out = {}
    local seen = {}
    for _, root in ipairs(detect.candidate_roots()) do
        if not is_vfox_managed(root.path) then
            if root.deep then
                find_deep(root.path, root.note, root.deep, seen, out)
            else
                find_shallow(root.path, root.note, seen, out)
            end
        end
    end
    table.sort(out, function(a, b)
        if a.version == b.version then
            return a.path < b.path
        end
        return less_version(a.version, b.version)
    end)
    return out
end

return detect
