---
title: "專案狀態"
description: "Wails v3 Beta 相容性、安全性支援與升級指南"
slug: "status"
sourcePath: "status.md"
---

## 目前狀態：Beta 版

請查看[變更記錄](/changelog/)以瞭解最新狀態。

我們的目標是推出穩定的 v3.0 版本。Wails v2 仍是目前的穩定版本，並持續獲得修正。部署前，請使用您的應用程式測試 Beta 版本。

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

請將可重現的問題回報為議題；新功能則應透過[WEP（Wails Enhancement Proposal）](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) PR 提出。

## 使用 Beta 版本

請固定 CLI、Go 模組和前端執行階段的確切版本，而非追蹤 `latest`。現有 Alpha 專案請遵循 [Alpha 升級至 Beta 指南](/migration/alpha-to-beta/)。

[安全性政策](https://github.com/wailsapp/wails/blob/master/SECURITY.md)將 v3 Beta 版本列為支援對象，Alpha 版本則不受支援。請透過[私密漏洞回報](https://github.com/wailsapp/wails/security/advisories/new)回報漏洞，不要建立公開議題。

## 追蹤中的工作

- [標有 v3 的未解決錯誤](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [標有 P0 或 P1 的未解決 v3 議題](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [發布里程碑](https://github.com/wailsapp/wails/milestones)

這些即時查詢依賴議題標籤，並非完整清單，也不承諾發布日期或範圍。請閱讀議題，評估對您專案的影響。
