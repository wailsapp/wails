---
title: "BulletinBoard"
description: "使用 Wails 构建的桌面应用程序"
slug: "community/showcase/bulletinboard"
sourcePath: "community/showcase/bulletinboard.md"
---

![BulletinBoard](/assets/showcase-images/bboard.webp)

[BulletinBoard](https://github.com/raguay/BulletinBoard) 应用程序是一款多用途留言板，可显示静态消息，也可通过对话框从用户处获取脚本所需的信息。它提供 TUI，用于创建之后可用来获取用户信息的新对话框。它的设计目标是在系统中保持运行，需要时显示信息，随后自动隐藏。我设置了一个进程来监视系统中的某个文件，并在文件发生变化时将其内容发送到 BulletinBoard。它非常适合我的工作流。此外，还有一个用于向该程序发送信息的[Alfred 工作流](https://github.com/raguay/MyAlfred/blob/master/Alfred%205/EmailIt.alfredworkflow)。该工作流也可与[EmailIt](https://github.com/raguay/EmailIt)配合使用。
