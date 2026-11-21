---
title: "Plattformübergreifendes Bauen"
description: "Für mehrere Plattformen auf einem einzigen Rechner bauen"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## Schnellstart

Wails v3 unterstützt das Bauen für Windows, macOS und Linux von jedem Host-Betriebssystem aus. Das Build-System erkennt Ihre Umgebung automatisch und wählt die richtige Kompilierungsmethode aus.

**Möchten Sie für macOS und Linux crosskompilieren?** Führen Sie diesen Befehl einmal aus, um die Docker-Images einzurichten (Download von ca. 800 MB):

```bash
wails3 task setup:docker
```

Bauen Sie anschließend für eine beliebige Plattform:

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

Windows ist das einfachste Ziel für die Crosskompilierung, da es standardmäßig kein CGO benötigt.

```bash
wails3 build GOOS=windows
```

Dies funktioniert ohne zusätzliche Einrichtung auf jedem Host-Betriebssystem. Die integrierte Crosskompilierung von Go übernimmt alles.

**Wenn Ihre App CGO benötigt** (z. B. weil Sie eine C-Bibliothek oder ein von CGO abhängiges Paket verwenden), benötigen Sie beim Bauen unter macOS oder Linux Docker:

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Das Taskfile erkennt `CGO_ENABLED=1` auf Nicht-Windows-Hosts und verwendet automatisch das Docker-Image.

### macOS

macOS-Builds benötigen CGO für die WebView-Integration. Daher erfordert die Crosskompilierung spezielle Werkzeuge.

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**Unter Linux oder Windows** müssen Sie zuerst Docker einrichten:

```bash
wails3 task setup:docker
```

Sobald die Images gebaut sind, erkennt das Build-System, dass Sie nicht unter macOS arbeiten, und verwendet automatisch Docker. Sie müssen Ihre Build-Befehle nicht ändern.

Beachten Sie, dass crosskompilierte macOS-Binärdateien nicht codesigniert sind. Vor der Verteilung müssen Sie sie unter macOS oder in der CI signieren.

### Linux

Linux-Builds benötigen CGO für die WebView-Integration.

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**Unter macOS oder Windows** müssen Sie zuerst Docker einrichten:

```bash
wails3 task setup:docker
```

Das Build-System erkennt, dass Sie nicht unter Linux arbeiten, und verwendet automatisch Docker.

**Unter Linux ohne C-Compiler** prüft das Build-System, ob `gcc` oder `clang` vorhanden ist. Wird keines von beiden gefunden, weicht es auf Docker aus. Dies ist für minimale Container oder Systeme ohne installierte Build-Werkzeuge nützlich. Sie können entweder:

1. Einen C-Compiler installieren: `sudo apt install build-essential` (Debian/Ubuntu) oder `sudo pacman -S base-devel` (Arch)
2. Das Docker-Image bauen und automatisch verwenden lassen

### ARM-Architektur

Alle Plattformen unterstützen die ARM64-Crosskompilierung mit `GOARCH`:

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

Das Docker-Image enthält Zig-Crosscompiler-Ziele für amd64 und arm64 auf allen Plattformen. Daher funktionieren ARM-Builds von jedem Host aus:

| ARM64 bauen für | Unter Windows | Unter macOS | Unter Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | Natives Go | Natives Go | Natives Go |
| **macOS ARM64** | Docker | Nativ | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Linux ARM64 unter Linux x86</em>64 verwendet Docker, da die CGO-Crosskompilierung eine andere Toolchain benötigt.

## Funktionsweise

### Crosskompilierungsmatrix

| Host → Ziel | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | Nativ | Docker | Docker |
| **macOS** | Natives Go | Nativ | Docker |
| **Linux** | Natives Go | Docker | Nativ |

- **Nativ** = Native Toolchain der Plattform, keine zusätzliche Einrichtung
- **Natives Go** = Integrierte Crosskompilierung von Go (`CGO_ENABLED=0`)
- **Docker** = Docker-Image mit Zig-Crosscompiler

### CGO-Anforderungen

| Ziel | CGO erforderlich | Crosskompilierungsmethode |
| --- | --- | --- |
| Windows | Nein (standardmäßig) | Natives Go. Docker nur bei `CGO_ENABLED=1` |
| macOS | Ja | Docker mit macOS SDK |
| Linux | Ja | Docker oder nativ, wenn ein C-Compiler verfügbar ist |

### Automatische Erkennung

Die Taskfiles wählen anhand Ihrer Umgebung automatisch die richtige Build-Methode aus:

- **Windows-Ziel:** Verwendet standardmäßig die native Go-Crosskompilierung. Wenn Sie auf einem Nicht-Windows-Host ausdrücklich `CGO_ENABLED=1` festlegen, wechselt das System zu Docker.
- **macOS-Ziel:** Verwendet außerhalb von macOS automatisch Docker. Kein manueller Eingriff erforderlich.
- **Linux-Ziel:** Prüft, ob `gcc` oder `clang` vorhanden ist. Falls ja, wird nativ kompiliert; andernfalls weicht das System auf Docker aus.

### Docker-Image

Wails verwendet ein einziges Docker-Image (`wails-cross`), das für alle Plattformen bauen kann. Es verwendet [Zig](https://ziglang.org/) als Crosscompiler, der von jedem Host aus jede Plattform als Ziel verwenden kann. Für darwin-Ziele ist das macOS SDK enthalten.

```bash
wails3 task setup:docker
```

Mit `wails3 doctor` können Sie prüfen, ob das Image bereit ist.

### macOS SDK

Das Docker-Image lädt das macOS SDK während des Image-Build-Prozesses von [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks) herunter. Dies ist erforderlich, da für die CGO-Kompilierung macOS-Header benötigt werden.

**Wichtig:** Wails vertreibt das macOS SDK nicht. Benutzer sind dafür verantwortlich, vor der Verwendung dieser Funktion die Lizenzbedingungen von Apple für das SDK zu prüfen.

## Eigenes Image erstellen

Wenn Sie das Docker-Image anpassen müssen, etwa um eine andere macOS-SDK-Version, zusätzliche Werkzeuge oder Ihr eigenes SDK zu verwenden, können Sie das Image selbst erstellen.

### Dockerfile

Erstellen Sie eine `Dockerfile` mit folgendem Inhalt:

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

### Image erstellen

Speichern Sie das Dockerfile und erstellen Sie das Image:

```bash
docker build -t wails-cross .
```

### Andere SDK-Version verwenden

Ändern Sie das Build-Argument `MACOS_SDK_VERSION`:

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

Die verfügbaren Optionen finden Sie unter [Verfügbare SDK-Versionen](https://github.com/wailsapp/macosx-sdks/releases).

### Eigenes SDK verwenden

Wenn Sie über ein eigenes macOS-SDK verfügen, etwa eines aus Xcode, können Sie das Dockerfile so ändern, dass statt eines Downloads eine lokale Datei verwendet wird:

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

Legen Sie das SDK-Tarball im selben Verzeichnis wie das Dockerfile ab und erstellen Sie das Image.

## CI/CD-Integration

Für Produktions-Releases empfehlen wir CI/CD mit nativen Runnern für jede Plattform. Dadurch entfällt die Cross-Kompilierung vollständig und Sie erhalten ordnungsgemäß signierte Binärdateien.

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

## Fehlerbehebung

### Docker-Image nicht gefunden

```
Docker image 'wails-cross' not found.
```

Führen Sie `wails3 task setup:docker` aus, um das Docker-Image zu erstellen. Dies ist nur einmal erforderlich.

### Docker-Daemon wird nicht ausgeführt

```
Docker is required for cross-compilation. Please install Docker.
```

Starten Sie Docker Desktop oder den Docker-Daemon. Unter Linux müssen Sie möglicherweise `sudo systemctl start docker` ausführen.

### Kein C-Compiler unter Linux

Wenn beim Erstellen unter Linux CGO-bezogene Fehler auftreten, haben Sie zwei Möglichkeiten:

1. **C-Compiler installieren:**
  - Debian/Ubuntu: `sudo apt install build-essential`
  - Arch Linux: `sudo pacman -S base-devel`
  - Fedora: `sudo dnf install gcc`


2. **Stattdessen Docker verwenden:** Führen Sie `wails3 task setup:docker` aus. Das Taskfile verwendet das Image automatisch, wenn kein Compiler erkannt wird.

### macOS-Binärdateien nicht signiert

Für macOS cross-kompilierte Binärdateien sind nicht codesigniert. Apple verlangt für die Verteilung eine Codesignatur. Daher müssen Sie entweder:

1. die Binärdatei auf einem macOS-Rechner signieren oder
2. sie in der CI mit einem macOS-Runner signieren

Weitere Informationen finden Sie unter [Anwendungen signieren](/guides/build/signing/).

### Universelle Binärdatei erstellen

Universelle Binärdateien, die arm64 und amd64 kombinieren, können auf jeder Plattform erstellt werden:

```bash
wails3 task darwin:build:universal
```

Unter Linux und Windows verwendet Wails den integrierten Befehl `wails3 tool lipo`, der auf [konoui/lipo](https://github.com/konoui/lipo) basiert, um die Binärdateien zusammenzuführen. Dadurch entsteht eine einzelne Binärdatei, die sowohl auf Apple-Silicon- als auch auf Intel-Macs nativ ausgeführt wird.

## Nächste Schritte

- [Anwendungen erstellen](/guides/build/building/) – Grundlegende Build-Befehle und Optionen
- [Anwendungen signieren](/guides/build/signing/) – Codesignatur für die Verteilung
