# Native Notes Swift control

This standalone AppKit program is the control build for the Wails native-window
performance experiment. It intentionally mirrors `../mac-native-editor`:

- accessory application with a status item;
- `NSWindow`, `NSToolbar`, and `NSSearchToolbarItem`;
- `NSSplitViewController` with an `NSOutlineView` source list;
- `NSScrollView` containing an editable `NSTextView`;
- the same sample documents, save action, filtering, and hide-on-close behaviour.

Build a release binary from the v3 directory:

```sh
swiftc -O -whole-module-optimization \
  examples/mac-native-editor-swift/main.swift \
  -framework AppKit -o /tmp/native-notes-swift
```

Both the Swift and Wails controls accept these benchmark-only environment
variables. They do not affect normal Wails API behaviour:

- `NATIVE_EDITOR_BENCH_HIDDEN=1` starts with only the status item visible.
- `NATIVE_EDITOR_BENCH_DOCUMENT_BYTES=N` creates and opens an N-byte text file.
- `NATIVE_EDITOR_BENCH_AUTO_QUIT_MS=N` exits N milliseconds after launch.
- `NATIVE_EDITOR_BENCH_READY_FILE=path` writes a marker after AppKit finished
  launching and the interface has been built.
- `NATIVE_EDITOR_BENCH_FORCE_GC=1` is a Wails-only diagnostic that forces Go to
  return reclaimable pages after launch. It is not used in the primary result.
