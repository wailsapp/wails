# Server Mode Example

> **Experimental** - This feature is experimental and may change in future releases.

This example demonstrates running a Wails application in server mode - without a native GUI window.

## What is Server Mode?

Server mode allows you to run your Wails application as a pure HTTP server. This enables:

- **Docker/Container deployments** - Deploy your Wails app without X11/Wayland dependencies
- **Server-side rendering** - Use your Wails app as a web server
- **Web-only access** - Share the same codebase between desktop and web deployments
- **CI/CD testing** - Run integration tests without a display server

## Building and Running

The recommended way to build server mode applications is using the Taskfile:

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

Or using Go directly:

```bash
# Build with server tag
go build -tags server -o myapp-server .

# Run
go run -tags server .
```

Then open http://localhost:8080 in your browser.

## Key Differences from Desktop Mode

1. **No native window** - The app runs as an HTTP server only
2. **Browser access** - Users access the app via their web browser
3. **No CGO required** - Can build without CGO dependencies
4. **Window APIs are no-ops** - Calls to window-related APIs are safely ignored
5. **Browser windows** - Each browser tab is represented as a "window" named `browser-1`, `browser-2`, etc.

## Events

Events work bidirectionally in server mode:

- **Frontend to Backend**: Events emitted from the browser are sent via HTTP and received by your Go event handlers
- **Backend to Frontend**: Events emitted from Go are broadcast to all connected browsers via WebSocket

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all browsers
app.Event.Emit("server-update", data)
```

## Configuration

Server mode is enabled by building with the `server` build tag. Configure the HTTP server options:

```go
app := application.New(application.Options{
    // Configure the HTTP server (used when built with -tags server)
    Server: application.ServerOptions{
        Host: "localhost",  // Use "0.0.0.0" for all interfaces
        Port: 8080,
    },

    // ... other options work the same as desktop mode
})
```

## Health Check

A health check endpoint is automatically available at `/health`:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Building for Production

```bash
# Using Taskfile (recommended)
wails3 task build:server

# Or using Go directly
go build -tags server -o myapp-server .
```

## Docker

Build and run with Docker using the built-in tasks:

```bash
# Build Docker image
task build:docker

# Build and run
task run:docker

# Run on a different port
task run:docker PORT=3000
```

Or build manually:

```bash
docker build -t server-example .
docker run --rm -p 8080:8080 server-example
```

## Limitations

Since server mode runs without a native GUI, the following features are not available:

- Native dialogs (file open/save, message boxes)
- System tray
- Native menus
- Window manipulation (resize, move, minimize, etc.)
- Clipboard access (use browser clipboard APIs instead)
- Screen information

These APIs are safe to call but will have no effect or return default values.
