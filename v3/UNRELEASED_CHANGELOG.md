# Unreleased Changes

<!--
AUTO-GENERATED FILE: do not edit in ordinary contribution PRs.

The Auto Changelog (v3) workflow generates entries after PRs merge.
Describe the user-visible change in the PR title and description instead.
Release tooling moves staged entries into docs/mpress/content/changelog.md
and resets this file. Do not delete or reset it by hand.

For changes to this automation or explicitly requested release-note corrections,
follow the changelog ownership policy in AGENTS.md and preserve existing entries.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->
- Link the 404 page to matching Wails v2 documentation pages on v2.wails.io in [PR](https://github.com/wailsapp/wails/pull/6366) by @taliesin-ai
- Move the v2 documentation site to v2.wails.io in [PR](https://github.com/wailsapp/wails/pull/6365) by @taliesin-ai
- Serve v3 documentation from wails.io and v2 documentation from v2.wails.io in [PR](https://github.com/wailsapp/wails/pull/6364) by @taliesin-ai
- Update v3.0.0 release notes and correct the security advisory claim in [PR](https://github.com/wailsapp/wails/pull/6363) by @taliesin-ai

## Fixed
<!-- Bug fixes -->
- Linux: detect when WebKitGTK's bwrap sandbox cannot start because AppArmor/the kernel restricts unprivileged user namespaces (Ubuntu 24.04), warn at startup with the AppArmor/sysctl fix instead of an opaque SIGTRAP, and add a check to `wails3 doctor`. The sandbox is never disabled automatically (#6172)
- Fix Windows webview not regaining focus on window activation, and Tab being trapped inside the webview instead of moving to the native window and back (#3783)
- Fix `internal/webview2/pkg/webview2` COM vtable structs (`ICoreWebView2_2`..`_27` and other derived interfaces) omitting the inherited method chain, which made calls land on the wrong vtable slot (#6129)
- Fix Windows question dialogs ignoring custom button labels and never firing button callbacks; custom buttons are now shown with a TaskDialog (falling back to MessageBox with positional callback mapping when comctl32 v6 is unavailable) (#5794)
- Fix `wails3 dev` hanging for 60s against Vite: the automatic frontend readiness gate now uses a TCP probe instead of an HTTP GET that sent no `Accept` header and received a 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Keep systray menus working when updated while open in [PR](https://github.com/wailsapp/wails/pull/6301) by @taliesin-ai
- Handle transient WebView2 COM errors gracefully in runtime operations in [PR](https://github.com/wailsapp/wails/pull/6302) by @leaanthony
- Preserve request bodies on Windows when WebView2 omits the Content-Length header in [PR](https://github.com/wailsapp/wails/pull/6326) by @strraaa
- Fix data race in event hook emission by snapshotting hooks under lock in [PR](https://github.com/wailsapp/wails/pull/6276) by @DevLumuz
- Fix `SystemTray.SetTooltip()` doing nothing on Linux: the StatusNotifierItem tooltip description is now set (#6017)
- Run changelog evidence check in the weekly release workflow in [PR](https://github.com/wailsapp/wails/pull/6374) by @taliesin-ai
- Fix release changelog evidence check to map runs to squash-merge commits in [PR](https://github.com/wailsapp/wails/pull/6373) by @taliesin-ai
- Fix Linux dark mode detection returning light on a dark desktop until the colour scheme changed, and `linux:SystemThemeChanged` never firing on the GTK4 backend. Both backends now share one implementation, backed by the `org.freedesktop.appearance` portal namespace (#6072)
- `Window.SetZoom` no longer raises factors below 1.0 to 1.0 on framed windows (Windows, macOS, Linux): WebView2 `ZoomFactor`, WKWebView `magnification` and WebKitGTK's zoom level accept any positive value, and the `Zoom` window option already passed such values through. Frameless windows keep the 1.0 floor because their drag and resize regions are not adjusted for zoom ([#4590](https://github.com/wailsapp/wails/issues/4590)); non-positive, NaN and Inf factors are now ignored instead of being raised to 1.0 — by @APshenkin

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---
