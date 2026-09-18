# contextmenus

This example shows how to create a context menu for your application.
It demonstrates window level and global context menus.

A simple menu is registered with the window and the application with the id "test".
In our frontend html, we then use the `--custom-contextmenu` style to attach the menu to an element.
We also use the `--custom-contextmenu-data` style to pass data to the menu callback which can be read in Go.
This is really useful when using components to distinguish between different elements.

```go

```html

<div class="region" id="123abc" style="--custom-contextmenu: test; --custom-contextmenu-data: 1">
    <h1>1</h1>
</div>
<div class="region" id="234abc" style="--custom-contextmenu: test; --custom-contextmenu-data: 2">
    <h1>2</h1>
</div>
```

# Status

| Platform | Status  |
|----------|---------|
| Mac      | Working |
| Windows  | Working |
| Linux    |         |

## Running on macOS with private APIs

This example configures a translucent macOS backdrop. The webview transparency needed to reveal that backdrop requires the `private_mac_apis` build tag. Without it, the example runs with an opaque webview above the native backdrop.

From this example directory, run:

```bash
go run -tags private_mac_apis .
```

Omit `-tags private_mac_apis` to run with public macOS APIs only. The tag has no effect on Windows, Linux, iOS, or Android. See the [shared private API guide](../README.md#private-macos-apis) for production builds and fallback details.
