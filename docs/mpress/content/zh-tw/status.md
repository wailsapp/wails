---
title: "開發藍圖"
description: "Wails v3 專案狀態、規劃中的功能，以及參與貢獻的方式"
slug: "status"
sourcePath: "status.md"
---

## 目前狀態：Beta 版

請查看[變更記錄](/changelog/)以瞭解最新狀態。

我們的目標是推出穩定的 v3.0 版本。本開發藍圖概述了正式發布前需要實作的主要功能和改進。請注意，這是一份持續更新的文件，可能會隨著優先順序調整或出現新的見解而更新。

## Beta 版相容性承諾

v3 Beta 版契約涵蓋桌面應用程式：

| 平台 | 支援的目標 | 需求與注意事項 |
| --- | --- | --- |
| Windows | amd64 和 arm64 | WebView2 執行階段 |
| macOS | Intel 和 Apple Silicon | 安裝指南中記載的 macOS 和 WebKit 版本 |
| Linux | amd64 和 arm64 | 預設使用 GTK4 + WebKitGTK 6.0；GTK3 + WebKit2GTK 4.1在 v3.0.x 期間仍是`-tags gtk3`舊版選項，並於 v3.1 移除 |

開發所有目標都需要 Go 1.25 或更新版本。Android 和 iOS 支援仍屬實驗性質，不會阻礙桌面版 Beta 的進度。Beta API 以穩定性為目標，但在 v3.0.0 之前，仍可能修正預先發布版本的缺陷以及已明確公告的變更。

## 參與貢獻的方式

- 測試最新的 Beta 版本，並回報可重現的錯誤
- 為文件和範例做出貢獻
- 參與討論，並針對 WEP 草案提供意見回饋
- 針對錯誤修正、文件或已接受的 WEP 提交提取要求

我們歡迎社群做出貢獻。如果您想協助達成這些目標，請加入社群討論。新功能提案應透過 WEP PR 提出，而非建立功能請求議題。

## 意見回饋與更新

本開發藍圖可能會根據社群意見回饋和專案優先順序而變更。我們會定期更新，以反映進度與方向調整。請將可重現的問題回報為議題；新功能則應透過[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) PR 提出。
