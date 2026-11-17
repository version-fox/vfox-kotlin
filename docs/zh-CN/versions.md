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
`kotlinc-jvm` 和 `kotlinc-js` 启动脚本就在该目录。要应用这些变量，必须启用 shell
hook，参见 vfox 的[快速开始](https://vfox.dev/guides/quick-start.html)。

## 校验和

从 1.9.0 起，JetBrains 会在每个编译器归档包旁发布 `.sha256` 文件。这些版本的
`PreInstall` 会把摘要交给 vfox，因此安装会显示 `Verifying checksum …`，下载与文件
不匹配时会以 `Checksum error` 失败。

1.9.0 之前的版本没有该文件。它们仍能安装，但 vfox 会打印
`WARNING: Checksum is not provided, skip verify...`。这是预期行为；若旧版本因其他
原因安装失败，参见[疑难解答](troubleshooting.md)。

## Java

Kotlin 编译器是 JVM 程序。运行 `kotlinc` 或 `kotlin` 前，`PATH` 上必须已有 Java
运行时（17 或更新）；本插件不会安装 Java。需要的话请使用
[vfox-java](https://github.com/version-fox/vfox-java) 插件。
