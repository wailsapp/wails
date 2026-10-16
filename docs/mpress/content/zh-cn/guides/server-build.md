---
title: "服务器构建"
description: "将 Wails 应用程序作为 HTTP 服务器运行，无需原生 GUI 窗口"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 支持服务器模式，可将应用程序作为纯 HTTP 服务器运行，而无需创建原生窗口或依赖 GUI。这使同一个 Wails 应用程序能够部署到服务器、容器和 Web 浏览器。

服务器模式适用于：

- **Docker/容器部署**——无需依赖 X11/Wayland 即可运行
- **服务器端应用程序**——部署为可通过浏览器访问的 Web 服务器
- **仅通过 Web 访问**——桌面端与 Web 端共享同一代码库
- **CI/CD 测试**——无需显示服务器即可运行集成测试
- **微服务**——在无头后端服务中使用 Wails 绑定

## 快速入门

通过 `server` 构建标签启用服务器模式。应用程序代码无需更改，只需使用该标签进行构建：

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

下面是一个最小示例：

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

同一份代码既可构建为桌面模式（不使用该标签），也可构建为服务器模式（使用 `-tags server`）。

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

    // Additional origins allowed to open WebSocket connections.
    // Same-origin connections are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket origin checks. Unsafe; default: false.
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

`/health` 处会自动提供健康检查端点：

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

此端点适用于：

- Kubernetes 存活探针和就绪探针
- 负载均衡器健康检查
- 监控系统

### 服务绑定

所有服务绑定的工作方式均与桌面模式完全相同：

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
- **后端到前端**：Go 发出的事件通过 WebSocket 广播到所有已连接的浏览器

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

在前端：

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
2. 等待活动请求完成（最长等待 `ShutdownTimeout`）
3. 运行 `OnShutdown` 钩子
4. 按相反顺序关闭服务

## 与桌面模式的区别

| 功能 | 桌面模式 | 服务器模式 |
| --- | --- | --- |
| 原生窗口 | 会创建 | 浏览器窗口（`browser-N`） |
| 系统托盘 | 可用 | 不可用 |
| 原生对话框 | 可用 | 不可用 |
| 应用程序菜单 | 可用 | 不可用 |
| 屏幕信息 | 可用 | 返回错误 |
| 服务绑定 | 正常工作 | 正常工作 |
| 事件 | 正常工作 | 正常工作（通过 WebSocket） |
| 资源 | 通过 WebView | 通过 HTTP |
| 需要 CGO | 是 | 否 |

### 窗口 API 的行为

在服务器模式下，与窗口相关的 API 会得到安全处理：

- `app.Window.NewWithOptions()`——记录警告并返回 nil
- `app.Hide()` / `app.Show()`——不执行任何操作
- `app.Screen.GetPrimary()`——返回错误

这样一来，引用窗口的代码可以继续运行而不会崩溃，但窗口操作不会产生任何效果。

## 生产环境构建

### 使用 Task（推荐）

使用`wails3 init`创建的项目包含一个`build:server`任务：

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

Wails 项目包含一套开箱即用的 Docker 配置。要在容器中构建并运行应用程序：

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

就是这样！你的应用程序将在`http://localhost:8080`上可用。

你可以使用以下几个选项自定义构建：

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

生成的`Dockerfile.server`会创建一个基于 distroless 的最小化镜像。它会自动处理网络绑定，因此可以从容器外部访问你的应用程序。

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
此健康检查示例使用`wget`。如果使用 distroless 基础镜像，则需要在镜像中包含健康检查二进制文件，或者使用外部健康检查机制（例如 Docker 的`curl`选项或 sidecar 容器）。

@end

### 自定义 Dockerfile

如果需要更精细的控制，可以创建自己的 Dockerfile。需要特别注意的是，应设置`WAILS_SERVER_HOST=0.0.0.0`，以便服务器接受来自容器外部的连接：

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
3. **置于反向代理之后**——使用 nginx/traefik 提供额外的安全保护
4. **保持 WebSocket 同源**——仅使用`WebSocketOriginPatterns`添加可信来源；避免使用`WebSocketAllowAllOrigins`
5. **验证所有输入**——遵循与任何 Web 应用程序相同的安全实践

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

这些环境变量的优先级高于代码中的`ServerOptions`，因此 Docker 示例会设置`WAILS_SERVER_HOST=0.0.0.0`——这样无需对应用程序进行任何更改，容器即可接受外部连接。

## 另请参阅

- [自定义传输](/guides/custom-transport/)——用于高级 IPC 自定义
- [服务](/features/bindings/services/)——服务绑定文档
- [事件](/guides/events-reference/)——事件系统文档

### 运行时请求大小

在进行 JSON 处理之前，发送到`/wails/runtime`的请求大小上限为64 MiB。超过此限制的普通请求会收到 HTTP 413，其中包括不含`Content-Length`的请求。分块运行时上传仍采用各自独立的限制：每个分块为1 MiB，组装后的有效负载为64 MiB。适当情况下，请使用应用程序中间件或反向代理设置更低的限制。
