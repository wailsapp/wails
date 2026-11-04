---
title: "您的第一個應用程式"
description: "逐步建立您的第一個 Wails 桌面應用程式"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

本指南將說明如何建立您的第一個 Wails v3應用程式，涵蓋專案設定、建置及開發工作流程。

<br/>

<br/>

@steps
### 建立新專案
開啟終端機並執行下列命令，以建立新的 Wails 專案：

```bash
wails3 init -n myfirstapp -t vanilla-js
```

本課程會明確選用 JavaScript Vanilla 範本。基本的`vanilla`範本使用 TypeScript。此命令會建立`myfirstapp`，其中包含所有必要檔案。

   <video src="/assets/wails_init.mp4" controls></video>

### 探索專案結構
前往`myfirstapp`目錄。您會看到數個檔案與資料夾：

@filetree
- build/           包含建置流程所使用的檔案
  - appicon.png  應用程式圖示
  - config.yml   建置設定
  - Taskfile.yml Build tasks
  - darwin/      macOS 專用建置檔案
    - Info.dev.plist Development configuration
    - Info.plist    正式環境設定
    - Taskfile.yml  macOS 建置工作
    - icons.icns    macOS 應用程式圖示
  - linux/       Linux 專用建置檔案
    - Taskfile.yml  Linux 建置工作
    - appimage/     AppImage 封裝
      - build.sh  AppImage 建置指令碼
    - nfpm/        NFPM 封裝
      - nfpm.yaml Package configuration
      - scripts/  建置指令碼
  - windows/     Windows 專用建置檔案
    - Taskfile.yml        Windows 建置工作
    - icon.ico           Windows 應用程式圖示
    - info.json          應用程式中繼資料
    - wails.exe.manifest Windows manifest file
    - nsis/              NSIS 安裝程式檔案
      - project.nsi                    NSIS 專案檔案
      - wails_tools.nsh               NSIS 輔助指令碼
- frontend/        前端應用程式檔案
  - index.html   主要 HTML 檔案
  - src/         前端原始碼
    - main.js    主要 JavaScript 檔案
  - package.json NPM package configuration
  - public/      靜態資源
  - Inter Font License.txt Font license
- .gitignore      Git 忽略檔案
- README.md       專案文件
- Taskfile.yml    專案工作
- go.mod          Go 模組檔案
- go.sum          Go 模組總和檢查碼
- greetservice.go Greeting service
- main.go         應用程式主要程式碼
@end

花點時間探索這些檔案並熟悉其結構。

@note{type="info"}
雖然 Wails v3預設使用[Task](https://taskfile.dev/)作為建置系統，但您仍可使用`make`或任何其他替代建置系統。

@end

### 建置應用程式
若要建置應用程式，請執行：

```bash
wails3 build
```

使用產生的 Taskfile 時，此命令會在`bin`中建立正式環境組建。開發模式請使用`wails3 dev`；若要明確建立開發組建，請使用`wails3 build DEV=true`。

@note{type="info"}
`wails3 build`是`wails3 task build`的簡寫，會執行`Taskfile.yml`中的`build`工作。

@end

     <video src="/assets/wails_build.mp4" controls></video>

建置完成後，您可以像執行一般應用程式一樣執行它：

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

您會看到入門使用者介面。正式環境組建不保證開發人員檢查器或主控台的行為；偵錯時請使用開發模式。

### 開發模式
您也可以在開發模式下執行應用程式。此模式可讓您變更前端程式碼，並立即在執行中的應用程式看到變更，無須重新建置整個應用程式。

1. 開啟新的終端機視窗。
2. 執行`wails3 dev`。應用程式會以偵錯模式編譯並執行。
3. 使用您偏好的編輯器開啟`frontend/index.html`。
4. 將名稱輸入欄位的`placeholder="Your name"`變更為`placeholder="Your full name"`，然後儲存檔案。
5. 儲存檔案。

這項變更會立即反映在應用程式中。

對後端程式碼所做的任何變更都會觸發重新建置：

1. 開啟`greetservice.go`。
2. 將包含`return "Hello " + name + "!"`的程式碼行變更為`return "Hello there " + name + "!"`。
3. 儲存檔案。

應用程式會在幾秒內更新。

     <video src="/assets/wails_dev.mp4" controls></video>

### 封裝應用程式
應用程式準備好發布後，您可以建立各平台專用的套件：

@tabs{sync-key="platform"}
[Mac]
若要建立`.app`套件組合：

```bash
wails3 package
```

這會建立正式環境組建，並將其封裝成`.app`套件組合，存放於`bin`目錄中。

[Windows]
若要建立 NSIS 安裝程式：

```bash
wails3 package
```

這會建立正式環境組建，並將其封裝成 NSIS 安裝程式，存放於`bin`目錄中。

[Linux]
Wails 支援多種用於發布 Linux 應用程式的套件格式：

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

如需封裝選項與設定的詳細資訊，請參閱我們的[建置與封裝指南](/guides/build/building/)。

### 設定版本控制與模組名稱
建立專案時會使用預留位置模組名稱`changeme`。建議將它更新為與您的儲存庫 URL 相符：

1. 在 GitHub（或您偏好的 Git 託管服務）上建立新的儲存庫
2. 在專案目錄中初始化 Git：
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. 設定遠端儲存庫（請替換為您的儲存庫 URL）：
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. 更新`go.mod`中的模組名稱，使其與您的儲存庫 URL 相符：
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. 推送您的程式碼：
  ```bash
  git push -u origin main
  ```


這可確保您的 Go 模組名稱符合 Go 的模組命名慣例，也能讓您更容易分享程式碼。

@note{type="tip" title="專業提示"}
建立專案時使用`-git`旗標，即可自動完成所有初始化步驟：

```bash
wails3 init -n myfirstapp -t vanilla-js -git github.com/username/myfirstapp
```

此功能支援多種 Git URL 格式：

- HTTPS：`https://github.com/username/project`
- SSH：`git@github.com:username/project`或`ssh://git@github.com/username/project`
- Git 通訊協定：`git://github.com/username/project`
- 檔案系統：`file:///path/to/project.git`

@end

@end

## 恭喜！

您剛剛已建立、開發並封裝了第一個 Wails 應用程式。這只是您運用 Wails v3實現各種可能的起點。

## 後續步驟

如果您剛開始使用 Wails，建議接著閱讀我們的教學；這些教學將透過實作，引導您瞭解 Wails 的各項功能。第一篇教學是[建立服務](/tutorials/01-creating-a-service/)。

如果您是較進階的使用者，請參閱[建置與封裝指南](/guides/build/building/)，進一步瞭解如何使用 Wails。
