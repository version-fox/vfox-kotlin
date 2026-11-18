# Changes

[中文](zh-CN/changes.md) | [English](changes.md)

What changed in each release. GitHub's release pages list the pull requests.

## Unreleased

- Fix `latest` failing with `Current version information not detected.`: the release page's first table lists release lines such as `2.4`, and `Available` listed them as versions, so `latest` resolved to `2.4`. `Available` now keeps only `X.Y.Z` releases, and `latest` installs the newest one.

## 0.1.2 (2026-11-17)

- Bundle Kotlin/Native with the JVM compiler. `PreInstall` adds JetBrains' prebuilt native archive as an additional package on hosts that have one, `EnvKeys` puts its `bin` directory on `PATH`, and `kotlinc-native`, `konanc` and `klib` are now available next to `kotlinc`. The asset name follows the release (it changed in Kotlin 2.0) and is verified with its own `.sha256` when published. Needs vfox 1.0.0 or later; hosts without a native build still get the JVM compiler. The Linux and Windows E2E suites now check `kotlinc-native -version`.
- Bilingual README and documentation in English and Chinese, covering versions, installing from source, development, E2E, troubleshooting and this changelog.
- Offline hook tests for `Available`, `PreInstall`, `EnvKeys` and the version helper, run by `vfox plugin test` in the **Test Plugin** workflow.
- E2E suites for Linux (Ubuntu container) and Windows (Windows container), with a GitHub Actions workflow for each. They install a pinned vfox, add the plugin from this checkout, install a compiler and check PATH and the version banner.
- Lint and formatting checks for Go, Lua and the workflows (`gofmt`, `go vet`, stylua, lua-language-server, actionlint), with pinned tooling and a `stylua.toml` / `.luarc.json`.
- Apache-2.0 header checks with Apache SkyWalking Eyes, plus `.licenserc.yaml`, `.license.tpl` and headers on every source file.
- `opencode.jsonc` allowing OpenCode sessions in this repository to run `git push`, `gh pr create` and `gh pr comment`.
- `vfox link --detect kotlin` adopts Kotlin compilers that are already on the machine — from SDKMAN, asdf, mise, Homebrew, IDE bundles, system directories and `KOTLIN_HOME`. A compiler home is recognised by `bin/kotlinc` and `build.txt`, and is named by the version in that file. Needs vfox 1.1.0 or later; documented in [link.md](link.md).
- Repository hygiene: `.gitattributes`, `.dockerignore`, an expanded `.gitignore` and a Linux dev container. The pinned vfox in tests and E2E is now 1.1.0.

## 0.1.1 (2026-09-12)

- Publishing moved to the shared `plugin-manifest-action` workflow. Maintainers release from **Actions → Plugin → Run workflow**; pull requests only run checks.

## 0.1.0 (2024-03-28)

- First release: install the Kotlin compiler from JetBrains' GitHub releases, resolve `latest` from the Kotlin release page, verify releases from 1.9.0 on against their `.sha256` asset, and add `<sdk>/bin` to `PATH`.
