---
title: "BMEditor"
description: "功能齐全的代码编辑器，打开即用，需要更多时还能扩展"
slug: "community/showcase/bmeditor"
sourcePath: "community/showcase/bmeditor.md"
---

![打开了 Go 项目的 BMEditor](/assets/showcase-images/bmeditor-editor.webp)

![调试 Go 程序，每个变量的值显示在代码旁边](/assets/showcase-images/bmeditor-debug.webp)

![Docker 面板，容器按 compose 项目分组](/assets/showcase-images/bmeditor-docker.webp)

![扩展窗口，BMEditor 自家的扩展随时可以安装](/assets/showcase-images/bmeditor-extensions.webp)

![侧边栏中的 Claude Code，它的修改正在编辑器中等待审阅](/assets/showcase-images/bmeditor-claude.webp)

BMEditor（Bare Minimum Editor (batteries included)）是一款面向 Windows 和 Linux 的桌面代码编辑器，打开就能用。调试器、终端、Docker、项目中的 Git 改动，以及 Markdown、HTML 和 OpenAPI 的实时预览都内置在编辑器里，不用先去找插件或做配置。Go 开箱即用，编辑器本身也可扩展，其他语言和工具都可以通过扩展加入。

在同一个窗口里，你可以打开 WSL 中的文件夹、通过 SSH 连接的服务器，或 Docker 容器内部，对方那边不需要安装任何东西。如果你用 Claude Code 或 GitHub Copilot CLI，它会在侧边栏中运行，它改过的文件会以审阅的形式回到你面前，每处改动都可以接受或拒绝。布局和快捷键沿用 VS Code，第一分钟就会觉得顺手。

它使用 Wails 3、Go 和 Vue 构建，可免费使用。

[官网](https://cn8001.github.io/BMEditor-Releases/) | [GitHub 发布版](https://github.com/Cn8001/BMEditor-Releases/releases)
