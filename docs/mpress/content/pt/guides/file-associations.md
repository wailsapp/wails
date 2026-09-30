---
title: "Associações de arquivos"
description: "Configure associações de arquivos para seu aplicativo Wails"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

Plataformas relevantes: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

As associações de arquivos permitem que seu aplicativo processe tipos específicos de arquivo quando os usuários os abrem. Isso é particularmente útil para editores de texto, visualizadores de imagens ou qualquer aplicativo que trabalhe com formatos de arquivo específicos. Este guia explica como implementar associações de arquivos em seu aplicativo Wails v3.

## Visão geral

Atualmente, o suporte a associações de arquivos no Wails v3 está disponível para:

- Windows (pacotes de instalação NSIS)
- macOS (pacotes de aplicativos)

## Configuração

As associações de arquivos são configuradas no arquivo `config.yml`, localizado no diretório `build` do seu projeto.

### Configuração básica

Para configurar associações de arquivos:

1. Abra `build/config.yml`
2. Adicione suas associações de arquivos na seção `fileAssociations`
3. Execute `wails3 update build-assets` para atualizar os recursos de compilação
4. Defina o campo `FileAssociations` nas opções do aplicativo
5. Empacote seu aplicativo usando `wails3 package`

Veja um exemplo de configuração:

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### Propriedades de configuração

| Propriedade | Descrição | Plataforma |
| --- | --- | --- |
| ext | Extensão do arquivo sem o ponto inicial (por exemplo, `txt`) | Todas |
| name | Nome de exibição do tipo de arquivo | Todas |
| description | Descrição exibida nas propriedades do arquivo | Windows |
| iconName | Nome do arquivo de ícone (sem a extensão) na pasta de compilação | Todas |
| role | Função do aplicativo para esse tipo de arquivo (por exemplo, `Editor`, `Viewer`) | macOS |
| mimeType | Tipo MIME do arquivo (por exemplo, `image/jpeg`) | macOS |

## Escuta de eventos de abertura de arquivo

Para processar eventos de abertura de arquivo em seu aplicativo, você pode escutar o evento `events.Common.ApplicationOpenedWithFile`:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## Tutorial passo a passo

Vamos configurar associações de arquivos para um editor de texto simples:

@steps
### Criar ícones
- Crie ícones para seu tipo de arquivo (tamanhos recomendados: 16x16, 32x32, 48x48, 256x256)
- Salve os ícones na pasta `build` do seu projeto
- Nomeie-os de acordo com sua configuração de `iconName` (por exemplo, `textFileIcon.png`)

@note{type="tip"}
Você pode usar `wails3 generate icons` para gerar os ícones necessários. Execute `wails3 generate icons --help` para obter mais informações.

@end

- No macOS, adicione uma instrução de cópia como `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources` à tarefa `create:app:bundle:`.

### Configurar associações de arquivos
Edite o arquivo `build/config.yml` para adicionar suas associações de arquivos:

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### Atualizar recursos de compilação
Execute o comando a seguir para atualizar os recursos de compilação:

```bash
wails3 update build-assets
```

### Definir associações de arquivos nas opções do aplicativo
No arquivo `main.go`, defina o campo `FileAssociations` nas opções do aplicativo:

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="Por que as extensões de arquivo são obrigatórias tanto na configuração do aplicativo quanto em config.yml?"}
No Windows, quando um arquivo é aberto por meio de uma associação de arquivos, o aplicativo é iniciado com o nome do arquivo como seu primeiro argumento. O aplicativo não tem como saber se o primeiro argumento é um arquivo ou um argumento de linha de comando. Por isso, ele usa o campo `FileAssociations` nas opções do aplicativo para determinar se o primeiro argumento é ou não um arquivo associado.

@end

### Empacotar seu aplicativo
Empacote seu aplicativo usando o comando a seguir:

```bash
wails3 package
```

O aplicativo empacotado será criado no diretório `bin`. Depois, você poderá instalar e testar o aplicativo.

## Observações adicionais

- Recomenda-se fornecer os ícones no formato PNG na pasta de compilação
- Para testar as associações de arquivos, é necessário instalar o aplicativo empacotado

@end
