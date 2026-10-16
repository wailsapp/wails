---
title: "技术概览"
description: "Wails v3 代码库的高层架构与路线图"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## 欢迎阅读 Wails v3 技术文档

本节<strong>并非</strong>介绍社区准则或如何创建拉取请求。 相反，本节将深入讲解<strong>Wails v3 的构建方式</strong>，帮助你快速熟悉代码库，并能放心地着手修改代码。

无论你打算修补运行时、扩展 CLI、制作新模板，还是 仅仅想了解其内部机制，后续页面都会提供你所需的技术背景。

---

## 高层架构

@cards{cols="2"}
◇ Go 后端
每个 Wails 应用的核心都是编译为原生可执行文件的 Go 代码。 它负责应用逻辑、系统集成以及对性能要求较高的操作。

---
▤ Web 前端
UI 使用标准 Web 技术（React、Vue、Svelte、Vanilla 等）编写， 并由轻量级系统 WebView 渲染（Linux/macOS 上使用 WebKit， Windows 上使用 WebView2）。

---
◆ 桥接层
零复制的内存桥接层支持<strong>Go⇄JavaScript</strong>调用，并能自动 转换类型、传播事件和转发错误。

---
▸ CLI 与工具链
`wails3`负责协调项目创建、支持实时重载的开发服务器、资源 打包、交叉编译以及软件包制作（deb、rpm、AppImage、msi、dmg 等）。

@end

---

## 架构概览

**Wails v3——端到端流程**

**[端到端流程图占位符]**

该图展示了<strong>端到端流程</strong>：

1. <strong>CLI</strong>负责驱动代码生成、开发服务器、编译和软件包制作。\
2. <strong>绑定系统</strong>生成粘合代码，使<strong>Web 前端</strong>能够调用<strong>Go 后端</strong>。\
3. 开发期间，<strong>资源服务器</strong>会将请求代理到框架开发服务器；在生产环境中，它则提供嵌入的文件。\
4. 运行时，<strong>桌面运行时</strong>负责管理窗口和操作系统 API，<strong>桥接层</strong>则在 Go 与 JavaScript 之间传递消息。

---

## 本文档涵盖的内容

| 主题 | 重要性 |
| --- | --- |
| **代码库布局** | `/v3`目录的分布以及模块之间的交互方式。 |
| **运行时内部机制** | 窗口管理、系统 API、消息处理器和平台适配层。 |
| **资源与开发服务器** | 开发环境中如何提供 Web 资源，以及生产环境中如何嵌入这些资源。 |
| **构建与软件包制作流水线** | 基于 Taskfile 的工作流、跨平台编译和安装程序生成。 |
| **绑定系统** | 生成类型安全的 Go⇄TS 绑定的静态分析流水线。 |
| **模板系统** | 为`wails3 init -t <framework>`提供支持的生成器架构。 |
| **测试与 CI** | 单元/集成测试框架、GitHub Actions 和竞态检测器使用指南。 |
| **扩展 Wails** | 添加服务、模板或 CLI 子命令。 |

后续各页面将通过具体的代码示例、图表以及相关源文件的引用， 深入讲解这些领域。

---

@note{type="info"}
先决条件：你应当熟悉<strong>Go 1.25+</strong>、TypeScript 基础知识 以及现代前端构建工具。如果你刚开始学习 Go，建议先快速浏览 官方教程。

@end

祝你探索愉快——欢迎深入了解 Wails v3 的内部机制！
