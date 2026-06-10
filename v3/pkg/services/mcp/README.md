# MCP Service

A [Model Context Protocol](https://modelcontextprotocol.io) server for Wails v3
applications. It lets LLM agents (Claude Code, IDE assistants, any MCP client)
test and control a **running** application: windows, DOM, JavaScript, bound Go
methods, events, and simulated mouse/keyboard input rendered with an animated
on-screen cursor.

## Quick start

```go
import "github.com/wailsapp/wails/v3/pkg/services/mcp"

app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(mcp.New(), application.ServiceOptions{
            Route: "/wails-mcp", // optional same-origin callback channel
        }),
    },
})
```

The service only compiles in with the `mcp` build tag; otherwise it is a no-op
stub, so registration can stay in place for all builds:

```shell
WAILS_MCP=1 wails3 dev      # CLI converts the env var into -tags mcp
WAILS_MCP=1 wails3 build
go run -tags mcp .          # plain Go
```

Connect a client to the logged endpoint (default `http://127.0.0.1:9099/mcp`,
streamable HTTP transport):

```shell
claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
```

## Tools

| Tool | Purpose |
|------|---------|
| `app_info`, `windows_list` | Application and window inventory |
| `window_control` | Focus, resize, move, fullscreen, reload, devtools, … |
| `js_eval` | Run JavaScript in a window and return its result |
| `dom_html`, `dom_query`, `screenshot_dom` | Inspect page structure |
| `mouse_move`, `mouse_click`, `mouse_drag`, `mouse_scroll` | Animated mouse input with hover, ripple and drag effects |
| `keyboard_type`, `keyboard_press` | Keyboard input with realistic per-character events |
| `call_bound_method` | Call bound Go methods through the runtime |
| `emit_event`, `wait_for_event` | Work with Wails application events |

## How it works

- A localhost HTTP listener serves the MCP streamable HTTP endpoint (`/mcp`).
- Tools that touch the page inject `inject.js` (idempotent) via `ExecJS` and
  receive results through a `fetch` callback: same-origin via the service
  `Route` when configured, otherwise the CORS-enabled localhost listener.
- Mouse tools animate a synthetic cursor overlay and dispatch the full
  pointer/mouse event sequences; keyboard tools dispatch per-character key and
  input events using native value setters so frameworks observe the changes.
- Animation degrades gracefully when the webview is not rendering (timer-based
  frames), so tools never hang on unfocused windows.

## Security

The server binds to `127.0.0.1` by default and rejects non-local browser
origins. Anything that can reach the port can fully control the application:
keep it on localhost and out of production builds (the default build without
the tag contains no server code).

## Configuration

See `Config` in [mcp.go](./mcp.go): `Host`, `Port` (`WAILS_MCP_PORT` overrides
at runtime, `-1` picks a free port), `EvalTimeout`, `HideCursor`.

A playground app lives at [`v3/examples/mcp`](../../../examples/mcp).
