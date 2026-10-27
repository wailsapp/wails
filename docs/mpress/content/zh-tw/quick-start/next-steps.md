---
title: "后续步骤"
description: "构建第一个应用后该做什么"
slug: "quick-start/next-steps"
sourcePath: "quick-start/next-steps.md"
---

你已经构建了第一个 Wails 应用。请按照下面的学习路径继续，或者，如果你已经知道自己想构建什么，可以直接跳到相应指南。

## 深入学习

每次通过一个完整的应用继续学习：

1. <strong>[创建二维码服务](/tutorials/01-creating-a-service/)。</strong>练习如何通过服务向前端公开 Go 逻辑。
2. <strong>[构建待办事项列表](/tutorials/02-todo-vanilla/)。</strong>通过服务调用添加、修改和删除应用数据。
3. <strong>[构建笔记应用](/tutorials/03-notes-vanilla/)。</strong>使用文件、原生对话框和已保存的数据。

准备好探索发布流程后，请按照[应用自动更新教程](/tutorials/04-self-update-a-wails-app/)操作。请参阅[所有教程](/tutorials/overview/)，了解各课的学习成果。

## 功能

<strong>已经在构建自己的应用了？</strong>请针对当前任务使用相应的操作指南：

| 我想要…… | 从这里开始 |
| --- | --- |
| 从前端调用 Go | [从 JavaScript 调用 Go](/features/bindings/methods/) |
| 打开或管理窗口 | [窗口基础知识](/features/windows/basics/) |
| 添加原生菜单或系统托盘图标 | [应用菜单](/features/menus/application/)或[系统托盘菜单](/features/menus/systray/) |
| 让用户打开或保存文件 | [文件对话框](/features/dialogs/file/) |
| 在组件之间发送事件 | [事件系统](/features/events/system/) |
| 复制文本或接收拖放的文件 | [剪贴板](/features/clipboard/basics/)或[文件拖放](/features/drag-and-drop/files/) |
| 添加键盘快捷键 | [键盘快捷键](/features/keyboard/shortcuts/) |
| 构建和打包应用 | 先参阅[构建应用](/guides/build/building/)，再为[Windows](/guides/build/windows/)、[macOS](/guides/build/macos/)或[Linux](/guides/build/linux/)打包 |

这些页面位于侧边栏的**操作指南**下。你可以直接打开相关主题，无须按顺序阅读整个章节。

## 了解 Wails 的工作原理

如果想了解代码背后的设计，请参阅**核心概念**：

- [Wails 的工作原理](/concepts/architecture/)介绍 Go 后端、前端和原生 WebView 之间的连接。
- [Go 与前端的桥接](/concepts/bridge/)说明调用如何在 Go 与 JavaScript 之间传递。
- [应用生命周期](/concepts/lifecycle/)说明启动和关闭过程。
- [构建系统](/concepts/build-system/)说明项目任务和前端工具的作用。

## 查阅 API

当你需要查找类型、方法、选项或命令，而不是学习教程时，请使用[API 参考](/reference/overview/)。对于命令行操作，请先查看[CLI 命令概览](/reference/cli/)；如需了解各个命令的详细信息，请参阅[CLI 标志和示例](/guides/cli/)。

## 获取帮助

查看[常见问题](/faq/)、浏览[可运行的示例](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples)，或在[Discord](https://discord.gg/JDdSxwjhGf)中提问。平台特定问题归类在侧边栏的**故障排除**下。
