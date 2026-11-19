# Versions

[中文](zh-CN/versions.md) | [English](versions.md)

A version is a Kotlin release number such as `2.0.0` or `1.9.24`. The plugin has
no distributions or suffixes: the version you type is the version of the
compiler you get.

## Installing

```bash
vfox install kotlin@2.1.0   # a specific release
vfox install kotlin@latest  # the newest release on the Kotlin release page
```

`latest` is resolved by `PreInstall`, which takes the first `X.Y.Z` release that
the `Available` hook lists from the
[Kotlin release details](https://kotlinlang.org/docs/releases.html#release-details)
table. If that page cannot be reached, or lists no release,
`vfox install kotlin@latest` fails rather than guessing a version.

## Listing

```bash
vfox search kotlin
```

`Available` lists the `X.Y.Z` releases from the same table in page order, newest
first, and marks releases that are already installed. The release lines at the
top of that page, such as `2.4`, are not versions and are skipped. The table is
upstream's, so a release appears here as soon as JetBrains adds it.

## Selecting

```bash
vfox use kotlin@2.1.0          # current shell session
vfox use --project kotlin@2.1.0 # this project
vfox use --global kotlin@2.1.0 # every new shell
vfox current kotlin            # the selected version
```

`EnvKeys` adds `<sdk>/bin` to `PATH`, which is where the `kotlinc`, `kotlin`,
`kotlinc-jvm` and `kotlinc-js` launchers live inside the installed archive. It
also adds the Kotlin/Native compiler's `bin` directory when one is installed.
A shell hook must be active for the variables to be applied; see vfox's
[quick start](https://vfox.dev/guides/quick-start.html).

## Kotlin/Native

Kotlin/Native is not part of the `kotlin-compiler` archive: JetBrains publishes
it as a separate asset. `PreInstall` adds it as an additional package whenever
a prebuilt compiler exists for the host, so a single `vfox install` provides
both toolchains:

```bash
kotlinc-native -version
kotlinc-native hello.kt -o hello
```

| Host | Native archive | From |
| --- | --- | --- |
| Linux x86_64 | `kotlin-native-prebuilt-linux-x86_64-<version>.tar.gz` | 2.0.0-Beta3 |
| Linux x86_64 | `kotlin-native-linux-x86_64-<version>.tar.gz` | 1.9.0 |
| macOS x86_64 | `kotlin-native-prebuilt-macos-x86_64-<version>.tar.gz` | 2.0.0-Beta3 |
| macOS arm64 | `kotlin-native-prebuilt-macos-aarch64-<version>.tar.gz` | 2.0.0-Beta3 |
| macOS x86_64 | `kotlin-native-macos-x86_64-<version>.tar.gz` | 1.9.0 |
| macOS arm64 | `kotlin-native-macos-aarch64-<version>.tar.gz` | 1.9.0 |
| Windows x86_64 | `kotlin-native-prebuilt-windows-x86_64-<version>.zip` | 2.0.0-Beta3 |
| Windows x86_64 | `kotlin-native-windows-x86_64-<version>.zip` | 1.9.0 |

The archive name changed in Kotlin 2.0.0-Beta3; 2.0.0-Beta1 and 2.0.0-Beta2 keep
the old name. The plugin chooses the canonical name for the release and verifies
it against the `.sha256` file when JetBrains publishes one. Releases before 1.9.0
have no native archive and install the JVM compiler only. There is no prebuilt
native compiler for Linux arm64, so those hosts also get the JVM compiler only.
Bundling the native compiler requires vfox 1.0.0 or later.

Like `kotlinc`, `kotlinc-native` is a JVM launcher and needs Java on `PATH`.

## Checksums

JetBrains publishes a `.sha256` file beside each compiler archive from 1.9.0-RC
on; 1.9.0-Beta has none. For those releases `PreInstall` returns the digest to
vfox, so the install ends with `Verifying checksum …` and a download that does
not match the file fails with `Checksum error`.

Releases before 1.9.0-RC have no such asset. They still install, but vfox prints
`WARNING: Checksum is not provided, skip verify...`. This is expected; see
[troubleshooting.md](troubleshooting.md) if an older release fails for another
reason.

## Already installed compilers

If a compiler from SDKMAN, asdf, mise, Homebrew, an IDE, the system package
manager or `KOTLIN_HOME` is already on the machine, adopt it instead of
downloading another copy with `vfox link --detect kotlin` (vfox 1.1.0 or later).
The locations it searches, how compilers are named and linking by hand are
described in [link.md](link.md).

## Java

The Kotlin compiler is a JVM program. A Java runtime (17 or newer) must already
be on `PATH` before `kotlinc` or `kotlin` can run; the plugin does not install
one. Use the [vfox-java](https://github.com/version-fox/vfox-java) plugin for
that.
