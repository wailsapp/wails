---
title: "Modal File Manager"
description: "使用 Wails 构建的桌面应用程序"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager)是一款使用 Web 技术的双窗格 文件管理器。我最初的设计基于 NW.js，可以在[此处](https://github.com/raguay/ModalFileManager-NWjs)找到。此 版本使用相同的基于 Svelte 的前端代码（但自从弃用 NW.js 后，代码已经过大幅修改）， 后端则采用[Wails 2](https://wails.io/)实现。使用此实现后，我不再 使用命令行`rm`、`cp`等命令，但系统中必须安装 git， 才能下载主题和扩展。它完全使用 Go 编写，运行速度比以前的版本快得多。

此文件管理器的设计遵循与 Vim 相同的原则：由状态控制的键盘操作。状态 数量并不固定，而且具有很强的可编程性。因此，可以创建和使用无限多种键盘 配置。这是它与其他文件管理器的主要区别。可从 GitHub 下载主题和扩展。
