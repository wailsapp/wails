---
title: "为什么选择 Wails？"
description: "了解为什么 Wails 是桌面应用的正确选择"
slug: "quick-start/why-wails"
sourcePath: "quick-start/why-wails.md"
---

Wails 将<strong>Go 的高性能与简洁性</strong>和<strong>现代 Web UI 的灵活性</strong>融为一体，让你能使用已经熟悉的工具构建美观的原生桌面应用。

## 用户可感知的性能

**Wails 应用：**

- **约 15MB 的二进制文件**（Electron 为 150MB）
- **约 10MB 的基准内存占用**（Electron 为 100MB 以上）
- **&lt;0.5 秒的启动时间**（Electron 为2-3 秒）
- 使用操作系统提供的 WebView 进行<strong>原生渲染</strong>

用户会感受到你的应用快速、轻量且专业。

## 开发者体验

**一次编写，随处运行：**

- 一套 Go 代码库支持 Windows、macOS 和 Linux
- 可使用任意 Web 框架（React、Vue、Svelte、原生 JS）
- 开发期间支持热重载
- 根据 Go 代码自动生成 TypeScript 绑定

需要维护的代码更少，交付速度更快。

## 可用于生产环境的功能

**所需功能一应俱全：**

- 多个窗口，各自拥有独立的生命周期
- 原生菜单（应用菜单、上下文菜单、系统托盘菜单）
- 具有平台原生 UI 的文件对话框
- 系统集成（通知、剪贴板、键盘快捷键）
- 支持所有平台的代码签名和打包

构建专业应用，而非原型。

## 更快的开发速度

- **一套代码库，三个平台**——一次编写，即可为 Windows、macOS 和 Linux 构建应用
- **运用现有技能**——后端使用 Go，UI 使用 HTML/CSS/JS
- **即时反馈**——开发期间支持热重载，编译时间以秒计
- **小型二进制文件**——15MB 的应用意味着构建更快、下载更快、迭代更快

## 何时选择 Wails

**Wails 非常适合：**

- **业务应用**（CRM、库存管理、仪表盘、管理工具）
- **开发者工具**（数据库客户端、API 测试工具、部署工具）
- **生产力应用**（笔记、任务管理器、时间跟踪工具）
- **创意工具**（图像编辑器、视频处理工具、设计实用工具）
- **内部工具**（公司专用应用、自动化工具）

## 真实案例

@note{type="tip" title="生产环境应用"}
Wails 为数千名用户使用的真实应用提供支持：

- 具有复杂 UI 的<strong>数据库管理工具</strong>
- 处理实时数据的<strong>金融仪表盘</strong>
- 具备原生性能的<strong>视频编辑工具</strong>
- 工程团队使用的<strong>开发实用工具</strong>

[查看应用展示 →](/community/showcase/)

@end

## Wails 的工作原理

Electron 会捆绑完整的浏览器和 Node.js 运行时，而 Wails 采用了截然不同的方法：你的 Go 代码会编译为原生二进制文件，UI 则在操作系统内置的 WebView 中运行。凭借这种架构，Wails 应用拥有体积小、启动快和内存占用低等优势，带来如同原生应用般的体验。

### 架构

Wails 应用由两个无缝通信的主要部分组成：负责业务逻辑和系统操作的 Go 后端，以及用于用户界面的 Web 前端。操作系统提供的 WebView 无需捆绑浏览器即可渲染 UI，而绑定层则为 Go 与 JavaScript 之间的通信提供类型安全保障。

<div style="display: flex; justify-content: center; align-items: center; margin: 2rem 0;">
  <img src="/img/architecture.svg" alt="Wails 架构：Go 后端与你的 Web UI 编译为一个原生二进制文件，通过生成的绑定相连接，并由操作系统的 WebView 渲染" style="max-width: 640px; width: 100%;" />
</div>

这种简单的架构让 JavaScript 代码可以直接调用 Go 函数（通过自动生成的绑定），同时 Go 也可以将事件和数据发回前端。两个层通过高效的内存桥接进行通信，开销低于一毫秒。

**Wails 实现高性能的方式：**

1. **不捆绑运行时**——使用 Go 编译生成的二进制文件
2. **原生 WebView**——使用操作系统提供的渲染引擎
3. **直接的 Go ↔ JS 桥接**——通过内存通信，无网络开销
4. **编译后的二进制文件**——即时启动，无需 JIT 编译

## 后续步骤

现在你已经了解 Wails 提供的功能，接下来开始配置开发环境：

1. **安装 Wails**——只需5分钟即可搭建开发环境 [安装指南 →](/quick-start/installation/)

2. **构建你的第一个应用**——创建一个可运行的应用并了解基础知识 [首个应用教程 →](/quick-start/first-app/)

3. **探索功能**——了解 Wails 能为你的应用做什么 [功能概览 →](/quick-start/next-steps/)

---

<strong>还有疑问？</strong>加入我们的[Discord 社区](https://discord.gg/JDdSxwjhGf)，直接向团队提问。
