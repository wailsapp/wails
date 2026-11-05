---
title: "Compilação de aplicativos"
description: "Compile e empacote seu aplicativo Wails"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

O Wails v3 usa o [Task](https://taskfile.dev) como sistema de compilação. Os comandos `wails3 build` e `wails3 package` são wrappers práticos para o Task.

## Compilação

Compile para a plataforma atual:

```bash
wails3 build
```

Compile para uma plataforma específica:

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

A saída é gravada no diretório `bin/`.

@note{type="tip"}
A compilação cruzada para macOS ou Linux a partir de outra plataforma requer o Docker. Consulte [Compilações multiplataforma](/guides/build/cross-platform/) para saber como configurá-la.

@end

## Desenvolvimento

Execute seu aplicativo com recarregamento automático:

```bash
wails3 dev
```

Isso inicia um monitor de arquivos que recompila e reinicia seu aplicativo quando há alterações. Por padrão, o servidor de desenvolvimento do frontend é executado na porta 9245.

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## Empacotamento

Empacote seu aplicativo para distribuição:

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

Isso cria pacotes específicos para cada plataforma:

- **Windows**: instalador NSIS — consulte [Empacotamento para Windows](/guides/build/windows/)
- **macOS**: pacote de aplicativo (`.app`) — consulte [Empacotamento para macOS](/guides/build/macos/)
- **Linux**: AppImage, deb e rpm — consulte [Empacotamento para Linux](/guides/build/linux/)

## Tags de compilação personalizadas

Passe tags de compilação personalizadas do Go com a opção `-tags`:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

As tags são encaminhadas como `EXTRA_TAGS` para o Taskfile subjacente. Consulte [Compilação do servidor](/guides/server-build/) e [Empacotamento para Linux — compatibilidade com o GTK3 legado](/guides/build/linux/#legacy-gtk3-support) para obter detalhes.

## Uso direto do Task

Para ter mais controle, use o Task diretamente:

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

Tarefas específicas de plataforma, como `linux:create:deb` ou `darwin:build:universal`, estão disponíveis somente por meio do Task.

## Geração de recursos

Gere novamente os ícones ou atualize a configuração de compilação:

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
