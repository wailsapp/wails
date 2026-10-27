---
title: "路线图"
description: "Wails v3 项目状态、计划功能及参与贡献的方式"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — 稳定版

查看[变更日志](/changelog/)，了解最新状态。

本文档涵盖 Wails v3.0.0。v3.0.0 发布后，修复和改进将通过每周版本发布；两次发布之间还会提供带日期的夜间预览版，并针对紧急问题发布热修复。请参阅[发布政策](/contributing/release-policy/)。

## 桌面端兼容性

v3.0.0兼容性约定涵盖桌面应用程序：

| 平台 | 支持的目标 | 要求和说明 |
| --- | --- | --- |
| Windows | amd64 和 arm64 | WebView2 运行时 |
| macOS | Intel 和 Apple Silicon | 安装指南中所列的 macOS 和 WebKit 版本 |
| Linux | amd64 和 arm64 | 默认使用 GTK4 + WebKitGTK 6.0；GTK3 + WebKit2GTK 4.1在v3.0.x 期间仍作为`-tags gtk3`旧版选项提供，并在v3.1中移除 |

开发所有目标均需要 Go 1.26或更高版本。Android 和 iOS 支持仍处于实验阶段，不属于稳定版桌面端兼容性约定。

## 如何参与贡献

- 测试 Wails v3.0.0并报告可复现的错误
- 为文档和示例作出贡献
- 参与讨论并就 WEP 草案提供反馈
- 提交用于错误修复、文档改进或已获接受的 WEP 的拉取请求

我们欢迎社区贡献。如果你愿意协助实现这些目标，请参与社区讨论。新功能提案应通过 WEP PR 提交，而不是作为功能请求议题提交。

## 反馈与更新

本路线图可能会根据社区反馈和项目优先级进行调整。我们将定期更新路线图，以反映进展和方向变化。请将可复现的问题报告为议题；新功能请通过[WEP（Wails 增强提案）](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md) PR 提出。
