## Systray

Wails 3 comes with a built-in systray. This is a fully featured systray that has
been designed to be as simple as possible to use. It is possible to set the
icon, tooltip and menu of the systray. It is possible to also "attach" a window
to the systray. Doing this will provide the following functionality:

- Clicking the systray icon with toggle the window visibility
- Right-clicking the systray will open the menu, if there is one

On macOS, if there is no attached window, the systray will use the default
method of displaying the menu (any button). If there is an attached window but
no menu, the systray will toggle the window regardless of the button pressed.