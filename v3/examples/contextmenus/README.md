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

