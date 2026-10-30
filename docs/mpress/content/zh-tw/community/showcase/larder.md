---
title: "Larder"
description: "本機優先的 Agent 資源檢查與安全管理工具"
slug: "community/showcase/larder"
sourcePath: "community/showcase/larder.md"
---

![Larder 截圖](/assets/showcase-images/larder-zh.webp)

[Larder](https://github.com/599yongyang/larder) 是一個本機優先的桌面應用程式，面向同時使用多個 AI Agent 的開發者。它會掃描磁碟上真實存在的 Skills、MCP Server、CLI 工具、指令檔案和自訂命令，並解釋每個 Agent 實際能看到什麼。

Larder 可以找出失效連結、重複 Skill 和過期設定。所有寫入都會先產生計畫和快照，修改可以回復。它還支援依分類匯出資料包，以及可選的加密 WebDAV/S3 備份。

使用 Go 1.26、Wails v3 beta.34、React 18、SQLite 和 Tailwind CSS 建立。MIT License。

[在 GitHub 上查看](https://github.com/599yongyang/larder)
