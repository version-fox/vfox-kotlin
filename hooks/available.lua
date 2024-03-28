local html = require("html")
local http = require("http")
--- Return all available versions provided by this plugin
--- @param ctx table Empty table used as context, for future extension
--- @return table Descriptions of available versions and accompanying tool descriptions
function PLUGIN:Available(ctx)
    local resp, err = http.get({
        url = "https://kotlinlang.org/docs/releases.html#release-details"
    })
    if err ~= nil or resp.status_code ~= 200 then
        return {}
    end
    local result = {}
    html.parse(resp.body):find("tbody tr"):each(function(i, sel)
        local version = sel:find("td b"):first():text()
        table.insert(result, {
            version = version,
            note = "",
        })
    end)
    return result
end