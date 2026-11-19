# E2E

[中文](zh-CN/e2e.md) | [English](e2e.md)

The offline Lua hook tests (`tests/*_test.lua`) run in the **Test Plugin**
workflow. The end-to-end suite runs in containers against real vfox on Linux
and Windows; each case gets its own throwaway container. CI runs it natively on
the stock runners (Docker plus Go from `actions/setup-go`). The dev container in
`.devcontainer/` reproduces the Linux path locally with the same pinned tools.

## Environment knobs

| Variable              | Default  | Meaning                                                     |
| --------------------- | -------- | ----------------------------------------------------------- |
| `KOTLIN_VERSION`      | `2.0.0`  | Release under test; its case runs every step (with checksum from 1.9.0-RC on) |
| `E2E_MAX_JOBS`        | `3`      | Containers (cases) that run at once in one `go test` run    |
| `E2E_RETRY_ATTEMPTS`  | `3`      | Attempts per network-touching step (10s linear backoff)      |

Each platform runs two cases, one container each:

- the release under test (`KOTLIN_VERSION`) runs every step below;
- `1.8.22`, the last stable release without a `.sha256` asset, runs steps 1 and 3 only.
  It gets no case of its own when it is the release under test.

With the default `E2E_MAX_JOBS=3` both cases run at once. `E2E_MAX_JOBS=1` runs
them one after the other.

## Linux (Ubuntu)

With Docker Engine and Go 1.27.1 on the host, from the repository root:

```bash
cd tests/e2e
go test -v -timeout 45m -run TestE2E ./linux
```

Rootless Docker users point the client at their socket with
`export DOCKER_HOST="unix://$XDG_RUNTIME_DIR/docker.sock"`. To iterate on one
release:

```bash
KOTLIN_VERSION=1.8.22 go test -v -timeout 45m -run TestE2E ./linux
```

The suite builds `tests/e2e/linux/Dockerfile` (Ubuntu 24.04, a JRE, curl and
zip), installs a pinned vfox with a checked SHA-256, adds the local plugin from
a ZIP, then drives the install and verification through `bash` exec calls.

## Windows 11

Windows containers need a Windows Docker daemon in Windows-container mode.
CI runs the suite on `windows-latest`: the job starts the `docker` service if
it is stopped, and fails unless the daemon reports Windows-container mode. From
a Windows host with Docker and Go:

```powershell
cd tests/e2e
go test -v -timeout 55m -run TestE2E ./windows
```

The image (`tests/e2e/windows/Dockerfile`) is servercore with PowerShell 7 and
a Temurin 17 JRE. The Go test drives it with `pwsh -NoProfile -Command`, using
the same pinned vfox archive and local plugin ZIP as Linux. Every command runs
with `VFOX_HOME` under a path containing a space, the case vfox handles since
1.0.13.

## Dev container

`.devcontainer/` opens a Linux environment with the pinned Go, stylua,
lua-language-server, license tools and vfox the workflows use. It mounts the
host Docker socket so the suite spawns the same target containers as CI
(Docker-out-of-Docker). Open it with a Dev Containers tool, for example:

```bash
devcontainer up --workspace-folder .
devcontainer exec --workspace-folder . bash
```

Windows containers cannot run in a Linux dev container, so the Windows suite
needs a Windows host as described above.

## Local checks

Execute from the dev container (all tools are preinstalled and pinned to the
versions the workflows use), or install the same tools on a host.

```bash
# Lua and Go checks (.github/workflows/lint.yml).
cd tests/e2e
go test -short ./...
test -z "$(gofmt -l .)"
go vet ./...
cd ../..

stylua --check .
lua-language-server --check . --checklevel=Warning
actionlint
```

```bash
# Offline hook tests (.github/workflows/test.yml).
vfox plugin test .
```

```bash
# License headers (.github/workflows/license.yml).
docker run --rm -v "$PWD:/work" -w /work apache/skywalking-eyes:0.9.0 header check
docker run --rm -v "$PWD:/work" -w /work apache/skywalking-eyes:0.9.0 header fix
```

## What the suite verifies

Both platforms run the same steps. Every case runs steps 1 and 3; the release
under test also runs steps 2 and 4 to 7:

1. vfox is installed from a pinned, checksum-checked archive and the local
   `kotlin` plugin is registered from this checkout — not from the registry.
2. `vfox search kotlin` lists the installed release.
3. `vfox install kotlin@<version>` installs under that exact name. For 1.9.0-RC
   and newer, vfox prints `Verifying checksum <64 hex digits>…`; for the others,
   such as 1.8.22 and 1.9.0-Beta, it reports that no checksum is provided. 2.0.0
   covers the first branch and 1.8.22 the second.
4. Installing a release JetBrains never published fails in `PreInstall`'s HEAD
   check.
5. `vfox use --global` selects the compiler. In a fresh shell the SDK's `bin`
   is on `PATH`, `vfox env` exports it, `kotlinc -version` and
   `kotlin -version` report the requested version, and `vfox current kotlin`
   names the selection.
6. On amd64, for releases from 1.9.0 on, `kotlinc-native` is on `PATH` and
   `kotlinc-native -version` reports the same release.
7. `vfox link --detect kotlin` finds a compiler copied outside vfox's own
   directories and pointed to by `KOTLIN_HOME`, links it, and the link shows up
   in `vfox list kotlin` as `(linked → …)`. The pinned vfox 1.1.0 is the first
   release with the `Detect` hook.
