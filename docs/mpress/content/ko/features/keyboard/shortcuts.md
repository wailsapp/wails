---
title: "키보드 단축키"
description: "기능에 빠르게 접근할 수 있도록 전역 키보드 단축키 등록"
slug: "features/keyboard/shortcuts"
sourcePath: "features/keyboard/shortcuts.md"
---

Wails는 애플리케이션의 모든 창에서 작동하는 전역 키보드 단축키를 등록할 수 있는 강력한 키 바인딩 시스템을 제공합니다. 이를 통해 사용자는 메뉴를 탐색하지 않고도 기능에 빠르게 접근할 수 있습니다.

## 키 바인딩 관리자에 접근하기

애플리케이션 인스턴스의 `KeyBindings` 속성을 통해 키 바인딩 관리자에 접근할 수 있습니다:

```go
app := application.New(application.Options{
    Name: "Keyboard Shortcuts Demo",
})

// Access the key binding manager
keyBindings := app.KeyBinding
```

## 키 바인딩 추가하기

### 기본 키 바인딩

간단한 키보드 단축키를 등록합니다:

```go
app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
    // Handle save action
    app.Logger.Info("Save shortcut triggered")
    // Perform save operation...
})
```

### 여러 키 바인딩

자주 사용하는 작업에 여러 단축키를 등록합니다:

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

## 키 바인딩 액셀러레이터

### 액셀러레이터 형식

키 바인딩은 보조 키와 일반 키로 구성된 표준 액셀러레이터 형식을 사용합니다:

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

### 플랫폼별 액셀러레이터

일반적인 단축키의 플랫폼별 차이를 처리합니다:

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

## 키 바인딩 관리하기

### 키 바인딩 제거하기

더 이상 필요하지 않은 키 바인딩을 제거합니다:

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

### 모든 키 바인딩 가져오기

등록된 모든 키 바인딩을 가져옵니다:

```go
allBindings := app.KeyBinding.GetAll()
for _, binding := range allBindings {
    app.Logger.Info("Key binding", "accelerator", binding.Accelerator)
}
```

## 고급 사용법

### 컨텍스트 인식 키 바인딩

애플리케이션 상태를 확인하여 키 바인딩이 컨텍스트를 인식하도록 만듭니다:

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

### 창별 동작

키 바인딩은 활성 창을 전달받으므로 창별 동작을 구현할 수 있습니다:

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

### 동적 키 바인딩 관리

애플리케이션 상태에 따라 키 바인딩을 동적으로 추가하고 제거합니다:

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

## 플랫폼 고려 사항

@tabs
[macOS]
macOS에서는 다음 사항을 고려하세요:

- 표준 단축키에는 `Ctrl` 대신 `Cmd`을 사용합니다
- `Cmd+Q`은 일반적으로 애플리케이션 종료용으로 예약되어 있습니다
- `Cmd+H`은 애플리케이션을 숨깁니다
- `Cmd+M`은 창을 최소화합니다
- 표준 macOS 키보드 단축키를 고려하세요

일반적인 macOS 패턴:

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
Windows에서는 다음 사항을 고려하세요:

- 표준 단축키에는 `Ctrl`을 사용합니다
- `Alt+F4`은 애플리케이션을 닫습니다
- `F1`은 일반적으로 도움말을 엽니다
- Windows 키보드 사용 규칙을 고려하세요

일반적인 Windows 패턴:

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
Linux에서는 다음 사항을 고려하세요:

- 일반적으로 `Ctrl`을 사용하는 Windows 규칙을 따릅니다
- 데스크톱 환경에 따라 다를 수 있습니다
- GNOME/KDE 표준 단축키를 고려하세요
- 일부 데스크톱 환경에서는 특정 단축키가 예약되어 있습니다

일반적인 Linux 패턴:

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

## 모범 사례

1. **표준 단축키 사용**: 일반적인 작업에는 플랫폼 규칙을 따르세요:
  ```go
  // Cross-platform save
  if runtime.GOOS == "darwin" {
      app.KeyBinding.Add("Cmd+S", saveHandler)
  } else {
      app.KeyBinding.Add("Ctrl+S", saveHandler)
  }
  ```


2. **시각적 피드백 제공**: 단축키가 실행되었음을 사용자가 알 수 있게 하세요:
  ```go
  app.KeyBinding.Add("Ctrl+S", func(window application.Window) {
      saveDocument()
      // Show brief notification
      window.EmitEvent("notification:show", "Document saved")
  })
  ```


3. **충돌 처리**: 중요한 시스템 단축키를 재정의하지 않도록 주의하세요:
  ```go
  // Avoid overriding system shortcuts like:
  // Ctrl+Alt+Del (Windows)
  // Cmd+Space (macOS Spotlight)
  // Alt+Tab (Window switching)
  ```


4. **단축키 문서화**: 사용할 수 있는 단축키에 관한 도움말이나 문서를 제공하세요:
  ```go
  app.KeyBinding.Add("F1", func(window application.Window) {
      // Show help dialog with available shortcuts
      showKeyboardShortcutsHelp()
  })
  ```


5. **정리**: 더 이상 필요하지 않은 임시 키 바인딩을 제거하세요:
  ```go
  func enterEditMode() {
      app.KeyBinding.Add("Escape", exitEditModeHandler)
  }

  func exitEditModeHandler(window application.Window) {
      exitEditMode()
      app.KeyBinding.Remove("Escape") // Clean up temporary binding
  }
  ```


## 전체 예제

다음은 키보드 단축키를 사용하는 텍스트 편집기의 전체 예제입니다:

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

@note{type="tip" title="전문가 팁"}
키 바인딩이 올바르게 작동하고 시스템 단축키와 충돌하지 않는지 모든 대상 플랫폼에서 테스트하세요.

@end

@note{type="danger" title="경고"}
중요한 시스템 단축키를 재정의하지 않도록 주의하세요. 일부 키 조합은 운영 체제에서 예약되어 있으므로 애플리케이션에서 감지할 수 없습니다.

@end
