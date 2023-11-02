# 应用程序

应用程序 API 用于使用 Wails 框架创建应用程序。

### New

API：`New(appOptions Options) *App`

`New(appOptions Options)` 使用给定的应用程序选项创建一个新的应用程序。它对未指定的选项应用默认值，将其与提供的选项合并，然后初始化并返回应用程序的实例。

如果在初始化过程中出现错误，应用程序将停止，并显示提供的错误消息。

需要注意的是，如果全局应用程序实例已经存在，将返回该实例，而不是创建新的实例。

```go title="main.go" hl_lines="6-9"
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name:        "WebviewWindow Demo",
		// 其他选项
    })

	// 其余的应用程序逻辑
}
```

### Get

`Get()` 返回全局应用程序实例。在代码的不同部分需要访问应用程序时非常有用。

```go
    // 获取应用程序实例
    app := application.Get()
```

### Capabilities

API：`Capabilities() capabilities.Capabilities`

`Capabilities()` 返回应用程序当前具有的功能的映射。这些功能可以是操作系统提供的不同功能，如 webview 功能。

```go
    // 获取应用程序的功能
    capabilities := app.Capabilities()
	if capabilities.HasNativeDrag {
		// 做一些事情
    }
```

### GetPID

API：`GetPID() int`

`GetPID()` 返回应用程序的进程 ID。

```go
    pid := app.GetPID()
```

### Run

API：`Run() error`

`Run()` 启动应用程序及其组件的执行。

```go
    app := application.New(application.Options{
	    // 选项
	})
    // 运行应用程序
    err := application.Run()
    if err != nil {
        // 处理错误
    }
```

### Quit

API：`Quit()`

`Quit()` 通过销毁窗口和可能的其他组件退出应用程序。

```go
    // 退出应用程序
    app.Quit()
```

### IsDarkMode

API：`IsDarkMode() bool`

`IsDarkMode()` 检查应用程序是否在暗模式下运行。它返回一个布尔值，指示是否启用了暗模式。

```go
    // 检查是否启用了暗模式
    if app.IsDarkMode() {
        // 做一些事情
    }
```

### Hide

API：`Hide()`

`Hide()` 隐藏应用程序窗口。

```go
    // 隐藏应用程序窗口
    app.Hide()
```

### Show

API：`Show()`

`Show()` 显示应用程序窗口。

```go
    // 显示应用程序窗口
    app.Show()
```

### NewWebviewWindow

API：`NewWebviewWindow() *WebviewWindow`

`NewWebviewWindow()` 使用默认选项创建一个新的 Webview 窗口，并返回它。

```go
    // 创建一个新的 Webview 窗口
    window := app.NewWebviewWindow()
```

### NewWebviewWindowWithOptions

API：`NewWebviewWindowWithOptions(windowOptions WebviewWindowOptions) *WebviewWindow`

`NewWebviewWindowWithOptions()` 使用自定义选项创建一个新的 Webview 窗口。新创建的窗口将添加到应用程序管理的窗口映射中。

```go
    // 使用自定义选项创建一个新的 Webview 窗口
    window := app.NewWebviewWindowWithOptions(WebviewWindowOptions{
		Name: "Main",
        Title: "My Window",
        Width: 800,
        Height: 600,
    })
```

### OnWindowCreation

API：`OnWindowCreation(callback func(window *WebviewWindow))`

`OnWindowCreation()` 注册一个回调函数，当创建窗口时调用该函数。

```go
    // 注册一个回调函数，当创建窗口时调用该函数
    app.OnWindowCreation(func(window *WebviewWindow) {
        // 做一些事情
    })
```

### GetWindowByName

API：`GetWindowByName(name string) *WebviewWindow`

`GetWindowByName()` 获取并返回具有特定名称的窗口。

```go
    // 通过名称获取窗口
    window := app.GetWindowByName("Main")
```

### CurrentWindow

API：`CurrentWindow() *WebviewWindow`

`CurrentWindow()` 获取并返回应用程序中当前活动窗口的指针。如果没有窗口，则返回 nil。

```go
    // 获取当前窗口
    window := app.CurrentWindow()
```

### RegisterContextMenu

API：`RegisterContextMenu(name string, menu *Menu)`

`RegisterContextMenu()` 注册具有给定名称的上下文菜单。稍后可以在应用程序中使用该菜单。

```go

    // 创建一个新的菜单
    ctxmenu := app.NewMenu()

    // 将菜单注册为上下文菜单
    app.RegisterContextMenu("MyContextMenu", ctxmenu)
```

### SetMenu

API：`SetMenu(menu *Menu)`

`SetMenu()` 设置应用程序的菜单。在 Mac 上，这将是全局菜单。对于 Windows 和 Linux，这将是任何新窗口的默认菜单。

```go
    // 创建一个新的菜单
    menu := app.NewMenu()

    // 设置应用程序的菜单
    app.SetMenu(menu)
```

### ShowAboutDialog

API：`ShowAboutDialog()`

`ShowAboutDialog()` 显示一个 "关于" 对话框。可以显示应用程序的名称、描述和图标。

```go
    // 显示关于对话框
    app.ShowAboutDialog()
```

### Info

API：`InfoDialog()`

`InfoDialog()` 创建并返回一个具有 `InfoDialogType` 的 `MessageDialog` 的新实例。此对话框通常用于向用户显示信息消息。

### Question

API：`QuestionDialog()`

`QuestionDialog()` 创建并返回一个具有 `QuestionDialogType` 的 `MessageDialog` 的新实例。此对话框通常用于向用户提问并期望回应。

### Warning

API：`WarningDialog()`

`WarningDialog()` 创建并返回一个具有 `WarningDialogType` 的 `MessageDialog` 的新实例。如其名称所示，此对话框主要用于向用户显示警告消息。

### Error

API：`ErrorDialog()`

`ErrorDialog()` 创建并返回一个具有 `ErrorDialogType` 的 `MessageDialog` 的新实例。此对话框设计用于在需要向用户显示错误消息时使用。

### OpenFile

API：`OpenFileDialog()`

`OpenFileDialog()` 创建并返回一个新的 `OpenFileDialogStruct`。此对话框提示用户从其文件系统中选择一个或多个文件。

### SaveFile

API：`SaveFileDialog()`

`SaveFileDialog()` 创建并返回一个新的 `SaveFileDialogStruct`。此对话框提示用户选择其文件系统上的位置以保存文件。

### OpenDirectory

API：`OpenDirectoryDialog()`

`OpenDirectoryDialog()` 创建并返回一个具有 `OpenDirectoryDialogType` 的 `MessageDialog` 的新实例。此对话框使用户能够从其文件系统中选择目录。

### On

API：`On(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`On()` 注册特定应用程序事件的事件侦听器。提供的回调函数将在相应事件发生时触发。该函数返回一个可调用的函数，用于删除侦听器。

### RegisterHook

API：`RegisterHook(eventType events.ApplicationEventType, callback func(event *Event)) func()`

`RegisterHook()` 注册要在特定事件期间作为钩子运行的回调函数。这些钩子在使用 `On()` 附加的侦听器之前运行。该函数返回一个可调用的函数，用于删除钩子。

### GetPrimaryScreen

API：`GetPrimaryScreen() (*Screen, error)`

`GetPrimaryScreen()` 返回系统的主屏幕。

### GetScreens

API：`GetScreens() ([]*Screen, error)`

`GetScreens()` 返回有关连接到系统的所有屏幕的信息。

这是提供的 `App` 结构中导出的方法的简要摘要。请注意，有关更详细的功能或注意事项，请参考实际的 Go 代码或进一步的内部文档。

## Options

```go title="application_options.go"
--8<--
../v3/pkg/application/options_application.go
--8<--
```

### Windows 选项

```go title="application_options_windows.go"
--8<--
../v3/pkg/application/options_application_win.go
--8<--
```

### Mac 选项

```go title="options_application_mac.go"
--8<--
../v3/pkg/application/options_application_mac.go
--8<--
```