# Liquid Glass Demo for Wails v3

This demo showcases the native Liquid Glass effect available in macOS 26.0+ with fallback to NSVisualEffectView for older systems.

## Features Demonstrated

### Window Styles

1. **Light Glass** - Clean, light appearance with no tint
2. **Dark Glass** - Dark themed glass effect
3. **Vibrant Glass** - Enhanced vibrant effect for maximum transparency
4. **Tinted Glass** - Blue tinted glass with custom RGBA color
5. **Sheet Material** - Using specific NSVisualEffectMaterialSheet
6. **HUD Window** - Ultra-light HUD window material
7. **Content Background** - Content background material with warm tint

### Customization Options

- **Style**: `LiquidGlassStyleAutomatic`, `LiquidGlassStyleLight`, `LiquidGlassStyleDark`, `LiquidGlassStyleVibrant`
- **Material**: Direct NSVisualEffectMaterial selection (when NSGlassEffectView is not available)
  - `NSVisualEffectMaterialAppearanceBased`
  - `NSVisualEffectMaterialLight`
  - `NSVisualEffectMaterialDark`
  - `NSVisualEffectMaterialSheet`
  - `NSVisualEffectMaterialHUDWindow`
  - `NSVisualEffectMaterialContentBackground`
  - `NSVisualEffectMaterialUnderWindowBackground`
  - `NSVisualEffectMaterialUnderPageBackground`
  - And more...
- **CornerRadius**: Rounded corners (0 for square corners)
- **TintColor**: Custom RGBA tint overlay
- **GroupID**: Groups multiple glass windows in `privatemacapis` builds
- **GroupSpacing**: Spacing between grouped windows in `privatemacapis` builds

### Running the Demo

```bash
go build -tags privatemacapis -o liquid-glass-demo .
./liquid-glass-demo
```

The build tag is required for this full-window demo because macOS does not
provide a public API for making `WKWebView` transparent. Without the tag, Wails
still creates the public native glass view and preserves the same Go API, but
the opaque webview does not reveal the effect behind it.

### Requirements

- macOS 10.14+ (best experience on macOS 26.0+ with native NSGlassEffectView)
- Wails v3

### Implementation Details

The implementation uses:
- Native `NSGlassEffectView` on macOS 26.0+ for authentic glass effect
- Falls back to `NSVisualEffectView` on older systems
- Runtime detection using `NSClassFromString` for compatibility
- Public regular and clear glass styles in default builds
- An isolated `privatemacapis` build variant for WKWebView transparency and legacy grouping

### Example Usage

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropLiquidGlass,
        InvisibleTitleBarHeight: 500, // Make window draggable
        LiquidGlass: application.MacLiquidGlass{
            Style:        application.LiquidGlassStyleLight,
            Material:     application.NSVisualEffectMaterialHUDWindow,
            CornerRadius: 20.0,
            TintColor:    &application.RGBA{Red: 0, Green: 100, Blue: 200, Alpha: 50},
        },
    },
})
```
