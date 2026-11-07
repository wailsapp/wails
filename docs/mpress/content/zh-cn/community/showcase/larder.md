---
title: "Larder"
description: "本地优先的 Agent 资源检查与安全管理工具"
slug: "community/showcase/larder"
sourcePath: "community/showcase/larder.md"
---

![Larder 截图](/assets/showcase-images/larder-zh.webp)

[Larder](https://github.com/599yongyang/larder) 是一个本地优先的桌面应用，面向同时使用多个 AI Agent 的开发者。它会扫描磁盘上真实存在的 Skills、MCP Server、CLI 工具、指令文件和自定义命令，并解释每个 Agent 实际能看到什么。

Larder 可以找出失效链接、重复 Skill 和过期配置。所有写入都会先生成计划和快照，修改可以回滚。它还支持按分类导出数据包，以及可选的加密 WebDAV/S3 备份。

使用 Go 1.26、Wails v3 beta.34、React 18、SQLite 和 Tailwind CSS 构建。MIT License。

[在 GitHub 上查看](https://github.com/599yongyang/larder)
