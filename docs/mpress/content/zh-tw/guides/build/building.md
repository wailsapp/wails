---
title: "建置應用程式"
description: "建置並封裝 Wails 應用程式"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3 使用[Task](https://taskfile.dev)作為建置系統。`wails3 build`和`wails3 package`命令是 Task 的便利包裝命令。

## 建置

為目前平台建置：

```bash
wails3 build
```

為特定平台建置：

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

輸出會寫入`bin/`目錄。

@note{type="tip"}
從其他平台交叉編譯至 macOS 或 Linux 需要 Docker。設定方式請參閱[跨平台建置](/guides/build/cross-platform/)。

@end

## 開發

以熱重載執行應用程式：

```bash
wails3 dev
```

這會啟動檔案監看程式，在偵測到變更時重新建置並重新啟動應用程式。前端開發伺服器預設在連接埠9245上執行。

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## 封裝

封裝應用程式以供發佈：

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

這會建立各平台專用的套件：

- **Windows**：NSIS 安裝程式 — 請參閱[Windows 封裝](/guides/build/windows/)
- **macOS**：應用程式套件組合（`.app`）— 請參閱[macOS 封裝](/guides/build/macos/)
- **Linux**：AppImage、deb 和 rpm — 請參閱[Linux 封裝](/guides/build/linux/)

## 自訂建置標籤

使用`-tags`旗標傳遞自訂 Go 建置標籤：

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

標籤會以`EXTRA_TAGS`的形式轉送至底層 Taskfile。詳細資訊請參閱[伺服器建置](/guides/server-build/)和[Linux 封裝 — 舊版 GTK3 支援](/guides/build/linux/#legacy-gtk3-support)。

## 直接使用 Task

如需更多控制，請直接使用 Task：

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

`linux:create:deb`或`darwin:build:universal`等平台專用工作只能透過 Task 使用。

## 產生資產

重新產生圖示或更新建置組態：

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
