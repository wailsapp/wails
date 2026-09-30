---
title: "為什麼選擇 Wails？"
description: "瞭解為什麼 Wails 是桌面應用程式的理想選擇"
slug: "quick-start/why-wails"
sourcePath: "quick-start/why-wails.md"
---

Wails 結合了<strong>Go 的效能與簡潔性</strong>以及<strong>現代 Web UI 的彈性</strong>，讓你能使用已熟悉的工具，建置美觀且原生的桌面應用程式。

## 使用者感受得到的效能

**Wails 應用程式：**

- **約 15MB 的二進位檔**（相較之下，Electron 為 150MB）
- **約 10MB 的基準記憶體用量**（相較之下，Electron 為 100MB 以上）
- **&lt;0.5 秒的啟動時間**（相較之下，Electron 為2-3 秒）
- 使用作業系統提供的 WebView 進行<strong>原生算繪</strong>

使用者會感受到你的應用程式快速、輕量且專業。

## 開發者體驗

**一次撰寫，隨處執行：**

- 以單一 Go 程式碼庫支援 Windows、macOS 和 Linux
- 可使用任何 Web 框架（React、Vue、Svelte、原生 JS）
- 開發期間支援熱重新載入
- 從 Go 程式碼自動產生 TypeScript 繫結

減少需要維護的程式碼，更快交付產品。

## 可投入正式環境的功能

**你所需的一切：**

- 多個視窗，各自擁有獨立的生命週期
- 原生選單（應用程式、內容選單、系統匣）
- 採用平台原生 UI 的檔案對話方塊
- 系統整合（通知、剪貼簿、鍵盤快速鍵）
- 支援所有平台的程式碼簽署與封裝

建置專業應用程式，而非原型。

## 加快開發速度

- **單一程式碼庫，支援三個平台**——一次撰寫，即可為 Windows、macOS 和 Linux 建置
- **運用既有技能**——後端使用 Go，UI 使用 HTML/CSS/JS
- **即時回饋**——開發期間支援熱重新載入，編譯時間以秒計
- **小型二進位檔**——15MB 的應用程式意味著建置、下載及反覆開發都更快

## 何時選擇 Wails

**Wails 非常適合：**

- **商務應用程式**（CRM、庫存管理、儀表板、管理工具）
- **開發者工具**（資料庫用戶端、API 測試工具、部署工具）
- **生產力應用程式**（筆記、工作管理、時間追蹤）
- **創作工具**（影像編輯器、影片處理工具、設計公用程式）
- **內部工具**（公司專用應用程式、自動化工具）

## 實際成功案例

@note{type="tip" title="正式環境應用程式"}
Wails 支援數千名使用者實際使用的應用程式：

- 具備複雜 UI 的<strong>資料庫管理工具</strong>
- 處理即時資料的<strong>財務儀表板</strong>
- 具備原生效能的<strong>影片編輯工具</strong>
- 工程團隊使用的<strong>開發公用程式</strong>

[查看案例展示 →](/community/showcase/)

@end

## Wails 的運作方式

Electron 會封裝完整的瀏覽器和 Node.js 執行階段，而 Wails 採用截然不同的方法：你的 Go 程式碼會編譯為原生二進位檔，UI 則在作業系統內建的 WebView 中執行。這種架構可提供小型二進位檔、快速啟動和低記憶體用量，讓 Wails 應用程式擁有原生應用程式般的體驗。

### 架構

Wails 應用程式由兩個能無縫通訊的主要部分組成：負責商務邏輯和系統作業的 Go 後端，以及提供使用者介面的 Web 前端。作業系統提供的 WebView 無須封裝瀏覽器即可算繪 UI，而繫結層則在 Go 與 JavaScript 之間提供型別安全的通訊。

<div style="display: flex; justify-content: center; align-items: center; margin: 2rem 0;">
  <img src="/img/architecture.svg" alt="Wails 架構：Go 後端與您的網頁 UI 編譯為單一原生二進位檔，透過產生的繫結連接，並由作業系統的 WebView 轉譯" style="max-width: 640px; width: 100%;" />
</div>

這種簡潔的架構讓 JavaScript 程式碼可透過自動產生的繫結直接呼叫 Go 函式，Go 也能將事件和資料傳回前端。兩個層級透過高效的記憶體內橋接器通訊，額外負擔不到一毫秒。

**Wails 實現高效能的方式：**

1. **不封裝執行階段**——使用 Go 編譯的二進位檔
2. **原生 WebView**——使用作業系統提供的算繪引擎
3. **直接的 Go ↔ JS 橋接**——在記憶體內通訊，沒有網路額外負擔
4. **編譯後的二進位檔**——立即啟動，無須 JIT 編譯

## 後續步驟

既然你已瞭解 Wails 提供的功能，接下來就開始設定環境：

1. **安裝 Wails**－在 5 分鐘內設定開發環境 [安裝指南 →](/quick-start/installation/)

2. **建置您的第一個應用程式**－建立可運作的應用程式並瞭解基礎知識 [第一個應用程式教學 →](/quick-start/first-app/)

3. **探索功能**－瞭解 Wails 能為您的應用程式提供哪些功能 [功能概覽 →](/quick-start/next-steps/)

---

<strong>還有疑問嗎？</strong>加入我們的 [Discord 社群](https://discord.gg/JDdSxwjhGf)，直接向團隊提問。
