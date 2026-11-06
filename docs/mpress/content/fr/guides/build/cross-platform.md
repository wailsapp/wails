---
title: "Compilation multiplateforme"
description: "Compilez pour plusieurs plateformes depuis une seule machine"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## Démarrage rapide

Wails v3 permet de compiler pour Windows, macOS et Linux depuis n’importe quel système d’exploitation hôte. Le système de compilation détecte automatiquement votre environnement et choisit la méthode de compilation appropriée.

**Vous souhaitez effectuer une compilation croisée pour macOS et Linux ?** Exécutez cette commande une fois pour configurer les images Docker (téléchargement d’environ 800 Mo) :

```bash
wails3 task setup:docker
```

Compilez ensuite pour n’importe quelle plateforme :

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

Windows est la cible de compilation croisée la plus simple, car elle ne nécessite pas CGO par défaut.

```bash
wails3 build GOOS=windows
```

Cela fonctionne depuis n’importe quel système d’exploitation hôte, sans configuration supplémentaire. La compilation croisée intégrée à Go s’occupe de tout.

**Si votre application nécessite CGO** (par exemple, si vous utilisez une bibliothèque C ou un paquet dépendant de CGO), vous aurez besoin de Docker pour compiler depuis macOS ou Linux :

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Le Taskfile détecte `CGO_ENABLED=1` sur les hôtes autres que Windows et utilise automatiquement l’image Docker.

### macOS

Les compilations macOS nécessitent CGO pour l’intégration de WebView, ce qui impose un outillage spécial pour la compilation croisée.

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**Depuis Linux ou Windows**, vous devez d’abord configurer Docker :

```bash
wails3 task setup:docker
```

Une fois les images créées, le système de compilation détecte que vous n’êtes pas sous macOS et utilise automatiquement Docker. Vous n’avez pas besoin de modifier vos commandes de compilation.

Notez que les binaires macOS issus d’une compilation croisée ne sont pas signés. Vous devrez les signer sous macOS ou dans votre environnement d’intégration continue avant leur distribution.

### Linux

Les compilations Linux nécessitent CGO pour l’intégration de WebView.

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**Depuis macOS ou Windows**, vous devez d’abord configurer Docker :

```bash
wails3 task setup:docker
```

Le système de compilation détecte que vous n’êtes pas sous Linux et utilise automatiquement Docker.

**Sous Linux sans compilateur C**, le système de compilation recherche `gcc` ou `clang`. S’il ne trouve ni l’un ni l’autre, il se rabat sur Docker. Cette solution est utile pour les conteneurs minimaux ou les systèmes sur lesquels aucun outil de compilation n’est installé. Vous pouvez :

1. Installer un compilateur C : `sudo apt install build-essential` (Debian/Ubuntu) ou `sudo pacman -S base-devel` (Arch)
2. Créer l’image Docker et laisser le système l’utiliser automatiquement

### Architecture ARM

Toutes les plateformes prennent en charge la compilation croisée pour ARM64 à l’aide de `GOARCH` :

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

L’image Docker inclut les cibles du compilateur croisé Zig pour amd64 et arm64 sur toutes les plateformes. Les compilations ARM fonctionnent donc depuis n’importe quel hôte :

| Compiler pour ARM64 sur | Depuis Windows | Depuis macOS | Depuis Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | Go natif | Go natif | Go natif |
| **macOS ARM64** | Docker | Natif | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Linux ARM64 depuis Linux x86</em>64 utilise Docker, car la compilation croisée avec CGO nécessite une chaîne d’outils différente.

## Fonctionnement

### Matrice de compilation croisée

| Hôte → cible | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | Natif | Docker | Docker |
| **macOS** | Go natif | Natif | Docker |
| **Linux** | Go natif | Docker | Natif |

- **Natif** = chaîne d’outils native de la plateforme, sans configuration supplémentaire
- **Go natif** = compilation croisée intégrée à Go (`CGO_ENABLED=0`)
- **Docker** = image Docker avec le compilateur croisé Zig

### Exigences de CGO

| Cible | CGO requis | Méthode de compilation croisée |
| --- | --- | --- |
| Windows | Non (par défaut) | Go natif. Docker uniquement si `CGO_ENABLED=1` |
| macOS | Oui | Docker avec le SDK macOS |
| Linux | Oui | Docker, ou compilation native si un compilateur C est disponible |

### Détection automatique

Les Taskfiles choisissent automatiquement la méthode de compilation appropriée en fonction de votre environnement :

- **Cible Windows :** utilise par défaut la compilation croisée native de Go. Si vous définissez explicitement `CGO_ENABLED=1` sur un hôte autre que Windows, le système passe à Docker.
- **Cible macOS :** utilise automatiquement Docker lorsque l’hôte n’exécute pas macOS. Aucune intervention manuelle n’est nécessaire.
- **Cible Linux :** recherche `gcc` ou `clang`. Si l’un des deux est trouvé, le système utilise la compilation native ; sinon, il se rabat sur Docker.

### Image Docker

Wails utilise une seule image Docker (`wails-cross`), capable de compiler pour toutes les plateformes. Elle utilise [Zig](https://ziglang.org/) comme compilateur croisé, ce qui permet de cibler n’importe quelle plateforme depuis n’importe quel hôte. Le SDK macOS est inclus pour les cibles darwin.

```bash
wails3 task setup:docker
```

Vous pouvez vérifier que l’image est prête en exécutant `wails3 doctor`.

### SDK macOS

Pendant la création de l’image, l’image Docker télécharge le SDK macOS depuis [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks). Cette opération est nécessaire, car les en-têtes macOS sont requis pour la compilation avec CGO.

**Important :** Wails ne distribue pas le SDK macOS. Il incombe aux utilisateurs de consulter les conditions de licence du SDK d’Apple avant d’utiliser cette fonctionnalité.

## Créer votre propre image

Si vous devez personnaliser l’image Docker, par exemple pour utiliser une autre version du SDK macOS, ajouter des outils ou utiliser votre propre SDK, vous pouvez créer vous-même l’image.

### Dockerfile

Créez un fichier `Dockerfile` avec le contenu suivant :

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

### Créer l’image

Enregistrez le Dockerfile, puis créez l’image :

```bash
docker build -t wails-cross .
```

### Utiliser une autre version du SDK

Modifiez l’argument de compilation `MACOS_SDK_VERSION` :

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

Consultez les [versions disponibles du SDK](https://github.com/wailsapp/macosx-sdks/releases) pour connaître les options.

### Utiliser votre propre SDK

Si vous disposez de votre propre SDK macOS, par exemple extrait de Xcode, vous pouvez modifier le Dockerfile pour utiliser un fichier local au lieu de le télécharger :

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

Placez l’archive tar de votre SDK dans le même répertoire que le Dockerfile, puis lancez la création de l’image.

## Intégration CI/CD

Pour les versions de production, nous vous recommandons d’utiliser une chaîne CI/CD avec des exécuteurs natifs pour chaque plateforme. Vous évitez ainsi entièrement la compilation croisée et obtenez des binaires correctement signés.

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

## Résolution des problèmes

### Image Docker introuvable

```
Docker image 'wails-cross' not found.
```

Exécutez `wails3 task setup:docker` pour créer l’image Docker. Vous ne devez le faire qu’une seule fois.

### Démon Docker non démarré

```
Docker is required for cross-compilation. Please install Docker.
```

Démarrez Docker Desktop ou le démon Docker. Sous Linux, vous devrez peut-être exécuter `sudo systemctl start docker`.

### Aucun compilateur C sous Linux

Si des erreurs liées à CGO apparaissent lors de la compilation sous Linux, deux possibilités s’offrent à vous :

1. **Installez un compilateur C :**
  - Debian/Ubuntu : `sudo apt install build-essential`
  - Arch Linux : `sudo pacman -S base-devel`
  - Fedora : `sudo dnf install gcc`


2. **Utilisez plutôt Docker :** exécutez `wails3 task setup:docker`. Le Taskfile l’utilisera automatiquement si aucun compilateur n’est détecté.

### Binaires macOS non signés

Les binaires macOS issus d’une compilation croisée ne sont pas signés. Apple exige la signature du code pour la distribution. Vous devrez donc :

1. Signer le binaire sur une machine macOS, ou
2. Le signer dans la chaîne CI à l’aide d’un exécuteur macOS

Consultez [Signer des applications](/guides/build/signing/) pour plus de détails.

### Créer un binaire universel

Vous pouvez créer des binaires universels combinant arm64 et amd64 sur n’importe quelle plateforme :

```bash
wails3 task darwin:build:universal
```

Sous Linux et Windows, Wails utilise sa commande `wails3 tool lipo` intégrée, basée sur [konoui/lipo](https://github.com/konoui/lipo), pour combiner les binaires. Vous obtenez ainsi un seul binaire qui s’exécute nativement sur les Mac équipés d’Apple Silicon comme sur ceux dotés d’un processeur Intel.

## Étapes suivantes

- [Créer des applications](/guides/build/building/) - Commandes et options de compilation de base
- [Signer des applications](/guides/build/signing/) - Signature du code pour la distribution
