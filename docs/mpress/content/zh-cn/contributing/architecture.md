---
title: "Wails v3 架构"
description: "深入解析 Wails v3 内部每个组成部分的图示与说明"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

Wails v3 是一个由 Go 运行时、 JavaScript 桥接层、任务驱动的工具链和一组模板组成的<strong>全栈桌面应用框架</strong>， 让你能够使用现代 Web 技术交付原生应用。

本页通过四幅图呈现<em>整体全貌</em>：

1. **整体架构**——各个子系统如何连接\
2. **运行时流程**——JS 调用 Go 以及 Go 调用 JS 时会发生什么\
3. **开发与生产环境**——资产服务器的两种模式\
4. **平台实现**——操作系统专用代码位于何处\

---

## 1 · 整体架构

**Wails v3——高层技术栈**

**[高层技术栈图占位符]**

---

## 2 · 运行时调用流程

**运行时——JavaScript ⇄ Go 调用路径**

**[运行时调用流程图占位符]**

要点：

- **不使用 HTTP / IPC**——桥接层使用原生 WebView 的内存通道\
- **方法 ID**——确定性的 FNV 哈希可在 Go 中实现 O(1) 查找\
- **Promise**——错误以拒绝的形式传播，并附带堆栈和错误码

---

## 3 · 开发与生产环境的资产流

**开发 ↔ 生产环境资产服务器**

**[资产流图占位符]**

- 在<strong>开发环境</strong>中，服务器会将未知路径代理到框架的实时重载服务器，并从磁盘提供静态资产。
- 在<strong>生产环境</strong>中，同一套 API 由`go:embed`提供支持，从而生成零依赖的二进制文件。

---

## 4 · 平台专用运行时拆分

**各操作系统的运行时文件**

**[平台拆分图占位符]**

每项功能都遵循以下模式：

1. `pkg/application`中的<strong>通用接口</strong>\
2. `pkg/application/messageprocessor_*.go`中的<strong>消息处理器</strong>入口\
3. `pkg/application/*_{darwin,linux,windows}.go`中的<strong>各操作系统专用实现</strong>（例如`webview_window_darwin.go`、`clipboard_linux.go`、`dialogs_windows.go`、`systemtray_*.go`、`mainthread_*.go`），由构建标签控制。Linux 还在`linux_cgo.go` / `linux_cgo_gtk4.{go,c,h}`中提供了 cgo 桥接层。

`internal/runtime/`仅包含少量`runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go` 构建标签粘合代码，以及`internal/runtime/desktop/`下嵌入的 JS 运行时。

`internal/capabilities/`用于声明各平台的能力集，但 不存在`ErrCapability`哨兵——功能门控通过普通的 构建标签和平台专用的桩返回值实现（例如`nil`或特定于功能的错误）。

---

## 总结

这些图概述了<strong>代码位于何处</strong>、**数据如何流动**，以及 **各层分别承担哪些职责**。 浏览后续详细页面时，请随时参考这些图——它们是你探索 Wails v3 源代码树的地图。
