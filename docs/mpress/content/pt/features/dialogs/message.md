---
title: "Caixas de diálogo de mensagens"
description: "Exiba informações, avisos, erros e perguntas"
slug: "features/dialogs/message"
sourcePath: "features/dialogs/message.md"
---

## Caixas de diálogo de mensagens

O Wails fornece **caixas de diálogo de mensagens nativas** com uma aparência apropriada para cada plataforma: caixas de diálogo de informações, avisos, erros e perguntas, com títulos, mensagens e botões personalizáveis. API simples, comportamento nativo e acessibilidade por padrão.

## Como criar caixas de diálogo

As caixas de diálogo de mensagens são acessadas por meio do gerenciador `app.Dialog`:

```go
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

Todos os métodos retornam um `*MessageDialog` que pode ser configurado por meio do encadeamento de métodos.

## Caixa de diálogo de informações

Exiba mensagens informativas:

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

**Casos de uso:**

- Confirmações de sucesso
- Avisos de conclusão
- Mensagens informativas
- Atualizações de status

**Exemplo — Confirmação de salvamento:**

```go
func saveFile(app *application.App, path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        return err
    }

    app.Dialog.Info().
        SetTitle("File Saved").
        SetMessage(fmt.Sprintf("Saved to %s", filepath.Base(path))).
        Show()

    return nil
}
```

## Caixa de diálogo de aviso

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
- Possíveis problemas

**Exemplo — Aviso de espaço em disco:**

```go
func checkDiskSpace(app *application.App) {
    available := getDiskSpace()

    if available < 100*1024*1024 { // Less than 100MB
        app.Dialog.Warning().
            SetTitle("Low Disk Space").
            SetMessage(fmt.Sprintf("Only %d MB available.", available/(1024*1024))).
            Show()
    }
}
```

## Caixa de diálogo de erro

Exiba erros:

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to connect to server.").
    Show()
```

**Casos de uso:**

- Mensagens de erro
- Notificações de falha
- Tratamento de exceções
- Problemas críticos

**Exemplo — Erro de rede:**

```go
func fetchData(app *application.App, url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        app.Dialog.Error().
            SetTitle("Network Error").
            SetMessage(fmt.Sprintf("Failed to connect: %v", err)).
            Show()
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

## Caixa de diálogo de pergunta

Faça perguntas aos usuários e trate as respostas por meio de callbacks dos botões:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Save changes before closing?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveChanges()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Don't close
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

**Casos de uso:**

- Confirmar ações
- Perguntas de Sim/Não
- Múltipla escolha
- Decisões do usuário

**Exemplo — Alterações não salvas:**

```go
func closeDocument(app *application.App) {
    if !hasUnsavedChanges() {
        doClose()
        return
    }

    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        if saveDocument() {
            doClose()
        }
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        doClose()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel button has no callback - just closes the dialog

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

## Opções da caixa de diálogo

### Título e mensagem

```go
dialog := app.Dialog.Info().
    SetTitle("Operation Complete").
    SetMessage("All files have been processed successfully.")
```

**Práticas recomendadas:**

- **Título:** curto e descritivo (2-5 palavras)
- **Mensagem:** clara, específica e acionável
- **Evite jargões:** use linguagem simples

### Botões

**Botão único (Informação/Aviso/Erro):**

As caixas de diálogo de informações, avisos e erros exibem um botão "OK" padrão:

```go
app.Dialog.Info().
    SetMessage("Done!").
    Show()
```

Você também pode adicionar botões personalizados:

```go
dialog := app.Dialog.Info().
    SetMessage("Done!")
ok := dialog.AddButton("Got it!")
dialog.SetDefaultButton(ok)
dialog.Show()
```

**Vários botões (Pergunta):**

Use `AddButton()` para adicionar botões. Esse método retorna um `*Button` que você pode configurar:

```go
dialog := app.Dialog.Question().
    SetMessage("Choose an action")

option1 := dialog.AddButton("Option 1")
option1.OnClick(func() {
    handleOption1()
})

option2 := dialog.AddButton("Option 2")
option2.OnClick(func() {
    handleOption2()
})

option3 := dialog.AddButton("Option 3")
option3.OnClick(func() {
    handleOption3()
})

dialog.Show()
```

**Botões padrão e Cancelar:**

Use `SetDefaultButton()` para especificar qual botão fica realçado e é acionado pela tecla Enter. Use `SetCancelButton()` para especificar qual botão é acionado pela tecla Escape.

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

deleteBtn := dialog.AddButton("Delete")
deleteBtn.OnClick(func() {
    performDelete()
})

cancelBtn := dialog.AddButton("Cancel")
// No callback needed - just dismisses dialog

dialog.SetDefaultButton(cancelBtn)  // Safe option as default
dialog.SetCancelButton(cancelBtn)   // Escape triggers Cancel
dialog.Show()
```

Você também pode usar os métodos fluentes `SetAsDefault()` e `SetAsCancel()` nos botões:

```go
dialog := app.Dialog.Question().
    SetMessage("Delete file?")

dialog.AddButton("Delete").OnClick(func() {
    performDelete()
})

dialog.AddButton("Cancel").SetAsDefault().SetAsCancel()

dialog.Show()
```

**Práticas recomendadas:**

- **1-3 botões:** não sobrecarregue os usuários
- **Rótulos claros:** "Salvar" em vez de "OK"
- **Opção padrão segura:** ação não destrutiva
- **A ordem importa:** coloque primeiro a ação mais provável (exceto Cancelar)

### Ícone personalizado

Defina um ícone personalizado para a caixa de diálogo:

```go
app.Dialog.Info().
    SetTitle("Custom Icon Example").
    SetMessage("Using a custom icon").
    SetIcon(myIconBytes).
    Show()
```

### Vinculação à janela

Vincule a caixa de diálogo a uma janela específica:

```go
dialog := app.Dialog.Question().
    SetMessage("Window-specific question").
    AttachToWindow(window)

dialog.AddButton("OK")
dialog.Show()
```

**Benefícios:**

- A caixa de diálogo aparece na janela correta
- A janela pai fica desabilitada enquanto a caixa de diálogo é exibida
- Melhor experiência do usuário com várias janelas

## Exemplos completos

### Confirmar uma ação destrutiva

```go
func deleteFiles(app *application.App, paths []string) {
    // Confirm deletion
    message := fmt.Sprintf("Delete %d file(s)?", len(paths))
    if len(paths) == 1 {
        message = fmt.Sprintf("Delete %s?", filepath.Base(paths[0]))
    }

    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage(message)

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        // Perform deletion
        var errs []error
        for _, path := range paths {
            if err := os.Remove(path); err != nil {
                errs = append(errs, err)
            }
        }

        // Show result
        if len(errs) > 0 {
            app.Dialog.Error().
                SetTitle("Delete Failed").
                SetMessage(fmt.Sprintf("Failed to delete %d file(s)", len(errs))).
                Show()
        } else {
            app.Dialog.Info().
                SetTitle("Delete Complete").
                SetMessage(fmt.Sprintf("Deleted %d file(s)", len(paths))).
                Show()
        }
    })

    cancelBtn := dialog.AddButton("Cancel")
    dialog.SetDefaultButton(cancelBtn)
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Confirmação de encerramento

```go
func confirmQuit(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Quit").
        SetMessage("You have unsaved work. Are you sure you want to quit?")

    yes := dialog.AddButton("Yes")
    yes.OnClick(func() {
        app.Quit()
    })

    no := dialog.AddButton("No")
    dialog.SetDefaultButton(no)
    dialog.Show()
}
```

### Caixa de diálogo de atualização com opção de download

```go
func showUpdateDialog(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Update").
        SetMessage("A new version is available. The cancel button is selected when pressing escape.")

    download := dialog.AddButton("📥 Download")
    download.OnClick(func() {
        app.Dialog.Info().SetMessage("Downloading...").Show()
    })

    cancel := dialog.AddButton("Cancel")

    dialog.SetDefaultButton(download)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Pergunta com ícone personalizado

```go
func showCustomIconQuestion(app *application.App, iconBytes []byte) {
    dialog := app.Dialog.Question().
        SetTitle("Custom Icon Example").
        SetMessage("Using a custom icon").
        SetIcon(iconBytes)

    likeIt := dialog.AddButton("I like it!")
    likeIt.OnClick(func() {
        app.Dialog.Info().SetMessage("Thanks!").Show()
    })

    notKeen := dialog.AddButton("Not so keen...")
    notKeen.OnClick(func() {
        app.Dialog.Info().SetMessage("Too bad!").Show()
    })

    dialog.SetDefaultButton(likeIt)
    dialog.Show()
}
```

## Práticas recomendadas

### ✅ Faça

- **Seja específico** — "Arquivo salvo em Documentos", não "Sucesso"
- **Use o tipo apropriado** — Erro para erros, Aviso para avisos
- **Forneça contexto** — Inclua os detalhes relevantes
- **Use rótulos claros nos botões** — "Excluir", não "OK"
- **Defina opções padrão seguras** — Use uma ação não destrutiva
- **Trate o cancelamento** — O usuário pode fechar a caixa de diálogo

### ❌ Não faça

- **Não use em excesso** — Isso interrompe o fluxo de trabalho
- **Não use para atualizações frequentes** — Use notificações em vez disso
- **Não use mensagens genéricas** — "Erro" não fornece nenhuma informação
- **Não ignore erros** — Trate os erros de dialog.Show()
- **Não bloqueie desnecessariamente** — Considere alternativas assíncronas
- **Não use jargão técnico** — Use linguagem simples

## Diferenças entre plataformas

### macOS

- Estilo de painel quando anexada a uma janela
- Atalhos de teclado padrão (⌘. para Cancelar)
- Segue automaticamente o tema do sistema
- Recursos de acessibilidade integrados

### Windows

- Caixas de diálogo modais
- Aparência de TaskDialog
- Esc para Cancelar
- Segue o tema do sistema

### Linux

- Caixas de diálogo GTK
- Varia conforme o ambiente de desktop
- Segue o tema do ambiente de desktop
- Navegação padrão pelo teclado

## Próximas etapas

@cards{cols="2"}
📖 Caixas de diálogo de arquivos
Abertura, salvamento e seleção de pastas.

[Saiba mais →](/features/dialogs/file/)

---
◆ Caixas de diálogo personalizadas
Crie janelas de diálogo personalizadas.

[Saiba mais →](/features/dialogs/custom/)

---
● Notificações
Notificações não intrusivas.

[Saiba mais →](/features/notifications/overview/)

---
★ Eventos
Use eventos para comunicação sem bloqueio.

[Saiba mais →](/features/events/system/)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de caixas de diálogo](https://github.com/wailsapp/wails/tree/master/v3/examples/dialogs).
