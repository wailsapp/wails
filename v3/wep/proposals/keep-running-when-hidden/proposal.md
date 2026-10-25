# Wails Enhancement Proposal (WEP)

## Keep WebViews running while hidden

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Julian Storer  
**Created**: 2026-08-04  
**Discussion**: None (this PR is the first discussion)  
**Implementor**: Julian Storer  
**Target**: Wails v3

## Summary

Add an opt-in per-window preference that keeps the WebView of a window that has
never been shown running, instead of suspended (macOS) or in efficiency mode
(Windows). It is exposed as `Mac.WebviewPreferences.KeepRunningWhenHidden` on
macOS and `Windows.KeepRunningWhenHidden` on Windows. Default behaviour is
unchanged.

## Motivation

On both platforms, the WebView of a window that is created hidden and never
shown effectively stops:

- **macOS**: under WebKit's default inactive scheduling policy (`Suspend`), a
  `WKWebView` in a window that has never been shown is suspended. Its timers do
  not fire at all (see Measured behaviour below).
- **Windows**: a WebView2 controller with `IsVisible=false` is put into
  "efficiency mode", driving the page's timers toward 0Hz (see
  wailsapp/wails#2861). Since #5743, a window that is created `Hidden` and never
  shown ends its first navigation with the controller hidden, so it is throttled
  from the start.

That is a reasonable default for a window that is not yet in use. It is wrong
for an application that deliberately keeps a WebView hidden and uses it as a
worker — for example a headless WebView that
runs plugin or engine code while a *different* window is the visible UI, or an
offscreen host used for tests. Today such a WebView never runs, and there is
no supported way to change that.

There is no pure-web workaround: the suspension happens below the page, so
timers, `postMessage` loops and workers hosted in that WebView all stop
regardless of what the page does.

This proposal is deliberately narrow. It keeps such a WebView *running*; it does
not exempt it from the platforms' ordinary throttling of hidden pages, and it
does not change windows that have been shown.

## Detailed Design

A per-window, opt-in preference, off by default.

### macOS

`WKPreferences.inactiveSchedulingPolicy` (macOS 14+ / iOS 17+) controls the
throttling. It has three values — `Suspend`, `Throttle` and `None` — and Apple
documents `WKInactiveSchedulingPolicySuspend` as the default. A new tristate
field is added to `MacWebviewPreferences`:

```go
// KeepRunningWhenHidden sets WKPreferences.inactiveSchedulingPolicy.
//   true  -> WKInactiveSchedulingPolicyNone
//   false -> WKInactiveSchedulingPolicyThrottle
//   unset -> not assigned; WebKit's default, WKInactiveSchedulingPolicySuspend
KeepRunningWhenHidden optional.Bool
```

Note that `false` is *milder* than the default: it selects `Throttle`, which
does not suspend the page. Both set values keep a never-shown WebView running;
see Measured behaviour for how they compare.

The value is assigned on the `WKWebViewConfiguration`'s preferences in
`windowNew`, alongside the other configuration-time preferences
(`tabFocusesLinks`, `textInteractionEnabled`, `elementFullscreenEnabled`), before
the `WKWebView` is allocated from that configuration. It is fixed for the
lifetime of the WebView; there is no runtime setter. The assignment is compiled
only against a macOS 14+ SDK (`MAC_OS_X_VERSION_MAX_ALLOWED >= 140000`) and runs
only under `@available(macOS 14.0, *)`. Left unset, the configuration is not
touched.

Apple's header describes the policy as applying to a web view that is "inactive
and detached from the view hierarchy". The measurements below show what that
means for a WebView in a hidden `NSWindow`.

#### Measured behaviour

Measured on macOS 27.0.1 (Apple silicon) with a standalone probe: a `WKWebView`
in one window runs a 100ms `setInterval`, in the page or in a dedicated worker,
while a second window of the same app is key. Counts are ticks in 10 seconds,
taken after 10 seconds to settle; 100 means unthrottled. Repeating the
measurement after three minutes gave the same counts. The probe is
[`measure.m`](measure.m), next to this proposal.

| Window | Policy | Page timer | Worker timer |
|---|---|---|---|
| Visible, key | any | 100 | 102 |
| Never shown | unset (`Suspend`) | **0** | **0** |
| Never shown | `Throttle` / `None` | 10 / 11 | 54 / 53 |
| Shown, then `orderOut:` | unset / `Throttle` / `None` | 10 / 10 / 10 | 53 / 54 / 54 |
| Miniaturised | unset / `None` | 11 / 10 | — / 53 |

What this shows:

- The policy matters for a window that has never been shown. There, the
  default suspends the page and either set value keeps it running.
- It does not lift WebKit's throttling of hidden pages. A hidden or
  miniaturised window's timers run at about 1Hz in the page and about half rate
  in a worker, under every policy. A busy loop in a worker got about 15% of its
  visible throughput under both `Throttle` and `None`, and
  `requestAnimationFrame` did not fire in any hidden state.
- A window that has been shown is not suspended once it is hidden, so the
  option changes nothing for it.
- `None` and `Throttle` behaved identically in every case measured (see the
  open question under Platform Considerations).

### Windows

WebView2 exposes visibility through the controller's `IsVisible`
(`Chromium.Show()`/`Hide()`). A new field is added to `WindowsWindow`:

```go
// KeepRunningWhenHidden leaves the WebView2 controller visible after the first
// navigation of a window that is not logically visible at that point (created
// Hidden and not yet shown, or hidden before that navigation completed), so
// WebView2 does not put it into efficiency mode. Showing and then hiding or
// minimising the window still hides the controller as usual.
KeepRunningWhenHidden bool
```

`navigationCompleted` always ends its first-paint nudge (`Hide()` then
`Show()`) with the controller visible. #5743 added a final `chromium.Hide()`
there for any window with `!windowShown && !showRequested`, so that a
never-shown window has no live input surface. When `KeepRunningWhenHidden` is
set, that final `Hide()` is skipped, and nothing else changes:

- `hide()` and the `SIZE_MINIMIZED` branch of `WM_SIZE` still call
  `chromium.Hide()` (#5742), so a window that is shown and then hidden or
  minimised is throttled as before. Keeping such windows running is out of
  scope.
- The `PutIsVisible(true)` TODO in the WebView2 setup is left alone.

**What this costs.** Under visual hosting, a WebView2 controller that stays
visible inside a hidden window leaves its DirectComposition input surface
hit-testing at the window's location. Visual hosting comes from either of two
independent options: the application-wide `Options.Windows.UseVisualHosting`,
which forces `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL` for every window, or
the per-window `WebviewWindowOptions.Windows.WebView2CompositionHosting`, which
creates that window's WebView2 through a composition controller. #5742 observed
the dead zone under `UseVisualHosting`; the Test Plan checks both. That is
an invisible desktop "dead zone" that swallows clicks, including right-clicks on
the desktop. Preventing it is exactly why #5742 and #5743 hide the controller.
An opted-in window accepts that dead zone in exchange for staying live, so the
option is meant only for windows that are never shown to the user. Positioning
such a window off-screen keeps the dead zone out of the way.

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
- On Windows, does not keep a window running after it has been shown and then
  hidden or minimised. It covers only windows that are hidden at their first
  navigation.
- No Linux implementation (see Platform Considerations).

## Platform Considerations

- **macOS**: requires macOS 14+ for `inactiveSchedulingPolicy`. On older
  systems the `@available` check skips the assignment, so the option is a no-op
  and the WebView keeps its existing behaviour. Building against an SDK older
  than 14 compiles the assignment out entirely. The field is a tristate because
  it maps Apple's three documented policies one to one: `None`, `Throttle`, or
  "leave the platform default (`Suspend`) alone". Apple documents `None` and
  `Throttle` as distinct, but they could not be told apart on macOS 27. If
  review prefers the smaller surface, a plain `bool` selecting `None` loses
  nothing that was measured.
- **Windows**: WebView2. Binary — the controller is either left visible after
  first navigation or not. Applies only to windows hidden at that point, and
  carries the DirectComposition dead-zone cost under visual hosting
  (`UseVisualHosting` or `WebView2CompositionHosting`; see Detailed Design).
- **Linux**: not implemented in this proposal, and the field is absent from the
  Linux options. WebKitGTK documents suspending the rendering of a hidden view,
  but not what happens to its JavaScript timers, on either supported backend
  (GTK4/WebKitGTK 6.0 or GTK3/WebKit2GTK 4.1). Whether Linux needs an
  equivalent is unestablished, and is left for a follow-up backed by a
  measurement like the one in the Test Plan.

The macOS field is an `optional.Bool` (three states) and the Windows field is a
plain `bool` (two states), because the underlying APIs differ: macOS exposes a
three-valued policy, Windows does not. This asymmetry is the main open question
for review. The alternatives are to make both `optional.Bool` for surface
consistency, or both `bool`, given the measurements above (see Alternatives).

## Pros/Cons

### Pros

- Enables a supported "background worker WebView" pattern that today does not
  run at all.
- Opt-in and per-window; zero effect on existing apps.
- Thin wrappers over first-party platform APIs, so low implementation and
  maintenance cost.

### Cons

- Two-platform-only; Linux has no equivalent, so behaviour is not uniform.
- macOS/Windows field types differ (tristate vs bool).
- The platforms cover different cases: macOS applies the policy for the
  WebView's whole lifetime, while Windows covers only windows hidden at first
  navigation.
- On Windows under visual hosting (`UseVisualHosting` or
  `WebView2CompositionHosting`), an opted-in hidden window leaves an invisible
  input dead zone at its location.
- It keeps the WebView running, not running at full rate. Hidden pages are
  still throttled, so work that needs its full rate does not belong in a
  hidden window.
- Keeping hidden WebViews running costs CPU/battery if misused.

## Alternatives Considered

- **Make both fields `optional.Bool`** for a consistent surface. Rejected for
  the reference implementation because the Windows path has no meaningful third
  state, but it is a reasonable outcome of review.
- **Make both fields `bool`**, with macOS `true` selecting `None`. This gives the
  smaller surface, and nothing measured is lost, but `Throttle` can no longer
  be selected. It is also a reasonable outcome of review.
- **A single cross-platform option on `WebviewWindowOptions`** instead of one
  field per platform block. Rejected: the semantics and availability differ per
  platform (macOS version gate, no Linux support), and Wails already groups
  platform-specific tuning under `Mac`/`Windows`/`Linux` blocks.
- **On Windows, also skip `chromium.Hide()` in `hide()` and on minimise**, so
  the option holds across every visibility transition. Rejected: it would
  reintroduce the #5742 dead zone for ordinary user-facing windows, and the
  pattern this targets uses windows that are never shown.
- **Do nothing / document the limitation.** Rejected: there is no page-level
  workaround, so an app that needs a background WebView has no path today.

## Backwards Compatibility

Behaviourally backwards compatible. Both fields default to "unset"/false, which
reproduces today's behaviour exactly, and no existing option is changed or
removed. At source level this is an additive API change: it adds exported fields
to `MacWebviewPreferences` and `WindowsWindow`, so code that builds either struct
with an unkeyed (positional) literal stops compiling. Keyed literals, the form
the documentation uses, are unaffected.

## Security and Privacy

No new capability is exposed to page content and no new data crosses the
bridge. The only effect is scheduling: an opted-in hidden WebView keeps
executing the code it was already running. There are no permission or
data-handling changes.

## Test Plan

The probe in every case is a hidden WebView running a periodic timer, with its
tick rate reported to a second, visible window that holds focus.

- **macOS option values**: in a `Hidden: true` window that is never shown,
  unset leaves the page suspended (no ticks), and `true` and `false` both keep
  it ticking. In a window that is shown and then hidden or miniaturised, all
  three values tick at the same throttled rate. These are the results in
  Measured behaviour.
- **macOS older than 14** (or a build against an SDK older than 14): setting the
  option compiles, does not error, and changes nothing.
- **Windows, never shown**: a `Hidden: true` window with the option keeps
  ticking after the first navigation; without it, it throttles as today.
- **Windows, shown → hidden and minimised → restored**: with or without the
  option, `hide()` and minimise hide the controller, and show/restore brings it
  back. This confirms the option does not leak into those paths.
- **Windows dead zone**: with `UseVisualHosting`, and separately with
  `WebView2CompositionHosting`, confirm the dead zone appears only for opted-in
  hidden windows and never for default ones.
- **Propagation**: unit tests check that the public option reaches the native
  side. On macOS, `getWebviewPreferences()` sets `KeepRunningWhenHidden` only
  when the option is set. On Windows, the `navigationCompleted` condition reads
  `options.Windows.KeepRunningWhenHidden`.
- An example under `v3/examples` demonstrating a background worker WebView can
  be added as part of implementation. It supplements the tests above; it does
  not replace them.

## Reference Implementation

A working reference implementation for both platforms is on the juggler-ai/wails
fork's branch [`feat/keep-running-when-hidden`](https://github.com/juggler-ai/wails/tree/feat/keep-running-when-hidden),
two commits on top of current `master`:

- macOS, `inactiveSchedulingPolicy`:
  [e5cb732c3ad5f3b3b964f3295d2acf0731de5ff3](https://github.com/juggler-ai/wails/commit/e5cb732c3ad5f3b3b964f3295d2acf0731de5ff3)
- Windows, skipping the never-shown `chromium.Hide()` in `navigationCompleted`:
  [dcf87ce03626a2d726b9ad06e7917293fb2590cc](https://github.com/juggler-ai/wails/commit/dcf87ce03626a2d726b9ad06e7917293fb2590cc)

The implementation PR will be linked here once opened.

## Maintenance Plan

The change is a thin wrapper over `WKPreferences.inactiveSchedulingPolicy` and
the WebView2 controller `IsVisible`, both stable first-party APIs, so ongoing
maintenance is minimal and tracks the platform SDKs. On Windows it is one extra
condition on the #5743 re-hide, so it needs revisiting whenever that visibility
logic changes. Maintained alongside the rest of the WebView option surface.

## Conclusion

`KeepRunningWhenHidden` makes the "hidden WebView as a background worker"
pattern viable on macOS and Windows with a small, opt-in, behaviourally
backwards-compatible addition, while leaving the default untouched for every
existing app. On both platforms it covers windows that are never shown, which
is the pattern it serves. It keeps them running but still throttled, and on
Windows it costs an input dead zone under visual hosting.
