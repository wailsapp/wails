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
- `Window.SetZoom` no longer raises factors below 1.0 to 1.0 on framed windows (Windows, macOS, Linux): WebView2 `ZoomFactor`, WKWebView `magnification` and WebKitGTK's zoom level accept any positive value, and the `Zoom` window option already passed such values through. Frameless windows keep the 1.0 floor because their drag and resize regions are not adjusted for zoom ([#4590](https://github.com/wailsapp/wails/issues/4590)); non-positive, NaN and Inf factors are now ignored instead of being raised to 1.0 — by @APshenkin
- Run changelog evidence check in the weekly release workflow in [PR](https://github.com/wailsapp/wails/pull/6374) by @taliesin-ai
- Fix release changelog evidence check to map runs to squash-merge commits in [PR](https://github.com/wailsapp/wails/pull/6373) by @taliesin-ai
- Fix Linux dark mode detection returning light on a dark desktop until the colour
  scheme changed, and `linux:SystemThemeChanged` never firing on the GTK4 backend.
  Both backends now share one implementation, backed by the
  `org.freedesktop.appearance` portal namespace (#6072)

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---
