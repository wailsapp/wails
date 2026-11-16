---
title: "教程"
description: "通过构建应用程序学习 Wails"
slug: "tutorials/overview"
sourcePath: "tutorials/overview.md"
---

如果你想通过实际构建项目来学习 Wails，请选择一篇引导式教程。每篇教程都围绕一个项目展开，最终实现具体成果。

如果这是你第一次访问，请先完成[安装](/quick-start/installation/)和[你的第一个应用](/quick-start/first-app/)。

## 选择教程

完成第一个应用后，请按以下顺序学习桌面应用教程：

| 教程 | 你将练习的内容 |
| --- | --- |
| [二维码服务](/tutorials/01-creating-a-service/) | 通过服务公开 Go 逻辑，并从前端调用它 |
| [TODO 列表](/tutorials/02-todo-vanilla/) | 创建、更新和删除应用程序数据 |
| [笔记应用](/tutorials/03-notes-vanilla/) | 使用原生对话框并将数据保存到文件 |
| [自更新应用](/tutorials/04-self-update-a-wails-app/) | 检查并应用 GitHub Releases 中的更新 |

如需学习可选的移动端教程，请参阅[你的第一个移动应用](/guides/mobile/first-mobile-app/)。移动端支持仍处于实验阶段。

如果你正在开发现有应用并需要完成特定任务，请改用[操作指南](/quick-start/next-steps/#features)。如需了解设计背景，请参阅[核心概念](/concepts/architecture/)。

@note{type="tip" title="刚开始学习 Go？"}
请先完成[Go 语言之旅](https://go.dev/tour/)，再开始学习这些教程。

@end

## 二维码服务

![二维码示例](/assets/qr1.png)

通过构建二维码生成器，学习 Wails 服务的基础知识。本教程将介绍如何把应用程序逻辑组织成可复用服务这一核心概念。

**你将学到：**

- 如何创建和组织 Wails 服务
- 管理外部 Go 依赖项
- 将 Go 方法绑定到前端
- 在 Go 与 JavaScript 之间传递数据
- 组织代码以提高可维护性

<strong>最适合：</strong>希望了解服务架构的 Wails 初学者

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/01-creating-a-service/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### TODO 列表

![TODO 列表应用程序](/assets/todo-app.png)

构建一个界面美观、现代且功能完整的 TODO 列表应用程序。本实践教程使用原生 JavaScript，通过贴近实际的实用应用程序讲解 Wails 的核心模式。

**你将学到：**

- 采用线程安全状态管理的服务式架构
- CRUD 操作（创建、读取、更新、删除）
- Go 与 JavaScript 之间的类型安全绑定
- 在不引入框架复杂性的情况下构建现代 UI
- 正确的错误处理和验证模式

<strong>完成时间：</strong>约20分钟

<strong>最适合：</strong>构建你的第一个完整 Wails 应用程序——非常适合在引入框架复杂性之前理解基础知识

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/02-todo-vanilla/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 笔记

![笔记应用程序](/assets/notes-app.png)

构建一个类似 Apple 备忘录的应用程序，配备原生文件对话框和自动保存功能。本教程演示文件操作、原生对话框和专业 UI 模式等桌面端特有功能。

**你将学到：**

- 原生文件对话框（保存、打开、信息）
- 基于 JSON 的数据持久化
- 防抖自动保存模式
- 专业的双栏桌面布局
- 在 Go 中处理文件系统操作

<strong>完成时间：</strong>约30分钟

<strong>最适合：</strong>学习文件操作和操作系统原生对话框等桌面端特有功能

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/03-notes-vanilla/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>

---

### 自更新 Wails 应用

从全新的`wails3 init`开始，为 Wails 应用添加应用内自更新功能，并完成签名版本验证和辅助模式下的二进制文件替换。使用 GitHub Releases 作为更新源。

**你将学到：**

- `app.Updater`如何集成到 Wails 应用中
- 配置 GitHub Releases 提供程序
- 使用`SHA256SUMS`发布版本，以进行摘要验证
- 添加 Ed25519 签名以防止篡改
- 通过 CSS、自定义 HTML 或自带方案来自定义默认窗口
- 使用`CheckInterval`定期执行后台检查

<strong>完成时间：</strong>约25分钟

<strong>最适合：</strong>发布可更新的桌面应用——涵盖完整的发布流程，而不仅仅是 API

<div style="text-align: center; margin-top: 1rem;">
<a class="mpress-link-button mpress-link-button-secondary" href="/tutorials/04-self-update-a-wails-app/"><span>开始学习</span><span class="mpress-link-button-icon" aria-hidden="true"><svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M17.92 11.62a1 1 0 0 0-.21-.33l-5-5a1 1 0 1 0-1.42 1.42l3.3 3.29H7a1 1 0 0 0 0 2h7.59l-3.3 3.29a1 1 0 1 0 1.42 1.42l5-5a1 1 0 0 0 .21-1.09Z"/></svg></span></a>
</div>
