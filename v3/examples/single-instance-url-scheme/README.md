# single-instance-url-scheme — reproduction for issue #5089

Reproduces [wailsapp/wails#5089](https://github.com/wailsapp/wails/issues/5089):
when `SingleInstance` is combined with a custom URL scheme on macOS, a
URL-scheme launch that arrives while a first instance is already running
is **dropped** — neither `OnSecondInstanceLaunch` nor
`ApplicationLaunchedWithUrl` sees it.

This example is the **failing** test case. No fix has been applied yet;
running these steps is expected to demonstrate the bug.

## Why this happens (short form)

On macOS, URL-scheme launches do not place the URL in `os.Args`.
LaunchServices delivers the URL via an Apple Event (`kAEGetURL`) that is
dispatched by the target process's `NSAppleEventManager` — which is only
wired up once `NSApplication.run` is executing.

When a second instance hits the flock in `newSingleInstanceManager` it
immediately calls `notifyFirstInstance()` and `os.Exit`. The Apple Event
handler is never installed, so the URL is discarded. The payload relayed
to the first instance (`SecondInstanceData{Args, WorkingDir, …}`) only
contains what was in `os.Args`, which on macOS does not include the URL.

On Windows and Linux the URL is passed through `argv`, so it surfaces
naturally as `SecondInstanceData.Args[1]`. macOS is the odd one out.

## Reproducing the bug

Requirements: macOS, Go, `wails3` CLI (`task` runner), Xcode command-line
tools. Uses `codesign --sign -` (ad-hoc) so no certificate needed.

```sh
cd v3/examples/single-instance-url-scheme

# 1. Build + package a dev .app bundle and run it. This also registers
#    the bundle's CFBundleURLTypes with LaunchServices so the custom
#    scheme is routed to the app.
wails3 task run

# In a separate terminal, tail the log the app writes:
tail -f /tmp/wails-single-instance-url.log
```

With the first instance window visible, trigger a URL-scheme launch:

```sh
# Either via the helper task:
wails3 task trigger URL='wails-single-url://hello?n=1'

# Or directly:
open 'wails-single-url://hello?n=1'
```

### macOS version behaviour differences

**macOS 14 / 15 (original report):** `open 'wails-single-url://...'` causes
LaunchServices to spawn a *new* process (because `LSMultipleInstancesProhibited`
is not set). That second process detects the flock, relays `os.Args`, and exits
before the Apple Event handler is registered. The URL is dropped.

**macOS 26+ (observed on 26.0/25A354, Apple Silicon):** LaunchServices routes
the Apple Event *directly* to the already-running instance without spawning a
new process. `ApplicationLaunchedWithUrl` fires correctly — the bug is not
visible via plain `open` on this OS version.

To reproduce the bug on macOS 26+, use `trigger:force` which forces a new
process via `open -n`:

```sh
wails3 task trigger:force URL='wails-single-url://hello?n=1'
```

### Expected (desired) behaviour

The first instance's log contains either a line from
`OnSecondInstanceLaunch` with the URL visible in `Args`, **or** a line
from `ApplicationLaunchedWithUrl`, or both.

### Actual (buggy) behaviour

Triggered via `trigger:force` (or via plain `trigger` on macOS 14/15):
you will see a second process start and exit (pid differs from the first
instance's pid). The first instance logs:

```
[first] OnSecondInstanceLaunch fired
[first]   Args           = [.../single-instance-url-scheme]
[first]   url-in-args?   = false  (url="")
```

`ApplicationLaunchedWithUrl` does **not** fire on either instance.

The UI mirrors this: the `ApplicationLaunchedWithUrl` row stays
`(never fired)`, and the `OnSecondInstanceLaunch` row shows
`"found": false`.

## What a fix must demonstrate

When re-running the same steps against a fixed Wails build, the first
instance must receive the URL `wails-single-url://hello?n=1` via one of
the existing callbacks (exact surface is still open — see issue #5089).

## Notes / gotchas for testing

- Running the raw binary (`go run .`) will **not** exercise the bug path.
  macOS only routes custom URL schemes to apps launched from a `.app`
  bundle registered with LaunchServices. Use `wails3 task run`.
- If `open 'wails-single-url://…'` launches a different app, the scheme
  is claimed by a previously-registered bundle. Re-register this one:
  `lsregister -f bin/single-instance-url-scheme.dev.app`.
- Fresh launches (first instance) already go through the
  `NSAppleEventManager` path and work correctly; the bug is specific to
  the second-instance relay.
- The issue was originally reported against `v3.0.0-alpha.67`. The code
  paths involved (`v3/pkg/application/single_instance.go`,
  `single_instance_darwin.go`, `application_darwin_delegate.m`) are
  unchanged on `v3-alpha` HEAD at the time this example was written.
