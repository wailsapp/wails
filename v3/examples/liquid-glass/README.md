# Liquid Glass Demo for Wails v3

This demo showcases the native Liquid Glass effect available in macOS 15.0+ with fallback to NSVisualEffectView for older systems.

## Running on macOS with private APIs

On macOS, this example uses private WebKit transparency and native Liquid Glass style values. Build with `private_mac_apis` to see the intended effects. Without it, the webview remains opaque and glass styles use public alternatives.

From this example directory, run:

```bash
go run -tags private_mac_apis .
```

Omit `-tags private_mac_apis` to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.

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
- **GroupID**: For grouping multiple glass windows (future feature)
- **GroupSpacing**: Spacing between grouped windows (future feature)

### Running the Demo

```bash
go build -tags private_mac_apis -o liquid-glass-demo .
./liquid-glass-demo
```

### Requirements

- macOS 10.14+ (best experience on macOS 26.0+ with native NSGlassEffectView)
- Wails v3

### Implementation Details

The implementation uses:
- Native `NSGlassEffectView` on macOS 26.0+ for authentic glass effect
- Falls back to `NSVisualEffectView` on older systems
- Runtime detection using `NSClassFromString` for compatibility
- Key-Value Coding (KVC) for dynamic property setting

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