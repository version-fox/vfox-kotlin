# 版本

[中文](versions.md) | [English](../versions.md)

版本号就是 Kotlin 的发布号，例如 `2.0.0` 或 `1.9.24`。本插件没有发行版或后缀：
你输入的版本号就是得到的编译器版本。

## 安装

```bash
vfox install kotlin@2.1.0   # 指定版本
vfox install kotlin@latest  # Kotlin 发布页上的最新版本
```

`latest` 由 `Available` hook 解析：它读取
[Kotlin 发布详情](https://kotlinlang.org/docs/releases.html#release-details)表格并取
列出的第一个版本。如果无法访问该页面，`vfox install kotlin@latest` 会失败，而不会猜测版本。

## 列出

```bash
vfox search kotlin
```

`Available` 按页面顺序列出同一张表格中的版本，最新的在前，并标出已安装的版本。
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
| Linux x86_64 | `kotlin-native-prebuilt-linux-x86_64-<version>.tar.gz` | 2.0.0 |
| Linux x86_64 | `kotlin-native-linux-x86_64-<version>.tar.gz` | 1.9.0 |
| macOS x86_64 / arm64 | `kotlin-native-prebuilt-macos-<arch>-<version>.tar.gz` | 2.0.0 |
| macOS x86_64 / arm64 | `kotlin-native-macos-<arch>-<version>.tar.gz` | 1.9.0 |
| Windows x86_64 | `kotlin-native-prebuilt-windows-x86_64-<version>.zip` | 2.0.0 |
| Windows x86_64 | `kotlin-native-windows-x86_64-<version>.zip` | 1.9.0 |

Kotlin 2.0 更改了归档包名称，插件会按版本选择规范名称，并在 JetBrains 提供 `.sha256`
时用它校验。1.9.0 之前的版本没有原生归档包，只会安装 JVM 编译器。Linux arm64 没有
预编译的原生编译器，这些主机同样只安装 JVM 编译器。附带原生编译器需要 vfox 1.0.0
或更新版本。

和 `kotlinc` 一样，`kotlinc-native` 也是 JVM 启动脚本，需要 `PATH` 上有 Java。

## 校验和

从 1.9.0 起，JetBrains 会在每个编译器归档包旁发布 `.sha256` 文件。这些版本的
`PreInstall` 会把摘要交给 vfox，因此安装会显示 `Verifying checksum …`，下载与文件
不匹配时会以 `Checksum error` 失败。

1.9.0 之前的版本没有该文件。它们仍能安装，但 vfox 会打印
`WARNING: Checksum is not provided, skip verify...`。这是预期行为；若旧版本因其他
原因安装失败，参见[疑难解答](troubleshooting.md)。

## 本机已有的编译器

如果机器上已有来自 SDKMAN、asdf、mise、Homebrew、IDE、系统包管理器或
`KOTLIN_HOME` 的编译器，vfox 可以直接纳入：

```bash
vfox link --detect kotlin
```

这需要 vfox 1.1.0 或更新版本（`vfox --version`）。vfox 会显示插件找到的编译器并让你
选择；加上 `-y` 则不问确认地全部纳入。之后 `vfox list kotlin` 会把这些版本标记为
`(linked → <path>)`，`vfox uninstall` 只删除链接，不会移动编译器。

插件通过 `bin/kotlinc` 启动脚本和 `build.txt` 识别编译器目录，并按文件中的版本命名，
因此检测到的编译器写法与安装的完全一致（例如 `2.0.0`）。vfox 自己目录下的条目
永远不会被提供。扫描位置包括：

- SDKMAN（`~/.sdkman/candidates/kotlin`）、asdf 和 mise 的安装目录
- Homebrew 和 Linuxbrew 的 keg，以及系统目录（`/usr/share/kotlinc`、`/opt`、`/usr/local`）
- IntelliJ IDEA 和 Toolbox 自带的编译器，位于 `plugins/kotlin/kotlinc`
- `KOTLIN_HOME`，以及 Windows 上 `Program Files` 和 `USERPROFILE` 下的对应位置

完整列表、命名规则和手工纳入见 [link.md](link.md)。要纳入未扫描位置的编译器，可以
手工指向其主目录（即包含 `bin` 的目录）：

```bash
vfox link kotlin@2.0.0-system /opt/my-kotlin
vfox use --global kotlin@2.0.0-system
```

## Java

Kotlin 编译器是 JVM 程序。运行 `kotlinc` 或 `kotlin` 前，`PATH` 上必须已有 Java
运行时（17 或更新）；本插件不会安装 Java。需要的话请使用
[vfox-java](https://github.com/version-fox/vfox-java) 插件。
