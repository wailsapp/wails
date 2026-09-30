---
title: "CLI 參考資料"
description: "Wails CLI 命令的完整參考資料"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

Wails CLI 提供一套完整的命令，協助您開發、建置及維護 Wails 應用程式。

## 核心命令

核心命令是用於建立、開發及建置專案的主要命令。

所有 CLI 命令皆採用以下格式：`wails3 <command>`。

### `init`

初始化新的 Wails 專案。在初始化期間，會執行`go mod tidy`命令，將專案套件更新至最新狀態。您可以對`init`命令使用`-skipgomodtidy`旗標來略過此步驟。

```bash
wails3 init [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-p` | Go 套件名稱 | `main` |
| `-t` | 範本名稱或 URL | `vanilla` |
| `-n` | 專案名稱 |  |
| `-d` | 專案目錄 | `.` |
| `-q` | 隱藏輸出 | `false` |
| `-l` | 列出範本 | `false` |
| `-mod` | Go 模組路徑（若省略，則根據`-git`計算） |  |
| `-git` | Git 儲存庫 URL |  |
| `-s` | 使用遠端範本時略過警告 | `false` |
| `-productname` | 產品名稱 | `My Product` |
| `-productdescription` | 產品說明 | `My Product Description` |
| `-productversion` | 產品版本 | `0.1.0` |
| `-productcompany` | 公司名稱 | `My Company` |
| `-productcopyright` | 著作權聲明 | `© now, My Company` |
| `-productcomments` | 檔案註解 | `This is a comment` |
| `-productidentifier` | 產品識別碼 |  |
| `-skipgomodtidy` | 略過 go mod tidy | `false` |

`-git`旗標接受多種 Git URL 格式：

- HTTPS：`https://github.com/username/project`
- SSH：`git@github.com:username/project`或`ssh://git@github.com/username/project`
- Git 通訊協定：`git://github.com/username/project`
- 檔案系統：`file:///path/to/project.git`

提供此旗標後，系統將：

1. 在專案目錄中初始化 Git 儲存庫
2. 將指定的 URL 設為遠端 origin
3. 更新`go.mod`中的模組名稱，使其與儲存庫 URL 相符
4. 加入所有檔案

### `dev`

以開發模式執行應用程式。您可以即時檢視前端程式碼，並在執行中的應用程式內看到所做的變更，而無須重新建置整個應用程式。系統也會偵測 Go 程式碼的變更，並自動重新建置及啟動應用程式。

```bash
wails3 dev [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-config` | 設定檔路徑 | `./build/config.yml` |
| `-port` | Vite 開發伺服器連接埠 | `9245` |
| `-s` | 啟用 HTTPS | `false` |

@note{type="info"}
這相當於執行`wails3 task dev`，並會執行專案主要 Taskfile 中的`dev`工作。您可以編輯`Taskfile.yml`檔案來自訂此行為。

@end

### `build`

建置應用程式的偵錯版本。預設會針對目前的平台及架構進行建置。

```bash
wails3 build [flags] [CLI variables...]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-tags` | 額外的 Go 建置標籤（以逗號分隔） |  |

您可以傳入 CLI 變數來自訂建置：

```bash
wails3 build PLATFORM=linux CONFIG=production
```

使用`-tags`旗標傳入自訂的 Go 建置標籤：

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

標籤會以`EXTRA_TAGS`形式轉送至底層 Taskfile。

@note{type="info"}
這等同於執行`wails3 task build`，它會執行專案主要 Taskfile 中的`build`任務。傳給`build`的所有 CLI 變數都會轉送至底層任務。您可以編輯`Taskfile.yml`檔案來自訂建置流程。

@end

### `package`

建立供散發使用的平台專屬套件。

```bash
wails3 package [CLI variables...]
```

您可以傳入 CLI 變數來自訂封裝：

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

#### 套件類型

各平台可使用下列套件類型：

| 平台 | 套件類型 |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
這等同於`wails3 task package`，它會執行專案主要 Taskfile 中的`package`任務。傳給`package`的所有 CLI 變數都會轉送至底層任務。您可以編輯`Taskfile.yml`檔案來自訂封裝流程。

@end

### `task`

執行專案 Taskfile.yml 中定義的任務。這是[Taskfile](https://taskfile.dev)的內嵌版本，可讓您定義並執行自訂的建置、測試與部署任務。

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### CLI 變數

您可以使用`KEY=VALUE`格式將變數傳給任務：

```bash
wails3 task build PLATFORM=linux CONFIG=production
wails3 task deploy ENV=staging VERSION=1.2.3
```

您可以在 Taskfile.yml 中使用 Go 範本語法存取這些變數：

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-h` | 顯示 Task 使用方式 | `false` |
| `-i` | 建立新的 Taskfile.yml | `false` |
| `-list` | 列出含說明的任務 | `false` |
| `-list-all` | 列出所有任務（無論是否有說明） | `false` |
| `-json` | 將任務清單格式化為 JSON | `false` |
| `-status` | 若任務不是最新狀態，則以非零狀態碼結束 | `false` |
| `-f` | 即使任務已是最新狀態，仍強制執行 | `false` |
| `-w` | 為指定任務啟用監看模式 | `false` |
| `-v` | 啟用詳細輸出模式 | `false` |
| `-version` | 顯示 Task 版本 | `false` |
| `-s` | 停用命令回顯 | `false` |
| `-p` | 平行執行任務 | `false` |
| `-dry` | 編譯並顯示任務，但不執行 | `false` |
| `-summary` | 顯示任務摘要 | `false` |
| `-x` | 直接傳回任務的結束代碼 | `false` |
| `-dir` | 設定執行目錄 |  |
| `-taskfile` | 選擇要執行的 Taskfile |  |
| `-output` | 設定輸出樣式：[interleaved|group|prefixed] |  |
| `-c` | 彩色輸出（預設啟用） | `true` |
| `-C` | 限制可同時執行的工作數量 |  |
| `-interval` | 監看變更的間隔（秒） |  |

#### 範例

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build PLATFORM=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

啟動 Wails 專案 MCP 伺服器，以便透過代理程式協助管理專案。此伺服器與編譯至執行中應用程式內的 MCP 伺服器不同：`wails3 mcp`管理專案檔案和生命週期命令，而應用程式 MCP 伺服器則控制執行中的 WebView。

```bash
wails3 mcp [flags]
```

傳輸方式會自動選取：

- 當 MCP 主機以管線連接的標準輸入／輸出啟動 Wails 時，伺服器會使用<strong>stdio</strong>。
- 在終端機中以互動方式執行時，伺服器會在`127.0.0.1`上使用<strong>Streamable HTTP</strong>，並向作業系統要求一個可用的連接埠。

使用`--stdio`或`--http`明確選取傳輸方式。在 HTTP 模式下，使用`--port 0`選擇可用的迴路連接埠。

#### MCP 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `--root` | 允許的專案根目錄。位於其外的路徑和符號連結會遭拒絕。 | 目前目錄 |
| `--token` | 供變更資料與控制程序的工具使用的工作階段／Bearer 權杖。若未設定，則改用`WAILS_MCP_TOKEN`。 | 安全地產生 |
| `--stdio` | 強制使用 stdio 傳輸。 | 自動 |
| `--http` | 強制使用 Streamable HTTP 傳輸。 | 自動 |
| `--port` | HTTP 連接埠；`0`會選取本機回送位址上可用的連接埠。 | `0` |

在 HTTP 模式下，Wails 會將端點和 Bearer 權杖輸出至 stderr。在 stdio 模式下，權杖會包含在 MCP 初始化指示中。伺服器不會提供任意 Shell 執行功能。遠端範本和 Git 遠端儲存庫必須透過工具的`allowExternal`輸入明確核准。

### `doctor`

執行系統檢查並顯示狀態報告。

```bash
wails3 doctor
```

## 產生命令

產生命令可協助建立繫結、圖示和建置檔案等各種專案資產。所有產生命令都使用基礎命令：`wails3 generate <command>`。

### `generate bindings`

為 Go 程式碼產生繫結和模型。

```bash
wails3 generate bindings [flags] [patterns...]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-f` | 其他 Go 建置旗標 |  |
| `-d` | 輸出目錄 | `frontend/bindings` |
| `-models` | 模型檔案名稱 | `models` |
| `-index` | 索引檔案名稱 | `index` |
| `-ts` | 產生 TypeScript | `false` |
| `-i` | 使用 TypeScript 介面 | `false` |
| `-b` | 使用隨附的執行階段 | `false` |
| `-names` | 使用名稱而非 ID | `false` |
| `-noindex` | 略過索引檔案 | `false` |
| `-noevents` | 略過產生事件相關繫結 | `false` |
| `-dry` | 試執行 | `false` |
| `-silent` | 靜默模式 | `false` |
| `-v` | 偵錯輸出 | `false` |
| `-clean` | 產生前清理輸出目錄 | `true` |

### `generate build-assets`

為您的應用程式產生建置資產。

```bash
wails3 generate build-assets [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-name` | 專案名稱 |  |
| `-dir` | 輸出目錄 | `build` |
| `-silent` | 隱藏輸出 | `false` |
| `-company` | 公司名稱 |  |
| `-productname` | 產品名稱 |  |
| `-description` | 產品說明 |  |
| `-version` | 產品版本 |  |
| `-identifier` | 產品識別碼 | `com.wails.[name]` |
| `-copyright` | 著作權聲明 |  |
| `-comments` | 檔案註解 |  |

### `generate icons`

產生應用程式圖示。

```bash
wails3 generate icons [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-input` | 輸入 PNG 檔案 | 必填 |
| `-windowsfilename` | Windows 輸出檔名 |  |
| `-macfilename` | macOS 輸出檔名 |  |
| `-sizes` | 圖示大小（以逗號分隔） | `256,128,64,48,32,16` |
| `-example` | 產生範例圖示 | `false` |
| `-iconcomposerinput` | 輸入 Icon Composer 檔案（`.icon`） |  |
| `-macassetdir` | Mac 資產的輸出目錄（Assets.car + icns） |  |

#### Icon Composer（macOS）

在 macOS 26+ 上，您可以使用 Icon Composer `.icon` 檔案產生 `Assets.car` 和 `icons.icns`：

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

這會使用 Apple 的 `actool` 命令編譯 `.icon` 檔案。需要安裝 Xcode，且 `actool` 版本須為 26 或更新版本。

使用 Icon Composer 時，請將 `build/config.yml` 中的 `cfBundleIconName` 設為與 `.icon` 檔名（不含副檔名）相同：

```yaml
info:
  cfBundleIconName: "appicon"
```

若未設定且 `Assets.car` 存在，則預設為 `"appicon"`。

### `generate syso`

產生 Windows .syso 檔案。

```bash
wails3 generate syso [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-manifest` | 資訊清單檔案的路徑 | 必填 |
| `-icon` | 圖示檔案的路徑 | 必填 |
| `-info` | 版本資訊檔案的路徑 |  |
| `-arch` | 目標架構 | 目前的 GOARCH |
| `-out` | 輸出檔名 | `rsrc_windows_[arch].syso` |

### `generate .desktop`

產生 Linux .desktop 檔案。

```bash
wails3 generate .desktop [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-name` | 應用程式名稱 | 必填 |
| `-exec` | 可執行檔路徑 | 必填 |
| `-icon` | 圖示路徑 |  |
| `-categories` | 應用程式類別 | `Utility` |
| `-comment` | 應用程式註解 |  |
| `-terminal` | 在終端機中執行 | `false` |
| `-keywords` | 搜尋關鍵字 |  |
| `-version` | 應用程式版本 |  |
| `-genericname` | 通用名稱 |  |
| `-startupnotify` | 顯示啟動通知 | `false` |
| `-mimetype` | 支援的 MIME 類型 |  |
| `-output` | 輸出檔名 | `[name].desktop` |

### `generate runtime`

產生預先建置的執行階段版本。

```bash
wails3 generate runtime
```

### `generate constants`

從 Go 程式碼產生 JavaScript 常數。

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

產生供散布使用的 Windows WebView2 啟動安裝程式。

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

建立新專案範本目錄的初始架構。

```bash
wails3 generate template [flags]
```

### `generate appimage`

產生 Linux AppImage。

```bash
wails3 generate appimage [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-binary` | 二進位檔的路徑 | 必填 |
| `-icon` | 圖示檔案的路徑 | 必填 |
| `-desktop` | .desktop 檔案的路徑 | 必填 |
| `-builddir` | 建置目錄 | 暫存目錄 |
| `-output` | 輸出目錄 | `.` |

## 服務命令

服務命令可協助管理 Wails 服務。所有服務命令都使用基礎命令：`wails3 service <command>`。

### `service init`

初始化新服務。

```bash
wails3 service init [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-n` | 服務名稱 | `example_service` |
| `-d` | 服務說明 | `Example service` |
| `-p` | 套件名稱 |  |
| `-o` | 輸出目錄 | `.` |
| `-q` | 隱藏輸出 | `false` |
| `-a` | 作者名稱 |  |
| `-v` | 版本 |  |
| `-w` | 網站 URL |  |
| `-r` | 儲存庫 URL |  |
| `-l` | 授權條款 |  |

## 工具命令

工具命令提供開發與偵錯用的公用工具。所有工具命令都使用基礎命令：`wails3 tool <command>`。

### `tool checkport`

檢查連接埠是否開啟。適合用來測試 vite 是否正在執行。

```bash
wails3 tool checkport [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-port` | 要檢查的連接埠 | `9245` |
| `-host` | 要檢查的主機 | `localhost` |

### `tool watcher`

監看檔案，並在檔案變更時執行命令。

```bash
wails3 tool watcher [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-config` | 設定檔路徑 | `./build/config.yml` |
| `-ignore` | 要忽略的模式 |  |
| `-include` | 要包含的模式 |  |

### `tool cp`

複製檔案。

```bash
wails3 tool cp
```

### `tool buildinfo`

顯示應用程式的建置資訊。

```bash
wails3 tool buildinfo
```

### `tool version`

根據提供的旗標遞增語意化版本。

```bash
wails3 tool version [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-v` | 要遞增的目前版本 |  |
| `-major` | 提升主要版本 | `false` |
| `-minor` | 提升次要版本 | `false` |
| `-patch` | 提升修補版本 | `false` |
| `-prerelease` | 提升預發行版本（例如從 alpha.5 提升至 alpha.6） | `false` |

此命令遵循以下優先順序：主要版本 > 次要版本 > 修補版本 > 預發行版本。若輸入版本含有「v」前綴，則會予以保留，並同時保留所有預發行和中繼資料組成部分。

用法範例：

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.0.0-alpha.5 -prerelease # Output: v3.0.0-alpha.6
```

### `tool package`

產生 Linux 套件（deb、rpm、archlinux）。

```bash
wails3 tool package [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-format` | 套件格式（deb、rpm、archlinux） | `deb` |
| `-name` | 可執行檔名稱 | `myapp` |
| `-config` | 設定檔路徑 |  |
| `-out` | 輸出目錄 | `.` |

### `tool lipo`

合併各架構專用的二進位檔，以建立 macOS 通用二進位檔。

```bash
wails3 tool lipo [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-output` | 輸出二進位檔路徑 |  |

### `tool capabilities`

檢查系統的建置能力（Linux 上是否有 GTK4/GTK3 可用）。

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

產生用於交叉編譯的 Docker 磁碟區掛載旗標。輸出 Go 模組快取的 `-v` 旗標，以及 `go.mod` 中所有本機 `replace` 指示詞的對應旗標，以供 Taskfile 的 `docker run` 命令使用。

```bash
wails3 tool docker-mounts
```

### `tool has`

檢查工具或功能是否可用，並將 `true` 或 `false` 輸出至標準輸出。此命令設計用於 Taskfile 的 `sh:` 變數，可作為 `command -v` 的跨平台替代方案。

使用 `|` 檢查多個替代項目中是否有任何一個可用。

```bash
wails3 tool has <tool>
```

#### 範例

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### 在 Taskfile 中使用

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="已棄用"}
`wails3 tool has-cc` 已棄用。請更新 Taskfile，改用 `wails3 tool has gcc|clang`。

@end

`wails3 tool has gcc|clang` 的向後相容別名。檢查 PATH 中是否有 `gcc` 或 `clang` 可用，並輸出 `true` 或 `false`。

```bash
wails3 tool has-cc
```

## 更新命令

更新命令可協助管理及更新專案資產。所有更新命令都使用以下基礎命令：`wails3 update <command>`。

### `update cli`

將 Wails CLI 更新至新版本。

```bash
wails3 update cli [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-pre` | 更新至最新預發行版本 | `false` |
| `-version` | 更新至指定版本 |  |
| `-nocolour` | 停用彩色輸出 | `false` |

update cli 命令可讓你更新 Wails CLI 安裝。預設會更新至最新穩定版本。 你可以使用 `-pre` 旗標更新至最新預發行版本，或使用 `-version` 旗標指定特定版本。

更新後，請記得更新專案的 go.mod 檔案，以使用相同版本：

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

使用指定的設定檔更新建置資產。

```bash
wails3 update build-assets [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-config` | 設定檔路徑 |  |
| `-dir` | 輸出目錄 | `build` |
| `-silent` | 隱藏輸出 | `false` |
| `-company` | 公司名稱 |  |
| `-productname` | 產品名稱 |  |
| `-description` | 產品說明 |  |
| `-version` | 產品版本 |  |
| `-identifier` | 產品識別碼 |  |
| `-copyright` | 著作權聲明 |  |
| `-comments` | 檔案註解 |  |

## 公用程式命令

公用程式命令提供常見工作所需的便利捷徑。請直接搭配基礎命令使用這些命令：`wails3 <command>`。

### `docs`

在預設瀏覽器中開啟 Wails 文件。

```bash
wails3 docs
```

### `releasenotes`

顯示目前版本或指定版本的版本資訊。

```bash
wails3 releasenotes [flags]
```

#### 旗標

| 旗標 | 說明 | 預設值 |
| --- | --- | --- |
| `-v` | 要顯示其版本資訊的版本 |  |
| `-n` | 停用彩色輸出 | `false` |

### `version`

輸出目前的 Wails 版本。

```bash
wails3 version
```

### `sponsor`

在預設瀏覽器中開啟 Wails 贊助頁面。

```bash
wails3 sponsor

```
