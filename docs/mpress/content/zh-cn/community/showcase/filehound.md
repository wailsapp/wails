---
title: "FileHound 导出工具"
description: "使用 Wails 构建的桌面应用程序"
slug: "community/showcase/filehound"
sourcePath: "community/showcase/filehound.md"
---

![FileHound 导出工具](/assets/showcase-images/filehound.webp)

[FileHound 导出工具](https://www.filehound.co.uk/) FileHound 是一个云端文档管理平台，提供安全的文件保留、业务流程自动化和 SmartCapture 功能。

FileHound 导出工具支持 FileHound 管理员运行安全的文档和数据提取任务，用于其他备份和恢复用途。此应用程序会根据您选择的筛选条件，下载 FileHound 中保存的所有文档和/或元数据。元数据将同时以 JSON 和 XML 格式导出。

后端使用以下技术构建：

- Go 1.15
- Wails 1.11.0
- go-sqlite3 1.14.6
- go-linq 3.2

前端使用以下技术：

- Vue 2.6.11
- Vuex 3.4.0
- TypeScript
- Tailwind 1.9.6
