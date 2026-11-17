# Troubleshooting

[中文](zh-CN/troubleshooting.md) | [English](troubleshooting.md)

## `kotlinc: command not found`

The plugin adds `<sdk>/bin` to `PATH`, but vfox's variables are applied by a
shell hook. Run `vfox init` once and restart the shell, or follow the
[quick start](https://vfox.dev/guides/quick-start.html) for your shell. Then
check that a version is selected with `vfox current kotlin` and that the
variables are exported with `vfox env -s bash` (or your shell).

## `kotlinc-native: command not found`

Kotlin/Native is a separate release asset and is only bundled when JetBrains
publishes a prebuilt compiler for your host (Linux x86_64, macOS x86_64/arm64
and Windows x86_64 from 1.9.0 on; there is no Linux arm64 build). If `kotlinc`
works but `kotlinc-native` does not, check `vfox --version`: bundling needs vfox
1.0.0 or later, so upgrade and reinstall the release. Versions before 1.9.0
have no Kotlin/Native archive at all; see [versions.md](versions.md#kotlinnative)
for the full platform list.

## `java: command not found` when running kotlinc

`kotlinc` and `kotlin` are JVM launchers; the plugin does not install Java.
Install a Java runtime (17 or newer) and make sure `java` is on `PATH`, for
example with the [vfox-java](https://github.com/version-fox/vfox-java) plugin.

## `WARNING: Checksum is not provided, skip verify...`

Kotlin releases before 1.9.0 have no `.sha256` asset on GitHub, so vfox cannot
verify them. The install still succeeds. From 1.9.0 on, a successful install
prints `Verifying checksum …` instead; see [versions.md](versions.md#checksums).

## `Checksum error`

The downloaded archive does not match the `.sha256` file. Re-run the install;
if it keeps failing, the download may be corrupted by a proxy. There is no
option to skip verification — choose a different release or network instead.

## `Current version information not detected.`

`PreInstall` sends a HEAD request to the GitHub release asset and this failed,
usually because the version does not exist or GitHub could not be reached.
Check the version with `vfox search kotlin` and retry; `vfox --debug install
kotlin@<version>` shows the request vfox made.

## `vfox install kotlin@latest` fails

`latest` reads the Kotlin release table at kotlinlang.org. If that page is
unavailable, install a version by number instead.

## The compiler is installed but the version is old

`vfox use` selects a version that is already installed. Install the newer
release first (`vfox install kotlin@<version>`), then select it. Use
`vfox list kotlin` to see what is installed.
