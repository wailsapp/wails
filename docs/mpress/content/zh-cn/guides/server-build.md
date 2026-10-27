---
title: "服务器构建"
description: "将 Wails 应用程序作为 HTTP 服务器运行，无需原生 GUI 窗口"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 支持服务器模式，可将应用程序作为纯 HTTP 服务器运行，无需创建原生窗口或依赖 GUI。这让同一个 Wails 应用程序可以部署到服务器、容器和 Web 浏览器。

服务器模式适用于：

- **Docker/容器部署**——无需依赖 X11/Wayland 即可运行
- **服务器端应用程序**——部署为可通过浏览器访问的 Web 服务器
- **仅限 Web 访问**——桌面端和 Web 端共用同一代码库
- **CI/CD 测试**——无需显示服务器即可运行集成测试
- **微服务**——在无头后端服务中使用 Wails 绑定

## 快速开始

服务器模式通过 `server` 构建标签启用。应用程序代码保持不变，只需使用该标签进行构建：

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

下面是一个最简示例：

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

同一份代码既可构建为桌面模式（不带该标签），也可构建为服务器模式（带 `-tags server`）。

## 配置

### ServerOptions

使用 `ServerOptions` 配置 HTTP 服务器：

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections and call
    // /wails/runtime. Same-origin requests are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket and runtime origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## 功能

### 健康检查端点

`/health` 会自动提供健康检查端点：

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

这适用于：

- Kubernetes 存活探针和就绪探针
- 负载均衡器健康检查
- 监控系统

### 服务绑定

所有服务绑定的工作方式都与桌面模式完全相同：

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

前端可以使用标准 Wails 运行时调用这些绑定：

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### 事件

在服务器模式下，事件支持双向传递：

- **前端到后端**：浏览器发出的事件通过 HTTP 发送，并由 Go 事件处理程序接收
- **后端到前端**：Go 发出的事件通过 WebSocket 广播给所有已连接的浏览器

每个浏览器标签页都表示为一个具有唯一名称（`browser-1`、`browser-2` 等）的“窗口”，可通过 `event.Sender` 访问：

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

在前端中：

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### 优雅关闭

服务器会妥善处理 `SIGINT` 和 `SIGTERM` 信号：

1. 停止接受新连接
2. 等待正在处理的请求完成（最长等待 `ShutdownTimeout`）
3. 运行 `OnShutdown` 钩子
4. 按相反顺序关闭服务

## 与桌面模式的区别

| 功能 | 桌面模式 | 服务器模式 |
| --- | --- | --- |
| 原生窗口 | 会创建 | 浏览器窗口（`browser-N`） |
| 系统托盘 | 可用 | 不可用 |
| 原生对话框 | 可用 | 不可用 |
| 应用程序菜单 | 可用 | 不可用 |
| 屏幕信息 | 可用 | 没有原生屏幕；`GetPrimary()` 返回 `nil` |
| 服务绑定 | 可用 | 可用 |
| 事件 | 可用 | 可用（通过 WebSocket） |
| 资源 | 通过 WebView | 通过 HTTP |
| 需要 CGO | 是 | 否 |

### 窗口 API 行为

在服务器模式下，窗口相关 API 会得到安全处理：

- `app.Window.NewWithOptions(options)`——返回一个非 nil 的 Go 窗口对象；这不会创建原生桌面窗口
- `app.Hide()` / `app.Show()`——不执行任何操作
- `app.Screen.GetPrimary()`——返回 `*application.Screen`；没有可用的主屏幕时，其值为 `nil`。解引用前请先检查

请区分 Go 窗口对象与原生窗口是否可用。构建服务器时，应对仅限桌面端的行为进行防护；Go 窗口对象非 nil 并不表示存在原生窗口或屏幕。

### 支持的平台

服务器模式使用纯 Go 实现，无需 CGO，即可针对 Wails 支持的所有 `GOOS`（包括 `android` 和 `ios`）进行编译。在移动端目标上，它只会构建无头服务器：

- 原生移动端 API `application.Android` 和 `application.IOS` 不可用。使用这些 API 的代码必须通过 `!server` 构建约束从服务器构建中排除。
- 默认日志记录器使用桌面端实现；在开发构建中，它会将日志写入 stderr，而不是 logcat 或 WebView 控制台。

## 生产环境构建

### 使用 Task（推荐）

使用 `wails3 init` 创建的项目包含一个 `build:server` 任务：

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### 手动构建

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Wails 项目包含开箱即用的 Docker 配置。要在容器中构建并运行应用程序：

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

就这么简单！应用程序将在 `http://localhost:8080` 上可用。

你可以使用几个选项自定义构建：

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

生成的 `Dockerfile.server` 会创建一个基于 distroless 的最小化镜像。它会自动处理网络绑定，因此可以从容器外部访问应用程序。

### Docker Compose

对于更复杂的部署，可以使用以下包含健康检查的 Docker Compose 配置：

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
此健康检查示例使用`wget`。如果使用 distroless 基础镜像，则需要在镜像中包含健康检查二进制文件，或使用外部健康检查机制（例如 Docker 的`curl`选项或边车容器）。

@end

### 自定义 Dockerfile

如果需要更强的控制能力，可以自行创建 Dockerfile。需要注意的关键事项是设置`WAILS_SERVER_HOST=0.0.0.0`，以便服务器接受来自容器外部的连接：

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## 安全注意事项

部署服务器模式应用程序时：

1. **默认绑定到 localhost**——仅在需要时使用`0.0.0.0`
2. **在生产环境中使用 TLS**——配置`ServerOptions.TLS`
3. **置于反向代理之后**——使用 nginx/traefik 增强安全性
4. **确保 WebSocket 和运行时调用同源**——仅使用`WebSocketOriginPatterns`添加可信来源；避免使用`WebSocketAllowAllOrigins`
5. **验证所有输入**——采用与任何 Web 应用程序相同的安全实践

## 示例

完整示例位于`v3/examples/server/`：

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## 环境变量

对于需要在不更改代码的情况下覆盖服务器配置的部署场景，Wails 可识别以下环境变量：

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | 要绑定的网络接口 | `localhost` |
| `WAILS_SERVER_PORT` | 要监听的端口 | `8080` |

这些环境变量的优先级高于代码中的`ServerOptions`，因此 Docker 示例会设置`WAILS_SERVER_HOST=0.0.0.0`——这使容器无需对应用程序进行任何更改即可接受外部连接。

## 另请参阅

- [自定义传输](/guides/custom-transport/)——用于高级 IPC 自定义
- [服务](/features/bindings/services/)——服务绑定文档
- [事件](/guides/events-reference/)——事件系统文档

### 运行时请求大小

在进行 JSON 处理之前，对`/wails/runtime`的请求大小上限为64 MiB。超过此限制的普通请求会收到 HTTP 413响应，包括不含`Content-Length`的请求。分块运行时上传仍分别遵循每个分块1 MiB 和组装后有效载荷64 MiB 的限制。适当时可使用应用程序中间件或反向代理施加更低的限制。

### 运行时请求来源

`/wails/runtime`仅接受`POST`请求；其他方法会收到 HTTP 405响应。如果浏览器请求的`Origin`既不是请求主机，也不匹配`WebSocketOriginPatterns`，则会收到 HTTP 403响应；缺少`Origin`且`Sec-Fetch-Site`为`cross-site`或`same-site`的请求也会收到该响应。既不携带这两个标头的请求（例如来自`curl`或其他非浏览器客户端的请求）会被接受。如果反向代理重写`Host`标头，请将公共来源添加到`WebSocketOriginPatterns`。

`WebSocketOriginPatterns`用于反向代理部署，此时公共来源与应用程序看到的`Host`不同。不支持从其他来源提供并直接调用`/wails/runtime`的前端：该端点不会发出 CORS 标头，并会对预检`OPTIONS`请求返回 HTTP 405。请从 Wails 服务器提供前端，或将二者置于同一个代理来源之后。

### 环回 Host 检查

当服务器绑定环回地址（默认的`localhost`、任何`127.0.0.0/8`地址或`::1`，无论通过`ServerOptions.Host`还是`WAILS_SERVER_HOST`设置）时，每个请求都必须携带`Host`标头，其中指定`localhost`、IP 地址字面量、绑定主机，或与某个`WebSocketOriginPatterns`条目的主机部分相匹配的主机；任何其他`Host`都会收到 HTTP 403响应。此检查可阻止 DNS 重绑定攻击：恶意网页将其自身的主机名指向`127.0.0.1`，从而使其请求通过上述同源检查。端口不会参与比较，并且会拒绝`localhost.`（带尾随点）或`127.1`等写法。即使设置了`WebSocketAllowAllOrigins`，此检查仍会执行。如果本地反向代理转发其公共`Host`，请将该主机添加到`WebSocketOriginPatterns`。绑定到其他地址（例如容器中的`0.0.0.0`）不受影响。
