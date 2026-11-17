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

`latest` is resolved by the `Available` hook: it reads the
[Kotlin release details](https://kotlinlang.org/docs/releases.html#release-details)
table and takes the first version listed. If that page cannot be reached,
`vfox install kotlin@latest` fails rather than guessing a version.

## Listing

```bash
vfox search kotlin
```

`Available` lists the releases from the same table in page order, newest first,
and marks releases that are already installed. The table is upstream's, so a
release appears here as soon as JetBrains adds it.

## Selecting

```bash
vfox use kotlin@2.1.0          # current shell session
vfox use --project kotlin@2.1.0 # this project
vfox use --global kotlin@2.1.0 # every new shell
vfox current kotlin            # the selected version
```

`EnvKeys` adds `<sdk>/bin` to `PATH`, which is where the `kotlinc`, `kotlin`,
`kotlinc-jvm` and `kotlinc-js` launchers live inside the installed archive.
A shell hook must be active for the variables to be applied; see vfox's
[quick start](https://vfox.dev/guides/quick-start.html).

## Checksums

JetBrains publishes a `.sha256` file beside each compiler archive from 1.9.0 on.
For those releases `PreInstall` returns the digest to vfox, so the install ends
with `Verifying checksum …` and a download that does not match the file fails
with `Checksum error`.

Releases before 1.9.0 have no such asset. They still install, but vfox prints
`WARNING: Checksum is not provided, skip verify...`. This is expected; see
[troubleshooting.md](troubleshooting.md) if an older release fails for another
reason.

## Already installed compilers

If a compiler from SDKMAN, asdf, mise, Homebrew, an IDE, the system package
manager or `KOTLIN_HOME` is already on the machine, vfox can adopt it:

```bash
vfox link --detect kotlin
```

This needs vfox 1.1.0 or later (`vfox --version`). vfox shows what the plugin
found and links what you pick; add `-y` to link everything without asking.
`vfox list kotlin` then marks those versions `(linked → <path>)` and
`vfox uninstall` removes only the link, leaving the compiler where it is.

The plugin recognises a compiler home by its `bin/kotlinc` launcher and its
`build.txt`, and names it after the version in that file, so a detected
compiler is spelled just like an installed one (for example `2.0.0`). Entries
under vfox's own directories are never offered. It looks in:

- SDKMAN (`~/.sdkman/candidates/kotlin`), asdf and mise installs
- Homebrew and Linuxbrew kegs, and the system directories (`/usr/share/kotlinc`, `/opt`, `/usr/local`)
- IntelliJ IDEA and Toolbox bundles, which keep the compiler under `plugins/kotlin/kotlinc`
- `KOTLIN_HOME`, and the equivalent Windows locations under `Program Files` and `USERPROFILE`

To adopt a compiler in a location that is not scanned, point at its home
directory by hand — the one that holds `bin`:

```bash
vfox link kotlin@2.0.0-system /opt/my-kotlin
vfox use --global kotlin@2.0.0-system
```

## Java

The Kotlin compiler is a JVM program. A Java runtime (17 or newer) must already
be on `PATH` before `kotlinc` or `kotlin` can run; the plugin does not install
one. Use the [vfox-java](https://github.com/version-fox/vfox-java) plugin for
that.
