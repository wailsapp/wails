---
title: "Foundation"
description: "可用於正式環境的 Wails v3 桌面應用程式腳手架"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Foundation 版面配置設定頁，可切換九種版面配置預設](/assets/showcase-images/foundation-layouts.webp)

![Foundation 深色主題首頁](/assets/showcase-images/foundation-home-dark.webp)

![Foundation 黑曜主題下的 X-Pro 圖表](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** 是一個可用於正式環境的 **Wails v3 桌面應用程式腳手架**，使用 **Go、React 19、TypeScript 與 Mantine 9** 打造。它為新專案準備好通常需要花上數週才能建置完成的部分——分層的 Go 後端、依平台區分的視窗外觀、資料持久化、主題、國際化與端對端測試——讓你能直接投入業務開發。

## 主要特色

- **各平台皆為原生視窗** — Windows 與 Linux 採用無邊框、由前端自繪的標題列；macOS 採用沉浸式標題列，並保留系統紅綠燈按鈕。
- **九種可切換的版面配置** — 圖示列、展開側邊欄、頂部導覽、融合標題列、右側欄、底部 Dock、可拖曳的懸浮 Dock、浮動卡片與專注模式，皆建構於 CSS Grid 外殼之上；只需一行程式碼即可註冊自訂預設。
- **主題與國際化** — 內建淺色、深色與黑曜預設，以及自訂調色盤編輯器；開箱即支援中文與英文；主視窗、子視窗與系統匣面板的主題與語言即時同步。
- **SQLite 資料持久化** — GORM AutoMigrate，無需遷移檔案；資料庫位置可由使用者選擇，並提供各資料表的空間使用圖表。
- **桌面整合** — 系統匣面板、子視窗、原生檔案與訊息對話方塊，以及透過原生「另存新檔」對話方塊匯出 CSV / XLSX。
- **可靠的工具層** — HTTP 用戶端、綁定 Windows Job Object 的子行程、AES-GCM 加密、輪替日誌與原子寫入。
- **端對端測試** — Playwright 驅動 Wails 伺服器模式的建置（真實的 Go 後端 + 真實的前端），同一套測試也會產生 README 螢幕擷圖。

Wails 讓這一切都封裝在各平台的單一執行檔中：Go 服務透過產生的繫結直接呼叫，前端則遵循嚴格的 MVVM 架構，功能增加後依然易於維護。

支援 Windows、macOS 與 Linux。

[GitHub](https://github.com/FxRayHughes/foundation) |
[發行版本](https://github.com/FxRayHughes/foundation/releases/latest)
