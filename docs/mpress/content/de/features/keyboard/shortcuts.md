---
title: "Tastenkombinationen"
description: "Globale Tastenkombinationen für den schnellen Zugriff auf Funktionen registrieren"
slug: "features/keyboard/shortcuts"
sourcePath: "features/keyboard/shortcuts.md"
---

Wails bietet ein leistungsfähiges System für Tastenkürzel, mit dem Sie globale Tastenkombinationen registrieren können, die in allen Fenstern Ihrer Anwendung funktionieren. So können Benutzer schnell auf Funktionen zugreifen, ohne durch Menüs navigieren zu müssen.

## Auf den Tastenkürzel-Manager zugreifen

Der Zugriff auf den Tastenkürzel-Manager erfolgt über die Eigenschaft `KeyBindings` Ihrer Anwendungsinstanz:

```go
app := application.New(application.Options{
    Name: "Keyboard Shortcuts Demo",
})

// Access the key binding manager
keyBindings := app.KeyBinding
```

## Tastenkürzel hinzufügen

### Einfaches Tastenkürzel

Registrieren Sie eine einfache Tastenkombination:

```go
app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
    // Handle save action
    app.Logger.Info("Save shortcut triggered")
    // Perform save operation...
})
```

### Mehrere Tastenkürzel

Registrieren Sie mehrere Tastenkombinationen für häufige Vorgänge:

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

## Accelerators für Tastenkürzel

### Accelerator-Format

Tastenkürzel verwenden ein standardisiertes Accelerator-Format mit Zusatztasten und Tasten:

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

### Plattformspezifische Accelerators

Berücksichtigen Sie Plattformunterschiede bei häufig verwendeten Tastenkombinationen:

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

## Tastenkürzel verwalten

### Tastenkürzel entfernen

Entfernen Sie Tastenkürzel, wenn sie nicht mehr benötigt werden:

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

### Alle Tastenkürzel abrufen

Rufen Sie alle registrierten Tastenkürzel ab:

```go
allBindings := app.KeyBinding.GetAll()
for _, binding := range allBindings {
    app.Logger.Info("Key binding", "accelerator", binding.Accelerator)
}
```

## Erweiterte Verwendung

### Kontextabhängige Tastenkürzel

Machen Sie Tastenkürzel kontextabhängig, indem Sie den Anwendungsstatus prüfen:

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

### Fensterspezifische Aktionen

Tastenkürzel erhalten das aktive Fenster und ermöglichen dadurch fensterspezifisches Verhalten:

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

### Dynamische Verwaltung von Tastenkürzeln

Fügen Sie abhängig vom Anwendungsstatus dynamisch Tastenkürzel hinzu oder entfernen Sie sie:

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

## Plattformspezifische Aspekte

@tabs
[macOS]
Unter macOS:

- Verwenden Sie für standardmäßige Tastenkombinationen `Cmd` anstelle von `Ctrl`
- `Cmd+Q` ist üblicherweise zum Beenden der Anwendung vorgesehen
- `Cmd+H` blendet die Anwendung aus
- `Cmd+M` minimiert Fenster
- Berücksichtigen Sie die üblichen macOS-Tastenkombinationen

Gängige macOS-Muster:

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
Unter Windows:

- Verwenden Sie `Ctrl` für standardmäßige Tastenkombinationen
- `Alt+F4` schließt Anwendungen
- `F1` öffnet üblicherweise die Hilfe
- Berücksichtigen Sie die Windows-Konventionen für Tastenkombinationen

Gängige Windows-Muster:

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
Unter Linux:

- Folgt im Allgemeinen den Windows-Konventionen mit `Ctrl`
- Kann je nach Desktop-Umgebung variieren
- Berücksichtigen Sie die Standardtastenkombinationen von GNOME/KDE
- Einige Desktop-Umgebungen reservieren bestimmte Tastenkombinationen

Gängige Linux-Muster:

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

## Bewährte Vorgehensweisen

1. **Standardtastenkombinationen verwenden**: Halten Sie sich bei häufigen Vorgängen an die Konventionen der jeweiligen Plattform:
  ```go
  // Cross-platform save
  if runtime.GOOS == "darwin" {
      app.KeyBinding.Add("Cmd+S", saveHandler)
  } else {
      app.KeyBinding.Add("Ctrl+S", saveHandler)
  }
  ```


2. **Visuelles Feedback geben**: Informieren Sie Benutzer, wenn Tastenkombinationen ausgelöst werden:
  ```go
  app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
      saveDocument()
      // Show brief notification
      window.EmitEvent("notification:show", "Document saved")
  })
  ```


3. **Konflikte behandeln**: Achten Sie darauf, keine wichtigen Systemtastenkombinationen zu überschreiben:
  ```go
  // Avoid overriding system shortcuts like:
  // Ctrl+Alt+Del (Windows)
  // Cmd+Space (macOS Spotlight)
  // Alt+Tab (Window switching)
  ```


4. **Tastenkombinationen dokumentieren**: Stellen Sie Hilfe oder Dokumentation zu den verfügbaren Tastenkombinationen bereit:
  ```go
  app.KeyBinding.Add("F1", func(window application.Window) {
      // Show help dialog with available shortcuts
      showKeyboardShortcutsHelp()
  })
  ```


5. **Bereinigen**: Entfernen Sie temporäre Tastenkürzel, wenn sie nicht mehr benötigt werden:
  ```go
  func enterEditMode() {
      app.KeyBinding.Add("Escape", exitEditModeHandler)
  }

  func exitEditModeHandler(window application.Window) {
      exitEditMode()
      app.KeyBinding.Remove("Escape") // Clean up temporary binding
  }
  ```


## Vollständiges Beispiel

Hier ist ein vollständiges Beispiel für einen Texteditor mit Tastenkombinationen:

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

@note{type="tip" title="Profi-Tipp"}
Testen Sie Ihre Tastenkürzel auf allen Zielplattformen, um sicherzustellen, dass sie ordnungsgemäß funktionieren und nicht mit Systemtastenkombinationen in Konflikt stehen.

@end

@note{type="danger" title="Warnung"}
Achten Sie darauf, keine kritischen Systemtastenkombinationen zu überschreiben. Einige Tastenkombinationen sind dem Betriebssystem vorbehalten und können von Anwendungen nicht erfasst werden.

@end
