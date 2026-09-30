---
title: "无边框窗口"
description: "使用无边框窗口创建自定义窗口装饰"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## 无边框窗口

Wails 提供基于 CSS 拖动区域并具有平台原生行为的<strong>无边框窗口支持</strong>。移除平台原生标题栏后，便可完全控制窗口装饰、创建自定义设计和独特的用户体验，同时保留拖动、调整大小和系统控件等基本功能。

![作为无边框窗口运行，并采用 macOS 原生圆角的默认 Wails v3 TypeScript 入门应用](/assets/screenshots/frameless-v3-native-corners-macos.png)

上例是启用了`Frameless: true`的默认 Wails v3 TypeScript 入门应用。

## 快速开始

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**可拖动标题栏的 CSS：**

```css
.titlebar {
    --wails-draggable: drag;
    height: 40px;
    background: #333;
}

.titlebar button {
    --wails-draggable: no-drag;
}
```

**HTML：**

```html
<div class="titlebar">
    <span>My Application</span>
    <button onclick="window.close()">×</button>
</div>
```

<strong>就这么简单！</strong>现在你已经有了自定义标题栏。

## 创建无边框窗口

### 圆角半径（macOS）

默认情况下，无边框窗口会保留 AppKit 的标准 macOS 圆角。设置`Mac.CornerRadius`以使用自定义半径（单位为点）：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

将`Mac.CornerType`设为`MacWindowCornerTypeSquare`可使用直角。此设置会忽略`CornerRadius`：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### 基本无边框窗口

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**你将获得：**

- 无标题栏
- 无窗口边框
- 无系统按钮
- 透明背景（可选）

**你需要实现：**

- 可拖动区域
- 关闭、最小化和最大化按钮
- 调整大小手柄（如果窗口可调整大小）

### 使用透明背景

<strong>macOS 上的私有 API：</strong>设置`Mac.Backdrop: application.MacBackdropTransparent`并使用`-tags private_mac_apis`构建，才能使 WebView 透明。如果没有该标签，即使 HTML/CSS 背景透明，原生 WebView 仍不透明。`Frameless`和`TitleBar.AppearsTransparent`本身使用公共 API。请参阅[macOS 私有 API](/guides/build/private-macos-apis/#webview-transparency-and-background)。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**使用场景：**

- 圆角
- 自定义形状
- 叠加窗口
- 启动画面

## 拖动区域

### 基于 CSS 的拖动

使用`--wails-draggable` CSS 属性：

```css
/* Draggable area */
.titlebar {
    --wails-draggable: drag;
}

/* Non-draggable elements within draggable area */
.titlebar button {
    --wails-draggable: no-drag;
}
```

**取值：**

- `drag`——区域可拖动
- `no-drag`——区域不可拖动（即使父元素可拖动）

### 完整标题栏示例

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #2c2c2c;
    color: white;
    padding: 0 16px;
}

.title {
    font-size: 14px;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: white;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
}
```

**按钮使用的 JavaScript：**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Windows 上的原生非客户区

Windows 可以将自定义标题栏的某些部分视为原生非客户区。这样，你既能使用任意 HTML/CSS 设计绘制标题栏和标题栏按钮，又能保留 Windows 原生行为：标题区域可拖动窗口，最大化按钮可显示 Windows 11的贴靠助手/贴靠布局，最小化、最大化和关闭按钮则可接收原生命中测试及鼠标状态。

下面的视频展示了一个使用 Windows 原生命中测试的自定义 HTML/CSS 标题栏，其中包括自定义最大化按钮上的 Windows 11贴靠助手/贴靠布局。

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails 支持两种 Windows 专用机制：

- `app-region`：通过 WebView2 的原生非客户区支持实现
- `--wails-non-client-region`：通过 Wails 运行时跟踪自定义标题栏按钮实现

### 选择模式

@note{type="caution" title="实验性功能"}
`WebView2CompositionHosting`会改变窗口在底层托管 WebView2 以及与其交互的方式。Wails 不再使用默认的 HWND 托管 WebView2 控制器，而是采用合成控制器托管并显式转发输入。此模式可能存在渲染、输入、焦点或 WebView2 Runtime 兼容性问题。仅当你需要自定义标题栏按钮具备原生行为时才启用此模式，并在你支持的 Windows 和 WebView2 Runtime 版本上仔细测试应用。

@end

请根据你需要的 Windows 功能进行选择：

- 如需通过 WebView2 的`app-region: drag`和`app-region: no-drag`实现简单的原生应用拖动，请使用`NonClientRegionSupport`。
- 如果你希望自定义的最小化、最大化和关闭按钮像 Windows 原生标题栏按钮一样工作，请使用`WebView2CompositionHosting`。
- 如果同一窗口既需要 WebView2 原生`app-region`支持，又需要由 Wails 管理的自定义标题栏按钮区域，请同时启用两者。

`NonClientRegionSupport`是 Wails `--wails-draggable`跟踪机制的一种轻量级原生替代方案。你通过 CSS 标记可拖动和不可拖动区域，由 WebView2 决定哪些像素属于标题区域，并在命中测试时由 Wails 向 WebView2 查询原生区域。

这就是此模式目前的全部功能范围。它不会让自定义的最小化、最大化或关闭按钮像 Windows 原生标题栏按钮一样工作，也不会为自定义最大化按钮启用 Windows 11贴靠助手/贴靠布局。如果你只需要简单的原生应用拖动，而不需要`--wails-draggable`的额外机制，请使用此模式。

`WebView2CompositionHosting`用于实现具有原生行为的自定义标题栏按钮。Wails 会跟踪以`--wails-non-client-region`标记的 DOM 矩形区域，将其映射到`HTMINBUTTON`、`HTMAXBUTTON`和`HTCLOSE`等 Windows 命中测试值，并将鼠标输入转发回由合成托管的 WebView2 表面。这样，自定义最大化按钮既能保留你选择的任意视觉设计，又能参与 Windows 11的贴靠助手和贴靠布局。

换句话说：`NonClientRegionSupport`是 WebView2 原生的 CSS 区域支持。`WebView2CompositionHosting`则表示由 Wails 负责宿主拥有的合成和自定义非客户区命中测试。

### WebView2 app-region

为窗口启用 WebView2 的原生非客户区支持：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

然后使用 CSS `app-region`属性标记可拖动区域：

```css
.titlebar {
    app-region: drag;
}

.titlebar button,
.titlebar input,
.titlebar select,
.titlebar textarea {
    app-region: no-drag;
}
```

如果你只需要原生标题栏拖动，并通过常规前端点击处理标题栏控件，请使用此模式。

此模式的限制在于它受 WebView2 自身的非客户区支持所约束。在当前的 WebView2 版本中，这意味着仅支持可拖动和不可拖动区域。它并非用于对具有不同原生最小化、最大化和关闭角色的完全自定义前端标题栏按钮进行建模。

### 具有原生行为的自定义标题栏按钮

对于应像系统标题栏按钮一样工作的自定义最小化、最大化和关闭按钮，请启用合成托管：

@note{type="caution" title="实验性功能"}
`WebView2CompositionHosting`使用基于 DirectComposition 的 WebView2 合成控制器托管。启用前，请参阅[选择模式](#heading-7)。

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

然后使用`--wails-non-client-region`标记每个前端区域：

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="window-controls">
        <button class="window-button minimize" aria-label="Minimize"></button>
        <button class="window-button maximize" aria-label="Maximize"></button>
        <button class="window-button close" aria-label="Close"></button>
    </div>
</div>
```

```css
.titlebar {
    --wails-non-client-region: caption;
    height: 40px;
}

.window-controls {
    display: flex;
    height: 100%;
}

.window-button {
    width: 46px;
    border: 0;
    background: transparent;
}

.window-button.minimize {
    --wails-non-client-region: minimize;
}

.window-button.maximize {
    --wails-non-client-region: maximize;
}

.window-button.close {
    --wails-non-client-region: close;
}
```

支持的`--wails-non-client-region`值：

- `caption` - 可拖动的标题栏区域
- `minimize` - 原生最小化按钮的命中目标
- `maximize` - 原生最大化按钮的命中目标，包括 Windows 11贴靠助手/贴靠布局的悬停行为
- `close` - 原生关闭按钮的命中目标

Wails 运行时会监测 DOM、样式、尺寸、滚动和视口的变化，然后将区域快照发送到原生窗口。区域几何尺寸以 CSS 像素为单位测量，并转换为物理像素，以供 Windows 进行命中测试。

视觉设计完全由你决定。这些区域只会告诉 Windows 每个矩形的含义；按钮的形状、图标、颜色、间距、悬停样式和布局仍由前端提供。

### 同时使用两种模式

如果希望同一窗口同时支持 WebView2 `app-region`和由 Wails 管理的标题栏按钮区域，可以同时启用这两个选项：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## 系统按钮

### 实现关闭/最小化/最大化

**Go 端：**

```go
type WindowControls struct {
    window *application.WebviewWindow
}

func (wc *WindowControls) Minimise() {
    wc.window.Minimise()
}

func (wc *WindowControls) Maximise() {
    if wc.window.IsMaximised() {
        wc.window.UnMaximise()
    } else {
        wc.window.Maximise()
    }
}

func (wc *WindowControls) Close() {
    wc.window.Close()
}
```

**JavaScript 端：**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**也可以使用运行时方法：**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### 切换最大化状态

跟踪最大化状态以更新按钮图标：

```javascript
import { Window } from '@wailsio/runtime'

async function toggleMaximise() {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
}

async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    const button = document.querySelector('.maximize')
    button.textContent = isMaximised ? '❐' : '□'
}
```

## 调整大小手柄

### 基于 CSS 调整大小

Wails 为无边框窗口提供自动调整大小手柄：

```css
/* Enable resize on all edges */
body {
    --wails-resize: all;
}

/* Or specific edges */
.resize-top {
    --wails-resize: top;
}

.resize-bottom {
    --wails-resize: bottom;
}

.resize-left {
    --wails-resize: left;
}

.resize-right {
    --wails-resize: right;
}

/* Corners */
.resize-top-left {
    --wails-resize: top-left;
}

.resize-top-right {
    --wails-resize: top-right;
}

.resize-bottom-left {
    --wails-resize: bottom-left;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
}
```

**可用值：**

- `all` - 从所有边缘调整大小
- `top`、`bottom`、`left`、`right` - 指定边缘
- `top-left`、`top-right`、`bottom-left`、`bottom-right` - 各个角
- `none` - 禁止调整大小

### 调整大小手柄示例

```html
<div class="window">
    <div class="titlebar">...</div>
    <div class="content">...</div>
    <div class="resize-handle resize-bottom-right"></div>
</div>
```

```css
.resize-handle {
    position: absolute;
    width: 16px;
    height: 16px;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
    bottom: 0;
    right: 0;
    cursor: nwse-resize;
}
```

## 特定于平台的行为

@tabs{sync-key="platform"}
[Windows]
**Windows 无边框窗口：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**功能：**

- 自动投射阴影
- 支持贴靠布局（Windows 11）
- 支持 Aero Snap
- DPI 缩放

**禁用窗口装饰：**

```go
Windows: application.WindowsWindow{
    DisableFramelessWindowDecorations: true,
},
```

**贴靠助手：**

```go
// Trigger Windows 11 Snap Assist
window.SnapAssist()
```

这会通过 Windows 热键路径触发贴靠布局。对于需要原生悬停贴靠布局的自定义 HTML 最大化按钮，请改用[Windows 上的原生非客户区](#windows-)。

**自定义标题栏高度：** Windows 会自动从 CSS 中检测拖动区域。

[macOS]
**macOS 无边框窗口：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        InvisibleTitleBarHeight: 40,
    },
})
```

**功能：**

- 原生全屏支持
- 交通灯按钮（可选）
- Vibrancy视觉效果
- 透明标题栏

**完全隐藏标题栏**（请使用`application`包导出的预设变体——不存在`TitleBarStyle`字段或`MacTitleBarStyleHidden`常量）：

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

其他预设包括`MacTitleBarDefault`、`MacTitleBarHiddenInset`和`MacTitleBarHiddenInsetUnified`。

**不可见标题栏：** 隐藏标题栏的同时仍允许拖动。仅当窗口为无边框窗口或使用`AppearsTransparent`时才会生效：

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Linux 无边框窗口：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**特性：**

- 基本的无边框支持
- CSS 拖动区域
- 因桌面环境而异

**桌面环境说明：**

- <strong>GNOME：</strong>支持良好
- <strong>KDE Plasma：</strong>支持良好
- <strong>XFCE：</strong>基本支持
- <strong>平铺式窗口管理器：</strong>支持有限

**需要合成器：** 透明效果需要合成器（大多数现代桌面环境都配有合成器）。

@end

## 常用模式

### 模式1：现代标题栏

```html
<div class="modern-titlebar">
    <div class="app-icon">
        <img src="/icon.png" alt="App Icon">
    </div>
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.modern-titlebar {
    --wails-draggable: drag;
    display: flex;
    align-items: center;
    height: 40px;
    background: linear-gradient(to bottom, #3a3a3a, #2c2c2c);
    border-bottom: 1px solid #1a1a1a;
    padding: 0 16px;
}

.app-icon {
    --wails-draggable: no-drag;
    width: 24px;
    height: 24px;
    margin-right: 12px;
}

.title {
    flex: 1;
    font-size: 13px;
    color: #e0e0e0;
    user-select: none;
}

.controls {
    display: flex;
    gap: 1px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 46px;
    height: 32px;
    border: none;
    background: transparent;
    color: #e0e0e0;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
    color: white;
}
```

### 模式2：启动画面

```go
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "Loading...",
    Width:          400,
    Height:         300,
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
}

.splash {
    background: white;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    padding: 40px;
    text-align: center;
}
```

### 模式3：圆角窗口

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    margin: 8px;
}

.window {
    background: white;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    height: calc(100vh - 16px);
}

.titlebar {
    --wails-draggable: drag;
    background: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
}
```

### 模式4：叠加窗口

```go
overlay := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
}

.overlay {
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 8px;
    padding: 20px;
}
```

## 完整示例

下面是一个可用于生产环境的无边框窗口：

**Go：**

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "Frameless App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Frameless Application",
        Width:     1000,
        Height:    700,
        MinWidth:  800,
        MinHeight: 600,
        Frameless: true,

        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            InvisibleTitleBarHeight: 40,
        },

        Windows: application.WindowsWindow{
            DisableFramelessWindowDecorations: false,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

**HTML：**

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="window">
        <div class="titlebar">
            <div class="title">Frameless Application</div>
            <div class="controls">
                <button class="minimize" title="Minimise">−</button>
                <button class="maximize" title="Maximise">□</button>
                <button class="close" title="Close">×</button>
            </div>
        </div>
        <div class="content">
            <h1>Hello from Frameless Window!</h1>
        </div>
    </div>
    <script src="/main.js" type="module"></script>
</body>
</html>
```

**CSS：**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f5f5;
}

.window {
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #ffffff;
    border-bottom: 1px solid #e0e0e0;
    padding: 0 16px;
}

.title {
    font-size: 13px;
    font-weight: 500;
    color: #333;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #666;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.controls button:hover {
    background: #f0f0f0;
    color: #333;
}

.controls .close:hover {
    background: #e81123;
    color: white;
}

.content {
    flex: 1;
    padding: 40px;
    overflow: auto;
}
```

**JavaScript：**

```javascript
import { Window } from '@wailsio/runtime'

// Minimise button
document.querySelector('.minimize').addEventListener('click', () => {
    Window.Minimise()
})

// Maximise/restore button
const maximiseBtn = document.querySelector('.maximize')
maximiseBtn.addEventListener('click', async () => {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
})

// Close button
document.querySelector('.close').addEventListener('click', () => {
    Window.Close()
})

// Update maximise button icon
async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    maximiseBtn.textContent = isMaximised ? '❐' : '□'
    maximiseBtn.title = isMaximised ? 'Restore' : 'Maximise'
}

// Initial state
updateMaximiseButton()
```

## 最佳实践

### ✅ 应该做

- **提供可拖动区域**——用户需要移动窗口
- **实现系统按钮**——关闭、最小化和最大化
- **设置最小尺寸**——避免布局变得无法使用
- **在所有平台上测试**——行为因平台而异
- **使用 CSS 定义拖动区域**——灵活且易于维护
- **提供视觉反馈**——为按钮添加悬停状态

### ❌ 不要做

- **不要忘记添加调整大小手柄**——如果窗口可调整大小
- **不要让整个窗口都可拖动**——这会阻止交互
- **不要忘记将按钮设为不可拖动**——否则按钮将无法使用
- **不要使用过小的拖动区域**——难以拖住
- **不要忽略平台差异**——请充分测试

## 故障排除

### 窗口无法拖动

<strong>原因：</strong>缺少`--wails-draggable: drag`

**解决方案：**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### 按钮无法使用

<strong>原因：</strong>按钮位于可拖动区域内

**解决方案：**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### 无法调整窗口大小

<strong>原因：</strong>缺少调整大小手柄

**解决方案：**

```css
body {
    --wails-resize: all;
}
```

## 后续步骤

@cards{cols="2"}
▣ 窗口基础
了解窗口管理的基础知识。

[了解更多 →](/features/windows/basics/)

---
⚙ 窗口选项
窗口选项的完整参考。

[了解更多 →](/features/windows/options/)

---
🚀 窗口事件
处理窗口生命周期事件。

[了解更多 →](/features/windows/events/)

---
◆ 多窗口
多窗口应用程序的常用模式。

[了解更多 →](/features/windows/multiple/)

@end

---

<strong>有问题？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[无边框窗口示例](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless)。
