---
title: "資產伺服器"
description: "Wails v3 如何在開發與正式環境中提供及嵌入 Web 資產"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## 概觀

每個 Wails 應用程式都以<strong>單一原生執行檔</strong>的形式發佈，其中結合了：

1. 您的<em>Go</em>後端
2. 一個<em>Web</em>前端（HTML + JS + CSS）

<strong>資產伺服器</strong>是實現此功能的黏合層。它有<strong>兩種運作模式</strong>，在編譯時透過 Go 建置標籤選取：

| 模式 | 標籤 | 用途 |
| --- | --- | --- |
| **開發** | `//go:build !production` | 透過熱重新載入快速反覆開發 |
| **正式環境** | `//go:build production` | 零相依性的嵌入式資產 |

實作位於`v3/internal/assetserver/`，並清楚拆分為不同檔案：

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## 開發模式

### 生命週期

1. `wails3 dev`啟動後，會執行`build/Taskfile.yml`中定義的工作（通常是`npm run dev`），以<strong>產生您的前端開發伺服器程序</strong>（Vite、SvelteKit、React-SWC……）。
2. CLI 將`WAILS_VITE_PORT`設為 Wails 開發連接埠，並將`FRONTEND_DEVSERVER_URL`設為指向框架執行中開發伺服器的<strong>完整</strong> URL（`http://host:port` / `https://host:port`）。請參閱`internal/commands/dev.go`。
3. 開發資產伺服器（透過`build_dev.go`中的`//go:build !production`編譯加入）會經由`GetDevServerURL()`讀取`FRONTEND_DEVSERVER_URL`，並將非執行階段流量反向代理至該伺服器。
4. 靜態檔案（`/assets/logo.svg`）可透過`asset_fileserver.go`**直接從磁碟提供**（速度較快），而任何未知內容都會<strong>代理</strong>至框架開發伺服器，讓您獲得<em>即時</em>的熱模組替換。

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### 功能

- **即時重新載入**——Vite、SvelteKit 等會透過 WebSocket 注入 HMR；開發資產伺服器會透明地代理它。
- **原始碼對應支援**——由於資產並未封裝，瀏覽器開發人員工具可將錯誤對應回原始碼。
- **無須重新編譯 Go**——只有前端會重新建置；在您變更`.go`檔案前，Go 程式碼會持續執行。

### 切換框架

開發代理<strong>不受框架限制</strong>。Wails CLI 啟動您的開發工作時，會發佈兩個環境變數：

| 環境變數 | 來源 | 含義 |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go`（`wailsVitePort`常數） | 預設開發連接埠（除非傳入`--port`，否則為9245）——您的 Vite 設定應採用此值 |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | Wails 將代理至的完整 URL；可在 Go 中透過`assetserver.GetDevServerURL()`（`build_dev.go`）讀取 |

v3 原始碼樹中沒有`VITE_PORT`、`FRONTEND_DEV_PORT`或`WAILSDEV_VERBOSE`環境變數。

新增範本 → 定義其開發工作 → 資產伺服器即可直接運作。

---

## 正式環境模式

執行`wails3 build`時，建置流程會：

1. 執行前端<strong>正式環境建置</strong>（`npm run build`），產生`frontend/dist/**`。
2. 透過您應用程式本身套件中的`go:embed`，將該目錄<strong>嵌入</strong>應用程式（通常是位於`main.go`旁的`//go:embed all:frontend/dist`）。
3. 使用`-tags production`編譯 Go 二進位檔（Taskfile 包裝函式會透過`EXTRA_TAGS`轉送此值）。

`internal/assetserver/build_production.go`是用來切換至正式環境程式碼路徑的建置標籤虛設常式。`internal/assetserver/bundled_assetserver.go`是<strong>手寫</strong>的；它會包裝位於`bundledassets/`中的執行階段 JS，並非產生的檔案。

### 要求處理

實際的處理常式是`internal/assetserver/assetserver.go` / `asset_fileserver.go`。概念上會：

1. 嘗試從要求的路徑取得嵌入式靜態資產。
2. 若未找到，則回退至`index.html`以支援 SPA 路由。
3. 若副檔名未知，則探測內容類型（`content_type_sniffer.go`）。
4. 設定合理的快取標頭。

- **MIME 偵測**——對於沒有副檔名的檔案，會從開頭約512個位元組探測其內容類型（`content_type_sniffer.go`），並將結果快取於`mimecache.go` / `ringqueue.go`。
- **安全性標頭**——禁止`file://`導覽，並設定`nosniff`。

由於所有內容均已嵌入，因此發佈的二進位檔<strong>沒有外部相依性</strong>（即使在 Windows 上亦然）。

---

## 銜接開發模式 ↔ 正式環境模式

從`pkg/application`的角度來看，兩種模式都公開<strong>相同的公用介面</strong>：一個具有`Handler http.Handler`的`AssetOptions`結構，以及`internal/assetserver/`內的中介軟體與生命週期接線。開發／正式環境模式完全透過 Go 建置標籤切換，因此兩種模式下的應用程式碼完全相同。

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## 前端框架的整合方式

### 範本

每個隨附的範本（React、Vue、Svelte、Solid、Vanilla……）都包含：

- `build/Taskfile.yml`
- `frontend/vite.config.ts`（或等效配置）

其 Vite 或等效設定會讀取`WAILS_VITE_PORT`，並將開發伺服器繫結至該連接埠。隨後，CLI 會發布實際的`FRONTEND_DEVSERVER_URL`，供應用程式內的代理使用。

前端框架與 Go 保持完全解耦：

- 建置時無需匯入任何 Wails JS SDK——`/wails/runtime.js`會在執行階段由資源伺服器提供。
- 任何具有 HTTP 開發伺服器的框架都可以接入。

---

## 擴充／自訂

需要自訂標頭、身分驗證或 gzip 嗎？

1. 定義一個`middleware.Middleware`（`func(http.Handler) http.Handler`的別名，宣告於`internal/assetserver/middleware.go`中）。
2. 透過`internal/assetserver/options.go`公開的設定，將其連接至你的`application.AssetOptions`。
3. 開發與正式環境中的行為完全相同——沒有依模式區分的中介軟體清單。

---

## 主要原始碼檔案

| 檔案 | 用途 |
| --- | --- |
| `build_dev.go` / `build_production.go` | 選擇開發或正式模式的建置標籤包裝層 |
| `assetserver.go` / `asset_fileserver.go` | 核心 HTTP 處理常式 |
| `assetserver_dev.go` | 指向`FRONTEND_DEVSERVER_URL`的反向代理 |
| `bundled_assetserver.go` | 以手寫方式實作、封裝`bundledassets/`的包裝層 |
| `options.go` | 面向`application.AssetOptions`的配置 |
| `mimecache.go` / `ringqueue.go` | MIME 快取與小型 LRU |

---

## 常見陷阱與除錯

- **正式環境出現白畫面**——通常是 SPA 路由問題：請確保開發伺服器會針對未知路徑提供`index.html`，且能進入內嵌正式環境處理常式的後備處理流程。
- **開發環境中出現404**——你的 Vite 設定未繫結至`WAILS_VITE_PORT`，或 CLI 無法連接開發伺服器，因此無法填入`FRONTEND_DEVSERVER_URL`。
- **大型資源**——內嵌會使二進位檔案膨脹。請從獨立來源提供大型媒體，或透過自訂`http.Handler`串流傳輸。

---

現在你已了解 Wails <strong>資源伺服器</strong>如何在<strong>開發</strong>與<strong>正式</strong>環境中，將網頁程式碼提供給原生視窗。掌握這一層後，你便能有把握地為載入問題除錯、加入中介軟體，甚至改用完全不同的前端工具鏈。
