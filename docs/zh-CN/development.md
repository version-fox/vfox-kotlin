# 开发

[中文](development.md) | [English](../development.md)

## 目录结构

```
metadata.lua                    插件名称、版本、主页和说明
hooks/available.lua             vfox search 和 kotlin@latest
hooks/pre_install.lua           vfox install：解析 URL 和校验和
hooks/post_install.lua          空操作
hooks/env_keys.lua              把 <sdk>/bin 加入 PATH
hooks/pre_use.lua               空操作
hooks/detect.lua                vfox link --detect：查找本机已有的编译器
lib/util.lua                    下载 URL 模板和版本比较
lib/detect_kotlin.lua           扫描根目录、build.txt 版本和编译器目录判定
tests/*_test.lua                离线 hook 测试，由 vfox plugin test 运行
tests/detect*_test.lua          Detect hook 的离线测试
tests/e2e/                      容器测试，见 e2e.md
Injection.lua                   vfox 注入对象的类型提示，不会被加载
README.md, docs/                英文文档；README_CN.md 和 docs/zh-CN/ 为中文版
opencode.jsonc                  允许本仓库中的 OpenCode 会话直接运行 git push、
                                gh pr create 和 gh pr comment
stylua.toml, .luarc.json        Lua 格式化和语言服务器配置
.licenserc.yaml, .license.tpl   Apache-2.0 头部检查及其文本模板
.devcontainer/                  E2E 套件的 Linux 开发环境
.github/workflows/              CI：发布、测试、lint、许可证和 E2E
```

vfox 加载 `metadata.lua`，在命令需要时调用 hook，并把 `PLUGIN` 和 `RUNTIME`
（`osType`、`archType`、`version`）放入作用域。`osType` 和 `archType` 是 vfox
二进制的 Go `GOOS` 和 `GOARCH`。

## Hook 约定

| 你运行 | Hook |
| --- | --- |
| `vfox search kotlin` | `Available` |
| `vfox install kotlin@2.0.0` | `PreInstall`；随后 vfox 下载、校验并解压 |
| `vfox use kotlin@2.0.0` | `PreUse`（空操作），然后 `EnvKeys` |
| `vfox link --detect kotlin` | `Detect` 返回本机已有的编译器 |
| shell hook | `EnvKeys` 应用于选中的版本 |

- `Available` 获取发布详情页面并解析 `tbody tr` / `td b`，返回
  `{ version, note }` 列表，只包含 `X.Y.Z` 版本：页面顶部的版本线（如 `2.4`）会被跳过。
  请求失败时抛出 `get Kotlin releases failed: HTTP <状态码>`（或连接错误），而不是返回空列表，
  这样 `vfox search` 和 `latest` 会说明为什么没有列出任何版本。页面解析后没有任何版本时才返回空列表。
- `PreInstall` 用版本号两次格式化 `lib/util.lua` 中的 URL。它发送 HEAD 请求，
  除非返回 200，否则抛出 `Current version information not detected.`。对 1.9.0
  及以上的版本，它 GET `<url>.sha256` 并把内容作为 `sha256` 返回，失败时抛出
  `Current version checksum not detected.`；否则不返回校验和。结果中的 `version`
  是 vfox 安装时使用的名称，因此 `latest` 在这里被解析为具体版本；列表为空时，
  `latest` 会抛出 `No Kotlin release is listed, so latest cannot be resolved.`。它还会探测
  `lib/util.lua` 中按平台生成的 Kotlin/Native 资源；
  HEAD 请求成功时，该归档包会作为 `addition` 返回（JetBrains 发布 `.sha256` 时也一并
  返回），由 vfox 作为第二个包下载并解压。缺少原生资源只会使 `kotlinc-native` 不可用，
  编译器仍会正常安装。
- `EnvKeys` 返回 `{ key = "PATH", value = ctx.path .. "/bin" }`，并为 `ctx.sdkInfo` 中
  每个非主包再返回一条 PATH 记录。vfox 会去掉归档包顶层的 `kotlinc` 目录，因此启动
  脚本直接位于 `<sdk>/bin`，而 Kotlin/Native 附加包的 `bin` 中有 `kotlinc-native`、
  `konanc` 和 `klib`。vfox 以运行时目录名（`add-native-<version>`）作为附加包的键，
  而不是 `PreInstall` 中发送的 `name`。
- `PreUse` 和 `PostInstall` 是空操作；vfox 无需插件帮助即可匹配已安装版本并应用
  `EnvKeys`。
- `Detect` 扫描 SDKMAN、asdf、mise、Homebrew、IDE 自带目录、系统目录和 `KOTLIN_HOME`，
  寻找同时包含 `bin/kotlinc` 和 `build.txt` 的编译器目录。它返回
  `{ version, path, note }`，版本取自 `build.txt`（缺失时回退到目录名），因此检测到的
  编译器命名与安装时一致。`.vfox` 下的路径会被跳过，结果去重并排序。该 hook 需要
  vfox 1.1.0；旧版 vfox 会直接忽略它。

## 测试

hook 测试在 vfox 1.0.13 起提供的 `vfox plugin test` 下离线运行。它会查找每个
`tests/**/*_test.lua` 并在独立的 VM 中运行，因此新增测试文件无需任何注册：

```bash
vfox plugin test .
vfox plugin test . --file tests/pre_install_test.lua
```

测试文件先从 `vfox.test` 获取 `test.load()`，然后通过它调用真正的 hook。
`test.load` 接受平台（`os`、`arch`）以及回答每个 HTTP 请求的函数（`http`）。例如：

```lua
local test = require("vfox.test")
local plugin = test.load({
    http = function(request)
        return { status_code = 200, body = "..." }
    end,
})
local result = plugin:PreInstall({ version = "2.0.0" })
```

上下文和结果都经过 vfox 的解码，因此返回类型错误的 hook 会使测试失败。错误路径
使用 `pcall` 调用 hook。

容器测试是 `tests/e2e/` 下使用 testcontainers-go 的 Go 程序。它们安装固定版本的
vfox，把本仓库作为本地插件添加，安装编译器，并在 Linux 和 Windows 上检查结果。
`go test -short` 会跳过它们；运行方式和检查内容见 [e2e.md](e2e.md)。

## 本地检查

lint、格式化、许可证头部和容器测试命令见
[本地检查](e2e.md#本地检查)。

## 发布

维护者在默认分支上从 **Actions → Plugin → Run workflow** 发布，填入不带 `v` 前缀的
稳定版本号。共享工作流
[plugin-manifest-action](https://github.com/version-fox/plugin-manifest-action)
会更新 `metadata.lua`、创建版本提交和标签，并在本仓库发布 ZIP 和 manifest。
无需本地打标签或额外的 release token。拉取请求只运行检查。

当 `PLUGIN.version` 与标签一致时，推送版本标签也可以发布。若发布失败，重新运行
失败的 job 即可续跑。工作流使用共享的 `@v1` 发布工具，更新该工具不会发布本插件。
之后 [vfox-plugins](https://github.com/version-fox/vfox-plugins) 中的
`plugins/kotlin.json` 会显示 `vfox add kotlin` 提供的版本。

### 发布之后

在 [changes.md](changes.md) 中把比上一个版本新的改动标为“未发布”，发布后连同日期
移到新版本下：

```bash
grep -rn -i -e unreleased -e 未发布 docs README*
```

## 编写文档

英文页面在 `docs/`，中文页面在 `docs/zh-CN/` 下同名存放（中文 README 是
`README_CN.md`）。每页开头都有指向另一种语言的链接，新增页面要加入两个 README 的
文档列表。CI 不检查正文，因此修改页面时请同时阅读两种语言。
