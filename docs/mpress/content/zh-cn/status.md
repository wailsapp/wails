---
title: "路线图"
description: "Wails v3 项目状态、计划功能以及参与贡献的方式"
slug: "status"
sourcePath: "status.md"
---

## 当前状态：Beta

请查看[更新日志](/changelog/)以了解最新状态。

我们的目标是发布稳定的 v3.0 版本。本路线图概述了正式发布前需要实现的主要功能和改进。请注意，本文档会持续更新，并可能随着优先事项的调整或新见解的出现而更改。

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

本路线图可能根据社区反馈和项目优先事项进行调整。我们会定期更新路线图，以反映进展和方向变化。请将可复现的问题报告为议题；新功能则应通过[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)拉取请求提出。
