# WebView Panel Example

This example demonstrates the **WebviewPanel** feature - embedding multiple independent webview panels within a single window. This is similar to Electron's BrowserView/WebContentsView and addresses [GitHub issue #1997](https://github.com/wailsapp/wails/issues/1997).

## Features

- Create multiple webview panels within a single window
- Panels are absolutely positioned with X, Y, Width, Height
- Each panel can load different URLs or HTML content
- Independent JavaScript execution in each panel
- Z-index support for panel stacking
- Layout helper methods for common patterns (DockLeft, DockRight, etc.)

## Running

```bash
cd v3/examples/webview-panel
go run main.go
```

## Use Cases

WebviewPanel is ideal for:

- **IDE-like layouts**: Editor + preview + terminal panels
- **Browser-style apps**: Tab bar + content area  
- **Dashboard apps**: Navigation sidebar + main content
- **Email clients**: Folder list + message list + preview pane
- **News readers**: Article list + external website viewer
- **Dev tools**: App preview + inspector panels

## API Overview

### Creating Panels

```go
// Create a panel with explicit positioning
panel := window.NewPanel(application.WebviewPanelOptions{
    Name:   "sidebar",
    X:      0,
    Y:      50,
    Width:  200,
    Height: 600,
    URL:    "https://example.com",
    // Or use HTML:
    // HTML: "<h1>Hello Panel!</h1>",
})
```

### Layout Helpers

```go
// Dock to edges
sidebar := window.NewPanel(opts).DockLeft(200)    // Left sidebar
inspector := window.NewPanel(opts).DockRight(300) // Right panel
toolbar := window.NewPanel(opts).DockTop(50)      // Top toolbar
statusBar := window.NewPanel(opts).DockBottom(30) // Bottom status

// Fill remaining space
content := window.NewPanel(opts).FillBeside(sidebar, "right")

// Fill entire window
fullPanel := window.NewPanel(opts).FillWindow()
```

### Panel Manipulation

```go
// Position and size
panel.SetBounds(application.Rect{X: 100, Y: 50, Width: 400, Height: 300})
panel.SetPosition(200, 100)
panel.SetSize(500, 400)

// Content
panel.SetURL("https://wails.io")
panel.SetHTML("<h1>Dynamic content</h1>")
panel.ExecJS("console.log('Hello from panel!')")
panel.Reload()

// Visibility
panel.Show()
panel.Hide()
visible := panel.IsVisible()

// Stacking order
panel.SetZIndex(10)

// Focus
panel.Focus()

// Zoom
panel.SetZoom(1.5)

// Developer tools
panel.OpenDevTools()

// Cleanup
panel.Destroy()
```

### Getting Panels

```go
// Get panel by name
panel := window.GetPanel("sidebar")

// Get panel by ID
panel := window.GetPanelByID(1)

// Get all panels
panels := window.GetPanels()

// Remove panel
window.RemovePanel("sidebar")
```

## TypeScript/Frontend API

```typescript
import { Panel } from '@wailsio/runtime';

// Get a reference to a panel
const panel = Panel.Get("content");

// Manipulate from frontend
await panel.SetBounds({ x: 100, y: 50, width: 500, height: 400 });
await panel.SetURL("https://wails.io");
await panel.ExecJS("document.body.style.background = 'red'");
await panel.Show();
await panel.Hide();
await panel.Focus();
```

## Key Differences from Windows

| Feature | WebviewWindow | WebviewPanel |
|---------|---------------|--------------|
| Has title bar | ✅ | ❌ |
| Can be minimized/maximized | ✅ | ❌ |
| Independent window | ✅ | ❌ (child of window) |
| Can show external URLs | ✅ | ✅ |
| Multiple per app | ✅ | ✅ (multiple per window) |
| Position relative to | Screen | Parent window |
