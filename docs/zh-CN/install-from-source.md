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

`vfox add kotlin` 会报错 `plugin named kotlin already exists`。你可以替换原有的插件，
也可以给源码构建的版本另起一个名字，与原插件并存。

替换要先执行 `vfox remove kotlin`，它会连同通过这个插件装的编译器一起删掉，之后需要重新安装。

想两个都留着，就给源码构建的版本起个别名：

```bash
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
vfox install kotlin-dev@latest
vfox remove -y kotlin-dev # 只删除 kotlin-dev 和它装的东西
```

## 更新

重新构建 ZIP 并再次添加。以原名添加的插件必须先删除，这会连已安装的编译器一起删掉，
所以经常更新时用别名更省事：

```bash
vfox remove -y kotlin-dev
vfox add kotlin --source ./kotlin.zip --alias kotlin-dev
```

`vfox update kotlin` 会拿 `metadata.lua` 中的版本和最新的已发布版本比较，因此不会替换
版本相同的源码构建。

参与插件开发请见[开发](development.md)。
