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

## Fixed
<!-- Bug fixes -->
- Fix Windows apps crashing with `panic: SetTextColor failed` when painting the application menu bar on some Windows builds: the menu bar draws on a valid device context or falls back to default drawing, `UAHMENUITEM.Position` now matches the C `int`, and menu painting never panics (#6311)
- Keep the staged update safe while a restart is pending: a download or a second restart started during shutdown cannot delete or contend for the payload the swap helper is waiting to install by @julianstorer
- Fix Windows apps crashing in `applyDarkMode` when `pkg/w32`'s init unloaded `uxtheme.dll` while its ordinal function pointers were still in use: the module is now kept loaded for the process lifetime (#6356)
- Validate the staged payload before the swap and the installed application before the backup is discarded, so a damaged update rolls back instead of replacing a working application by @julianstorer

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->
- Remove Crowdin integration after final translation sync in [PR](https://github.com/wailsapp/wails/pull/5736) by @github-actions[bot]

## Security
<!-- Security-related changes -->
- Block reserved `wails:` messages from cross-origin child frames on Linux (GTK3 and GTK4) by injecting a guard script into all frames, mitigating GHSA-2p5f-2gx8-mh8m
- Update `github.com/klauspost/compress` to 1.18.7 (GHSA-259r-337f-4rfw) and `golang.org/x/image` to 0.45.0, and build with Go 1.26.10 to pick up the `crypto/x509` fix (GO-2026-6809)
- Update `vite` to 8.2.4 or later in the runtime package, setup wizard, examples and project templates (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Known limitations

- **Linux:** the child-frame guard cannot stop a *same-origin* child frame (for example a `srcdoc` or `about:blank` iframe) from reaching the native message handler through its parent, because WebKitGTK gives the handler no frame information. Cross-origin frames and sandboxed frames without `allow-same-origin` are blocked. Do not embed untrusted same-origin content in a Wails window.
- **Android (experimental):** the `window.wails` JavaScript bridge is exposed to every frame in the WebView, with no main-frame or origin check. Content in an embedded iframe, including cross-origin content, can call bound methods and send reserved `wails:` messages. Do not embed untrusted third-party content in the app's WebView. This differs from macOS and iOS, which ignore `wails:` messages from child frames. A fix is planned for 3.0.1.
- **All platforms:** bound-method calls from child frames travel over HTTP to `/wails/runtime`, not the native message handler, so the guard does not restrict them. Treat any iframe you embed as having access to your bound Go methods.

---
