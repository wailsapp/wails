---
title: "Кроссплатформенная сборка"
description: "Сборка для нескольких платформ на одном компьютере"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## Быстрый старт

Wails v3 поддерживает сборку для Windows, macOS и Linux в любой операционной системе хоста. Система сборки автоматически определяет окружение и выбирает подходящий способ компиляции.

**Хотите выполнять кросс-компиляцию для macOS и Linux?** Один раз выполните эту команду, чтобы подготовить образы Docker (будет загружено около 800 МБ):

```bash
wails3 task setup:docker
```

Затем выполняйте сборку для любой платформы:

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

Windows — самая простая целевая платформа для кросс-компиляции, поскольку по умолчанию она не требует CGO.

```bash
wails3 build GOOS=windows
```

Это работает в любой ОС хоста без дополнительной настройки. Встроенные средства кросс-компиляции Go выполняют всю необходимую работу.

**Если вашему приложению требуется CGO** (например, вы используете библиотеку C или пакет, зависящий от CGO), то для сборки в macOS или Linux потребуется Docker:

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

На хостах без Windows файл Taskfile обнаруживает `CGO_ENABLED=1` и автоматически использует образ Docker.

### macOS

Для интеграции с WebView при сборке для macOS требуется CGO, поэтому для кросс-компиляции нужны специальные инструменты.

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**При сборке в Linux или Windows** сначала необходимо настроить Docker:

```bash
wails3 task setup:docker
```

После создания образов система сборки определяет, что вы работаете не в macOS, и автоматически использует Docker. Изменять команды сборки не требуется.

Обратите внимание: собранные кросс-компиляцией двоичные файлы macOS не имеют подписи кода. Перед распространением их необходимо подписать в macOS или в среде CI.

### Linux

Для интеграции с WebView при сборке для Linux требуется CGO.

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**При сборке в macOS или Windows** сначала необходимо настроить Docker:

```bash
wails3 task setup:docker
```

Система сборки определяет, что вы работаете не в Linux, и автоматически использует Docker.

**В Linux без компилятора C** система сборки проверяет наличие `gcc` или `clang`. Если не найден ни один из них, она переключается на Docker. Это удобно для минимальных контейнеров и систем без установленных инструментов сборки. Вы можете:

1. Установить компилятор C: `sudo apt install build-essential` (Debian/Ubuntu) или `sudo pacman -S base-devel` (Arch)
2. Собрать образ Docker, после чего он будет использоваться автоматически

### Архитектура ARM

Все платформы поддерживают кросс-компиляцию для ARM64 с помощью `GOARCH`:

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

Образ Docker содержит целевые конфигурации кросс-компилятора Zig для amd64 и arm64 на всех платформах, поэтому сборки для ARM можно выполнять на любом хосте:

| Сборка ARM64 для | Из Windows | Из macOS | Из Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | Средствами Go | Средствами Go | Средствами Go |
| **macOS ARM64** | Docker | Нативная сборка | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Для сборки Linux ARM64 в Linux x86</em>64 используется Docker, поскольку для кросс-компиляции с CGO требуется другой набор инструментов.

## Принцип работы

### Матрица кросс-компиляции

| Хост → целевая платформа | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | Нативная сборка | Docker | Docker |
| **macOS** | Средствами Go | Нативная сборка | Docker |
| **Linux** | Средствами Go | Docker | Нативная сборка |

- **Нативная сборка** = нативный набор инструментов платформы, дополнительная настройка не требуется
- **Средствами Go** = встроенные средства кросс-компиляции Go (`CGO_ENABLED=0`)
- **Docker** = образ Docker с кросс-компилятором Zig

### Требования CGO

| Целевая платформа | Требуется CGO | Способ кросс-компиляции |
| --- | --- | --- |
| Windows | Нет (по умолчанию) | Средствами Go. Docker — только при `CGO_ENABLED=1` |
| macOS | Да | Docker с macOS SDK |
| Linux | Да | Docker или нативная сборка при наличии компилятора C |

### Автоматическое определение

Taskfile-файлы автоматически выбирают подходящий способ сборки в зависимости от вашей среды:

- **Целевая платформа Windows:** по умолчанию используется нативная кросс-компиляция средствами Go. Если на хосте не под управлением Windows явно задать `CGO_ENABLED=1`, система переключится на Docker.
- **Целевая платформа macOS:** вне macOS Docker используется автоматически. Вмешательство вручную не требуется.
- **Целевая платформа Linux:** выполняется проверка наличия `gcc` или `clang`. Если один из них найден, используется нативная компиляция, иначе система переключается на Docker.

### Образ Docker

Wails использует единый образ Docker (`wails-cross`), позволяющий выполнять сборку для всех платформ. В качестве кросс-компилятора используется [Zig](https://ziglang.org/), который позволяет собирать приложения для любой целевой платформы на любом хосте. Для целевых платформ darwin в образ включён macOS SDK.

```bash
wails3 task setup:docker
```

Чтобы проверить готовность образа, выполните `wails3 doctor`.

### macOS SDK

При сборке образа Docker macOS SDK загружается из репозитория [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks). Это необходимо, поскольку для компиляции с помощью CGO требуются заголовочные файлы macOS.

**Важно:** Wails не распространяет macOS SDK. Перед использованием этой функции пользователи должны самостоятельно ознакомиться с условиями лицензии Apple на SDK.

## Сборка собственного образа

Если вам необходимо настроить образ Docker (например, использовать другую версию macOS SDK, добавить дополнительные инструменты или использовать собственный SDK), вы можете собрать образ самостоятельно.

### Dockerfile

Создайте файл `Dockerfile` со следующим содержимым:

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

### Сборка образа

Сохраните Dockerfile и соберите образ:

```bash
docker build -t wails-cross .
```

### Использование другой версии SDK

Измените аргумент сборки `MACOS_SDK_VERSION`:

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

Доступные варианты перечислены в разделе [«Доступные версии SDK»](https://github.com/wailsapp/macosx-sdks/releases).

### Использование собственного SDK

Если у вас есть собственный macOS SDK (например, извлечённый из Xcode), вы можете изменить Dockerfile так, чтобы вместо загрузки использовался локальный файл:

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

Поместите tar-архив SDK в один каталог с Dockerfile и выполните сборку.

## Интеграция с CI/CD

Для производственных выпусков рекомендуется использовать CI/CD с нативными исполнителями для каждой платформы. Это полностью исключает кросс-компиляцию и обеспечивает получение правильно подписанных двоичных файлов.

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

## Устранение неполадок

### Образ Docker не найден

```
Docker image 'wails-cross' not found.
```

Чтобы собрать образ Docker, выполните `wails3 task setup:docker`. Это потребуется сделать только один раз.

### Демон Docker не запущен

```
Docker is required for cross-compilation. Please install Docker.
```

Запустите Docker Desktop или демон Docker. В Linux может потребоваться выполнить `sudo systemctl start docker`.

### В Linux отсутствует компилятор C

Если при сборке в Linux возникают ошибки, связанные с CGO, есть два варианта:

1. **Установить компилятор C:**
  - Debian/Ubuntu: `sudo apt install build-essential`
  - Arch Linux: `sudo pacman -S base-devel`
  - Fedora: `sudo dnf install gcc`


2. **Использовать Docker:** выполните `wails3 task setup:docker`, и Taskfile автоматически задействует Docker, если компилятор не обнаружен.

### Двоичные файлы macOS не подписаны

Двоичные файлы macOS, полученные путём кросс-компиляции, не имеют подписи кода. Apple требует подписывать код для распространения, поэтому необходимо:

1. Подписать двоичный файл на компьютере с macOS или
2. Подписать его в CI с помощью исполнителя macOS

Подробности см. в разделе [«Подписание приложений»](/guides/build/signing/).

### Создание универсального двоичного файла

Универсальные двоичные файлы (объединяющие arm64 и amd64) можно собирать на любой платформе:

```bash
wails3 task darwin:build:universal
```

В Linux и Windows Wails объединяет двоичные файлы с помощью встроенной команды `wails3 tool lipo` (на базе [konoui/lipo](https://github.com/konoui/lipo)). В результате создаётся единый двоичный файл, который выполняется нативно как на компьютерах Mac с Apple Silicon, так и на компьютерах Mac с процессорами Intel.

## Дальнейшие действия

- [Сборка приложений](/guides/build/building/) — основные команды и параметры сборки
- [Подписание приложений](/guides/build/signing/) — подписание кода для распространения
