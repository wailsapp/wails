# Native titlebar accessories

This example is a plain single-WebView window whose titlebar carries three
native AppKit control strips built with `MacAccessory`. Nothing in the
titlebar is HTML: every control is a real `NSSearchField`,
`NSSegmentedControl`, `NSButton`, or `NSTextField` hosted by an
`NSTitlebarAccessoryViewController`, so it gets the system's titlebar
material, full-screen behavior, and scroll-edge treatment for free.

The window shows a small mailbox. The accessories filter, sort, and summarise
it from Go; the page only renders what it is sent.

It demonstrates:

- `MacAccessoryLayoutLeading`: a segmented folder switcher
  (`AddSegmented` with `SetSegmentSymbols`) next to the window buttons;
- `MacAccessoryLayoutTrailing`: an incremental search field (`AddSearch`,
  `SetIncremental`, `SetWidth`), a sort menu built from an ordinary `Menu`
  (`AddMenuButton`), and a symbol-only compose button (`AddSymbolButton`);
- `MacAccessoryLayoutBottom`: a full-width status strip with a secondary
  label that carries an SF Symbol (`AddLabel`, `SetSymbol`, `SetText`), a
  flexible space, and two push buttons; its height is set with `SetHeight`
  and it asks for a soft scroll-edge effect through
  `SetPreferredScrollEdgeEffectStyle`;
- live state on the handles: `SetEnabled` on the Mark All Read button,
  `SetText` on the status label and the search field;
- attaching before the window exists (queued and installed with the window);
- the accessory lifecycle from the View menu: `SetHidden` collapses the
  status strip in place, and `Remove` followed by another
  `AddTitlebarAccessory` detaches and reattaches the search tools.

`MacAccessory.Controller()` exposes the type-checked
`MacAccessoryViewController` wrapper once attached, for native code that
needs the AppKit controller itself.

## Run the demo

```sh
GOWORK=off go run .
```

Choose a folder with the segmented control, type in the search field to
filter as you go, pick a sort order from the menu button, and use
**View → Toggle Status Bar** (Cmd+/) and **View → Toggle Search Tools**
(Cmd+Shift+F) to hide, detach, and reattach the strips.

## Status

| Platform | Status |
|----------|--------|
| Mac      | Working |
| Windows  | N/A (`AddTitlebarAccessory` returns `ErrMacAccessoryUnsupported`) |
| Linux    | N/A (`AddTitlebarAccessory` returns `ErrMacAccessoryUnsupported`) |
