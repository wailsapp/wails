# Liquid Glass Effect Example

This example demonstrates the Liquid Glass backdrop effect for macOS windows in Wails v3.

## Features

The Liquid Glass effect provides:
- **Dynamic glass material** that reflects and refracts light
- **Adaptive appearance** that responds to the content behind it
- **Liquid merge effects** when windows are grouped and placed close together
- **Customizable styles** including Light, Dark, and Vibrant modes
- **Tint colors** for adding subtle color overlays
- **Corner radius** for rounded glass effects

## Running the Example

```bash
cd v3/examples/liquid-glass
go run .
```

## Window Configurations

The example creates three windows to showcase different configurations:

### Window 1: Simple Liquid Glass
Uses the simplest configuration with just the backdrop type set:
```go
Mac: application.MacWindow{
    Backdrop: application.MacBackdropLiquidGlass,
}
```

### Window 2: Advanced Configuration
Shows advanced options with custom style, corner radius, and tint:
```go
Mac: application.MacWindow{
    Backdrop: application.MacBackdropLiquidGlass,
    LiquidGlass: application.MacLiquidGlass{
        Style:        application.LiquidGlassStyleVibrant,
        CornerRadius: 16.0,
        TintColor:    &application.RGBA{0, 122, 255, 50},
        GroupID:      "main-group",
        GroupSpacing: 8.0,
    },
}
```

### Window 3: Dark Style
Demonstrates the dark glass style with magenta tint:
```go
Mac: application.MacWindow{
    Backdrop: application.MacBackdropLiquidGlass,
    LiquidGlass: application.MacLiquidGlass{
        Style:        application.LiquidGlassStyleDark,
        CornerRadius: 20.0,
        TintColor:    &application.RGBA{255, 0, 255, 30},
        GroupID:      "secondary-group",
    },
}
```

## Compatibility

- **macOS 15.0+**: Full Liquid Glass effect with enhanced NSVisualEffectView
- **macOS 10.10-14.x**: Automatic fallback to standard translucent effect
- **Other platforms**: Not applicable (macOS-only feature)

## CSS Considerations

For best results with Liquid Glass:
1. Use `background: transparent` on the body
2. Apply semi-transparent backgrounds to containers
3. Use `backdrop-filter` for additional blur effects
4. Ensure text contrast with shadows or appropriate colors

## Grouping Windows

Windows with the same `GroupID` will exhibit liquid merge effects when positioned close together. The `GroupSpacing` property controls how close windows need to be to trigger the merge effect.

## Performance

The Liquid Glass effect uses GPU acceleration. For better performance on battery-powered devices, you can:
- Set `ReduceMotion: true` to reduce visual effects
- Use `StaticMode: true` for windows with static content
- Limit the number of grouped windows