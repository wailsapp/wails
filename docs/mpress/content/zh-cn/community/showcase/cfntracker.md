---
title: "CFN Tracker"
description: "使用 Wails 构建的桌面应用程序"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker)——跟踪任意《Street Fighter 6》或《Street Fighter V》CFN 个人资料中的实时对战。请访问 [网站](https://cfn.williamsjokvist.se/)开始使用。

## 功能

- 实时对战跟踪
- 存储对战日志和统计数据
- 支持通过浏览器源向 OBS 显示实时统计数据
- 同时支持 SF6 和 SFV
- 用户可使用 CSS 创建自己的 OBS 浏览器主题

### 与 Wails 搭配使用的主要技术

- [Task](https://github.com/go-task/task)——封装 Wails CLI，让常用命令更易于使用
- [React](https://github.com/facebook/react)——因其生态系统丰富（radix、framer-motion）而选用
- [Bun](https://github.com/oven-sh/bun)——用于实现快速的依赖项解析和构建
- [Rod](https://github.com/go-rod/rod)——用于身份验证和轮询变更的无头浏览器自动化工具
- [SQLite](https://github.com/mattn/go-sqlite3)——用于存储对战、会话和个人资料
- [服务器发送事件](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)——通过 HTTP 流将跟踪更新发送到 OBS 浏览器源
- [i18next](https://github.com/i18next/)——搭配后端连接器使用，从 Go 层提供本地化对象
- [xstate](https://github.com/statelyai/xstate)——用于身份验证流程和跟踪的状态机
