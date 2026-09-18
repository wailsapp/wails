# Clipboard and Drag (macOS)

Demonstrates three pieces of macOS platform integration:

- **Rich clipboard**: `app.Clipboard.SetImage/Image`, `SetFiles/Files`,
  `SetHTML/HTML`, `SetRTF/RTF`, `SetData/Data` for custom type identifiers,
  `Types`, `Clear`, `ChangeCount` and `OnChange` (polled every 500 ms while a
  listener is registered).
- **Drag out**: `window.StartDrag(application.DragItems{...})` starts a system
  drag from the window with files, file promises (written on demand when the
  destination accepts the drop) or text, with an optional PNG drag image.
  `window.OnDragEnd` reports the operation the destination performed.
- **Non-file drops**: `WebviewWindowOptions.DropTypes` adds `DropText`,
  `DropURLs` and `DropImages` to the default `DropFiles`; those drops are
  delivered to `window.OnDrop` while files keep using the
  `WindowFilesDropped` event.

## Running

```shell
go run .
```

Drag the cards into Finder, Mail or a text editor. Drop a text selection, a
link from Safari or an image from Photos onto the drop zone. Use the buttons
to put different content on the clipboard and Inspect to read it back; change
the clipboard from another app to see `OnChange` fire.

## Drag out and the mouse gesture

The OS only lets an application start a drag while a mouse button is down, so
`StartDrag` must be called during a gesture: the page's `mousedown` (or
`pointerdown`/`dragstart`) handler calls a bound Go method which calls
`StartDrag`. Set `draggable="false"` on the element so WebKit does not start
its own HTML5 drag at the same time. Calling `StartDrag` with no button held
down returns `application.ErrDragOutNoGesture`.

## Drop types and the page

Registering `DropText`, `DropURLs` or `DropImages` routes those drops to Go
instead of the page. The page's own HTML5 `drop` handlers no longer receive
them; if a page needs to handle text or image drops itself, leave those
types out.

Off macOS the rich clipboard methods return
`application.ErrClipboardNotSupported`, `StartDrag` returns
`application.ErrDragOutUnsupported` and non-file drop types are ignored.
