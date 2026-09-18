---
title: "クロスプラットフォームビルド"
description: "1台のマシンから複数のプラットフォーム向けにビルド"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## クイックスタート

Wails v3では、どのホストOSからでもWindows、macOS、Linux向けにビルドできます。ビルドシステムが環境を自動的に検出し、適切なコンパイル方法を選択します。

<strong>macOSとLinux向けにクロスコンパイルする場合</strong>は、次のコマンドを一度実行してDockerイメージをセットアップします（ダウンロードサイズは約800MB）。

```bash
wails3 task setup:docker
```

その後、任意のプラットフォーム向けにビルドします。

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

Windowsは、デフォルトではCGOを必要としないため、最も簡単なクロスコンパイルターゲットです。

```bash
wails3 build GOOS=windows
```

追加のセットアップなしで、どのホストOSからでもビルドできます。Goに組み込まれたクロスコンパイル機能がすべてを処理します。

**アプリでCGOが必要な場合**（CライブラリやCGO依存パッケージを使用している場合など）、macOSまたはLinuxからビルドするにはDockerが必要です。

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Windows以外のホストでは、Taskfileが`CGO_ENABLED=1`を検出し、Dockerイメージを自動的に使用します。

### macOS

macOS向けのビルドではWebViewとの統合にCGOが必要なため、クロスコンパイルには専用のツールが必要です。

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

<strong>LinuxまたはWindowsからビルドする場合</strong>は、最初にDockerをセットアップする必要があります。

```bash
wails3 task setup:docker
```

イメージをビルドすると、ビルドシステムがmacOS以外の環境であることを検出し、Dockerを自動的に使用します。ビルドコマンドを変更する必要はありません。

クロスコンパイルしたmacOSバイナリはコード署名されていないことに注意してください。配布前にmacOS上またはCIで署名する必要があります。

### Linux

Linux向けのビルドでは、WebViewとの統合にCGOが必要です。

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

<strong>macOSまたはWindowsからビルドする場合</strong>は、最初にDockerをセットアップする必要があります。

```bash
wails3 task setup:docker
```

ビルドシステムがLinux以外の環境であることを検出し、Dockerを自動的に使用します。

<strong>CコンパイラがないLinux環境</strong>では、ビルドシステムが`gcc`または`clang`を確認します。どちらも見つからない場合は、Dockerを使用します。これは、最小構成のコンテナやビルドツールがインストールされていないシステムで便利です。次のいずれかを選択できます。

1. Cコンパイラをインストールする：`sudo apt install build-essential`（Debian/Ubuntu）または`sudo pacman -S base-devel`（Arch）
2. Dockerイメージをビルドし、自動的に使用されるようにする

### ARMアーキテクチャ

すべてのプラットフォームで、`GOARCH`を使用したARM64クロスコンパイルがサポートされています。

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

Dockerイメージには、すべてのプラットフォーム向けにamd64とarm64の両方のZigクロスコンパイラターゲットが含まれているため、どのホストからでもARM向けにビルドできます。

| ARM64のビルドターゲット | Windowsから | macOSから | Linuxから |
| --- | --- | --- | --- |
| **Windows ARM64** | Goネイティブ | Goネイティブ | Goネイティブ |
| **macOS ARM64** | Docker | ネイティブ | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Linux x86</em>64からLinux ARM64向けにビルドする場合は、CGOのクロスコンパイルに別のツールチェーンが必要なため、Dockerを使用します。

## 仕組み

### クロスコンパイル対応表

| ホスト → ターゲット | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | ネイティブ | Docker | Docker |
| **macOS** | Goネイティブ | ネイティブ | Docker |
| **Linux** | Goネイティブ | Docker | ネイティブ |

- **ネイティブ** = プラットフォームのネイティブツールチェーン。追加のセットアップは不要
- **Goネイティブ** = Goに組み込まれたクロスコンパイル機能（`CGO_ENABLED=0`）
- **Docker** = Zigクロスコンパイラを含むDockerイメージ

### CGOの要件

| ターゲット | CGOの要否 | クロスコンパイル方法 |
| --- | --- | --- |
| Windows | 不要（デフォルト） | Goネイティブ。`CGO_ENABLED=1`の場合のみDocker |
| macOS | 必要 | macOS SDK を使用する Docker |
| Linux | 必要 | Docker、または C コンパイラが利用可能な場合はネイティブビルド |

### 自動検出

Taskfile は、使用環境に応じて適切なビルド方法を自動的に選択します。

- <strong>Windows ターゲット：</strong>デフォルトでは Go のネイティブなクロスコンパイルを使用します。Windows 以外のホストで `CGO_ENABLED=1` を明示的に設定すると、Docker に切り替わります。
- <strong>macOS ターゲット：</strong>macOS 以外では Docker を自動的に使用します。手動での操作は不要です。
- **Linux ターゲット：**`gcc` または `clang` を確認します。見つかった場合はネイティブコンパイルを使用し、見つからない場合は Docker にフォールバックします。

### Docker イメージ

Wails は、すべてのプラットフォーム向けにビルドできる単一の Docker イメージ（`wails-cross`）を使用します。クロスコンパイラには [Zig](https://ziglang.org/) を使用し、どのホストからでも任意のプラットフォームをターゲットにできます。darwin ターゲット用に macOS SDK も含まれています。

```bash
wails3 task setup:docker
```

`wails3 doctor` を実行すると、イメージの準備ができているか確認できます。

### macOS SDK

Docker イメージは、イメージのビルド時に [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks) から macOS SDK をダウンロードします。CGO のコンパイルには macOS のヘッダーファイルが必要なためです。

<strong>重要：</strong>Wails は macOS SDK を配布しません。この機能を使用する前に、Apple の SDK ライセンス条項を確認する責任はユーザーにあります。

## 独自イメージのビルド

Docker イメージをカスタマイズする必要がある場合（別のバージョンの macOS SDK を使用する、追加ツールを導入する、独自の SDK を使用するなど）は、イメージを自分でビルドできます。

### Dockerfile

次の内容で `Dockerfile` を作成します。

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

### イメージのビルド

Dockerfile を保存し、イメージをビルドします。

```bash
docker build -t wails-cross .
```

### 別の SDK バージョンの使用

ビルド引数 `MACOS_SDK_VERSION` を変更します。

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

選択肢については、[利用可能な SDK バージョン](https://github.com/wailsapp/macosx-sdks/releases)を参照してください。

### 独自 SDK の使用

独自の macOS SDK（Xcode から取得したものなど）がある場合は、ダウンロードする代わりにローカルファイルを使用するよう Dockerfile を変更できます。

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

SDK の tarball を Dockerfile と同じディレクトリに配置し、ビルドします。

## CI/CD との統合

本番リリースでは、プラットフォームごとにネイティブランナーを用意した CI/CD の使用を推奨します。これによりクロスコンパイルを完全に回避でき、適切に署名されたバイナリを確実に生成できます。

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

## トラブルシューティング

### Docker イメージが見つからない

```
Docker image 'wails-cross' not found.
```

`wails3 task setup:docker` を実行して Docker イメージをビルドします。この操作が必要なのは一度だけです。

### Docker デーモンが実行されていない

```
Docker is required for cross-compilation. Please install Docker.
```

Docker Desktop または Docker デーモンを起動します。Linux では `sudo systemctl start docker` の実行が必要になる場合があります。

### Linux に C コンパイラがない

Linux でのビルド時に CGO 関連のエラーが表示された場合は、次の 2 つの方法があります。

1. **C コンパイラをインストールする：**
  - Debian/Ubuntu：`sudo apt install build-essential`
  - Arch Linux：`sudo pacman -S base-devel`
  - Fedora：`sudo dnf install gcc`


2. **代わりに Docker を使用する：**`wails3 task setup:docker` を実行すると、コンパイラが検出されない場合に Taskfile が Docker を自動的に使用します。

### macOS バイナリが署名されていない

クロスコンパイルされた macOS バイナリはコード署名されていません。Apple は配布時にコード署名を必須としているため、次のいずれかを行う必要があります。

1. macOS マシン上でバイナリに署名する、または
2. macOS ランナーを使用して CI で署名する

詳細については、[アプリケーションの署名](/guides/build/signing/)を参照してください。

### ユニバーサルバイナリの作成

ユニバーサルバイナリ（arm64 と amd64 を統合したもの）は、どのプラットフォームでもビルドできます。

```bash
wails3 task darwin:build:universal
```

Linux と Windows では、Wails は組み込みの `wails3 tool lipo` コマンド（[konoui/lipo](https://github.com/konoui/lipo) を利用）を使用してバイナリを統合します。これにより、Apple Silicon Mac と Intel Mac の両方でネイティブに動作する単一のバイナリが作成されます。

## 次のステップ

- [アプリケーションのビルド](/guides/build/building/) — 基本的なビルドコマンドとオプション
- [アプリケーションの署名](/guides/build/signing/) — 配布用のコード署名
