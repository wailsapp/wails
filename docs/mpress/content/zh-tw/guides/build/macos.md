---
title: "macOS 封裝"
description: "封裝 Wails 應用程式以供 macOS 發佈"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## 私有 macOS API

Wails v3預設使用公開的 macOS API。若要啟用需要未公開 Apple API 的功能，請使用單一 Go 建置標籤`private_mac_apis`建置應用程式：

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

若直接使用 Go 建置，請使用`go build -tags private_mac_apis .`（正式環境則使用`-tags production,private_mac_apis`）。依賴私有行為的現有應用程式必須加入此標籤，才能保留該行為。此標籤僅適用於 macOS 桌面版建置。

如需受影響功能與選項值的完整清單、公開版本建置時的確切備援行為、Liquid Glass 樣式對應，以及檢查器建置組合，請參閱[私有 macOS API](/guides/build/private-macos-apis/)。若未加上此標籤，僅限私有 API 的操作不會執行任何動作；公開 Go API 則維持不變。

## 應用程式套件

將應用程式封裝為標準 macOS `.app`套件：

```bash
wails3 package GOOS=darwin
```

這會建立包含下列內容的`bin/<AppName>.app`：

- 位於`Contents/MacOS/`的已編譯二進位檔
- 位於`Contents/Resources/`的應用程式圖示（來自`icons.icns`；若有資產目錄`Assets.car`，則來自該目錄）
- 包含應用程式中繼資料的`Info.plist`

### 同步資料夾中的專案

部分 macOS 檔案提供者會將 Finder 中繼資料加入`.app`目錄。這會導致`codesign`拒絕該套件，或使簽署後的簽章失效。如果提供者會再次加入中繼資料，僅移除中繼資料並不足以解決問題。

產生的 macOS 建置工作會檢查此行為。偵測到此行為時，Wails 會將專案產生的輸出目錄移至`~/Library/Application Support/Wails/build-output/<project>-<hash>`，並以符號連結取代原始輸出目錄。現有輸出檔案會予以保留，封裝、簽署和開發模式仍會使用`bin/<AppName>.app`路徑。CLI 執行移動時會顯示本機儲存位置。

設定`WAILS_MACOS_OUTPUT_ROOT`可選擇其他未同步的本機儲存位置。重新安置的根目錄必須與專案位於相同的檔案系統。若明確設定的`BIN_DIR`位於專案外部，則會依設定使用該位置。現有專案必須採用更新後的`build/darwin/Taskfile.yml`，才能使用此處理方式。

若要發佈，請複製`bin/<AppName>.app`或建立 DMG。僅複製`bin`符號連結不會包含建置成品。如果直接執行`wails3 tool sign`命令時，其輸入資料夾會重新附加中繼資料，命令將回報錯誤；請先將該套件移至未同步的資料夾，再進行簽署。

## 套件資源

`Contents/Resources/`是隨 macOS 應用程式提供唯讀檔案的標準位置。大型範本、初始資料、媒體、語言套件或其他應於需要時才開啟，而非使用`embed`編譯至 Go 可執行檔的內容，皆可置於此處。

Wails 已將應用程式圖示放在此目錄中。若要加入自己的檔案，請將檔案放入`build/resources/`之類的來源目錄，然後在`build/darwin/Taskfile.yml`的`create:app:bundle`工作中加入複製步驟：

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

如果使用 Taskfile 的`darwin:run`工作，請將對應命令加入其`run`工作，並將目標設為`{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`。

### 從 Go 讀取資源

匯入 macOS 平台套件：

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

對於小型檔案，請使用`LoadResource`：

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

對於大型檔案，請使用`ResourceFS`。它會傳回根目錄為`Contents/Resources`的`io/fs.FS`，讓呼叫端無須先將整個資源載入 Go 位元組切片，即可開啟並以串流方式讀取資源：

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

資源名稱是相對於`Contents/Resources`、以斜線分隔的路徑。除非可執行檔從`.app/Contents/MacOS`執行，否則`ResourceFS`和`LoadResource`會傳回`mac.ErrNotInAppBundle`。

請將套件資源視為不可變。變更已簽署應用程式內的檔案會使其程式碼簽章失效；下載、產生或可由使用者編輯的資料，應改存於使用者的 Application Support 目錄。

### 通用二進位檔

同時為 Apple Silicon 和 Intel Mac 建置：

```bash
wails3 task darwin:package:universal
```

這會建立單一`.app`，可在兩種架構上原生執行。通用二進位檔可在任何平台上建置；在 Linux 和 Windows 上會自動使用`wails3 tool lipo`。

## 自訂套件

編輯`build/darwin/Info.plist`以自訂：

- 套件識別碼（`CFBundleIdentifier`）
- 應用程式名稱與版本
- 最低 macOS 版本
- 檔案關聯
- URL 配置

應用程式圖示由`build/`目錄中的資產產生。請使用`generate:icons`工作：

```bash
wails3 task common:generate:icons
```

此工作使用`build/appicon.png`產生`darwin/icons.icns`和`windows/icon.ico`。在 macOS 上，您也可以提供`build/appicon.icon`（Icon Composer 格式）：此工作會傳遞`-iconcomposerinput appicon.icon -macassetdir darwin`，以便從`.icon`檔案產生`Assets.car`和`darwin/icons.icns`（在非 macOS 平台上會略過）。若存在`Assets.car`，請執行`update:build-assets`工作，讓`Info.plist`和`CFBundleIconName`隨之更新：

```bash
wails3 task common:update:build-assets
```

若要從`build/`目錄手動執行圖示命令：

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## 程式碼簽署

簽署應用程式以供發佈：

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

在`build/darwin/Taskfile.yml`中設定簽署：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### 公證

對於在 Mac App Store 以外發佈的應用程式，Apple 要求進行公證：

```bash
wails3 task darwin:sign:notarize
```

首先，儲存您的憑證。您可以執行互動式精靈（`wails3 signing`），或直接呼叫`notarytool`：

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

在`build/darwin/Taskfile.yml`中設定：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

詳情請參閱[簽署應用程式](/guides/build/signing/)。

## DMG 安裝程式

隨 Wails 提供的3範本包含`wails3 task darwin:package:dmg`。它會先建立`.app`，再使用 DMG 程式庫建置具有樣式的 DMG。DMG 預設使用 Wails 品牌漸層背景，並顯示紅色飛龍標誌與 WAILS 字標。

```bash
wails3 task darwin:package:dmg
```

較底層的`darwin:create:dmg`工作會從現有`.app`套件建立 DMG，且可直接在 Taskfile 中設定：

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### 預設版面配置

產生的 DMG 包含：

- 左側的應用程式套件
- 右側的`Applications`連結
- 大小為540×380像素的 Finder 視窗
- 圖示大小為96點，且每個圖示下方均顯示標籤
- 來自`build/darwin/dmg-background.png`的 Wails 品牌背景

應用程式和`Applications`圖示會依設定的視窗尺寸定位，因此變更`DMG_WINDOW_WIDTH`或`DMG_WINDOW_HEIGHT`時，預設的雙圖示版面配置仍會維持成比例的間距。為獲得最佳效果，請使用與 Finder 視窗像素尺寸相同的背景圖片。

### 取代 DMG 資產

`build/darwin/`下產生的檔案是一般專案資產，可以取代：

- `DMG_BACKGROUND`控制 Finder 視窗內容後方顯示的圖片。
- `DMG_VOLUME_ICON`控制掛載卷宗所顯示的圖示。
- `DMG_FILE_ICON`控制產生的`.dmg`檔案在 Finder 中顯示的圖示。

卷宗圖示與 DMG 檔案圖示是不同的資源。取代應用程式圖示不會自動取代其中任何一個圖示。

### 加入額外檔案

使用 `DMG_FILES`，在應用程式旁加入安裝程式指令碼、版本資訊、授權條款或其他資源。此值是以逗號分隔的 `name=path` 配對清單：

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

`=` 前的名稱是在 DMG 內顯示的檔名。`=` 後的路徑是專案中的來源檔案。開頭與結尾的空白會被忽略。

每個顯示名稱都必須是唯一的。額外檔案不能取代封裝工具已建立的項目，包括應用程式套件或 `Applications` 項目。名稱衝突時，封裝會失敗並顯示錯誤，而不會產生損壞的 DMG。

@note{type="note"}
只有 macOS 支援建立 DMG，因為此程序會使用 macOS 的磁碟映像檔與 Finder 工具。您可以在其他平台建立交叉編譯的 `.app` 套件，但最終的 DMG 必須在 Mac 上產生。

@end

## 疑難排解

### 「App 已損壞，無法打開」

此 App 未簽署。請使用 Developer ID 憑證簽署，使用者也可以略過 Gatekeeper：

```bash
xattr -cr /path/to/YourApp.app
```

### 公證失敗

常見問題：

- **憑證無效**：重新執行 `xcrun notarytool store-credentials`（或 `wails3 signing`）
- **需要強化執行階段**：如有需要，請確認權利設定包含 `com.apple.security.cs.allow-unsigned-executable-memory`
- **缺少時間戳記**：簽署程序應會自動加入時間戳記

### 交叉編譯的 App 無法執行

交叉編譯的 macOS 二進位檔未經簽署。請將其傳輸至 Mac，並在測試前完成簽署：

```bash
codesign --force --deep --sign - YourApp.app
```
