---
title: "Compilação multiplataforma"
description: "Compile para várias plataformas em uma única máquina"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## Início rápido

O Wails v3 permite compilar para Windows, macOS e Linux a partir de qualquer sistema operacional host. O sistema de compilação detecta automaticamente o ambiente e escolhe o método de compilação adequado.

**Quer fazer compilação cruzada para macOS e Linux?** Execute isto uma vez para configurar as imagens do Docker (download de aproximadamente 800 MB):

```bash
wails3 task setup:docker
```

Depois, compile para qualquer plataforma:

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

O Windows é o destino de compilação cruzada mais simples, pois, por padrão, não requer CGO.

```bash
wails3 build GOOS=windows
```

Isso funciona em qualquer sistema operacional host sem nenhuma configuração adicional. A compilação cruzada integrada do Go cuida de tudo.

**Se o seu aplicativo requer CGO** (por exemplo, se você usa uma biblioteca C ou um pacote que depende de CGO), será necessário usar o Docker ao compilar no macOS ou Linux:

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

O Taskfile detecta `CGO_ENABLED=1` em hosts que não executam Windows e usa automaticamente a imagem do Docker.

### macOS

As compilações para macOS requerem CGO para a integração com o WebView, o que significa que a compilação cruzada exige ferramentas especiais.

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**No Linux ou Windows**, primeiro será necessário configurar o Docker:

```bash
wails3 task setup:docker
```

Depois que as imagens forem criadas, o sistema de compilação detectará que você não está no macOS e usará o Docker automaticamente. Não é necessário alterar os comandos de compilação.

Observe que os binários do macOS gerados por compilação cruzada não têm assinatura de código. Antes de distribuí-los, será necessário assiná-los no macOS ou na CI.

### Linux

As compilações para Linux requerem CGO para a integração com o WebView.

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**No macOS ou Windows**, primeiro será necessário configurar o Docker:

```bash
wails3 task setup:docker
```

O sistema de compilação detecta que você não está no Linux e usa o Docker automaticamente.

**No Linux sem um compilador C**, o sistema de compilação procura `gcc` ou `clang`. Se nenhum deles for encontrado, ele recorre ao Docker. Isso é útil em contêineres mínimos ou sistemas sem ferramentas de compilação instaladas. Você pode:

1. Instalar um compilador C: `sudo apt install build-essential` (Debian/Ubuntu) ou `sudo pacman -S base-devel` (Arch)
2. Criar a imagem do Docker e permitir que ela seja usada automaticamente

### Arquitetura ARM

Todas as plataformas permitem compilação cruzada para ARM64 usando `GOARCH`:

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

A imagem do Docker inclui alvos do compilador cruzado Zig para amd64 e arm64 em todas as plataformas. Portanto, as compilações para ARM funcionam em qualquer host:

| Compilar para ARM64 com destino a | No Windows | No macOS | No Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | Go nativo | Go nativo | Go nativo |
| **macOS ARM64** | Docker | Nativa | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Linux ARM64 a partir do Linux x86</em>64 usa o Docker porque a compilação cruzada com CGO requer uma cadeia de ferramentas diferente.

## Como funciona

### Matriz de compilação cruzada

| Host → Destino | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | Nativa | Docker | Docker |
| **macOS** | Go nativo | Nativa | Docker |
| **Linux** | Go nativo | Docker | Nativa |

- **Nativa** = cadeia de ferramentas nativa da plataforma, sem configuração adicional
- **Go nativo** = compilação cruzada integrada do Go (`CGO_ENABLED=0`)
- **Docker** = imagem do Docker com o compilador cruzado Zig

### Requisitos de CGO

| Destino | CGO necessário | Método de compilação cruzada |
| --- | --- | --- |
| Windows | Não (por padrão) | Go nativo. Docker somente se `CGO_ENABLED=1` |
| macOS | Sim | Docker com o SDK do macOS |
| Linux | Sim | Docker ou compilação nativa, se houver um compilador C disponível |

### Detecção automática

Os Taskfiles escolhem automaticamente o método de compilação correto com base no seu ambiente:

- **Destino Windows:** Por padrão, usa a compilação cruzada nativa do Go. Se você definir explicitamente `CGO_ENABLED=1` em um host que não seja Windows, passará a usar o Docker.
- **Destino macOS:** Usa o Docker automaticamente quando não está no macOS. Nenhuma intervenção manual é necessária.
- **Destino Linux:** Verifica se há `gcc` ou `clang`. Se encontrar, usa a compilação nativa; caso contrário, recorre ao Docker.

### Imagem do Docker

O Wails usa uma única imagem do Docker (`wails-cross`), capaz de compilar para todas as plataformas. Ela usa o [Zig](https://ziglang.org/) como compilador cruzado, permitindo gerar código para qualquer plataforma a partir de qualquer host. O SDK do macOS está incluído para destinos darwin.

```bash
wails3 task setup:docker
```

Você pode verificar se a imagem está pronta executando `wails3 doctor`.

### SDK do macOS

Durante o processo de criação da imagem, a imagem do Docker baixa o SDK do macOS de [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks). Isso é necessário porque os cabeçalhos do macOS são exigidos para a compilação com CGO.

**Importante:** O Wails não distribui o SDK do macOS. Os usuários são responsáveis por analisar os termos da licença do SDK da Apple antes de usar esse recurso.

## Crie sua própria imagem

Se precisar personalizar a imagem do Docker (por exemplo, usar outra versão do SDK do macOS, adicionar ferramentas ou usar seu próprio SDK), você poderá criar a imagem por conta própria.

### Dockerfile

Crie um `Dockerfile` com o seguinte conteúdo:

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

### Criação da imagem

Salve o Dockerfile e crie a imagem:

```bash
docker build -t wails-cross .
```

### Como usar outra versão do SDK

Altere o argumento de compilação `MACOS_SDK_VERSION`:

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

Consulte as [versões disponíveis do SDK](https://github.com/wailsapp/macosx-sdks/releases) para conhecer as opções.

### Como usar seu próprio SDK

Se você tiver seu próprio SDK do macOS (por exemplo, extraído do Xcode), poderá modificar o Dockerfile para usar um arquivo local em vez de baixá-lo:

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

Coloque o tarball do SDK no mesmo diretório que o Dockerfile e execute a compilação.

## Integração com CI/CD

Para versões de produção, recomendamos usar CI/CD com runners nativos de cada plataforma. Isso elimina totalmente a compilação cruzada e garante a obtenção de binários devidamente assinados.

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

## Solução de problemas

### Imagem do Docker não encontrada

```
Docker image 'wails-cross' not found.
```

Execute `wails3 task setup:docker` para criar a imagem do Docker. Você só precisa fazer isso uma vez.

### Daemon do Docker não está em execução

```
Docker is required for cross-compilation. Please install Docker.
```

Inicie o Docker Desktop ou o daemon do Docker. No Linux, talvez seja necessário executar `sudo systemctl start docker`.

### Nenhum compilador C no Linux

Se você encontrar erros relacionados ao CGO ao compilar no Linux, terá duas opções:

1. **Instalar um compilador C:**
  - Debian/Ubuntu: `sudo apt install build-essential`
  - Arch Linux: `sudo pacman -S base-devel`
  - Fedora: `sudo dnf install gcc`


2. **Usar o Docker:** Execute `wails3 task setup:docker`, e o Taskfile usará o Docker automaticamente quando nenhum compilador for detectado.

### Binários do macOS não assinados

Os binários do macOS gerados por compilação cruzada não têm assinatura de código. A Apple exige a assinatura de código para distribuição, portanto, você precisará:

1. Assinar o binário em uma máquina macOS ou
2. Assinar no CI usando um runner do macOS

Consulte [Assinatura de aplicativos](/guides/build/signing/) para obter detalhes.

### Criação de binários universais

Binários universais (arm64 + amd64 combinados) podem ser criados em qualquer plataforma:

```bash
wails3 task darwin:build:universal
```

No Linux e no Windows, o Wails usa seu comando `wails3 tool lipo` integrado (baseado em [konoui/lipo](https://github.com/konoui/lipo)) para combinar os binários. Isso cria um único binário que é executado nativamente tanto em Macs com Apple Silicon quanto em Macs com processadores Intel.

## Próximas etapas

- [Compilação de aplicativos](/guides/build/building/) — Comandos e opções básicos de compilação
- [Assinatura de aplicativos](/guides/build/signing/) — Assinatura de código para distribuição
