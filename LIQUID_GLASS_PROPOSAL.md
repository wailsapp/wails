# Liquid Glass Implementation Proposal for Wails v3

## Executive Summary

This proposal outlines the implementation of Apple's Liquid Glass effect in Wails v3, providing a modern, translucent glass material that reflects and refracts light to create depth and dynamism in the user interface. The implementation will provide a simple, DX-friendly API while maintaining cross-platform compatibility where possible.

## Background

Apple introduced the Liquid Glass material as part of their new design language in 2025. It represents a significant evolution from the existing `NSVisualEffectView`, offering:
- Dynamic glass effects that adapt to surroundings
- Liquid visual effects when multiple glass elements merge
- Enhanced legibility through automatic content treatment
- Sophisticated light refraction and reflection

## Current State Analysis

### Existing Wails v3 Architecture

Wails v3 currently supports:
- **macOS**: `MacBackdrop` enum with `Normal`, `Transparent`, and `Translucent` options
- **Windows**: `BackdropType` with `Auto`, `None`, `Mica`, `Acrylic`, and `Tabbed` options
- **Linux**: Basic translucency through `WindowIsTranslucent` flag

The current implementation uses:
- macOS: `NSVisualEffectView` for translucent effects
- Windows: DWM APIs for Mica/Acrylic effects (Windows 11) or fallback blur
- Linux: GTK transparency

## Technical Requirements

### macOS Requirements
- **Minimum OS**: macOS Tahoe (26.0) for full Liquid Glass support
- **SDK**: Xcode 26 with updated AppKit headers
- **APIs**: `NSGlassEffectView` and `NSGlassEffectContainerView`

### Key API Components

```objc
// NSGlassEffectView - Primary glass effect view
@interface NSGlassEffectView : NSView
@property (nonatomic, strong) NSView *contentView;
@property (nonatomic) CGFloat cornerRadius;
@property (nonatomic, strong) NSColor *tintColor;
@property (nonatomic) NSGlassEffectStyle style;
@end

// NSGlassEffectContainerView - Groups multiple glass elements
@interface NSGlassEffectContainerView : NSView
@property (nonatomic) CGFloat spacing;
@end

// Style enum (inferred from documentation)
typedef NS_ENUM(NSInteger, NSGlassEffectStyle) {
    NSGlassEffectStyleAutomatic = 0,
    NSGlassEffectStyleLight = 1,
    NSGlassEffectStyleDark = 2,
    NSGlassEffectStyleVibrant = 3
};
```

## Proposed Implementation

### 1. API Design

#### Window Options Enhancement

```go
// pkg/application/webview_window_options.go

// Add to MacBackdrop enum
const (
    MacBackdropNormal      MacBackdrop = iota
    MacBackdropTransparent
    MacBackdropTranslucent
    MacBackdropLiquidGlass // New option
)

// Add new LiquidGlass configuration
type MacLiquidGlass struct {
    // Enable liquid glass effect
    Enabled bool
    
    // Style of the glass effect
    Style MacLiquidGlassStyle
    
    // Corner radius for the glass effect
    CornerRadius float64
    
    // Tint color for the glass (optional)
    TintColor *RGBA
    
    // Group with other windows for liquid merging effect
    GroupID string
    
    // Spacing for grouped glass elements (in points)
    GroupSpacing float64
}

type MacLiquidGlassStyle int

const (
    LiquidGlassStyleAutomatic MacLiquidGlassStyle = iota
    LiquidGlassStyleLight
    LiquidGlassStyleDark
    LiquidGlassStyleVibrant
)

// Update MacWindow struct
type MacWindow struct {
    // ... existing fields ...
    
    // Liquid Glass configuration
    LiquidGlass MacLiquidGlass
}
```

### 2. Developer Experience (DX)

#### Simple Toggle Option
For developers who want a quick implementation:

```go
app.NewWebviewWindowWithOptions(webview.WebviewWindowOptions{
    Title:  "My App",
    Width:  800,
    Height: 600,
    Mac: MacWindow{
        Backdrop: MacBackdropLiquidGlass, // Simple toggle
    },
})
```

#### Advanced Configuration
For fine-tuned control:

```go
app.NewWebviewWindowWithOptions(webview.WebviewWindowOptions{
    Title:  "My App",
    Width:  800,
    Height: 600,
    Mac: MacWindow{
        LiquidGlass: MacLiquidGlass{
            Enabled:      true,
            Style:        LiquidGlassStyleVibrant,
            CornerRadius: 12.0,
            TintColor:    &RGBA{0, 122, 255, 128}, // Semi-transparent blue
            GroupID:      "main-window-group",
            GroupSpacing: 8.0,
        },
    },
})
```

### 3. Implementation Details

#### macOS Implementation (`webview_window_darwin.go`)

```go
// Add to window creation
func (w *macosWebviewWindow) applyLiquidGlass(options MacLiquidGlass) {
    if !options.Enabled {
        return
    }
    
    // Check OS compatibility
    if !C.isLiquidGlassSupported() {
        // Fallback to NSVisualEffectView
        w.applyTranslucentBackdrop()
        w.app.debug("Liquid Glass not supported, falling back to translucent backdrop")
        return
    }
    
    // Apply liquid glass
    C.windowSetLiquidGlass(
        w.nsWindow,
        C.int(options.Style),
        C.double(options.CornerRadius),
        C.int(options.TintColor.Red),
        C.int(options.TintColor.Green),
        C.int(options.TintColor.Blue),
        C.int(options.TintColor.Alpha),
        C.CString(options.GroupID),
        C.double(options.GroupSpacing),
    )
}
```

#### Objective-C Implementation (`webview_window_darwin.m`)

```objc
// Check for Liquid Glass support
bool isLiquidGlassSupported() {
    if (@available(macOS 26.0, *)) {
        return [NSGlassEffectView class] != nil;
    }
    return false;
}

// Apply Liquid Glass effect
void windowSetLiquidGlass(void* nsWindow, int style, double cornerRadius, 
                          int r, int g, int b, int a, 
                          const char* groupID, double groupSpacing) {
    WebviewWindow* window = (WebviewWindow*)nsWindow;
    
    // Remove existing effect views
    [window removeVisualEffects];
    
    // Create glass effect view
    NSGlassEffectView* glassView = [[NSGlassEffectView alloc] init];
    glassView.style = (NSGlassEffectStyle)style;
    glassView.cornerRadius = cornerRadius;
    
    if (a > 0) {
        glassView.tintColor = [NSColor colorWithRed:r/255.0 
                                              green:g/255.0 
                                               blue:b/255.0 
                                              alpha:a/255.0];
    }
    
    // Set the webview as content
    glassView.contentView = window.webView;
    
    // Handle grouping if specified
    if (groupID && strlen(groupID) > 0) {
        NSGlassEffectContainerView* container = [windowGroupManager getOrCreateContainer:groupID];
        container.spacing = groupSpacing;
        [container addSubview:glassView];
        [window.contentView addSubview:container];
    } else {
        [window.contentView addSubview:glassView];
    }
    
    // Ensure proper layout
    [glassView setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
}
```

### 4. WebView Integration Considerations

#### Challenges
1. **Z-ordering**: Glass effect must be behind WebView content
2. **Performance**: Multiple glass layers may impact rendering
3. **Transparency**: WebView background must be transparent
4. **Event handling**: Glass effect shouldn't interfere with WebView interactions

#### Solutions

```objc
// Ensure proper WebView configuration
void configureWebViewForLiquidGlass(WKWebView* webView) {
    // Make WebView background transparent
    [webView setValue:@NO forKey:@"drawsBackground"];
    [webView setValue:[NSColor clearColor] forKey:@"backgroundColor"];
    
    // Ensure WebView is above glass layer
    webView.layer.zPosition = 1.0;
    
    // Optimize for performance
    webView.layer.shouldRasterize = YES;
    webView.layer.rasterizationScale = [[NSScreen mainScreen] backingScaleFactor];
}
```

### 5. Cross-Platform Strategy

#### Windows Fallback
Map Liquid Glass to the closest Windows equivalent:

```go
// Windows mapping
if options.Mac.LiquidGlass.Enabled {
    // Use Acrylic as closest equivalent
    options.Windows.BackdropType = Acrylic
    
    // Apply tint color if specified
    if options.Mac.LiquidGlass.TintColor != nil {
        options.Windows.CustomTheme.DarkModeActive.TitleBarColour = 
            NewRGBPtr(
                options.Mac.LiquidGlass.TintColor.Red,
                options.Mac.LiquidGlass.TintColor.Green,
                options.Mac.LiquidGlass.TintColor.Blue,
            )
    }
}
```

#### Linux Fallback
Use translucent background with blur if available:

```go
// Linux mapping
if options.Mac.LiquidGlass.Enabled {
    options.Linux.WindowIsTranslucent = true
    // Note: Actual blur implementation depends on compositor
}
```

## Migration Path

### Backward Compatibility
- Existing `MacBackdropTranslucent` continues to use `NSVisualEffectView`
- New `MacBackdropLiquidGlass` uses `NSGlassEffectView` when available
- Automatic fallback to `NSVisualEffectView` on older macOS versions

### Version Detection

```go
func (w *macosWebviewWindow) selectBackdropImplementation(backdrop MacBackdrop) {
    switch backdrop {
    case MacBackdropLiquidGlass:
        if w.isLiquidGlassAvailable() {
            w.applyLiquidGlass(w.options.Mac.LiquidGlass)
        } else {
            w.applyTranslucentBackdrop() // Fallback
            w.app.warn("Liquid Glass requires macOS 26.0+, using translucent backdrop")
        }
    case MacBackdropTranslucent:
        w.applyTranslucentBackdrop()
    // ... other cases
    }
}
```

## Performance Considerations

1. **Rendering Cost**: Liquid Glass has higher GPU requirements
2. **Memory Usage**: Each glass layer adds to compositor memory
3. **Battery Impact**: Continuous visual updates may impact battery life

### Optimization Strategies

```go
type MacLiquidGlass struct {
    // ... existing fields ...
    
    // Performance options
    Performance LiquidGlassPerformance
}

type LiquidGlassPerformance struct {
    // Reduce visual quality for better performance
    ReduceMotion bool
    
    // Disable liquid merging animations
    DisableLiquidMerge bool
    
    // Use static glass (no dynamic updates)
    StaticMode bool
}
```

## Testing Strategy

### Unit Tests
- Verify option parsing and validation
- Test fallback logic for unsupported systems
- Validate cross-platform mapping

### Integration Tests
- Test on macOS 26.0+ with real `NSGlassEffectView`
- Verify WebView content remains interactive
- Test window grouping and liquid merge effects
- Performance benchmarks with multiple windows

### Manual Testing Checklist
- [ ] Glass effect renders correctly
- [ ] WebView content is legible
- [ ] Tint color applies properly
- [ ] Corner radius works as expected
- [ ] Multiple windows merge when grouped
- [ ] Fallback works on older macOS versions
- [ ] Windows/Linux alternatives function correctly

## Documentation

### API Documentation
```go
// MacBackdropLiquidGlass enables Apple's Liquid Glass effect on macOS 26.0+.
// This creates a dynamic glass material that reflects and refracts light,
// providing depth and visual interest to your application windows.
// 
// On older macOS versions, this automatically falls back to MacBackdropTranslucent.
// On Windows, this maps to the Acrylic backdrop type.
// On Linux, this enables window translucency.
//
// Example:
//   options := WebviewWindowOptions{
//       Mac: MacWindow{
//           Backdrop: MacBackdropLiquidGlass,
//       },
//   }
```

### Migration Guide
```markdown
# Migrating to Liquid Glass

## From Translucent Backdrop
If you're currently using `MacBackdropTranslucent`:

Before:
```go
Mac: MacWindow{
    Backdrop: MacBackdropTranslucent,
}
```

After:
```go
Mac: MacWindow{
    Backdrop: MacBackdropLiquidGlass,
}
```

The migration is seamless - Liquid Glass will automatically fall back to 
translucent on older systems.
```

## Timeline

### Phase 1: Preparation (Week 1)
- [ ] Set up development environment with Xcode 26 beta
- [ ] Create feature branch for liquid glass implementation
- [ ] Implement OS version detection

### Phase 2: Core Implementation (Week 2-3)
- [ ] Implement NSGlassEffectView wrapper
- [ ] Add window options and configuration
- [ ] Implement fallback logic

### Phase 3: Integration (Week 4)
- [ ] WebView integration and testing
- [ ] Cross-platform fallback implementation
- [ ] Performance optimization

### Phase 4: Testing & Documentation (Week 5)
- [ ] Comprehensive testing on multiple macOS versions
- [ ] API documentation
- [ ] Example applications
- [ ] Migration guide

### Phase 5: Release (Week 6)
- [ ] Code review
- [ ] Beta release
- [ ] Community feedback incorporation
- [ ] Final release

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| NSGlassEffectView API changes | High | Monitor Apple's beta releases, maintain fallback |
| Performance issues with multiple windows | Medium | Implement performance modes, provide guidelines |
| WebView rendering conflicts | Medium | Extensive testing, clear documentation |
| Limited adoption due to OS requirements | Low | Automatic fallback, clear benefits communication |

## Conclusion

Implementing Liquid Glass in Wails v3 will provide developers with access to Apple's latest design language while maintaining the framework's commitment to simplicity and cross-platform compatibility. The proposed API design balances ease of use with flexibility, allowing both quick adoption and fine-tuned control.

The implementation leverages Wails' existing architecture, making it a natural extension of the current backdrop system. With proper fallback mechanisms and cross-platform alternatives, applications can adopt Liquid Glass without sacrificing compatibility with older systems or other platforms.

## Appendix A: Example Application

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:        "Liquid Glass Demo",
        Description: "Demonstrates Liquid Glass effect",
    })

    // Main window with liquid glass
    app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
        Title:  "Liquid Glass Window",
        Width:  800,
        Height: 600,
        Mac: application.MacWindow{
            Backdrop: application.MacBackdropLiquidGlass,
            LiquidGlass: application.MacLiquidGlass{
                Enabled:      true,
                Style:        application.LiquidGlassStyleVibrant,
                CornerRadius: 16.0,
                GroupID:      "main-group",
            },
        },
        // Fallback for other platforms
        Windows: application.WindowsWindow{
            BackdropType: application.Acrylic,
        },
        Linux: application.LinuxWindow{
            WindowIsTranslucent: true,
        },
    })

    app.Run()
}
```

## Appendix B: CSS Considerations

For optimal visual results with Liquid Glass:

```css
/* Recommended CSS for Liquid Glass windows */
body {
    background: transparent;
    /* Avoid solid backgrounds that hide the glass effect */
}

.glass-container {
    /* Use semi-transparent backgrounds */
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 16px;
}

/* Ensure text remains legible */
.glass-content {
    color: #000;
    text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

@media (prefers-color-scheme: dark) {
    .glass-content {
        color: #fff;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
    }
}
```