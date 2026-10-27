---
title: "自訂建置"
description: "使用 Task 和 Taskfile.yml 自訂建置流程"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## 概觀

Wails 建置系統是一套靈活且強大的工具，旨在簡化 Wails 應用程式的建置流程。它採用任務執行器[Task](https://taskfile.dev)，讓您能輕鬆定義及執行任務。雖然v3建置系統是預設選項，但 Wails 鼓勵「自備工具」的做法，讓開發人員可視需要自訂建置流程。

如需進一步瞭解 Task 的用法，請參閱[官方文件](https://taskfile.dev/usage/)。

## Task：建置系統的核心

[Task](https://taskfile.dev)是以 Go 撰寫的現代化 Make 替代方案。它使用 YAML 檔案定義任務及其相依性。在 Wails 建置系統中，[Task](https://taskfile.dev)負責統籌建置流程，扮演核心角色。

主要的`Taskfile.yml`位於專案根目錄，而各平台專屬的任務則定義於`build/<platform>/Taskfile.yml`檔案中。`build`目錄內的共用`Taskfile.yml`檔案包含各平台共享的通用任務。

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

專案根目錄中的`Taskfile.yml`檔案是建置系統的主要進入點，其中定義了任務及其相依性。以下是預設的`Taskfile.yml`檔案：

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

每個平台都有自己的 Taskfile，位於`build`目錄下的各平台目錄中。這些檔案定義該平台的核心任務。每個 Taskfile 都會納入`build/Taskfile.yml`檔案中的通用任務。

### Windows

位置：`build/windows/Taskfile.yml`

Windows 專屬 Taskfile 包含在 Windows 上建置、封裝及執行應用程式的任務。主要功能包括：

- 使用選用的正式環境旗標進行建置
- 產生`.ico`圖示檔案
- 產生 Windows `.syso`檔案
- 建立用於封裝的 NSIS 安裝程式

### Linux

位置：`build/linux/Taskfile.yml`

Linux 專屬 Taskfile 包含在 Linux 上建置、封裝及執行應用程式的任務。主要功能包括：

- 使用選用的正式環境旗標進行建置
- 建立 AppImage、deb、rpm 及 Arch Linux 套件
- 為 Linux 應用程式產生`.desktop`檔案

### macOS

位置：`build/darwin/Taskfile.yml`

macOS 專屬 Taskfile 包含在 macOS 上建置、封裝及執行應用程式的任務。主要功能包括：

- 建置適用於 amd64、arm64 及通用（兩者皆支援）架構的二進位檔
- 產生`.icns`圖示檔案
- 建立用於發佈的`.app`套件組合
- 對`.app`套件組合進行臨時簽署
- 設定 macOS 專屬的建置旗標與環境變數

## 執行任務與命令別名

`wails3 task`命令是[Taskfile](https://taskfile.dev)的內嵌版本，用於執行您在`Taskfile.yml`中定義的任務。

`wails3 build`和`wails3 package`命令分別是`wails3 task build`和`wails3 task package`的別名。執行這些命令時，Wails 會在內部將其轉換為對應的任務執行命令：

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### 將參數傳遞給任務

您可以使用`KEY=VALUE`格式將 CLI 變數傳遞給任務。這些變數會透過別名命令轉送：

```bash
# These are equivalent:
wails3 build GOOS=linux
wails3 task build GOOS=linux
```

內建任務會使用`GOOS`、`ARCH`和`DEV`。`VERSION`、`OUTPUT`、`PLATFORM`或`CONFIG`等變數，只有在您自己的 Taskfile 讀取它們時才會影響任務。以下是自訂 Taskfile 的範例，並非 Wails 內建的建置任務：

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

如果您的自訂封裝工作會讀取`VERSION`和`OUTPUT`，則可用相同方式傳入這些變數：

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

## 通用建置流程

在所有平台上，建置流程通常包含下列步驟：

1. 整理 Go 模組
2. 建置前端
3. 產生圖示
4. 使用平台專屬旗標編譯 Go 程式碼
5. 封裝應用程式（依平台而異）

## 自訂建置流程

雖然v3建置系統提供穩健的預設設定，但您可以輕鬆加以自訂，以符合專案需求。修改`Taskfile.yml`和平台專屬 Taskfile 後，您可以：

- 新增任務
- 修改現有任務
- 變更任務執行順序
- 與其他工具及指令碼整合

這種靈活性讓您能依具體需求調整建置流程，同時繼續受益於 Wails 建置系統所提供的架構。

@note{type="tip" title="學習 Taskfile"}
強烈建議閱讀[Taskfile](https://taskfile.dev)文件，以瞭解如何有效運用 Taskfile。執行`wails3 task --version`即可查詢 Wails CLI 內嵌的 Taskfile 版本。

@end

## 開發模式

Wails 建置系統包含功能強大的開發模式，可提供即時重新載入與熱模組替換，改善開發人員體驗。使用`wails3 dev`命令即可啟用此模式。

### 運作方式

執行`wails3 dev`時，會進行下列流程：

1. 命令會檢查可用的連接埠；若未指定，預設使用9245。
2. 設定前端開發伺服器（Vite）的環境變數。
3. 使用[refresh](https://github.com/atterpac/refresh)程式庫啟動檔案監看器。

[refresh](https://github.com/atterpac/refresh)程式庫負責監控檔案變更並觸發重新建置。它使用`./build/config.yml`檔案中`dev_mode`鍵下定義的設定。您可以將其設定為忽略特定目錄與檔案、決定要監看的檔案，以及偵測到變更時要採取的動作。預設設定已相當實用，但您仍可依需求自行調整。

### 設定

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

此設定檔可讓您：

- 設定檔案監看的根路徑
- 設定記錄層級
- 設定檔案變更事件的防彈跳時間
- 忽略特定目錄、檔案或副檔名
- 定義檔案變更時要執行的命令

### 自訂開發模式

您可以修改`config.yml`檔案中的這些值，自訂開發模式的使用體驗。

可採用的自訂方式包括：

1. 變更要監看變更的目錄或檔案
2. 調整防抖時間，以控制系統回應變更的速度
3. 新增或修改執行指令，以符合專案需求

### 使用瀏覽器進行開發

雖然 Wails v2完全支援使用瀏覽器進行開發，但這造成了許多混淆。能在瀏覽器中運作的應用程式不一定能在桌面應用程式中運作，因為 WebView 並未提供所有瀏覽器 API。

進行以 UI 為主的開發工作時，您仍可在v3中靈活使用瀏覽器：在開發模式下存取位於`http://localhost:9245`的 Vite URL。這樣一來，您在處理樣式和版面配置時便可使用功能強大的瀏覽器開發工具。請注意，在此模式下，Go 繫結*將無法運作*。準備測試繫結和事件等功能時，只需切換至桌面檢視，即可確保所有功能在生產環境中正常運作。
