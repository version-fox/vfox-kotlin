# 疑难解答

[中文](troubleshooting.md) | [English](../troubleshooting.md)

## `kotlinc: command not found`

插件会把 `<sdk>/bin` 加入 `PATH`，但 vfox 的变量由 shell hook 应用。请用 `vfox activate`
把 vfox 接入 shell：bash 下在 `~/.bashrc` 中加入 `eval "$(vfox activate bash)"`，然后重启 shell。
其他 shell 请参照[快速开始](https://vfox.dev/guides/quick-start.html)。然后用
`vfox current kotlin` 检查是否已选中版本，用 `vfox env -s bash`（或你的 shell）查看导出的变量。

## `kotlinc-native: command not found`

Kotlin/Native 是单独的发布资源，只有当 JetBrains 为主机提供预编译编译器时才会附带
（从 1.9.0 起支持 Linux x86_64、macOS x86_64/arm64 和 Windows x86_64；没有 Linux
arm64 版本）。如果 `kotlinc` 可用但 `kotlinc-native` 不可用，请用 `vfox --version`
检查版本：附带原生编译器需要 vfox 1.0.0 或更新版本，升级后请重新安装该版本。
1.9.0 之前的版本根本没有 Kotlin/Native 归档包；完整平台列表见
[versions.md](versions.md#kotlinnative)。

## 运行 kotlinc 时出现 `java: command not found`

`kotlinc` 和 `kotlin` 是 JVM 启动脚本；本插件不安装 Java。请安装 Java 运行时
（17 或更新）并确保 `java` 在 `PATH` 中，例如使用
[vfox-java](https://github.com/version-fox/vfox-java) 插件。

## `WARNING: Checksum is not provided, skip verify...`

1.9.0 之前的 Kotlin 版本在 GitHub 上没有 `.sha256` 文件，因此 vfox 无法校验。
安装仍会成功。从 1.9.0 起，安装成功时会显示 `Verifying checksum …`；
详见[版本](versions.md#校验和)。

## `Checksum error`

下载到的归档包与 `.sha256` 文件不一致。重新安装试试；如果持续失败，可能是代理
损坏了下载。没有跳过校验的选项——请换一个版本或换网络。

## `Current version information not detected.`

`PreInstall` 向 GitHub 发布资源发送 HEAD 请求，但请求失败，通常是版本不存在或
无法访问 GitHub。请用 `vfox search kotlin` 检查版本后重试；
`vfox --debug install kotlin@<version>` 会显示 vfox 发出的请求。

## `Current version checksum not detected.`

从 1.9.0 起，`PreInstall` 还会下载 JetBrains 在每个归档包旁发布的 `.sha256` 文件。下载失败时
安装会以此错误停止。通常是网络问题，或代理拦截了 GitHub 的访问；请重试，或用
`vfox --debug install kotlin@<version>` 查看请求。

## `vfox install kotlin@latest` 失败

`latest` 会读取 kotlinlang.org 上的 Kotlin 发布表格。如果无法读取该页面，vfox 会
以 `get Kotlin releases failed: HTTP <状态码>`（或连接错误）停止。如果页面可以读取，但没有
列出任何 `X.Y.Z` 发布版本，vfox 会以 `No Kotlin release is listed, so latest cannot be resolved.`
停止。这两种情况都请改用具体版本号安装。

## 编译器已安装但版本较旧

`vfox use` 选择的是已安装的版本。请先安装新版本
（`vfox install kotlin@<version>`），再选中它。用 `vfox list kotlin` 查看已安装的版本。
