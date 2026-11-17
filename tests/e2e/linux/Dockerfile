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

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# NOTE: no heredoc here. testcontainers-go builds through the classic builder,
# which accepts `RUN <<EOF` but silently drops its body. The chained-`;`
# form works on both builders.
#
# A JRE is required: kotlinc and kotlin are JVM launchers. It is not managed
# by vfox, so the compiler on PATH can only come from what vfox installs.
RUN set -o errexit -o nounset -o xtrace; \
    apt-get update; \
    apt-get install --assume-yes --no-install-recommends \
        ca-certificates \
        curl \
        openjdk-17-jre-headless \
        unzip \
        zip; \
    rm --recursive --force /var/lib/apt/lists/*

COPY metadata.lua /e2e/metadata.lua
COPY hooks /e2e/hooks
COPY lib /e2e/lib
COPY Injection.lua /e2e/Injection.lua

# Long-lived container: the Go test drives every phase through exec calls.
CMD ["sleep", "infinity"]
