---
title: "程式碼庫配置"
description: "Wails v3 儲存庫的組織方式，以及各部分如何相互配合"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3 位於一個包含框架執行階段、CLI、範例、文件和建置工具鏈的<strong>單一儲存庫</strong>中。 本頁將逐一介紹深入研究內部實作時需要瞭解的<em>目錄結構</em>。

## 頂層概覽

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

接下來，我們將深入查看<strong>`v3/`</strong>目錄樹。

## `v3/`根目錄

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> 專案範本位於`internal/templates/`下（每個框架
>
> 技術棧各有一個資料夾，另加`base/`、`_common/`和`ios/`）。頂層沒有`v3/templates/`
>
> 目錄。

### 思維模型

1. <strong>`pkg/`</strong>公開<em>應用程式開發人員匯入的內容</em>\
2. <strong>`internal/`</strong>包含<em>底層功能的實作方式</em>\
3. <strong>`cmd/wails3`</strong>驅動<em>專案生命週期與建置流程</em>\

其他所有部分都為這三大支柱提供支援。

---

## `cmd/` – 命令

| 路徑 | 備註 |
| --- | --- |
| `v3/cmd/wails3` | **CLI 進入點**。精簡的`main.go`會將所有邏輯委派給`internal/commands`中的套件。 |
| `internal/commands/*` | 子命令（init、dev、build、doctor 等）。每個子命令各自位於單獨的檔案中，方便查找。 |
| `internal/commands/task_wrapper.go` | 銜接 CLI 旗標與 Taskfile 建置管線。 |

CLI 負責：

- **專案鷹架建立**（`init`、產生範本）\
- **開發伺服器協調**（`dev`、即時重新載入）\
- **正式環境建置與封裝**（`build`、`package`、平台包裝程式）\
- **診斷**（`doctor`）\

---

## `internal/` – 核心機房

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### 主要子套件

| 套件 | 職責 | 連接位置 |
| --- | --- | --- |
| `runtime` | 包含精簡的`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`建置標籤黏合程式碼，以及`runtime/desktop/`下的嵌入式 JS 執行階段。各作業系統實際使用的視窗、剪貼簿、對話方塊和系統匣程式碼位於`pkg/application/*_{darwin,linux,windows}.go`。 | 透過`pkg/application`間接匯入。 |
| `assetserver` | 雙模式檔案伺服器：<br />• 開發：從磁碟提供檔案並代理 Vite（`build_dev.go`）<br />• 正式環境：透過`go:embed`嵌入資產（`build_production.go`） | 啟動期間由`pkg/application`初始化。 |
| `generator` | 剖析 Go 原始碼以建立<strong>繫結中繼資料</strong>，之後再由此產生 TypeScript/JS 虛設常式檔案和事件常數。進入點：`collect/`與`render/`之上的`generator.Generate`／`generator.Generator`。 | 由`wails3 generate bindings`觸發。 |
| `packager` | 用於產生 Linux `deb`／`rpm`／`archlinux`成品的`nfpm`包裝程式（由`internal/commands/`下的`myapp.DEB`／`.RPM`／`.ARCHLINUX` nfpm 設定驅動）。 | 由`wails3 tool package`叫用。macOS DMG／Windows MSIX 位於`internal/commands/{dmg,msix.go,webview2/}`下。 |

輔助工具（例如`s/`、`hash/`、`flags/`）可讓各項內部關切事項保持解耦。

---

## `pkg/` – 公開 API

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> 不存在`pkg/runtime/`、`pkg/options/`或`pkg/menu/`套件。視窗／選單
>
> 選項與`pkg/application`位於同一處（例如`WebviewWindowOptions`、`Menu`、
>
> `MenuItem`），而`assetserver/`位於`internal/`下。

`pkg/application`會啟動 Wails 程式：

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

其底層會：

1. 連接`internal/runtime`建置標籤黏合程式碼與`pkg/application/`中各作業系統專用的程式碼
2. 設定`internal/assetserver`執行個體
3. 註冊所有由繫結驅動的訊息處理器
4. 進入作業系統主執行緒

---

## `internal/templates/` – 鷹架藍圖

`internal/templates/`隨附<strong>基礎範本</strong>（Go 配置位於`base/`、 `_common/`、`ios/`下）和<strong>前端外觀範本</strong>（`vanilla[-ts]`、`react[-ts]`、 `react-swc[-ts]`、`lit[-ts]`、`preact[-ts]`、`qwik[-ts]`、`solid[-ts]`、 `svelte[-ts]`、`sveltekit[-ts]`、`vue[-ts]`）。

在`wails3 init -t react`時，CLI 會：

1. 複製`_common` Go 檔案
2. 合併所需的前端套件
3. 執行`go mod tidy`（可使用`--skipgomodtidy`略過）

編輯範本<strong>不會</strong>影響現有應用程式，只會影響日後的`init`。公開範例位於`v3/examples/`下；這些範例不能取代貢獻者文件中所述的自動化測試套件。

---

## `tasks/` – 發行自動化

Taskfile 封裝了複雜的跨平台編譯、版本更新和變更記錄產生作業。`internal/commands/task.go`會以程式化方式使用這些 Taskfile，因此<strong>CLI</strong>和<strong>CI</strong>採用相同的邏輯。

---

## 各部分如何互動

```d2
direction: down
CLI: wails3 CLI
Generator: internal/generator
AssetDev: assetserver（開發）
Packager: internal/packager
AppRuntime: {
  label: 應用程式執行階段
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: 作業系統 API
}
CLI -> Generator: 建置／產生
CLI -> AssetDev: 開發
CLI -> Packager: 封裝
Generator -> ApplicationPkg: 繫結
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

<em>CLI → 產生器 → 執行階段</em>構成從<strong>原始碼</strong>到<strong>執行中的桌面應用程式</strong>的核心路徑。

---

## 快速導覽提示

| 需要瞭解… | 請查看… |
| --- | --- |
| 平台轉接層 | `pkg/application/*_darwin.go`、`*_linux.go`、`*_windows.go`（視窗、剪貼簿、對話方塊、系統匣、主執行緒、events_common）。Linux cgo：`pkg/application/linux_cgo*.go`。 |
| 橋接通訊協定 | `pkg/application/messageprocessor*.go` |
| 資產工作流程 | `internal/assetserver/`（`build_dev.go`與`build_production.go`） |
| 封裝流程 | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`、`internal/packager/` |
| 範本引擎 | `internal/templates/`（`templates.Install`、`templates.GetDefaultTemplates`） |
| 靜態分析 | `internal/generator/{generate.go,collect/,render/}` |

---

現在你已經掌握儲存庫的<strong>概念藍圖</strong>。搭配`ripgrep`、IDE 的「前往檔案／符號」功能和範例應用程式，即可深入探索任何功能。祝開發愉快！
