---
title: "跨平台构建"
description: "在一台机器上为多个平台构建"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## 快速开始

Wails v3 支持从任意宿主操作系统为 Windows、macOS 和 Linux 构建。构建系统会自动检测你的环境并选择正确的编译方式。

<strong>想要交叉编译到 macOS 和 Linux？</strong>运行一次以下命令来设置 Docker 镜像（下载约 800MB）：

```bash
wails3 task setup:docker
```

然后即可为任意平台构建：

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

Windows 是最简单的交叉编译目标，因为默认不需要 CGO。

```bash
wails3 build GOOS=windows
```

这可以在任意宿主操作系统上运行，无需额外设置。Go 内置的交叉编译功能会处理一切。

**如果你的应用需要 CGO**（例如使用了 C 库或依赖 CGO 的包），从 macOS 或 Linux 构建时需要使用 Docker：

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Taskfile 会在非 Windows 宿主机上检测 `CGO_ENABLED=1`，并自动使用 Docker 镜像。

### macOS

macOS 构建需要使用 CGO 集成 WebView，因此交叉编译需要特殊工具。

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**从 Linux 或 Windows 构建时**，需要先设置 Docker：

```bash
wails3 task setup:docker
```

镜像构建完成后，构建系统会检测到你当前未使用 macOS，并自动使用 Docker。无需更改构建命令。

请注意，交叉编译得到的 macOS 二进制文件没有代码签名。分发前需要在 macOS 或 CI 中对其签名。

### Linux

Linux 构建需要使用 CGO 集成 WebView。

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**从 macOS 或 Windows 构建时**，需要先设置 Docker：

```bash
wails3 task setup:docker
```

构建系统会检测到你当前未使用 Linux，并自动使用 Docker。

**在没有 C 编译器的 Linux 上**，构建系统会检查 `gcc` 或 `clang`。如果两者均未找到，则回退到 Docker。这对最小化容器或未安装构建工具的系统很有用。你可以选择：

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

Docker 镜像包含适用于所有平台上 amd64 和 arm64 的 Zig 交叉编译目标，因此可以从任意宿主机进行 ARM 构建：

| 构建 ARM64 目标 | 从 Windows | 从 macOS | 从 Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | 原生 Go | 原生 Go | 原生 Go |
| **macOS ARM64** | Docker | 原生 | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>从 Linux x86 构建 Linux ARM64</em>64时会使用 Docker，因为 CGO 交叉编译需要不同的工具链。

## 工作原理

### 交叉编译矩阵

| 宿主机 → 目标平台 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | 原生 | Docker | Docker |
| **macOS** | 原生 Go | 原生 | Docker |
| **Linux** | 原生 Go | Docker | 原生 |

- **原生** = 使用平台的原生工具链，无需额外设置
- **原生 Go** = 使用 Go 内置的交叉编译功能（`CGO_ENABLED=0`）
- **Docker** = 包含 Zig 交叉编译器的 Docker 镜像

### CGO 要求

| 目标平台 | 是否需要 CGO | 交叉编译方式 |
| --- | --- | --- |
| Windows | 否（默认） | 原生 Go。仅在设置了 `CGO_ENABLED=1` 时使用 Docker |
| macOS | 是 | 使用包含 macOS SDK 的 Docker |
| Linux | 是 | 使用 Docker；如果有可用的 C 编译器，则使用原生编译 |

### 自动检测

Taskfile 会根据你的环境自动选择正确的构建方式：

- <strong>Windows 目标平台：</strong>默认使用原生 Go 交叉编译。如果你在非 Windows 宿主机上显式设置 `CGO_ENABLED=1`，则切换为 Docker。
- <strong>macOS 目标平台：</strong>不在 macOS 上时自动使用 Docker。无需手动干预。
- <strong>Linux 目标平台：</strong>检查 `gcc` 或 `clang`。如果找到，则使用原生编译；否则回退到 Docker。

### Docker 镜像

Wails 使用一个 Docker 镜像（`wails-cross`）即可为所有平台构建。该镜像使用 [Zig](https://ziglang.org/) 作为交叉编译器，可以从任意宿主机面向任意平台进行编译。镜像中包含用于 darwin 目标的 macOS SDK。

```bash
wails3 task setup:docker
```

可以运行 `wails3 doctor` 检查镜像是否已准备就绪。

### macOS SDK

构建 Docker 镜像时，该镜像会从 [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks) 下载 macOS SDK。之所以需要这样做，是因为 CGO 编译需要 macOS 头文件。

<strong>重要：</strong>Wails 不分发 macOS SDK。用户在使用此功能前有责任审阅 Apple 的 SDK 许可条款。

## 构建自己的镜像

如果需要自定义 Docker 镜像（例如使用其他版本的 macOS SDK、添加额外工具或使用自己的 SDK），可以自行构建镜像。

### Dockerfile

创建一个包含以下内容的`Dockerfile`：

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

可用选项请参阅[可用的 SDK 版本](https://github.com/wailsapp/macosx-sdks/releases)。

### 使用自己的 SDK

如果你有自己的 macOS SDK（例如从 Xcode 中提取的 SDK），可以修改 Dockerfile，使其使用本地文件而不是下载文件：

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

将 SDK tarball 放在与 Dockerfile 相同的目录中，然后进行构建。

## CI/CD 集成

对于生产版本，我们建议使用 CI/CD，并为每个平台使用原生运行器。这样可以完全避免交叉编译，并确保生成的二进制文件已正确签名。

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

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

运行`wails3 task setup:docker`以构建 Docker 镜像。此操作只需执行一次。

### Docker 守护进程未运行

```
Docker is required for cross-compilation. Please install Docker.
```

启动 Docker Desktop 或 Docker 守护进程。在 Linux 上，可能需要运行`sudo systemctl start docker`。

### Linux 上没有 C 编译器

如果在 Linux 上构建时出现与 CGO 相关的错误，有以下两种选择：

1. **安装 C 编译器：**
  - Debian/Ubuntu：`sudo apt install build-essential`
  - Arch Linux：`sudo pacman -S base-devel`
  - Fedora：`sudo dnf install gcc`


2. <strong>改用 Docker：</strong>运行`wails3 task setup:docker`；如果未检测到编译器，Taskfile 将自动使用它。

### macOS 二进制文件未签名

交叉编译生成的 macOS 二进制文件没有代码签名。Apple 要求分发的软件必须进行代码签名，因此需要：

1. 在 macOS 计算机上对二进制文件进行签名，或者
2. 在 CI 中使用 macOS 运行器进行签名

有关详细信息，请参阅[应用签名](/guides/build/signing/)。

### 创建通用二进制文件

可以在任何平台上构建通用二进制文件（合并 arm64 和 amd64）：

```bash
wails3 task darwin:build:universal
```

在 Linux 和 Windows 上，Wails 使用其内置的`wails3 tool lipo`命令（由[konoui/lipo](https://github.com/konoui/lipo)提供支持）合并二进制文件。这样会创建一个可在 Apple 芯片和 Intel Mac 上原生运行的二进制文件。

## 后续步骤

- [构建应用程序](/guides/build/building/) - 基本构建命令和选项
- [应用签名](/guides/build/signing/) - 用于分发的代码签名
