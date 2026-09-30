---
title: "開發環境設定"
description: "設定 Wails v3 開發環境"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## 開發環境設定

本指南將逐步說明如何設定完整的 Wails v3 開發環境。

## 必要工具

### Go 開發環境

1. **安裝 Go 1.25或更新版本：**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **設定 Go 環境：**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **安裝實用的 Go 工具：**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js 與 npm

僅前端整合範例需要。

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### 平台專用相依套件

**macOS：**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows：**

1. 安裝[MSYS2](https://www.msys2.org/)以取得類 Unix 環境
2. WebView2 Runtime（Windows 11已預先安裝；若使用 Windows 10，請[下載](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)）
3. 選用：安裝[Git for Windows](https://git-scm.com/download/win)

**Linux（Debian/Ubuntu）：**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux（Fedora/RHEL）：**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux（Arch）：**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## 儲存庫設定

### 複製並設定

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### 建置 Wails CLI

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### 加入 PATH（選用）

**Linux/macOS：**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows：**

透過「系統內容」將 Wails 目錄加入 PATH 環境變數。

## IDE 設定

### VS Code（建議）

1. **安裝 VS Code：**[下載](https://code.visualstudio.com/)

2. **安裝擴充功能：**
  - Go（由 Google 的 Go 團隊開發）
  - ESLint
  - Prettier
  - MDX（用於文件）


3. **設定工作區組態**（`.vscode/settings.json`）：
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **安裝 GoLand：**[下載](https://www.jetbrains.com/go/)

2. **進行設定：**
  - 啟用 Go modules 支援
  - 為`goimports`設定檔案監看器
  - 設定符合專案慣例的程式碼樣式


## 驗證設定

執行下列命令，確認所有項目皆可正常運作：

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

如果測試應用程式能夠成功建置並執行，表示環境已準備就緒！

## 執行測試

### 單元測試

```bash
cd v3
go test ./...
cd ..
```

### 特定套件測試

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### 執行測試並產生涵蓋率

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 使用競爭偵測器執行測試

```bash
cd v3
go test ./... -race
```

## 處理文件

Wails v3 文件使用 M-Press 編寫。英文原始檔位於  
`docs/mpress/content/`；翻譯檔位於各語言目錄，例如  
`fr/` 和 `id/`。

請從儲存庫根目錄預覽並驗證文件變更：

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

正式環境網站是靜態網站。在本機處理文件時，不需要 Node.js、翻譯服務供應商或 Cloudflare  
憑證。

## 偵錯

### 偵錯 Go 程式碼

**VS Code：**

建立`.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**命令列：**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### 偵錯平台程式碼

偵錯平台專用程式碼需要使用對應的平台工具：

- <strong>macOS：</strong>Xcode Instruments
- <strong>Windows：</strong>Visual Studio Debugger
- <strong>Linux：</strong>GDB

## 常見問題

### 「command not found: wails3」

將 Wails 目錄加入 PATH，或從專案根目錄使用 `./wails3`。

### 「找不到 webkitgtk-6.0」或「找不到 webkit2gtk」（Linux）

針對您要建置的技術堆疊安裝相應的開發套件：

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### 因 Go 模組錯誤導致建置失敗

```bash
cd v3
go mod tidy
go mod download
```

### Windows 上的「CGO_ENABLED」錯誤

請確認 PATH 中包含 C 編譯器（透過 MSYS2 安裝的 MinGW-w64）。

## 後續步驟

- 檢閱[程式碼撰寫標準](/contributing/standards/)
- 瀏覽[技術文件](/contributing/)
- 尋找可著手處理的議題：[適合新手的議題](https://github.com/wailsapp/wails/labels/good%20first%20issue)
