---
title: "在 Wails 中自訂視窗"
description: "自訂 Wails 應用程式中的視窗外觀與行為"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

適用平台：<span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails 提供 API，可控制視窗控制項的外觀與功能。這項功能適用於 Windows 和 macOS，但不適用於 Linux。

## 設定視窗按鈕狀態

按鈕狀態由 `ButtonState` 列舉定義：

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`：按鈕已啟用且可見。
- `ButtonDisabled`：按鈕可見但已停用（呈灰色）。
- `ButtonHidden`：按鈕在標題列中隱藏。

按鈕狀態可在建立視窗時或執行階段設定。

### 建立視窗時設定按鈕狀態

建立新視窗時，可使用 `WebviewWindowOptions` 結構設定按鈕的初始狀態：

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

在上述範例中，最小化按鈕已隱藏，最大化按鈕未啟用（呈灰色），而關閉按鈕已啟用。

### 在執行階段設定按鈕狀態

也可以在執行階段使用 `Window` 介面上的下列方法變更按鈕狀態：

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS：MaximiseButtonState 和 FullscreenButtonState 共用一個按鈕

在 macOS 上，綠色交通號誌按鈕（`NSWindowZoomButton`）是最大化與全螢幕功能共用的同一個實體控制項。若在建立視窗時將 `MaximiseButtonState` 和 `FullscreenButtonState` 設為不同的值，原本會發生無提示的「最後寫入者優先」覆寫。

為避免這種情況，Wails 會在初始化時套用兩種狀態中<strong>限制較嚴格</strong>的狀態，其順序為 `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden`。

| `MaximiseButtonState` | `FullscreenButtonState` | macOS 上的實際狀態 |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

在執行階段，`SetMaximiseButtonState` 和 `SetFullscreenButtonState` 在 macOS 上都以 `NSWindowZoomButton` 為目標，因此以最後一次呼叫為準。

### 平台差異

按鈕狀態功能在 Windows 和 macOS 上的行為略有不同：

|  | Windows | Mac |
| --- | --- | --- |
| 停用最小化／最大化／關閉按鈕 | 停用最小化／最大化／關閉按鈕 | 停用最小化／最大化／關閉按鈕 |
| 隱藏最小化按鈕 | 停用最小化按鈕 | 隱藏最小化按鈕 |
| 隱藏最大化按鈕 | 停用最大化按鈕 | 隱藏最大化按鈕 |
| 隱藏關閉按鈕 | 隱藏所有控制項 | 隱藏關閉按鈕 |
| `FullscreenButtonState` | 不執行任何操作 | 以縮放（綠色）按鈕為目標 |

注意：在 Windows 上，無法個別隱藏最小化／最大化按鈕。 不過，同時停用兩者會隱藏這兩個控制項，只顯示關閉按鈕。Windows 的標準標題列沒有專用的全螢幕按鈕，因此`FullscreenButtonState`在該平台上不會有任何作用。

### 控制視窗樣式（Windows）

若要控制 Windows 上的標題列樣式，可以使用`WebviewWindowOptions`結構中的`ExStyle`欄位：

範例：

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

這項設定會覆寫其他影響視窗擴充樣式的選項：

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
