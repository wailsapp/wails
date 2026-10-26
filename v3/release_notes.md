# Wails v3.0.0

The first stable release of Wails v3: build desktop and mobile apps with Go and the web frontend of your choice, with multiple windows, native menus and system tray, generated TypeScript bindings, an updater, packaging and signing, and a new documentation site.

## Install

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Wails v3 requires **Go 1.26 or newer**. Go 1.26.10 and 1.27.3 (or later) include a `crypto/x509` security fix, so we recommend building with one of those. Run `wails3 doctor` to check your setup.

## Security

All published Wails v3 security advisories are fixed in this release. If you are on `v3.0.0-beta.30` or earlier, please upgrade.

- Block reserved `wails:` messages from cross-origin child frames on Linux (GTK3 and GTK4) by injecting a guard script into all frames, mitigating GHSA-2p5f-2gx8-mh8m
- Update `github.com/klauspost/compress` to 1.18.7 (GHSA-259r-337f-4rfw) and `golang.org/x/image` to 0.45.0, and build with Go 1.26.10 to pick up the `crypto/x509` fix (GO-2026-6809)
- Update `vite` to 8.2.4 or later in the runtime package, setup wizard, examples and project templates (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Known limitations

- **Linux:** the child-frame guard cannot stop a *same-origin* child frame (for example a `srcdoc` or `about:blank` iframe) from reaching the native message handler through its parent, because WebKitGTK gives the handler no frame information. Cross-origin frames and sandboxed frames without `allow-same-origin` are blocked. Do not embed untrusted same-origin content in a Wails window.
- **Android (experimental):** the `window.wails` JavaScript bridge is exposed to every frame in the WebView, with no main-frame or origin check. Content in an embedded iframe, including cross-origin content, can call bound methods and send reserved `wails:` messages. Do not embed untrusted third-party content in the app's WebView. A fix is planned for 3.0.1.
- **All platforms:** bound-method calls from child frames travel over HTTP to `/wails/runtime`, not the native message handler, so the guard does not restrict them. Treat any iframe you embed as having access to your bound Go methods.

## Changes since v3.0.0-beta.34

- Fix Windows apps crashing with `panic: SetTextColor failed` when painting the application menu bar on some Windows builds (#6311)
- Fix Windows apps crashing in `applyDarkMode` when `uxtheme.dll` was unloaded while its functions were still in use (#6358)
- Keep a staged update safe while a restart is pending, and validate the update before replacing an installed application, so a damaged update rolls back instead of breaking a working app (#6350, #6360)
- Fix `wails3 update cli` not finding v3 releases
- Project templates pin `@wailsio/runtime` to the version of the CLI that generated them
- Remove the Crowdin integration after the final translation sync

## Releases from here on

After v3.0.0, **weekly** releases bump the minor version (`vX.Y.0`), dated **nightly** pre-releases (`vX.Y.0-nightly.YYYYMMDD`) preview the next weekly, and hotfixes bump the patch version. Go's `@latest` resolves to the newest weekly or hotfix release, never a nightly. Pin an exact version in production builds.

Full history: [changelog](https://github.com/wailsapp/wails/blob/master/docs/mpress/content/changelog.md).
