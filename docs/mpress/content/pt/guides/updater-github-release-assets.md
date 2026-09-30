---
title: "Ativos de versões do GitHub para o atualizador"
description: "Como o atualizador do Wails seleciona artefatos do aplicativo nas versões do GitHub e evita pacotes de instalação."
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

O provedor do GitHub Releases seleciona um ativo da versão usando o `AssetMatcher` configurado. Quando `AssetMatcher` é `nil`, ele usa `github.DefaultAssetMatcher`.

## Correspondência padrão

O mecanismo de correspondência padrão procura a plataforma e a arquitetura atuais no nome de arquivo de cada ativo. Ele reconhece aliases comuns de arquitetura, incluindo:

- `amd64`, `x86_64` e `x64`
- `arm64` e `aarch64`
- `386`, `i386`, `x86` e `ia32`

Arquivos auxiliares, como assinaturas e somas de verificação, são ignorados.

## Ativos de instalação

Uma versão do GitHub pode conter tanto o binário do aplicativo usado pelo atualizador quanto um instalador convencional destinado à primeira instalação. O mecanismo de correspondência padrão ignora os ativos cujo nome de arquivo, convertido para letras minúsculas:

- contém `-installer.`
- contém `_installer.`
- é exatamente `installer.exe`

Por exemplo, considerando estes ativos para Windows:

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher` seleciona `myapp-windows-amd64.exe` e ignora o instalador. Isso impede que o atualizador substitua o aplicativo em execução por um executável de instalação do NSIS ou de outro instalador empacotado de forma semelhante.

A verificação é deliberadamente restrita. Nomes de aplicativos que apenas contêm a palavra `installer` continuam válidos, incluindo:

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## Esquemas de nomenclatura personalizados

Configure `AssetMatcher` quando os ativos da versão não seguirem a convenção de nomenclatura baseada em plataforma e arquitetura ou quando você precisar de uma filtragem diferente para instaladores:

```go
import (
    "strings"

    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

gh, err := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, asset := range assets {
            name := strings.ToLower(asset.Name)
            if strings.Contains(name, req.Platform) &&
                strings.Contains(name, req.Arch) &&
                !strings.Contains(name, "-setup.") {
                return i
            }
        }
        return -1
    },
})
```

Um mecanismo de correspondência personalizado substitui completamente `DefaultAssetMatcher` e, portanto, é responsável por excluir assinaturas, somas de verificação, instaladores e quaisquer outros ativos que não devam ser instalados como uma atualização do aplicativo.
