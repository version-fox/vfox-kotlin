# 从源码安装

[中文](install-from-source.md) | [English](../install-from-source.md)

`vfox add kotlin` 安装最新发布的插件。要试用尚未提交的改动，可以从本仓库构建插件
ZIP 并从该文件添加。插件由 `metadata.lua`、`hooks/` 和 `lib/` 目录组成；
`Injection.lua` 只是类型提示，可以省略。

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

在 PowerShell 7 中：

```powershell
winget install --id Git.Git --source winget --exact
git clone --depth 1 https://github.com/version-fox/vfox-kotlin.git
cd ./vfox-kotlin/
tar -a -cf ./kotlin.zip metadata.lua hooks lib Injection.lua
vfox add kotlin --source ./kotlin.zip
Remove-Item ./kotlin.zip
```

## 验证

```bash
vfox install kotlin@latest
vfox use --global kotlin@latest
kotlinc -version
```

请记住 `kotlinc` 需要 `PATH` 上有 Java（见[版本](versions.md#java)）。

## 已添加 kotlin 时

`vfox add kotlin` 会替换现有插件。要同时保留已发布版本，可以给源码构建指定别名：

```bash
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
vfox install kotlin-dev@latest
```

## 更新

重新构建 ZIP 并再次添加。使用别名时只更新该插件，已发布版本不受影响：

```bash
vfox remove kotlin-dev
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
```

`vfox update kotlin` 会比较 `metadata.lua` 中的版本和已安装版本，因此不会替换来自
同一提交的源码构建。

参与插件开发请见[开发](development.md)。
