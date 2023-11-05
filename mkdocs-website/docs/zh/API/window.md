# 窗口

要创建一个窗口，可以使
用[Application.NewWebviewWindow](application.md#newwebviewwindow)或[Application.NewWebviewWindowWithOptions](application.md#newwebviewwindowwithoptions)。
前者创建一个具有默认选项的窗口，而后者允许您指定自定义选项。

这些方法可在返回的WebviewWindow对象上调用：

### SetTitle

API: `SetTitle(title string) *WebviewWindow`

此方法将窗口标题更新为提供的字符串。它返回WebviewWindow对象，允许进行方法链接。

### Name

API: `Name() string`

此函数返回WebviewWindow的名称。

### SetSize

API: `SetSize(width, height int) *WebviewWindow`

此方法将WebviewWindow的大小设置为提供的宽度和高度参数。如果提供的尺寸超过约束条
件，它们将被相应调整。

### SetAlwaysOnTop

API: `SetAlwaysOnTop(b bool) *WebviewWindow`

此函数根据提供的布尔标志设置窗口始终置顶。

### Show

API: `Show() *WebviewWindow`

`Show`方法用于使窗口可见。如果窗口未运行，它首先调用`run`方法启动窗口，然后使其
可见。

### Hide

API: `Hide() *WebviewWindow`

`Hide`方法用于隐藏窗口。它将窗口的隐藏状态设置为true，并触发窗口隐藏事件。

### SetURL

API: `SetURL(s string) *WebviewWindow`

`SetURL`方法用于将窗口的URL设置为给定的URL字符串。

### SetZoom

API: `SetZoom(magnification float64) *WebviewWindow`

`SetZoom`方法将窗口内容的缩放级别设置为提供的放大倍数。

### GetZoom

API: `GetZoom() float64`

`GetZoom`函数返回窗口内容的当前缩放级别。

### GetScreen

API: `GetScreen() (*Screen, error)`

`GetScreen`方法返回窗口所显示的屏幕。

#### SetFrameless

API: `SetFrameless(frameless bool) *WebviewWindow`

此函数用于移除窗口边框和标题栏。它根据提供的布尔值（true表示无边框，false表示有
边框）切换窗口的无边框状态。

#### RegisterContextMenu

API: `RegisterContextMenu(name string, menu *Menu)`

此函数用于注册上下文菜单并为其指定给定的名称。

#### NativeWindowHandle

API: `NativeWindowHandle() (uintptr, error)`

此函数用于获取窗口的平台本机窗口句柄。

#### Focus

API: `Focus()`

此函数用于将焦点设置到窗口。

#### SetEnabled

API: `SetEnabled(enabled bool)`

此函数用于根据提供的布尔值启用/禁用窗口。

#### SetAbsolutePosition

API: `SetAbsolutePosition(x int, y int)`

此函数设置窗口在屏幕上的绝对位置。
