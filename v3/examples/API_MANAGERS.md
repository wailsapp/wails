# Wails v3 Manager API Guide

This document demonstrates the new manager-based API structure introduced in Wails v3 for better organization and discoverability.

## Overview

The App API has been reorganized into focused manager structs while maintaining full backward compatibility:

- **app.Windows** - Window management
- **app.ContextMenus** - Context menu operations  
- **app.KeyBindings** - Key binding management
- **app.Browser** - Browser operations
- **app.Env** - Environment information
- **app.Dialogs** - Dialog operations
- **app.Events** - Event handling
- **app.Menus** - Application menu management
- **app.Screens** - Screen and display management

## API Comparison

### Window Management

```go
// Traditional API (still works)
window := app.GetWindowByName("main")
app.OnWindowCreation(func(w Window) { ... })
newWindow := app.NewWebviewWindow()

// NEW: Manager API (recommended with consistent return patterns)
window, exists := app.Windows.GetByName("main")
window, exists := app.Windows.GetByID(123)
current, exists := app.Windows.GetCurrent()
app.Windows.OnCreate(func(w Window) { ... })
newWindow := app.Windows.New()
```

### Event Handling

```go
// Traditional API (still works)
app.EmitEvent("custom", data)
app.OnEvent("custom", func(e *CustomEvent) { ... })
app.OffEvent("custom")
app.ResetEvents()

// NEW: Manager API (recommended)
app.Events.Emit("custom", data)
app.Events.On("custom", func(e *CustomEvent) { ... })
app.Events.Off("custom")
app.Events.Reset()
```

### Browser Operations

```go
// Traditional API (still works)
app.BrowserOpenURL("https://wails.io")
app.BrowserOpenFile("/path/to/file")

// NEW: Manager API (recommended)
app.Browser.OpenURL("https://wails.io")
app.Browser.OpenFile("/path/to/file")
```

### Menu Management

```go
// Traditional API (still works)
menu := app.NewMenu()
app.SetMenu(menu)
app.ShowAboutDialog()

// NEW: Manager API (recommended)
menu := app.Menus.New()
app.Menus.Set(menu)
app.Menus.ShowAbout()
```

### Environment Information

```go
// Traditional API (still works)
env := app.Environment()
darkMode := app.IsDarkMode()

// NEW: Manager API (recommended)
env := app.Env.Info()
darkMode := app.Env.IsDarkMode()
```

### Dialog Operations

```go
// Traditional API (global functions, still works)
result := application.OpenFileDialog().PromptForSingleSelection()
application.InfoDialog().SetMessage("Hello").Show()

// NEW: Manager API (clearer method names with Show prefix)
result := app.Dialogs.ShowOpenFileDialog().PromptForSingleSelection()
app.Dialogs.ShowInfoDialog().SetMessage("Hello").Show()
app.Dialogs.ShowWarningDialog().SetMessage("Warning!").Show()
app.Dialogs.ShowErrorDialog().SetMessage("Error occurred").Show()
```

### Key Bindings

```go
// Traditional API (private methods, accessed via options)
app := application.New(application.Options{
    KeyBindings: map[string]func(window *application.WebviewWindow){
        "ctrl+n": func(window *application.WebviewWindow) {
            // Handle key binding
        },
    },
})

// NEW: Manager API (recommended, clearer parameter naming)
app.KeyBindings.Add("ctrl+n", func(window *application.WebviewWindow) {
    // Handle key binding with clear 'accelerator' parameter
})
app.KeyBindings.Remove("ctrl+n")
bindings := app.KeyBindings.GetAll() // Returns []*KeyBinding slice
```

### Context Menus

```go
// Traditional API (still works)
contextMenu := application.NewContextMenu("myMenu")
// Context menus were already well-organized

// NEW: Manager API (consistent Add/Remove verbs)
app.ContextMenus.Add("myMenu", contextMenu)
menu, exists := app.ContextMenus.Get("myMenu")
app.ContextMenus.Remove("myMenu")
menus := app.ContextMenus.GetAll() // Returns []*ContextMenu slice
```

### Screen Management

```go
// Traditional API (still works)
screens, err := app.GetScreens()
screen, err := app.GetPrimaryScreen()

// NEW: Manager API (recommended, no error handling needed)
screens := app.Screens.GetAll()
primaryScreen := app.Screens.GetPrimary()

// Advanced screen operations (always available through manager)
dipPoint := application.Point{X: 100, Y: 100}
physicalPoint := app.Screens.DipToPhysicalPoint(dipPoint)
nearestScreen := app.Screens.ScreenNearestDipPoint(dipPoint)
```

## Migration Strategy

1. **Existing code continues to work** - no immediate changes required
2. **New projects should use manager APIs** - better organization and discoverability
3. **Gradual migration recommended** - update methods as you encounter them
4. **IDE support improved** - autocomplete shows organized API surface

## Benefits

- **Better discoverability** - related methods grouped together
- **Improved IDE experience** - easier to find relevant APIs
- **Cleaner code organization** - separation of concerns
- **Future extensibility** - easier to add new features to specific areas
- **Backward compatibility** - existing code continues to work unchanged

## Example

See `examples/events/main_with_managers.go` for a complete example showing both traditional and manager APIs side by side.