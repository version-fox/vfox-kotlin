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

--- Each SDK may have different environment variable configurations.
--- This allows plugins to define custom environment variables (including PATH settings)
--- Note: Be sure to distinguish between environment variable settings for different platforms!
--- @param ctx table Context information; ctx.path is the SDK installation directory
function PLUGIN:EnvKeys(ctx)
    --- this variable is same as ctx.sdkInfo['plugin-name'].path
    local version_path = ctx.path
    local result = {
        {
            key = "PATH",
            value = version_path .. "/bin",
        },
    }
    -- PreInstall installs Kotlin/Native as an additional package when the
    -- host has a prebuilt archive. vfox keys additions in sdkInfo by their
    -- runtime directory name (add-native-<version>), so add the bin of every
    -- non-main package; that is where kotlinc-native, konanc and klib live.
    -- Linked and older installs have no such addition.
    for name, sdk in pairs(ctx.sdkInfo or {}) do
        if name ~= "kotlin" and sdk.path ~= nil and sdk.path ~= "" then
            table.insert(result, {
                key = "PATH",
                value = sdk.path .. "/bin",
            })
        end
    end
    return result
end
