# Install from source

[中文](zh-CN/install-from-source.md) | [English](install-from-source.md)

`vfox add kotlin` installs the latest published plugin. To try a change that is
not released yet, build a plugin ZIP from this checkout and add it from that
file. A plugin is `metadata.lua` plus the `hooks/` and `lib/` directories;
`Injection.lua` is only type hints and may be left out.

## Ubuntu 24.04

```bash
sudo apt install git zip
git clone --depth 1 https://github.com/version-fox/vfox-kotlin.git
cd vfox-kotlin
zip -qr ./kotlin.zip metadata.lua hooks lib Injection.lua
vfox add kotlin --source ./kotlin.zip
rm ./kotlin.zip
```

## Windows 11

In PowerShell 7:

```powershell
winget install --id Git.Git --source winget --exact
git clone --depth 1 https://github.com/version-fox/vfox-kotlin.git
cd ./vfox-kotlin/
tar -a -cf ./kotlin.zip metadata.lua hooks lib Injection.lua
vfox add kotlin --source ./kotlin.zip
Remove-Item ./kotlin.zip
```

## Check it

```bash
vfox install kotlin@latest
vfox use --global kotlin@latest
kotlinc -version
```

Remember that `kotlinc` needs Java on `PATH` (see [versions.md](versions.md#java)).

## When kotlin is already added

`vfox add kotlin` replaces the existing plugin. To keep the published one as
well, give the source build an alias:

```bash
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
vfox install kotlin-dev@latest
```

## Updating

Rebuild the ZIP and add it again. With an alias, update only that plugin so the
published one is left alone:

```bash
vfox remove kotlin-dev
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
```

`vfox update kotlin` compares the version in `metadata.lua` with the installed
one, so it does not replace a source build from the same commit.

To work on the plugin, see [development.md](development.md).
