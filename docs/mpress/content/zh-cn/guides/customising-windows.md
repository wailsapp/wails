---
title: "在 Wails 中自定义窗口"
description: "自定义 Wails 应用程序中的窗口外观和行为"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

相关平台：<span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails 提供了用于控制窗口控件外观和功能的 API。此功能可用于 Windows 和 macOS，但不支持 Linux。

## 设置窗口按钮状态

按钮状态由 `ButtonState` 枚举定义：

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`：按钮已启用且可见。
- `ButtonDisabled`：按钮可见但已禁用（显示为灰色）。
- `ButtonHidden`：按钮在标题栏中隐藏。

可以在创建窗口时或运行时设置按钮状态。

### 创建窗口时设置按钮状态

创建新窗口时，可以使用 `WebviewWindowOptions` 结构体设置按钮的初始状态：

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

在上面的示例中，最小化按钮处于隐藏状态，最大化按钮处于非活动状态（显示为灰色），关闭按钮处于活动状态。

### 在运行时设置按钮状态

还可以在运行时使用 `Window` 接口上的以下方法更改按钮状态：

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS：MaximiseButtonState 和 FullscreenButtonState 共用一个按钮

在macOS上，绿色交通灯按钮（`NSWindowZoomButton`）是最大化和全屏功能共用的同一个物理控件。在创建窗口时将`MaximiseButtonState`和`FullscreenButtonState`设置为不同的值，原本会导致后写入的值覆盖先写入的值，且不作任何提示。

为避免这种情况，Wails 会在初始化时采用这两个状态中<strong>限制更严格</strong>的一个，严格程度依次为 `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden`。

| `MaximiseButtonState` | `FullscreenButtonState` | 在 macOS 上的实际状态 |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

在运行时，`SetMaximiseButtonState` 和 `SetFullscreenButtonState` 在 macOS 上都以 `NSWindowZoomButton` 为目标，因此最后一次调用生效。

### 平台差异

按钮状态功能在 Windows 和 macOS 上的行为略有不同：

|  | Windows | Mac |
| --- | --- | --- |
| 禁用最小化/最大化/关闭按钮 | 禁用最小化/最大化/关闭按钮 | 禁用最小化/最大化/关闭按钮 |
| 隐藏最小化按钮 | 禁用最小化按钮 | 隐藏最小化按钮 |
| 隐藏最大化按钮 | 禁用最大化按钮 | 隐藏最大化按钮 |
| 隐藏关闭按钮 | 隐藏所有控件 | 隐藏关闭按钮 |
| `FullscreenButtonState` | 无操作 | 针对缩放（绿色）按钮 |

注意：在 Windows 上，无法单独隐藏最小化/最大化按钮。 但是，同时禁用这两个按钮会将它们全部隐藏，仅显示关闭按钮。Windows 的标准标题栏中没有专用的全屏按钮，因此 `FullscreenButtonState` 在该平台上不起作用。

### 控制窗口样式（Windows）

要控制 Windows 上的标题栏样式，可以使用 `WebviewWindowOptions` 结构体中的 `ExStyle` 字段：

示例：

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

此设置将覆盖其他影响窗口扩展样式的选项：

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
