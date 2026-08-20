# Wails Enhancement Proposal (WEP)

## Keep WebViews running while hidden

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Julian Storer  
**Created**: 2026-08-04  
**Discussion**: [optional link to any prior discussion of the idea]  
**Implementor**: Julian Storer  
**Target**: Wails v3

## Summary

Add an opt-in per-window preference that stops the platform WebView from
throttling JavaScript while the window is hidden. It is exposed as
`Mac.WebviewPreferences.KeepRunningWhenHidden` on macOS and
`Windows.KeepRunningWhenHidden` on Windows. Default behaviour is unchanged.

## Motivation

Both platform WebViews aggressively throttle a hidden window's JavaScript:

- **macOS**: once another window in the application becomes active, a hidden
  `WKWebView` has its scheduling reduced all the way to 0Hz.
- **Windows**: WebView2 moves a hidden window's controller into "efficiency
  mode" (WebView2 issue #2861), driving the page's timers toward 0Hz.

That is the right default for a normal window the user has minimised or
switched away from. It is wrong for an application that deliberately keeps a
WebView hidden and uses it as a worker — for example a headless WebView that
runs plugin or engine code while a *different* window is the visible UI, or an
offscreen host used for tests. Today such a WebView stalls the moment any other
window takes focus, and there is no supported way to keep it running.

There is no reliable pure-web workaround: the throttling happens below the page,
so timers, `postMessage` loops and workers hosted in that WebView all slow down
regardless of what the page does.

## Detailed Design

A per-window, opt-in preference, off by default.

### macOS

`WKPreferences.inactiveSchedulingPolicy` (macOS 14+ / iOS 17+) controls the
throttling. A new tristate field is added to `MacWebviewPreferences`:

```go
// KeepRunningWhenHidden controls WKPreferences.inactiveSchedulingPolicy.
//   true  -> WKInactiveSchedulingPolicyNone     (never throttle)
//   false -> WKInactiveSchedulingPolicyThrottle (always throttle)
//   unset -> platform default (unchanged)
KeepRunningWhenHidden optional.Bool
```

`inactiveSchedulingPolicy` is read from the `WKWebViewConfiguration`'s
preferences at navigation time, so it is assigned alongside the other
configuration-time preferences (`tabFocusesLinks`, `textInteractionEnabled`,
`elementFullscreenEnabled`) before the `WKWebView` is allocated. Left unset, the
configuration is not touched and the OS default (automatic) applies.

### Windows

WebView2 exposes visibility through the controller's `IsVisible`. A new field is
added to `WindowsWindow`:

```go
// KeepRunningWhenHidden keeps the WebView2 controller IsVisible=true even when
// the window is hidden, so WebView2 never enters efficiency mode.
KeepRunningWhenHidden bool
```

When set, a hidden-and-never-shown window keeps its controller `IsVisible=true`
in `navigationCompleted` instead of hiding it. This fills in the
`PutIsVisible(true)` path the code already anticipated with a TODO referencing
issue #2861, now that `Chromium.Show()`/`Hide()` exist.

### Usage

```go
app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
    Hidden: true,
    Mac: application.MacWindow{
        WebviewPreferences: application.MacWebviewPreferences{
            KeepRunningWhenHidden: optional.NewBool(true),
        },
    },
    Windows: application.WindowsWindow{
        KeepRunningWhenHidden: true,
    },
})
```

## Non-Goals

- Not a global/application-wide setting; it is per-window and opt-in.
- Does not change the default. A window that does not set it behaves exactly as
  before, including full throttling when hidden.
- Not a power-management feature: keeping a hidden WebView running has a battery
  cost, which is the caller's decision to make.
- No Linux implementation (see Platform Considerations).

## Platform Considerations

- **macOS**: requires macOS 14+ for `inactiveSchedulingPolicy`. On older
  systems the property is unavailable and the option is a no-op; the WebView
  keeps the OS default behaviour. Tristate so callers can request "never
  throttle", "always throttle", or "leave the platform default alone".
- **Windows**: WebView2. Binary — the controller is either kept visible or not.
- **Linux**: not implemented in this proposal. WebKitGTK does not throttle a
  hidden WebView in the same way, so there is no equivalent knob to expose yet;
  the field is simply absent from the Linux options.

The macOS field is an `optional.Bool` (three states) and the Windows field is a
plain `bool` (two states), because the underlying APIs differ: macOS has a
genuine third "platform default" state worth preserving, Windows does not. This
asymmetry is deliberate but is the main open question for review — an
alternative is to make both `optional.Bool` for surface consistency (see
Alternatives).

## Pros/Cons

### Pros

- Enables a supported "background worker WebView" pattern that is currently
  impossible without stalling.
- Opt-in and per-window; zero effect on existing apps.
- Thin wrappers over first-party platform APIs, so low implementation and
  maintenance cost.

### Cons

- Two-platform-only; Linux has no equivalent, so behaviour is not uniform.
- macOS/Windows field types differ (tristate vs bool).
- Keeping hidden WebViews running costs CPU/battery if misused.

## Alternatives Considered

- **Make both fields `optional.Bool`** for a consistent surface. Rejected for
  the reference implementation because the Windows path has no meaningful third
  state, but it is a reasonable outcome of review.
- **A single cross-platform option on `WebviewWindowOptions`** instead of one
  field per platform block. Rejected: the semantics and availability differ per
  platform (macOS version gate, no Linux support), and Wails already groups
  platform-specific tuning under `Mac`/`Windows`/`Linux` blocks.
- **Do nothing / document the limitation.** Rejected: there is no page-level
  workaround, so an app that needs a background WebView has no path today.

## Backwards Compatibility

Fully backwards compatible. Both fields default to "unset"/false, which
reproduces today's behaviour exactly. No existing option is changed or removed.

## Security and Privacy

No new capability is exposed to page content and no new data crosses the
bridge. The only effect is scheduling: an opted-in hidden WebView keeps
executing the code it was already running. There are no permission or
data-handling changes.

## Test Plan

- Manual: a hidden WebView running a periodic timer (e.g. `postMessage` or a
  visible tick counter in a second window). With the option off, the tick rate
  collapses once another window is focused; with it on, the rate holds.
- macOS: verify the option is a no-op on < macOS 14 and does not error.
- Windows: verify a hidden-never-shown window keeps ticking and that a normally
  shown/hidden window is unaffected.
- An example under `v3/examples` demonstrating a background worker WebView can
  be added as part of implementation.

## Reference Implementation

A working reference implementation exists for both platforms (macOS via
`inactiveSchedulingPolicy`, Windows via the controller `IsVisible` path). The
implementation PR will be linked here once opened.

## Maintenance Plan

The change is a thin wrapper over `WKPreferences.inactiveSchedulingPolicy` and
the WebView2 controller `IsVisible`, both stable first-party APIs, so ongoing
maintenance is minimal and tracks the platform SDKs. Maintained alongside the
rest of the WebView option surface.

## Conclusion

`KeepRunningWhenHidden` makes the "hidden WebView as a background worker"
pattern viable on macOS and Windows with a small, opt-in, backwards-compatible
addition, while leaving the throttling default untouched for every existing app.
