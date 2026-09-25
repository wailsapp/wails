---
title: "Visão geral das caixas de diálogo"
description: "Exiba caixas de diálogo nativas do sistema em seu aplicativo"
slug: "features/dialogs/overview"
sourcePath: "features/dialogs/overview.md"
---

## Caixas de diálogo nativas

O Wails fornece **caixas de diálogo nativas do sistema** que funcionam em todas as plataformas: caixas de mensagem (informação, aviso, erro e pergunta), caixas de diálogo de arquivo (abrir, salvar e pasta) e janelas de diálogo personalizadas com aparência e comportamento nativos da plataforma.

![Uma caixa de diálogo de pergunta do Wails no macOS com os botões Cancelar e Descartar](/assets/screenshots/dialog-question-macos.png)

A mesma API é renderizada de acordo com as convenções de cada plataforma compatível. Este exemplo no macOS mostra uma caixa de diálogo de pergunta anexada à janela do Wails, incluindo os botões padrão e de cancelamento.

## Início rápido

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()

// Question dialog with button callbacks
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()

// File open dialog
path, _ := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()
```

**É só isso!** Caixas de diálogo nativas com o mínimo de código.

## Como acessar as caixas de diálogo

As caixas de diálogo são acessadas por meio do gerenciador `app.Dialog`:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
app.Dialog.OpenFile()
app.Dialog.SaveFile()
```

## Tipos de caixa de diálogo

### Caixa de diálogo de informação

Exiba mensagens simples:

```go
app.Dialog.Info().
    SetTitle("Welcome").
    SetMessage("Welcome to our application!").
    Show()
```

**Casos de uso:**

- Mensagens de sucesso
- Avisos informativos
- Confirmações de conclusão

### Caixa de diálogo de aviso

Exiba avisos:

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

**Casos de uso:**

- Avisos não críticos
- Avisos de descontinuação
- Mensagens de atenção

### Caixa de diálogo de erro

Exiba erros:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

**Casos de uso:**

- Mensagens de erro
- Notificações de falha
- Tratamento de exceções

### Caixa de diálogo de pergunta

Faça perguntas aos usuários e trate as respostas por meio de callbacks dos botões:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm Delete").
    SetMessage("Are you sure you want to delete this file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    deleteFile()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)
dialog.SetCancelButton(cancelBtn)
dialog.Show()
```

**Casos de uso:**

- Confirmar ações
- Perguntas de Sim/Não
- Múltipla escolha

## Caixas de diálogo de arquivo

### Caixa de diálogo Abrir Arquivo

Selecione arquivos para abrir:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    openFile(path)
}
```

**Seleção múltipla:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err == nil {
    for _, path := range paths {
        processFile(path)
    }
}
```

### Caixa de diálogo Salvar Arquivo

Escolha onde salvar:

```go
path, err := app.Dialog.SaveFile().
    SetFilename("document.txt").
    AddFilter("Text Files", "*.txt").
    AddFilter("All Files", "*.*").
    PromptForSingleSelection()

if err == nil && path != "" {
    saveFile(path)
}
```

### Caixa de diálogo Selecionar Pasta

Escolha um diretório usando a caixa de diálogo Abrir Arquivo com a seleção de diretórios habilitada:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err == nil && path != "" {
    exportToFolder(path)
}
```

## Opções da caixa de diálogo

### Título e mensagem

```go
dialog := app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed successfully!")
```

### Botões

**Botão padrão para caixas de diálogo simples:**

As caixas de diálogo de informação, aviso e erro exibem um botão "OK" padrão:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

**Botões personalizados para caixas de diálogo de pergunta:**

Use `AddButton()` para adicionar botões. Esse método retorna um `*Button` que você pode configurar com callbacks:

```go
dialog := app.Dialog.Question().
    SetMessage("Choose action")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    discardChanges()
})

cancel := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Botões padrão e Cancelar:**

Use `SetDefaultButton()` para especificar qual botão será realçado e acionado pela tecla Enter. Use `SetCancelButton()` para especificar qual botão será acionado pela tecla Escape.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
dialog.SetDefaultButton(cancelBtn)  // Safe option highlighted by default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

### Vinculação à janela

Anexe a caixa de diálogo a uma janela específica:

```go
dialog := app.Dialog.Info().
    SetMessage("Window-specific message").
    AttachToWindow(window)

dialog.Show()
```

**Comportamento:**

- A caixa de diálogo aparece centralizada na janela pai
- A janela pai fica desabilitada enquanto a caixa de diálogo é exibida
- A caixa de diálogo se move com a janela pai (macOS)

## Comportamento por plataforma

@tabs{sync-key="platform"}
[macOS]
**Caixas de diálogo no macOS:**

- Aparência nativa do NSAlert
- Seguem o tema do sistema (claro/escuro)
- Oferecem suporte à navegação pelo teclado
- Atalhos padrão (⌘. para Cancelar)
- Recursos de acessibilidade integrados
- Estilo de folha quando anexadas à janela

**Exemplo:**

```go
// Appears as sheet on macOS
dialog := app.Dialog.Question().
    SetMessage("Save changes?").
    AttachToWindow(window)
dialog.AddButton("Yes")
dialog.AddButton("No")
dialog.Show()
```

[Windows]
**Caixas de diálogo no Windows:**

- Aparência nativa do TaskDialog
- Seguem o tema do sistema
- Oferecem suporte à navegação pelo teclado
- Atalhos padrão (Esc para Cancelar)
- Recursos de acessibilidade integrados
- Modal em relação à janela pai

**Exemplo:**

```go
// Modal dialog on Windows
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Operation failed").
    Show()
```

[Linux]
**Caixas de diálogo no Linux:**

- Aparência das caixas de diálogo do GTK
- Seguem o tema da área de trabalho
- Oferecem suporte à navegação pelo teclado
- Integração com o ambiente de desktop
- Varia conforme o ambiente de desktop (GNOME, KDE etc.)

**Exemplo:**

```go
// GTK dialog on Linux
app.Dialog.Info().
    SetMessage("Update complete").
    Show()
```

@end

## Padrões comuns

### Confirmar antes de uma ação destrutiva

```go
func deleteFile(app *application.App, path string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(fmt.Sprintf("Delete %s?", filepath.Base(path)))

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        if err := os.Remove(path); err != nil {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(err.Error()).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Tratamento de erros com caixas de diálogo

```go
func saveDocument(app *application.App, path string, data []byte) {
    if err := os.WriteFile(path, data, 0644); err != nil {
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return
    }

    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()
}
```

### Seleção de arquivos com validação

```go
func selectImageFile(app *application.App) (string, error) {
    path, err := app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    if path == "" {
        return "", errors.New("no file selected")
    }

    // Validate file
    if !isValidImage(path) {
        app.Dialog.Error().
            SetTitle("Invalid File").
            SetMessage("Selected file is not a valid image.").
            Show()
        return "", errors.New("invalid image")
    }

    return path, nil
}
```

### Fluxo de caixas de diálogo em várias etapas

```go
func exportData(app *application.App) {
    // Step 1: Confirm export
    dialog := app.Dialog.Question().
        SetTitle("Export Data").
        SetMessage("Export all data to CSV?")

    exportBtn := dialog.AddButton("Export")
    exportBtn.OnClick(func() {
        // Step 2: Select destination
        path, err := app.Dialog.SaveFile().
            SetFilename("export.csv").
            AddFilter("CSV Files", "*.csv").
            PromptForSingleSelection()

        if err != nil || path == "" {
            return
        }

        // Step 3: Perform export
        if err := performExport(path); err != nil {
            app.Dialog.Error().
                SetTitle("Export Failed").
                SetMessage(err.Error()).
                Show()
            return
        }

        // Step 4: Success
        app.Dialog.Info().
            SetTitle("Export Complete").
            SetMessage("Data exported successfully!").
            Show()
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

## Boas práticas

### ✅ Recomendações

- **Use caixas de diálogo nativas** — Oferecem uma experiência do usuário melhor do que as personalizadas
- **Forneça mensagens claras** — Seja específico
- **Defina títulos adequados** — O contexto é importante
- **Use os botões padrão com cuidado** — Defina a opção segura como padrão
- **Trate o cancelamento** — O usuário pode cancelar
- **Valide os arquivos selecionados** — Verifique os tipos de arquivo

### ❌ O que não fazer

- **Não use caixas de diálogo em excesso** — Elas interrompem o fluxo de trabalho
- **Não as use para mensagens frequentes** — Use notificações
- **Não se esqueça do tratamento de erros** — O usuário pode cancelar
- **Não bloqueie desnecessariamente** — Considere alternativas
- **Não use mensagens genéricas** — Seja específico
- **Não ignore as diferenças entre plataformas** — Teste em todas as plataformas

## Próximas etapas

@cards{cols="2"}
ℹ Caixas de diálogo de mensagem
Caixas de diálogo de informação, aviso e erro.

[Saiba mais →](/features/dialogs/message/)

---
📖 Caixas de diálogo de arquivo
Abertura, salvamento e seleção de pastas.

[Saiba mais →](/features/dialogs/file/)

---
◆ Caixas de diálogo personalizadas
Crie janelas de diálogo personalizadas.

[Saiba mais →](/features/dialogs/custom/)

---
▣ Janelas
Saiba mais sobre o gerenciamento de janelas.

[Saiba mais →](/features/windows/basics/)

@end

---

**Tem alguma dúvida?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de caixas de diálogo](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
