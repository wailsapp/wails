# Mac Split Sidebar Example

This example demonstrates `application.NewSplitWindow`: a real `NSSplitViewController`
window with three independent panes, each its own `WKWebView`, not a single webview with
absolutely positioned content.

A Reader app with:

- A **sidebar** pane (`PaneBehaviorSidebar`) listing articles
- A **content** pane (`PaneBehaviorDefault`) showing the selected article
- An **inspector** pane (`PaneBehaviorInspector`) showing its metadata
- A toolbar `ToolbarSidebarToggle` item that collapses/expands the sidebar via AppKit's
  built-in `toggleSidebar:` responder-chain action — no extra wiring needed
- `AutosaveName` so the divider positions you drag persist across relaunches

Selecting an article in the sidebar emits a `reader:select` event; the content and inspector
panes, running in their own independent webviews, pick it up via `Events.On`. This is the
app's normal event bus, not special split-window plumbing: every pane routes IPC through the
same window ID as the window's own webview.

Sidebar and inspector get their translucent material background from AppKit's
`+sidebarWithViewController:` / `+inspectorWithViewController:` factory methods — nothing in
this example's Go or HTML sets that up.

## Running the example

```bash
go run .
```

# Status

| Platform | Status  |
|----------|---------|
| Mac      | Working |
| Windows  | N/A (MacSplitWindow is macOS only) |
| Linux    | N/A (MacSplitWindow is macOS only) |
