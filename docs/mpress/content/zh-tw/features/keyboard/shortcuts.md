---
title: "鍵盤快速鍵"
description: "註冊全域鍵盤快速鍵，以便快速存取功能"
slug: "features/keyboard/shortcuts"
sourcePath: "features/keyboard/shortcuts.md"
---

Wails 提供功能強大的按鍵綁定系統，可讓你註冊在應用程式所有視窗中皆可使用的全域鍵盤快速鍵。使用者無須瀏覽選單，即可快速存取功能。

## 存取按鍵綁定管理器

透過應用程式執行個體的 `KeyBindings` 屬性存取按鍵綁定管理器：

```go
app := application.New(application.Options{
    Name: "Keyboard Shortcuts Demo",
})

// Access the key binding manager
keyBindings := app.KeyBinding
```

## 新增按鍵綁定

### 基本按鍵綁定

註冊簡單的鍵盤快速鍵：

```go
app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
    // Handle save action
    app.Logger.Info("Save shortcut triggered")
    // Perform save operation...
})
```

### 多個按鍵綁定

為常用操作註冊多個快速鍵：

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

## 按鍵綁定加速鍵

### 加速鍵格式

按鍵綁定使用由輔助鍵和按鍵組成的標準加速鍵格式：

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

### 平台特定的加速鍵

處理常用快速鍵在不同平台上的差異：

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

## 管理按鍵綁定

### 移除按鍵綁定

不再需要按鍵綁定時將其移除：

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

### 取得所有按鍵綁定

擷取所有已註冊的按鍵綁定：

```go
allBindings := app.KeyBinding.GetAll()
for _, binding := range allBindings {
    app.Logger.Info("Key binding", "accelerator", binding.Accelerator)
}
```

## 進階用法

### 可感知情境的按鍵綁定

檢查應用程式狀態，使按鍵綁定能感知情境：

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

### 視窗特定操作

按鍵綁定會接收作用中的視窗，因此可實作視窗特定的行為：

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

### 動態管理按鍵綁定

根據應用程式狀態動態新增及移除按鍵綁定：

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

## 平台注意事項

@tabs
[macOS]
在 macOS 上：

- 標準快速鍵請使用 `Cmd`，而非 `Ctrl`
- `Cmd+Q` 通常保留用於結束應用程式
- `Cmd+H` 會隱藏應用程式
- `Cmd+M` 會將視窗最小化
- 請考量標準的 macOS 鍵盤快速鍵

常見的 macOS 模式：

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
在 Windows 上：

- 標準快速鍵請使用 `Ctrl`
- `Alt+F4` 會關閉應用程式
- `F1` 通常會開啟說明
- 請考量 Windows 鍵盤操作慣例

常見的 Windows 模式：

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
在 Linux 上：

- 通常會遵循使用 `Ctrl` 的 Windows 慣例
- 可能因桌面環境而異
- 請考量 GNOME/KDE 的標準快速鍵
- 某些桌面環境會保留特定快速鍵

常見的 Linux 模式：

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

## 最佳實務

1. **使用標準快速鍵**：常用操作請遵循平台慣例：
  ```go
  // Cross-platform save
  if runtime.GOOS == "darwin" {
      app.KeyBinding.Add("Cmd+S", saveHandler)
  } else {
      app.KeyBinding.Add("Ctrl+S", saveHandler)
  }
  ```


2. **提供視覺回饋**：觸發快速鍵時通知使用者：
  ```go
  app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
      saveDocument()
      // Show brief notification
      window.EmitEvent("notification:show", "Document saved")
  })
  ```


3. **處理衝突**：請小心，不要覆寫重要的系統快速鍵：
  ```go
  // Avoid overriding system shortcuts like:
  // Ctrl+Alt+Del (Windows)
  // Cmd+Space (macOS Spotlight)
  // Alt+Tab (Window switching)
  ```


4. **記錄快速鍵**：為可用的快速鍵提供說明或文件：
  ```go
  app.KeyBinding.Add("F1", func(window application.Window) {
      // Show help dialog with available shortcuts
      showKeyboardShortcutsHelp()
  })
  ```


5. **清理**：不再需要暫時按鍵綁定時將其移除：
  ```go
  func enterEditMode() {
      app.KeyBinding.Add("Escape", exitEditModeHandler)
  }

  func exitEditModeHandler(window application.Window) {
      exitEditMode()
      app.KeyBinding.Remove("Escape") // Clean up temporary binding
  }
  ```


## 完整範例

以下是具有鍵盤快速鍵之文字編輯器的完整範例：

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

@note{type="tip" title="專業提示"}
請在所有目標平台上測試按鍵綁定，以確保其正常運作，且不會與系統快速鍵發生衝突。

@end

@note{type="danger" title="警告"}
請小心，不要覆寫關鍵的系統快速鍵。某些按鍵組合由作業系統保留，應用程式無法擷取。

@end
