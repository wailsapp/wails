# global-shortcuts

This example demonstrates how to use the global shortcuts API to register
system-wide keyboard shortcuts that fire even when the application is not
focused.

Press `Ctrl+Shift+G` (`Cmd+Shift+G` on macOS) from any application to bring the
window to the front, and `Ctrl+Shift+H` (`Cmd+Shift+H`) to hide it. The example
also shows that registering the same shortcut twice is rejected.

# Status

| Platform | Status  |
|----------|---------|
| Mac      | Working (Carbon hot keys) |
| Windows  | Working (`RegisterHotKey`) |
| Linux X11 | Working (`XGrabKey`) |
| Linux Wayland | Implemented via the XDG global shortcuts portal. Requires a desktop that provides the portal (recent GNOME or KDE Plasma). On Wayland the compositor, and ultimately the user, chooses the final key binding. |
