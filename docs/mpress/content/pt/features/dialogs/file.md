---
title: "Caixas de diálogo de arquivos"
description: "Caixas de diálogo para abrir, salvar e selecionar pastas"
slug: "features/dialogs/file"
sourcePath: "features/dialogs/file.md"
---

## Caixas de diálogo de arquivos

O Wails fornece **caixas de diálogo de arquivos nativas**, com a aparência apropriada para cada plataforma, para abrir e salvar arquivos e selecionar pastas. Uma API simples com filtragem por tipo de arquivo, suporte à seleção de vários arquivos e locais padrão.

![Um seletor de arquivos nativo do macOS aberto por um aplicativo Wails](/assets/screenshots/file-dialog-macos.png)

A API do Wails delega a seleção ao seletor do sistema operacional, preservando a navegação, a filtragem e o comportamento de seleção conhecidos.

## Como criar caixas de diálogo de arquivos

As caixas de diálogo de arquivos são acessadas por meio do gerenciador `app.Dialog`:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Caixa de diálogo Abrir arquivo

Selecione os arquivos que deseja abrir:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

openFile(path)
```

**Casos de uso:**

- Abrir documentos
- Importar arquivos
- Carregar imagens
- Selecionar arquivos de configuração

### Seleção de um único arquivo

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open Document").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    // User cancelled or error occurred
    return
}

// Use selected file
data, _ := os.ReadFile(path)
```

### Seleção de vários arquivos

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
    PromptForMultipleSelection()

if err != nil {
    return
}

// Process all selected files
for _, path := range paths {
    processFile(path)
}
```

### Com diretório padrão

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Open File").
    SetDirectory("/Users/me/Documents").
    PromptForSingleSelection()
```

## Caixa de diálogo Salvar arquivo

Escolha onde salvar:

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

saveFile(path, data)
```

**Casos de uso:**

- Salvar documentos
- Exportar dados
- Criar novos arquivos
- Salvar como...

### Com nome de arquivo padrão

```go
path, err := app.Dialog.SaveFile().
    SetFilename("export.csv").
    AddFilter("CSV Files", "*.csv").
    PromptForSingleSelection()
```

### Com diretório padrão

```go
path, err := app.Dialog.SaveFile().
    SetDirectory("/Users/me/Documents").
    SetFilename("untitled.txt").
    PromptForSingleSelection()
```

### Confirmação de substituição

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

// Check if file exists
if _, err := os.Stat(path); err == nil {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Overwrite").
        SetMessage("File already exists. Overwrite?")

    overwriteBtn := dialog.AddButton("Overwrite")
    overwriteBtn.OnClick(func() {
        saveFile(path, data)
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
    return
}

saveFile(path, data)
```

## Caixa de diálogo Selecionar pasta

Escolha um diretório usando a caixa de diálogo para abrir arquivos com a seleção de diretórios habilitada:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil || path == "" {
    return
}

exportToFolder(path)
```

**Casos de uso:**

- Escolher o diretório de saída
- Selecionar o espaço de trabalho
- Escolher o local do backup
- Escolher o diretório de instalação

### Com diretório padrão

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    SetDirectory("/Users/me/Documents").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Filtros de arquivos

Use o método `AddFilter()` para adicionar filtros de tipo de arquivo às caixas de diálogo. Cada chamada adiciona uma nova opção de filtro.

### Filtros básicos

```go
path, _ := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()
```

### Várias extensões

Use pontos e vírgulas para especificar várias extensões em um único filtro:

```go
dialog := app.Dialog.OpenFile().
    AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
    AddFilter("Documents", "*.txt;*.doc;*.docx;*.pdf").
    AddFilter("All Files", "*.*")
```

### Formato do padrão

Use **pontos e vírgulas** para separar várias extensões em um único filtro:

```go
// Multiple extensions separated by semicolons
AddFilter("Images", "*.png;*.jpg;*.gif")
```

## Exemplos completos

### Abrir arquivo de imagem

```go
func openImage(app *application.App) (image.Image, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.bmp").
        PromptForSingleSelection()

    if err != nil {
        return nil, err
    }

    if path == "" {
        return nil, errors.New("no file selected")
    }

    // Open and decode image
    file, err := os.Open(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Open Failed").
            SetMessage(err.Error()).
            Show()
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Image").
            SetMessage("Could not decode image file.").
            Show()
        return nil, err
    }

    return img, nil
}
```

### Salvar documento com validação

```go
func saveDocument(app *application.App, content string) {
    path, err := app.Dialog.SaveFile().
        SetFilename("document.txt").
        AddFilter("Text Files", "*.txt").
        AddFilter("Markdown Files", "*.md").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Validate extension
    ext := filepath.Ext(path)
    if ext != ".txt" && ext != ".md" {
        dialog := app.Dialog.Question().
            SetTitle("Confirm Extension").
            SetMessage(fmt.Sprintf("Save as %s file?", ext))

        saveBtn := dialog.AddButton("Save")
        saveBtn.OnClick(func() {
            doSave(app, path, content)
        })

        cancelBtn := dialog.AddButton("Cancel")
        dialog.SetDefaultButton(cancelBtn)
        dialog.SetCancelButton(cancelBtn)
        dialog.Show()
        return
    }

    doSave(app, path, content)
}

func doSave(app *application.App, path, content string) {
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Saved").
        SetMessage("Document saved successfully!").
        Show()
}
```

### Processamento de arquivos em lote

```go
func processMultipleFiles(app *application.App) {
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil || len(paths) == 0 {
        return
    }

    // Confirm processing
    dialog := app.Dialog.Question().
        SetTitle("Confirm Processing").
        SetMessage(fmt.Sprintf("Process %d file(s)?", len(paths)))

    processBtn := dialog.AddButton("Process")
    processBtn.OnClick(func() {
        // Process files
        var errs []error
        for i, path := range paths {
            if err := processFile(path); err != nil {
                errs = append(errs, err)
            }

            // Update progress
            // app.Event.Emit("progress", map[string]interface{}{
            //     "current": i + 1,
            //     "total":   len(paths),
            // })
            _ = i // suppress unused variable warning in example
        }

        // Show results
        if len(errs) > 0 {
            app.Dialog.Warning().
                SetTitle("Processing Complete").
                SetMessage(fmt.Sprintf("Processed %d files with %d errors.",
                    len(paths), len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Success").
                SetMessage(fmt.Sprintf("Processed %d files successfully!", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Exportar com seleção de pasta

```go
func exportData(app *application.App, data []byte) {
    // Select output folder
    folder, err := app.Dialog.OpenFile().
        SetTitle("Select Export Folder").
        SetDirectory(getDefaultExportFolder()).
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil || folder == "" {
        return
    }

    // Generate filename
    filename := fmt.Sprintf("export_%s.csv",
        time.Now().Format("2006-01-02_15-04-05"))
    path := filepath.Join(folder, filename)

    // Save file
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Export Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Show success with option to open folder
    dialog := app.Dialog.Question().
        SetTitle("Export Complete").
        SetMessage(fmt.Sprintf("Exported to %s", filename))

    openBtn := dialog.AddButton("Open Folder")
    openBtn.OnClick(func() {
        openFolder(folder)
    })

    dialog.AddButton("OK")
    dialog.Show()
}
```

### Importar com validação

```go
func importConfiguration(app *application.App) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Import Configuration").
        AddFilter("JSON Files", "*.json").
        AddFilter("YAML Files", "*.yaml;*.yml").
        PromptForSingleSelection()

    if err != nil || path == "" {
        return
    }

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Read Failed").
            SetMessage(err.Error()).
            Show()
        return
    }

    // Validate configuration
    config, err := parseConfig(data)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Invalid Configuration").
            SetMessage("File is not a valid configuration.").
            Show()
        return
    }

    // Confirm import
    dialog := app.Dialog.Question().
        SetTitle("Confirm Import").
        SetMessage("Import this configuration?")

    importBtn := dialog.AddButton("Import")
    importBtn.OnClick(func() {
        // Apply configuration
        if err := applyConfig(config); err != nil {
            app.Dialog.Error().
                SetTitle("Import Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        app.Dialog.Info().
            SetTitle("Success").
            SetMessage("Configuration imported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Práticas recomendadas

### ✅ Faça

- **Forneça filtros de arquivos** — Ajude os usuários a encontrar arquivos
- **Defina títulos adequados** — Forneça um contexto claro
- **Use diretórios padrão** — Comece em um local lógico
- **Valide as seleções** — Verifique os tipos de arquivo
- **Trate o cancelamento** — O usuário pode cancelar
- **Exiba uma confirmação** — Para ações destrutivas
- **Forneça feedback** — Mensagens de sucesso ou erro

### ❌ Não faça

- **Não ignore a validação** — Verifique os tipos de arquivo
- **Não ignore os erros** — Trate o cancelamento
- **Não use filtros genéricos** — Seja específico
- **Não se esqueça de "Todos os arquivos"** — Sempre inclua essa opção
- **Não codifique caminhos diretamente** — Use o diretório pessoal do usuário
- **Não presuma que o arquivo existe** — Verifique antes de abrir

## Diferenças entre plataformas

### macOS

- NSOpenPanel/NSSavePanel nativo
- Exibido como painel quando anexado a uma janela
- Segue o tema do sistema
- Compatível com a visualização do Quick Look
- Integração com etiquetas e favoritos

### Windows

- Caixas de diálogo nativas para abrir e salvar arquivos
- Segue o tema do sistema
- Integração com arquivos recentes
- Compatibilidade com locais de rede

### Linux

- Seletor de arquivos do GTK
- Varia de acordo com o ambiente de desktop
- Segue o tema do ambiente de desktop
- Compatibilidade com arquivos recentes

## Próximas etapas

@cards{cols="2"}
ℹ Caixas de diálogo de mensagens
Caixas de diálogo de informações, avisos e erros.

[Saiba mais →](/features/dialogs/message/)

---
◆ Caixas de diálogo personalizadas
Crie janelas de diálogo personalizadas.

[Saiba mais →](/features/dialogs/custom/)

---
🚀 Bindings
Chame funções Go a partir do JavaScript.

[Saiba mais →](/features/bindings/methods/)

---
★ Eventos
Use eventos para atualizações de progresso.

[Saiba mais →](/features/events/system/)

@end

---

**Tem alguma dúvida?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de caixas de diálogo de arquivos](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
