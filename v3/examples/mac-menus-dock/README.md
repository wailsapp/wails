# macOS menus and Dock

This example exercises the macOS-only menu and Dock extras from a plain Wails
window. Everything shown is native AppKit; the window content only logs what
happens.

It demonstrates:

- `MenuItem.SetSymbol` to put SF Symbol images on menu items (macOS 11+);
- `MenuItem.SetBadge` / `SetBadgeText` (`NSMenuItemBadge`, macOS 14+);
- `Menu.AddSectionHeader` (`NSMenuItem sectionHeaderWithTitle:`, macOS 14+);
- `Menu.AddPalette`, a colour palette submenu built with
  `NSMenu paletteMenuWithColors:` (macOS 14+);
- `MenuItem.SetMixed` for the partially checked state, `SetAlternate` for an
  Option-key alternate of the item above it, and `SetIndentationLevel`;
- `app.Menu.SetDockMenu` and `app.Menu.OnDockMenu` for a static or dynamic
  right-click menu on the Dock icon;
- the `OpenRecent` role plus `app.Menu.AddRecentDocument`,
  `ClearRecentDocuments` and `RecentDocuments` backed by
  `NSDocumentController`; opening an entry emits the usual
  `ApplicationOpenedWithFile` event;
- `dock.SetProgress` / `ClearProgress`, a progress bar drawn on the Dock tile
  that coexists with the badge label.

Features that need a newer macOS than the one running degrade to a logged
no-op (or a plain disabled item for section headers).

## Running

```shell
go run .
```

Open the **Demo** menu, hold Option to reveal the alternate item, right-click
the Dock icon to drive the badge and progress bar, and use **File > Open
Recent** after adding a document.
