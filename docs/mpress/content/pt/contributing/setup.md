---
title: "Configuração do ambiente de desenvolvimento"
description: "Configure seu ambiente de desenvolvimento para trabalhar no Wails v3"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## Configuração do ambiente de desenvolvimento

Este guia orienta você na configuração de um ambiente de desenvolvimento completo para trabalhar no Wails v3.

## Ferramentas necessárias

### Desenvolvimento em Go

1. **Instale o Go 1.26 ou posterior:**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **Configure o ambiente Go:**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **Instale ferramentas úteis do Go:**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js e npm

Necessário apenas para exemplos de integração com o frontend.

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### Dependências específicas da plataforma

**macOS:**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows:**

1. Instale o [MSYS2](https://www.msys2.org/) para ter um ambiente semelhante ao Unix
2. WebView2 Runtime (pré-instalado no Windows 11; para o Windows 10, faça o [download](https://developer.microsoft.com/en-us/microsoft-edge/webview2/))
3. Opcional: instale o [Git for Windows](https://git-scm.com/download/win)

**Linux (Debian/Ubuntu):**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux (Fedora/RHEL):**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux (Arch):**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## Configuração do repositório

### Clonar e configurar

```bash
# Clone your fork
git clone --branch releases/v3.0.0 https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Compilar a CLI do Wails

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### Adicionar ao PATH (opcional)

**Linux/macOS:**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows:**

Adicione o diretório do Wails à variável de ambiente PATH por meio das Propriedades do Sistema.

## Configuração do IDE

### VS Code (recomendado)

1. **Instale o VS Code:** [Baixar](https://code.visualstudio.com/)

2. **Instale as extensões:**
  - Go (da equipe do Go no Google)
  - ESLint
  - Prettier
  - MDX (para a documentação)


3. **Configure as definições do espaço de trabalho** (`.vscode/settings.json`):
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **Instale o GoLand:** [Baixar](https://www.jetbrains.com/go/)

2. **Configure:**
  - Ative o suporte a módulos Go
  - Configure observadores de arquivos para `goimports`
  - Configure o estilo do código de acordo com as convenções do projeto


## Verificar a configuração

Execute estes comandos para verificar se tudo está funcionando:

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

Se o aplicativo de teste for compilado e executado, seu ambiente estará pronto!

## Executar testes

### Testes unitários

```bash
cd v3
go test ./...
cd ..
```

### Testes de pacotes específicos

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### Executar com cobertura

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Executar com o detector de condições de corrida

```bash
cd v3
go test ./... -race
```

## Trabalhar com a documentação

Quem contribui para a documentação pode usar o M-Press sem instalar a CLI do Wails nem as dependências de compilação para desktop. Consulte [Escrever documentação com o M-Press](/contributing/documentation/) para obter instruções sobre downloads dos executáveis com versões fixadas, instalação, visualização ao vivo, validação e fluxo de trabalho de PRs para branches de lançamento.

Os arquivos Markdown em inglês ficam em `docs/mpress/content/`; as traduções ficam nos diretórios dos idiomas. Na raiz do repositório, `mpress dev` disponibiliza a visualização em `http://127.0.0.1:3000`. Interrompa-a com Ctrl+C antes de executar os comandos de validação do guia.

## Depuração

### Depurar código Go

**VS Code:**

Crie `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**Linha de comando:**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### Depurar código específico da plataforma

A depuração específica da plataforma exige ferramentas próprias da plataforma:

- **macOS:** Xcode Instruments
- **Windows:** Visual Studio Debugger
- **Linux:** GDB

## Problemas comuns

### "command not found: wails3"

Adicione o diretório do Wails ao PATH ou use `./wails3` na raiz do projeto.

### "webkitgtk-6.0 not found" ou "webkit2gtk not found" (Linux)

Instale os pacotes de desenvolvimento da pilha para a qual você está compilando:

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### A compilação falha com erros de módulos Go

```bash
cd v3
go mod tidy
go mod download
```

### Erros de "CGO_ENABLED" no Windows

Verifique se há um compilador C (MinGW-w64 via MSYS2) no PATH.

## Próximas etapas

- Consulte os [Padrões de codificação](/contributing/standards/)
- Explore a [Documentação técnica](/contributing/)
- Encontre uma issue na qual trabalhar: [Issues para iniciantes](https://github.com/wailsapp/wails/labels/good%20first%20issue)
