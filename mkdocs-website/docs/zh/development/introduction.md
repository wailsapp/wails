# 介绍

!!! note 这个指南仍在制作中。

感谢您想要帮助开发Wails！本指南将帮助您入门。

## 入门指南

- Git 克隆此存储库。切换到 `v3-alpha` 分支。
- 安装 CLI：`cd v3/cmd/wails3 && go install`

- 可选：如果您想要使用构建系统构建前端代码，您需要安装 [npm](https://nodejs.org/en/download)。

## 构建

对于简单的程序，您可以使用标准的 `go build` 命令。也可以使用 `go run`。

Wails 还配备了一个构建系统，可用于构建更复杂的项目。它使用了强大的 [Task](https://taskfile.dev) 构建系统。要了解更多信息，请查看任务主页或运行 `wails task --help`。

## 项目结构

该项目具有以下结构：

    ```
    v3
    ├── cmd/wails3                  // CLI
    ├── examples                   // Wails 应用示例
    ├── internal                   // 内部包
    |   ├── runtime                // Wails JS 运行时
    |   └── templates              // 支持的项目模板
    ├── pkg
    |   ├── application            // 核心 Wails 库
    |   └── events                 // 事件定义
    |   └── mac                    // 由插件使用的 macOS 特定代码
    |   └── w32                    // Windows 特定代码
    ├── plugins                    // 支持的插件
    ├── tasks                      // 通用任务
    └── Taskfile.yaml              // 开发任务配置
    ```

## 开发

### 添加窗口功能

添加窗口功能的首选方法是在 `pkg/application/webview_window.go` 文件中添加一个新函数。这应该实现所有平台所需的功能。任何特定于平台的代码都应通过 `webviewWindowImpl` 接口方法调用。该接口由每个目标平台实现，以提供平台特定的功能。在某些情况下，这可能不执行任何操作。添加接口方法后，请确保每个平台都实现了它。一个很好的例子是 `SetMinSize` 方法。

- Mac: `webview_window_darwin.go`
- Windows: `webview_window_windows.go`
- Linux: `webview_window_linux.go`

大多数，如果不是全部，特定于平台的代码应在主线程上运行。为了简化这一点，在 `application.go` 中定义了一些 `invokeSync` 方法。

### 更新运行时

运行时位于 `v3/internal/runtime`。更新运行时时，需要执行以下步骤：

```shell
wails3 task runtime:build
```

### 事件

事件定义在 `v3/pkg/events` 中。当添加新事件时，需要执行以下步骤：

- 将事件添加到 `events.txt` 文件中
- 运行 `wails3 task events:generate`

有几种类型的事件：特定于平台的应用程序和窗口事件 + 通用事件。通用事件对于跨平台事件处理很有用，但您不必局限于“最低公共分母”。如果需要，可以使用特定于平台的事件。

添加通用事件时，请确保映射了特定于平台的事件。一个示例是在 `window_webview_darwin.go` 中：

```go
		// 将 ShouldClose 转化为通用的 WindowClosing 事件
		w.parent.On(events.Mac.WindowShouldClose, func(_ *WindowEventContext) {
			w.parent.emit(events.Common.WindowClosing)
		})
```

注意：我们可能会尝试通过将映射添加到事件定义中来自动化此过程。

### 插件

插件是扩展 Wails 应用功能的一种方式。

#### 创建插件

插件是符合以下接口的标准 Go 结构：

```go
type Plugin interface {
    Name() string
    Init(*application.App) error
    Shutdown()
    CallableByJS() []string
    InjectJS() string
}
```

`Name()` 方法返回插件的名称。这用于日志记录。

`Init(*application.App) error` 方法在加载插件时调用。`*application.App` 参数是加载插件的应用程序。任何错误都将阻止应用程序启动。

`Shutdown()` 方法在应用程序关闭时调用。

`CallableByJS()` 方法返回可以从前端调用的导出函数列表。这些方法的名称必须与插件导出的方法的名称完全匹配。

`InjectJS()` 方法返回应注入到所有窗口中的 JavaScript。这对于添加与插件相补充的自定义 JavaScript 函数非常有用。

内置插件可以在 `v3/plugins` 目录中找到。参考它们以获得灵感。

## 任务

Wails CLI 使用 [Task](https://taskfile.dev) 构建系统。它作为库导入并用于运行 `Taskfile.yaml` 中定义的任务。与 Task 的主要交互发生在 `v3/internal/commands/task.go` 中。

### 升级 Taskfile

要检查是否有 Taskfile 的升级，请运行 `wails3 task -version` 并检查 Task 网站。

要升级使用的 Taskfile 版本，请运行：

```shell
wails3 task taskfile:upgrade
```

如果存在不兼容性，则应在 `v3/internal/commands/task.go` 文件中显示。

通常，修复不兼容性的最佳方法是克隆 `https://github.com/go-task/task` 上的任务存储库，并查看 git 历史记录以确定发生了什么变化以及原因。

要检查所有更改是否正确工作，请重新安装 CLI 并再次检查版本：

```shell
wails3 task cli:install
wails3 task -version
```

## 打开 PR

确保所有 PR 都有与之关联的工单，以提供更改的上下文。如果没有工单，请先创建一个。确保所有 PR 都已使用所做的更改更新了 CHANGELOG.md 文件。CHANGELOG.md 文件位于 `mkdocs-website/docs` 目录中。

## 其他任务

### 升级 Taskfile

Wails CLI 使用 [Task](https://taskfile.dev) 构建系统。它作为库导入并用于运行 `Taskfile.yaml` 中定义的任务。与 Task 的主要交互发生在 `v3/internal/commands/task.go` 中。

要检查是否有 Taskfile 的升级，请运行 `wails3 task -version` 并检查 Task 网站。

要升级使用的 Taskfile 版本，请运行：

```shell
wails3 task taskfile:upgrade
```

如果存在不兼容性，则应在 `v3/internal/commands/task.go` 文件中显示。

通常，修复不兼容性的最佳方法是克隆 `https://github.com/go-task/task` 上的任务存储库，并查看 git 历史记录以确定发生了什么变化以及原因。

要检查所有更改是否正确工作，请重新安装 CLI 并再次检查版本：

```shell
wails3 task cli:install
wails3 task -version
```