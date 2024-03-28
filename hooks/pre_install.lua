local util = require("util")
local http = require("http")

--- Returns some pre-installed information, such as version number, download address, local files, etc.
--- If checksum is provided, vfox will automatically check it for you.
--- @param ctx table
--- @field ctx.version string User-input version
--- @return table Version information
function PLUGIN:PreInstall(ctx)
    local version = ctx.version

    if version == "latest" then
        local lists = self:Available({})
        version = lists[1].version
    end
    local url = util.DownloadURL:format(version, version)
    local resp, err = http.head({
        url = url
    })
    if err ~= nil or resp.status_code ~= 200 then
        error("Current version information not detected.")
    end
    if util:compare_versions({ version = version },{version = "1.9.0"}) or version=="1.9.0" then
        resp, err = http.get({
            url = url .. ".sha256"
        })
        if err ~= nil or resp.status_code ~= 200 then
            error("Current version checksum not detected.")
        end
        return {
            version = version,
            url = url,
            sha256 = resp.body
        }
    else
        return {
            version = version,
            url = url
        }
    end
end