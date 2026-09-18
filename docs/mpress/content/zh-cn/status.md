---
title: "项目状态"
description: "Wails v3 Beta 兼容性、安全支持及升级指南"
slug: "status"
sourcePath: "status.md"
---

## 当前状态：Beta

请查看[更新日志](/changelog/)以了解最新状态。

我们的目标是发布稳定的 v3.0 版本。Wails v2 仍是当前稳定版，并继续接收修复。部署前，请使用你的应用程序测试 Beta 版本。

## Beta 兼容性承诺

v3 Beta 兼容性约定涵盖桌面应用程序：

| 平台 | 支持的目标 | 要求和说明 |
| --- | --- | --- |
| Windows | amd64 和 arm64 | WebView2 运行时 |
| macOS | Intel 和 Apple Silicon | 安装指南中注明的 macOS 和 WebKit 版本 |
| Linux | amd64 和 arm64 | 默认使用 GTK4 + WebKitGTK 6.0；GTK3 + WebKit2GTK 4.1在整个 v3.0.x 期间仍作为`-tags gtk3`旧版选项提供，并将在 v3.1 中移除 |

开发所有目标平台的应用程序都需要 Go 1.25 或更高版本。Android 和 iOS 支持仍处于实验阶段，不会阻碍桌面版进入 Beta。Beta API 以稳定为目标，但预发布版本中的缺陷以及明确公布的变更仍可能在 v3.0.0 之前得到修正。

## 如何参与贡献

- 测试最新 Beta 版本并报告可复现的错误
- 为文档和示例作出贡献
- 参与讨论并就 WEP 草案提供反馈
- 提交修复错误、改进文档或实施已接受 WEP 的拉取请求

我们欢迎社区贡献。如果你愿意协助实现这些目标，请参与社区讨论。新功能提案应通过 WEP 拉取请求提交，而不是作为功能请求议题提交。

## 反馈和更新

请将可复现的问题报告为议题；新功能则应通过[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)拉取请求提出。

## 使用 Beta 版本

请固定 CLI、Go 模块和前端运行时的确切版本，而不是跟踪 `latest`。现有 Alpha 项目请遵循 [Alpha 到 Beta 升级指南](/migration/alpha-to-beta/)。

[安全策略](https://github.com/wailsapp/wails/blob/master/SECURITY.md)将 v3 Beta 版本列为受支持版本，Alpha 版本则不受支持。请通过[私密漏洞报告](https://github.com/wailsapp/wails/security/advisories/new)报告漏洞，不要提交公开议题。

## 跟踪中的工作

- [带有 v3 标签的未解决错误](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [带有 P0 或 P1 标签的未解决 v3 议题](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [发布里程碑](https://github.com/wailsapp/wails/milestones)

这些实时查询依赖议题标签，既不是完整列表，也不承诺发布日期或范围。请阅读议题，评估它们对项目的影响。
