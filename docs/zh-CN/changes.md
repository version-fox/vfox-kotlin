# 更新日志

[中文](changes.md) | [English](../changes.md)

每个版本的用户可感知变更。GitHub 发布页列出了拉取请求。

## 未发布

- 中英文 README 和文档，涵盖版本、从源码安装、开发、E2E、疑难解答以及本更新日志。
- `Available`、`PreInstall`、`EnvKeys` 和版本辅助函数的离线 hook 测试，由 **Test Plugin** 工作流中的 `vfox plugin test` 运行。
- Linux（Ubuntu 容器）和 Windows（Windows 容器）的 E2E 套件，各有对应的 GitHub Actions 工作流。它们安装固定版本的 vfox，从本仓库添加插件，安装编译器并检查 PATH 和版本输出。
- Go、Lua 和工作流的 lint 与格式检查（`gofmt`、`go vet`、stylua、lua-language-server、actionlint），工具版本固定，并提供 `stylua.toml` / `.luarc.json`。
- 使用 Apache SkyWalking Eyes 进行 Apache-2.0 头部检查，提供 `.licenserc.yaml`、`.license.tpl` 并为每个源文件添加头部。
- `opencode.jsonc` 允许本仓库中的 OpenCode 会话运行 `git push`、`gh pr create` 和 `gh pr comment`。
- `vfox link --detect kotlin` 纳入本机已有的 Kotlin 编译器——来自 SDKMAN、asdf、mise、Homebrew、IDE 自带目录、系统目录和 `KOTLIN_HOME`。编译器目录通过 `bin/kotlinc` 和 `build.txt` 识别，并按文件中的版本命名。需要 vfox 1.1.0 或更新版本；详见[版本](versions.md#本机已有的编译器)。
- 仓库规范：`.gitattributes`、`.dockerignore`、扩充的 `.gitignore` 以及 Linux 开发容器。测试和 E2E 中固定的 vfox 版本现为 1.1.0。

## 0.1.1（2026-09-12）

- 发布迁移到共享的 `plugin-manifest-action` 工作流。维护者从 **Actions → Plugin → Run workflow** 发布；拉取请求只运行检查。

## 0.1.0（2024-03-28）

- 首个版本：从 JetBrains 的 GitHub 发布页安装 Kotlin 编译器，从 Kotlin 发布页解析 `latest`，从 1.9.0 起使用 `.sha256` 文件校验，并把 `<sdk>/bin` 加入 `PATH`。
