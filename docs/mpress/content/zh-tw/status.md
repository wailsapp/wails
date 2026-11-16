---
title: "路線圖"
description: "Wails v3 專案狀態、規劃功能及參與貢獻方式"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — 穩定版

請查看[變更記錄](/changelog/)以瞭解最新狀態。

本文件涵蓋 Wails v3.0.0。v3.0.0之後，修正與改進會透過每週版本發布，期間則提供標註日期的每夜預覽版，緊急問題另以緊急修正版處理。請參閱[發布政策](/contributing/release-policy/)。

## 桌面平台相容性

v3.0.0 相容性承諾涵蓋桌面應用程式：

| 平台 | 支援的目標平台 | 需求與備註 |
| --- | --- | --- |
| Windows | amd64 與 arm64 | WebView2 執行階段 |
| macOS | Intel 與 Apple Silicon | 安裝指南中記載的 macOS 與 WebKit 版本 |
| Linux | amd64 與 arm64 | 預設使用 GTK4 + WebKitGTK 6.0；GTK3 + WebKit2GTK 4.1在 v3.0.x 期間仍是`-tags gtk3`舊版選項，並於 v3.1 移除 |

所有目標平台的開發作業都需要 Go 1.26 或更新版本。Android 與 iOS 支援仍處於實驗階段，不屬於穩定版桌面平台相容性承諾的範圍。

## 如何參與貢獻

- 測試 Wails v3.0.0 並回報可重現的錯誤
- 為文件與範例做出貢獻
- 參與討論，並針對 WEP 草案提供意見
- 針對錯誤修正、文件或已接受的 WEP 提交拉取請求

我們歡迎社群貢獻。如果您想協助達成這些目標，請加入社群討論。新功能提案應以 WEP PR 提出，而非建立功能請求議題。

## 意見回饋與更新

本路線圖可能會依社群意見與專案優先順序調整。我們會定期更新，以反映進度與方向變化。請將可重現的問題回報為議題；新功能則應透過[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/v3.0.0/v3/wep/README.md) PR 提案。
