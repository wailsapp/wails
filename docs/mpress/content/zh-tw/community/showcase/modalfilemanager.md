---
title: "Modal File Manager"
description: "使用 Wails 建置的桌面應用程式"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager)是一款使用 Web 技術的雙窗格檔案管理器。我的原始設計以 NW.js 為基礎，可在[此處](https://github.com/raguay/ModalFileManager-NWjs)找到。此版本沿用相同的 Svelte 前端程式碼（但自從不再使用 NW.js 後已經過大幅修改），後端則採用[Wails 2](https://wails.io/)實作。採用此實作後，我不再使用命令列的`rm`、`cp`等命令，但系統仍必須安裝 git，才能下載佈景主題和擴充功能。它完全使用 Go 編寫，執行速度比先前版本快得多。

此檔案管理器採用與 Vim 相同的設計原則：由狀態控制鍵盤操作。狀態數量並非固定，而是具有高度可程式化的彈性。因此，可以建立並使用無限多種鍵盤設定。這是它與其他檔案管理器最主要的差異。另有可從 GitHub 下載的佈景主題和擴充功能。
