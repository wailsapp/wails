---
title: "私有 macOS API"
description: "依赖私有 macOS API 的所有 Wails 功能和选项，包括用于选择启用的命令以及公开构建的回退方案。"
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Wails 默认使用公开的 macOS API。单个 Go 构建标签`private_mac_apis`可启用本页列出的私有 WebKit 和 AppKit 调用。无论采用哪种构建，所有公开的 Go 选项和方法均保持可用。如果未添加该标签，仅限私有 API 的操作不会执行任何操作；存在公开替代方案的功能则会使用这些替代方案。

@note{type="caution" title="选择启用 macOS 私有行为"}
设置窗口选项不会启用私有 API。请将`private_mac_apis`添加到构建命令中以启用这些 API。该标签仅适用于 macOS 桌面构建，不适用于 iOS、Android、Windows、Linux 或服务器构建。

@end

## 启用私有 API

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

对于直接的生产构建，请使用`go build -tags production,private_mac_apis .`。对于前端示例，请按照其 README 先构建绑定和资源，再运行示例。自定义或较旧的 Taskfile 必须将`EXTRA_TAGS`传递给 Go 编译器。

## 功能清单

| 功能或值 | `private_mac_apis`启用的功能 | 未添加该标签时 |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | 位于原生窗口之上的透明 WKWebView | 原生窗口会完成配置，但 WebView 仍不透明 |
| `Mac.Backdrop: MacBackdropTranslucent` | 透明的 WKWebView，使底层的原生模糊效果能够透出 | 模糊效果会配置在不透明的 WebView 后方 |
| `Mac.Backdrop: MacBackdropLiquidGlass` | 位于玻璃层之上的透明 WKWebView，并使用私有 WebView 背景控制 | 玻璃层会配置在不透明的 WebView 后方；样式使用公开替代方案 |
| 设置 Liquid Glass 时清除 WebView 背景 | 私有 WebKit `backgroundColor`控制 | 在 macOS 12 及更高版本上使用公开的`underPageBackgroundColor`，在较旧的 macOS 上则使用图层颜色；不会使 WebView 透明 |
| `app.Window.NewNotchWindow(...)` | 异形刘海面板内的透明 WebView | 面板仍可正常工作，包括定位和动画，但其中的 WebView 仍不透明 |
| `Mac.LiquidGlass.Style` | Wails 现有的原生样式映射，包括一个未记录的深色样式值 | 使用公开的常规/透明样式和浅色/深色外观；请参阅下方的值表 |
| `Mac.LiquidGlass.GroupID` | 当标识符非空时，请求私有玻璃分组 | 忽略；不请求分组 |
| `Mac.LiquidGlass.GroupSpacing` | 当值大于零时，请求私有分组间距 | 忽略 |
| `window.OpenDevTools()`和 JavaScript `Window.OpenDevTools()` | 在 macOS 12 及更高版本上以编程方式打开 WebKit 检查器 | 不执行任何操作 |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | 在 macOS 12 及更高版本上首次显示窗口时，请求以编程方式打开检查器 | 不执行任何操作 |
| macOS 13.3 之前版本的旧式检查器启用方式 | 编译时包含检查器支持后，启用 WebKit 开发者扩展功能 | 不执行任何操作；通过公开的 Safari 检查功能进行检查需要 macOS 13.3 或更高版本 |

## WebView 透明度和背景

**需要私有 API：**`MacBackdropTransparent`、`MacBackdropTranslucent`、`MacBackdropLiquidGlass`和刘海窗口所使用的 WebView 透明效果。在内部，Wails 会设置 WebKit 的私有`drawsBackground`键。仅将 HTML 或 CSS 背景设为透明，无法使不透明的原生 WKWebView 变为透明。

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

私有 WebView 背景颜色操作使用 WebKit 的`backgroundColor`键。Liquid Glass 设置过程使用此键清除 WebView 背景。如果未添加该标签，此内部操作会在 macOS 12 及更高版本上使用公开的`underPageBackgroundColor`，在较旧的 macOS 上则使用视图的图层。这些替代方案不会使 WebView 透明。

在 macOS 上，`WebviewWindowOptions.BackgroundColour`和`window.SetBackgroundColour()`设置<strong>原生窗口</strong>的颜色，本身不需要私有 API。同样，`Frameless`和`Mac.TitleBar.AppearsTransparent`使用公开的 AppKit API；依赖私有 API 的是 WebView 透明度，而不是标题栏透明度。对于 macOS 背景效果，请配置`Mac.Backdrop`，不要仅依赖`BackgroundType`。

请参阅[窗口选项](/features/windows/options/#mac-options)、[无边框窗口](/features/windows/frameless/#with-transparent-background)和[刘海窗口](/features/windows/notch-windows/)。

## Liquid Glass 值

当原生`NSGlassEffectView`可用时（macOS 26 及更高版本），将应用以下映射。只有原生样式值`0`（常规）和`1`（透明）有正式文档。Go 常量在两种构建中均保留其现有值。

| `MacLiquidGlassStyle`值 | 添加`private_mac_apis`时 | 未添加该标签时 |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic`（`0`） | 原生常规样式（`0`） | 原生常规样式（`0`） |
| `LiquidGlassStyleLight`（`1`） | 现有的原生样式映射（`1`，透明） | 采用 Aqua 外观的原生常规样式（`0`） |
| `LiquidGlassStyleDark`（`2`） | **未记录的原生样式值`2`** | 采用深色 Aqua 外观的原生常规样式（`0`） |
| `LiquidGlassStyleVibrant`（`3`） | 映射到现有的浅色/原生透明样式（`1`） | 原生透明样式（`1`） |

Automatic 和 Vibrant 使用已有文档说明的原生样式值，但要让 Liquid Glass 背景覆盖整个窗口，仍需通过该标签启用<strong>网页视图透明效果</strong>。Light 在两种构建中的外观不同。私有样式映射无法保证未来的 macOS 版本会渲染出相同效果。

**始终属于私有 API：**`GroupID`和`GroupSpacing`。请求分组前，Wails 会检查私有的`setGroupIdentifier:`、`setGroupName:`和`setGroupSpacing:`选择器。如果没有该标签，这些操作不会产生任何效果。启用该标签并不能保证当前运行的 macOS 版本支持这些选择器。

`MacLiquidGlass.Material`、`CornerRadius`和`TintColor`本身不需要私有 API。在不提供原生 Liquid Glass 支持的 macOS 版本上，Wails 会配置半透明回退效果；要透过网页视图看到该效果，仍需启用该标签。

## 网页检查器

<strong>需要私有 API：</strong>调用`OpenDevTools()`或设置`OpenInspectorOnStartup: true`，以从应用程序中打开 WebKit 检查器。Wails 使用私有的`_inspector`选择器。如果没有`private_mac_apis`，这些操作不会产生任何效果，也不会发出提示。

还必须在编译时加入检查器支持。现有的`production`和`devtools`标签保留其原有含义：

| 构建标签 | 以编程方式打开检查器 | macOS 13.3+ 上的公开 Safari 检查功能 |
| --- | --- | --- |
| 无 | 无操作 | 已启用 |
| `private_mac_apis` | 在 macOS 12+ 上启用 | 已启用 |
| `production` | 无操作 | 已禁用 |
| `production,private_mac_apis` | 无操作 | 已禁用 |
| `production,devtools` | 无操作 | 已启用 |
| `production,devtools,private_mac_apis` | 在 macOS 12+ 上启用 | 已启用 |

在 macOS 13.3+ 上，Wails 使用公开的`WKWebView.inspectable`启用 Safari 检查功能；这不需要私有 API。在 macOS 13.3之前的版本中，回退方案使用私有的`developerExtrasEnabled`偏好设置，因此除了需要在编译时加入检查器支持外，还需要`private_mac_apis`。

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## 维护此列表

所有私有原生调用都隔离在`v3/pkg/application/mac_private_api_darwin.go`中；默认构建选择`mac_public_api_darwin.go`。此清单涵盖透明效果、网页视图背景颜色、玻璃样式、玻璃分组、打开检查器以及启用旧版检查器。更改这些实现时，应同时更新本页面和受影响的选项文档。
