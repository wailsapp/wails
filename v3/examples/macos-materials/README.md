# macOS Materials Support

This directory demonstrates the comprehensive macOS materials support added to Wails v3, which provides access to all native Apple visual effect materials as documented in [Apple's Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/materials).

## Overview

macOS materials are translucent visual effects that allow content beneath windows to show through, creating depth and visual hierarchy. This implementation provides access to all 14 official Apple materials with full configuration options.

## Available Materials

### Basic Materials
- `MacBackdropNormal` - Default opaque background
- `MacBackdropTransparent` - Fully transparent background  
- `MacBackdropTranslucent` - Basic translucent effect (legacy)

### Apple-Defined Materials
- `MacBackdropMaterial` - Default system material
- `MacBackdropSidebar` - Designed for sidebar backgrounds
- `MacBackdropMenu` - Optimized for menu backgrounds
- `MacBackdropPopover` - Perfect for popover windows
- `MacBackdropTitlebar` - Titlebar material (deprecated in macOS 10.14+)
- `MacBackdropHeaderView` - Header view backgrounds
- `MacBackdropSheet` - Sheet dialog backgrounds
- `MacBackdropWindowBackground` - Main window backgrounds
- `MacBackdropUnderWindowBackground` - Background layers
- `MacBackdropContentBackground` - Content area backgrounds
- `MacBackdropUnderPageBackground` - Page background layers
- `MacBackdropTooltip` - Tooltip backgrounds
- `MacBackdropFullScreenUI` - Full-screen interface backgrounds
- `MacBackdropHUDWindow` - HUD (heads-up display) backgrounds

## Configuration Options

### Material Blending Modes
- `MacMaterialBlendingModeBehindWindow` - Blends with content behind the window
- `MacMaterialBlendingModeWithinWindow` - Blends with content within the window

### Material States  
- `MacMaterialStateFollowsWindowActiveState` - Material follows window's active state
- `MacMaterialStateActive` - Material is always active
- `MacMaterialStateInactive` - Material is always inactive

### Enhanced Options
- `EmphasizedAppearance` - Enables enhanced appearance (macOS 10.14+)

## Usage Example

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Materials Demo",
    })

    window := app.NewWebviewWindow().SetOptions(application.WebviewWindowOptions{
        Title:  "Sidebar Material Window",
        Width:  400,
        Height: 300,
        Mac: application.MacWindow{
            Backdrop: application.MacBackdropSidebar,
            MaterialOptions: application.MacMaterialOptions{
                BlendingMode:         application.MacMaterialBlendingModeBehindWindow,
                State:               application.MacMaterialStateFollowsWindowActiveState,
                EmphasizedAppearance: true,
            },
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
                FullSizeContent:    true,
            },
        },
        HTML: `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            background: rgba(255, 255, 255, 0.1);
            font-family: -apple-system, BlinkMacSystemFont, sans-serif;
            color: #333;
        }
        .content {
            background: rgba(255, 255, 255, 0.7);
            padding: 20px;
            border-radius: 10px;
            backdrop-filter: blur(10px);
        }
    </style>
</head>
<body>
    <div class="content">
        <h1>Sidebar Material</h1>
        <p>This demonstrates native macOS Sidebar material.</p>
    </div>
</body>
</html>
        `,
    })

    window.Show()
    app.Run()
}
```

## Implementation Details

### C Implementation
The implementation uses `NSVisualEffectView` with proper macOS version compatibility:

- **macOS 10.10+**: Basic materials (Sidebar, Menu, Popover, Titlebar)
- **macOS 10.14+**: Enhanced materials and emphasized appearance
- **Runtime checks**: Ensures compatibility across macOS versions

### CSS Integration
For best results, combine materials with CSS:

```css
body {
    background: rgba(255, 255, 255, 0.1); /* Semi-transparent background */
}

.glass-panel {
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: blur(10px);
    border-radius: 10px;
}
```

## Design Guidelines

### When to Use Materials

1. **Sidebar Material** - Navigation panels, file browsers
2. **Menu Material** - Context menus, dropdown menus  
3. **Popover Material** - Floating panels, tooltips
4. **Header View Material** - Toolbar backgrounds
5. **Sheet Material** - Modal dialogs, settings panels
6. **Window Background** - Main content areas
7. **HUD Window** - Overlay interfaces, controls

### Best Practices

1. **Layer appropriately** - Use correct material for UI hierarchy
2. **Consider contrast** - Ensure text remains readable
3. **Test across themes** - Verify appearance in light/dark modes
4. **Performance** - Materials can impact rendering performance
5. **Accessibility** - Maintain sufficient contrast ratios

## Version Compatibility

- **macOS 10.10+**: Basic material support
- **macOS 10.14+**: Full material library and emphasized appearance
- **macOS 11.0+**: Optimized performance and additional refinements

## Troubleshooting

### Material Not Appearing
1. Ensure `webviewSetTransparent()` is called
2. Verify CSS uses semi-transparent backgrounds
3. Check macOS version compatibility

### Poor Performance
1. Limit number of material views
2. Use appropriate blending modes
3. Consider simpler alternatives for complex layouts

### Accessibility Issues
1. Test with VoiceOver enabled
2. Verify contrast ratios meet guidelines
3. Provide high-contrast alternatives when needed

## Related Documentation

- [Apple's Materials Documentation](https://developer.apple.com/design/human-interface-guidelines/materials)
- [NSVisualEffectView Reference](https://developer.apple.com/documentation/appkit/nsvisualeffectview)
- [Wails v3 Window Options](../../../docs/window-options.md) 