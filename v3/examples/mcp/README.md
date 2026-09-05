# MCP Service Example

A playground application that an LLM can test and control through the
[Model Context Protocol](https://modelcontextprotocol.io) using the
built-in MCP server (`v3/pkg/application`, enabled with `-tags mcp`).

The page contains a counter, a text input wired to a bound Go method,
a drag & drop area and a scrollable list — one widget for each MCP tool
category. Simulated mouse input is rendered with an animated cursor inside
the window, so you can watch the LLM drive the app.

## Running

The MCP service is compiled in only with the `mcp` build tag:

```shell
go run -tags mcp .
```

In a real project the tag is added automatically by setting an environment
variable when building:

```shell
WAILS_MCP=1 wails3 build   # or: wails3 build -tags mcp
WAILS_MCP=1 wails3 dev
```

Without the tag the service is a no-op stub and the app behaves normally.

## Connecting a client

The server listens on `http://127.0.0.1:9099/mcp` (streamable HTTP).
Configure the port with the `WAILS_MCP_PORT` environment variable.

For Claude Code:

```shell
claude mcp add --transport http mcp-demo http://127.0.0.1:9099/mcp
```

Then ask things like:

- "List the windows of the app and take a DOM snapshot."
- "Click the Increment button three times and verify the counter shows 3."
- "Type 'Wails' into the name input, click Greet and tell me the result."
- "Drag the blue box onto the drop target."
- "Call the bound method main.GreetService.Add with 2 and 40."

## Available tools

`app_info`, `windows_list`, `window_control`, `js_eval`, `dom_html`,
`dom_query`, `screenshot_dom`, `mouse_move`, `mouse_click`, `mouse_drag`,
`mouse_scroll`, `keyboard_type`, `keyboard_press`, `call_bound_method`,
`emit_event`, `wait_for_event`.
