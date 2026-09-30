---
title: "LLM 控制（MCP）"
description: "让 LLM 智能体通过模型上下文协议测试和控制你的应用"
slug: "guides/mcp-service"
sourcePath: "guides/mcp-service.md"
---

@note{type="caution" title="实验性功能"}
内置 MCP 服务器尚属实验性功能，其 API 在未来版本中可能会发生变化。

@end

Wails v3 内置了[模型上下文协议](https://modelcontextprotocol.io)（MCP）服务器，LLM 智能体（Claude Code、IDE 助手或任何 MCP 客户端）可通过它检查、测试和操控<strong>正在运行的</strong> Wails 应用。

如需自动化项目生命周期，请使用单独的`wails3 mcp` CLI 服务器。智能体可通过它检查和初始化项目、运行诊断、启动构建和开发任务、生成绑定、运行指定的 Taskfile 任务，以及获取有界的任务输出。默认情况下，CLI 服务器仅限访问当前目录，且不提供任意 shell 命令执行功能。有关传输、身份验证和工具的详细信息，请参阅[CLI MCP 文档](/guides/cli/#mcp)。

启用后，连接到应用的智能体可以：

- **列出并控制窗口**——调整大小和位置、聚焦、全屏显示、打开开发者工具、重新加载等
- **检查 DOM**——查询元素、获取 HTML、生成结构快照
- **执行 JavaScript**——在任意窗口内运行任意代码并获取结果
- **模拟用户输入**——通过<strong>带动画效果的屏幕光标</strong>呈现鼠标移动、单击、拖动和滚动，让你可以看到智能体的操作过程
- **输入文本和按键**——生成逼真的逐字符事件，可用于 React 受控输入框
- **调用已绑定的 Go 方法**，以及发出或等待应用事件

## 工作原理

仅当存在<strong>`mcp`构建标签</strong>时，MCP 服务器才会被编译到应用中。如果没有该标签，服务器代码将完全不包含在二进制文件中——没有运行时开销，不会开放端口，也不会增加攻击面。

存在该标签时，服务器会在`App.Run()`中自动启动，默认绑定到`127.0.0.1:9099`，并记录其端点。无需编写任何用户代码。

## 教程

### 第1步——编写常规 Wails 应用

MCP 无需导入或注册。按通常方式创建应用即可：

```go {title="main.go"}
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Width: 1024, Height: 768,
    })

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### 第2步——使用`mcp`标签构建或运行

@tabs
[Wails CLI（推荐）]
设置`WAILS_MCP=1`后，Wails CLI 会自动为你添加`mcp`标签：

```shell
# Development
WAILS_MCP=1 wails3 dev

# Production build
WAILS_MCP=1 wails3 build
```

[直接使用 Go]
将该标签直接传递给`go run`或`go build`：

```shell
go run -tags mcp .
go build -tags mcp -o myapp .
```

[Windows（PowerShell）]
```powershell
$env:WAILS_MCP = "1"
wails3 dev
# or
wails3 build
```

@end

应用启动时会记录 MCP 端点：

```
INFO MCP server started. Connect MCP clients using the streamable HTTP transport.
     url=http://127.0.0.1:9099/mcp
```

### 第3步——连接客户端

服务器使用<strong>MCP Streamable HTTP 传输</strong>。可使用任何兼容 MCP 的客户端连接。

@tabs
[Claude Code]
```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

然后让 Claude 与你的应用交互：

```
Click the "Submit" button, then verify a success toast appears.
```

[VS Code（GitHub Copilot）]
添加到`.vscode/settings.json`：

```json
{
  "github.copilot.chat.mcp.enabled": true,
  "mcp": {
    "servers": {
      "my-wails-app": {
        "type": "http",
        "url": "http://127.0.0.1:9099/mcp"
      }
    }
  }
}
```

[其他客户端]
将任何支持 Streamable HTTP 传输的 MCP 客户端指向：

```
http://127.0.0.1:9099/mcp
```

@end

### 第4步——运行测试会话

让智能体操作并测试你的应用。下面是一些示例提示词：

```
Take a DOM snapshot of the main window.
```

```
Click the "Add item" button, type "Hello world" in the input field,
then press Enter and verify the item appears in the list.
```

```
Call the bound method main.GreetService.Greet with argument ["World"]
and return the result.
```

```
Wait for the event "save:complete" while clicking the Save button.
```

## 配置

所有配置均通过环境变量完成，无需更改代码。

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `WAILS_MCP` | （未设置） | 使用 Wails CLI 时，将其设为`1`、`true`、`on`或`yes`，即可自动添加`mcp`构建标签。 |
| `WAILS_MCP_HOST` | `127.0.0.1` | 要绑定的网络接口。绑定到非环回接口时必须设置`WAILS_MCP_TOKEN`。 |
| `WAILS_MCP_TOKEN` | 未设置 | 在环回接口上可选择使用的持有者令牌；绑定到其他地址时必须设置。客户端发送`Authorization: Bearer <token>`。 |
| `WAILS_MCP_PORT` | `9099` | 监听端口。设为`0`时，将随机分配一个空闲端口（并在日志中输出）。 |
| `WAILS_MCP_TIMEOUT` | `30000` | JavaScript 求值的默认超时时间，以<strong>毫秒</strong>为单位。 |
| `WAILS_MCP_HIDE_CURSOR` | （未设置） | 设为`1`或`true`可禁用动画光标叠加层。 |

示例——自定义端口和60秒超时时间：

```shell
WAILS_MCP=1 WAILS_MCP_PORT=9200 WAILS_MCP_TIMEOUT=60000 wails3 dev
```

## 可用工具

| 工具 | 用途 |
| --- | --- |
| `app_info` | 应用信息：平台、架构、所有窗口和 MCP 端点 |
| `windows_list` | 列出所有窗口及其几何信息和状态 |
| `window_control` | 聚焦、调整大小、移动、全屏显示、打开开发者工具、重新加载、设置 URL 等（22项操作） |
| `js_eval` | 在窗口中执行 JavaScript（异步函数体，使用`return`返回值） |
| `dom_html` | 获取页面或特定元素的 HTML |
| `dom_query` | 通过 CSS 选择器查找元素——获取标签、文本、边界和可见性 |
| `screenshot_dom` | 可见页面的结构化快照（基于 DOM，不含像素） |
| `mouse_move` | 以动画方式将光标移动到某个点或 CSS 选择器所指向的元素 |
| `mouse_click` | 使用动画光标单击（左键、右键或中键；支持双击和修饰键） |
| `mouse_drag` | 使用动画光标拖动（支持 HTML5 拖放元素） |
| `mouse_scroll` | 在某个点或元素处滚动 |
| `keyboard_type` | 逐字符输入文本，并触发真实的事件 |
| `keyboard_press` | 按下单个按键（Enter、Tab、Escape、ArrowDown 等），可选用修饰键 |
| `call_bound_method` | 调用已绑定的 Go 服务方法，例如`main.GreetService.Greet` |
| `emit_event` | 发出 Wails 应用程序事件 |
| `wait_for_event` | 等待 Wails 应用程序事件并返回其数据 |

### 多窗口支持

所有对窗口执行操作的工具都接受一个可选的`window`参数，其中包含窗口的<strong>名称</strong>（通过`WebviewWindowOptions.Name`设置）。如果省略该参数，工具会以当前获得焦点的窗口为目标；如果没有窗口获得焦点，则以第一个窗口为目标。

```
List all windows, then click the "New" button in the window named "editor".
```

### 选择元素

鼠标和键盘工具接受以下任一方式：

- **CSS 选择器**——`selector: "#submit-btn"`（元素会自动滚动到可见区域）
- **坐标**——`x: 400, y: 300`（相对于视口的 CSS 像素）

对于拖动操作，请分别添加前缀`from_`和`to_`：

```
Drag from selector: ".card" to selector: ".dropzone"
```

## 安全性

@note{type="caution"}
MCP 服务器可通过程序完全控制你的应用程序。任何能够访问其工具的人都可以读取 DOM、执行 JavaScript、单击按钮以及调用 Go 方法。

@end

- 默认情况下，服务器绑定到`127.0.0.1`。浏览器来源必须是 HTTP(S) 环回来源；不透明（`null`）、格式错误及外部来源均会被拒绝。
- 仍然支持不带标头的原生 MCP 客户端。未设置`WAILS_MCP_TOKEN`时，本地进程和获准的本地浏览器来源会受到信任；来源检查并非身份验证。请设置一个高熵令牌，要求所有`/mcp`调用都进行持有者身份验证。在客户端的`Authorization: Bearer <token>`标头中配置同一个令牌。预检请求无需提供该令牌。
- `/eval-result`回调使用每次求值时生成的不可预测 ID，而不使用客户端持有者令牌，因此仍可兼容 WebView 结果传递。
- 生产构建<strong>不应</strong>包含`mcp`标签。只有明确设置`WAILS_MCP=1`时，Wails CLI 才会添加该标签；默认的`wails3 build`不含任何服务器代码。
- 如果需要在非环回接口上公开服务器（例如用于局域网测试），请设置`WAILS_MCP_HOST=0.0.0.0`和高熵`WAILS_MCP_TOKEN`；未设置令牌时，服务器将无法启动。在不受信任的网络上，请使用加密隧道或终止 TLS 的代理，因为内置监听器使用 HTTP。

## 示例应用程序

[`v3/examples/mcp`](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples/mcp)提供了一个演示所有工具的完整试验场应用程序。它包括：

- 带有递增和重置按钮的计数器
- 姓名输入框，以及已绑定的 Greet、Add、Shout 方法
- HTML5 拖放源和目标
- 可滚动列表（包含50个项目）
- 事件日志

使用以下命令运行：

```shell
cd v3/examples/mcp
go run -tags mcp .
```

然后将 Claude Code 或任意 MCP 客户端连接到`http://127.0.0.1:9099/mcp`，并让它操作该用户界面。

## 反馈

内置 MCP 服务器是一项实验，其未来方向取决于你的反馈。如果你试用了它，我们很想了解你使用了哪些客户端和工具、你的预期与实际结果，以及让代理操作你的应用程序是否有用——最有价值的报告会准确说明你运行了什么。请在[MCP 服务器反馈讨论](https://github.com/wailsapp/wails/discussions/5692)中告诉我们。
