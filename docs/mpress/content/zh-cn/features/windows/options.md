---
title: "窗口选项"
description: "WebviewWindowOptions 完整参考"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## 窗口配置选项

Wails 提供全面的窗口配置，包含数十个用于设置大小、位置、外观和行为的选项。本参考文档是`WebviewWindowOptions`的<strong>完整参考</strong>，涵盖 Windows、macOS 和 Linux 上的所有可用选项，并提供每个选项、每个平台的示例和约束。

## WebviewWindowOptions 结构

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window *WebviewWindow)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

`WebviewWindowOptions`**没有**`Parent`字段——如需建立父窗口/模态窗口关系，请使用`parentWindow.AttachModal(childWindow)`。它也<strong>没有</strong>`Assets`字段——资源配置位于`application.Options`（`Assets AssetOptions`）中。

完整源代码：[`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_options.go)。

## 核心选项

### Name

**类型：**`string` <strong>默认值：</strong>自动生成的 UUID <strong>平台：</strong>全部

```go
Name: "main-window"
```

<strong>用途：</strong>用于稍后查找窗口的唯一标识符。

**最佳实践：**

- 使用描述性名称：`"main"`、`"settings"`、`"about"`
- 使用 kebab-case：`"file-browser"`、`"color-picker"`
- 名称应简短且易于记忆

**示例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### Title

**类型：**`string` <strong>默认值：</strong>应用程序名称 <strong>平台：</strong>全部

```go
Title: "My Application"
```

<strong>用途：</strong>显示在标题栏和任务栏中的文本。

**动态更新：**

```go
window.SetTitle("My Application - Document.txt")
```

### Width / Height

**类型：**`int`（像素） <strong>默认值：</strong>800 x 600 <strong>平台：</strong>全部 <strong>约束：</strong>必须为正数

```go
Width:  1200,
Height: 800,
```

<strong>用途：</strong>以逻辑像素指定窗口的初始大小。

**注意事项：**

- Wails 会自动处理 DPI 缩放
- 使用逻辑像素，而非物理像素
- 考虑最低屏幕分辨率（1024x768）

**尺寸示例：**

| 使用场景 | 宽度 | 高度 |
| --- | --- | --- |
| 小型实用程序 | 400 | 300 |
| 标准应用程序 | 1024 | 768 |
| 大型应用程序 | 1440 | 900 |
| 全高清 | 1920 | 1080 |

### X/Y

**类型：**`int`（像素） <strong>默认值：</strong>在屏幕上居中 <strong>平台：</strong>全部

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

<strong>用途：</strong>窗口的初始位置。

**坐标系：**

- （0、0）是主屏幕的左上角
- X 正方向向右
- Y 正方向向下

**示例：**

只有设置了`InitialPosition: application.WindowXY`，`X`和`Y`才会生效。否则，`InitialPosition`默认为`WindowCentered`，并忽略`X`/`Y`。

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

<strong>最佳实践：</strong>如果不需要指定具体坐标，请在创建窗口后使用`Center()`将其居中：

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**类型：**`int`（像素） <strong>默认值：</strong>0（无最小值） <strong>平台：</strong>全部

```go
MinWidth:  400,
MinHeight: 300,
```

<strong>用途：</strong>防止窗口过小。

**使用场景：**

- 防止布局错乱
- 确保可用性
- 保持宽高比

**示例：**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### MaxWidth / MaxHeight

**类型：**`int`（像素） <strong>默认值：</strong>0（无最大值） <strong>平台：</strong>全部

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

<strong>用途：</strong>防止窗口过大。

**使用场景：**

- 固定大小的应用程序
- 防止资源使用过多
- 保持设计约束

## 状态选项

### Hidden

**类型：**`bool` **默认值：**`false` <strong>平台：</strong>所有平台

```go
Hidden: true,
```

<strong>用途：</strong>创建窗口但不显示。

**使用场景：**

- 后台窗口
- 按需显示的窗口
- 启动画面（创建、加载，然后显示）
- 防止加载内容时出现白屏闪烁

**平台改进：**

- <strong>Windows：</strong>已修复窗口白屏闪烁问题——调用`Show()`之前，窗口会一直保持不可见
- <strong>macOS：</strong>完全支持
- <strong>Linux：</strong>完全支持

**实现平滑加载的推荐模式：**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**示例：**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Frameless

**类型：**`bool` **默认值：**`false` <strong>平台：</strong>所有平台

```go
Frameless: true,
```

<strong>用途：</strong>移除标题栏和窗口边框。

**使用场景：**

- 自定义窗口装饰
- 启动画面
- 自助终端应用
- 采用自定义设计的窗口

<strong>重要提示：</strong>需要自行实现：

- 窗口拖动
- 关闭、最小化和最大化按钮
- 调整大小手柄（如果窗口可调整大小）

**有关详细信息，请参阅[无边框窗口](/features/windows/frameless/)。**

### DisableResize

**类型：**`bool` **默认值：**`false`（默认情况下窗口可调整大小） <strong>平台：</strong>所有平台

```go
DisableResize: true,
```

<strong>用途：</strong>阻止调整窗口大小。请注意，此字段与v2的`Resizable`含义<strong>相反</strong>——设置`DisableResize: true`可使窗口无法调整大小。

**使用场景：**

- 固定尺寸的应用
- 启动画面
- 对话框

<strong>注意：</strong>除非还通过`MaximiseButtonState`或集合行为禁用最大化和全屏，否则用户仍可将窗口最大化或切换到全屏。

### AlwaysOnTop

**类型：**`bool` **默认值：**`false` <strong>平台：</strong>所有平台

```go
AlwaysOnTop: true,
```

<strong>用途：</strong>使窗口保持在所有其他窗口之上。

**使用场景：**

- 浮动工具栏
- 通知
- 画中画
- 计时器

**平台说明：**

- <strong>macOS：</strong>完全支持
- <strong>Windows：</strong>完全支持
- <strong>Linux：</strong>取决于窗口管理器

### StartState

**类型：**`WindowState`枚举 **默认值：**`WindowStateNormal` <strong>平台：</strong>所有平台

```go
StartState: application.WindowStateMaximised,
```

<strong>用途：</strong>窗口显示时的初始状态。

**可选值：**

- `WindowStateNormal` - 普通窗口
- `WindowStateMinimised` - 最小化
- `WindowStateMaximised` - 最大化
- `WindowStateFullscreen` - 全屏

不存在`WindowStateHidden`常量——要让窗口在启动时不可见，请使用`Hidden`布尔字段。

**在运行时切换全屏模式：**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## 外观选项

### BackgroundColour

**类型：**`RGBA`结构体 <strong>默认值：</strong>白色 <strong>平台：</strong>全部

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

`RGBA`字段均为`Red, Green, Blue, Alpha`（uint8）。建议使用辅助函数`application.NewRGB(r, g, b)`（alpha 为255）或`application.NewRGBA(r, g, b, a)`。

<strong>用途：</strong>内容加载前的窗口背景颜色。

**使用场景：**

- 与应用的主题保持一致
- 防止使用深色主题时出现白屏闪烁
- 提供流畅的加载体验

**示例：**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**辅助方法：**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**类型：**`BackgroundType`枚举 **默认值：**`BackgroundTypeSolid` <strong>平台：</strong>macOS、Windows（部分支持）

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**取值：**

- `BackgroundTypeSolid` - 纯色
- `BackgroundTypeTransparent` - 完全透明
- `BackgroundTypeTranslucent` - 半透明模糊

**平台支持：**

- <strong>macOS：</strong>配置`Mac.Backdrop`；要使 WebView 透明，必须设置[`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background)。否则，WebView 将保持不透明。
- <strong>Windows：</strong>支持透明和半透明（Windows 11及更高版本）
- <strong>Linux：</strong>仅支持纯色

**示例（macOS）：**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup 和 OpenDevTools

**macOS 上的私有 API：**`OpenInspectorOnStartup: true`、Go `window.OpenDevTools()`和 JavaScript `Window.OpenDevTools()`需要`private_mac_apis`才能以编程方式打开检查器。若未设置，这些操作不会产生任何效果。生产构建还需要`devtools`。macOS 13.3及更高版本上的公开 Safari 检查功能不需要私有 API；在较旧版本的 macOS 上启用检查器则需要。请参阅[Web 检查器构建矩阵](/guides/build/private-macos-apis/#web-inspector)。

## 内容选项

### URL

**类型：**`string` <strong>默认值：</strong>空（从 Assets 加载） <strong>平台：</strong>全部

```go
URL: "https://example.com",
```

<strong>用途：</strong>加载外部 URL，而不是嵌入式资源。

**使用场景：**

- 开发（从开发服务器加载）
- 基于 Web 的应用
- 混合应用

**示例：**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

生产环境下的应用级代码片段：

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**类型：**`string` <strong>默认值：</strong>空 <strong>平台：</strong>全部

```go
HTML: "<h1>Hello World</h1>",
```

<strong>用途：</strong>直接加载 HTML 字符串。

**使用场景：**

- 简单窗口
- 生成的内容
- 测试

**示例：**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### Assets（仅限应用级）

资源配置<strong>不是</strong>`WebviewWindowOptions`字段。前端资源由应用本身通过`application.Options.Assets`（`AssetOptions`）提供；每个窗口都继承该资源服务器。

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**有关详细信息，请参阅[构建系统](/concepts/build-system/)。**

### UseApplicationMenu

**类型：**`bool` **默认值：**`false` <strong>平台：</strong>Windows、Linux（在 macOS 上无效）

```go
UseApplicationMenu: true,
```

<strong>用途：</strong>为此窗口使用应用菜单（通过`app.Menu.Set()`设置）。

在<strong>macOS</strong>上，此选项无效，因为 macOS 始终使用屏幕顶部的全局应用菜单。

在<strong>Windows</strong>和<strong>Linux</strong>上，窗口默认不显示菜单。设置`UseApplicationMenu: true`后，窗口将使用应用级菜单，从而提供一种简单的跨平台解决方案。

**示例：**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**注意事项：**

- 如果同时设置了`UseApplicationMenu`和窗口专用菜单，则窗口专用菜单优先
- 这无需在运行时检查操作系统，从而简化跨平台代码
- 有关菜单的完整文档，请参阅[应用菜单](/features/menus/application/)

## 输入选项

### EnableFileDrop

**类型：**`bool` **默认值：**`false` <strong>平台：</strong>所有平台

```go
EnableFileDrop: true,
```

<strong>用途：</strong>允许将操作系统中的文件拖放到窗口中。

启用后：

- 可以将文件管理器中的文件拖放到应用程序中
- 拖放文件时会触发`WindowFilesDropped`事件，并提供被拖放文件的路径
- 带有`data-file-drop-target`属性的元素会提供详细的拖放信息

**使用场景：**

- 文件上传界面
- 文档编辑器
- 媒体导入工具
- 任何接受文件的应用

**示例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**HTML 放置区：**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**有关完整文档，请参阅[文件拖放](/features/drag-and-drop/files/)。**

## 安全选项

### ContentProtectionEnabled

**类型：**`bool` **默认值：**`false` <strong>平台：</strong>Windows（10及更高版本）、macOS

```go
ContentProtectionEnabled: true,
```

<strong>用途：</strong>防止捕获窗口内容的屏幕画面。

**平台支持：**

- <strong>Windows：</strong>Windows 10内部版本19041及更高版本（完全支持），旧版本（部分支持）
- <strong>macOS：</strong>完全支持
- <strong>Linux：</strong>不支持

**使用场景：**

- 银行应用程序
- 密码管理器
- 医疗记录
- 机密文档

**重要说明：**

1. 无法防止使用物理相机拍摄
2. 某些工具可能会绕过保护
3. 这只是综合安全措施的一部分，不能作为唯一的保护措施
4. DevTools 窗口不会自动受到保护

**示例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### 权限

**类型：**`map[PermissionType]Permission` **默认值：**`nil`（由平台按默认方式处理） <strong>平台：</strong>Linux、Windows（macOS 交由 TCC 处理）

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

<strong>用途：</strong>以声明方式控制如何处理窗口中网页内容发出的功能请求（摄像头、麦克风、地理位置、通知和读取剪贴板），无需编写平台专用代码。

**PermissionType 值：**`PermissionMicrophone`、`PermissionCamera`、`PermissionGeolocation`、`PermissionNotifications`、`PermissionClipboardRead`

**权限值：**

- `PermissionDefault`（0）— 使用平台的原生处理方式：在 macOS/Windows 上显示 OS/WebView2 提示；在 Linux 上，允许使用摄像头和麦克风，拒绝其他所有请求
- `PermissionAllow`（1）— 不显示提示，直接授予权限（Linux：仅实现了摄像头和麦克风权限；其他类型仍会被拒绝）
- `PermissionDeny`（2）— 不显示提示，直接拒绝权限

<strong>重要说明 — Windows：</strong>在此选项出现之前，Wails 会静默授予所有 WebView2 功能权限。现在，只要在`Permissions`中设置任何条目，就会禁用这种一律授予权限的行为。未列出的功能将显示 WebView2 的原生提示，而不会被自动允许。请明确列出应用所需的每项功能。

**有关完整指南、平台支持矩阵和示例，请参阅[权限](/features/windows/permissions/)。**

## 窗口生命周期事件

使用`OnWindowEvent`和`RegisterHook`处理窗口生命周期事件。通过这些方法，可以精细控制窗口关闭和销毁行为。

### 取消关闭窗口

要阻止窗口关闭（例如存在未保存的更改），请将`RegisterHook`与`WindowClosing`事件结合使用，并调用`event.Cancel()`：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**要点：**

- `RegisterHook`在事件发生前将其拦截
- 调用`event.Cancel()`可阻止窗口关闭
- 取消关闭后，窗口将保持打开

### 处理窗口关闭

要在窗口关闭时执行清理，请将`OnWindowEvent`与`WindowClosing`事件结合使用：

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**要点：**

- `OnWindowEvent`处理即将发生的事件
- 清理操作在窗口销毁前运行
- 无法在此处取消关闭（请使用`RegisterHook`）

### 单例窗口清理模式

对于单例窗口（确保仅有一个实例），请使用`WindowClosing`清理引用：

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## 平台特定选项

### Mac 选项

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**TitleBar**（`MacTitleBar`）

- `AppearsTransparent` - 将标题栏设为透明，内容延伸至标题栏区域
- `Hide` - 完全隐藏标题栏
- `HideTitle` - 仅隐藏标题文本
- `FullSizeContent` - 将内容延伸至整个窗口

在 macOS 上，WebView 透明效果和通过编程方式打开检查器需要`private_mac_apis`构建标签。没有该标签时，这些选项仍然有效，但仅限私有 API 的操作不会产生任何效果。Liquid Glass 分组也会被忽略，样式则使用公开的替代方案。有关构建命令和确切行为，请参阅[私有 macOS API](/guides/build/private-macos-apis/)。

**Backdrop**（`MacBackdrop`）

- `MacBackdropNormal` - 标准不透明背景
- `MacBackdropTranslucent` - <strong>WebView 透明效果需要私有 API。</strong>没有该标签时，原生模糊效果仍位于不透明的 WebView 后方。
- `MacBackdropTransparent` - <strong>WebView 透明效果需要私有 API。</strong>没有该标签时，WebView 仍为不透明。
- `MacBackdropLiquidGlass` - <strong>WebView 透明效果需要私有 API。</strong>没有该标签时，玻璃效果层仍位于不透明的 WebView 后方；样式则使用公开的替代方案。

**LiquidGlass**（`MacLiquidGlass`）

| 字段或值 | 在 macOS 上对私有 API 的依赖 |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | 原生常规样式是公开的；背景效果的 WebView 透明效果需要`private_mac_apis`。 |
| `Style: LiquidGlassStyleLight` | 该标签会保留现有的原生透明样式映射；没有该标签时，Wails 使用采用 Aqua 外观的常规玻璃效果。 |
| `Style: LiquidGlassStyleDark` | <strong>私有 API：</strong>未公开的原生样式值`2`；没有该标签时，使用采用 Dark Aqua 外观的常规玻璃效果。 |
| `Style: LiquidGlassStyleVibrant` | 原生透明样式映射是公开的；背景效果的 WebView 透明效果需要该标签。 |
| `GroupID` | <strong>私有 API：</strong>非空值会请求分组；没有该标签时将被忽略。 |
| `GroupSpacing` | <strong>私有 API：</strong>正值会请求设置组间距；没有该标签时将被忽略。 |
| `Material`、`CornerRadius`、`TintColor` | 它们本身不依赖私有 API。 |

有关原生样式值和操作系统可用性，请参阅[Liquid Glass 值](/guides/build/private-macos-apis/#liquid-glass-values)。

**InvisibleTitleBarHeight**（`int`）

- 不可见标题栏区域的高度（用于拖动）
- 仅当原生标题栏拖动区域被隐藏时生效，即窗口无边框（`Frameless: true`）或使用透明标题栏（`AppearsTransparent: true`）时
- 对标题栏可见的标准窗口无效

**WindowClass**（`MacWindowClass`）

- `MacWindowClassWindow` - 标准`NSWindow`行为（默认）
- `MacWindowClassPanel` - 一种永远不会成为应用程序主窗口的辅助`NSPanel`

`PanelPreferences`仅适用于`MacWindowClassPanel`：

- `NonActivating`会添加`NSWindowStyleMaskNonactivatingPanel`。显示面板或使其获得焦点不会激活 Wails 应用程序，但面板仍可成为按键窗口，以便操作控件和输入文本。
- `FloatingPanel`会启用 AppKit 的浮动面板行为。
- 仅当所点击的视图请求键盘输入时，`BecomesKeyOnlyIfNeeded`才会获得按键窗口状态。
- `UtilityWindow`会应用原生实用工具窗口样式。

Wails 面板在应用程序停用时仍保持可见，并在关闭时释放，这与`WebviewWindow`预期的生命周期一致。这些设置有意覆盖`NSPanel`中相反的默认值。

窗口类、窗口层级、激活策略和集合行为分别解决不同的问题：

- `WindowClass`用于选择`NSWindow`或`NSPanel`，并控制主窗口和按键窗口的语义。
- `WindowLevel`控制 Z 轴顺序。对于菜单栏覆盖层，请使用`MacWindowLevelPopUpMenu`。
- `MacOptions.ActivationPolicy`控制整个应用程序，包括 Dock 和菜单栏的呈现方式。非激活面板不需要辅助应用激活策略，但应用仍可使用该策略隐藏其 Dock 图标。
- `CollectionBehavior`控制窗口是否参与 Spaces 和全屏模式。

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**WindowLevel**（`MacWindowLevel`）

- `MacWindowLevelNormal` - 标准窗口层级（默认）
- `MacWindowLevelFloating` - 浮动于普通窗口上方
- `MacWindowLevelTornOffMenu` - 分离式菜单层级
- `MacWindowLevelModalPanel` - 模态面板层级
- `MacWindowLevelMainMenu` - 主菜单层级
- `MacWindowLevelStatus` - 状态窗口层级
- `MacWindowLevelPopUpMenu` - 弹出式菜单层级
- `MacWindowLevelScreenSaver` - 屏幕保护程序层级

显式设置的`WindowLevel`优先于`AlwaysOnTop`和`PanelPreferences.FloatingPanel`。如果未显式设置层级，则`AlwaysOnTop`或浮动面板会解析为`MacWindowLevelFloating`；否则，层级为`MacWindowLevelNormal`。此后调用`SetAlwaysOnTop`仍属于显式的运行时更改。

**CollectionBehavior**（`MacWindowCollectionBehavior`）

控制窗口在 macOS Spaces 和全屏模式下的行为。这些值是位掩码，可以使用按位或（`|`）组合。

**Space 行为：**

- `MacWindowCollectionBehaviorDefault` - 使用 FullScreenPrimary（默认，向后兼容）
- `MacWindowCollectionBehaviorCanJoinAllSpaces` - 窗口显示在所有 Space 上
- `MacWindowCollectionBehaviorMoveToActiveSpace` - 显示时移至当前活动的 Space
- `MacWindowCollectionBehaviorManaged` - 默认的托管窗口行为
- `MacWindowCollectionBehaviorTransient` - 临时/瞬态窗口
- `MacWindowCollectionBehaviorStationary` - 切换 Space 时保持原位

**窗口循环切换：**

- `MacWindowCollectionBehaviorParticipatesInCycle` - 包含在 Cmd+` 循环切换中
- `MacWindowCollectionBehaviorIgnoresCycle` - 不包含在 Cmd+` 循环切换中

**全屏行为：**

- `MacWindowCollectionBehaviorFullScreenPrimary` - 可以进入全屏模式
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - 可以覆盖全屏应用
- `MacWindowCollectionBehaviorFullScreenNone` - 禁用全屏功能
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` - 允许并排平铺（macOS 10.11+）
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` - 阻止平铺（macOS 10.11+）

**示例 - 类似 Spotlight 的窗口：**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**示例 - 单一行为：**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode**（`MacWindowTabbingMode`）

控制 macOS 10.12及更高版本上的窗口标签页行为。窗口标签页功能可将多个窗口组合为标签页。

**选项：**

- `MacWindowTabbingModeDefault` - 零值哨兵（未显式设置）。运行时默认为禁止使用标签页
- `MacWindowTabbingModeAutomatic` - 由系统决定标签页行为
- `MacWindowTabbingModePreferred` - 窗口倾向于使用标签页模式
- `MacWindowTabbingModeDisallowed` - 禁用窗口标签页功能

**示例 - 禁用窗口标签页功能：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**示例 - 优先使用窗口标签页：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences**（`MacWebviewPreferences`）

精细控制底层`WKWebView`配置。所有字段均为可选；未设置的字段不会改变 WebKit 的默认值。

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — 当设为`true`时，按 Tab 键会将焦点移至链接和表单控件（默认值：`false`）
- `TextInteractionEnabled` — 当设为`true`时，用户可以选择网页视图中的文本并与其交互（默认值：`true`）
- `FullscreenEnabled` — 当设为`true`时，网页内容可以通过 HTML Fullscreen API 进入全屏模式（默认值：`false`）。需要 macOS 12.3 或更高版本。
- `AllowsBackForwardNavigationGestures` — 当设为`true`时，水平轻扫手势会触发后退/前进导航（默认值：`false`）
- `AllowsMagnification` — 当设为`true`时，网页视图将启用双指缩放（默认值：`false`）
- `AllowsAirPlayForMediaPlayback` — 当设为`true`时，可以将媒体流式传输到 AirPlay 设备（默认值：`true`）
- `JavaScriptCanOpenWindowsAutomatically` — 当设为`true`时，JavaScript 无需用户手势即可打开新窗口（默认值：`false`）
- `MinimumFontSize` — 以点为单位的最小字体大小。使用`optional.NewVar(12.0)`进行设置。未设置时保留 WebKit 默认值。
- `ApplicationNameForUserAgent` — 覆盖 WebKit 用户代理字符串中的应用程序名称后缀。当网站拒绝默认的`"wails.io"`标识符时（例如 YouTube 嵌入内容），此选项很有用。留空则保留默认值。
- `EnableAutoplayWithoutUserAction` — 当设为`true`时，音频和视频无需用户手势即可自动播放。映射到`WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone`（默认值：`false`）

### Windows 选项（每个窗口）

每个窗口使用的结构体是`application.WindowsWindow`，**不是**`WindowsOptions`（后者是<em>应用程序</em>级结构体）。

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon**（`bool`）

- 从标题栏中移除图标。

**DisableMenu**（`bool`）

- 禁用窗口的菜单栏。当设为`true`时，即使已配置菜单栏，窗口也不会显示它。
- 默认值：`false`

**BackdropType**（`BackdropType`）

- `application.Auto` — 系统默认值
- `application.None` — 无背景材质
- `application.Mica` — Mica 材质（Windows 11）
- `application.Acrylic` — Acrylic 材质（Windows 11）
- `application.Tabbed` — Tabbed 材质（Windows 11）

没有`WindowsBackdropTypeMica`样式的常量，请使用`application.Mica`等。

**CustomTheme**（`ThemeSettings`）

- 值（不是指针）。为窗口边框、标题栏文本和背景以及菜单栏自定义深色/浅色模式颜色。

**DisableFramelessWindowDecorations**（`bool`）

- 禁用默认的无边框窗口装饰（Aero 阴影、圆角）。

**NonClientRegionSupport**（`bool`）

- 为无边框自定义标题栏启用 WebView2 原生的`app-region: drag`/`app-region: no-drag`支持。
- 此选项仅用于简单的原生应用拖动。它不为自定义标题按钮提供原生行为，也不为自定义最大化按钮提供 Windows 11 Snap Assist / Snap Layouts。

**WebView2CompositionHosting**（`bool`）

- 为具有原生 Windows 行为的自定义标题按钮启用由 Wails 管理的`--wails-non-client-region`支持，包括自定义最大化按钮上的 Windows 11 Snap Assist / Snap Layouts。
- 实验性功能。此选项通过`ICoreWebView2CompositionController`和 DirectComposition 托管 WebView2，而不是使用默认的 HWND 托管控制器。
- 当窗口同时需要 WebView2 原生的`app-region`支持和由 Wails 管理的自定义标题按钮区域时，可以将其与`NonClientRegionSupport`结合使用。

**示例：**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**示例 — 自定义 Windows 标题栏区域：**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

有关详细行为、权衡以及配套 CSS，请参阅[无边框窗口](/features/windows/frameless/#native-non-client-regions-on-windows)。

### Linux 选项（每个窗口）

每个窗口使用的结构体是`application.LinuxWindow`，**不是**`LinuxOptions`。

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon**（`[]byte`）

- 窗口图标（PNG 格式）。

**WindowIsTranslucent**（`bool`）

- 需要合成器支持。

**示例：**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## 应用级 Windows 选项

某些 Windows 专用选项必须在应用级而非逐窗口配置。这是因为对于每个用户数据路径，WebView2 仅共享一个浏览器环境。

### 浏览器标志

WebView2 浏览器标志用于控制应用中<strong>所有窗口</strong>的实验性功能和行为。必须在`application.Options.Windows`中设置这些标志：

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures**（`[]string`）

- 要启用的 WebView2 功能标志列表
- 有关可用标志，请参阅[WebView2 浏览器标志](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)
- 示例：`"msWebView2EnableDraggableRegions"`

**DisabledFeatures**（`[]string`）

- 要禁用的 WebView2 功能标志列表
- Wails 会自动禁用`msSmartScreenProtection`
- 示例：`"msExperimentalFeature"`

**AdditionalBrowserArgs**（`[]string`）

- 传递给浏览器进程的 Chromium 命令行参数
- 必须包含`--`前缀（例如`"--remote-debugging-port=9222"`）
- 有关可用参数，请参阅[Chromium 命令行开关](https://peter.sh/experiments/chromium-command-line-switches/)

@note{type="caution" title="重要"}
由于对于每个用户数据路径，WebView2 仅共享一个浏览器环境，因此这些标志会全局应用于所有窗口。无法为不同窗口设置不同的浏览器标志。

@end

**完整示例：**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## 完整示例

下面是一份可用于生产环境的窗口配置：

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

前端资源在应用级（`application.Options.Assets`）提供，而不是逐窗口提供。

## 后续步骤

- [窗口基础](/features/windows/basics/)——创建和控制窗口
- [多窗口](/features/windows/multiple/)——多窗口模式
- [无边框窗口](/features/windows/frameless/)——自定义窗口装饰
- [窗口事件](/features/windows/events/)——生命周期事件

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[示例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
