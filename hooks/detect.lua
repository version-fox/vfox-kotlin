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

local detect = require("detect_kotlin")

--- Finds Kotlin compilers that are already on the machine, for
--- `vfox link --detect kotlin`: SDKMAN, asdf, mise, Homebrew, IDE bundles,
--- the system directories and KOTLIN_HOME. Each entry is named the way
--- `vfox install` names it, from the compiler's build.txt.
function PLUGIN:Detect(ctx)
    return detect.find_all()
end
