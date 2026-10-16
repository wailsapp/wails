---
title: "範本系統"
description: "Wails v3 如何為新專案建立骨架、範本的組織方式，以及如何建立自己的範本。"
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails 隨附<strong>範本系統</strong>，可讓`wails3 init`產生可直接執行的專案。內建範本刻意只涵蓋少數框架（Vanilla、React、Vue、Svelte）；若要使用其他框架，可以[自備前端](/guides/dev/frontend-frameworks/)或發佈[自訂範本](/guides/advanced/custom-templates/)。

本頁涵蓋：

1. 範本目錄配置
2. CLI 如何選擇及轉譯範本
3. 逐步建立新範本
4. 更新或覆寫現有範本
5. 疑難排解與最佳實務

---

## 1. 範本的存放位置

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — 合併至每個專案的通用樣板（Taskfile、`build/`目錄及共用基礎架構）。
- **`base/`** — 每個範本都以此 Go 端為起點。請注意：`base/`本身<strong>不</strong>包含`template.json`；該檔案位於各框架專用的範本中。
- **框架資料夾** — 包含前端（`frontend/`）、框架設定，以及用來描述範本中繼資料的`template.json`。
- 資料夾名稱會與傳給 CLI 的<strong>範本 ID</strong>（`wails3 init -t react`）相符。
- <strong>語言慣例：</strong>TypeScript 是預設語言，並使用不帶後綴的名稱（`react`）；若有 JavaScript 變體，其名稱會加上`-js`後綴（`react-js`）。內建範本會在`template.yaml`中以`typescript: true|false`明確宣告其語言。社群範本仍可使用舊版`-ts`後綴，系統會將其視為備援依據。

> 整個`internal/templates/`目錄都會編譯進 CLI 二進位檔
>
> 這是透過`//go:embed *`完成，因此使用者可以離線建立專案骨架。

---

## 2. `wails3 init`如何使用範本

呼叫鏈（沒有`cmd/wails3/init.go` — CLI 直接在`cmd/wails3/main.go`中連接）：

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

沒有`Template.Load()`／`Template.CopyTo()`／`Template.Validate()` API；擷取作業由`gosod`（`github.com/leaanthony/gosod`）針對內嵌的`fs.FS`進行處理。

### `wails3 init`旗標

定義於`internal/flags/init.go`：

| 旗標 | 用途 | 預設值 |
| --- | --- | --- |
| `-p` | 套件名稱 | `main` |
| `-t` | 內建範本名稱、本機路徑或 URL | `vanilla` |
| `-n` | 專案名稱 | （空白） |
| `-d` | 專案目錄 | `.` |
| `-q` | 隱藏主控台輸出 | false |
| `-l` | 列出範本 | false |
| `-skipgomodtidy` | 擷取後略過執行`go mod tidy` | false |
| `-git` | 要初始化的 Git 儲存庫 URL | （空白） |
| `-mod` | Go 模組路徑（若未設定，則從`-git`衍生） | （空白） |
| `-s` | 使用遠端範本時略過警告 | false |
| `-productname`／`-productdescription`／`-productversion`／`-productcompany`／`-productcopyright`／`-productcomments`／`-productidentifier` | 內建於產生之建置資產中的中繼資料 | 合理的預設值 |

沒有<strong>`-list`長別名</strong>（只有`-l`），也沒有各範本專用的`--help`。

### 替換

預留位置是標準 Go 範本指令；開頭的`.`是欄位存取器的一部分：

| 預留位置 | 範例 | 來源 |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | `-n`旗標／目錄名稱 |
| `{{.ModulePath}}` | `github.com/me/myapp` | `-mod`旗標，或衍生自`-git` |
| `{{.WailsVersion}}` | `v3.0.0-…` | 來自`internal/version`的編譯時內嵌常數 |
| `{{.ProductName}}`、`{{.ProductDescription}}`、`{{.ProductVersion}}`、`{{.ProductCompany}}`、`{{.ProductCopyright}}`、`{{.ProductComments}}`、`{{.ProductIdentifier}}` | 建置時中繼資料 | 對應的`-product*`旗標 |

如果需要新的預留位置，請在`internal/templates/templates.go`的範本資料中新增欄位，並在`internal/flags/init.go`中新增相符的欄位／旗標（或從`internal/commands/init.go`設定）。

### 複製後掛鉤

`gosod`完成範本解壓縮後，CLI 會執行：

```
go mod tidy
```

除非傳入`-skipgomodtidy`。此處沒有`task deps`步驟。

---

## 3. 建立新範本

> 範例：新增<strong>Solid</strong>範本

### 3.1資料夾與 ID

```
internal/templates/solid/
```

資料夾名稱就是範本 ID。請維持<strong>kebab-case</strong>格式。

### 3.2最小檔案集

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

先複製`react`，再刪減檔案。別忘了編寫`template.yaml`；若為 TypeScript 範本，請設定`typescript: true`。只有`base/`資料夾沒有這個檔案。

### 3.3更新預留位置

搜尋範例常值，並將其取代為 Go 範本指令，例如：

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4整合範本

由於`templates.go`會在初始化時走訪嵌入式檔案系統，因此通常只要在`internal/templates/<id>/`下新增資料夾即可，不需要手動呼叫註冊函式。如果需要額外邏輯（自訂驗證、複製後步驟），請將其新增至`internal/templates/templates.go`中的`templates.Install`。

### 3.5測試

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

請確認：

- 開發伺服器會在`WAILS_VITE_PORT`中公布的連接埠上啟動
- 產生的繫結會出現在`frontend/bindings/...`下
- 熱重新載入正常運作

---

## 4. 修改現有範本

1. 編輯`internal/templates/<id>/`下的檔案。
2. 重新建置 CLI（`cd v3 && go build -o ../wails3 ./cmd/wails3`）；`//go:embed *`指令會納入新內容。
3. 更新`frontend/package.json`與`Taskfile.yml`中的<strong>相依套件版本</strong>。
4. 若行為有所變更，請更新範本的`template.json`說明。

### 常見調整

| 工作 | 位置 |
| --- | --- |
| 變更開發伺服器連接埠 | `frontend/vite.config.ts`——讀取`WAILS_VITE_PORT` |
| 新增環境變數 | `build/Taskfile.yml`或`frontend/.env` |
| 更換 JavaScript 套件管理器 | 在`build/Taskfile.yml`中將`npm`替換為`pnpm`／`bun` |

---

## 5. 範本編寫提示

- **保持前端通用性**——避免參照 Wails 專用的全域變數；執行階段會由資產伺服器提供`/wails/runtime.js`。
- **不要包含編譯產物**——從嵌入目錄排除`node_modules`、`dist`及`.DS_Store`（或透過`.gitignore`忽略這些檔案或目錄，確保永遠不會提交）。
- **記錄先決條件**——在`template.json`或`NEXTSTEPS.md`中記錄 Node 版本、額外的 CLI 工具等。
- **避免破壞性變更**——如果改動幅度很大，請建立新的範本 ID，不要直接變更現有範本。

---

## 6. 疑難排解

| 症狀 | 原因 | 解決方法 |
| --- | --- | --- |
| `unknown template name` | `-t` 中有拼寫錯誤，或範本未嵌入 | 執行 `wails3 init -l` 以列出可用的範本 |
| 預留位置未被取代 | 使用了 `{{ProjectName}}`，而不是 `{{.ProjectName}}` | 加入開頭的 `.`（Go 範本欄位存取語法） |
| 開發伺服器開啟空白頁面 | Vite 設定未讀取 `WAILS_VITE_PORT` | 檢查您的 `vite.config.ts` |
| 前端在正式環境中建置失敗 | 忘記設定 Vite 的 `base` 路徑 | 在 `vite.config.ts` 中設定 `base: "./"` |

---

## 7. 主要原始碼檔案對照表

| 檔案 | 職責 |
| --- | --- |
| `internal/templates/templates.go` | 嵌入範本檔案系統，並公開 `Install(options *flags.Init) error`、`GetDefaultTemplates()`、`ValidTemplateName(name)` |
| `internal/templates/<id>/**` | 實際的範本內容 |
| `internal/commands/init.go` | CLI 銜接程式碼：選取範本、填入中繼資料，並呼叫 `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — 將現有專案<em>匯出</em>回範本的工具（更新時很方便） |
| `internal/flags/init.go` | `wails3 init` 的旗標定義 |

---

## 8. 重點回顧

- 範本位於 **`internal/templates/`**，並透過 `//go:embed *` 內建至 CLI。
- `wails3 init -t <id>` 透過 `gosod` 解出範本，並執行 `go mod tidy`（可使用 `-skipgomodtidy` 跳過）。
- 建立範本非常簡單：**建立資料夾**、加入檔案和 `template.json`，再使用 `{{.ProjectName}}` 形式的預留位置即可。
- 此系統<strong>可擴充</strong>且<strong>自成一體</strong>，非常適合與團隊或社群分享自訂技術堆疊。

祝您範本製作順利！
