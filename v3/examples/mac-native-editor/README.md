# Native Notes

Native Notes is an experimental, WebView-free Wails v3 application. It opens
its editor once on launch, then remains available from a macOS status item. The
close button hides the window so the tray icon can reopen it. The editor uses
only AppKit:

- `NSWindow`
- `NSToolbar` and `NSSearchToolbarItem`
- `NSSplitViewController`
- `NSOutlineView` source-list sidebar
- `NSScrollView` and `NSTextView`

Run it from the v3 directory:

```sh
GOWORK=off go run -tags wails_native ./examples/mac-native-editor
```

The demo copies its three embedded text files into a temporary directory and
prints that directory on startup. Edits made with the toolbar or Command-S are
written to those files.

For a quick process-level sample, find the compiled process and inspect it over
several seconds:

```sh
pid=$(pgrep -n mac-native-editor)
ps -o pid,rss,%cpu,time,command -p "$pid"
top -l 5 -pid "$pid" -stats pid,command,cpu,mem,threads,wq
```

There should be no WebKit WebContent child process associated with Native
Notes. Compare the same measurements with `examples/mac-toolbar` to quantify
the cost of the WKWebView path on the same machine and OS release.

Use the `wails_native` build tag for a genuinely native binary. `NativeOnly`
without the tag suppresses frontend startup at runtime, but the monolithic
application package still links WebKit and includes the frontend and updater
dependency graphs.

Single-instance support is also omitted from the lean native build. Add
`wails_single_instance` when it is required:

```sh
GOWORK=off go run -tags 'wails_native wails_single_instance' \
  ./examples/mac-native-editor
```

For the controlled Swift/AppKit comparison, build both release binaries and
run the complete launch, idle, hidden-window, and document-scaling matrix:

```sh
examples/mac-native-editor/benchmark/run.sh
```

The script prints its results directory. `LAUNCH_ITERATIONS` and `IDLE_SECONDS`
can reduce or extend the default 20 launches and 30-second idle observations.
The checked methodology and August 2026 native/Swift result are documented in
`wep/proposals/macos-native-window/performance.md`. The wider investigation of
normal Wails builds, package reachability, WebKit process cost, GC behaviour,
and optimization priorities is in
`wep/proposals/macos-native-window/performance-investigation.md`.
