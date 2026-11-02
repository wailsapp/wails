---
title: "Zentro"
description: "面向 PostgreSQL、MySQL、SQL Server 和 SQLite 的以项目为中心的桌面 SQL IDE"
slug: "community/showcase/zentro"
sourcePath: "community/showcase/zentro.md"
---

![Zentro SQL 编辑器，带有基于架构的自动补全、架构浏览器和记录面板](/assets/showcase-images/zentro.png)

**[Zentro](https://zentro-ide.com)** 是一款面向 PostgreSQL、MySQL、SQL Server 和 SQLite 的 **以项目为中心的桌面 SQL IDE** ，使用 **Go、Wails v2 和 React**构建。每个项目都清楚自己指向哪个环境，在执行任何不可撤销的操作之前，Zentro 都会先询问，因此原本用于预发布环境的查询不会落到生产环境。

## 主要亮点

- **环境感知** - 项目自带环境（开发、预发布、生产），并有醒目的生产环境标识、语句风险分级，以及破坏性写入前的受影响行预览
- **SQL 编辑器** - 基于架构的自动补全、行内诊断和快速修复
- **架构浏览器和 ERD** - 浏览表、列、键和索引，并将关系可视化
- **基于 Git 的跟踪** - 已保存的脚本和架构快照像代码一样进行版本管理，带有历史记录和差异对比
- **SSH 隧道和内置终端** - 通过隧道访问数据库，并在查询旁保留本地或 SSH 终端
- **使用自己密钥的 AI 助手** - 连接到你配置的模型端点，绝不经过 Zentro 运营的服务器，并提供面向外部智能体的 MCP 服务器
- **键盘优先** - 在浏览器和结果网格中按焦点导航，提供命令面板和组合键序列
- **默认私密** - 数据库凭据保存在操作系统的钥匙串中

Zentro 为 Windows、macOS 和 Linux 的每个平台提供单个原生二进制文件。核心在 GPL-3.0 下开源，可选的付费 Pro 版增加生产环境防护等功能，并提供 30 天免费试用；Community Edition 的源代码是公开的。

[网站](https://zentro-ide.com) |
[GitHub](https://github.com/ZentroIDE/Zentro) |
[下载](https://zentro-ide.com/download)
