# 版本

[中文](versions.md) | [English](../versions.md)

版本号就是 Kotlin 的发布号，例如 `2.0.0` 或 `1.9.24`。本插件没有发行版或后缀：
你输入的版本号就是得到的编译器版本。

## 安装

```bash
vfox install kotlin@2.1.0   # 指定版本
vfox install kotlin@latest  # Kotlin 发布页上的最新版本
```

`latest` 由 `PreInstall` 解析：它取 `Available` hook 从
[Kotlin 发布详情](https://kotlinlang.org/docs/releases.html#release-details)表格中列出的
第一个 `X.Y.Z` 发布版本。如果无法访问该页面，或该页面未列出任何版本，`vfox install kotlin@latest`
会失败，而不会猜测版本。

## 列出

```bash
vfox search kotlin
```

`Available` 按页面顺序列出同一张表格中的 `X.Y.Z` 发布版本，最新的在前，并标出已安装的版本。
页面顶部的版本线（如 `2.4`）不是发布版本，会被跳过。
表格来自上游，因此 JetBrains 一发布新版本，它就会出现在这里。

## 选择

```bash
vfox use kotlin@2.1.0           # 当前 shell 会话
vfox use --project kotlin@2.1.0 # 当前项目
vfox use --global kotlin@2.1.0  # 所有新 shell
vfox current kotlin             # 当前选中的版本
```

`EnvKeys` 把 `<sdk>/bin` 加入 `PATH`，归档包中的 `kotlinc`、`kotlin`、
`kotlinc-jvm` 和 `kotlinc-js` 启动脚本就在该目录；如果安装了 Kotlin/Native，它的
`bin` 目录也会加入。要应用这些变量，必须启用 shell hook，参见 vfox 的
[快速开始](https://vfox.dev/guides/quick-start.html)。

## Kotlin/Native

Kotlin/Native 不在 `kotlin-compiler` 归档包中，而是 JetBrains 单独发布的资源。
只要主机有对应的预编译编译器，`PreInstall` 就会把它作为附加包一起下载，因此一次
`vfox install` 就能同时获得两套工具链：

```bash
kotlinc-native -version
kotlinc-native hello.kt -o hello
```

| 主机 | 原生归档包 | 起始版本 |
| --- | --- | --- |
| Linux x86_64 | `kotlin-native-prebuilt-linux-x86_64-<version>.tar.gz` | 2.0.0-Beta3 |
| Linux x86_64 | `kotlin-native-linux-x86_64-<version>.tar.gz` | 1.9.0 |
| macOS x86_64 | `kotlin-native-prebuilt-macos-x86_64-<version>.tar.gz` | 2.0.0-Beta3 |
| macOS arm64 | `kotlin-native-prebuilt-macos-aarch64-<version>.tar.gz` | 2.0.0-Beta3 |
| macOS x86_64 | `kotlin-native-macos-x86_64-<version>.tar.gz` | 1.9.0 |
| macOS arm64 | `kotlin-native-macos-aarch64-<version>.tar.gz` | 1.9.0 |
| Windows x86_64 | `kotlin-native-prebuilt-windows-x86_64-<version>.zip` | 2.0.0-Beta3 |
| Windows x86_64 | `kotlin-native-windows-x86_64-<version>.zip` | 1.9.0 |

Kotlin 2.0.0-Beta3 更改了归档包名称（2.0.0-Beta1 和 2.0.0-Beta2 仍使用旧名称），插件会按版本选择
规范名称，并在 JetBrains 提供 `.sha256` 时用它校验。1.9.0 之前的版本没有原生归档包，只会安装 JVM 编译器。Linux arm64 没有
预编译的原生编译器，这些主机同样只安装 JVM 编译器。附带原生编译器需要 vfox 1.0.0
或更新版本。

和 `kotlinc` 一样，`kotlinc-native` 也是 JVM 启动脚本，需要 `PATH` 上有 Java。

## 校验和

从 1.9.0-RC 起，JetBrains 会在每个编译器归档包旁发布 `.sha256` 文件（1.9.0-Beta 没有）。这些版本的
`PreInstall` 会把摘要交给 vfox，因此安装会显示 `Verifying checksum …`，下载与文件
不匹配时会以 `Checksum error` 失败。

1.9.0-RC 之前的版本没有该文件。它们仍能安装，但 vfox 会打印
`WARNING: Checksum is not provided, skip verify...`。这是预期行为；若旧版本因其他
原因安装失败，参见[疑难解答](troubleshooting.md)。

## 本机已有的编译器

机器上若已有来自 SDKMAN、asdf、mise、Homebrew、IDE、系统包管理器或 `KOTLIN_HOME` 的编译器，
可以用 `vfox link --detect kotlin` 直接纳入，而不必再下载一份（需要 vfox 1.1.0 或更新版本）。
插件会扫描哪些位置、如何命名，以及手工纳入的方法，参见 [link.md](link.md)。

## Java

Kotlin 编译器是 JVM 程序。运行 `kotlinc` 或 `kotlin` 前，`PATH` 上必须已有 Java
运行时（17 或更新）；本插件不会安装 Java。需要的话请使用
[vfox-java](https://github.com/version-fox/vfox-java) 插件。
