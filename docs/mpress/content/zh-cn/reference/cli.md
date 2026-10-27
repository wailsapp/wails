---
title: "CLI 命令概览"
description: "按用途查找 Wails CLI 命令：项目、构建、签名、工具和服务"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

从这里开始按用途查找命令。有关特定命令的标志和用法示例，请参阅[CLI 标志和示例](/guides/cli/)。

## 概览

Wails CLI（`wails3`）是创建、开发、构建、签名、打包和检查 Wails 3应用程序的命令行入口。大多数构建编排工作都委托给各项目的 Taskfile（位于项目的`build/`下）——许多`wails3`命令只是调用特定任务的轻量封装。

要获取任何命令的最新帮助，请运行：

```bash
wails3 --help
wails3 <command> --help
```

## 项目生命周期

| 命令 | 说明 |
| --- | --- |
| `wails3 init` | 根据模板创建新项目。标志：`-n`（项目名称）、`-t`（模板，默认为`vanilla`）、`-p`（Go 包名，默认为`main`）、`-d`（项目目录，默认为`.`）、`-q`（静默模式）、`-l`（列出模板）、`-mod`（Go 模块路径）、`--git`（Git 仓库 URL）、`--skipgomodtidy`、`-s`（跳过远程模板警告）、`--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`。 |
| `wails3 dev` | 以开发模式运行应用程序，并启用前端热重载。标志：`--config`（默认为`./build/config.yml`）、`--port`（Vite 开发端口）、`-s`（启用 HTTPS）。 |
| `wails3 build` | 构建项目。这是对 Taskfile 中`build`任务的轻量封装。标志：`--tags`（以`EXTRA_TAGS=`转发）、`--obfuscated`（使用 Garble 构建；请参阅[混淆构建](/guides/build/obfuscation/)）、`--garbleargs`（在`build`子命令之前转发给`garble`的额外标志）。 |
| `wails3 package` | 运行特定于平台的`package` Taskfile 任务。 |
| `wails3 task [name]` | 运行任意 Taskfile 任务；未指定名称时，`--list`会显示所有已注册的任务。 |
| `wails3 mcp` | 启动项目 MCP 服务器。对于由代理启动的进程，会自动使用 stdio；在交互式终端中使用时，则自动使用基于环回地址的 Streamable HTTP。 |
| `wails3 doctor` | 输出环境诊断报告。 |
| `wails3 doctor-ng` | `doctor`的新版 TUI 变体。 |
| `wails3 version` | 输出 CLI 版本。 |
| `wails3 releasenotes` | 输出近期的发行说明。 |
| `wails3 docs` | 在浏览器中打开文档网站。 |
| `wails3 sponsor` | 打开赞助页面。 |

## 生成

`wails3 generate <subcommand>`:

| 子命令 | 说明 |
| --- | --- |
| `generate bindings` | 生成从 Go 到前端的绑定。标志：`-d`（输出目录）、`-models`、`-index`、`-ts`、`-i`（接口）、`-b`（捆绑）、`-names`（生成`Call.ByName`）、`-noevents`、`-noindex`、`-dry`、`-silent`、`-v`、`-clean`（默认为`true`）、`-f`、`-obfuscated`（生成带有稳定绑定 ID 的`wails_obfuscated.gen.go`，用于 Garble 构建；请参阅[混淆构建](/guides/build/obfuscation/)）、`-obfuscated-output`（生成文件的目录；默认为 main 包目录）。接受包模式（例如`./...`）；如果未提供，则回退到当前目录。 |
| `generate icons` | 将源 PNG 转换为各平台的图标格式。标志：`-input`、`-windowsfilename`、`-macfilename`、`-iconcomposerinput`、`-macassetdir`。 |
| `generate build-assets` | 根据`build/config.yml`生成`build/`目录的内容（Taskfile 片段、NSIS 文件、`Info.plist`、`.desktop`模板等）。 |
| `generate runtime` | 重新生成随 webview 一同提供的预构建`/wails/runtime.js`。 |
| `generate syso` | 生成 Windows `.syso`资源文件（图标、清单和版本信息）。 |
| `generate webview2bootstrapper` | 生成适用于 Windows 的 WebView2 引导安装程序。 |
| `generate constants` | 根据 Go 事件类型生成 JS 事件名称常量。 |
| `generate template` | 搭建新的项目模板框架。 |
| `generate .desktop` | 生成 Linux `.desktop`文件（供 AppImage/DEB/RPM 使用）。 |
| `generate appimage` | 生成 AppImage 构建目录。 |

## 更新

`wails3 update <subcommand>`:

| 子命令 | 说明 |
| --- | --- |
| `update build-assets` | 根据`build/config.yml`刷新`build/`目录（尽可能保留用户所做的编辑）。 |
| `update cli` | 对`wails3`二进制文件执行自更新。 |

## 代码签名和打包

| 命令 | 说明 |
| --- | --- |
| `wails3 signing` | 交互式向导，为其在`build/`中检测到的平台配置签名。标志：`--platform`（可重复指定；默认根据构建目录自动检测）。 |
| `wails3 entitlements` | 用于配置 macOS 权利声明的交互式向导。标志：`--output`（路径；默认为`build/darwin/entitlements.plist`）。 |
| `wails3 sign [GOOS=…]` | 封装命令，用于运行当前操作系统（或通过`GOOS`指定的操作系统）所对应的`*:sign` Taskfile 任务。 |
| `wails3 tool sign` | 底层直接签名入口。标志：`--input`、`--output`、`--verbose`、`--certificate`、`--password`、`--thumbprint`、`--timestamp`、`--identity`、`--entitlements`、`--hardened-runtime`、`--notarize`、`--keychain-profile`、`--pgp-key`、`--pgp-password`、`--role`。 |

在此版本中，CLI 解析器在顶层公开了`signing`和`entitlements`。请查看`wails3 signing --help`和`wails3 entitlements --help`；`wails3 setup signing`不会选择签名向导。若要直接管理凭据，请使用`xcrun notarytool store-credentials`或`gpg`。

## 工具

`wails3 tool <subcommand>`:

| 子命令 | 说明 |
| --- | --- |
| `tool checkport` | 检查 TCP 端口是否开放（适用于等待 Vite）。 |
| `tool watcher` | 每当监视的文件发生更改时运行命令。 |
| `tool cp` | 跨平台复制文件。 |
| `tool buildinfo` | 输出二进制文件中嵌入的 Go 构建信息。 |
| `tool package` | 根据`build/linux/nfpm`构建 Linux 软件包（`deb`、`rpm`、`archlinux`）。 |
| `tool version` | 递增项目的语义化版本号。 |
| `tool lipo` | 将多个 macOS 架构的二进制文件合并为通用二进制文件。 |
| `tool capabilities` | 探测系统是否支持 GTK3/GTK4/WebKit。 |
| `tool sign` | （请参阅[代码签名与打包](#heading-4)。） |

## 服务

`wails3 service <subcommand>`:

| 子命令 | 说明 |
| --- | --- |
| `service init` | 搭建新的服务包框架。 |

## iOS

`wails3 ios <subcommand>`:

| 子命令 | 说明 |
| --- | --- |
| `ios overlay:gen` | 为 iOS 桥接适配层生成 Go overlay。 |
| `ios xcode:gen` | 在输出目录中生成 Xcode 项目。 |

## 构建输出路径

- 原生二进制文件位于`bin/<APP_NAME>`（在 Windows 上则位于`bin/<APP_NAME>.exe`）。不存在`build/bin/`。
- 打包产物（`.app`、`.dmg`、NSIS 安装程序、MSIX、DEB/RPM/AppImage）同样位于`bin/`，或位于相关 Taskfile 任务创建的特定于平台的子目录中。

## 通用标志

| 标志 | 适用于 | 说明 |
| --- | --- | --- |
| `-nocolour` | 在`--help`中公开此标志的命令（例如`init`和`update cli`） | 在 CLI 输出中禁用 ANSI 颜色。并非所有命令都接受此标志，其中包括`version`。 |

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[示例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)。
