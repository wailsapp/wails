---
title: "Сочетания клавиш"
description: "Регистрация глобальных сочетаний клавиш для быстрого доступа к функциям"
slug: "features/keyboard/shortcuts"
sourcePath: "features/keyboard/shortcuts.md"
---

Wails предоставляет мощную систему привязок клавиш, которая позволяет регистрировать глобальные сочетания клавиш, работающие во всех окнах приложения. Благодаря этому пользователи могут быстро обращаться к функциям, не переходя по меню.

## Доступ к диспетчеру привязок клавиш

Доступ к диспетчеру привязок клавиш осуществляется через свойство `KeyBindings` экземпляра приложения:

```go
app := application.New(application.Options{
    Name: "Keyboard Shortcuts Demo",
})

// Access the key binding manager
keyBindings := app.KeyBinding
```

## Добавление привязок клавиш

### Простая привязка клавиш

Зарегистрируйте простое сочетание клавиш:

```go
app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
    // Handle save action
    app.Logger.Info("Save shortcut triggered")
    // Perform save operation...
})
```

### Несколько привязок клавиш

Зарегистрируйте несколько сочетаний клавиш для распространённых операций:

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

## Акселераторы привязок клавиш

### Формат акселератора

Для привязок клавиш используется стандартный формат акселератора с клавишами-модификаторами и обычными клавишами:

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

### Акселераторы для разных платформ

Учитывайте различия между платформами для распространённых сочетаний клавиш:

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

## Управление привязками клавиш

### Удаление привязок клавиш

Удаляйте привязки клавиш, когда они больше не нужны:

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

### Получение всех привязок клавиш

Получите все зарегистрированные привязки клавиш:

```go
allBindings := app.KeyBinding.GetAll()
for _, binding := range allBindings {
    app.Logger.Info("Key binding", "accelerator", binding.Accelerator)
}
```

## Расширенное использование

### Контекстно-зависимые привязки клавиш

Сделайте привязки клавиш контекстно-зависимыми, проверяя состояние приложения:

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

### Действия для конкретного окна

Привязки клавиш получают активное окно, что позволяет реализовать поведение для конкретного окна:

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

### Динамическое управление привязками клавиш

Динамически добавляйте и удаляйте привязки клавиш в зависимости от состояния приложения:

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

## Особенности платформ

@tabs
[macOS]
В macOS:

- Для стандартных сочетаний клавиш используйте `Cmd` вместо `Ctrl`
- `Cmd+Q` обычно зарезервировано для завершения работы приложения
- `Cmd+H` скрывает приложение
- `Cmd+M` сворачивает окна
- Учитывайте стандартные сочетания клавиш macOS

Распространённые сочетания в macOS:

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
В Windows:

- Для стандартных сочетаний клавиш используйте `Ctrl`
- `Alt+F4` закрывает приложения
- `F1` обычно открывает справку
- Учитывайте принятые в Windows соглашения о сочетаниях клавиш

Распространённые сочетания в Windows:

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
В Linux:

- Обычно используются принятые в Windows соглашения с `Ctrl`
- Поведение может различаться в зависимости от среды рабочего стола
- Учитывайте стандартные сочетания клавиш GNOME/KDE
- Некоторые среды рабочего стола резервируют определённые сочетания клавиш

Распространённые сочетания в Linux:

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

## Рекомендации

1. **Используйте стандартные сочетания клавиш**: следуйте соглашениям платформы для распространённых операций:
  ```go
  // Cross-platform save
  if runtime.GOOS == "darwin" {
      app.KeyBinding.Add("Cmd+S", saveHandler)
  } else {
      app.KeyBinding.Add("Ctrl+S", saveHandler)
  }
  ```


2. **Обеспечьте визуальную обратную связь**: сообщайте пользователям о срабатывании сочетаний клавиш:
  ```go
  app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
      saveDocument()
      // Show brief notification
      window.EmitEvent("notification:show", "Document saved")
  })
  ```


3. **Обрабатывайте конфликты**: следите за тем, чтобы не переопределить важные системные сочетания клавиш:
  ```go
  // Avoid overriding system shortcuts like:
  // Ctrl+Alt+Del (Windows)
  // Cmd+Space (macOS Spotlight)
  // Alt+Tab (Window switching)
  ```


4. **Документируйте сочетания клавиш**: предоставьте справку или документацию по доступным сочетаниям клавиш:
  ```go
  app.KeyBinding.Add("F1", func(window application.Window) {
      // Show help dialog with available shortcuts
      showKeyboardShortcutsHelp()
  })
  ```


5. **Выполняйте очистку**: удаляйте временные привязки клавиш, когда они больше не нужны:
  ```go
  func enterEditMode() {
      app.KeyBinding.Add("Escape", exitEditModeHandler)
  }

  func exitEditModeHandler(window application.Window) {
      exitEditMode()
      app.KeyBinding.Remove("Escape") // Clean up temporary binding
  }
  ```


## Полный пример

Ниже приведён полный пример текстового редактора с сочетаниями клавиш:

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

@note{type="tip" title="Полезный совет"}
Проверяйте привязки клавиш на всех целевых платформах, чтобы убедиться, что они работают правильно и не конфликтуют с системными сочетаниями клавиш.

@end

@note{type="danger" title="Предупреждение"}
Следите за тем, чтобы не переопределить критически важные системные сочетания клавиш. Некоторые комбинации клавиш зарезервированы операционной системой и не могут быть перехвачены приложениями.

@end
