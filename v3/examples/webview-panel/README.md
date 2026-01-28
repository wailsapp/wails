# WebView Panel Example

This example demonstrates the **WebviewPanel** feature - embedding multiple independent webview panels within a single window. This is similar to Electron's BrowserView/WebContentsView and addresses [GitHub issue #1997](https://github.com/wailsapp/wails/issues/1997).

## Features Demonstrated

- **URL Loading**: Load external websites (like <https://wails.io>) in an embedded panel
- **Responsive Layout**: Panel automatically resizes with the window using anchor settings
- **Dynamic Switching**: Switch panel content between different URLs at runtime
- **Custom UI Integration**: Panel embedded within a custom HTML interface

## Running

```bash
cd v3/examples/webview-panel
go run .
```

## What This Example Shows

1. **Embedded Webview**: The main window displays a custom UI with a header and navigation buttons
2. **Panel Inside Window**: An embedded webview panel shows <https://wails.io> inside the window
3. **URL Switching**: Click the "Wails.io" or "Google.com" buttons to switch the panel content
4. **Responsive Behavior**: Resize the window to see the panel automatically adjust its size

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
// Create a panel with URL and positioning
panel := window.NewPanel(application.WebviewPanelOptions{
    Name:   "content",
    URL:    "https://example.com",
    X:      20,
    Y:      60,
    Width:  800,
    Height: 500,
})

// Panel with anchoring (responsive to window resize)
panel := window.NewPanel(application.WebviewPanelOptions{
    Name:   "sidebar",
    URL:    "/sidebar.html",
    X:      0,
    Y:      0,
    Width:  200,
    Height: 600,
    Anchor: application.AnchorTop | application.AnchorBottom | application.AnchorLeft,
})

// Panel that fills the entire window
panel := window.NewPanel(application.WebviewPanelOptions{
    Name:   "fullscreen",
    URL:    "https://wails.io",
    X:      0,
    Y:      0,
    Width:  800,
    Height: 600,
    Anchor: application.AnchorFill,
})
```

### Anchor Types

Anchors control how panels respond to window resizing:

| Anchor | Behavior |
|--------|----------|
| `AnchorNone` | Fixed position and size |
| `AnchorTop` | Maintains distance from top edge |
| `AnchorBottom` | Maintains distance from bottom edge |
| `AnchorLeft` | Maintains distance from left edge |
| `AnchorRight` | Maintains distance from right edge |
| `AnchorFill` | Anchored to all edges (fills window with margins) |

Combine anchors with `|` for complex layouts:
```go
// Left sidebar that stretches vertically
Anchor: application.AnchorTop | application.AnchorBottom | application.AnchorLeft
```

### Panel Options

```go
application.WebviewPanelOptions{
    // Identity
    Name: "panel-name",           // Unique identifier
    
    // Content
    URL:     "https://example.com", // URL to load
    Headers: map[string]string{     // Custom HTTP headers (optional)
        "Authorization": "Bearer token",
    },
    UserAgent: "Custom UA",         // Custom user agent (optional)
    
    // Position & Size
    X:      100,                    // X position (CSS pixels)
    Y:      50,                     // Y position (CSS pixels)  
    Width:  800,                    // Width (CSS pixels)
    Height: 600,                    // Height (CSS pixels)
    ZIndex: 1,                      // Stacking order
    Anchor: application.AnchorFill, // Resize behavior
    
    // Appearance
    Visible:         boolPtr(true),  // Initially visible
    BackgroundColour: application.NewRGB(255, 255, 255),
    Transparent:     false,          // Transparent background
    Zoom:            1.0,            // Zoom level (1.0 = 100%)
    
    // Developer
    DevToolsEnabled:        boolPtr(true),
    OpenInspectorOnStartup: false,
}
```

### Panel Manipulation

```go
// Position and size
panel.SetBounds(application.Rect{X: 100, Y: 50, Width: 400, Height: 300})
bounds := panel.Bounds()

// Content
panel.SetURL("https://wails.io")
panel.Reload()
panel.ForceReload() // Bypass cache

// Visibility
panel.Show()
panel.Hide()
visible := panel.IsVisible()

// Stacking order
panel.SetZIndex(10)

// Focus
panel.Focus()
focused := panel.IsFocused()

// Zoom
panel.SetZoom(1.5)
zoom := panel.GetZoom()

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

## Key Differences from Windows

| Feature | WebviewWindow | WebviewPanel |
|---------|---------------|--------------|
| Has title bar | ✅ | ❌ |
| Can be minimized/maximized | ✅ | ❌ |
| Independent window | ✅ | ❌ (child of window) |
| Can show external URLs | ✅ | ✅ |
| Multiple per app | ✅ | ✅ (multiple per window) |
| Position relative to | Screen | Parent window |
| Responsive anchoring | ❌ | ✅ |
