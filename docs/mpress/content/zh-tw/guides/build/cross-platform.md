---
title: "跨平台建置"
description: "從單一機器為多個平台建置"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## 快速入門

Wails v3 支援從任何主機作業系統為 Windows、macOS 和 Linux 建置。建置系統會自動偵測您的環境，並選擇正確的編譯方式。

<strong>想要交叉編譯至 macOS 和 Linux 嗎？</strong>執行一次以下操作以設定 Docker 映像檔（下載量約 800MB）：

```bash
wails3 task setup:docker
```

接著即可為任何平台建置：

```bash
# Build for current platform (production by default)
wails3 build

# Build for specific platforms
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# Build for ARM64 architecture
wails3 build GOOS=windows GOARCH=arm64
wails3 build GOOS=darwin GOARCH=arm64
wails3 build GOOS=linux GOARCH=arm64

# Environment variable style also works
GOOS=darwin GOARCH=arm64 wails3 build
```

### Windows

Windows 是最簡單的交叉編譯目標，因為預設不需要 CGO。

```bash
wails3 build GOOS=windows
```

這可在任何主機作業系統上運作，不需要額外設定。Go 內建的交叉編譯功能會處理所有作業。

**如果您的應用程式需要 CGO**（例如使用 C 程式庫或依賴 CGO 的套件），從 macOS 或 Linux 建置時便需要 Docker：

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Taskfile 會在非 Windows 主機上偵測`CGO_ENABLED=1`，並自動使用 Docker 映像檔。

### macOS

macOS 建置需要使用 CGO 來整合 WebView，因此交叉編譯需要特殊工具。

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**從 Linux 或 Windows 建置時**，您需要先設定 Docker：

```bash
wails3 task setup:docker
```

映像檔建置完成後，建置系統會偵測到您並非位於 macOS，並自動使用 Docker。您不需要變更建置命令。

請注意，交叉編譯的 macOS 二進位檔未經程式碼簽署。散布前，您需要在 macOS 或 CI 中簽署這些檔案。

### Linux

Linux 建置需要使用 CGO 來整合 WebView。

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**從 macOS 或 Windows 建置時**，您需要先設定 Docker：

```bash
wails3 task setup:docker
```

建置系統會偵測到您並非位於 Linux，並自動使用 Docker。

**在沒有 C 編譯器的 Linux 上**，建置系統會檢查`gcc`或`clang`。若兩者皆未找到，便會改用 Docker。這對未安裝建置工具的精簡容器或系統很有用。您可以選擇：

1. 安裝 C 編譯器：`sudo apt install build-essential`（Debian/Ubuntu）或`sudo pacman -S base-devel`（Arch）
2. 建置 Docker 映像檔，讓系統自動使用

### ARM 架構

所有平台都支援使用`GOARCH`進行 ARM64 交叉編譯：

```bash
# Windows ARM64 (Surface Pro X, Windows on ARM)
wails3 build GOOS=windows GOARCH=arm64

# Linux ARM64 (Raspberry Pi 4/5, AWS Graviton)
wails3 build GOOS=linux GOARCH=arm64

# macOS ARM64 (Apple Silicon - this is the default on macOS)
wails3 build GOOS=darwin GOARCH=arm64

# macOS Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64
```

Docker 映像檔包含所有平台的 amd64 和 arm64 Zig 交叉編譯器目標，因此可從任何主機進行 ARM 建置：

| 建置 ARM64 目標 | 從 Windows | 從 macOS | 從 Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | 原生 Go | 原生 Go | 原生 Go |
| **macOS ARM64** | Docker | 原生 | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>在 Linux x86</em>64 上建置 Linux ARM64 時會使用 Docker，因為 CGO 交叉編譯需要不同的工具鏈。

## 運作方式

### 交叉編譯矩陣

| 主機 → 目標 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | 原生 | Docker | Docker |
| **macOS** | 原生 Go | 原生 | Docker |
| **Linux** | 原生 Go | Docker | 原生 |

- **原生** = 平台的原生工具鏈，不需要額外設定
- **原生 Go** = Go 內建的交叉編譯（`CGO_ENABLED=0`）
- **Docker** = 包含 Zig 交叉編譯器的 Docker 映像檔

### CGO 需求

| 目標 | 是否需要 CGO | 交叉編譯方式 |
| --- | --- | --- |
| Windows | 否（預設） | 原生 Go。僅在`CGO_ENABLED=1`時使用 Docker |
| macOS | 是 | 使用含 macOS SDK 的 Docker |
| Linux | 是 | Docker；若有可用的 C 編譯器，也可使用原生編譯 |

### 自動偵測

Taskfile 會根據您的環境自動選擇適當的建置方式：

- <strong>Windows 目標：</strong>預設使用 Go 原生交叉編譯。若您在非 Windows 主機上明確設定 `CGO_ENABLED=1`，則會改用 Docker。
- <strong>macOS 目標：</strong>若主機不是 macOS，會自動使用 Docker，無須手動介入。
- <strong>Linux 目標：</strong>檢查是否有 `gcc` 或 `clang`。若找到則使用原生編譯，否則改用 Docker。

### Docker 映像檔

Wails 使用單一 Docker 映像檔（`wails-cross`）建置所有平台的版本。它使用 [Zig](https://ziglang.org/) 作為交叉編譯器，可從任何主機以任何平台為目標進行編譯。映像檔包含供 darwin 目標使用的 macOS SDK。

```bash
wails3 task setup:docker
```

您可以執行 `wails3 doctor`，檢查映像檔是否已準備就緒。

### macOS SDK

Docker 映像檔會在映像檔建置過程中，從 [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks) 下載 macOS SDK。這是必要步驟，因為 CGO 編譯需要 macOS 標頭檔。

<strong>重要：</strong>Wails 不散布 macOS SDK。使用者在使用此功能前，有責任查閱 Apple 的 SDK 授權條款。

## 自行建置映像檔

若需要自訂 Docker 映像檔（例如使用不同版本的 macOS SDK、新增其他工具，或使用您自己的 SDK），可以自行建置映像檔。

### Dockerfile

建立一個 `Dockerfile`，內容如下：

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.24-alpine

ARG ZIG_VERSION=0.14.0
ARG MACOS_SDK_VERSION=14.5
ARG IMAGE_VERSION=1.0.0

LABEL org.opencontainers.image.title="Wails Cross-Compiler"
LABEL org.opencontainers.image.description="Cross-compile Wails v3 apps to macOS, Linux, and Windows"
LABEL org.opencontainers.image.source="https://github.com/wailsapp/wails"
LABEL org.opencontainers.image.vendor="Wails"
LABEL org.opencontainers.image.version="${IMAGE_VERSION}"
LABEL io.wails.sdk.version="${MACOS_SDK_VERSION}"
LABEL io.wails.zig.version="${ZIG_VERSION}"

RUN apk add --no-cache curl xz nodejs npm gcompat

RUN curl -L "https://ziglang.org/download/${ZIG_VERSION}/zig-linux-x86_64-${ZIG_VERSION}.tar.xz" \
    | tar -xJ -C /opt \
    && ln -s /opt/zig-linux-x86_64-${ZIG_VERSION}/zig /usr/local/bin/zig

RUN curl -fL --retry 3 --retry-delay 5 -o /tmp/sdk.tar.xz \
    "https://github.com/wailsapp/macosx-sdks/releases/download/${MACOS_SDK_VERSION}/MacOSX${MACOS_SDK_VERSION}.sdk.tar.xz" \
    && tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX${MACOS_SDK_VERSION}.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz

ENV MACOS_SDK_PATH=/opt/macos-sdk

# Create zig cc wrappers for each target
# Darwin arm64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-darwin-arm64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -mmacosx-version-min=*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target aarch64-macos-none -isysroot /opt/macos-sdk -I/opt/macos-sdk/usr/include -L/opt/macos-sdk/usr/lib -F/opt/macos-sdk/System/Library/Frameworks -w $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-darwin-arm64

# Darwin amd64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-darwin-amd64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -mmacosx-version-min=*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target x86_64-macos-none -isysroot /opt/macos-sdk -I/opt/macos-sdk/usr/include -L/opt/macos-sdk/usr/lib -F/opt/macos-sdk/System/Library/Frameworks -w $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-darwin-amd64

# Linux amd64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-linux-amd64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target x86_64-linux-musl $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-linux-amd64

# Linux arm64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-linux-arm64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target aarch64-linux-musl $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-linux-arm64

# Windows amd64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-windows-amd64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -Wl,*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target x86_64-windows-gnu $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-windows-amd64

# Windows arm64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-windows-arm64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -Wl,*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target aarch64-windows-gnu $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-windows-arm64

# Build script
COPY <<'SCRIPT' /usr/local/bin/build.sh
#!/bin/sh
set -e

OS=${1:-darwin}
ARCH=${2:-arm64}

case "${OS}-${ARCH}" in
    darwin-arm64|darwin-aarch64) export CC=zcc-darwin-arm64; export GOARCH=arm64; export GOOS=darwin ;;
    darwin-amd64|darwin-x86_64)  export CC=zcc-darwin-amd64; export GOARCH=amd64; export GOOS=darwin ;;
    linux-arm64|linux-aarch64)   export CC=zcc-linux-arm64;  export GOARCH=arm64; export GOOS=linux ;;
    linux-amd64|linux-x86_64)    export CC=zcc-linux-amd64;  export GOARCH=amd64; export GOOS=linux ;;
    windows-arm64|windows-aarch64) export CC=zcc-windows-arm64; export GOARCH=arm64; export GOOS=windows ;;
    windows-amd64|windows-x86_64)  export CC=zcc-windows-amd64; export GOARCH=amd64; export GOOS=windows ;;
    *) echo "Usage: <os> <arch>"; echo "  os: darwin, linux, windows"; echo "  arch: amd64, arm64"; exit 1 ;;
esac

export CGO_ENABLED=1
export CGO_CFLAGS="-w"

# Build frontend if exists and not already built (host may have built it)
if [ -d "frontend" ] && [ -f "frontend/package.json" ] && [ ! -d "frontend/dist" ]; then
    (cd frontend && npm install --silent && npm run build --silent)
fi

# Build
APP=${APP_NAME:-$(basename $(pwd))}
mkdir -p bin

EXT=""
LDFLAGS="-s -w"
if [ "$GOOS" = "windows" ]; then
    EXT=".exe"
    LDFLAGS="-s -w -H windowsgui"
fi

go build -ldflags="$LDFLAGS" -o bin/${APP}-${GOOS}-${GOARCH}${EXT} .
echo "Built: bin/${APP}-${GOOS}-${GOARCH}${EXT}"
SCRIPT
RUN chmod +x /usr/local/bin/build.sh

WORKDIR /app
ENTRYPOINT ["/usr/local/bin/build.sh"]
CMD ["darwin", "arm64"]
```

### 建置映像檔

儲存 Dockerfile，然後建置映像檔：

```bash
docker build -t wails-cross .
```

### 使用不同的 SDK 版本

變更 `MACOS_SDK_VERSION` 建置引數：

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

如需可用選項，請參閱[可用的 SDK 版本](https://github.com/wailsapp/macosx-sdks/releases)。

### 使用您自己的 SDK

若您有自己的 macOS SDK（例如從 Xcode 擷取），可以修改 Dockerfile，改用本機檔案而不下載：

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

將 SDK tarball 放在 Dockerfile 所在的目錄中，然後進行建置。

## CI/CD 整合

針對正式環境版本，我們建議使用 CI/CD，並為每個平台採用原生執行器。這可完全避免交叉編譯，並確保產生已正確簽署的二進位檔。

```yaml
name: Build

on:
  push:
    branches: [main]

jobs:
  build:
    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            goos: linux
          - os: macos-latest
            goos: darwin
          - os: windows-latest
            goos: windows

    runs-on: ${{ matrix.os }}

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Wails CLI
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Install Task
        uses: arduino/setup-task@v2

      - name: Build
        run: wails3 build

      - uses: actions/upload-artifact@v4
        with:
          name: app-${{ matrix.goos }}
          path: bin/
```

## 疑難排解

### 找不到 Docker 映像檔

```
Docker image 'wails-cross' not found.
```

執行 `wails3 task setup:docker` 以建置 Docker 映像檔。此操作只需執行一次。

### Docker 常駐程式未執行

```
Docker is required for cross-compilation. Please install Docker.
```

啟動 Docker Desktop 或 Docker 常駐程式。在 Linux 上，您可能需要執行 `sudo systemctl start docker`。

### Linux 上沒有 C 編譯器

若在 Linux 上建置時看到 CGO 相關錯誤，有以下兩種選擇：

1. **安裝 C 編譯器：**
  - Debian/Ubuntu：`sudo apt install build-essential`
  - Arch Linux：`sudo pacman -S base-devel`
  - Fedora：`sudo dnf install gcc`


2. <strong>改用 Docker：</strong>執行 `wails3 task setup:docker`；未偵測到編譯器時，Taskfile 會自動使用該映像檔。

### macOS 二進位檔未簽署

交叉編譯的 macOS 二進位檔未經程式碼簽署。Apple 要求散布時必須進行程式碼簽署，因此您需要：

1. 在 macOS 電腦上簽署二進位檔，或
2. 在 CI 中使用 macOS 執行器進行簽署

如需詳細資訊，請參閱[簽署應用程式](/guides/build/signing/)。

### 建立通用二進位檔

可在任何平台上建置通用二進位檔（合併 arm64 與 amd64）：

```bash
wails3 task darwin:build:universal
```

在 Linux 和 Windows 上，Wails 使用其內建的 `wails3 tool lipo` 命令（由 [konoui/lipo](https://github.com/konoui/lipo) 提供支援）合併二進位檔。這會建立單一二進位檔，可在 Apple Silicon 與 Intel Mac 上原生執行。

## 後續步驟

- [建置應用程式](/guides/build/building/) — 基本建置命令與選項
- [簽署應用程式](/guides/build/signing/) — 用於散布的程式碼簽署
