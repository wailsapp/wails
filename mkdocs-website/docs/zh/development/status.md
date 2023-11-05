将以下文本翻译为中文，并不要翻译 `!!! note` 或 `!!! tip` 或以此格式的字符串：

# 状态

v3版功能的状态。

!!! note

        此列表包含公有和内部API支持的混合内容。<br/>
        它不完整且可能不是最新的。

## 已知问题

- Linux尚未与Windows/Mac达到功能平衡

## 应用程序

应用程序接口方法

| 方法                                                          | Windows | Linux | Mac | 备注 |
| ------------------------------------------------------------- | ------- | ----- | --- | ---- |
| run() error                                                   | Y       | Y     | Y   |      |
| destroy()                                                     |         | Y     | Y   |      |
| setApplicationMenu(menu \*Menu)                               | Y       | Y     | Y   |      |
| name() string                                                 |         | Y     | Y   |      |
| getCurrentWindowID() uint                                     | Y       | Y     | Y   |      |
| showAboutDialog(name string, description string, icon []byte) |         | Y     | Y   |      |
| setIcon(icon []byte)                                          | -       | Y     | Y   |      |
| on(id uint)                                                   |         |       | Y   |      |
| dispatchOnMainThread(fn func())                               | Y       | Y     | Y   |      |
| hide()                                                        | Y       | Y     | Y   |      |
| show()                                                        | Y       | Y     | Y   |      |
| getPrimaryScreen() (\*Screen, error)                          |         | Y     | Y   |      |
| getScreens() ([]\*Screen, error)                              |         | Y     | Y   |      |

## Webview 窗口

Webview 窗口接口方法

| 方法                                               | Windows | Linux | Mac | 备注                 |
| -------------------------------------------------- | ------- | ----- | --- | -------------------- |
| center()                                           | Y       | Y     | Y   |                      |
| close()                                            | y       | Y     | Y   |                      |
| destroy()                                          |         | Y     | Y   |                      |
| execJS(js string)                                  | y       | Y     | Y   |                      |
| focus()                                            | Y       | Y     |     |                      |
| forceReload()                                      |         | Y     | Y   |                      |
| fullscreen()                                       | Y       | Y     | Y   |                      |
| getScreen() (\*Screen, error)                      | y       | Y     | Y   |                      |
| getZoom() float64                                  |         | Y     | Y   |                      |
| height() int                                       | Y       | Y     | Y   |                      |
| hide()                                             | Y       | Y     | Y   |                      |
| isFullscreen() bool                                | Y       | Y     | Y   |                      |
| isMaximised() bool                                 | Y       | Y     | Y   |                      |
| isMinimised() bool                                 | Y       | Y     | Y   |                      |
| maximise()                                         | Y       | Y     | Y   |                      |
| minimise()                                         | Y       | Y     | Y   |                      |
| nativeWindowHandle() (uintptr, error)              | Y       | Y     |     |                      |
| on(eventID uint)                                   | y       |       | Y   |                      |
| openContextMenu(menu *Menu, data *ContextMenuData) | y       |       | Y   |                      |
| relativePosition() (int, int)                      | Y       | Y     | Y   |                      |
| reload()                                           | y       | Y     | Y   |                      |
| run()                                              | Y       | Y     | Y   |                      |
| setAlwaysOnTop(alwaysOnTop bool)                   | Y       | Y     | Y   |                      |
| setBackgroundColour(color RGBA)                    | Y       | Y     | Y   |                      |
| setEnabled(bool)                                   |         | Y     | Y   |                      |
| setFrameless(bool)                                 |         | Y     | Y   |                      |
| setFullscreenButtonEnabled(enabled bool)           | -       | Y     | Y   | Windows 没有全屏按钮 |
| setHTML(html string)                               | Y       | Y     | Y   |                      |
| setMaxSize(width, height int)                      | Y       | Y     | Y   |                      |
| setMinSize(width, height int)                      | Y       | Y     | Y   |                      |
| setRelativePosition(x int, y int)                  | Y       | Y     | Y   |                      |
| setResizable(resizable bool)                       | Y       | Y     | Y   |                      |
| setSize(width, height int)                         | Y       | Y     | Y   |                      |
| setTitle(title string)                             | Y       | Y     | Y   |                      |
| setURL(url string)                                 | Y       | Y     | Y   |                      |
| setZoom(zoom float64)                              | Y       | Y     | Y   |                      |
| show()                                             | Y       | Y     | Y   |                      |
| size() (int, int)                                  | Y       | Y     | Y   |                      |
| toggleDevTools()                                   | Y       | Y     | Y   |                      |
| unfullscreen()                                     | Y       | Y     | Y   |                      |
| unmaximise()                                       | Y       | Y     | Y   |                      |
| unminimise()                                       | Y       | Y     | Y   |                      |
| width() int                                        | Y       | Y     | Y   |                      |
| zoom()                                             |         | Y     | Y   |                      |
| zoomIn()                                           | Y       | Y     | Y   |                      |
| zoomOut()                                          | Y       | Y     | Y   |                      |
| zoomReset()                                        | Y       | Y     | Y   |                      |

## 运行时

### 应用程序

| 功能 | Windows | Linux | Mac | 备注 |
| ---- | ------- | ----- | --- | ---- |
| 退出 | Y       | Y     | Y   |      |
| 隐藏 | Y       |       | Y   |      |
| 显示 | Y       |       | Y   |      |

### 对话框

| 功能     | Windows | Linux | Mac | 备注 |
| -------- | ------- | ----- | --- | ---- |
| 信息     | Y       | Y     | Y   |      |
| 警告     | Y       | Y     | Y   |      |
| 错误     | Y       | Y     | Y   |      |
| 问题     | Y       | Y     | Y   |      |
| 打开文件 | Y       |       | Y   |      |
| 保存文件 | Y       |       | Y   |      |

### 剪贴板

| 功能     | Windows | Linux | Mac | 备注 |
| -------- | ------- | ----- | --- | ---- |
| 设置文本 | Y       |       | Y   |      |
| 文本     | Y       |       | Y   |      |

### 上下文菜单

| 功能           | Windows | Linux | Mac | 备注 |
| -------------- | ------- | ----- | --- | ---- |
| 打开上下文菜单 | Y       |       | Y   |      |
| 默认开启       |         |       |     |      |
| 通过 HTML 控制 | Y       |       |     |      |

默认上下文菜单默认对所有`contentEditable: true`、`<input>`或`<textarea>`标签的元
素或具有`--default-contextmenu: true`样式的元素启
用。`--default-contextmenu: show`样式将始终显示上下文菜
单。`--default-contextmenu: hide`样式将始终隐藏上下文菜单。

嵌套在带有`--default-contextmenu: hide`样式的标签下的任何内容，除非使
用`--default-contextmenu: show`进行显式设置，否则不会显示上下文菜单。

### 屏幕

| 功能         | Windows | Linux | Mac | 备注 |
| ------------ | ------- | ----- | --- | ---- |
| 获取所有     | Y       | Y     | Y   |      |
| 获取主屏幕   | Y       | Y     | Y   |      |
| 获取当前屏幕 | Y       | Y     | Y   |      |

### 系统

| 功能         | Windows | Linux | Mac | 备注 |
| ------------ | ------- | ----- | --- | ---- |
| 是否为暗模式 |         |       | Y   |      |

### 窗口

Y = 支持 U = 未经测试

- = 不可用

| 功能               | Windows | Linux | Mac | 备注                                                                                 |
| ------------------ | ------- | ----- | --- | ------------------------------------------------------------------------------------ |
| 居中               | Y       | Y     | Y   |                                                                                      |
| 获得焦点           | Y       | Y     |     |                                                                                      |
| 全屏               | Y       | Y     | Y   |                                                                                      |
| 获得缩放比例       | Y       | Y     | Y   | 获取当前视图比例                                                                     |
| 高度               | Y       | Y     | Y   |                                                                                      |
| 隐藏               | Y       | Y     | Y   |                                                                                      |
| 最大化             | Y       | Y     | Y   |                                                                                      |
| 最小化             | Y       | Y     | Y   |                                                                                      |
| 相对位置           | Y       | Y     | Y   |                                                                                      |
| 屏幕               | Y       | Y     | Y   | 获取窗口的屏幕                                                                       |
| 设置始终在顶部     | Y       | Y     | Y   |                                                                                      |
| 设置背景颜色       | Y       | Y     | Y   | https://github.com/MicrosoftEdge/WebView2Feedback/issues/1621#issuecomment-938234294 |
| 设置启用状态       | Y       | U     | -   | 设置窗口是否可用                                                                     |
| 设置最大尺寸       | Y       | Y     | Y   |                                                                                      |
| 设置最小尺寸       | Y       | Y     | Y   |                                                                                      |
| 设置相对位置       | Y       | Y     | Y   |                                                                                      |
| 设置是否可调整大小 | Y       | Y     | Y   |                                                                                      |
| 设置大小           | Y       | Y     | Y   |                                                                                      |
| 设置标题           | Y       | Y     | Y   |                                                                                      |
| 设置缩放比例       | Y       | Y     | Y   | 设置视图比例                                                                         |
| 显示               | Y       | Y     | Y   |                                                                                      |
| 尺寸               | Y       | Y     | Y   |                                                                                      |
| 取消全屏           | Y       | Y     | Y   |                                                                                      |
| 取消最大化         | Y       | Y     | Y   |                                                                                      |
| 取消最小化         | Y       | Y     | Y   |                                                                                      |
| 宽度               | Y       | Y     | Y   |                                                                                      |
| 缩放               |         | Y     | Y   |                                                                                      |
| 放大               | Y       | Y     | Y   | 增加视图比例                                                                         |
| 缩小               | Y       | Y     | Y   | 减小视图比例                                                                         |
| 重置缩放           | Y       | Y     | Y   | 重置视图比例                                                                         |

### 窗口选项

下表中的'Y'表示已经测试并且在窗口创建时应用了该选项。'X'表示该平台不支持该选项。

| 功能             | Windows | Linux | Mac | 备注                                 |
| ---------------- | ------- | ----- | --- | ------------------------------------ |
| 始终在顶部       | Y       |       |     |                                      |
| 背景颜色         | Y       | Y     |     |                                      |
| 背景类型         |         |       |     | 默认情况下，亚克力效果有效，其他无效 |
| CSS              | Y       | Y     |     |                                      |
| DevToolsEnabled  | Y       | Y     | Y   |                                      |
| DisableResize    | Y       | Y     |     |                                      |
| 启用拖放         |         | Y     |     |                                      |
| 启用欺诈网站警告 |         |       |     |                                      |
| 获得焦点         | Y       | Y     |     |                                      |
| 无边框           | Y       | Y     |     |                                      |
| 启用全屏按钮     | Y       |       |     | Windows上没有全屏按钮                |
| HTML             | Y       | Y     |     |                                      |
| JS               | Y       | Y     |     |                                      |
| Mac              | -       | -     |     |                                      |
| 最大高度         | Y       | Y     |     |                                      |
| 最大宽度         | Y       | Y     |     |                                      |
| 最小高度         | Y       | Y     |     |                                      |
| 最小宽度         | Y       | Y     |     |                                      |
| 名称             | Y       | Y     |     |                                      |
| 启动时打开检查器 |         |       |     |                                      |
| 启动状态         | Y       |       |     |                                      |
| 标题             | Y       | Y     |     |                                      |
| URL              | Y       | Y     |     |                                      |
| 宽度             | Y       | Y     |     |                                      |
| Windows          | Y       | -     | -   |                                      |
| X                | Y       | Y     |     |                                      |
| Y                | Y       | Y     |     |                                      |
| 缩放             |         |       |     |                                      |
| 启用缩放控件     |         |       |     |                                      |

### 日志

要记录还是不要记录？系统日志器与自定义日志器。

## 菜单

| 事件             | Windows | Linux | Mac | 备注 |
| ---------------- | ------- | ----- | --- | ---- |
| 默认应用程序菜单 | Y       | Y     | Y   |      |

## 托盘菜单

| 功能           | Windows | Linux | Mac | 备注                                               |
| -------------- | ------- | ----- | --- | -------------------------------------------------- |
| 图标           | Y       |       | Y   | Windows具有默认的浅色/深色模式图标并支持PNG或ICO。 |
| 标签           | -       |       | Y   |                                                    |
| 标签（ANSI码） | -       |       |     |                                                    |
| 菜单           | Y       |       | Y   |                                                    |

### 方法

| 方法                          | Windows | Linux | Mac | 备注 |
| ----------------------------- | ------- | ----- | --- | ---- |
| setLabel(label string)        | -       |       | Y   |      |
| run()                         | Y       |       | Y   |      |
| setIcon(icon []byte)          | Y       |       | Y   |      |
| setMenu(menu \*Menu)          | Y       |       | Y   |      |
| setIconPosition(position int) | -       |       | Y   |      |
| setTemplateIcon(icon []byte)  | -       |       | Y   |      |
| destroy()                     | Y       |       | Y   |      |
| setDarkModeIcon(icon []byte)  | Y       |       | Y   |      |

## 跨平台事件

将本机事件映射到跨平台事件。

| 事件                     | Windows | Linux | Mac             | 备注 |
| ------------------------ | ------- | ----- | --------------- | ---- |
| WindowWillClose          |         |       | WindowWillClose |      |
| WindowDidClose           |         |       |                 |      |
| WindowDidResize          |         |       |                 |      |
| WindowDidHide            |         |       |                 |      |
| ApplicationWillTerminate |         |       |                 |      |

... 添加更多

## 绑定生成

工作正常。

## 模型生成

工作正常。

## 任务文件

包含很多开发所需的内容。

## 主题

| 模式 | Windows | Linux | Mac | 备注 |
| ---- | ------- | ----- | --- | ---- |
| 暗   | Y       |       |     |      |
| 亮   | Y       |       |     |      |
| 系统 | Y       |       |     |      |

## NSIS安装程序

待定

## 模板

所有模板都可用。

## 插件

内置插件支持：

| 插件       | Windows | Linux | Mac | 备注 |
| ---------- | ------- | ----- | --- | ---- |
| 浏览器     | Y       |       | Y   |      |
| KV 存储    | Y       | Y     | Y   |      |
| 日志       | Y       | Y     | Y   |      |
| 单实例     | Y       |       | Y   |      |
| SQLite     | Y       | Y     | Y   |      |
| 开机自启动 | Y       |       | Y   |      |
| 服务器     |         |       |     |      |

待办事项：

- 确保每个插件都有一个可以注入到窗口中的JS包装器。

## 打包

|                    | Windows | Linux | Mac | 备注 |
| ------------------ | ------- | ----- | --- | ---- |
| 图标生成           | Y       |       | Y   |      |
| 图标嵌入           | Y       |       | Y   |      |
| Info.plist         | -       |       | Y   |      |
| NSIS 安装程序      |         |       | -   |      |
| Mac 包             | -       |       | Y   |      |
| Windows 可执行文件 | Y       |       | -   |      |

## 无边框窗口

| 功能     | Windows | Linux | Mac | 备注                               |
| -------- | ------- | ----- | --- | ---------------------------------- |
| 调整大小 | Y       |       | Y   |                                    |
| 拖拽     | Y       | Y     | Y   | Linux-始终可以使用 `Meta`+左键拖拽 |

## Mac 特定

- [x] 半透明

### Mac 选项

| 功能           | 默认值            | 备注                               |
| -------------- | ----------------- | ---------------------------------- |
| 背景           | MacBackdropNormal | 标准的实心窗口                     |
| 禁用阴影       | false             |                                    |
| 标题栏         |                   | 默认情况下使用标准的窗口装饰       |
| 外观           | DefaultAppearance |                                    |
| 隐藏标题栏高度 | 0                 | 为无边框窗口创建一个不可见的标题栏 |
| 禁用阴影       | false             | 禁用窗口投影阴影                   |

## Windows 特定

- [x] 半透明
- [x] 自定义主题

### Windows 选项

| 功能               | 默认值        | 备注                 |
| ------------------ | ------------- | -------------------- |
| 背景类型           | Solid         |                      |
| 禁用图标           | false         |                      |
| 主题               | SystemDefault |                      |
| 自定义主题         | nil           |                      |
| 禁用无边框窗口装饰 | false         |                      |
| 窗口遮罩           | nil           | 使窗口成为位图的内容 |

## Linux 特定

由`*_linux.go`文件使用的函数的实现详细信息位于以下文件中：

- linux_cgo.go：CGo 实现
- linux_purego.go：PureGo 实现

### CGO

默认情况下，使用 CGO 编译 Linux 端口。这会阻止轻松的交叉编译，因此同时也正在同时
开发 PureGo 实现。

### Purego

可以使用以下命令编译示例：

    CGO_ENABLED=0 go build -tags purego

注意：重构之后的功能目前无法正常工作。
