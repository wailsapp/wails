---
title: "自訂建置"
description: "使用 Task 和 Taskfile.yml 自訂建置流程"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## 概觀

Wails 建置系統是一套靈活且強大的工具，旨在簡化 Wails 應用程式的建置流程。它採用工作執行器[Task](https://taskfile.dev)，讓您能輕鬆定義及執行工作。雖然 v3 建置系統是預設選項，但 Wails 鼓勵「自備工具」的做法，讓開發人員可依需求自訂建置流程。

如需進一步瞭解 Task 的使用方式，請參閱[官方文件](https://taskfile.dev/usage/)。

## Task：建置系統的核心

[Task](https://taskfile.dev) 是以 Go 編寫的現代化 Make 替代方案。它使用 YAML 檔案定義工作及其相依性。在 Wails 建置系統中，[Task](https://taskfile.dev)負責協調建置流程，是其中的核心。

主要的`Taskfile.yml`位於專案根目錄，而平台專屬工作則定義於`build/<platform>/Taskfile.yml`檔案中。`build`目錄中的共用`Taskfile.yml`檔案包含各平台共用的工作。

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

專案根目錄中的`Taskfile.yml`檔案是建置系統的主要進入點，其中定義了工作及其相依性。以下是預設的`Taskfile.yml`檔案：

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## 平台專屬 Taskfile

每個平台都有自己的 Taskfile，位於`build`目錄下的平台目錄中。這些檔案定義該平台的核心工作。每個 Taskfile 都會納入`build/Taskfile.yml`檔案中的共用工作。

### Windows

位置：`build/windows/Taskfile.yml`

Windows 專屬 Taskfile 包含在 Windows 上建置、封裝及執行應用程式的工作。主要功能包括：

- 使用選用的正式環境旗標進行建置
- 產生`.ico`圖示檔案
- 產生 Windows `.syso`檔案
- 建立用於封裝的 NSIS 安裝程式

### Linux

位置：`build/linux/Taskfile.yml`

Linux 專屬 Taskfile 包含在 Linux 上建置、封裝及執行應用程式的工作。主要功能包括：

- 使用選用的正式環境旗標進行建置
- 建立 AppImage、deb、rpm 及 Arch Linux 套件
- 產生供 Linux 應用程式使用的`.desktop`檔案

### macOS

位置：`build/darwin/Taskfile.yml`

macOS 專屬 Taskfile 包含在 macOS 上建置、封裝及執行應用程式的工作。主要功能包括：

- 為 amd64、arm64 及通用（兩者）架構建置二進位檔
- 產生`.icns`圖示檔案
- 建立用於散佈的`.app`套件組合
- 對`.app`套件組合進行臨時簽署
- 設定 macOS 專屬建置旗標和環境變數

## 工作執行與命令別名

`wails3 task`命令是[Taskfile](https://taskfile.dev)的內嵌版本，可執行`Taskfile.yml`中定義的工作。

`wails3 build`和`wails3 package`命令分別是`wails3 task build`和`wails3 task package`的別名。執行這些命令時，Wails 會在內部將其轉換為適當的工作執行命令：

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### 將參數傳遞給工作

您可以使用`KEY=VALUE`格式將 CLI 變數傳遞給工作。這些變數會透過別名命令轉送：

```bash
# These are equivalent:
wails3 build PLATFORM=linux CONFIG=production
wails3 task build PLATFORM=linux CONFIG=production

# Package with custom version:
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

在`Taskfile.yml`中，您可以使用 Go 範本語法存取這些變數：

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

### Go 建置自訂變數

產生的專案提供用來追加自訂 Go 建置標籤、連結器旗標及設定 CGO 的變數。標籤值必須是以逗號分隔的名稱，不包含 `-tags` 選項；Wails 會將它們與所選建置模式及平台所需的標籤合併。

| 變數 | 適用範圍 |
| --- | --- |
| `APP_TAGS` | 所有透過 Taskfile 執行的建置 |
| `APP_TAGS_LINUX` | Linux 建置 |
| `APP_TAGS_DARWIN` | macOS 建置 |
| `APP_TAGS_WINDOWS` | Windows 建置 |
| `APP_TAGS_ANDROID` | Android 建置 |
| `APP_TAGS_IOS` | iOS 建置 |
| `APP_TAGS_SERVER` | 伺服器模式建置 |
| `APP_LDFLAGS` | 為所有透過 Taskfile 執行的建置追加連結器旗標 |
| `APP_CGO_ENABLED` | 以 `0` 或 `1` 覆寫桌面及伺服器建置的 CGO 設定 |
| `EXTRA_TAGS` | 為單次呼叫追加標籤 |

例如：

```bash
wails3 build APP_TAGS=sqlite_fts5,netgo APP_TAGS_LINUX=myapp_linux
wails3 build EXTRA_TAGS=diagnostics
APP_LDFLAGS='-X example.com/myapp/internal/version.Value=1.2.3' wails3 build
```

若要設定持續使用的專案預設值，請在根目錄的 `Taskfile.yml` 中設定 `APP_*` 值。隨命令傳入的 `KEY=value` 優先順序最高，會在該次呼叫中取代變數的值。根 Taskfile 中的字面值優先於同名的處理程序環境變數。若保留下方所示的自我參照預設值運算式，則未提供命令列值時會使用處理程序環境。若只想為一次建置追加標籤而不取代持續使用的 `APP_TAGS` 值，請使用 `EXTRA_TAGS`。

當 `APP_CGO_ENABLED` 為空時，Linux 與 macOS 預設為 `1`，Windows 預設為 `0`，原生伺服器建置保留 Go 在主機上的預設值，伺服器 Docker 建置則預設為 `0`。Android 與 iOS 建置一律使用 CGO，並維持 `CGO_ENABLED=1`；`APP_CGO_ENABLED` 不會覆寫這些行動工具鏈的設定。以 Docker 執行的桌面交叉編譯與伺服器建置，會接收與對應原生 Taskfile 建置相同的適用 `APP_*` 值。

@note{type="info" title="現有專案"}
根目錄的 `Taskfile.yml` 由專案管理，因此 `wails3 update build-assets` 不會覆寫它。在引入這些變數之前建立的專案，必須手動將下列項目加入根 Taskfile 的 `vars` 區塊：

```yaml
vars:
  APP_TAGS: '{{.APP_TAGS | default ""}}'
  APP_TAGS_LINUX: '{{.APP_TAGS_LINUX | default ""}}'
  APP_TAGS_DARWIN: '{{.APP_TAGS_DARWIN | default ""}}'
  APP_TAGS_WINDOWS: '{{.APP_TAGS_WINDOWS | default ""}}'
  APP_TAGS_ANDROID: '{{.APP_TAGS_ANDROID | default ""}}'
  APP_TAGS_IOS: '{{.APP_TAGS_IOS | default ""}}'
  APP_TAGS_SERVER: '{{.APP_TAGS_SERVER | default ""}}'
  APP_LDFLAGS: '{{.APP_LDFLAGS | default ""}}'
  APP_CGO_ENABLED: '{{.APP_CGO_ENABLED | default ""}}'
```

將自我參照預設值運算式替換成字面值，即可設為持續使用的專案預設值。

@end

這些變數由透過 Taskfile 執行的建置使用。產生的 Xcode 專案在建置階段直接呼叫 Go，不會經過這套 Taskfile 自訂流程，因此從 Xcode 啟動的建置目前不會使用這些變數。

## 通用建置流程

在所有平台上，建置流程通常包含以下步驟：

1. 整理 Go 模組
2. 建置前端
3. 產生圖示
4. 使用平台專屬旗標編譯 Go 程式碼
5. 封裝應用程式（依平台而異）

## 自訂建置流程

雖然 v3 建置系統提供穩健的預設組態，您仍可輕鬆自訂，以符合專案需求。修改`Taskfile.yml`及平台專屬 Taskfile 後，您可以：

- 新增工作
- 修改現有工作
- 變更工作執行順序
- 與其他工具及指令碼整合

這種靈活性讓您能依特定需求調整建置流程，同時仍可受益於 Wails 建置系統所提供的架構。

@note{type="tip" title="學習 Taskfile"}
強烈建議閱讀[Taskfile](https://taskfile.dev)文件，以瞭解如何有效使用 Taskfile。執行`wails3 task --version`即可查詢 Wails CLI 內嵌的 Taskfile 版本。

@end

## 開發模式

Wails 建置系統包含功能強大的開發模式，可提供即時重新載入及熱模組替換，進而改善開發體驗。使用`wails3 dev`命令即可啟用此模式。

### 運作方式

執行`wails3 dev`時，會進行以下流程：

1. 此命令會檢查可用的連接埠；若未指定，預設為9245。
2. 它會為前端開發伺服器（Vite）設定環境變數。
3. 它會使用[refresh](https://github.com/atterpac/refresh)程式庫啟動檔案監看器。

[refresh](https://github.com/atterpac/refresh)程式庫負責監控檔案變更並觸發重新建置。它使用`./build/config.yml`檔案中`dev_mode`鍵下定義的組態。您可以將它設定為忽略特定目錄和檔案、決定要監看哪些檔案，以及偵測到變更時要採取哪些動作。預設組態已相當實用，但您可以依需求自由自訂。

### 組態

以下是其結構範例：

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

此組態檔可讓您：

- 設定檔案監看的根路徑
- 設定記錄層級
- 設定檔案變更事件的防彈跳時間
- 忽略特定目錄、檔案或副檔名
- 定義檔案變更時要執行的命令

### 自訂開發模式

您可以修改`config.yml`檔案中的這些值，自訂開發模式的使用體驗。

自訂方式包括：

1. 變更要監看的目錄或檔案
2. 調整防彈跳時間，以控制系統回應變更的速度
3. 新增或修改執行命令，以符合專案需求

### 使用瀏覽器進行開發

雖然 Wails v2 完整支援使用瀏覽器進行開發，但這造成了許多混淆。能在瀏覽器中運作的應用程式不一定能在桌面應用程式中運作，因為 WebView 並不提供所有瀏覽器 API。

進行以 UI 為主的開發工作時，v3 仍可讓您彈性地使用瀏覽器：在開發模式下存取`http://localhost:9245`所示的 Vite URL。這可讓您在調整樣式和版面配置時使用功能強大的瀏覽器開發工具。請注意，在此模式下，Go 繫結<em>無法運作</em>。準備測試繫結和事件等功能時，只需切換至桌面檢視，即可確保一切在生產環境中正常運作。
