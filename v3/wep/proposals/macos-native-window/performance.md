# Native editor performance: Wails versus Swift/AppKit

This document contains the controlled native-editor measurements. The broader
package, linker, WebKit, garbage-collector, and normal-Wails analysis is in
[`performance-investigation.md`](performance-investigation.md).

## Result

The experimental native Wails path is viable for a small tray utility. With a
normal text document it is genuinely idle, creates no WebKit helper process or
network listener, and has a physical footprint about 4.7 MiB above an
equivalent Swift/AppKit executable at the same 980×680 window size.

The current whole-string text bridge is not yet suitable for large documents.
At 10 MiB, Wails settles about 25.4 MiB above Swift. This is a memory-scaling
issue, not a CPU or responsiveness issue: launch-and-load CPU time remained
roughly comparable.

## Control application

`examples/mac-native-editor-swift/main.swift` is a standalone Swift control. It
uses the same visible AppKit structure and behaviour as the Wails example:

- accessory application and status item;
- 980×680 `NSWindow`, unified `NSToolbar`, sidebar toggle, tracking separator,
  `NSSearchToolbarItem`, and Save item;
- `NSSplitViewController`, source-list `NSOutlineView`, `NSScrollView`, and
  editable `NSTextView`;
- the same sample documents, generated benchmark document, save, search,
  first-responder, and close-to-hide behaviour.

Neither executable was bundled or signed. The benchmark validates process and
framework cost, not application bundle launch services.

## Test system and builds

- Apple M4 Max, 36 GiB RAM
- macOS 26.4.1 (25E253), arm64
- Xcode 26.6
- Go 1.26.2
- Swift 6.3.3

Release commands:

```sh
GOWORK=off go build -tags 'production wails_native' \
  -trimpath -ldflags='-s -w' \
  -o native-notes-wails ./examples/mac-native-editor

swiftc -O -whole-module-optimization \
  examples/mac-native-editor-swift/main.swift \
  -framework AppKit -o native-notes-swift
```

The repeatable harness is `examples/mac-native-editor/benchmark/run.sh`. It
uses `proc_pid_rusage(RUSAGE_INFO_V4)` for physical footprint, CPU, wakeups,
page-ins, instructions, and cycles. It uses `/usr/bin/time -l` for launch peak
footprint and `top` for threads and ports.

## Measurements

Physical footprint is the primary memory measure. RSS is not used for the
comparison because it includes shared framework mappings and makes both tiny
AppKit applications appear to consume more than 100 MiB.

### Steady state after settling

Each row is one 30-second observation after a two-second settling period.

| State | Wails footprint | Swift footprint | Difference | Wails threads | Swift threads |
|---|---:|---:|---:|---:|---:|
| Visible, sample note | 25.84 MiB | 21.16 MiB | +4.68 MiB | 12 | 6 |
| Hidden, sample note | 24.45 MiB | 19.81 MiB | +4.64 MiB | 10 | 7 |
| Visible, 1 MiB text | 32.59 MiB | 24.30 MiB | +8.29 MiB | 14 | 6 |
| Visible, 10 MiB text | 59.84 MiB | 34.42 MiB | +25.42 MiB | 14 | 7 |

All observations used less than 6 ms of CPU over 30 seconds. Most used less
than 1.3 ms. AppKit timer wakeups varied with focus and text layout, especially
in Swift, so the exact counts are not stable enough to rank the runtimes. Both
were effectively idle; observed totals ranged from single digits to roughly
500 wakeups over 30 seconds.

### Repeated launch and document load

Ten runs per document size. The auto-quit timer starts after AppKit has
finished launching and the complete interface has been built, then deliberately
keeps the process alive for one second. Elapsed values therefore include that
one-second dwell. CPU and peak footprint are medians.

| Document | Wails elapsed | Swift elapsed | Wails CPU | Swift CPU | Wails peak | Swift peak |
|---|---:|---:|---:|---:|---:|---:|
| Sample note | 1.11 s | 1.24 s | 0.22 s | 0.24 s | 25.73 MiB | 21.37 MiB |
| 1 MiB | 1.11 s | 1.25 s | 0.23 s | 0.24 s | 32.88 MiB | 24.34 MiB |
| 10 MiB | 1.11 s | 1.27 s | 0.23 s | 0.26 s | 63.69 MiB | 46.19 MiB |

These are repeated warm launches; filesystem and framework caches were not
purged. The result supports “no material load-time regression,” not a claim
that Wails universally launches faster than Swift.

### Binary and process shape

| Metric | Wails | Swift |
|---|---:|---:|
| Stripped executable | 3.33 MiB | 0.17 MiB |
| Child processes | 0 | 0 |
| Network listeners/connections | 0 | 0 |

The lean native Wails binary is about 20 times larger. It does not link WebKit;
the `wails_native` build boundary also removes the frontend transport, asset
server, updater, and default encrypted single-instance dependency graph. The
Swift control links AppKit/Foundation and the system Swift runtime, while Go
must include its runtime and metadata in the executable. Launching either
control created no WebContent or Networking helper process.

## Problems found and fixed while profiling

1. `MacTextEditor` retained a complete Go copy after the text had been accepted
   by `NSTextView`. Attached editors now treat AppKit as authoritative and keep
   Go text only as pre-install staging state.
2. A stale native-acceptance completion could have cleared a newer `SetText`.
   Text staging is now versioned and only the accepted version is released.
3. Native installation seeded the editor and then replayed identical text to
   defend against an installation race. Replay is now conditional on the text
   version changing. This reduced the 10 MiB lifetime peak from about 82 MiB to
   about 66 MiB.
4. Initial focus was requested before the `NSWindow` existed. The native window
   now focuses its primary editor after becoming key.
5. The Swift control initially adopted the split controller's 640×420 fitting
   size. It now restores 980×680 after attaching the hierarchy, matching Wails.

## Interpretation and next work

For a tray app, settings utility, small editor, or other ordinary native
window, the current overhead is acceptable: about 4.6–4.7 MiB physical memory,
several Go runtime threads, negligible idle CPU, and no WebKit processes.

Two optimizations should precede claiming production-quality large-document
support:

1. Add a file-oriented editor path (`LoadFile`/`SaveFile`, or an equivalent
   `NSData`/length-aware provider) so Foundation can load and save without
   `[]byte -> string -> C string -> NSString` whole-document copies.
2. Continue reducing the native build graph. The new `wails_native` boundary
   already excludes WebKit, frontend transport, the asset server, updater, and
   optional encrypted single-instance support. The remaining executable is
   3.33 MiB, compared with a measured 1.13 MiB minimal Go+AppKit floor.

The benchmark does not cover sustained typing, syntax highlighting, attributed
text, multiple large documents, autosave, or `NSDocument`. Those belong to the
next editor/content API experiment rather than this minimal native-window test.
