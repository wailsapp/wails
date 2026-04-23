# systray-stress

Stress harness for `SystemTray.SetMenu` and `SystemTray.OpenMenu` on Windows.
Doubles as a regression test for the handle-leak and crash fixes that go with
this example.

## Workloads

- `-mode churn` — background goroutine loops `SetMenu(buildMenu(...))` without
  ever showing a popup. Isolates per-rebuild handle churn. Watch the
  `handles_delta` / `gdi_delta` columns on stderr.
- `-mode show` — opens and dismisses the tray popup on a timer using
  synthetic `VK_ESCAPE` via `SendInput`. Baseline; should never crash. Note
  that `SendInput` targets whatever window holds focus, so don't leave this
  running during interactive work.
- `-mode churn+show` — both workloads concurrently. Exercises the
  `SetMenu`-during-`TrackPopupMenuEx` interaction: `InvokeSync` posted from
  the churn goroutine is dispatched by the modal message loop that
  `TrackPopupMenuEx` runs.

Every `-log-every` iterations (default 500) the harness logs
`GetGuiResources(GR_USEROBJECTS)` and `GR_GDIOBJECTS` deltas. It exits
cleanly when either delta exceeds `-handle-cap` (default 5000) or iterations
reach `-iters`, whichever comes first.

## Flags

| Flag | Default | Purpose |
|---|---|---|
| `-mode` | `churn` | `churn`, `show`, or `churn+show` |
| `-iters` | `50000` | max SetMenu iterations before clean exit (`0` = unbounded) |
| `-handle-cap` | `5000` | exit if user-object or GDI delta exceeds this |
| `-churn-gap` | `2ms` | sleep between SetMenu calls |
| `-show-gap` | `80ms` | sleep between OpenMenu calls |
| `-dismiss-gap` | `30ms` | delay between popup open and synthetic ESC |
| `-log-every` | `500` | emit a progress line every N iterations |
| `-duration` | `0` | wall-clock cap (`0` = unbounded) |
| `-bitmaps` | `false` | attach a bitmap icon to a subset of menu items — exercises the GDI path |

## Running

```pwsh
# From v3/
go build -o harness.exe ./examples/systray-stress
.\harness.exe -mode churn -duration 30s -churn-gap 10ms -bitmaps
```

Expected on a fixed build: `handles_delta` and `gdi_delta` stay within a
small bounded range (tens) across tens of thousands of iterations, with no
`error adding menu item` fatal.

## For manual testing of the popup-open race

Run with `-mode churn` and right-click the tray icon to open the popup
yourself. The label on the first item (`Iter #N`) ticks with each rebuild,
so you can see updates landing. The popup should stay responsive even while
the churn goroutine is firing SetMenu calls through the modal loop.
