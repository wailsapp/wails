# Native macOS toolbar

This example is a normal titled Wails window with a real AppKit `NSToolbar` above it.
Nothing in the toolbar is HTML.

It demonstrates:

- a persistent `NSToolbarItemGroup` for Write/Preview mode;
- an `NSSearchToolbarItem` with a Go callback;
- native New, Share, Details, Save, and Focus controls;
- generated internal item identifiers, with no IDs in application code;
- live item handles that update labels, badges, selection, and visibility;
- a working notes editor with a note list, write/preview modes, search, sharing,
  local save state, details, focus mode, and keyboard save.

Run it on macOS from this directory:

```sh
go run .
```
