---
title: "路線圖"
description: "Wails v3 專案狀態、規劃功能及參與貢獻方式"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — 穩定版

請查看[變更日誌](/changelog/)，瞭解最新狀態。

本文件涵蓋 Wails v3.0.0。v3.0.0之後，修正與改進會透過每週版本發布；兩次發布之間會提供標註日期的每夜預覽版，緊急問題則會發布 Hotfix。請參閱[發布政策](/contributing/release-policy/)。

## 桌面平台相容性

v3.0.0相容性契約涵蓋桌面應用程式：

| 平台 | 支援的目標 | 需求與注意事項 |
| --- | --- | --- |
| Windows | amd64 和 arm64 | WebView2 執行階段 |
| macOS | Intel 和 Apple Silicon | 安裝指南中所列的 macOS 和 WebKit 版本 |
| Linux | amd64 和 arm64 | 預設使用 GTK4 + WebKitGTK 6.0；GTK3 + WebKit2GTK 4.1在v3.0.x 期間仍是`-tags gtk3`舊版選項，並於v3.1移除 |

所有目標在開發時都需要 Go 1.26或更新版本。Android 和 iOS 支援仍處於實驗階段，不在穩定版桌面平台相容性契約的範圍內。

## 如何參與貢獻

- 測試 Wails v3.0.0並回報可重現的錯誤
- 為文件和範例做出貢獻
- 參與討論並針對 WEP 草案提供意見
- 針對錯誤修正、文件或已接受的 WEP 提交提取要求

我們歡迎社群貢獻。如果您想協助達成這些目標，請加入社群討論。新功能提案應透過 WEP PR 提出，而非建立功能要求議題。

## 意見回饋與更新

本路線圖可能會因社群意見和專案優先順序而變更。我們會定期更新，以反映進度及方向調整。請將可重現的問題回報為議題；若要提議新功能，請提交[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md) PR。
