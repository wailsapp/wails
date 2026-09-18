---
title: "Wails v3 架構"
description: "深入解析 Wails v3 內部每個運作元件的圖表與說明"
slug: "contributing/architecture"
sourcePath: "contributing/architecture.md"
---

Wails v3 是一套<strong>全端桌面應用程式框架</strong>，由 Go 執行階段、 JavaScript 橋接層、工作驅動的工具鏈，以及一組範本所構成，讓您能以現代 Web 技術打造並發布原生應用程式。

本頁透過四張圖呈現<em>整體概觀</em>：

1. **整體架構**－各子系統如何連接\
2. **執行階段流程**－JS 呼叫 Go 以及 Go 呼叫 JS 時會發生什麼事\
3. **開發環境與正式環境**－資產伺服器的兩種模式\
4. **平台實作**－作業系統特定程式碼所在的位置\

---

## 1 · 整體架構

**Wails v3－高階技術堆疊**

**[高階技術堆疊圖預留位置]**

---

## 2 · 執行階段呼叫流程

**執行階段－JavaScript ⇄ Go 呼叫路徑**

**[執行階段呼叫流程圖預留位置]**

重點：

- **不使用 HTTP／IPC**－橋接層使用原生 WebView 的記憶體內通道\
- **方法 ID**－確定性的 FNV 雜湊可在 Go 中進行 O(1) 查詢\
- **Promise**－錯誤會以拒絕狀態傳遞，並附帶堆疊與錯誤碼

---

## 3 · 開發與正式環境的資產流程

**開發 ↔ 正式環境資產伺服器**

**[資產流程圖預留位置]**

- 在<strong>開發環境</strong>中，伺服器會將未知路徑代理至框架的即時重新載入伺服器，並從磁碟提供靜態資產。
- 在<strong>正式環境</strong>中，同一套 API 由`go:embed`提供支援，產生不含任何相依項目的二進位檔。

---

## 4 · 平台特定執行階段的拆分方式

**各作業系統的執行階段檔案**

**[平台拆分圖預留位置]**

每項功能都遵循以下模式：

1. `pkg/application`中的<strong>共用介面</strong>\
2. `pkg/application/messageprocessor_*.go`中的<strong>訊息處理器</strong>進入點\
3. `pkg/application/*_{darwin,linux,windows}.go`中的<strong>各作業系統專屬實作</strong>（例如`webview_window_darwin.go`、`clipboard_linux.go`、`dialogs_windows.go`、`systemtray_*.go`、`mainthread_*.go`），並由建置標籤保護。Linux 另在`linux_cgo.go`／`linux_cgo_gtk4.{go,c,h}`中提供 cgo 橋接層。

`internal/runtime/`僅包含少量的`runtime{_darwin,_linux,_windows,_android,_dev,_prod}.go` 建置標籤黏合程式碼，以及`internal/runtime/desktop/`下內嵌的 JS 執行階段。

`internal/capabilities/`用於宣告各平台的功能集，但 並不存在`ErrCapability`哨兵值－功能閘控是透過一般的 建置標籤和平台特定的虛設回傳值完成（例如`nil`或功能特定錯誤）。

---

## 摘要

這些圖表概述了<strong>程式碼所在的位置</strong>、**資料如何流動**，以及 **各層分別負責哪些職責**。 瀏覽後續的詳細頁面時，請將這些圖表放在手邊－它們就是探索 Wails v3 原始碼樹的地圖。
