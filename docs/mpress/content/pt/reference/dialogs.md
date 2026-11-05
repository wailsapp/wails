---
title: "API de caixas de diálogo"
description: "Referência completa das APIs de caixas de diálogo nativas"
slug: "reference/dialogs"
sourcePath: "reference/dialogs.md"
---

## Visão geral

A API de caixas de diálogo fornece métodos para exibir caixas de diálogo nativas de arquivos e de mensagens. Acesse-as por meio do gerenciador `app.Dialog`.

**Tipos de caixa de diálogo:**

- **Caixas de diálogo de arquivos** — caixas de diálogo para abrir e salvar arquivos
- **Caixas de diálogo de mensagens** — caixas de diálogo de informação, erro, aviso e pergunta

Todas as caixas de diálogo são **nativas do sistema operacional** e seguem a aparência e o comportamento da plataforma.

## Como acessar as caixas de diálogo

As caixas de diálogo são acessadas por meio do gerenciador `app.Dialog`:

```go
app.Dialog.OpenFile()
app.Dialog.SaveFile()
app.Dialog.Info()
app.Dialog.Question()
app.Dialog.Warning()
app.Dialog.Error()
```

## Caixas de diálogo de arquivos

### OpenFile()

Cria uma caixa de diálogo para abrir arquivos.

```go
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct
```

**Exemplo:**

```go
dialog := app.Dialog.OpenFile()
```

### Métodos de OpenFileDialogStruct

#### SetTitle()

Define o título da caixa de diálogo.

```go
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct
```

**Exemplo:**

```go
dialog.SetTitle("Select Image")
```

#### AddFilter()

Adiciona um filtro de tipo de arquivo.

```go
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct
```

**Parâmetros:**

- `displayName` — descrição do filtro exibida ao usuário (por exemplo, "Imagens", "Documentos")
- `pattern` — lista de extensões separadas por ponto e vírgula (por exemplo, "*.png;*.jpg")

**Exemplo:**

```go
dialog.AddFilter("Images", "*.png;*.jpg;*.gif").
    AddFilter("Documents", "*.pdf;*.docx").
    AddFilter("All Files", "*.*")
```

#### SetDirectory()

Define o diretório inicial.

```go
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct
```

**Exemplo:**

```go
homeDir, _ := os.UserHomeDir()
dialog.SetDirectory(homeDir)
```

#### CanChooseDirectories()

Ativa ou desativa a seleção de diretórios.

```go
func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct
```

**Exemplo (seleção de pasta):**

```go
// Select folders instead of files
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

#### CanChooseFiles()

Ativa ou desativa a seleção de arquivos.

```go
func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct
```

#### CanCreateDirectories()

Ativa ou desativa a criação de novos diretórios.

```go
func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct
```

#### ShowHiddenFiles()

Exibe ou oculta arquivos ocultos.

```go
func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct
```

#### AttachToWindow()

Vincula a caixa de diálogo a uma janela específica.

```go
func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct
```

#### PromptForSingleSelection()

Exibe a caixa de diálogo e retorna o arquivo selecionado.

```go
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Retornos:**

- `string` — caminho do arquivo selecionado. Considere uma string vazia como "nenhuma seleção" (dependendo do sistema operacional, as implementações da plataforma podem retornar uma string vazia ou um erro não nulo em caso de cancelamento).
- `error` — não nulo se não foi possível exibir a própria caixa de diálogo.

**Exemplo:**

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Image").
    AddFilter("Images", "*.png;*.jpg;*.gif").
    PromptForSingleSelection()

if err != nil {
    // The dialog failed to present (rare).
    return
}
if path == "" {
    // User cancelled.
    return
}

// Use the selected file
processFile(path)
```

#### PromptForMultipleSelection()

Exibe a caixa de diálogo e retorna vários arquivos selecionados.

```go
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error)
```

**Retornos:**

- `[]string` — array de caminhos dos arquivos selecionados
- `error` — erro se a caixa de diálogo falhar

**Exemplo:**

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Select Images").
    AddFilter("Images", "*.png;*.jpg").
    PromptForMultipleSelection()

if err != nil {
    return
}

for _, path := range paths {
    processFile(path)
}
```

### SaveFile()

Cria uma caixa de diálogo para salvar arquivos.

```go
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct
```

**Exemplo:**

```go
dialog := app.Dialog.SaveFile()
```

### Métodos de SaveFileDialogStruct

#### SetTitle()

Define o título da caixa de diálogo.

```go
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct
```

#### SetFilename()

Define o nome de arquivo padrão.

```go
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct
```

**Exemplo:**

```go
dialog.SetFilename("document.pdf")
```

#### AddFilter()

Adiciona um filtro de tipo de arquivo.

```go
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct
```

**Exemplo:**

```go
dialog.AddFilter("PDF Document", "*.pdf").
    AddFilter("Text Document", "*.txt")
```

#### SetDirectory()

Define o diretório inicial.

```go
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct
```

#### AttachToWindow()

Vincula a caixa de diálogo a uma janela específica.

```go
func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct
```

#### PromptForSingleSelection()

Exibe a caixa de diálogo e retorna o caminho de salvamento.

```go
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error)
```

**Exemplo:**

```go
path, err := app.Dialog.SaveFile().
    SetTitle("Save Document").
    SetFilename("untitled.pdf").
    AddFilter("PDF Document", "*.pdf").
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Save to the selected path
saveDocument(path)
```

### Seleção de pasta

Não há uma `SelectFolderDialog` separada. Use `OpenFile()` com as opções de diretório:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select Output Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

if err != nil {
    // User cancelled
    return
}

// Use the selected folder
outputDir = path
```

## Caixas de diálogo de mensagem

Todas as caixas de diálogo de mensagem retornam `*MessageDialog` e compartilham os mesmos métodos.

### Info()

Cria uma caixa de diálogo informativa.

```go
func (dm *DialogManager) Info() *MessageDialog
```

**Exemplo:**

```go
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("File saved successfully!").
    Show()
```

### Error()

Cria uma caixa de diálogo de erro.

```go
func (dm *DialogManager) Error() *MessageDialog
```

**Exemplo:**

```go
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Failed to save file: " + err.Error()).
    Show()
```

### Warning()

Cria uma caixa de diálogo de aviso.

```go
func (dm *DialogManager) Warning() *MessageDialog
```

**Exemplo:**

```go
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Question()

Cria uma caixa de diálogo de pergunta com botões personalizados.

```go
func (dm *DialogManager) Question() *MessageDialog
```

**Exemplo:**

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Do you want to save changes?")

save := dialog.AddButton("Save")
save.OnClick(func() {
    saveDocument()
})

dontSave := dialog.AddButton("Don't Save")
dontSave.OnClick(func() {
    // Continue without saving
})

cancel := dialog.AddButton("Cancel")
cancel.OnClick(func() {
    // Do nothing
})

dialog.SetDefaultButton(save)
dialog.SetCancelButton(cancel)
dialog.Show()
```

### Métodos de MessageDialog

#### SetTitle()

Define o título da caixa de diálogo.

```go
func (d *MessageDialog) SetTitle(title string) *MessageDialog
```

#### SetMessage()

Define a mensagem da caixa de diálogo.

```go
func (d *MessageDialog) SetMessage(message string) *MessageDialog
```

#### SetIcon()

Define um ícone personalizado para a caixa de diálogo.

```go
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog
```

#### AddButton()

Adiciona um botão à caixa de diálogo e retorna esse botão para configuração.

```go
func (d *MessageDialog) AddButton(label string) *Button
```

**Retorna:** `*Button` — A instância do botão para configuração adicional

**Exemplo:**

```go
button := dialog.AddButton("OK")
button.OnClick(func() {
    // Handle click
})
```

#### SetDefaultButton()

Define qual botão será o padrão (ativado ao pressionar Enter).

```go
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog
```

**Exemplo:**

```go
yes := dialog.AddButton("Yes")
no := dialog.AddButton("No")
dialog.SetDefaultButton(yes)
```

#### SetCancelButton()

Define qual botão será o de cancelamento (ativado ao pressionar Escape).

```go
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog
```

**Exemplo:**

```go
ok := dialog.AddButton("OK")
cancel := dialog.AddButton("Cancel")
dialog.SetCancelButton(cancel)
```

#### AttachToWindow()

Vincula a caixa de diálogo a uma janela específica.

```go
func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog
```

#### Show()

Exibe a caixa de diálogo. Os callbacks dos botões processam as respostas do usuário.

```go
func (d *MessageDialog) Show()
```

**Observação:** `Show()` não retorna um valor. Use os callbacks dos botões para processar as respostas do usuário.

### Métodos de botão

#### OnClick()

Define a função de callback executada quando o botão é clicado.

```go
func (b *Button) OnClick(callback func()) *Button
```

#### SetAsDefault()

Marca este botão como o botão padrão.

```go
func (b *Button) SetAsDefault() *Button
```

#### SetAsCancel()

Marca este botão como o botão de cancelamento.

```go
func (b *Button) SetAsCancel() *Button
```

## Exemplos completos

### Exemplo de seleção de arquivo

```go
type FileService struct {
    app *application.App
}

func (s *FileService) OpenImage() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Image").
        AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif").
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SaveDocument(defaultName string) (string, error) {
    path, err := s.app.Dialog.SaveFile().
        SetTitle("Save Document").
        SetFilename(defaultName).
        AddFilter("PDF Document", "*.pdf").
        AddFilter("Text Document", "*.txt").
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}

func (s *FileService) SelectOutputFolder() (string, error) {
    path, err := s.app.Dialog.OpenFile().
        SetTitle("Select Output Folder").
        CanChooseDirectories(true).
        CanChooseFiles(false).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Exemplo de caixa de diálogo de confirmação

```go
func (s *Service) DeleteItem(app *application.App, id string) {
    dialog := app.Dialog.Question().
        SetTitle("Confirm Delete").
        SetMessage("Are you sure you want to delete this item?")

    deleteBtn := dialog.AddButton("Delete")
    deleteBtn.OnClick(func() {
        deleteFromDatabase(id)
    })

    cancelBtn := dialog.AddButton("Cancel")
    // Cancel does nothing

    dialog.SetDefaultButton(cancelBtn) // Default to Cancel for safety
    dialog.SetCancelButton(cancelBtn)
    dialog.Show()
}
```

### Caixa de diálogo para salvar alterações

```go
func (s *Editor) PromptSaveChanges(app *application.App) {
    dialog := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Do you want to save your changes before closing?")

    save := dialog.AddButton("Save")
    save.OnClick(func() {
        s.Save()
        s.Close()
    })

    dontSave := dialog.AddButton("Don't Save")
    dontSave.OnClick(func() {
        s.Close()
    })

    cancel := dialog.AddButton("Cancel")
    // Cancel does nothing, dialog closes

    dialog.SetDefaultButton(save)
    dialog.SetCancelButton(cancel)
    dialog.Show()
}
```

### Processamento de vários arquivos

```go
func (s *Service) ProcessMultipleFiles(app *application.App) error {
    // Select multiple files
    paths, err := app.Dialog.OpenFile().
        SetTitle("Select Files to Process").
        AddFilter("Images", "*.png;*.jpg").
        PromptForMultipleSelection()

    if err != nil {
        return err
    }

    if len(paths) == 0 {
        app.Dialog.Info().
            SetTitle("No Files Selected").
            SetMessage("Please select at least one file.").
            Show()
        return nil
    }

    // Process files
    for _, path := range paths {
        err := processFile(path)
        if err != nil {
            app.Dialog.Error().
                SetTitle("Processing Error").
                SetMessage(fmt.Sprintf("Failed to process %s: %v", path, err)).
                Show()
            continue
        }
    }

    // Show completion
    app.Dialog.Info().
        SetTitle("Complete").
        SetMessage(fmt.Sprintf("Successfully processed %d files", len(paths))).
        Show()

    return nil
}
```

### Tratamento de erros com caixas de diálogo

```go
func (s *Service) SaveFile(app *application.App, data []byte) error {
    // Select save location
    path, err := app.Dialog.SaveFile().
        SetTitle("Save File").
        SetFilename("data.json").
        AddFilter("JSON File", "*.json").
        PromptForSingleSelection()

    if err != nil {
        // User cancelled - not an error
        return nil
    }

    // Attempt to save
    err = os.WriteFile(path, data, 0644)
    if err != nil {
        // Show error dialog
        app.Dialog.Error().
            SetTitle("Save Failed").
            SetMessage(fmt.Sprintf("Could not save file: %v", err)).
            Show()
        return err
    }

    // Show success
    app.Dialog.Info().
        SetTitle("Success").
        SetMessage("File saved successfully!").
        Show()

    return nil
}
```

### Valores padrão específicos da plataforma

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func (s *Service) GetDefaultDirectory() string {
    homeDir, _ := os.UserHomeDir()

    switch runtime.GOOS {
    case "windows":
        return filepath.Join(homeDir, "Documents")
    case "darwin":
        return filepath.Join(homeDir, "Documents")
    case "linux":
        return filepath.Join(homeDir, "Documents")
    default:
        return homeDir
    }
}

func (s *Service) OpenWithDefaults(app *application.App) (string, error) {
    return app.Dialog.OpenFile().
        SetTitle("Open File").
        SetDirectory(s.GetDefaultDirectory()).
        AddFilter("All Files", "*.*").
        PromptForSingleSelection()
}
```

## Boas práticas

### O que fazer

- **Use caixas de diálogo nativas** — Elas seguem a aparência e o comportamento da plataforma
- **Forneça títulos claros** — Ajude os usuários a entender a finalidade
- **Defina filtros adequados** — Oriente os usuários até os tipos de arquivo corretos
- **Trate o cancelamento** — Verifique se há erros (o usuário pode cancelar)
- **Solicite confirmação para ações destrutivas** — Use caixas de diálogo Question
- **Forneça feedback** — Use caixas de diálogo Info para mensagens de sucesso
- **Defina valores padrão apropriados** — Diretório padrão, nome do arquivo etc.
- **Use callbacks para ações dos botões** — Processe adequadamente as respostas do usuário

### O que não fazer

- **Não ignore erros** — O cancelamento pelo usuário retorna um erro
- **Não use rótulos ambíguos nos botões** — Seja específico: "Salvar"/"Cancelar"
- **Não use diálogos em excesso** — Eles interrompem o fluxo de trabalho
- **Não exiba erros em caso de cancelamento** — Essa é uma ação normal
- **Não se esqueça dos filtros de arquivo** — Ajude os usuários a encontrar os arquivos corretos
- **Não codifique caminhos diretamente** — Use os.UserHomeDir() ou uma alternativa semelhante

## Tipos de diálogo por plataforma

### macOS

- Os diálogos deslizam para baixo a partir da barra de título
- Estilo "folha" anexado à janela pai
- Aparência nativa do macOS

### Windows

- Diálogos padrão do Windows
- Segue as diretrizes de design do Windows
- Aparência moderna do Windows 10/11

### Linux

- Diálogos GTK em sistemas baseados em GTK
- Diálogos Qt em sistemas baseados em Qt
- Adapta-se ao ambiente de desktop

#### Comportamento dos diálogos no Linux

No Linux, a compilação GTK4 padrão usa o **xdg-desktop-portal** para diálogos de arquivo, o que proporciona integração nativa com o ambiente de desktop, mas faz com que algumas opções não tenham efeito. O caminho legado do GTK3 (`-tags gtk3`) mantém controle programático total sobre estas opções:

| Opção | GTK3 (`-tags gtk3`) | GTK4 (padrão) | Observações |
| --- | --- | --- | --- |
| `ShowHiddenFiles()` | ✅ Funciona | ❌ Sem efeito | O usuário controla essa opção pela interface do diálogo (Ctrl+H ou menu) |
| `CanCreateDirectories()` | ✅ Funciona | ❌ Sem efeito | Sempre habilitado no portal |
| `ResolvesAliases()` | ✅ Funciona | ❌ Sem efeito | O portal processa a resolução de links simbólicos |
| `SetButtonText()` | ✅ Funciona | ✅ Funciona | O texto personalizado do botão de confirmação funciona |

**Por que essas limitações existem:** os diálogos do GTK4 baseados em portal delegam o controle da interface ao ambiente de desktop (GNOME, KDE etc.). Isso ocorre por design — o portal oferece uma experiência de usuário consistente entre aplicativos e respeita as preferências do usuário.

@note{type="info"}
A compilação GTK4 padrão usa os diálogos fornecidos pelo portal. Se o seu aplicativo precisar de controle programático total sobre as opções de diálogo acima, compile usando o caminho legado `-tags gtk3` (com suporte até a v3.0.x; removido na v3.1) — consulte [Empacotamento no Linux — Suporte legado ao GTK3](/guides/build/linux/#legacy-gtk3-support).

@end

## Padrões comuns

### Padrão "Salvar como"

```go
func (s *Service) SaveAs(app *application.App, currentPath string) (string, error) {
    // Extract filename from current path
    filename := filepath.Base(currentPath)

    // Show save dialog
    path, err := app.Dialog.SaveFile().
        SetTitle("Save As").
        SetFilename(filename).
        PromptForSingleSelection()

    if err != nil {
        return "", err
    }

    return path, nil
}
```

### Padrão "Abrir recente"

```go
func (s *Service) OpenRecent(app *application.App, recentPath string) error {
    // Check if file still exists
    if _, err := os.Stat(recentPath); os.IsNotExist(err) {
        dialog := app.Dialog.Question().
            SetTitle("File Not Found").
            SetMessage("The file no longer exists. Remove from recent files?")

        remove := dialog.AddButton("Remove")
        remove.OnClick(func() {
            s.removeFromRecent(recentPath)
        })

        cancel := dialog.AddButton("Cancel")
        dialog.SetCancelButton(cancel)
        dialog.Show()

        return err
    }

    return s.openFile(recentPath)
}
```
