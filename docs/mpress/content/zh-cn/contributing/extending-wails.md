---
title: "扩展 Wails"
description: "向 Wails v3 添加新功能和新平台的实用指南"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> Wails 的设计理念是<strong>易于改造</strong>。
>
> 每个主要子系统都以 Go 代码实现，你可以阅读、修改并发布这些代码。
>
> 本页介绍在执行以下操作时应从<em>何处</em>着手，以及<em>如何</em>保持跨平台兼容：

- 添加<strong>服务</strong>（通知、KV 存储、自定义 IPC 等）
- 创建<strong>新的 CLI 命令</strong>（`wails3 <foo>`）
- 扩展<strong>运行时</strong>（窗口 API、对话框、事件）
- 引入<strong>平台能力</strong>（Wayland 等）
- 在不被`//go:build`标签淹没的情况下保持<strong>跨平台兼容性</strong>

---

## 1. 添加服务

v3 中的“服务”是用户提供的 Go 类型，通过 `application.Options.Services`注册，并借助生成的绑定向 JS 公开。v3 代码库提供：

- `internal/service/` — 用于`wails3 generate service`的脚手架：
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — 可立即注册使用的现成服务（通知、kvstore、sqlite、日志、文件服务器、dock 等）。

旧版草稿中引用的生成器和 CLI 文件 `internal/service/template/template.go`及 `internal/generator/collect/services.go`并不存在——脚手架工具是 `internal/service/service.go`（入口点为`service.Install`），服务的绑定 元数据则收集在`internal/generator/collect/service.go`中。

### 1.1 定义服务

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 实现生命周期接口（可选）

服务可以选择实现以下接口（来自`pkg/application`）：

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **重要提示：**`ServiceShutdown`**不接受任何参数**。具有以下
>
> 签名的方法`ServiceShutdown(ctx context.Context) error`<strong>不</strong>满足
>
> 该接口，因此绝不会被调用，也不会产生任何提示。

### 1.3 向应用程序注册服务

不存在全局`services.Register(...)`调用。服务在运行时通过 `application.Options.Services`注册：

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

注册后，`wails3 generate bindings`会在 `frontend/bindings/<your import path>/...`下生成封装了导出方法的 ES 模块。

### 1.4 从 JS 调用

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

v3 中不存在全局`window.backend.*`——调用需通过生成的 ES 模块进行，这些模块继而调用`/wails/runtime.js`中的`Call.ByID(...)`。

---

## 2. 编写新的 CLI 命令

v3 CLI 使用<strong>`github.com/leaanthony/clir`</strong>（而非 cobra）。相关装配代码位于 `v3/cmd/wails3/main.go`中：

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

不存在基于`init()`的自动注册。请将新子命令添加到 `cmd/wails3/main.go`，并在`internal/commands/`中添加相应的支持函数（如果该命令接受选项，还需在 `internal/flags/`下添加标志结构体）。重新构建 CLI：

```
cd v3
go install ./cmd/wails3
wails3 hello
```

如果命令需要 Taskfile 粘合代码，请复用 `internal/commands/task_wrapper.go`中的辅助函数（`wrapTask("yourtask", args)`）。

---

## 3. 修改运行时

常见原因：

- 新的窗口功能（`SetOpacity`、`Shake`等）
- 额外的对话框（`ColorPicker`）
- 系统级 API（屏幕亮度）

### 3.1 公共 API

将该方法添加到`pkg/application/webview_window.go`（接口位于 `window.go`中）：

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

使用现有的`InvokeSync`/`InvokeAsync`辅助函数，确保调用在 主线程上运行。

### 3.2 消息处理器

如果 JS 需要调用新方法，请扩展相应的 `pkg/application/messageprocessor_*.go`文件。消息处理器使用 `MessageProcessor`上基于 switch 的方法，而不是全局`register(...)`调用：

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

代码库中<strong>没有</strong>`messageprocessor_window_opacity.go`文件，也没有基于`init()`的 `register(MsgSetOpacity, ...)`模式。

### 3.3 平台实现

在`pkg/application/`下每个操作系统对应的文件中添加实现：

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

如果某个平台不支持该功能，请编写一个无操作存根。该框架 没有`ErrCapability`哨兵值——请在文档中说明支持情况；如有 需要，也通过`Options`或特定于平台的选项结构体中的 相关布尔字段公开支持情况。

### 3.4 能力标志（可选）

`internal/capabilities/`包用于声明各平台的 能力集。不存在公共`application.HasCapability`/ `application.CapOpacity` API。如果希望某项能力可在运行时检查， 请将其添加到`internal/capabilities/`下，并从 `pkg/application`公开一个类型化的 getter。

---

## 4. 添加新的平台能力

示例：Linux 上的可选 Wayland 支持。

1. 将相关的`pkg/application/*_linux.go`文件拆分为`*_linux_x11.go`（`//go:build linux && !wayland`）和`*_linux_wayland.go`（`//go:build linux && wayland`）。
2. 让用户通过`wails3 build --tags wayland`选择启用。利用`internal/commands/task_wrapper.go`中现有的`EXTRA_TAGS`传递机制转发额外标签。不存在`dev`级别的`--tags wayland`标志——`wails3 dev`仅接受`--config`、`--port`和`-s`。
3. 更新文档以及`pkg/application/`下所有特定于平台的 README。

> 尽量减少默认构建标签；仅为小众功能使用需选择启用的标签。

---

## 5. 跨平台兼容性检查清单

| ✅ 步骤 | 原因 |
| --- | --- |
| 在所有平台文件中提供<strong>每个</strong>公共方法（即使只是存根） | 确保在每种操作系统上都能成功构建 |
| 记录各操作系统上的优雅降级行为 | 应用可以根据`runtime.GOOS`进行分支处理，而不会出现隐藏错误 |
| 优先使用<strong>纯 Go</strong>，仅在必要时使用 Cgo | 简化交叉编译（Linux 已经承担了使用 Cgo 的代价） |
| 运行`task test:cli`、`task test:generator`和`task test:templates` | 在本地复现 CI 测试流程 |
| 在贡献者文档／模板 README 中记录新的构建标签 | 用户必须了解需要主动启用的功能 |

---

## 6. 调试构建与迭代速度

- 使用`Options.LogLevel = slog.LevelDebug`（`Options.Logger = slog.Default()`）输出详细的运行时活动。不存在`WAILS_LOG_LEVEL`环境变量。
- `wails3 dev`的标志包括`--config`、`--port`和`-s`。不存在`-race`或`-verbose`标志——请使用`go test -race ./...`运行竞态检测器，或者对应用执行`go build -race`后直接运行。
- 竞态／Cgo 测试指南位于`v3/TESTING.md`（旧版草稿指向了并不存在的`pkg/application/RACE.md`）。

---

## 7. 向上游贡献

1. 如需新增功能或更改公共行为，请提交<strong>WEP（Wails Enhancement Proposal）</strong>草稿 PR，以便讨论构想和设计。只有可复现的缺陷或文档问题才应使用 issue。
2. 按照上述方法进行实现。
3. 添加：
  - 单元测试（`*_test.go`）
  - 文档（此文件或相关的`docs/...`页面）
  - 如果修改了绑定生成器，请在`internal/generator/testcases/`下添加回归测试

4. 推送前，请在本地运行`task precommit`以及相关的`task test:*`目标。

---

### 快速链接

| 领域 | 位置 |
| --- | --- |
| 内置服务 | `pkg/services/` |
| 服务脚手架生成器 | `internal/service/` |
| CLI 命令的注册与接入 | `v3/cmd/wails3/main.go` |
| CLI 命令主体 | `internal/commands/` |
| 各操作系统的运行时实现 | `pkg/application/*_{darwin,linux,windows}.go` |
| 能力声明 | `internal/capabilities/` |
| Taskfile DSL | `v3/Taskfile.yaml` |
| 事件常量生成器 | `v3/tasks/events/generate.go` |

---

现在，你已经有了一份随心改造 Wails 的<strong>路线图</strong>——可以添加服务、 为 CLI 注入魔法、改造运行时，或引入全新的操作系统功能。 祝扩展愉快！
