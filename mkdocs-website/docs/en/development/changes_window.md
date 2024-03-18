## Window

The Window API has largely remained the same, however the methods are now on an
instance of a window rather than the runtime. Some notable differences are:

- Windows now have a Name that identifies them. This is used to identify the
  window when emitting events.
- Windows have even more methods on the that were previously unavailable, such
  as `AbsolutePosition` and `ToggleDevTools`.
- Windows can now accept files via native drag and drop. See the Drag and Drop
  section for more details.

### BackgroundColour

In v2, this was a pointer to an `RGBA` struct. In v3, this is an `RGBA` struct
value.

### WindowIsTranslucent

This flag has been removed. Now there is a `BackgroundType` flag that can be
used to set the type of background the window should have. This flag can be set
to any of the following values:

- `BackgroundTypeSolid` - The window will have a solid background
- `BackgroundTypeTransparent` - The window will have a transparent background
- `BackgroundTypeTranslucent` - The window will have a translucent background

On Windows, if the `BackgroundType` is set to `BackgroundTypeTranslucent`, the
type of translucency can be set using the `BackdropType` flag in the
`WindowsWindow` options. This can be set to any of the following values:

- `Auto` - The window will use an effect determined by the system
- `None` - The window will have no background
- `Mica` - The window will use the Mica effect
- `Acrylic` - The window will use the acrylic effect
- `Tabbed` - The window will use the tabbed effect