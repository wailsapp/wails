---
title: "Atalhos de teclado"
description: "Registre atalhos de teclado globais para acessar funcionalidades rapidamente"
slug: "features/keyboard/shortcuts"
sourcePath: "features/keyboard/shortcuts.md"
---

O Wails oferece um sistema avançado de associações de teclas que permite registrar atalhos de teclado globais, válidos em todas as janelas do aplicativo. Assim, os usuários podem acessar funcionalidades rapidamente sem navegar por menus.

## Como acessar o gerenciador de associações de teclas

Acesse o gerenciador de associações de teclas pela propriedade `KeyBindings` da instância do aplicativo:

```go
app := application.New(application.Options{
    Name: "Keyboard Shortcuts Demo",
})

// Access the key binding manager
keyBindings := app.KeyBinding
```

## Como adicionar associações de teclas

### Associação de teclas básica

Registre um atalho de teclado simples:

```go
app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
    // Handle save action
    app.Logger.Info("Save shortcut triggered")
    // Perform save operation...
})
```

### Várias associações de teclas

Registre vários atalhos para operações comuns:

```go
// File operations
app.KeyBinding.Add("Ctrl+N", func(window application.Window) {
    // New file
    window.EmitEvent("file:new", nil)
})

app.KeyBinding.Add("Ctrl+O", func(window application.Window) {
    // Open file
    dialog := app.Dialog.OpenFile()
    if file, err := dialog.PromptForSingleSelection(); err == nil {
        window.EmitEvent("file:open", file)
    }
})

app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
    // Save file
    window.EmitEvent("file:save", nil)
})

// Edit operations
app.KeyBinding.Add("Ctrl+Z", func(window application.Window) {
    // Undo
    window.EmitEvent("edit:undo", nil)
})

app.KeyBinding.Add("Ctrl+Y", func(window application.Window) {
    // Redo (Windows/Linux)
    window.EmitEvent("edit:redo", nil)
})

app.KeyBinding.Add("Cmd+Shift+Z", func(window application.Window) {
    // Redo (macOS)
    window.EmitEvent("edit:redo", nil)
})
```

## Aceleradores de associações de teclas

### Formato dos aceleradores

As associações de teclas usam um formato padrão de aceleradores com teclas modificadoras e outras teclas:

```go
// Modifier keys
"Ctrl+S"        // Control + S
"Cmd+S"         // Command + S (macOS)
"Alt+F4"        // Alt + F4
"Shift+Ctrl+Z"  // Shift + Control + Z

// Function keys
"F1"            // F1 key
"Ctrl+F5"       // Control + F5

// Special keys
"Escape"        // Escape key
"Enter"         // Enter key
"Space"         // Spacebar
"Tab"           // Tab key
"Backspace"     // Backspace key
"Delete"        // Delete key

// Arrow keys
"Up"            // Up arrow
"Down"          // Down arrow
"Left"          // Left arrow
"Right"         // Right arrow
```

### Aceleradores específicos da plataforma

Trate as diferenças entre plataformas para atalhos comuns:

```go
import "runtime"

// Cross-platform save shortcut
if runtime.GOOS == "darwin" {
    app.KeyBinding.Add("Cmd+S", saveHandler)
} else {
    app.KeyBinding.Add("Ctrl+S", saveHandler)
}

// Or register both
app.KeyBinding.Add("Ctrl+S", saveHandler)
app.KeyBinding.Add("Cmd+S", saveHandler)
```

## Como gerenciar associações de teclas

### Como remover associações de teclas

Remova as associações de teclas quando elas não forem mais necessárias:

```go
// Remove a specific key binding
app.KeyBinding.Remove("Ctrl+S")

// Example: Temporary key binding for a modal
app.KeyBinding.Add("Escape", func(window application.Window) {
    // Close modal
    window.EmitEvent("modal:close", nil)
    // Remove this temporary binding
    app.KeyBinding.Remove("Escape")
})
```

### Como obter todas as associações de teclas

Recupere todas as associações de teclas registradas:

```go
allBindings := app.KeyBinding.GetAll()
for _, binding := range allBindings {
    app.Logger.Info("Key binding", "accelerator", binding.Accelerator)
}
```

## Uso avançado

### Associações de teclas sensíveis ao contexto

Faça com que as associações de teclas considerem o contexto verificando o estado do aplicativo:

```go
app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
    // Check current application state
    if isEditMode() {
        // Save document
        saveDocument()
    } else if isInSettings() {
        // Save settings
        saveSettings()
    } else {
        app.Logger.Info("Save not available in current context")
    }
})
```

### Ações específicas da janela

As associações de teclas recebem a janela ativa, permitindo comportamentos específicos para cada janela:

```go
app.KeyBinding.Add("F11", func(window application.Window) {
    // Toggle fullscreen for the active window
    window.ToggleFullscreen()
})

app.KeyBinding.Add("Ctrl+W", func(window application.Window) {
    // Close the active window
    window.Close()
})
```

### Gerenciamento dinâmico de associações de teclas

Adicione e remova associações de teclas dinamicamente com base no estado do aplicativo:

```go
func enableEditMode() {
    // Add edit-specific key bindings
    app.KeyBinding.Add("Ctrl+B", func(window application.Window) {
        window.EmitEvent("format:bold", nil)
    })
    
    app.KeyBinding.Add("Ctrl+I", func(window application.Window) {
        window.EmitEvent("format:italic", nil)
    })
    
    app.KeyBinding.Add("Ctrl+U", func(window application.Window) {
        window.EmitEvent("format:underline", nil)
    })
}

func disableEditMode() {
    // Remove edit-specific key bindings
    app.KeyBinding.Remove("Ctrl+B")
    app.KeyBinding.Remove("Ctrl+I")
    app.KeyBinding.Remove("Ctrl+U")
}
```

## Considerações específicas da plataforma

@tabs
[macOS]
No macOS:

- Use `Cmd` em vez de `Ctrl` para atalhos padrão
- `Cmd+Q` costuma ser reservado para encerrar o aplicativo
- `Cmd+H` oculta o aplicativo
- `Cmd+M` minimiza as janelas
- Considere os atalhos de teclado padrão do macOS

Padrões comuns do macOS:

```go
app.KeyBinding.Add("Cmd+N", newFileHandler)      // New
app.KeyBinding.Add("Cmd+O", openFileHandler)     // Open
app.KeyBinding.Add("Cmd+S", saveFileHandler)     // Save
app.KeyBinding.Add("Cmd+Z", undoHandler)         // Undo
app.KeyBinding.Add("Cmd+Shift+Z", redoHandler)   // Redo
app.KeyBinding.Add("Cmd+C", copyHandler)         // Copy
app.KeyBinding.Add("Cmd+V", pasteHandler)        // Paste
```

[Windows]
No Windows:

- Use `Ctrl` para atalhos padrão
- `Alt+F4` fecha aplicativos
- `F1` normalmente abre a ajuda
- Considere as convenções de teclado do Windows

Padrões comuns do Windows:

```go
app.KeyBinding.Add("Ctrl+N", newFileHandler)     // New
app.KeyBinding.Add("Ctrl+O", openFileHandler)    // Open
app.KeyBinding.Add("Ctrl+S", saveFileHandler)    // Save
app.KeyBinding.Add("Ctrl+Z", undoHandler)        // Undo
app.KeyBinding.Add("Ctrl+Y", redoHandler)        // Redo
app.KeyBinding.Add("Ctrl+C", copyHandler)        // Copy
app.KeyBinding.Add("Ctrl+V", pasteHandler)       // Paste
app.KeyBinding.Add("F1", helpHandler)            // Help
```

[Linux]
No Linux:

- Em geral, segue as convenções do Windows com `Ctrl`
- Pode variar conforme o ambiente de desktop
- Considere os atalhos padrão do GNOME/KDE
- Alguns ambientes de desktop reservam determinados atalhos

Padrões comuns do Linux:

```go
app.KeyBinding.Add("Ctrl+N", newFileHandler)     // New
app.KeyBinding.Add("Ctrl+O", openFileHandler)    // Open
app.KeyBinding.Add("Ctrl+S", saveFileHandler)    // Save
app.KeyBinding.Add("Ctrl+Z", undoHandler)        // Undo
app.KeyBinding.Add("Ctrl+Shift+Z", redoHandler)  // Redo
app.KeyBinding.Add("Ctrl+C", copyHandler)        // Copy
app.KeyBinding.Add("Ctrl+V", pasteHandler)       // Paste
```

@end

## Práticas recomendadas

1. **Use atalhos padrão**: siga as convenções da plataforma para operações comuns:
  ```go
  // Cross-platform save
  if runtime.GOOS == "darwin" {
      app.KeyBinding.Add("Cmd+S", saveHandler)
  } else {
      app.KeyBinding.Add("Ctrl+S", saveHandler)
  }
  ```


2. **Forneça feedback visual**: informe aos usuários quando os atalhos forem acionados:
  ```go
  app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
      saveDocument()
      // Show brief notification
      window.EmitEvent("notification:show", "Document saved")
  })
  ```


3. **Trate os conflitos**: tome cuidado para não substituir atalhos importantes do sistema:
  ```go
  // Avoid overriding system shortcuts like:
  // Ctrl+Alt+Del (Windows)
  // Cmd+Space (macOS Spotlight)
  // Alt+Tab (Window switching)
  ```


4. **Documente os atalhos**: forneça ajuda ou documentação sobre os atalhos disponíveis:
  ```go
  app.KeyBinding.Add("F1", func(window application.Window) {
      // Show help dialog with available shortcuts
      showKeyboardShortcutsHelp()
  })
  ```


5. **Faça a limpeza**: remova associações de teclas temporárias quando elas não forem mais necessárias:
  ```go
  func enterEditMode() {
      app.KeyBinding.Add("Escape", exitEditModeHandler)
  }

  func exitEditModeHandler(window application.Window) {
      exitEditMode()
      app.KeyBinding.Remove("Escape") // Clean up temporary binding
  }
  ```


## Exemplo completo

Veja um exemplo completo de um editor de texto com atalhos de teclado:

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Text Editor with Shortcuts",
    })

    // File operations
    if runtime.GOOS == "darwin" {
        app.KeyBinding.Add("Cmd+N", func(window application.Window) {
            window.EmitEvent("file:new", nil)
        })
        app.KeyBinding.Add("Cmd+O", func(window application.Window) {
            openFile(app, window)
        })
        app.KeyBinding.Add("Cmd+S", func(window application.Window) {
            window.EmitEvent("file:save", nil)
        })
    } else {
        app.KeyBinding.Add("Ctrl+N", func(window application.Window) {
            window.EmitEvent("file:new", nil)
        })
        app.KeyBinding.Add("Ctrl+O", func(window application.Window) {
            openFile(app, window)
        })
        app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
            window.EmitEvent("file:save", nil)
        })
    }

    // View operations
    app.KeyBinding.Add("F11", func(window application.Window) {
        window.ToggleFullscreen()
    })

    app.KeyBinding.Add("F1", func(window application.Window) {
        showKeyboardShortcuts(window)
    })

    // Create main window
    window := app.Window.New()
    window.SetTitle("Text Editor")

    err := app.Run()
    if err != nil {
        panic(err)
    }
}

func openFile(app *application.App, window application.Window) {
    dialog := app.Dialog.OpenFile()
    dialog.AddFilter("Text Files", "*.txt;*.md")
    
    if file, err := dialog.PromptForSingleSelection(); err == nil {
        window.EmitEvent("file:open", file)
    }
}

func showKeyboardShortcuts(window application.Window) {
    shortcuts := `
Keyboard Shortcuts:
- Ctrl/Cmd+N: New file
- Ctrl/Cmd+O: Open file
- Ctrl/Cmd+S: Save file
- F11: Toggle fullscreen
- F1: Show this help
`
    window.EmitEvent("help:show", shortcuts)
}
```

@note{type="tip" title="Dica profissional"}
Teste as associações de teclas em todas as plataformas de destino para garantir que funcionem corretamente e não entrem em conflito com atalhos do sistema.

@end

@note{type="danger" title="Aviso"}
Tenha cuidado para não substituir atalhos essenciais do sistema. Algumas combinações de teclas são reservadas pelo sistema operacional e não podem ser capturadas pelos aplicativos.

@end
