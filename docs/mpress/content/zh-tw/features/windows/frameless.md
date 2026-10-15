---
title: "無框視窗"
description: "使用無框視窗建立自訂視窗框架"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## 無框視窗

Wails 提供<strong>無框視窗支援</strong>，包含以 CSS 為基礎的拖曳區域及平台原生行為。移除平台原生標題列，即可完全掌控視窗框架、自訂設計及獨特的使用者體驗，同時保留拖曳、調整大小及系統控制項等必要功能。

![以無框視窗執行並具有原生 macOS 圓角的預設 Wails v3 TypeScript 入門應用程式](/assets/screenshots/frameless-v3-native-corners-macos.png)

上例是啟用`Frameless: true`的預設 Wails v3 TypeScript 入門應用程式。

## 快速開始

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**可拖曳標題列的 CSS：**

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

<strong>就這麼簡單！</strong>現在您已有自訂標題列。

## 建立無框視窗

### 圓角半徑（macOS）

無框視窗預設會保留 AppKit 標準的 macOS 圓角。設定`Mac.CornerRadius`即可使用自訂半徑（以點為單位）：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

將`Mac.CornerType`設為`MacWindowCornerTypeSquare`可使用直角。此設定會忽略`CornerRadius`：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### 基本無框視窗

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**您會得到：**

- 無標題列
- 無視窗邊框
- 無系統按鈕
- 透明背景（選用）

**您需要實作：**

- 可拖曳區域
- 關閉／最小化／最大化按鈕
- 調整大小控點（若可調整大小）

### 使用透明背景

<strong>macOS 上的私有 API：</strong>若要讓 WebView 透明，請設定`Mac.Backdrop: application.MacBackdropTransparent`，並使用`-tags private_mac_apis`建置。若未使用此標記，即使 HTML/CSS 背景透明，原生 WebView 仍會保持不透明。`Frameless`與`TitleBar.AppearsTransparent`本身使用公用 API。請參閱[macOS 私有 API](/guides/build/private-macos-apis/#webview-transparency-and-background)。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**使用案例：**

- 圓角
- 自訂形狀
- 覆疊視窗
- 啟動畫面

## 拖曳區域

### 以 CSS 為基礎的拖曳

使用`--wails-draggable` CSS 屬性：

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

**值：**

- `drag`－區域可拖曳
- `no-drag`－區域不可拖曳（即使父元素可拖曳）

### 完整標題列範例

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

**按鈕的 JavaScript：**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Windows 上的原生非用戶區域

Windows 可將自訂標題列的部分區域視為原生非用戶區域。如此一來，您便能使用任意 HTML/CSS 設計繪製標題列及標題按鈕，同時保留 Windows 原生行為：標題區域可拖曳視窗、最大化按鈕可顯示 Windows 11貼齊小幫手／貼齊配置，而最小化、最大化及關閉按鈕則可接收原生命中測試與滑鼠狀態。

下方影片展示使用 Windows 原生命中測試的自訂 HTML/CSS 標題列，其中包括自訂最大化按鈕上的 Windows 11貼齊小幫手／貼齊配置。

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails 支援兩種 Windows 專用機制：

- 透過 WebView2 原生非用戶區域支援的`app-region`
- 透過 Wails 執行階段追蹤自訂標題按鈕的`--wails-non-client-region`

### 選擇模式

@note{type="caution" title="實驗性功能"}
`WebView2CompositionHosting`會從底層變更視窗託管 WebView2 及與其互動的方式。Wails 不再使用預設由 HWND 託管的 WebView2 控制器，而是改用合成控制器託管並明確轉送輸入。此模式可能存在算繪、輸入、焦點或 WebView2 Runtime 相容性問題。只有在需要自訂標題按鈕的原生行為時才啟用此模式，並在您支援的 Windows 與 WebView2 Runtime 版本上仔細測試應用程式。

@end

請依您需要的 Windows 功能選擇：

- 若要透過 WebView2 的`app-region: drag`與`app-region: no-drag`實現簡單的原生應用程式拖曳，請使用`NonClientRegionSupport`。
- 若希望自訂的最小化、最大化及關閉按鈕能如同 Windows 原生標題按鈕般運作，請使用`WebView2CompositionHosting`。
- 若同一個視窗同時需要 WebView2 原生的`app-region`支援，以及由 Wails 管理的自訂標題按鈕區域，請同時啟用兩者。

`NonClientRegionSupport`是 Wails `--wails-draggable`追蹤機制的輕量原生替代方案。您使用 CSS 標記可拖曳與不可拖曳區域，由 WebView2 判斷哪些像素屬於標題區域，而 Wails 會在執行命中測試時向 WebView2 查詢原生區域。

目前此模式的完整範圍僅限於此。它不會讓自訂的最小化、最大化或關閉按鈕如同 Windows 原生標題按鈕般運作，也不會為自訂最大化按鈕啟用 Windows 11貼齊小幫手／貼齊配置。若您只需要簡單的原生應用程式拖曳，而不需要`--wails-draggable`的額外機制，請使用此模式。

`WebView2CompositionHosting`適用於具備原生行為的自訂標題列按鈕。Wails 會追蹤以`--wails-non-client-region`標記的 DOM 矩形區域，將其對應至`HTMINBUTTON`、`HTMAXBUTTON`和`HTCLOSE`等 Windows 命中測試值，再將滑鼠輸入轉送回以合成模式託管的 WebView2 介面。因此，自訂最大化按鈕可以支援 Windows 11的貼齊小幫手／貼齊版面配置，同時保留您選擇的任何視覺設計。

換句話說：`NonClientRegionSupport`是 WebView2 原生的 CSS 區域支援；`WebView2CompositionHosting`則是由 Wails 負責主機端合成與自訂非工作區命中測試。

### WebView2 app-region

為視窗啟用 WebView2 的原生非工作區支援：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

接著，使用 CSS `app-region`屬性標記可拖曳區域：

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

如果您只需要原生標題列拖曳功能，而標題列控制項以一般前端點擊處理，請使用此模式。

此模式受限於 WebView2 本身的非工作區支援。目前的 WebView2 版本只支援可拖曳和不可拖曳區域，並非用來建模具有不同原生最小化、最大化和關閉角色的完全自訂前端標題列按鈕。

### 具備原生行為的自訂標題列按鈕

若要讓自訂的最小化、最大化和關閉按鈕像系統標題列按鈕一樣運作，請啟用合成託管：

@note{type="caution" title="實驗性功能"}
`WebView2CompositionHosting`使用搭配 DirectComposition 的 WebView2 合成控制器託管。啟用前，請參閱[選擇模式](#heading-7)。

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

接著，以`--wails-non-client-region`標記各個前端區域：

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

支援的`--wails-non-client-region`值：

- `caption`－可拖曳的標題列區域
- `minimize`－原生最小化按鈕的點擊目標
- `maximize`－原生最大化按鈕的點擊目標，包括 Windows 11貼齊小幫手／貼齊版面配置的暫留行為
- `close`－原生關閉按鈕的點擊目標

Wails 執行階段會監控 DOM、樣式、大小、捲動和檢視區的變更，再將區域快照傳送至原生視窗。區域幾何資訊以 CSS 像素測量，並轉換為實體像素，以供 Windows 進行命中測試。

視覺設計仍完全由您掌控。這些區域只會告知 Windows 各矩形所代表的意義；按鈕形狀、圖示、色彩、間距、暫留樣式和版面配置仍由前端提供。

### 結合兩種模式

若要在同一個視窗中同時使用 WebView2 `app-region`支援和由 Wails 管理的標題列按鈕區域，可以同時啟用這兩個選項：

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## 系統按鈕

### 實作關閉／最小化／最大化

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

**或者使用執行階段方法：**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### 切換最大化狀態

追蹤最大化狀態，以更新按鈕圖示：

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

## 調整大小控點

### 使用 CSS 調整大小

Wails 為無框線視窗提供自動調整大小控點：

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

**值：**

- `all`－從所有邊緣調整大小
- `top`、`bottom`、`left`、`right`－指定邊緣
- `top-left`、`top-right`、`bottom-left`、`bottom-right`－角落
- `none`－不允許調整大小

### 調整大小控點範例

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

## 平台特定行為

@tabs{sync-key="platform"}
[Windows]
**Windows 無框線視窗：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**功能：**

- 自動陰影
- 支援貼齊版面配置（Windows 11）
- 支援 Aero Snap
- DPI 縮放

**停用視窗裝飾：**

```go
Windows: application.WindowsWindow{
    DisableFramelessWindowDecorations: true,
},
```

**貼齊小幫手：**

```go
// Trigger Windows 11 Snap Assist
window.SnapAssist()
```

這會透過 Windows 快速鍵路徑觸發貼齊版面配置。若自訂 HTML 最大化按鈕需要具備原生暫留貼齊版面配置，請改用[Windows 上的原生非工作區](#windows-)。

**自訂標題列高度：** Windows 會自動從 CSS 偵測拖曳區域。

[macOS]
**macOS 無框線視窗：**

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

- 原生全螢幕支援
- 紅綠燈按鈕（選用）
- 毛玻璃效果
- 透明標題列

**完全隱藏標題列**（請使用`application`套件匯出的預設變體——不存在`TitleBarStyle`欄位或`MacTitleBarStyleHidden`常數）：

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

其他預設包括`MacTitleBarDefault`、`MacTitleBarHiddenInset`和`MacTitleBarHiddenInsetUnified`。

**不可見的標題列：** 隱藏標題列的同時仍可拖曳視窗。這僅在視窗為無框線或使用`AppearsTransparent`時生效：

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Linux 無框線視窗：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**功能：**

- 基本無框線支援
- CSS 拖曳區域
- 依桌面環境而異

**桌面環境注意事項：**

- <strong>GNOME：</strong>支援良好
- <strong>KDE Plasma：</strong>支援良好
- <strong>XFCE：</strong>基本支援
- <strong>平鋪式視窗管理員：</strong>支援有限

**需要合成器：** 透明效果需要合成器（大多數現代桌面環境均有提供）。

@end

## 常見模式

### 模式1：現代化標題列

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

### 模式2：啟動畫面

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

### 模式3：圓角視窗

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

### 模式4：浮層視窗

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

## 完整範例

以下是可用於正式環境的無框線視窗：

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

## 最佳實務

### ✅ 建議做法

- **提供可拖曳區域**——使用者需要移動視窗
- **實作系統按鈕**——關閉、最小化和最大化
- **設定最小尺寸**——避免版面無法使用
- **在所有平台上測試**——行為因平台而異
- **使用 CSS 定義拖曳區域**——靈活且易於維護
- **提供視覺回饋**——為按鈕設定游標懸停狀態

### ❌ 請勿

- **若視窗可調整大小，請勿忘記加入調整大小控點**
- **請勿讓整個視窗都可拖曳**——這會妨礙互動
- **請勿忘記將按鈕設為不可拖曳**——否則按鈕將無法運作
- **請勿使用過小的拖曳區域**——難以拖曳
- **請勿忽略平台差異**——請徹底測試

## 疑難排解

### 無法拖曳視窗

<strong>原因：</strong>缺少`--wails-draggable: drag`

**解決方法：**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### 按鈕無法運作

<strong>原因：</strong>按鈕位於可拖曳區域內

**解決方法：**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### 無法調整視窗大小

<strong>原因：</strong>缺少調整大小控點

**解決方法：**

```css
body {
    --wails-resize: all;
}
```

## 後續步驟

@cards{cols="2"}
▣ 視窗基礎
瞭解視窗管理的基礎知識。

[深入瞭解 →](/features/windows/basics/)

---
⚙ 視窗選項
視窗選項的完整參考資料。

[深入瞭解 →](/features/windows/options/)

---
🚀 視窗事件
處理視窗生命週期事件。

[深入瞭解 →](/features/windows/events/)

---
◆ 多視窗
多視窗應用程式的模式。

[深入瞭解 →](/features/windows/multiple/)

@end

---

<strong>有問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[無框線範例](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless)。
