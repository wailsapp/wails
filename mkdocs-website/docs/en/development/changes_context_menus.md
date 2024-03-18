## Context Menus

Context menus are contextual menus that are shown when the user right-clicks on
an element. Creating a context menu is the same as creating a standard menu , by
using `app.NewMenu()`. To make the context menu available to a window, call
`window.RegisterContextMenu(name, menu)`. The name will be the id of the context
menu and used by the frontend.

To indicate that an element has a context menu, add the `data-contextmenu`
attribute to the element. The value of this attribute should be the name of a
context menu previously registered with the window.

It is possible to register a context menu at the application level, making it
available to all windows. This can be done using
`app.RegisterContextMenu(name, menu)`. If a context menu cannot be found at the
window level, the application context menus will be checked. A demo of this can
be found in `v3/examples/contextmenus`.