# Window Menu Example

*** Windows Only ***

This example demonstrates how to create a window with a menu bar that can be toggled using the window.ToggleMenuBar() method.

## Features

- Default menu bar with File, Edit, and Help menus
- F1 key to toggle menu bar visibility 
- Simple HTML interface with instructions

## Running the Example

```bash
cd v3/examples/window-menu
go run .
```

## How it Works

The example creates a window with a default menu and binds the F10 key to toggle the menu bar's visibility. The menu bar will hide when F10 is pressed and show when F10 is released.

Note: The menu bar toggling functionality only works on Windows. On other platforms, the F10 key binding will have no effect.
