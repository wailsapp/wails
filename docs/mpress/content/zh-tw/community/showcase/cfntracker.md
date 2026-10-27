---
title: "CFN Tracker"
description: "使用 Wails 建置的桌面應用程式"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker) — 追蹤任何 Street Fighter 6 或 V CFN 個人檔案的即時對戰。前往 [網站](https://cfn.williamsjokvist.se/)即可開始使用。

## 功能

- 即時對戰追蹤
- 儲存對戰記錄與統計資料
- 支援透過瀏覽器來源向 OBS 顯示即時統計資料
- 同時支援 SF6 與 SFV
- 使用者可使用 CSS 建立自己的 OBS 瀏覽器主題

### 與 Wails 搭配使用的主要技術

- [Task](https://github.com/go-task/task) — 封裝 Wails CLI，讓常用命令更容易使用
- [React](https://github.com/facebook/react) — 因其豐富的生態系統（radix、framer-motion）而採用
- [Bun](https://github.com/oven-sh/bun) — 用於快速解析相依套件及縮短建置時間
- [Rod](https://github.com/go-rod/rod) — 用於驗證及輪詢變更的無頭瀏覽器自動化工具
- [SQLite](https://github.com/mattn/go-sqlite3) — 用於儲存對戰、工作階段及個人檔案
- [伺服器傳送事件](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) — 用於將追蹤更新傳送至 OBS 瀏覽器來源的 HTTP 串流
- [i18next](https://github.com/i18next/) — 搭配後端連接器，從 Go 層提供本地化物件
- [xstate](https://github.com/statelyai/xstate) — 用於驗證流程與追蹤的狀態機
