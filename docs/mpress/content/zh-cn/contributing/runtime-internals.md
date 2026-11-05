---
title: "运行时内部机制"
description: "深入了解 Wails v3 如何启动、运行以及与操作系统通信"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

<strong>运行时</strong>是将普通 Go 函数转变为跨平台桌面应用程序的层。 本文档介绍跟踪源代码时会遇到的各个组成部分。

---

## 1. 应用程序生命周期

| 阶段 | 代码路径 | 执行内容 |
| --- | --- | --- |
| **引导启动** | `pkg/application/application.go:init()` | 注册构建时数据，并创建一个全局`application`单例。 |
| **New()** | `application.New(...)` | 验证`Options`，启动<strong>AssetServer</strong>，并初始化日志记录。 |
| **Run()** | `application.(*App).Run()` | 1。调用平台的`mainthread.X()`以进入操作系统 UI 线程。<br />2。启动<strong>运行时</strong>（`internal/runtime`）。<br />3。阻塞，直到最后一个窗口关闭或`Quit()`被调用。 |
| **关闭** | `application.(*App).Quit()` | 广播`application:shutdown`事件，刷新日志，并销毁窗口和服务。 |

生命周期严格遵循<strong>单次进入</strong>原则：可以创建多个窗口，但应用程序对象本身只初始化一次。

---

## 2。窗口管理

### 公共 API

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()`**不接受任何参数**；请在以下情况下使用`NewWithOptions(...)`：
>
> 需要按值传递一个`application.WebviewWindowOptions`结构体。

`app.Window.New[WithOptions]()` 委托给包含平台特定实现的 `pkg/application/webview_window_*.go`：

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

每个文件：

1. 创建原生 WebView（WKWebView、WebKitGTK、WebView2）。
2. 注册<strong>消息处理器</strong>回调（`pkg/application/messageprocessor*.go`）。
3. 将 Wails 事件（`WindowDidResize`、`WindowFocus`、`WindowFilesDropped`……）映射到`pkg/events`中的常量。

`internal/runtime/`保留用于少量构建标签粘合代码 （`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`）以及`internal/runtime/desktop/`下嵌入的 JS 运行时。

活动窗口由`pkg/application/window_manager.go`/ `webview_window.go`跟踪。`pkg/application/screenmanager.go`用于存储<strong>显示器</strong> 元数据（分辨率、缩放比例、工作区），不负责管理窗口。

---

## 3. 消息处理管线

JavaScript 与 Go 之间的桥接由`pkg/application/messageprocessor_*.go`中的<strong>消息 处理器</strong>系列实现。

流程：

1. <strong>JavaScript</strong>从`/wails/runtime.js`调用`Call.ByID(<fnv-id>, ...args)`（在`internal/runtime/desktop/@wailsio/runtime/src/calls.ts`中实现）；对于名称模式构建，则调用`Call.ByName("pkg.Struct.Method", ...args)`。
2. 运行时辅助程序会封装该调用，并通过各平台的原生桥接将其分派给 Go。
3. <strong>Go</strong>在`pkg/application/messageprocessor_call.go`中接收该消息。
4. 处理器在`pkg/application/bindings.go`（手工编写，基于`reflect`）中查找已绑定的方法并调用它。
5. 结果或错误会被编组并传回 JS，在那里相应的`Promise`会兑现或拒绝。

> 确切的 JSON 信封格式由 JS 端的运行时辅助程序和
>
> Go 端的`messageprocessor_call.go`共同定义；本页的旧版草稿
>
> 曾引用一种`{"t":"c","id":"123","m":"Greet","p":[…]}`结构，但它与
>
> 当前实现不符。排查
>
> 数据传输格式错误。

专用处理器：

| 文件 | 用途 |
| --- | --- |
| `messageprocessor_window.go` | 窗口操作（隐藏、最大化等） |
| `messageprocessor_dialog.go` | 原生对话框（`OpenFile`、`MessageBox`等） |
| `messageprocessor_clipboard.go` | 剪贴板读写 |
| `messageprocessor_events.go` | 事件订阅/触发 |
| `messageprocessor_browser.go` | 浏览器导航、开发者工具 |

处理器是<strong>无状态的</strong>，所需的一切都从每条消息随附的 `ApplicationContext`中获取。

---

## 4. 事件系统

事件是带命名空间的字符串，跨三个层级分派：

1. **应用程序事件**：全局生命周期（`application:ready`、`application:shutdown`）。
2. **窗口事件**：每个窗口独立（`window:focus`、`window:resize`）。
3. **自定义事件**：由用户定义（`chat:new-message`）。

实现细节：

- 事件常量位于`pkg/events/`中（`defaults.go`、`known_events.go`、`events.txt`）。它们由`v3/tasks/events/generate.go`生成，并公开为`events.Common.*`、`events.Mac.*`、`events.Windows.*`和`events.Linux.*`。可以使用`wails3 generate constants`重新生成它们。
- Go 端（应用程序事件）：
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- Go 端（窗口事件）：
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- Go 端（自定义事件）：
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- JS 端：
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


应用程序事件、窗口事件和自定义事件均经由 `pkg/application/event_manager.go`传递。窗口事件订阅的作用域限定在其所属窗口内，因此关闭窗口时会自动注销其处理程序。

---

## 5. 平台特定实现

条件编译在隐藏操作系统差异的同时，使公共 API 保持一致。

| 关注点 | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| 主线程 | `mainthread_darwin.go`（通过 Cgo 调用 Foundation） | `mainthread_linux.go`（GTK） | `mainthread_windows.go`（Win32 `AttachThreadInput`） |
| 对话框 | `dialogs_darwin.*`（NSAlert） | `dialogs_linux.go`（GtkFileChooser） | `dialogs_windows.go`（IFileOpenDialog） |
| 剪贴板 | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| 托盘图标 | `systemtray_darwin.*` | `systemtray_linux.go`（DBus） | `systemtray_windows.go`（Shell_NotifyIcon） |

关键原则：

- <strong>macOS 和 Windows</strong>仅少量使用 Cgo（主要通过`pkg/mac/`以及`pkg/w32`中的`w32` Win32 封装器）。
- 受实际需要所限，**Linux 大量使用 Cgo**——`pkg/application/linux_cgo.go`（约69 KB）和`linux_cgo_gtk4.{c,go,h}`（约50 KB 以上）直接驱动 GTK/WebKitGTK。
- 使用<strong>构建标签</strong>（`//go:build darwin`、`//go:build linux`……）来保持各操作系统专用文件的可读性。
- `internal/capabilities/`用于定义各平台的能力标志，但框架<strong>不会</strong>导出`ErrCapability`哨兵值——功能门控通过平台特定的桩实现返回值完成。

---

## 6. 文件指南

| 文件 | 何时需要修改 |
| --- | --- |
| `internal/runtime/runtime_*.go` | 修改小型构建标签桩层（开发与生产环境、操作系统特定的衔接代码）。 |
| `pkg/application/webview_window_*.go` | 实现新的窗口提示或行为。 |
| `pkg/application/messageprocessor*.go` | 添加可从 JS 调用的新桥接命令。 |
| `pkg/events/*.go` | 扩展内置事件定义（然后重新运行`wails3 generate constants`）。 |
| `internal/assetserver/*` | 调整开发/生产环境下的资源处理。 |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | 编辑嵌入式 JS 运行时（调用/事件分派、对话框、拖动等）。 |

---

## 7. 调试技巧

- 配置`Options.LogLevel`（例如`slog.LevelDebug`）并检查`Options.Logger`输出——不存在`WAILS_LOG_LEVEL`环境变量。
- `wails3 dev`标志包括`--config`、`--port`、`-s`（启用 HTTPS）以及全局`--no-colour`。不存在`-verbose`标志。
- 在 macOS 上，请在`lldb --`下运行，以便尽早捕获 Objective-C 异常。
- 对于 Windows 上的 Chromium 问题，请启用 WebView2 调试日志：`set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. 扩展运行时

1. 可选：在`internal/capabilities/`中声明任何新的能力标志。
2. 使用构建标签在每个`pkg/application/*_{darwin,linux,windows}.go`变体中实现该功能。对于不支持该功能的平台，请提供桩实现。
3. 在`pkg/application`中添加公共 API（接口和具体的`WebviewWindow`方法、选项结构体等）。
4. 如果 JS 需要调用该功能，请注册新的消息处理器方法（`pkg/application/messageprocessor*.go`），并在 JS 运行时中添加与之匹配的辅助函数。
5. 如果添加事件，请在`pkg/events/`中声明其常量，并运行`wails3 generate constants`以刷新生成的文件。

按照此检查清单操作，即可保持跨平台契约不变。

---

## 9. 拖放

在所有平台上，文件拖放都采用<strong>以 JavaScript 为主的方式</strong>。原生层会拦截操作系统的拖动事件，但实际的放置处理和 DOM 交互在 JavaScript 中进行。

### 流程

1. 用户从操作系统中将文件拖到 Wails 窗口上方
2. 原生层检测到拖动，并通知 JavaScript 显示悬停效果
3. 用户放下文件
4. 原生层将文件路径和坐标发送给 JavaScript
5. JavaScript 查找放置目标元素（`data-file-drop-target`）
6. JavaScript 将文件路径和元素详细信息发送给 Go 后端
7. Go 发出包含完整上下文的`WindowFilesDropped`事件

### 平台实现

| 平台 | 原生层 | 主要难点 |
| --- | --- | --- |
| **Windows** | WebView2 的内置拖动支持 | 坐标以 CSS 像素为单位，无需转换 |
| **macOS** | NSWindow 拖动委托 | 将窗口相对坐标转换为 WebView 相对坐标 |
| **Linux** | GTK3 拖动信号 | 必须区分文件拖动和内部 HTML5 拖动 |

### Linux：区分拖动类型

GTK 和 WebKit 都希望处理拖动事件。关键在于检查拖动目标类型：

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

对于内部拖动，信号处理程序返回`FALSE`（交由 WebKit 处理）；对于文件拖动，则返回`TRUE`（由我们自行处理）。

### 阻止文件放置

当`EnableFileDrop`为`false`时，我们仍需阻止浏览器导航到被放置的文件。各平台的处理方式不同：

- **Windows**：JavaScript 对拖动事件调用`preventDefault()`
- **macOS**：JavaScript 对拖动事件调用`preventDefault()`\
- **Linux**：GTK 信号处理程序在原生层拦截并拒绝文件拖动

### 关键文件

| 文件 | 用途 |
| --- | --- |
| `pkg/application/linux_cgo.go` | GTK 拖动信号处理程序（cgo 前导部分中的 C 代码） |
| `pkg/application/webview_window_darwin.go` | macOS 拖动委托 |
| `pkg/application/webview_window_windows.go` | WebView2 消息处理 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | JavaScript 放置处理 |

### 调试

- **Linux**：在 C 代码中添加`printf`（别忘了`fflush(stdout)`）
- **Windows**：使用`globalApplication.debug()`
- **JavaScript**：检查浏览器控制台并启用调试模式

常见问题：

1. **内部 HTML5 拖动不起作用**：原生处理程序拦截了该拖动（对于非文件拖动应返回`FALSE`）
2. **未显示悬停效果**：JavaScript 处理程序未被调用
3. **坐标错误**：检查坐标空间转换

---

现在，你已经在引导下了解了运行时的内部机制。结合这些知识以及<strong>代码库布局</strong>示意图和<strong>资源服务器</strong>文档，你便能自信地浏览代码库并做出富有影响力的贡献。祝编码愉快！
