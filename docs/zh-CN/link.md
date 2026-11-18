# 使用本机已有的编译器

[中文](link.md) | [English](../link.md)

如果机器上已有 Kotlin 编译器——来自 SDKMAN、asdf、mise、Homebrew、IDE 自带目录、
系统包管理器或 `KOTLIN_HOME`——`vfox link --detect kotlin` 可以直接纳入，而不必再
下载一份。这需要 vfox 1.1.0 或更新版本（`vfox --version`）；旧版 vfox 会忽略插件的
`Detect` hook。

```bash
vfox link --detect kotlin
```

vfox 会列出插件找到的编译器并让你选择。加上 `-y` 则不问确认地全部纳入：

```bash
vfox link --detect -y kotlin
```

之后 `vfox use kotlin@...` 照常工作：`vfox list kotlin` 会把纳入的版本标记为
`(linked → <path>)`，`vfox uninstall` 只删除链接，编译器本身不会被移动或删除。

## 它能找到什么

- SDKMAN（`~/.sdkman/candidates/kotlin`）、asdf 和 mise 的安装目录
- Homebrew 和 Linuxbrew 的 keg，以及系统目录
  （`/usr/share/kotlinc`、`/opt`、`/usr/local`）
- IntelliJ IDEA 和 JetBrains Toolbox 自带的编译器，位于 `plugins/kotlin/kotlinc`
- Scoop 应用，以及 Windows 上 `Program Files` 和 `USERPROFILE` 下的常见位置
- `KOTLIN_HOME` 指向的路径

vfox 自己目录下的路径永远不会被提供，因此已由 vfox 安装的编译器不会被重复检测。

## 版本命名

找到的编译器按 `vfox install` 的方式命名：取自 `build.txt` 的版本，例如 `1.9.0`
或 `2.0.0`。一个目录只有同时包含 `bin/kotlinc` 启动脚本（Windows 上是
`kotlinc.bat`）和 `build.txt` 或 `lib/kotlin-compiler.jar` 时才算编译器目录。
如果没有 `build.txt`，则使用以版本号开头的目录名（如 `2.0.0` 或 `2.0.0-system`）中的版本；两者都没有的目录会被跳过。结果会去重并按版本排序，因此列表是稳定的。

## 手工纳入

若编译器不在 `--detect` 的扫描范围内，可以手工指向其主目录，也就是包含 `bin` 的目录：

```bash
vfox link kotlin@2.0.0-system /opt/my-kotlin
vfox use --global kotlin@2.0.0-system
```

请记住 `kotlinc` 需要 `PATH` 上有 Java 运行时（见[版本](versions.md#java)）。
