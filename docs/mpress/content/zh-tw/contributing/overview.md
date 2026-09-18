---
title: "技術概覽"
description: "Wails v3 程式碼庫的高階架構與導覽路線"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## 歡迎閱讀 Wails v3 技術文件

本節<strong>不是</strong>在介紹社群準則或如何建立提取請求。 相反地，本節將深入說明<strong>Wails v3 的建構方式</strong>，協助你快速熟悉 程式碼庫，並有信心地開始修改程式碼。

無論你打算修補執行階段、擴充 CLI、製作新範本，還是 只想瞭解其內部機制，後續頁面都會提供你所需的技術 背景。

---

## 高階架構

@cards{cols="2"}
◇ Go 後端
每個 Wails 應用程式的核心，都是編譯為原生可執行檔的 Go 程式碼。 它負責應用程式邏輯、系統整合，以及對效能要求嚴苛的 作業。

---
▤ Web 前端
UI 使用標準 Web 技術（React、Vue、Svelte、Vanilla 等）編寫， 並由輕量的系統 WebView 呈現（Linux/macOS 使用 WebKit，Windows 使用 WebView2）。

---
◆ 橋接層
零複製的記憶體內橋接層可讓<strong>Go⇄JavaScript</strong>互相呼叫，並自動進行 型別轉換、事件傳播及錯誤轉送。

---
▸ CLI 與工具
`wails3`負責協調專案建立、即時重新載入開發伺服器、資源 打包、交叉編譯及封裝（deb、rpm、AppImage、msi、dmg 等）。

@end

---

## 架構概覽

**Wails v3 端對端流程**

**[端對端流程圖預留位置]**

此圖顯示<strong>端對端流程</strong>：

1. <strong>CLI</strong>負責驅動程式碼產生、開發伺服器、編譯及封裝。\
2. <strong>繫結系統</strong>會產生黏合程式碼，讓<strong>Web 前端</strong>能夠呼叫<strong>Go 後端</strong>。\
3. 開發期間，<strong>資源伺服器</strong>會將請求代理至框架的開發伺服器；正式環境中則提供內嵌檔案。\
4. 執行期間，<strong>桌面執行階段</strong>負責管理視窗與作業系統 API，而<strong>橋接層</strong>則在 Go 與 JavaScript 之間傳遞訊息。

---

## 本文件涵蓋的內容

| 主題 | 重要性 |
| --- | --- |
| **程式碼庫配置** | `/v3`目錄的配置圖，以及模組之間的互動方式。 |
| **執行階段內部機制** | 視窗管理、系統 API、訊息處理器及平台轉接層。 |
| **資源與開發伺服器** | Web 資源在開發環境中的提供方式，以及在正式環境中的內嵌方式。 |
| **建構與封裝管線** | 以 Taskfile 為基礎的工作流程、跨平台編譯及安裝程式產生。 |
| **繫結系統** | 產生型別安全 Go⇄TS 繫結的靜態分析管線。 |
| **範本系統** | 驅動`wails3 init -t <framework>`的產生器架構。 |
| **測試與 CI** | 單元／整合測試框架、GitHub Actions 及競爭偵測器指引。 |
| **擴充 Wails** | 新增服務、範本或 CLI 子命令。 |

後續各頁會透過具體的程式碼範例、圖表，以及相關原始碼檔案的 參照，深入探討這些領域。

---

@note{type="info"}
先備知識：你應熟悉<strong>Go 1.25+</strong>、基本 TypeScript， 以及現代前端建構工具。如果你剛開始接觸 Go，建議先快速瀏覽 官方導覽。

@end

祝你探索愉快，也歡迎深入瞭解 Wails v3 的內部機制！
