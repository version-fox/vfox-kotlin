# escape=`
# Copyright 2026 Han Li and contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
ARG BASE_IMAGE=mcr.microsoft.com/windows/servercore:ltsc2025-KB5122871-amd64
FROM ${BASE_IMAGE}

SHELL ["powershell", "-NoProfile", "-Command"]

# PowerShell 7: the Go test drives the container with `pwsh -NoProfile
# -Command`. The stock servercore only has Windows PowerShell 5.1.
ENV POWERSHELL_TELEMETRY_OPTOUT=1

RUN $ErrorActionPreference = 'Stop'; `
    curl.exe -fsSL --retry 3 --retry-delay 5 --retry-all-errors -o C:\pwsh.zip 'https://github.com/PowerShell/PowerShell/releases/download/v7.6.6/PowerShell-7.6.6-win-x64.zip'; `
    New-Item -ItemType Directory -Force -Path 'C:\Program Files\PowerShell\7' | Out-Null; `
    tar.exe -xf C:\pwsh.zip -C 'C:\Program Files\PowerShell\7'; `
    Remove-Item C:\pwsh.zip; `
    if (-not (Test-Path 'C:\Program Files\PowerShell\7\pwsh.exe')) { throw 'pwsh.exe is missing after extraction' }; `
    [Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'Machine') + ';C:\Program Files\PowerShell\7', 'Machine'); `
    New-ItemProperty -Path 'HKLM:\SYSTEM\CurrentControlSet\Control\FileSystem' -Name LongPathsEnabled -Value 1 -PropertyType DWord -Force | Out-Null

# A JRE is required: kotlinc and kotlin are JVM launchers. Pin the Temurin 17
# JRE from Adoptium's stable API and add its bin to the machine PATH.
RUN $ErrorActionPreference = 'Stop'; `
    Invoke-WebRequest -UseBasicParsing -MaximumRedirection 5 -Uri 'https://api.adoptium.net/v3/binary/latest/17/ga/windows/x64/jre/hotspot/normal/eclipse' -OutFile C:\jre.zip; `
    Expand-Archive -Path C:\jre.zip -DestinationPath C:\jre -Force; `
    $jreBin = (Get-ChildItem -Path C:\jre -Recurse -Filter java.exe | Select-Object -First 1).Directory.FullName; `
    if (-not $jreBin) { throw 'No java.exe found in the JRE archive' }; `
    [Environment]::SetEnvironmentVariable('Path', $jreBin + ';' + [Environment]::GetEnvironmentVariable('Path', 'Machine'), 'Machine'); `
    Remove-Item C:\jre.zip

WORKDIR C:\e2e

COPY metadata.lua C:\e2e\metadata.lua
COPY hooks C:\e2e\hooks
COPY lib C:\e2e\lib
COPY Injection.lua C:\e2e\Injection.lua

ENTRYPOINT ["C:/Program Files/PowerShell/7/pwsh.exe", "-NoProfile", "-Command"]
CMD ["while ($true) { Start-Sleep 3600 }"]
