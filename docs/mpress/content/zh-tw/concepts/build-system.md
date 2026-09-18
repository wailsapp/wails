---
title: "建置系統"
description: "了解 Wails 如何建置及封裝您的應用程式"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## 統一建置系統

Wails 提供<strong>統一建置系統</strong>，只需一道命令即可編譯 Go 程式碼、封裝前端資產、將所有內容嵌入單一執行檔，並處理平台專屬建置。

```bash
wails3 build
```

<strong>輸出：</strong>嵌入所有內容的原生執行檔。

## 建置流程概覽

**[建置流程圖預留位置]**

## 建置階段

### 1. 分析階段

Wails 會掃描您的 Go 程式碼，以了解其中的服務：

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Wails 擷取的內容：**

- 服務名稱：`GreetService`
- 方法名稱：`Greet`
- 參數型別：`string`
- 傳回型別：`string`

<strong>用途：</strong>產生 TypeScript 繫結

### 2. 產生階段

#### TypeScript 繫結

Wails 會產生型別安全的繫結：

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**優點：**

- 完整的型別安全性
- IDE 自動完成
- 編譯時期錯誤
- JSDoc 註解

#### 前端建置

系統會執行您的前端封裝工具（Vite、webpack 等）：

```bash
# Vite example
vite build --outDir dist
```

**執行內容：**

- 編譯 JavaScript/TypeScript
- 處理並壓縮 CSS
- 最佳化資產
- 產生來源對應表（僅限開發環境）
- 輸出至`frontend/dist/`

### 3. 編譯階段

#### Go 編譯

以最佳化選項編譯 Go 程式碼：

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**旗標：**

- `-s`：移除符號表
- `-w`：移除 DWARF 偵錯資訊
- 結果：較小的二進位檔（縮減約30%）

**平台專屬處理：**

- Windows：嵌入圖示的`.exe`
- macOS：`.app`套件結構
- Linux：ELF 二進位檔

#### 資產嵌入

前端資產會嵌入 Go 二進位檔：

```go
//go:embed frontend/dist
var assets embed.FS
```

<strong>結果：</strong>包含所有內容的單一執行檔。

### 4. 輸出

**單一原生二進位檔：**

- Windows：`myapp.exe`（約 15MB）
- macOS：`myapp.app`（約 15MB）
- Linux：`myapp`（約 15MB）

**無相依項目**（系統 WebView 除外）。

## 開發環境與正式環境

@tabs{sync-key="mode"}
[開發環境（wails3 dev）]
**針對速度最佳化：**

```bash
wails3 dev
```

**執行內容：**

1. 啟動前端開發伺服器（Vite 預設使用連接埠9245）
2. 不進行最佳化即編譯 Go
3. 啟動應用程式並將其指向開發伺服器
4. 啟用熱重載
5. 包含來源對應表

**特性：**

- **快速重新建置**（前端變更耗時&lt;1 秒）
- **不嵌入資產**（由開發伺服器提供）
- 包含<strong>偵錯符號</strong>
- 啟用<strong>來源對應表</strong>
- **詳細記錄**

<strong>檔案大小：</strong>較大（包含偵錯符號時約 50MB）

[正式環境（wails3 build）]
**針對大小和效能最佳化：**

```bash
wails3 build
```

**執行內容：**

1. 為正式環境建置前端（最小化）
2. 啟用最佳化來編譯 Go
3. 移除偵錯符號
4. 嵌入資產
5. 建立單一二進位檔

**特性：**

- **最佳化的程式碼**（已最小化並移除未使用的程式碼）
- **已嵌入資產**（無外部檔案）
- **已移除偵錯符號**
- **無原始碼對應檔**
- **最少量的記錄**

<strong>檔案大小：</strong>較小（約 15MB）

@end

## 建置命令

### 基本建置

```bash
wails3 build
```

**輸出：**`bin/<APP_NAME>`（在 Windows 上則為`bin/<APP_NAME>.exe`）。`bin/`目錄位於專案根目錄。

`wails3 build`是`wails3 task build`的薄包裝。它唯一會轉送的建置階段旗標是`--tags`，該旗標會成為`EXTRA_TAGS` Taskfile 變數：

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build`沒有`-platform`、`-o`、`-skipbindings`、`-clean`、`-debug`、`-devbuild`、`-icon`、`-ldflags`或`-package`旗標。交叉編譯、輸出路徑、圖示與封裝均透過專案的 Taskfile（`Taskfile.yml` + `build/config.yml`）控制。

### 跨平台與平台專屬建置

平台建置會以 Taskfile 任務的形式公開於`darwin:` / `windows:` / `linux:`命名空間下（定義於`build/Taskfile.<platform>.yml`）。例如：

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

若要查看目前專案中的所有可用任務：

```bash
wails3 task --list
```

### 圖示與封裝

從來源 PNG 產生平台圖示（`build/icons.icns`、`build/icon.ico`等）：

```bash
wails3 generate icons -input appicon.png
```

建置平台專屬的安裝程式／套件：

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## 建置設定

### Taskfile.yml

Wails 3專案使用[Taskfile](https://taskfile.dev/)協調建置。根目錄中的`Taskfile.yml`會包含來自`build/`的各平台任務檔案：

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

使用`wails3 task <name>`或`task <name>`執行任務：

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### 專案設定：`build/config.yml`

專案中繼資料（名稱、識別碼、版本、info-plist 值、NSIS 設定、`.desktop`欄位、自訂通訊協定等）位於`build/config.yml`。Taskfile 會在產生圖示、資訊清單、安裝程式等項目時讀取此檔案。Wails 3中<strong>沒有</strong>`build/build.json`檔案。

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

執行`wails3 generate build-assets`（或`wails3 update build-assets`），依據此設定重新產生平台專屬的建置資產。

## 資產嵌入

### 運作方式

Wails 使用 Go 的`embed`套件：

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**建置時：**

1. 將前端建置至`frontend/dist/`
2. `//go:embed`指令會納入檔案
3. 將檔案編譯進二進位檔
4. 二進位檔包含所有內容

**執行階段：**

1. 應用程式啟動
2. 從記憶體提供資產
3. 資產不需要磁碟 I/O
4. 快速載入

### 自訂資產

嵌入其他檔案：

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## 建置最佳化

### 前端最佳化

**Vite（預設）：**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**結果：**

- JavaScript 已最小化（縮減約70%）
- CSS 已最小化（縮減約60%）
- 圖片已最佳化
- 已套用未使用程式碼移除

### Go 最佳化

**編譯器旗標：**

```bash
-ldflags="-s -w"
```

- `-s`：移除符號表（大小減少約10%）
- `-w`：移除 DWARF 偵錯資訊（大小減少約20%）

**其他最佳化：**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`：在建置時設定變數值
- 適用於版本號、建置日期

### 二進位檔壓縮

**UPX（選用）：**

```bash
# After building
upx --best bin/myapp.exe
```

**結果：**

- 大小減少約50%
- 啟動速度稍慢（約 100ms）
- 不建議用於 macOS（程式碼簽署問題）

## 平台特定建置

### Windows

**輸出：**`myapp.exe`

**包含：**

- 應用程式圖示
- 版本資訊
- 資訊清單（UAC 設定）

**圖示：**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

接著，Windows 的`tool package`步驟會將產生的`.ico`嵌入可執行檔。

**資訊清單：**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**輸出：**`myapp.app`（應用程式套件）

**結構：**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist：**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**通用二進位檔：**

macOS Taskfile 提供`darwin:build:universal`（以及`darwin:package:universal`）任務，可建置這兩種架構，並透過`wails3 tool lipo`將其合併：

```bash
wails3 task darwin:build:universal
```

### Linux

**輸出：**`myapp`（ELF 二進位檔）

**相依項目：**

- GTK3
- WebKitGTK

**桌面項目檔：**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**安裝：**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## 建置效能

### 一般建置時間

| 階段 | 時間 | 備註 |
| --- | --- | --- |
| 分析 | &lt;1s | 掃描 Go 程式碼 |
| 繫結產生 | &lt;1s | 產生 TypeScript 程式碼 |
| 前端建置 | 5-30s | 取決於專案大小 |
| Go 編譯 | 2-10s | 取決於程式碼大小 |
| 資產嵌入 | &lt;1s | 嵌入前端資產 |
| **總計** | **10-45s** | 首次建置 |
| **增量建置** | **5-15s** | 後續建置 |

### 加快建置速度

**1. 使用建置快取：**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. 僅執行所需項目：**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. 平行建置（多部機器／CI）：**

在 v3 中，Linux、Windows 與 macOS 之間的交叉編譯通常會在 Docker `wails-cross`容器中進行，或使用各平台的專用執行器；`wails3 build`本身以主機作業系統為目標。請參閱[跨平台建置](/guides/build/cross-platform/)，瞭解支援的工作流程。

**4。使用速度更快的工具：**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## 疑難排解

### 建置失敗

**症狀：**`wails3 build`結束時發生錯誤

**常見原因：**

1. **Go 編譯錯誤**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **前端建置錯誤**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **缺少相依套件**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### 二進位檔過大

<strong>症狀：</strong>二進位檔大於 50 MB

**解決方法：**

1. **移除偵錯符號**（隨附的 Taskfile 已將`-ldflags="-s -w"`傳給`go build`）。

2. **檢查嵌入的資產**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **使用 UPX 壓縮**
  ```bash
  upx --best bin/myapp.exe
  ```


### 建置速度緩慢

<strong>症狀：</strong>建置耗時超過1分鐘

**解決方法：**

1. **使用建置快取**
  - Go 會自動使用快取
  - 前端快取（Vite）會自動啟用


2. **只執行所需的工作**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **最佳化前端建置**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## 最佳實務

### ✅ 建議做法

- **在開發期間使用`wails3 dev`** — 快速迭代
- **發行時使用`wails3 build`** — 產生最佳化輸出
- **為建置加上版本資訊** — 使用`-ldflags`嵌入版本
- **在目標平台上測試建置** — 交叉編譯並不完美
- **維持前端建置速度** — 最佳化打包工具設定
- **使用建置快取** — 加快後續建置速度

### ❌ 請勿

- **請勿提交`build/`目錄** — 將其加入`.gitignore`
- **請勿略過建置測試** — 發行前務必測試
- **請勿嵌入不必要的資產** — 縮小二進位檔
- **請勿在正式環境中使用偵錯建置** — 請使用最佳化建置
- **請勿忘記程式碼簽署** — 發佈時必須進行簽署

## 後續步驟

**建置應用程式** — 建置與封裝的詳細指南 [深入瞭解 →](/guides/build/building/)

**跨平台建置** — 從一台電腦為所有平台進行建置 [深入瞭解 →](/guides/build/cross-platform/)

**建立安裝程式** — 為使用者建立安裝程式 [深入瞭解 →](/guides/installers/)

---

<strong>對建置有疑問嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[建置範例](https://github.com/wailsapp/wails/tree/master/v3/examples/build)。
