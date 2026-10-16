---
title: "크로스 플랫폼 빌드"
description: "한 대의 컴퓨터에서 여러 플랫폼용으로 빌드"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## 빠른 시작

Wails v3는 어떤 호스트 운영 체제에서든 Windows, macOS 및 Linux용 빌드를 지원합니다. 빌드 시스템은 환경을 자동으로 감지하여 적절한 컴파일 방식을 선택합니다.

**macOS 및 Linux용으로 크로스 컴파일하려고 하나요?** Docker 이미지를 설정하려면 다음 명령을 한 번 실행하세요(약 800MB 다운로드):

```bash
wails3 task setup:docker
```

그런 다음 원하는 플랫폼용으로 빌드하세요:

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

Windows는 기본적으로 CGO가 필요하지 않으므로 가장 간단한 크로스 컴파일 대상입니다.

```bash
wails3 build GOOS=windows
```

추가 설정 없이 어떤 호스트 OS에서든 빌드할 수 있습니다. Go의 기본 제공 크로스 컴파일 기능이 모든 작업을 처리합니다.

**앱에 CGO가 필요한 경우**(예: C 라이브러리 또는 CGO 종속 패키지를 사용하는 경우) macOS나 Linux에서 빌드하려면 Docker가 필요합니다:

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Taskfile은 Windows가 아닌 호스트에서 `CGO_ENABLED=1`을 감지하고 Docker 이미지를 자동으로 사용합니다.

### macOS

macOS 빌드는 WebView 통합에 CGO가 필요하므로 크로스 컴파일하려면 특수 도구가 필요합니다.

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**Linux 또는 Windows에서 빌드하는 경우** 먼저 Docker를 설정해야 합니다:

```bash
wails3 task setup:docker
```

이미지를 빌드하고 나면 빌드 시스템은 현재 환경이 macOS가 아님을 감지하고 Docker를 자동으로 사용합니다. 빌드 명령을 변경할 필요는 없습니다.

크로스 컴파일된 macOS 바이너리는 코드 서명되지 않습니다. 배포하기 전에 macOS 또는 CI에서 서명해야 합니다.

### Linux

Linux 빌드는 WebView 통합에 CGO가 필요합니다.

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**macOS 또는 Windows에서 빌드하는 경우** 먼저 Docker를 설정해야 합니다:

```bash
wails3 task setup:docker
```

빌드 시스템은 현재 환경이 Linux가 아님을 감지하고 Docker를 자동으로 사용합니다.

**C 컴파일러가 없는 Linux에서는** 빌드 시스템이 `gcc` 또는 `clang`을 확인합니다. 둘 다 없으면 Docker를 대신 사용합니다. 이는 최소 구성 컨테이너나 빌드 도구가 설치되지 않은 시스템에 유용합니다. 다음 중 하나를 선택할 수 있습니다:

1. C 컴파일러를 설치하세요: `sudo apt install build-essential`(Debian/Ubuntu) 또는 `sudo pacman -S base-devel`(Arch)
2. Docker 이미지를 빌드하여 자동으로 사용되도록 하세요

### ARM 아키텍처

모든 플랫폼은 `GOARCH`을 사용한 ARM64 크로스 컴파일을 지원합니다:

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

Docker 이미지에는 모든 플랫폼의 amd64 및 arm64용 Zig 크로스 컴파일러 대상이 포함되어 있으므로 어떤 호스트에서든 ARM 빌드를 수행할 수 있습니다:

| ARM64 빌드 대상 | Windows에서 | macOS에서 | Linux에서 |
| --- | --- | --- | --- |
| **Windows ARM64** | Go 기본 기능 | Go 기본 기능 | Go 기본 기능 |
| **macOS ARM64** | Docker | 네이티브 | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Linux x86</em>64에서 Linux ARM64용으로 빌드할 때는 CGO 크로스 컴파일에 다른 툴체인이 필요하므로 Docker를 사용합니다.

## 작동 방식

### 크로스 컴파일 매트릭스

| 호스트 → 대상 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | 네이티브 | Docker | Docker |
| **macOS** | Go 기본 기능 | 네이티브 | Docker |
| **Linux** | Go 기본 기능 | Docker | 네이티브 |

- **네이티브** = 플랫폼의 네이티브 툴체인, 추가 설정 불필요
- **Go 기본 기능** = Go의 기본 제공 크로스 컴파일(`CGO_ENABLED=0`)
- **Docker** = Zig 크로스 컴파일러가 포함된 Docker 이미지

### CGO 요구 사항

| 대상 | CGO 필요 여부 | 크로스 컴파일 방식 |
| --- | --- | --- |
| Windows | 아니요(기본값) | Go 기본 기능. `CGO_ENABLED=1`인 경우에만 Docker 사용 |
| macOS | 예 | macOS SDK가 포함된 Docker |
| Linux | 예 | Docker 또는 C 컴파일러를 사용할 수 있는 경우 네이티브 빌드 |

### 자동 감지

Taskfile은 환경에 따라 적절한 빌드 방식을 자동으로 선택합니다:

- **Windows 대상:** 기본적으로 네이티브 Go 교차 컴파일을 사용합니다. Windows가 아닌 호스트에서 `CGO_ENABLED=1`을(를) 명시적으로 설정하면 Docker로 전환합니다.
- **macOS 대상:** macOS가 아닌 환경에서는 Docker를 자동으로 사용합니다. 수동으로 개입할 필요가 없습니다.
- **Linux 대상:** `gcc` 또는 `clang`이(가) 있는지 확인합니다. 발견되면 네이티브 컴파일을 사용하고, 그렇지 않으면 Docker를 사용합니다.

### Docker 이미지

Wails는 모든 플랫폼용으로 빌드할 수 있는 단일 Docker 이미지(`wails-cross`)를 사용합니다. 모든 호스트에서 모든 플랫폼을 대상으로 지정할 수 있는 교차 컴파일러인 [Zig](https://ziglang.org/)를 사용합니다. darwin 대상용 macOS SDK도 포함되어 있습니다.

```bash
wails3 task setup:docker
```

`wails3 doctor`을(를) 실행하여 이미지가 준비되었는지 확인할 수 있습니다.

### macOS SDK

Docker 이미지는 이미지 빌드 과정에서 [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks)에서 macOS SDK를 다운로드합니다. CGO 컴파일에는 macOS 헤더가 필요하므로 이 SDK가 필요합니다.

**중요:** Wails는 macOS SDK를 배포하지 않습니다. 이 기능을 사용하기 전에 Apple의 SDK 라이선스 조건을 검토할 책임은 사용자에게 있습니다.

## 직접 이미지 빌드하기

Docker 이미지를 사용자 지정해야 하는 경우(예: 다른 macOS SDK 버전 사용, 추가 도구 설치 또는 자체 SDK 사용) 이미지를 직접 빌드할 수 있습니다.

### Dockerfile

다음 내용으로 `Dockerfile`을(를) 만드세요:

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

### 이미지 빌드하기

Dockerfile을 저장하고 이미지를 빌드하세요:

```bash
docker build -t wails-cross .
```

### 다른 SDK 버전 사용하기

`MACOS_SDK_VERSION` 빌드 인수를 변경하세요:

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

선택할 수 있는 버전은 [사용 가능한 SDK 버전](https://github.com/wailsapp/macosx-sdks/releases)에서 확인하세요.

### 자체 SDK 사용하기

자체 macOS SDK(예: Xcode에서 추출한 SDK)가 있는 경우 다운로드 대신 로컬 파일을 사용하도록 Dockerfile을 수정할 수 있습니다:

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

SDK tarball을 Dockerfile과 같은 디렉터리에 두고 빌드하세요.

## CI/CD 통합

프로덕션 릴리스에는 플랫폼별 네이티브 러너를 사용하는 CI/CD를 권장합니다. 이렇게 하면 교차 컴파일을 완전히 피할 수 있으며 올바르게 서명된 바이너리를 얻을 수 있습니다.

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

## 문제 해결

### Docker 이미지를 찾을 수 없음

```
Docker image 'wails-cross' not found.
```

`wails3 task setup:docker`을(를) 실행하여 Docker 이미지를 빌드하세요. 이 작업은 한 번만 하면 됩니다.

### Docker 데몬이 실행되고 있지 않음

```
Docker is required for cross-compilation. Please install Docker.
```

Docker Desktop 또는 Docker 데몬을 시작하세요. Linux에서는 `sudo systemctl start docker`을(를) 실행해야 할 수 있습니다.

### Linux에 C 컴파일러가 없음

Linux에서 빌드할 때 CGO 관련 오류가 발생하면 다음 두 가지 방법 중 하나를 선택할 수 있습니다:

1. **C 컴파일러 설치:**
  - Debian/Ubuntu: `sudo apt install build-essential`
  - Arch Linux: `sudo pacman -S base-devel`
  - Fedora: `sudo dnf install gcc`


2. **Docker 사용:** `wails3 task setup:docker`을(를) 실행하세요. 컴파일러가 감지되지 않으면 Taskfile이 Docker를 자동으로 사용합니다.

### macOS 바이너리가 서명되지 않음

교차 컴파일된 macOS 바이너리에는 코드 서명이 적용되지 않습니다. Apple은 배포 시 코드 서명을 요구하므로 다음 중 하나를 수행해야 합니다:

1. macOS 머신에서 바이너리에 서명하거나
2. macOS 러너를 사용하는 CI에서 서명하세요.

자세한 내용은 [애플리케이션 서명](/guides/build/signing/)을 참조하세요.

### 유니버설 바이너리 만들기

유니버설 바이너리(arm64와 amd64 결합)는 모든 플랫폼에서 빌드할 수 있습니다:

```bash
wails3 task darwin:build:universal
```

Linux와 Windows에서 Wails는 내장 `wails3 tool lipo` 명령([konoui/lipo](https://github.com/konoui/lipo) 기반)을 사용하여 바이너리를 결합합니다. 이렇게 생성된 단일 바이너리는 Apple Silicon Mac과 Intel Mac 모두에서 네이티브로 실행됩니다.

## 다음 단계

- [애플리케이션 빌드](/guides/build/building/) - 기본 빌드 명령 및 옵션
- [애플리케이션 서명](/guides/build/signing/) - 배포를 위한 코드 서명
