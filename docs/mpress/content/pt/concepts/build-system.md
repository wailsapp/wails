---
title: "Sistema de compilação"
description: "Entenda como o Wails compila e empacota seu aplicativo"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## Sistema de compilação unificado

O Wails oferece um **sistema de compilação unificado** que compila o código Go, agrupa os recursos do frontend, incorpora tudo em um único executável e gerencia compilações específicas de cada plataforma — tudo com um único comando.

```bash
wails3 build
```

**Saída:** executável nativo com tudo incorporado.

## Visão geral do processo de compilação

**[Espaço reservado para o diagrama do processo de compilação]**

## Fases da compilação

### 1. Fase de análise

O Wails examina seu código Go para compreender seus serviços:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**O que o Wails extrai:**

- Nome do serviço: `GreetService`
- Nome do método: `Greet`
- Tipos dos parâmetros: `string`
- Tipos de retorno: `string`

**Usado para:** gerar bindings TypeScript

### 2. Fase de geração

#### Bindings TypeScript

O Wails gera bindings com segurança de tipos:

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**Benefícios:**

- Segurança de tipos completa
- Preenchimento automático na IDE
- Erros em tempo de compilação
- Comentários JSDoc

#### Compilação do frontend

Seu empacotador de frontend é executado (Vite, webpack etc.):

```bash
# Vite example
vite build --outDir dist
```

**O que acontece:**

- JavaScript/TypeScript compilado
- CSS processado e minificado
- Recursos otimizados
- Mapas de código-fonte gerados (somente em desenvolvimento)
- Saída em `frontend/dist/`

### 3. Fase de compilação

#### Compilação do Go

O código Go é compilado com otimizações:

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**Opções:**

- `-s`: remove a tabela de símbolos
- `-w`: remove as informações de depuração DWARF
- Resultado: binário menor (redução de aproximadamente 30%)

**Específico da plataforma:**

- Windows: `.exe` com ícone incorporado
- macOS: estrutura de pacote `.app`
- Linux: binário ELF

#### Incorporação de recursos

Os recursos do frontend são incorporados ao binário Go:

```go
//go:embed frontend/dist
var assets embed.FS
```

**Resultado:** um único executável com tudo dentro dele.

### 4. Saída

**Um único binário nativo:**

- Windows: `myapp.exe` (aproximadamente 15 MB)
- macOS: `myapp.app` (aproximadamente 15 MB)
- Linux: `myapp` (aproximadamente 15 MB)

**Sem dependências** (exceto a WebView do sistema).

## Desenvolvimento versus produção

@tabs{sync-key="mode"}
[Desenvolvimento (wails3 dev)]
**Otimizado para velocidade:**

```bash
wails3 dev
```

**O que acontece:**

1. Inicia o servidor de desenvolvimento do frontend (por padrão, Vite na porta 9245)
2. Compila o Go sem otimizações
3. Inicia o aplicativo apontando para o servidor de desenvolvimento
4. Ativa o recarregamento automático
5. Inclui mapas de código-fonte

**Características:**

- **Recompilações rápidas** (&lt;1 s para alterações no frontend)
- **Sem incorporação de recursos** (servidos pelo servidor de desenvolvimento)
- **Símbolos de depuração** incluídos
- **Mapas de código-fonte** ativados
- **Logs detalhados**

**Tamanho do arquivo:** maior (aproximadamente 50 MB com símbolos de depuração)

[Produção (wails3 build)]
**Otimizado para tamanho e desempenho:**

```bash
wails3 build
```

**O que acontece:**

1. Compila o frontend para produção (minificado)
2. Compila o Go com otimizações
3. Remove os símbolos de depuração
4. Incorpora os recursos
5. Cria um único binário

**Características:**

- **Código otimizado** (minificado e com tree shaking)
- **Recursos incorporados** (sem arquivos externos)
- **Símbolos de depuração removidos**
- **Sem mapas de código-fonte**
- **Logs mínimos**

**Tamanho do arquivo:** menor (aproximadamente 15 MB)

@end

## Comandos de compilação

### Compilação básica

```bash
wails3 build
```

**Saída:** `bin/<APP_NAME>` (ou `bin/<APP_NAME>.exe` no Windows). O diretório `bin/` fica na raiz do projeto.

`wails3 build` é um wrapper simples para `wails3 task build`. A única opção de compilação que ele encaminha é `--tags`, que se torna a variável `EXTRA_TAGS` do Taskfile:

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build` não aceita as opções `-platform`, `-o`, `-skipbindings`, `-clean`, `-debug`, `-devbuild`, `-icon`, `-ldflags` nem `-package`. A compilação cruzada, os caminhos de saída, os ícones e o empacotamento são controlados pelo Taskfile do projeto (`Taskfile.yml` + `build/config.yml`).

### Compilações multiplataforma e específicas de cada plataforma

As compilações para cada plataforma são disponibilizadas como tarefas do Taskfile nos namespaces `darwin:` / `windows:` / `linux:` (definidos em `build/Taskfile.<platform>.yml`). Por exemplo:

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

Para ver todas as tarefas disponíveis no projeto atual:

```bash
wails3 task --list
```

### Ícones e empacotamento

Gere ícones para as plataformas (`build/icons.icns`, `build/icon.ico` etc.) a partir de um PNG de origem:

```bash
wails3 generate icons -input appicon.png
```

Crie instaladores/pacotes específicos de cada plataforma:

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## Configuração da compilação

### Taskfile.yml

Os projetos Wails 3 usam o [Taskfile](https://taskfile.dev/) para orquestrar a compilação. O `Taskfile.yml` raiz inclui arquivos de tarefas específicos de cada plataforma provenientes de `build/`:

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

Execute as tarefas com `wails3 task <name>` ou `task <name>`:

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### Configuração do projeto: `build/config.yml`

Os metadados do projeto (nome, identificador, versão, valores de info-plist, configurações do NSIS, campos de `.desktop`, protocolos personalizados etc.) ficam em `build/config.yml`. O Taskfile lê esse arquivo ao gerar ícones, manifestos, instaladores e itens semelhantes. No Wails 3, **não** existe um arquivo `build/build.json`.

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

Execute `wails3 generate build-assets` (ou `wails3 update build-assets`) para atualizar, com base nessa configuração, os recursos de compilação específicos de cada plataforma.

## Incorporação de recursos

### Como funciona

O Wails usa o pacote `embed` do Go:

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**Durante a compilação:**

1. O frontend é compilado em `frontend/dist/`
2. A diretiva `//go:embed` inclui os arquivos
3. Os arquivos são compilados no binário
4. O binário contém tudo

**Em tempo de execução:**

1. O aplicativo é iniciado
2. Os recursos são servidos da memória
3. Não há E/S de disco para os recursos
4. Carregamento rápido

### Recursos personalizados

Incorpore arquivos adicionais:

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## Otimizações da compilação

### Otimizações do frontend

**Vite (padrão):**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**Resultados:**

- JavaScript minificado (redução de aproximadamente 70%)
- CSS minificado (redução de aproximadamente 60%)
- Imagens otimizadas
- Tree-shaking aplicado

### Otimizações do Go

**Flags do compilador:**

```bash
-ldflags="-s -w"
```

- `-s`: remove a tabela de símbolos (redução de aproximadamente 10%)
- `-w`: remove as informações de depuração DWARF (redução de aproximadamente 20%)

**Otimizações adicionais:**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`: define valores de variáveis durante a compilação
- Útil para números de versão e datas de compilação

### Compactação do binário

**UPX (opcional):**

```bash
# After building
upx --best bin/myapp.exe
```

**Resultados:**

- Redução de tamanho de aproximadamente 50%
- Inicialização ligeiramente mais lenta (aproximadamente 100 ms)
- Não recomendado para macOS (problemas com assinatura de código)

## Compilações específicas de cada plataforma

### Windows

**Saída:** `myapp.exe`

**Inclui:**

- Ícone do aplicativo
- Informações da versão
- Manifesto (configurações do UAC)

**Ícone:**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

Em seguida, a etapa `tool package` do Windows incorpora o `.ico` gerado ao executável.

**Manifesto:**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**Saída:** `myapp.app` (pacote do aplicativo)

**Estrutura:**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**Binário universal:**

O Taskfile do macOS fornece uma tarefa `darwin:build:universal` (e `darwin:package:universal`) que compila as duas arquiteturas e as combina por meio de `wails3 tool lipo`:

```bash
wails3 task darwin:build:universal
```

### Linux

**Saída:** `myapp` (binário ELF)

**Dependências:**

- GTK3
- WebKitGTK

**Arquivo desktop:**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**Instalação:**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## Desempenho da compilação

### Tempos típicos de compilação

| Fase | Tempo | Observações |
| --- | --- | --- |
| Análise | &lt;1 s | Varredura do código Go |
| Geração de bindings | &lt;1 s | Geração de TypeScript |
| Compilação do frontend | 5-30 s | Depende do tamanho do projeto |
| Compilação do Go | 2-10 s | Depende do tamanho do código |
| Incorporação de recursos | &lt;1 s | Incorporação do frontend |
| **Total** | **10-45 s** | Primeira compilação |
| **Incremental** | **5-15 s** | Compilações subsequentes |

### Como acelerar as compilações

**1. Use o cache de compilação:**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. Execute apenas o necessário:**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. Compilações paralelas (várias máquinas/CI):**

A compilação cruzada entre Linux/Windows/macOS no v3 geralmente ocorre em um contêiner Docker `wails-cross` ou em executores dedicados para cada plataforma — o próprio `wails3 build` tem como destino o sistema operacional do host. Consulte [Compilações multiplataforma](/guides/build/cross-platform/) para conhecer os fluxos de trabalho compatíveis.

**4. Use ferramentas mais rápidas:**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## Solução de problemas

### Falha na compilação

**Sintoma:** `wails3 build` é encerrado com erro

**Causas comuns:**

1. **Erro de compilação do Go**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **Erro de compilação do frontend**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **Dependências ausentes**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### Binário grande demais

**Sintoma:** o binário tem mais de 50 MB

**Soluções:**

1. **Remova os símbolos de depuração** (o Taskfile fornecido já passa `-ldflags="-s -w"` para `go build`).

2. **Verifique os assets incorporados**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **Use a compactação UPX**
  ```bash
  upx --best bin/myapp.exe
  ```


### Compilações lentas

**Sintoma:** as compilações levam mais de 1 minuto

**Soluções:**

1. **Use o cache de compilação**
  - O cache do Go é automático
  - O cache do frontend (Vite) é automático


2. **Execute apenas a tarefa necessária**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **Otimize a compilação do frontend**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## Práticas recomendadas

### ✅ Faça

- **Use `wails3 dev` durante o desenvolvimento** — Iteração rápida
- **Use `wails3 build` para lançamentos** — Saída otimizada
- **Atribua versões às compilações** — Use `-ldflags` para incorporar a versão
- **Teste as compilações nas plataformas de destino** — A compilação cruzada não é perfeita
- **Mantenha rápidas as compilações do frontend** — Otimize a configuração do empacotador
- **Use o cache de compilação** — Acelera as compilações subsequentes

### ❌ Não faça

- **Não faça commit do diretório `build/`** — Adicione-o a `.gitignore`
- **Não deixe de testar as compilações** — Sempre teste antes do lançamento
- **Não incorpore assets desnecessários** — Mantenha os binários pequenos
- **Não use compilações de depuração em produção** — Use compilações otimizadas
- **Não se esqueça da assinatura de código** — Ela é obrigatória para distribuição

## Próximas etapas

**Compilação de aplicações** — Guia detalhado para compilar e empacotar [Saiba mais →](/guides/build/building/)

**Compilações multiplataforma** — Compile para todas as plataformas em uma única máquina [Saiba mais →](/guides/build/cross-platform/)

**Criação de instaladores** — Crie instaladores para usuários finais [Saiba mais →](/guides/installers/)

---

**Tem dúvidas sobre compilação?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de compilação](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/build).
