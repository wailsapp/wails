# Mac Titlebar Accessory Example

This example demonstrates `WebviewWindow.AddTitlebarAccessory`: independent, own-webview
controls pinned to the titlebar itself, not to the window's content.

A plain single-webview window (no toolbar, no split layout) with:

- A **Ping** button pinned `AccessoryLeading`, in its own webview
- A **status pill** pinned `AccessoryTrailing`, also its own webview, driven entirely from Go
  via `ExecJS` on a timer -- nothing in the main window's content is aware of it

Clicking Ping emits an event that Go logs; the status pill cycles Idle → Syncing → Synced on its
own, proving both accessories render and update independently, without touching the main
window's page at all.

## Running the example

```bash
go run .
```

# Status

| Platform | Status  |
|----------|---------|
| Mac      | Working |
| Windows  | N/A (titlebar accessories are macOS only) |
| Linux    | N/A (titlebar accessories are macOS only) |
