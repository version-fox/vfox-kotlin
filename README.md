# vfox-kotlin

[中文](README_CN.md) | [English](README.md)

Kotlin plugin for [vfox](https://vfox.dev/).

It installs the Kotlin compiler from [JetBrains' releases](https://github.com/JetBrains/kotlin/releases)
onto your `PATH`. Releases from 1.9.0 on are verified against the `.sha256`
file JetBrains publishes beside the archive.

## Install

After installing [vfox](https://github.com/version-fox/vfox), install the plugin,
a compiler, and select it:

```bash
vfox add kotlin
vfox install kotlin@2.0.0
vfox use --global kotlin@2.0.0
kotlinc -version
```

Use `kotlin@latest` for the newest release on the
[Kotlin release page](https://kotlinlang.org/docs/releases.html#release-details),
and `vfox search kotlin` to list every release. For `kotlinc` and `kotlin` to be
on your `PATH`, vfox has to be hooked into your shell — see its
[quick start](https://vfox.dev/guides/quick-start.html).

## Versions

A version is a Kotlin release number, such as `2.0.0` or `1.9.24`, with no
distribution suffix:

```bash
vfox install kotlin@2.1.0      # a specific release
vfox install kotlin@latest     # the newest release
vfox use kotlin@2.1.0          # the current shell session
vfox use --global kotlin@2.1.0 # every new shell
```

See [versions.md](docs/versions.md) for details, including the checksum
behaviour of releases before 1.9.0.

## Requirements

`kotlinc` and `kotlin` are JVM launchers, so a Java runtime (17 or newer) must
already be on your `PATH`. The plugin does not install Java; use the
[vfox-java](https://github.com/version-fox/vfox-java) plugin if you need one.

## Already installed compilers

If a Kotlin compiler is already on the machine — from SDKMAN, asdf, mise,
Homebrew, an IDE bundle, the system package manager or `KOTLIN_HOME` — adopt it
instead of downloading another copy:

```bash
vfox link --detect kotlin
```

This needs vfox 1.1.0 or later. vfox lists what the plugin found and links what
you pick; `vfox link --detect -y kotlin` links everything without asking.
See [versions.md](docs/versions.md#already-installed-compilers).

## Documentation

**Using the plugin**

- [Versions](docs/versions.md): version syntax, installing, searching and selecting a compiler
- [Install from source](docs/install-from-source.md): try an unreleased change from this checkout
- [Troubleshooting](docs/troubleshooting.md): Java, shell hooks, checksums and downloads

**Working on the plugin**

- [Development](docs/development.md): layout, hook contracts, tests and releasing
- [E2E](docs/e2e.md): local checks and the Linux/Windows container suites
- [Changes](docs/changes.md): what changed in each release

## Releasing this plugin

Maintainers can publish from **Actions → Plugin → Run workflow** on the default
branch by entering a stable plugin version without the `v` prefix. The shared
workflow updates `metadata.lua`, creates the version commit and tag, and publishes
the ZIP and manifest in this repository. No local tag or extra release token is
needed. Pull requests run checks only; PR titles no longer trigger publication.

Existing version-tag pushes are supported when `PLUGIN.version` already matches
the tag. If publication fails, re-run the original failed job to resume it.

The workflow follows the shared `@v1` release-tool version. Updating the tool does
not release this plugin. See the [shared workflow documentation](https://github.com/version-fox/plugin-manifest-action)
for the package contract and first-rollout requirements.

## License

Apache-2.0. See [LICENSE](LICENSE).
