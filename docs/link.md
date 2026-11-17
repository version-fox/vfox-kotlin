# Using compilers that are already on the machine

[中文](zh-CN/link.md) | [English](link.md)

If a Kotlin compiler is already on the machine — from SDKMAN, asdf, mise,
Homebrew, an IDE bundle, the system package manager or `KOTLIN_HOME` —
`vfox link --detect kotlin` adopts it instead of downloading another copy. It
needs vfox 1.1.0 or later (`vfox --version`); older vfox releases ignore the
plugin's `Detect` hook.

```bash
vfox link --detect kotlin
```

vfox lists what the plugin found and links what you pick. Add `-y` to link
everything without asking:

```bash
vfox link --detect -y kotlin
```

Afterwards `vfox use kotlin@...` works as usual: `vfox list kotlin` marks
linked versions `(linked → <path>)`, and `vfox uninstall` removes only the
link, leaving the compiler where it is.

## What it finds

- SDKMAN (`~/.sdkman/candidates/kotlin`), asdf and mise installs
- Homebrew and Linuxbrew kegs, and the system directories
  (`/usr/share/kotlinc`, `/opt`, `/usr/local`)
- IntelliJ IDEA and JetBrains Toolbox bundles, which keep the compiler under
  `plugins/kotlin/kotlinc`
- Scoop apps and the usual Windows locations under `Program Files` and
  `USERPROFILE`
- The path pointed to by `KOTLIN_HOME`

Paths under vfox's own directories are never offered, so a compiler vfox
already installed is not detected again.

## Version names

A found compiler is named the way `vfox install` names it: the version in its
`bin/../build.txt`, such as `1.9.0` or `2.0.0`. A directory counts as a
compiler home when it holds both a `bin/kotlinc` launcher (or `kotlinc.bat`
on Windows) and a `build.txt` or `lib/kotlin-compiler.jar`. If `build.txt` is
missing, the version in the directory name is used. Results are deduplicated
and sorted by version, so the list is stable.

## Linking by hand

Where `--detect` does not look, point at the directory by hand. It is the
compiler home itself, the directory that holds `bin`:

```bash
vfox link kotlin@2.0.0-system /opt/my-kotlin
vfox use --global kotlin@2.0.0-system
```

Remember that `kotlinc` needs a Java runtime on `PATH` (see
[versions.md](versions.md#java)).
