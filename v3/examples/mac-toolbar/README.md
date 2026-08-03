# Mac Toolbar Example

This example demonstrates `application.MacToolbar` / `WebviewWindow.SetToolbar`: a real `NSToolbar`
attached above the webview, not a DOM element under the traffic lights.

A Notes-style app with:

- An Edit/Preview segmented group (`ToolbarGroup`)
- A flexible space (`ToolbarFlexibleSpace`)
- A native search field (`ToolbarSearchField`) whose query is submitted to Go and relayed back to
  the frontend to highlight matches
- A Share button (`ToolbarButton`, bordered)
- An Info button (`ToolbarButton`, bordered + prominent + tinted) whose `BadgeCount` increments on
  every click by rebuilding the toolbar and calling `SetToolbar` again — there is no per-item update
  method, so this is the pattern for any dynamic toolbar state

Every toolbar interaction is logged in the activity panel at the bottom of the window so each Go
callback firing is visible without opening a debugger.

## Running the example

```bash
go run .
```

# Status

| Platform | Status  |
|----------|---------|
| Mac      | Working |
| Windows  | N/A (MacToolbar is a no-op) |
| Linux    | N/A (MacToolbar is a no-op) |
