# 系统托盘

系统托盘位于桌面环境的通知区域，可以包含当前运行应用程序的图标和特定系统通知。

您可以通过调用 `app.NewSystemTray()` 来创建一个系统托盘：

```go
    // 创建一个新的系统托盘
tray := app.NewSystemTray()
```

`SystemTray` 类型上提供了以下方法：

### SetLabel

API：`SetLabel(label string)`

`SetLabel` 方法设置托盘的标签。

### Label

API：`Label() string`

`Label` 方法获取托盘的标签。

### PositionWindow

API：`PositionWindow(*WebviewWindow, offset int) error`

`PositionWindow` 方法调用了 `AttachWindow` 和 `WindowOffset` 方法。

### SetIcon

API：`SetIcon(icon []byte) *SystemTray`

`SetIcon` 方法设置系统托盘的图标。

### SetDarkModeIcon

API：`SetDarkModeIcon(icon []byte) *SystemTray`

`SetDarkModeIcon` 方法设置暗黑模式下系统托盘的图标。

### SetMenu

API：`SetMenu(menu *Menu) *SystemTray`

`SetMenu` 方法设置系统托盘的菜单。

### Destroy

API：`Destroy()`

`Destroy` 方法销毁系统托盘实例。

### OnClick

API：`OnClick(handler func()) *SystemTray`

`OnClick` 方法设置点击托盘图标时执行的函数。

### OnRightClick

API：`OnRightClick(handler func()) *SystemTray`

`OnRightClick` 方法设置右键点击托盘图标时执行的函数。

### OnDoubleClick

API：`OnDoubleClick(handler func()) *SystemTray`

`OnDoubleClick` 方法设置双击托盘图标时执行的函数。

### OnRightDoubleClick

API：`OnRightDoubleClick(handler func()) *SystemTray`

`OnRightDoubleClick` 方法设置右键双击托盘图标时执行的函数。

### AttachWindow

API：`AttachWindow(window *WebviewWindow) *SystemTray`

`AttachWindow` 方法将窗口附加到系统托盘。当点击系统托盘图标时，窗口将显示出来。

### WindowOffset

API：`WindowOffset(offset int) *SystemTray`

`WindowOffset` 方法设置系统托盘与窗口之间的像素间隔。

### WindowDebounce

API：`WindowDebounce(debounce time.Duration) *SystemTray`

`WindowDebounce` 方法设置防抖时间。在 Windows 上，它用于指定在响应通知图标上的鼠
标松开事件之前等待多长时间。

### OpenMenu

API：`OpenMenu()`

`OpenMenu` 方法打开与系统托盘关联的菜单。
