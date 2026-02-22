# iOS Features TODO (Prioritized)

This document lists potential iOS features and platform options to enhance the Wails v3 iOS runtime. Items are ordered by importance for typical app development workflows.

## Top Priority (Implement first)

1) Input accessory bar control
- Status: Implemented as `IOSOptions.DisableInputAccessoryView` (default false = shown). Native toggle + WKWebView subclass.

2) Scrolling and bounce behavior
- Options:
  - `DisableScroll` (default true in current runtime to preserve no-scroll template behavior)
  - `DisableBounce` (default true in current runtime)
  - `HideScrollIndicators` (default true in current runtime)
- Purpose: control elastic bounce, page scrolling, and indicators.

3) Web Inspector / Debug
- Options:
  - `DisableInspectable` (default false; inspector enabled by default in dev)
- Purpose: enable/disable WKWebView inspector.

4) Back/forward navigation gestures
- Options:
  - `AllowsBackForwardNavigationGestures` (default false)
- Purpose: enable iOS edge-swipe navigation.

5) Link previews
- Options:
  - `DisableLinkPreview` (default false)
- Purpose: allow long-press link previews.

6) Media autoplay and inline playback
- Options:
  - `DisableInlineMediaPlayback` (default false)
  - `RequireUserActionForMediaPlayback` (default false)
- Purpose: control media playback UX.

7) User agent customization
- Options:
  - `UserAgent` (string)
  - `ApplicationNameForUserAgent` (string; default "wails.io")
- Purpose: customize UA / identify app.

8) Keyboard behavior
- Options:
  - Already: `DisableInputAccessoryView`
  - Future: `KeyboardDismissMode` (none | onDrag | interactive)
- Purpose: refine keyboard UX.

9) Safe-area and content inset behavior
- Options (future):
  - `ContentInsetAdjustment` (automatic | never | always)
  - `UseSafeArea` (bool)
- Purpose: fine-tune layout under notch/home indicator.

10) Data detectors (future feasibility)
- Options: `DataDetectorTypes []string` (phoneNumber, link, address)
- Note: Not all are directly available on WKWebView; feasibility TBD.

## Medium Priority

11) Pull-to-refresh (custom)
12) File picker / photo access bridges
13) Haptics feedback helpers
14) Clipboard read/write helpers (partially present)
15) Share sheet / activity view bridges
16) Background audio / PiP controls
17) App lifecycle event hooks (background/foreground)
18) Permissions prompts helpers (camera, mic, photos)
19) Open in external browser vs in-app policy
20) Cookie / storage policy helpers

## Low Priority

21) Theme/dynamic color helpers bridging to CSS vars
22) Orientation lock helpers per window
23) Status bar style control from Go
24) Network reachability events bridge
25) Push notifications

---

# Implementation Plan (Top 10)

Implement the following immediately:
- DisableScroll, DisableBounce, HideScrollIndicators
- AllowsBackForwardNavigationGestures
- DisableLinkPreview
- DisableInlineMediaPlayback
- RequireUserActionForMediaPlayback
- DisableInspectable
- UserAgent
- ApplicationNameForUserAgent

Approach:
- Extend `IOSOptions` in `pkg/application/application_options.go` with these fields.
- Add native globals + C setters in `pkg/application/application_ios.h/.m`.
- Apply options in `pkg/application/webview_window_ios.m` during WKWebView configuration and on the scrollView.
- Wire from Go in `pkg/application/application_ios.go`.
- Maintain current template behavior as defaults (no scroll/bounce/indicators) to avoid regressions in existing tests.
