---
title: Server Build
description: Run Wails applications as HTTP servers without a native GUI window
sidebar:
    order: 52
    badge:
        text: Experimental
        variant: caution
---

import { Aside } from '@astrojs/starlight/components';

<Aside type="caution" title="Experimental Feature">
Server mode is experimental and may change in future releases.
</Aside>

Wails v3 supports server mode, allowing you to run your application as a pure HTTP server without creating native windows or requiring GUI dependencies. This enables deploying the same Wails application to servers, containers, and web browsers.

## Overview

Server mode is useful for:

- **Docker/Container deployments** - Run without X11/Wayland dependencies
- **Server-side applications** - Deploy as a web server accessible via browser
- **Web-only access** - Share the same codebase between desktop and web
- **CI/CD testing** - Run integration tests without a display server
- **Microservices** - Use Wails bindings in headless backend services

## Quick Start

Server mode is enabled via the `server` build tag. Your application code remains the same - you just build with the tag:

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

Here's a minimal example:

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

The same code can be built for desktop (without the tag) or server mode (with `-tags server`).

## Configuration

### ServerOptions

Configure the HTTP server with `ServerOptions`:

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

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## Features

### Health Check Endpoint

A health check endpoint is automatically available at `/health`:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

This is useful for:
- Kubernetes liveness/readiness probes
- Load balancer health checks
- Monitoring systems

### Service Bindings

All service bindings work identically to desktop mode:

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

The frontend can call these bindings using the standard Wails runtime:

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### Events

Events work bidirectionally in server mode:

- **Frontend to Backend**: Events emitted from the browser are sent via HTTP and received by your Go event handlers
- **Backend to Frontend**: Events emitted from Go are broadcast to all connected browsers via WebSocket

Each browser tab is represented as a "window" with a unique name (`browser-1`, `browser-2`, etc.), accessible via `event.Sender`:

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

From the frontend:

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### Graceful Shutdown

The server handles `SIGINT` and `SIGTERM` signals gracefully:

1. Stops accepting new connections
2. Waits for active requests to complete (up to `ShutdownTimeout`)
3. Runs `OnShutdown` hooks
4. Shuts down services in reverse order

## Differences from Desktop Mode

| Feature | Desktop Mode | Server Mode |
|---------|-------------|-------------|
| Native windows | Created | Browser windows (`browser-N`) |
| System tray | Available | Not available |
| Native dialogs | Available | Not available |
| Application menu | Available | Not available |
| Screen info | Available | Returns error |
| Service bindings | Works | Works |
| Events | Works | Works (via WebSocket) |
| Assets | Via webview | Via HTTP |
| CGO required | Yes | No |

### Window API Behavior

In server mode, window-related APIs are safely handled:

- `app.Window.NewWithOptions()` - Logs a warning, returns nil
- `app.Hide()` / `app.Show()` - No-op
- `app.Screen.GetPrimary()` - Returns error

This allows code that references windows to run without crashing, though window operations have no effect.

## Building for Production

### Using Task (Recommended)

Projects created with `wails3 init` include a `build:server` task:

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### Manual Build

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Wails projects include a ready-to-use Docker setup. To build and run your application in a container:

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

That's it! Your application will be available at `http://localhost:8080`.

You can customise the build with a few options:

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

The generated `Dockerfile.server` creates a minimal image based on distroless. It handles the network binding automatically, so your application will be accessible from outside the container.

### Docker Compose

For more complex deployments, here's a Docker Compose configuration with health checks:

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

:::note
The healthcheck example uses `wget`. If you're using a distroless base image, you'll need to either include a healthcheck binary in your image or use an external health check mechanism (e.g., Docker's `curl` option or a sidecar container).
:::

### Custom Dockerfile

If you need more control, you can create your own Dockerfile. The key thing to remember is setting `WAILS_SERVER_HOST=0.0.0.0` so the server accepts connections from outside the container:

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

## Security Considerations

When deploying server mode applications:

1. **Bind to localhost by default** - Only use `0.0.0.0` when needed
2. **Use TLS in production** - Configure `ServerOptions.TLS`
3. **Place behind reverse proxy** - Use nginx/traefik for additional security
4. **Validate all inputs** - Same security practices as any web application

## Example

A complete example is available at `v3/examples/server/`:

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## Environment Variables

For deployment scenarios where you need to override the server configuration without changing code, Wails recognises these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `WAILS_SERVER_HOST` | Network interface to bind to | `localhost` |
| `WAILS_SERVER_PORT` | Port to listen on | `8080` |

These take precedence over the `ServerOptions` in your code, which is why the Docker examples set `WAILS_SERVER_HOST=0.0.0.0` - it allows the container to accept external connections without requiring any changes to your application.

## See Also

- [Custom Transport](/docs/guides/custom-transport) - For advanced IPC customization
- [Services](/docs/concepts/services) - Service binding documentation
- [Events](/docs/guides/events-reference) - Event system documentation
