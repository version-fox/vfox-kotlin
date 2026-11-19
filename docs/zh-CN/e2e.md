# E2E

[中文](e2e.md) | [English](../e2e.md)

离线 Lua hook 测试（`tests/*_test.lua`）在 **Test Plugin** 工作流中运行。端到端套件
则在容器中对真实的 vfox 运行，覆盖 Linux 和 Windows；每个用例使用独立的临时容器。
CI 在标准 runner 上原生运行（Docker 加上 `actions/setup-go` 安装的 Go）。
`.devcontainer/` 中的开发容器使用相同固定版本的工具，在本地复现 Linux 流程。

## 环境变量

| 变量                 | 默认值   | 含义                                                     |
| -------------------- | -------- | -------------------------------------------------------- |
| `KOTLIN_VERSION`     | `2.0.0`  | 被测版本，它的用例执行全部步骤（1.9.0-RC 起带校验和）         |
| `E2E_MAX_JOBS`       | `3`      | 一次 `go test` 运行中同时运行的容器（用例）数              |
| `E2E_RETRY_ATTEMPTS` | `3`      | 每个涉及网络的步骤的重试次数（线性退避 10 秒）            |

每个平台运行两个用例，各自一个容器：

- 被测版本（`KOTLIN_VERSION`）执行下文的全部步骤；
- `1.8.22`（最后一个没有 `.sha256` 资产的稳定版本）只执行步骤 1 和 3。若它本身就是被测版本，则不单独运行。

默认 `E2E_MAX_JOBS=3` 时两个用例同时运行；`E2E_MAX_JOBS=1` 则依次运行。

## Linux（Ubuntu）

主机上装有 Docker Engine 和 Go 1.27.1 时，在仓库根目录运行：

```bash
cd tests/e2e
go test -v -timeout 45m -run TestE2E ./linux
```

rootless Docker 用户需要用
`export DOCKER_HOST="unix://$XDG_RUNTIME_DIR/docker.sock"` 指向自己的 socket。
迭代单个版本时：

```bash
KOTLIN_VERSION=1.8.22 go test -v -timeout 45m -run TestE2E ./linux
```

套件会构建 `tests/e2e/linux/Dockerfile`（Ubuntu 24.04、JRE、curl 和 zip），安装
带 SHA-256 校验的固定版本 vfox，从 ZIP 添加本地插件，然后通过 `bash` exec 调用
驱动安装与验证。

## Windows 11

Windows 容器需要处于 Windows 容器模式的 Windows Docker 守护进程。CI 在
`windows-latest` 上运行：若 `docker` 服务已停止，job 会先启动它；守护进程若不处于
Windows 容器模式，job 会失败。
在装有 Docker 和 Go 的 Windows 主机上：

```powershell
cd tests/e2e
go test -v -timeout 55m -run TestE2E ./windows
```

镜像（`tests/e2e/windows/Dockerfile`）基于 servercore，带有 PowerShell 7 和
Temurin 17 JRE。Go 测试通过 `pwsh -NoProfile -Command` 驱动它，使用与 Linux 相同
的固定 vfox 归档和本地插件 ZIP。每条命令都在包含空格的 `VFOX_HOME` 路径下运行，
这是 vfox 自 1.0.13 起支持的情况。

## 开发容器

`.devcontainer/` 打开一个 Linux 环境，带有工作流使用的固定版本 Go、stylua、
lua-language-server、许可证工具和 vfox。它挂载宿主 Docker socket，使套件生成与 CI
相同的目标容器（Docker-out-of-Docker）。可用 Dev Containers 工具打开，例如：

```bash
devcontainer up --workspace-folder .
devcontainer exec --workspace-folder . bash
```

Windows 容器无法在 Linux 开发容器中运行，因此 Windows 套件需要按上文所述使用
Windows 主机。

## 本地检查

在开发容器中执行（所有工具均已预装并固定到工作流使用的版本），或在主机上安装相同工具。

```bash
# Lua 和 Go 检查（.github/workflows/lint.yml）。
cd tests/e2e
go test -short ./...
test -z "$(gofmt -l .)"
go vet ./...
cd ../..

stylua --check .
lua-language-server --check . --checklevel=Warning
actionlint
```

```bash
# 离线 hook 测试（.github/workflows/test.yml）。
vfox plugin test .
```

```bash
# 许可证头部（.github/workflows/license.yml）。
docker run --rm -v "$PWD:/work" -w /work apache/skywalking-eyes:0.9.0 header check
docker run --rm -v "$PWD:/work" -w /work apache/skywalking-eyes:0.9.0 header fix
```

## 套件验证什么

两个平台执行相同的步骤。每个用例都执行步骤 1 和 3；被测版本还执行步骤 2 和 4 至 7：

1. 从经过校验和检查的固定归档安装 vfox，并从本仓库注册本地 `kotlin` 插件——不是
   从插件注册表。
2. `vfox search kotlin` 列出已安装的版本。
3. `vfox install kotlin@<version>` 以该名称安装。1.9.0-RC 及以上版本，vfox 会打印
   `Verifying checksum <64 位十六进制>…`；1.8.22、1.9.0-Beta 等其他版本则报告未提供校验和。
   2.0.0 覆盖前一种情况，1.8.22 覆盖后一种情况。
4. 安装 JetBrains 从未发布的版本会在 `PreInstall` 的 HEAD 检查中失败。
5. `vfox use --global` 选中编译器。在新的 shell 中，SDK 的 `bin` 位于 `PATH` 中，
   `vfox env` 导出它，`kotlinc -version` 和 `kotlin -version` 报告请求的版本，
   `vfox current kotlin` 给出选中的版本。
6. 在 amd64 上，对 1.9.0-Beta 及以上的版本，`kotlinc-native` 位于 `PATH` 中，并且
   `kotlinc-native -version` 报告相同的版本。
7. `vfox link --detect kotlin` 找到被复制到 vfox 目录之外并由 `KOTLIN_HOME` 指向的
   编译器，将其纳入，链接会出现在 `vfox list kotlin` 中并标为 `(linked → …)`。
   固定使用的 vfox 1.1.0 是首个支持 `Detect` hook 的版本。
