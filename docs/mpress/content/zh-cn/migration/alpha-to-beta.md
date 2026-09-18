---
title: "从 v3 Alpha 升级"
description: "将现有 Wails v3 Alpha 项目升级到固定的 Beta 版本"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

本指南适用于现有 v3 Alpha 项目。Wails v2 用户请使用 [v2 到 v3 指南](/migration/v2-to-v3/)。

## 升级之前

提交或备份项目。阅读从当前 Alpha 版本到所选 Beta 版本之间的[更新日志](/changelog/)：可能需要修改源代码、API 或构建配置。请检查[桌面兼容性策略](/status/)和平台要求。

以下命令使用已发布的 `v3.0.0-beta.23` 作为固定版本的示例，并非建议始终跟踪最新版本。如果选择其他版本，请核实其 CLI、Go 模块和 npm 运行时版本，并一起更新命令。在此示例中，npm 版本与去掉 `v` 前缀的 Go 版本相同。

## 1. 更新 CLI

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

确认 `wails3 version` 显示已安装的版本。`PATH` 中位置更靠前的旧可执行文件可能会遮蔽新的 CLI。

## 2. 更新 Go 模块

在项目根目录执行。检查依赖变更，不要批量升级无关模块。

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. 更新前端运行时

对于使用 npm 和 `frontend` 目录的项目：

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

保留锁文件并检查其变更。如果前端使用其他包管理器或目录，请相应调整此步骤，同时保持运行时版本固定。

## 4. 重新生成、构建和测试

在项目根目录，从 Go 服务重新生成绑定并构建：

```sh
wails3 generate bindings
wails3 build
```

运行构建后的应用程序，并在你发布的每个受支持平台上测试工作流程。检查并一起提交源代码、生成的绑定、模块文件和前端锁文件的变更。

## 升级失败时

检查 `PATH` 中的 CLI，使用 `go list -m github.com/wailsapp/wails/v3` 检查模块版本，使用 `npm --prefix frontend ls @wailsio/runtime` 检查已安装的运行时。解决版本不匹配后重新生成绑定。不要假设所有 Alpha 版本都能在不修改代码的情况下升级。

如果问题仍然存在，请[报告可复现的问题](https://github.com/wailsapp/wails/issues/new/choose)，附上新旧版本、准确的错误信息及 `wails3 doctor` 输出。漏洞报告请遵循[安全策略](https://github.com/wailsapp/wails/blob/master/SECURITY.md)。
