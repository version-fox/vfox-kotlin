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

`vfox add kotlin` refuses with `plugin named kotlin already exists`. You can
replace the plugin, or add the source build next to it.

To replace it, `vfox remove kotlin` comes first. It also removes the compilers
installed through the plugin, so install them again afterwards.

To keep both, give the source build another name:

```bash
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
vfox install kotlin-dev@latest
vfox remove -y kotlin-dev # removes only kotlin-dev and what it installed
```

## Updating

Rebuild the ZIP and add it again. A plugin added under its own name has to be
removed first, which takes its installed compilers with it, so the alias is the
easier route when you update often:

```bash
vfox remove -y kotlin-dev
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
```

`vfox update kotlin` compares the version in `metadata.lua` with the latest
published one, so it leaves a source build of the same version alone.

To work on the plugin, see [development.md](development.md).
