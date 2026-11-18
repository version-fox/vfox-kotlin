# vfox-kotlin

[中文](README_CN.md) | [English](README.md)

[vfox](https://vfox.dev/) 的 Kotlin 插件。

它从 [JetBrains 的发布页](https://github.com/JetBrains/kotlin/releases) 安装
Kotlin 编译器并加入 `PATH`；当 JetBrains 为你的平台提供预编译包时，也会一并安装
Kotlin/Native（`kotlinc-native`）。从 1.9.0 起，安装时会用 JetBrains 同时发布的
`.sha256` 文件校验下载。

## 安装

安装好 [vfox](https://github.com/version-fox/vfox) 后，安装插件和编译器并选中它：

```bash
vfox add kotlin
vfox install kotlin@2.0.0
vfox use --global kotlin@2.0.0
kotlinc -version
```

用 `kotlin@latest` 安装 [Kotlin 发布页](https://kotlinlang.org/docs/releases.html#release-details)
上的最新版本，用 `vfox search kotlin` 列出所有版本。要让 `kotlinc` 和 `kotlin`
出现在 `PATH` 中，需要把 vfox 接入你的 shell，参见其
[快速开始](https://vfox.dev/guides/quick-start.html)。

## 版本

版本号就是 Kotlin 的发布号，例如 `2.0.0` 或 `1.9.24`，没有发行版后缀：

```bash
vfox install kotlin@2.1.0      # 指定版本
vfox install kotlin@latest     # 最新版本
vfox use kotlin@2.1.0          # 仅当前 shell 会话
vfox use --global kotlin@2.1.0 # 所有新 shell
```

详见 [versions.md](docs/zh-CN/versions.md)，包括 1.9.0 之前版本的校验行为。

## Kotlin/Native

Kotlin/Native 作为单独的发布资源提供。从 1.9.0 起，插件会下载适用于 Linux x86_64、
macOS x86_64/arm64 或 Windows x86_64 的预编译编译器，并把它的 `bin` 目录与 JVM 编译器
一起加入 `PATH`：

```bash
kotlinc-native -version
kotlinc-native hello.kt -o hello
```

当 JetBrains 提供 `.sha256` 时，原生归档包会用该文件单独校验。在没有预编译原生
编译器的主机上（例如 Linux arm64）仍可正常安装 JVM 编译器。附带原生编译器需要
vfox 1.0.0 或更新版本，并且和 `kotlinc` 一样需要 `PATH` 上有 Java 运行时。

## 运行要求

`kotlinc` 和 `kotlin` 是 JVM 启动脚本，因此 `PATH` 上必须已有 Java 运行时（17 或更新）。
本插件不安装 Java；如果需要，请使用 [vfox-java](https://github.com/version-fox/vfox-java) 插件。

## 本机已有的编译器

如果机器上已有 Kotlin 编译器——来自 SDKMAN、asdf、mise、Homebrew、IDE 自带、
系统包管理器或 `KOTLIN_HOME`——可以直接纳入而不必再下载一份：

```bash
vfox link --detect kotlin
```

这需要 vfox 1.1.0 或更新版本。vfox 会列出插件找到的编译器并让你选择；
`vfox link --detect -y kotlin` 会不问确认地全部纳入。详见
[link.md](docs/zh-CN/link.md)。

## 文档

**使用插件**

- [版本](docs/zh-CN/versions.md)：版本写法、安装、搜索和选择编译器
- [本机已有的编译器](docs/zh-CN/link.md)：用 `vfox link --detect kotlin` 纳入本机编译器
- [从源码安装](docs/zh-CN/install-from-source.md)：试用本仓库中尚未发布的改动
- [疑难解答](docs/zh-CN/troubleshooting.md)：Java、shell hook、校验和下载问题

**参与开发**

- [开发](docs/zh-CN/development.md)：目录结构、hook 约定、测试和发布
- [E2E](docs/zh-CN/e2e.md)：本地检查以及 Linux/Windows 容器测试
- [更新日志](docs/zh-CN/changes.md)：每个版本的变更

## 发布插件

维护者可以在默认分支上从 **Actions → Plugin → Run workflow** 发布，填入不带 `v`
前缀的稳定版本号。共享工作流会更新 `metadata.lua`、创建版本提交和标签，并在本仓库
发布 ZIP 和 manifest。无需本地打标签或额外的 release token。拉取请求只运行检查。

当 `PLUGIN.version` 与标签一致时，推送版本标签也可以发布。若发布失败，重新运行
原来失败的 job 即可续跑。

工作流使用共享的 `@v1` 发布工具版本。更新该工具不会发布本插件。包结构和首次发布
要求见[共享工作流文档](https://github.com/version-fox/plugin-manifest-action)。

## 许可证

Apache-2.0，见 [LICENSE](LICENSE)。
