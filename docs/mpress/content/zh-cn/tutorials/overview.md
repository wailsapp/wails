---
title: "教程"
description: "通过构建应用程序学习 Wails"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

通过分步构建完整的应用程序来学习 Wails 概念。每篇教程都包含可运行的代码、说明和实用模式。

@note{type="tip" title="刚接触 Go？"}
开始学习教程前，请先完成[Go 之旅](https://go.dev/tour/)。

@end

## 二维码服务

![二维码示例](/assets/qr1.png)

通过构建二维码生成器，学习 Wails 服务的基础知识。本教程将介绍如何把应用程序逻辑组织成可复用服务这一核心概念。

**你将学到：**

- 如何创建并组织 Wails 服务
- 管理外部 Go 依赖项
- 将 Go 方法绑定到前端
- 在 Go 与 JavaScript 之间传递数据
- 以便于维护的方式组织代码

<strong>最适合：</strong>希望了解服务架构的 Wails 初学者

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 待办事项列表

![待办事项列表应用程序](/assets/todo-app.png)

构建一个界面美观、现代的完整待办事项列表应用程序。这篇实践教程使用原生 JavaScript，通过一个贴近实际的实用应用程序讲解 Wails 核心模式。

**你将学到：**

- 采用线程安全状态管理的服务式架构
- CRUD 操作（创建、读取、更新、删除）
- Go 与 JavaScript 之间的类型安全绑定
- 无需引入框架复杂性即可构建现代 UI
- 正确的错误处理与验证模式

<strong>完成时间：</strong>约 20 分钟

<strong>最适合：</strong>你的第一个完整 Wails 应用程序——非常适合在引入框架复杂性之前理解基础知识

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 笔记

![笔记应用程序](/assets/notes-app.png)

构建一款采用 Apple 备忘录风格、支持原生文件对话框和自动保存功能的应用程序。本教程演示文件操作、原生对话框和专业 UI 模式等桌面端特有功能。

**你将学到：**

- 原生文件对话框（保存、打开、信息）
- 基于 JSON 的数据持久化
- 采用防抖的自动保存模式
- 专业的双栏桌面布局
- 在 Go 中进行文件系统操作

<strong>完成时间：</strong>约 30 分钟

<strong>最适合：</strong>学习文件操作和原生操作系统对话框等桌面端特有功能

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 可自行更新的 Wails 应用程序

从全新的`wails3 init`开始，为 Wails 应用程序添加应用内自行更新功能，直至实现签名发布版本验证和辅助模式二进制文件替换。使用 GitHub Releases 作为更新源。

**你将学到：**

- `app.Updater`如何集成到 Wails 应用程序中
- 配置 GitHub Releases 提供程序
- 使用`SHA256SUMS`发布版本，以便进行摘要验证
- 添加 Ed25519 签名以防篡改
- 通过 CSS、自定义 HTML 或 BYO 自定义默认窗口
- 使用`CheckInterval`定期执行后台检查

<strong>完成时间：</strong>约 25 分钟

<strong>最适合：</strong>交付可更新的桌面应用程序——涵盖完整的发布流水线，而不仅仅是 API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
