# single-instance-url-scheme — fix for issue #5089

Demonstrates (and formerly reproduced) [wailsapp/wails#5089](https://github.com/wailsapp/wails/issues/5089):
combining `SingleInstance` with a custom URL scheme on macOS.

## What the bug was

On macOS, URL-scheme launches do not place the URL in `os.Args`.
LaunchServices delivers the URL via an Apple Event (`kAEGetURL`) that is
dispatched by the target process's `NSAppleEventManager` — which is only
wired up once `NSApplication.run` is executing.

When a second instance hits the flock in `newSingleInstanceManager` it
immediately called `notifyFirstInstance()` and `os.Exit`. The Apple Event
handler was never installed, so the URL was discarded. The payload relayed
to the first instance (`SecondInstanceData{Args, WorkingDir, …}`) only
contained what was in `os.Args`, which on macOS does not include the URL.

On Windows and Linux the URL is passed through `argv`, so it surfaces
naturally as `SecondInstanceData.Args[1]`. macOS was the odd one out.

## The fix

`v3/pkg/application/single_instance_darwin_url.go` adds a `captureLaunchURL()`
helper that:

1. Creates `[NSApplication sharedApplication]` with `NSApplicationActivationPolicyProhibited` (no dock icon).
2. Registers a `kAEGetURL` Apple Event handler **before** calling `[NSApp run]`.
3. Calls `[NSApp run]` — which triggers `finishLaunching`, signalling to
   LaunchServices that this process is ready to receive Apple Events.
4. Stops the run loop immediately when the URL event arrives (or after a
   300 ms safety-net timeout if no event arrives).
5. Returns the captured URL (or `""` on timeout).

`notifyFirstInstance()` (in `single_instance.go`) calls `captureLaunchURL()`
on darwin and appends the URL to `SecondInstanceData.Args` before notifying
the first instance, matching the Windows/Linux behaviour.

`ApplicationLaunchedWithUrl` is **not** fired on the second-instance relay
path, consistent with Windows and Linux behaviour.

## Running the example

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

**macOS 14 / 15:** `open 'wails-single-url://...'` causes LaunchServices to
spawn a second process (because `LSMultipleInstancesProhibited` is not set).
That second process detects the flock, captures the URL via the fix, and
relays it to the first instance.

**macOS 26+ (observed on 26.0/25A354, Apple Silicon):** `open 'wails-single-url://…'`
without `-n` routes the Apple Event directly to the running first instance
(`ApplicationLaunchedWithUrl` fires). No second process is spawned.

To trigger the second-instance relay path on macOS 26+, use `trigger:force`
which forces a new process via `open -n`:

```sh
wails3 task trigger:force URL='wails-single-url://hello?n=1'
```

### Fixed behaviour

After applying the fix, triggered via `trigger:force` (or via plain `trigger`
on macOS 14/15), the first instance log shows:

```
[first] OnSecondInstanceLaunch fired
[first]   Args           = [.../single-instance-url-scheme wails-single-url://hello?n=1]
[first]   url-in-args?   = true  (url="wails-single-url://hello?n=1")
```

`ApplicationLaunchedWithUrl` does **not** fire on either instance for the
second-instance relay path.

Timing (measured on macOS 26, Apple Silicon):
- URL captured: second instance exits in **~120–160 ms** (early-exit on event arrival).
- No URL (e.g. `open -n app.app`): second instance exits in **~320–430 ms** (300 ms timeout).

## Unfixed behaviour (pre-PR, for historical reference)

Triggered via `trigger:force` on the unfixed branch, the first instance logged:

```
[first] OnSecondInstanceLaunch fired
[first]   Args           = [.../single-instance-url-scheme]
[first]   url-in-args?   = false  (url="")
```

`ApplicationLaunchedWithUrl` did **not** fire on either instance.

## Notes / gotchas for testing

- Running the raw binary (`go run .`) will **not** exercise the bug path.
  macOS only routes custom URL schemes to apps launched from a `.app`
  bundle registered with LaunchServices. Use `wails3 task run`.
- If `open 'wails-single-url://…'` launches a different app, the scheme
  is claimed by a previously-registered bundle. Re-register this one:
  `lsregister -f bin/single-instance-url-scheme.dev.app`.
- Fresh launches (first instance) already go through the
  `NSAppleEventManager` path and work correctly; the fix is specific to
  the second-instance relay.
