# Wails Enhancement Proposal (WEP)

## Expose the WebKit inactive scheduling policy on macOS

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Andrey Pshenkin  
**Created**: 2026-11-09  
**Discussion**: #5920 proposes an overlapping option; see Alternatives Considered  
**Implementor**: Andrey Pshenkin  
**Target**: Wails v3

## Summary

Add `Mac.WebviewPreferences.InactiveSchedulingPolicy`, a per-window, opt-in
option that maps one to one to `WKPreferences.inactiveSchedulingPolicy`
(macOS 14+). It is applied before the `WKWebView` is created, on both macOS
backends. The zero value leaves the configuration untouched, so default
behaviour does not change.

## Motivation

A window created with `Hidden: true` and shown later by its own page is a common
way to avoid a blank window during startup. On macOS 14+ that pattern can hang:

- Under WebKit's default policy (`Suspend`), the WebContent process of a window
  that has never been shown is suspended. In a desktop app using this pattern on
  macOS 27, the process was fully suspended about 7.5 seconds after launch,
  before a 10 second in-page `setTimeout` could fire. The page's plan to reveal
  or redirect itself never ran. `log stream` with
  `subsystem == "com.apple.WebKit" AND category == "ProcessSuspension"` shows the
  sequence.
- The failure depends on timing and on what the page waits for. Immediate work
  finishes during navigation, before the suspension. `WebviewWindow.ExecJS` takes
  a foreground activity and resumes the process, so flows driven from Go keep
  working. Only in-page timers and promises stall, which makes the cause hard to
  find.
- The behaviour depends on the SDK the app is linked against. In the WebKit
  source, `linkedOnOrAfterSDKWithBehavior` gates RunningBoard throttling on the
  macOS 14 SDK and immediate suspension of background content on the macOS 15
  SDK. An app that worked can start hanging after an Xcode update, with no code
  change.
- An application cannot fix this itself. WebKit latches the policy when the page
  is created, and Wails creates the `WKWebView` internally. Setting the policy
  from a native hook after creation protected only processes created later by a
  cross-site process swap, never the first page. The preference is not
  overridable with `defaults write`.

## Detailed Design

A new type and a new field on `MacWebviewPreferences`:

```go
// Values map to WKInactiveSchedulingPolicy, offset by 1 so the zero value means "unset".
type MacInactiveSchedulingPolicy int

const (
    MacInactiveSchedulingPolicyDefault MacInactiveSchedulingPolicy = iota
    MacInactiveSchedulingPolicySuspend
    MacInactiveSchedulingPolicyThrottle
    MacInactiveSchedulingPolicyNone
)

type MacWebviewPreferences struct {
    // ...
    InactiveSchedulingPolicy MacInactiveSchedulingPolicy
}
```

The offset follows the existing `MacWindowTabbingMode`. `Default` does not touch
the configuration. The other three values select the WebKit policy of the same
name, so the option name says what is set.

The value is assigned on the `WKWebViewConfiguration`'s preferences before
`initWithFrame:configuration:`, and is fixed for the lifetime of the webview.

- **cgo backend**: in `windowNew`, next to the other configuration-time
  preferences, compiled against a macOS 14+ SDK and guarded by
  `@available(macOS 14.0, *)`.
- **CGO-free backend**: in the window setup, before the `WKWebView` is
  allocated, guarded by `respondsToSelector:`. It cannot go in
  `applyWebviewPreferences`, which runs after the webview exists and would be too
  late for this preference.

No existing option or default changes.

## Non-Goals

- Does not lift WebKit's throttling of hidden pages. A hidden page still has
  `visibilityState: hidden`, aligned timers and no `requestAnimationFrame`.
- No runtime setter, because WebKit latches the value at page creation.
- No Windows or Linux behaviour. WebView2 has no equivalent policy, and its
  controller visibility trade-offs are a separate decision.
- Does not change the default for any window.

## Platform Considerations

- **macOS 14+**: the option applies on both backends.
- **macOS older than 14**, or a build against an older SDK: the assignment is
  skipped or compiled out, and the option has no effect.
- **Windows and Linux**: the field is part of the macOS options and has no
  effect there. WebView2 does not suspend a hidden webview by itself unless the
  app calls `TrySuspend()`, which Wails does not; it applies background
  throttling instead.
- **iOS**: out of scope for this proposal.

## Pros/Cons

### Pros

- Fixes the hidden-start hang with a single opt-in preference.
- A thin mapping of a first-party API, with every WebKit value reachable by name.
- Covers both macOS backends.
- Independent of any Windows design, so it can ship on its own.

### Cons

- Overlaps with #5920, so maintainers have to choose or combine.
- macOS only.
- Measurements in #5920 could not tell `Throttle` and `None` apart on macOS 27,
  so one of the named values may be redundant in practice today.
- Keeping a hidden webview running has a CPU and battery cost, which is the
  caller's decision.

## Alternatives Considered

- **`KeepRunningWhenHidden` (#5920)**. It uses the same WebKit preference on
  macOS and adds a Windows option. Its measurements agree with what we saw: the
  policy matters for a window that has never been shown. The differences are:
  - In #5920 `false` selects `Throttle`, which is milder than leaving the option
    unset. A named value avoids a boolean whose `false` is not the default.
  - `Suspend` can be pinned explicitly here, independent of future WebKit
    defaults.
  - The macOS change does not depend on the Windows dead-zone trade-off.
  - #5920 describes the cgo path. The CGO-free backend needs the placement
    described above.

  If #5920 is accepted instead, the SDK gating and placement findings here still
  apply to its macOS half.
- **A boolean selecting `None`**. Smaller surface, but it hides two of the three
  WebKit policies.
- **Setting the policy after the webview is created**, for example in
  `applyWebviewPreferences` or a native hook. Rejected: too late, see Motivation.
- **`Hidden: false` as a workaround**. It avoids the suspension but brings back
  the blank window at startup that `Hidden: true` exists to prevent.

## Backwards Compatibility

Behaviourally compatible: the zero value reproduces today's behaviour. At source
level this adds an exported field to `MacWebviewPreferences`, so an unkeyed
struct literal of that type would stop compiling. Keyed literals, the documented
form, are unaffected.

## Security and Privacy

No new capability is exposed to page content and no new data crosses the bridge.
The only effect is on scheduling. There are no permission or data-handling
changes.

## Test Plan

- In a `Hidden: true` window that is never shown, on macOS 14+ and an app linked
  against a current SDK, run a 10 second in-page timer and watch the
  `ProcessSuspension` log. `Default` and `Suspend` suspend the page. `Throttle`
  and `None` let the timer fire.
- Repeat on the CGO-free backend (`purego` tag).
- On a window that has been shown and then hidden, confirm that the option does
  not change behaviour beyond what the measurements in #5920 report.
- On macOS older than 14, setting the option does not error and changes nothing.
- Unit test that `getWebviewPreferences()` sets the native value only when the
  option differs from `Default`, with the offset applied.

## Reference Implementation

Branch [`darwin-inactive-policy`](https://github.com/APshenkin/wails/tree/darwin-inactive-policy)
on APshenkin/wails, one commit on top of `master`. The cgo path has been in use
in a desktop app since October 2026. The CGO-free path has not been run yet and
will be measured before an implementation PR is opened.

## Maintenance Plan

The change is a thin wrapper over a stable first-party preference, and tracks the
macOS SDK. It is maintained with the rest of `MacWebviewPreferences`.

## Conclusion

A small, opt-in, macOS-only preference removes a hang that applications cannot
fix on their own, and leaves the default untouched for every existing app.
