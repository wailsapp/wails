---
title: "macOS 封裝"
description: "封裝 Wails 應用程式以供 macOS 發佈"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## 私有 macOS API

Wails v3 預設使用公開的 macOS API。若要選用需要未公開 Apple API 的功能，請使用單一 Go 建置標籤`private_mac_apis`建置應用程式：

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

若直接使用 Go 建置，請使用`go build -tags private_mac_apis .`（正式環境則使用`-tags production,private_mac_apis`）。依賴私有行為的現有應用程式必須加入此標籤，才能保留該行為。此標籤僅適用於 macOS 桌面版建置。

如需受影響功能與選項值的完整清單、公開版本建置的確切備援行為、Liquid Glass 樣式對應，以及檢查器建置組合，請參閱[私有 macOS API](/guides/build/private-macos-apis/)。若未使用此標籤，僅限私有 API 的操作不會執行任何動作；公開 Go API 則維持不變。

## 應用程式套件組合

將應用程式封裝為標準 macOS`.app`套件組合：

```bash
wails3 package GOOS=darwin
```

這會建立包含下列項目的`bin/<AppName>.app`：

- 位於`Contents/MacOS/`中的已編譯二進位檔
- 位於`Contents/Resources/`中的應用程式圖示（來自`icons.icns`；若有資產目錄`Assets.car`，則從該目錄取得）
- 包含應用程式中繼資料的`Info.plist`

## 套件組合資源

`Contents/Resources/`是 macOS 應用程式隨附唯讀檔案的標準位置。請將其用於較大的範本、種子資料、媒體、語言套件，或其他應按需開啟、而非透過`embed`編譯至 Go 可執行檔中的內容。

Wails 已將應用程式圖示放在此目錄中。若要加入自己的檔案，請將檔案放入`build/resources/`之類的來源目錄，然後在`build/darwin/Taskfile.yml`的`create:app:bundle`任務中加入複製步驟：

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

若使用 Taskfile 的`darwin:run`任務，請將對等的命令加入其`run`任務，並將目標設為`{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`。

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

對於較大的檔案，請使用`ResourceFS`。它會傳回以`Contents/Resources`為根目錄的`io/fs.FS`，讓呼叫端無須先將整個資源載入 Go 位元組切片，即可開啟並以串流方式處理資源：

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

資源名稱是相對於`Contents/Resources`、以斜線分隔的路徑。除非可執行檔從`.app/Contents/MacOS`執行，否則`ResourceFS`與`LoadResource`會傳回`mac.ErrNotInAppBundle`。

請將套件組合資源視為不可變。變更已簽署應用程式內的檔案會使其程式碼簽章失效；下載、產生或可由使用者編輯的資料應改存於使用者的 Application Support 目錄。

### 通用二進位檔

同時為 Apple Silicon 與 Intel Mac 建置：

```bash
wails3 task darwin:package:universal
```

這會建立可在兩種架構上原生執行的單一`.app`。通用二進位檔可在任何平台上建置；在 Linux 與 Windows 上，系統會自動使用`wails3 tool lipo`。

## 自訂套件組合

編輯`build/darwin/Info.plist`以自訂下列項目：

- 套件組合識別碼（`CFBundleIdentifier`）
- 應用程式名稱與版本
- 最低 macOS 版本
- 檔案關聯
- URL 協定名稱

應用程式圖示由`build/`目錄中的資產產生。請使用`generate:icons`任務：

```bash
wails3 task common:generate:icons
```

此任務使用`build/appicon.png`產生`darwin/icons.icns`與`windows/icon.ico`。在 macOS 上，您也可以提供`build/appicon.icon`（Icon Composer 格式）：此任務會傳遞`-iconcomposerinput appicon.icon -macassetdir darwin`，以從`.icon`檔案產生`Assets.car`與`darwin/icons.icns`（在非 macOS 平台上會略過）。若存在`Assets.car`，請執行`update:build-assets`任務，以便相應更新`Info.plist`與`CFBundleIconName`：

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

首先儲存您的認證資訊。您可以執行互動式精靈（`wails3 setup signing`），或直接呼叫`notarytool`：

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

在`build/darwin/Taskfile.yml`中進行設定：

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

如需詳細資訊，請參閱[簽署應用程式](/guides/build/signing/)。

## DMG 安裝程式

Wails 3隨附的範本提供`wails3 task darwin:package:dmg`。它會先建立`.app`，然後使用 DMG 程式庫建置具樣式的 DMG。DMG 預設使用帶有紅色龍形標誌與 WAILS 文字標誌的 Wails 品牌漸層背景。

```bash
wails3 task darwin:package:dmg
```

較低階的`darwin:create:dmg`任務會從現有的`.app`套件組合建立 DMG，並可直接從 Taskfile 進行設定：

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

### 預設配置

產生的 DMG 包含：

- 左側的應用程式套件組合
- 右側的`Applications`連結
- 大小為540×380像素的 Finder 視窗
- 圖示大小為96點，且每個圖示下方都有標籤
- 來自`build/darwin/dmg-background.png`的 Wails 品牌背景

應用程式與`Applications`圖示的位置會根據所設定的視窗尺寸決定，因此變更`DMG_WINDOW_WIDTH`或`DMG_WINDOW_HEIGHT`時，預設的雙圖示配置仍會維持等比例間距。為獲得最佳效果，請使用像素尺寸與 Finder 視窗相同的背景圖片。

### 取代 DMG 資產

`build/darwin/`下產生的檔案是一般的專案資源，可以替換：

- `DMG_BACKGROUND`控制顯示在 Finder 視窗內容後方的影像。
- `DMG_VOLUME_ICON`控制掛載卷宗所顯示的圖示。
- `DMG_FILE_ICON`控制產生的`.dmg`檔案在 Finder 中顯示的圖示。

卷宗圖示與 DMG 檔案圖示是不同的資源。替換應用程式圖示不會自動替換其中任何一個。

### 新增額外檔案

使用`DMG_FILES`，將安裝程式指令碼、版本資訊、授權條款或其他資源連同應用程式一起加入。其值是以逗號分隔的`name=path`配對清單：

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

`=`之前的名稱是在 DMG 內顯示的檔名。`=`之後的路徑是專案中的來源檔案。開頭與結尾的空白會被忽略。

每個顯示名稱都必須是唯一的。額外檔案不能取代封裝工具已建立的項目，包括應用程式套件或`Applications`項目。名稱衝突會導致封裝失敗並顯示錯誤，而不會產生損壞的 DMG。

@note{type="note"}
由於 DMG 建立程序會使用 macOS 的磁碟映像檔與 Finder 工具，因此僅支援在 macOS 上建立。交叉編譯的`.app`套件可在其他平台上建立，但最終的 DMG 必須在 Mac 上產生。

@end

## 疑難排解

### 「應用程式已損毀，無法開啟」

此應用程式尚未簽署。您可以使用 Developer ID 憑證簽署，也可以讓使用者略過 Gatekeeper：

```bash
xattr -cr /path/to/YourApp.app
```

### 公證失敗

常見問題：

- **認證資訊無效**：重新執行`xcrun notarytool store-credentials`（或`wails3 setup signing`）
- **需要強化執行階段**：如有需要，請確認權利設定包含`com.apple.security.cs.allow-unsigned-executable-memory`
- **缺少時間戳記**：簽署程序應會自動加入時間戳記

### 交叉編譯的應用程式無法執行

交叉編譯的 macOS 二進位檔未經簽署。請將其傳輸至 Mac，並在測試前完成簽署：

```bash
codesign --force --deep --sign - YourApp.app
```
