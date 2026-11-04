---
title: "範本系統"
description: "Wails v3 如何建立新專案的基本架構、範本的組織方式，以及如何建立自己的範本。"
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails 隨附一套**範本系統**，可讓`wails3 init`產生可立即執行的專案。內建範本刻意只涵蓋少數框架（Vanilla、React、Vue、Svelte）；如要使用其他框架，可以[自備前端](/guides/dev/frontend-frameworks/)，或發布[自訂範本](/guides/advanced/custom-templates/)。

本頁涵蓋：

1. 範本目錄配置
2. CLI 如何選擇及算繪範本
3. 逐步建立新範本
4. 更新或覆寫現有範本
5. 疑難排解與最佳實務

---

## 1. 範本的存放位置

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Minimal frontend skeleton (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — 合併至每個專案的通用樣板（Taskfile、`build/`目錄及共用基礎設施）。
- **`base/`** — 用於自訂範本的最小前端骨架，並非自動產生的後端層。一般的專案工作副本會將所選框架範本疊加至`_common/`。
- **框架資料夾** — 包含前端（`frontend/`）、框架設定，以及說明範本中繼資料的`template.yaml`。
- 資料夾名稱須與傳給 CLI 的**範本 ID**（`wails3 init -t react`）相符。
- <strong>語言慣例：</strong>TypeScript 是預設語言，使用不含後綴的名稱（`react`）；若有 JavaScript 版本，則加上`-js`後綴（`react-js`）。內建範本會在`template.yaml`中使用`typescript: true|false`明確宣告其語言。社群範本仍可使用舊版`-ts`後綴，系統會將其視為備援判斷依據。

> 整個`internal/templates/`目錄會編譯至 CLI 二進位檔中
>
> 這是透過`//go:embed *`完成，因此使用者可以離線建立專案的基本架構。

---

載入器也接受舊版`template.json`中繼資料。新的內建範本使用 YAML，並包含明確的`typescript`設定。

## 2. `wails3 init`如何使用範本

呼叫鏈（沒有`cmd/wails3/init.go`，CLI 直接在`cmd/wails3/main.go`中接線）：

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

不存在`Template.Load()` / `Template.CopyTo()` / `Template.Validate()` API；工作副本由`gosod`（`github.com/leaanthony/gosod`）根據嵌入的`fs.FS`處理。

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
| `-skipgomodtidy` | 建立工作副本後，略過執行`go mod tidy` | false |
| `-git` | 要初始化的 Git 儲存庫 URL | （空白） |
| `-mod` | Go 模組路徑（若未設定，則衍生自`-git`） | （空白） |
| `-s` | 使用遠端範本時略過警告 | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | 嵌入所產生建置資產中的中繼資料 | 合理的預設值 |

沒有`-list`長別名（只有`-l`），也沒有各範本專屬的`--help`。

### 替換值

預留位置是標準的 Go 範本指令；開頭的`.`是欄位存取器的一部分：

| 預留位置 | 範例 | 來源 |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | `-n`旗標／目錄名稱 |
| `{{.ModulePath}}` | `github.com/me/myapp` | `-mod`旗標，或衍生自`-git` |
| `{{.WailsVersion}}` | `v3.0.0` | 來自`internal/version`的編譯時常數 |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | 建置時中繼資料 | 對應的`-product*`旗標 |

如需新的預留位置，請在`internal/templates/templates.go`的範本資料中新增欄位，並在`internal/flags/init.go`中新增相符的欄位／旗標（或透過`internal/commands/init.go`設定）。

### 複製後掛鉤

`gosod`完成建立範本工作副本後，CLI 會執行：

```
go mod tidy
```

除非傳入`-skipgomodtidy`。沒有`task deps`步驟。

---

## 3. 建立新範本

> 範例：新增**Solid**範本

### 3.1 資料夾與 ID

```
internal/templates/solid/
```

資料夾名稱 = 範本 ID。請使用**kebab-case**。

### 3.2 最小檔案集

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

先複製`react`，再刪除不需要的檔案。別忘了建立`template.yaml`；若是 TypeScript 範本，請設定`typescript: true`。`base/`是唯一沒有此檔案的資料夾。

### 3.3更新預留位置

搜尋範例常值，並將其取代為 Go 範本指令，例如：

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4接入範本

由於`templates.go`會在初始化時巡覽嵌入式檔案系統，因此通常只要在`internal/templates/<id>/`下新增資料夾即可，無須手動呼叫註冊函式。若需要額外邏輯（自訂驗證、複製後步驟），請將其加入`internal/templates/templates.go`中的`templates.Install`。

### 3.5測試

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

請確認：

- 開發伺服器會在`WAILS_VITE_PORT`所公布的連接埠上啟動
- 產生的繫結會出現在`frontend/bindings/...`下
- 熱重新載入可正常運作

---

## 4. 修改現有範本

1. 編輯`internal/templates/<id>/`下的檔案。
2. 重新建置 CLI（`cd v3 && go build -o ../wails3 ./cmd/wails3`）；`//go:embed *`指令會納入新內容。
3. 更新`frontend/package.json`和`Taskfile.yml`中的所有**相依套件版本**。
4. 若行為有變更，請更新範本的`template.yaml`說明。

### 常見調整

| 工作 | 位置 |
| --- | --- |
| 變更開發伺服器連接埠 | `frontend/vite.config.ts`—讀取`WAILS_VITE_PORT` |
| 新增環境變數 | `build/Taskfile.yml`或`frontend/.env` |
| 更換 JavaScript 套件管理員 | 在`build/Taskfile.yml`中將`npm`換成`pnpm`/`bun` |

---

## 5. 範本編寫技巧

- **保持前端通用**—避免參照 Wails 專屬的全域變數；執行階段會由資產伺服器提供`/wails/runtime.js`。
- **不要包含編譯產物**—從嵌入式目錄中排除`node_modules`、`dist`和`.DS_Store`（或將它們加入`.gitignore`，確保絕不提交）。
- **記錄必要條件**—在`template.yaml`或`NEXTSTEPS.md`中說明 Node 版本、額外的 CLI 工具等。
- **避免破壞性變更**—若大幅改造範本，請建立新的範本 ID，而不要修改現有範本。

---

## 6. 疑難排解

| 症狀 | 原因 | 修正方式 |
| --- | --- | --- |
| `unknown template name` | `-t`中有拼字錯誤，或範本未嵌入 | 執行`wails3 init -l`以列出可用範本 |
| 預留位置未被取代 | 使用了`{{ProjectName}}`，而不是`{{.ProjectName}}` | 加上開頭的`.`（Go 範本欄位存取） |
| 開發伺服器開啟空白頁面 | Vite 設定未讀取`WAILS_VITE_PORT` | 檢查`vite.config.ts` |
| 前端無法在生產環境中建置 | 遺漏 Vite 的`base`路徑 | 在`vite.config.ts`中設定`base: "./"` |

---

## 7. 主要原始碼檔案對照表

| 檔案 | 職責 |
| --- | --- |
| `internal/templates/templates.go` | 嵌入範本檔案系統，並公開`Install(options *flags.Init) error`、`GetDefaultTemplates()`和`ValidTemplateName(name)` |
| `internal/templates/<id>/**` | 實際的範本內容 |
| `internal/commands/init.go` | CLI 銜接程式碼：選取範本、填入中繼資料，並呼叫`templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template`—用來將即時專案*匯出*回範本的公用工具（更新時很方便） |
| `internal/flags/init.go` | `wails3 init`的旗標定義 |

---

## 8. 重點回顧

- 範本位於<strong>`internal/templates/`</strong>，並透過`//go:embed *`內建至 CLI。
- `wails3 init -t <id>`會透過`gosod`擷取範本，並執行`go mod tidy`（可使用`-skipgomodtidy`略過）。
- 建立範本很簡單：**建立資料夾**、加入檔案和`template.yaml`，再使用`{{.ProjectName}}`樣式的預留位置即可。
- 此系統**可擴充**且**自給自足**，非常適合與團隊或社群分享自訂技術堆疊。

祝你建立範本順利！
