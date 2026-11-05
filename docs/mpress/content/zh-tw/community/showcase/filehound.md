---
title: "FileHound 匯出工具"
description: "使用 Wails 建置的桌面應用程式"
slug: "community/showcase/filehound"
sourcePath: "community/showcase/filehound.md"
---

![FileHound 匯出工具](/assets/showcase-images/filehound.webp)

[FileHound 匯出工具](https://www.filehound.co.uk/) FileHound 是一個雲端 文件管理平台，專為安全的檔案保留、業務流程 自動化及 SmartCapture 功能而打造。

FileHound 匯出工具可讓 FileHound 管理員執行 安全的文件與資料擷取工作，以供替代備份與復原 之用。此應用程式會根據您選擇的篩選條件，下載儲存於 FileHound 中的所有文件及／或中繼資料。中繼資料將同時匯出為 JSON 與 XML 格式。

後端使用以下技術建置：

- Go 1.15
- Wails 1.11.0
- go-sqlite3 1.14.6
- go-linq 3.2

前端使用以下技術：

- Vue 2.6.11
- Vuex 3.4.0
- TypeScript
- Tailwind 1.9.6
