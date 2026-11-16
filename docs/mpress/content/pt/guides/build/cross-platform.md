---
title: "Compilação multiplataforma"
description: "Compile para várias plataformas em uma única máquina"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## Início rápido

O Wails v3 permite compilar para Windows, macOS e Linux a partir de qualquer sistema operacional host. O sistema de compilação detecta automaticamente seu ambiente e escolhe o método de compilação correto.

**Quer fazer compilação cruzada para macOS e Linux?** Execute este comando uma vez para configurar as imagens Docker (download de aproximadamente 800 MB):

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

O Windows é o destino de compilação cruzada mais simples, pois não exige CGO por padrão.

```bash
wails3 build GOOS=windows
```

Isso funciona em qualquer sistema operacional host sem configuração adicional. A compilação cruzada integrada do Go cuida de tudo.

**Se seu aplicativo exigir CGO** (por exemplo, se você usar uma biblioteca C ou um pacote que dependa de CGO), será necessário usar o Docker ao compilar no macOS ou Linux:

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

O Taskfile detecta `CGO_ENABLED=1` em hosts que não executam Windows e usa automaticamente a imagem Docker.

### macOS

As compilações para macOS exigem CGO para a integração com o WebView, o que significa que a compilação cruzada requer ferramentas especiais.

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**No Linux ou Windows**, primeiro você precisará configurar o Docker:

```bash
wails3 task setup:docker
```

Depois que as imagens forem criadas, o sistema de compilação detectará que você não está no macOS e usará o Docker automaticamente. Não é necessário alterar seus comandos de compilação.

Observe que os binários do macOS gerados por compilação cruzada não são assinados digitalmente. Você precisará assiná-los no macOS ou no ambiente de CI antes da distribuição.

### Linux

As compilações para Linux exigem CGO para a integração com o WebView.

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**No macOS ou Windows**, primeiro você precisará configurar o Docker:

```bash
wails3 task setup:docker
```

O sistema de compilação detecta que você não está no Linux e usa o Docker automaticamente.

**No Linux sem um compilador C**, o sistema de compilação procura `gcc` ou `clang`. Se nenhum deles for encontrado, ele recorrerá ao Docker. Isso é útil para contêineres mínimos ou sistemas sem ferramentas de compilação instaladas. Você pode:

1. Instalar um compilador C: `sudo apt install build-essential` (Debian/Ubuntu) ou `sudo pacman -S base-devel` (Arch)
2. Criar a imagem Docker e permitir que ela seja usada automaticamente

### Arquitetura ARM

Todas as plataformas oferecem suporte à compilação cruzada para ARM64 usando `GOARCH`:

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

A imagem Docker inclui destinos do compilador cruzado Zig para amd64 e arm64 em todas as plataformas, portanto as compilações para ARM funcionam em qualquer host:

| Compilar ARM64 para | No Windows | No macOS | No Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | Go nativo | Go nativo | Go nativo |
| **macOS ARM64** | Docker | Nativo | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Linux ARM64 a partir do Linux x86</em>64 usa o Docker porque a compilação cruzada com CGO exige uma cadeia de ferramentas diferente.

## Como funciona

### Matriz de compilação cruzada

| Host → Destino | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | Nativo | Docker | Docker |
| **macOS** | Go nativo | Nativo | Docker |
| **Linux** | Go nativo | Docker | Nativo |

- **Nativo** = cadeia de ferramentas nativa da plataforma, sem configuração adicional
- **Go nativo** = compilação cruzada integrada do Go (`CGO_ENABLED=0`)
- **Docker** = imagem Docker com o compilador cruzado Zig

### Requisitos do CGO

| Destino | CGO obrigatório | Método de compilação cruzada |
| --- | --- | --- |
| Windows | Não (por padrão) | Go nativo. Docker somente se `CGO_ENABLED=1` |
| macOS | Sim | Docker com o SDK do macOS |
| Linux | Sim | Docker ou compilação nativa, se houver um compilador C disponível |

### Detecção automática

Os Taskfiles escolhem automaticamente o método de compilação correto com base no seu ambiente:

- **Destino Windows:** usa a compilação cruzada nativa do Go por padrão. Se você definir explicitamente `CGO_ENABLED=1` em um host que não execute Windows, o sistema passará a usar o Docker.
- **Destino macOS:** usa o Docker automaticamente quando o host não executa macOS. Nenhuma intervenção manual é necessária.
- **Destino Linux:** procura `gcc` ou `clang`. Se encontrar um deles, usa a compilação nativa; caso contrário, recorre ao Docker.

### Imagem Docker

O Wails usa uma única imagem Docker (`wails-cross`), capaz de compilar para todas as plataformas. Ela usa o [Zig](https://ziglang.org/) como compilador cruzado, permitindo direcionar a compilação a qualquer plataforma a partir de qualquer host. O SDK do macOS está incluído para destinos darwin.

```bash
wails3 task setup:docker
```

Você pode verificar se a imagem está pronta executando `wails3 doctor`.

### SDK do macOS

A imagem Docker baixa o SDK do macOS de [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks) durante o processo de criação da imagem. Isso é necessário porque os cabeçalhos do macOS são exigidos para a compilação com CGO.

**Importante:** o Wails não distribui o SDK do macOS. Os usuários são responsáveis por analisar os termos da licença do SDK da Apple antes de usar esse recurso.

## Crie sua própria imagem

Se precisar personalizar a imagem do Docker (por exemplo, usar uma versão diferente do SDK do macOS, adicionar outras ferramentas ou usar seu próprio SDK), você poderá criar a imagem por conta própria.

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

### Como criar a imagem

Salve o Dockerfile e crie a imagem:

```bash
docker build -t wails-cross .
```

### Como usar outra versão do SDK

Altere o argumento de compilação `MACOS_SDK_VERSION`:

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

Consulte as [versões disponíveis do SDK](https://github.com/wailsapp/macosx-sdks/releases) para ver as opções.

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

Coloque o arquivo tar do SDK no mesmo diretório do Dockerfile e crie a imagem.

## Integração com CI/CD

Para versões de produção, recomendamos usar CI/CD com executores nativos para cada plataforma. Isso elimina totalmente a compilação cruzada e garante a geração de binários devidamente assinados.

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

Se ocorrerem erros relacionados ao CGO durante a compilação no Linux, há duas opções:

1. **Instalar um compilador C:**
  - Debian/Ubuntu: `sudo apt install build-essential`
  - Arch Linux: `sudo pacman -S base-devel`
  - Fedora: `sudo dnf install gcc`


2. **Usar o Docker:** execute `wails3 task setup:docker`, e o Taskfile o usará automaticamente quando nenhum compilador for detectado.

### Binários do macOS não assinados

Os binários do macOS gerados por compilação cruzada não têm assinatura de código. A Apple exige assinatura de código para distribuição, portanto você precisará:

1. Assinar o binário em uma máquina com macOS; ou
2. Assiná-lo no CI usando um executor do macOS

Consulte [Assinatura de aplicativos](/guides/build/signing/) para obter detalhes.

### Criação de binários universais

É possível criar binários universais (arm64 + amd64 combinados) em qualquer plataforma:

```bash
wails3 task darwin:build:universal
```

No Linux e no Windows, o Wails usa o comando `wails3 tool lipo` integrado (baseado em [konoui/lipo](https://github.com/konoui/lipo)) para combinar os binários. Isso cria um único binário que é executado nativamente em Macs com Apple Silicon e Intel.

## Próximas etapas

- [Compilação de aplicativos](/guides/build/building/) — Comandos e opções básicos de compilação
- [Assinatura de aplicativos](/guides/build/signing/) — Assinatura de código para distribuição
