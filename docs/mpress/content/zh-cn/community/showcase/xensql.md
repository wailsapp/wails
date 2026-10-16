---
title: "XenSQL"
description: "本地优先的 SQL 工作台，支持多种数据库并提供高级查询工具"
slug: "community/showcase/xensql"
sourcePath: "community/showcase/xensql.md"
---

![支持行内编辑的 XenSQL 表编辑器](/assets/showcase-images/xensql-1.png) ![提供表和列建议的 XenSQL 查询编辑器](/assets/showcase-images/xensql-2.png) ![XenSQL 多事务查询](/assets/showcase-images/xensql-3.png)

**[XenSQL](https://github.com/Bare7a/XenSQL)** 是一款使用<strong>Go、Wails 和 React</strong>构建的<strong>快速、本地优先的 SQL 桌面工作台</strong>。它将 PostgreSQL、MySQL/MariaDB 和 SQLite 集成到一个简洁、具有原生体验的界面中——无需云服务，不收集遥测数据，也无需账户。

## 主要亮点

- **强大的 SQL 编辑器**——基于 Monaco，提供智能的架构感知自动补全、多语句执行、流式结果，以及每条语句独立的结果标签页
- **高级数据查看器**——交互式 JSON 检查器、可识别语法的单元格编辑器（JSON、XML、HTML、文本）、行内编辑和完整记录检查
- **流畅的数据编辑**——浏览表、在行内暂存更改、执行批量操作，以及安全地执行`INSERT`/`UPDATE`/`DELETE`操作，并支持`RETURNING`
- **效率功能**——架构浏览器、已保存的查询、查询历史记录、快速搜索（`Ctrl+P`），以及键盘优先的工作流
- **导出选项**——CSV、JSON、Markdown、SQL INSERT 语句

完全离线且便于携带。所有内容都存储在一个随应用一起移动的本地`XenSQL-data/`文件夹中。

**支持的数据库**：PostgreSQL、MySQL、MariaDB 和 SQLite（支持只读模式和安全传输选项）。

专为追求速度、清晰度和掌控力，又不想承受传统 SQL 工具臃肿负担的开发者而设计。
