# Development

[中文](zh-CN/development.md) | [English](development.md)

## Layout

```
metadata.lua                    plugin name, version, homepage and notes
hooks/available.lua             vfox search and kotlin@latest
hooks/pre_install.lua           vfox install: resolve the URL and checksum
hooks/post_install.lua          no-op
hooks/env_keys.lua              adds <sdk>/bin to PATH
hooks/pre_use.lua               no-op
hooks/detect.lua                vfox link --detect: find compilers already installed
lib/util.lua                    download URL template and version comparison
lib/detect_kotlin.lua           scanning roots, build.txt version and compiler home checks
tests/*_test.lua                offline hook tests, run by vfox plugin test
tests/detect*_test.lua          offline tests of the Detect hook
tests/e2e/                      container tests, see e2e.md
Injection.lua                   type hints for what vfox injects, never loaded
README.md, docs/                English documentation; README_CN.md and
                                docs/zh-CN/ are the Chinese version
opencode.jsonc                  lets OpenCode agent sessions run git push,
                                gh pr create and gh pr comment without asking
stylua.toml, .luarc.json        Lua formatting and language-server settings
.licenserc.yaml, .license.tpl   Apache-2.0 header check and its text template
.devcontainer/                  Linux development environment for the E2E suite
.github/workflows/              CI: publish, test, lint, license and E2E
```

vfox loads `metadata.lua`, calls a hook when a command needs it, and puts
`PLUGIN` and `RUNTIME` (`osType`, `archType`, `version`) in scope. `osType` and
`archType` are Go's `GOOS` and `GOARCH` of the vfox binary.

## Hook contracts

| You run | Hooks |
| --- | --- |
| `vfox search kotlin` | `Available` |
| `vfox install kotlin@2.0.0` | `PreInstall`; vfox downloads, verifies and unpacks |
| `vfox use kotlin@2.0.0` | `PreUse` (no-op), then `EnvKeys` |
| `vfox link --detect kotlin` | `Detect` returns compilers already on the machine |
| a shell hook | `EnvKeys` for the selected version |

- `Available` fetches the release details page and parses `tbody tr` / `td b`.
  It returns a list of `{ version, note }`. A failed request returns an empty
  list rather than raising, so `vfox search` stays usable offline.
- `PreInstall` formats `lib/util.lua`'s URL with the version twice. It sends a
  HEAD request and raises `Current version information not detected.` unless
  the answer is 200. For a version at or above 1.9.0 it GETs `<url>.sha256` and
  returns its body as `sha256`; otherwise it returns no checksum. `version` in
  the result is the name vfox installs under, which is why `latest` is resolved
  to a concrete release here. It also probes the per-platform Kotlin/Native
  asset from `lib/util.lua`; when the HEAD request succeeds, the archive is
  returned in `addition` (with its own `.sha256` when one exists), so vfox
  downloads and extracts it as a second package. An absent native asset only
  leaves out `kotlinc-native`; the compiler still installs.
- `EnvKeys` returns `{ key = "PATH", value = ctx.path .. "/bin" }` and one more
  PATH entry for every non-main package in `ctx.sdkInfo`. vfox strips the
  top-level `kotlinc` directory of the archive, so the launchers sit directly
  in `<sdk>/bin`, and the Kotlin/Native addition's `bin` holds
  `kotlinc-native`, `konanc` and `klib`. vfox keys additions by their runtime
  directory name (`add-native-<version>`), not by the `name` sent in
  `PreInstall`.
- `PreUse` and `PostInstall` are no-ops; vfox matches an installed version and
  applies `EnvKeys` without help from the plugin.
- `Detect` scans SDKMAN, asdf, mise, Homebrew, IDE bundles, the system
  directories and `KOTLIN_HOME` for compiler homes with `bin/kotlinc` and
  `build.txt`. It returns `{ version, path, note }`; the version comes from
  `build.txt` (falling back to the directory name), so a detected compiler is
  named exactly like an installed one. Paths under `.vfox` are skipped, and
  results are deduplicated and sorted. The hook needs vfox 1.1.0; older vfox
  simply ignores it.

## Tests

The hook tests run offline under `vfox plugin test`, available in vfox 1.0.13
and later. It finds every `tests/**/*_test.lua` and runs each file in a VM of
its own, so a new test file needs no entry anywhere:

```bash
vfox plugin test .
vfox plugin test . --file tests/pre_install_test.lua
```

A test file calls `test.load()` from `vfox.test` once, then calls the real
hooks through what it returns. `test.load` takes the platform (`os`, `arch`)
and a function that answers every HTTP request (`http`). For example:

```lua
local test = require("vfox.test")
local plugin = test.load({
    http = function(request)
        return { status_code = 200, body = "..." }
    end,
})
local result = plugin:PreInstall({ version = "2.0.0" })
```

The context and result go through the decoding vfox does, so a hook that
returns the wrong type fails the test. Error paths use `pcall` on the hook.

The container tests are Go programs under `tests/e2e/` using testcontainers-go.
They install a pinned vfox, add this checkout as a local plugin, install a
compiler and check the result on Linux and Windows. They are skipped by
`go test -short`; see [e2e.md](e2e.md) for how to run them and what they check.

## Local checks

Lint, formatting, license headers and the container tests are documented under
[Local checks](e2e.md#local-checks).

## Releasing

A maintainer publishes from **Actions → Plugin → Run workflow** on the default
branch, entering a stable plugin version without the `v` prefix. The shared
workflow, [plugin-manifest-action](https://github.com/version-fox/plugin-manifest-action),
updates `metadata.lua`, creates the version commit and the tag, and publishes
the ZIP and the manifest in this repository. No local tag or extra release
token is needed. Pull requests run the checks only.

Pushing a version tag publishes too, when `PLUGIN.version` already matches the
tag. If publication fails, re-run the failed job to resume it. The workflow
follows the shared `@v1` release tool, and updating the tool does not release
this plugin. Afterwards `plugins/kotlin.json` in
[vfox-plugins](https://github.com/version-fox/vfox-plugins) shows the version
that `vfox add kotlin` serves.

### After a release

Mark changes newer than the last release as "Unreleased" ("未发布" in Chinese)
in [changes.md](changes.md), then move them under the new version with the
release date:

```bash
grep -rn -i -e unreleased -e 未发布 docs README*
```

## Writing the docs

English pages are in `docs/`, Chinese ones in `docs/zh-CN/` under the same
names (the Chinese README is `README_CN.md`). Each page starts with a link to
the other language, and every new page is added to the Documentation list of
both READMEs. Nothing in CI checks the prose, so read both languages when you
change a page.
