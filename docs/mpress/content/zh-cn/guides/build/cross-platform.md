---
title: "跨平台构建"
description: "在一台计算机上为多个平台构建"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## 快速开始

Wails v3 支持从任意主机操作系统为 Windows、macOS 和 Linux 构建。构建系统会自动检测你的环境并选择正确的编译方式。

<strong>想要交叉编译到 macOS 和 Linux？</strong>运行一次以下命令来设置 Docker 镜像（下载约 800MB）：

```bash
wails3 task setup:docker
```

然后为任意平台构建：

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

Windows 是最简单的交叉编译目标，因为默认情况下不需要 CGO。

```bash
wails3 build GOOS=windows
```

无需任何额外设置，即可从任意主机操作系统构建。Go 的内置交叉编译会处理一切。

**如果你的应用需要 CGO**（例如使用了 C 库或依赖 CGO 的软件包），从 macOS 或 Linux 构建时需要使用 Docker：

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Taskfile 会在非 Windows 主机上检测 `CGO_ENABLED=1`，并自动使用 Docker 镜像。

### macOS

macOS 构建需要使用 CGO 来集成 WebView，因此交叉编译需要特殊工具。

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

<strong>从 Linux 或 Windows 构建</strong>时，需要先设置 Docker：

```bash
wails3 task setup:docker
```

镜像构建完成后，构建系统会检测到你使用的不是 macOS，并自动使用 Docker。无需更改构建命令。

请注意，交叉编译生成的 macOS 二进制文件没有代码签名。分发前需要在 macOS 上或 CI 中对其签名。

### Linux

Linux 构建需要使用 CGO 来集成 WebView。

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

<strong>从 macOS 或 Windows 构建</strong>时，需要先设置 Docker：

```bash
wails3 task setup:docker
```

构建系统会检测到你使用的不是 Linux，并自动使用 Docker。

**在没有 C 编译器的 Linux 上**，构建系统会检查是否存在 `gcc` 或 `clang`。如果两者都未找到，则会回退到 Docker。这适用于未安装构建工具的最小化容器或系统。你可以选择：

1. 安装 C 编译器：`sudo apt install build-essential`（Debian/Ubuntu）或 `sudo pacman -S base-devel`（Arch）
2. 构建 Docker 镜像并让系统自动使用它

### ARM 架构

所有平台都支持使用 `GOARCH` 进行 ARM64 交叉编译：

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

Docker 镜像包含适用于所有平台的 amd64 和 arm64 Zig 交叉编译目标，因此可以从任意主机构建 ARM 版本：

| 构建以下平台的 ARM64 版本 | 从 Windows | 从 macOS | 从 Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | Go 原生交叉编译 | Go 原生交叉编译 | Go 原生交叉编译 |
| **macOS ARM64** | Docker | 原生 | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>从 Linux x86</em>64 构建 Linux ARM64 时会使用 Docker，因为 CGO 交叉编译需要不同的工具链。

## 工作原理

### 交叉编译矩阵

| 主机 → 目标平台 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | 原生 | Docker | Docker |
| **macOS** | Go 原生交叉编译 | 原生 | Docker |
| **Linux** | Go 原生交叉编译 | Docker | 原生 |

- **原生** = 平台的原生工具链，无需额外设置
- **Go 原生交叉编译** = Go 的内置交叉编译（`CGO_ENABLED=0`）
- **Docker** = 包含 Zig 交叉编译器的 Docker 镜像

### CGO 要求

| 目标平台 | 是否需要 CGO | 交叉编译方式 |
| --- | --- | --- |
| Windows | 否（默认） | Go 原生交叉编译。仅当 `CGO_ENABLED=1` 时使用 Docker |
| macOS | 是 | 包含 macOS SDK 的 Docker |
| Linux | 是 | Docker；如果有可用的 C 编译器，也可使用原生编译 |

### 自动检测

Taskfile 会根据你的环境自动选择合适的构建方式：

- <strong>Windows 目标平台：</strong>默认使用 Go 原生交叉编译。如果你在非 Windows 主机上显式设置了`CGO_ENABLED=1`，则会切换到 Docker。
- <strong>macOS 目标平台：</strong>不在 macOS 上时自动使用 Docker，无需手动干预。
- <strong>Linux 目标平台：</strong>检查是否存在`gcc`或`clang`。如果找到，则使用原生编译；否则回退到 Docker。

### Docker 镜像

Wails 使用一个可为所有平台构建的 Docker 镜像（`wails-cross`）。该镜像使用[Zig](https://ziglang.org/)作为交叉编译器，可以在任何主机上面向任何平台进行构建。镜像中包含用于 darwin 目标平台的 macOS SDK。

```bash
wails3 task setup:docker
```

你可以运行`wails3 doctor`检查镜像是否已准备就绪。

### macOS SDK

构建 Docker 镜像时，会从[wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks)下载 macOS SDK。CGO 编译需要 macOS 头文件，因此必须执行此操作。

<strong>重要提示：</strong>Wails 不分发 macOS SDK。用户在使用此功能前，应自行查阅 Apple 的 SDK 许可条款。

## 自行构建镜像

如果需要自定义 Docker 镜像（例如使用其他版本的 macOS SDK、添加额外工具或使用你自己的 SDK），可以自行构建镜像。

### Dockerfile

创建一个`Dockerfile`，内容如下：

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

### 构建镜像

保存 Dockerfile，然后构建镜像：

```bash
docker build -t wails-cross .
```

### 使用其他 SDK 版本

更改`MACOS_SDK_VERSION`构建参数：

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

有关可选版本，请参阅[可用的 SDK 版本](https://github.com/wailsapp/macosx-sdks/releases)。

### 使用你自己的 SDK

如果你有自己的 macOS SDK（例如从 Xcode 中提取的 SDK），可以修改 Dockerfile，改用本地文件而不进行下载：

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

将 SDK tarball 放在 Dockerfile 所在的目录中，然后进行构建。

## CI/CD 集成

对于生产版本，我们建议使用 CI/CD，并为每个平台使用原生运行器。这样可以完全避免交叉编译，并确保获得正确签名的二进制文件。

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

## 故障排除

### 找不到 Docker 镜像

```
Docker image 'wails-cross' not found.
```

运行`wails3 task setup:docker`构建 Docker 镜像。此操作只需执行一次。

### Docker 守护进程未运行

```
Docker is required for cross-compilation. Please install Docker.
```

启动 Docker Desktop 或 Docker 守护进程。在 Linux 上，你可能需要运行`sudo systemctl start docker`。

### Linux 上没有 C 编译器

如果在 Linux 上构建时出现与 CGO 相关的错误，你有两种选择：

1. **安装 C 编译器：**
  - Debian/Ubuntu：`sudo apt install build-essential`
  - Arch Linux：`sudo pacman -S base-devel`
  - Fedora：`sudo dnf install gcc`


2. <strong>改用 Docker：</strong>运行`wails3 task setup:docker`；未检测到编译器时，Taskfile 会自动使用该镜像。

### macOS 二进制文件未签名

交叉编译的 macOS 二进制文件没有代码签名。Apple 要求分发的软件必须经过代码签名，因此你需要：

1. 在 macOS 计算机上为二进制文件签名，或
2. 在 CI 中使用 macOS 运行器进行签名

有关详细信息，请参阅[为应用程序签名](/guides/build/signing/)。

### 创建通用二进制文件

可以在任何平台上构建通用二进制文件（合并 arm64 和 amd64）：

```bash
wails3 task darwin:build:universal
```

在 Linux 和 Windows 上，Wails 使用其内置的`wails3 tool lipo`命令（由[konoui/lipo](https://github.com/konoui/lipo)提供支持）合并二进制文件。这样会生成一个可在 Apple Silicon Mac 和 Intel Mac 上原生运行的二进制文件。

## 后续步骤

- [构建应用程序](/guides/build/building/) - 基本构建命令和选项
- [为应用程序签名](/guides/build/signing/) - 用于分发的代码签名
