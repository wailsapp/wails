---
title: "Foundation"
description: "基於 Wails v3.1、React 19 與 Mantine 9 的桌面應用程式腳手架"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Foundation 版面配置設定頁，可切換九種版面配置預設](/assets/showcase-images/foundation-layouts.webp)

![Foundation 深色主題首頁](/assets/showcase-images/foundation-home-dark.webp)

![Foundation 黑曜主題下的 X-Pro 圖表](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** 是一個 **Wails v3.1 桌面應用程式腳手架**，使用 **Go、React 19、TypeScript 與 Mantine 9** 打造。它提供分層的 Go 後端、依平台區分的視窗外觀、SQLite 資料持久化、主題、國際化與端對端測試，作為新專案的起點。

## 功能

- **依平台區分的視窗** — Windows 與 Linux 採用無邊框視窗與前端自繪的標題列；macOS 採用沉浸式標題列，並保留系統紅綠燈按鈕。
- **九種版面配置預設** — 圖示列、展開側邊欄、頂部導覽、融合標題列、右側欄、底部 Dock、可拖曳的懸浮 Dock、浮動卡片與專注模式，皆建構於 CSS Grid 外殼之上；可註冊自訂預設。
- **主題與國際化** — 內建淺色、深色與黑曜預設以及自訂調色盤編輯器；介面支援中文與英文；主視窗、子視窗與系統匣面板的主題與語言保持同步。
- **SQLite 資料持久化** — GORM AutoMigrate，無需遷移檔案；資料庫位置可由使用者選擇，並提供各資料表的空間使用圖表。
- **桌面整合** — 系統匣面板、子視窗、原生檔案與訊息對話方塊，以及透過原生「另存新檔」對話方塊匯出 CSV / XLSX。
- **工具層** — HTTP 用戶端、綁定 Windows Job Object 的子行程、AES-GCM 加密、輪替日誌與原子寫入。
- **端對端測試** — Playwright 針對 Wails 伺服器模式的建置執行（真實的 Go 後端與前端）；同一套測試也會產生 README 螢幕擷圖。

前端透過 Wails 產生的繫結呼叫 Go 服務，每個平台皆建置為單一執行檔。前端採用 MVVM 架構，README 中附有將腳手架改名為你自己專案的逐步指南。

可在 Windows、macOS 與 Linux 上執行。

[GitHub](https://github.com/FxRayHughes/foundation) |
[發行版本](https://github.com/FxRayHughes/foundation/releases/latest)
