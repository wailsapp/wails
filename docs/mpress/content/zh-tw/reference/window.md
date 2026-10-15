---
title: "視窗 API"
description: "視窗 API 完整參考文件"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## 概觀

視窗 API 提供控制視窗外觀、行為與生命週期的方法。可透過視窗執行個體或`app.Window`管理器存取。

`Window`是由`*application.WebviewWindow`實作的介面；以下方法簽章屬於`*WebviewWindow`。許多修改方法會傳回`Window`以便鏈式呼叫，各方法的傳回值均記載於其說明中。

**常見操作：**

- 建立並顯示視窗
- 控制大小、位置與狀態
- 處理視窗事件
- 管理視窗內容
- 設定外觀與行為

## 可見性

### Show()

顯示視窗。如果視窗原先隱藏，則會變為可見。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) Show() Window
```

**範例：**

```go
window := app.Window.New()
window.Show()
```

### Hide()

隱藏視窗，但不將其關閉。視窗會保留在記憶體中，之後可再次顯示。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) Hide() Window
```

**範例：**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**使用情境：**

- 隱藏至系統匣的應用程式
- 重複使用視窗的精靈流程
- 在作業期間暫時隱藏視窗

### Close()

關閉視窗。這會觸發`WindowClosing`事件。

```go
func (w *WebviewWindow) Close()
```

**範例：**

```go
window.Close()
```

<strong>注意：</strong>如果已註冊的鉤子呼叫`event.Cancel()`，便會阻止視窗關閉。

## 視窗屬性

### SetTitle()

設定視窗標題列文字。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**參數：**

- `title` - 新的視窗標題

**範例：**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

傳回視窗的唯一名稱識別碼。

```go
func (w *WebviewWindow) Name() string
```

**範例：**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## 大小與位置

### SetSize()

以像素為單位設定視窗尺寸。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**參數：**

- `width` - 視窗寬度（像素）
- `height` - 視窗高度（像素）

**範例：**

```go
window.SetSize(1024, 768)
```

### Size()

傳回目前的視窗尺寸。

```go
func (w *WebviewWindow) Size() (width, height int)
```

**範例：**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

設定視窗的最小與最大尺寸。兩者都會傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**範例：**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

設定視窗相對於螢幕左上角的位置。

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**參數：**

- `x` - 水平位置（像素）
- `y` - 垂直位置（像素）

**範例：**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

傳回目前的視窗位置。

```go
func (w *WebviewWindow) Position() (x, y int)
```

**範例：**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

將視窗置於螢幕中央。

```go
func (w *WebviewWindow) Center()
```

**範例：**

```go
window := app.Window.New()
window.Center()
window.Show()
```

<strong>注意：</strong>視窗會置於主要顯示器中央。若使用多顯示器設定，請參閱螢幕 API。

### Focus()

將視窗移至最前方，並使其取得鍵盤焦點。

```go
func (w *WebviewWindow) Focus()
```

**範例：**

```go
// Bring window to front
window.Focus()
```

## 視窗狀態

### Minimise() / UnMinimise()

將視窗最小化至工作列／Dock，或將其還原。`Minimise()`會傳回接收者以便鏈式呼叫；`UnMinimise()`不會傳回任何值。

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**範例：**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

將視窗最大化以填滿螢幕，或還原至先前的大小。`Maximise()`會傳回接收者以便鏈式呼叫；`UnMaximise()`不會傳回任何值。

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**範例：**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

進入或退出全螢幕模式。`Fullscreen()`會傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**範例：**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

沒有`SetFullscreen(bool)`方法。

### IsMinimised() / IsMaximised() / IsFullscreen()

檢查目前的視窗狀態。

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**範例：**

```go
if window.IsMinimised() {
    window.UnMinimise()
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    window.UnFullscreen()
}
```

## 視窗內容

### SetURL()

在視窗內導覽至指定的 URL。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**參數：**

- `url` - 要導覽至的 URL（內嵌資產可使用`http://wails.localhost/`）

**範例：**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

直接使用 HTML 字串設定視窗內容。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**參數：**

- `html` - 要顯示的 HTML 內容

**範例：**

```go
html := `
<!DOCTYPE html>
<html>
<head><title>Dynamic Content</title></head>
<body>
    <h1>Hello from Go!</h1>
    <p>This content was generated dynamically.</p>
</body>
</html>
`
window.SetHTML(html)
```

**使用情境：**

- 動態產生內容
- 不使用前端建置流程的簡易視窗
- 錯誤頁面或啟動畫面

### Reload()

重新載入目前的視窗內容。

```go
func (w *WebviewWindow) Reload()
```

**範例：**

```go
// Reload current page
window.Reload()
```

<strong>注意：</strong>適合在開發期間或需要重新整理內容時使用。

## 視窗事件

Wails 提供兩種處理視窗事件的方法：

- **OnWindowEvent()** - 監聽視窗事件（無法阻止事件）。
- **RegisterHook()** - 攔截視窗事件（可呼叫`event.Cancel()`阻止事件）。

### OnWindowEvent()

為視窗事件註冊回呼函式。傳回取消訂閱函式。

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**範例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

// Listen for window lost focus
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window lost focus")
})

// Listen for window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    app.Logger.Info("Window resized")
})
```

**常見視窗事件：**

- `events.Common.WindowClosing` - 視窗即將關閉
- `events.Common.WindowFocus` - 視窗取得焦點
- `events.Common.WindowLostFocus` - 視窗失去焦點
- `events.Common.WindowDidMove` - 視窗已移動
- `events.Common.WindowDidResize` - 視窗大小已調整
- `events.Common.WindowMinimise` - 視窗已最小化
- `events.Common.WindowMaximise` - 視窗已最大化
- `events.Common.WindowFullscreen` - 視窗已進入全螢幕模式
- `events.Common.WindowRuntimeReady` - 視窗內執行階段已初始化

### RegisterHook()

為視窗事件註冊鉤子。鉤子會在監聽器之前執行，並可呼叫`event.Cancel()`阻止事件。傳回取消訂閱函式。

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**範例 - 阻止視窗關閉：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    confirm := app.Dialog.Question().
        SetTitle("Confirm Close").
        SetMessage("Are you sure you want to close this window?")

    yes := confirm.AddButton("Yes")
    no := confirm.AddButton("No")
    confirm.SetDefaultButton(yes)
    confirm.SetCancelButton(no)

    no.OnClick(func() {
        e.Cancel() // Prevent window from closing
    })

    confirm.Show()
})
```

**範例 - 關閉前儲存：**

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges {
        return
    }

    dlg := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save changes before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Don't Save")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveData() })
    cancel.OnClick(func() { e.Cancel() })
    _ = discard // allow close

    dlg.Show()
})
```

### EmitEvent()

向視窗的前端發出自訂事件。如果發出動作被鉤子取消，則傳回`true`。

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**參數：**

- `name` - 事件名稱
- `data` - 要隨事件傳送的選用資料

**範例：**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**前端（JavaScript）：**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## 其他方法

### SetEnabled()

啟用或停用使用者與視窗的互動。

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**範例：**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

設定視窗的背景色彩（在內容載入前顯示）。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA`是`application.RGBA{Red, Green, Blue, Alpha uint8}`。請使用輔助函式`application.NewRGB(r, g, b)`（Alpha 值為255）或`application.NewRGBA(r, g, b, a)`。

**範例：**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

控制使用者是否可調整視窗大小。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**範例：**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

設定視窗是否保持在其他視窗之上。傳回接收者以便鏈式呼叫。

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**範例：**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

開啟視窗內容的原生列印對話方塊。

```go
func (w *WebviewWindow) Print() error
```

<strong>傳回值：</strong>列印失敗時傳回錯誤。

**範例：**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

將第二個視窗附加為附屬於父視窗的模態對話框。

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**參數：**

- `modalWindow`－要附加為模態視窗的視窗

**平台支援：**

- **macOS**：完整支援（以附屬於父視窗的對話框形式顯示）
- **Windows**：不支援
- **Linux**：不支援

**範例：**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## 平台特定選項

### Linux

Linux 視窗可透過`LinuxWindow`使用下列平台特定選項：

#### MenuStyle

控制應用程式選單的顯示方式。此選項可用於預設的 GTK4 組建，但在舊版`-tags gtk3`組建中會被忽略。

| 值 | 說明 |
| --- | --- |
| `LinuxMenuStyleMenuBar` | 標題列下方的傳統選單列（預設） |
| `LinuxMenuStylePrimaryMenu` | 標頭列中的主要選單按鈕（GNOME 風格） |

**範例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

<strong>注意：</strong>主要選單樣式會依循 GNOME 人機介面指南，在標頭列中顯示漢堡選單按鈕（☰）。這是現代 GNOME 應用程式的建議樣式。

## 完整範例

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "Window API Demo",
    })

    // Create window with options
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My Application",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    // Configure window behaviour
    window.SetResizable(true)
    window.SetMinSize(800, 600)
    window.SetMaxSize(1920, 1080)

    // Confirm-before-close hook
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        dlg := app.Dialog.Question().
            SetTitle("Confirm Close").
            SetMessage("Are you sure you want to close this window?")

        yes := dlg.AddButton("Yes")
        no := dlg.AddButton("No")
        dlg.SetDefaultButton(yes)
        dlg.SetCancelButton(no)
        no.OnClick(func() { e.Cancel() })

        dlg.Show()
    })

    // Listen for window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application (Active)")
        app.Logger.Info("Window gained focus")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application")
        app.Logger.Info("Window lost focus")
    })

    // Position and show window
    window.Center()
    window.Show()

    app.Run()
}
```
