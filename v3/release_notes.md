# Wails v3.0.0

**Wails v3 is here.** Build desktop and mobile applications with Go and the web frontend you already know, using each platform's native WebView instead of a bundled browser.

v3 is a ground-up rework of Wails. Applications, windows, services, events and platform features are now explicit objects you create and talk to directly, so multiple windows, system trays and background services are a normal part of the model rather than workarounds. It has been nearly four years since the first v3 alpha in January 2023. Thank you to everyone who tested, reported, translated, reviewed and contributed along the way.

📣 **Read the announcement:** [Wails v3.0.0 Released](https://wails.io/blog/wails-v3-released/)
📚 **Documentation:** [wails.io](https://wails.io/) · **Migrating from v2:** [v2 to v3 guide](https://wails.io/migration/v2-to-v3/)

## Install

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
wails3 setup    # guided setup: checks your system and installs what Wails needs
wails3 init     # create a new project
```

Wails v3 requires **Go 1.26 or newer**. We recommend Go 1.26.10 or 1.27.3 (or later), which include a `crypto/x509` security fix. Run `wails3 doctor` at any time to check your setup.

## What's new in v3

### Application model
- **Multiple windows.** Create, show, hide and close any number of windows at runtime, each with its own options, events and lifecycle.
- **Services.** Your backend is ordinary Go: register services with the application and their exported methods become callable from the frontend. Services can have startup and shutdown hooks and can serve HTTP routes.
- **Manager APIs.** Windows, menus, dialogs, events, screens, the clipboard and more are reached through clear manager objects on the application, instead of a single context-bound runtime package.
- **Events.** Application, window and custom events, with typed payloads, available from both Go and JavaScript, plus platform system events.

### Bindings
- **Generated with static analysis.** `wails3 generate bindings` reads your Go source, so doc comments and real parameter names carry through to the generated TypeScript or JavaScript.
- **Models and enums.** Go structs become TypeScript classes or interfaces, and Go enums become TypeScript enums.
- **Streams and raw messages.** Stream data between Go and the frontend, or use raw messages and custom transports when you need full control.

### Native user interface
- **System tray.** Tray icons with menus, light and dark icons, macOS template icons, and the option to attach a window to the tray icon.
- **Menus.** Application menus, context menus (including from HTML with WML attributes) and tray menus, all updatable at runtime.
- **Dialogs.** Native message, question, file open and save dialogs, plus custom dialogs.
- **Frameless and custom windows.** Frameless windows with native corners and shadows, custom title bars, drag regions and translucency, including **native Liquid Glass on macOS 26** with a fallback on earlier versions.
- **Notifications.** Native desktop notifications with actions.
- **Keyboard shortcuts.** Window and application key bindings and **global shortcuts**; on macOS, letter accelerators follow the labelled key on AZERTY, QWERTZ, Dvorak and other layouts.
- **And more:**
  - clipboard;
  - file drag and drop, plus HTML drag and drop;
  - screen and environment information;
  - dock badges on macOS;
  - start at login;
  - opening URLs and files in the default browser;
  - **notch windows** on macOS;
  - a per-window **permissions** map for the camera, microphone, geolocation and notifications.

### Platforms
- **Windows** (amd64 and arm64), **macOS** (Intel and Apple Silicon) and **Linux** (amd64 and arm64).
- **Linux uses GTK4 and WebKitGTK 6.0 by default.** The legacy GTK3 stack is still available with `-tags gtk3` until v3.1.
- **iOS and Android (experimental).** The same `main.go` and frontend can run on mobile. Try it and tell us what you find; it isn't yet covered by the stability promise.
- **Server builds.** Run the same application and services as an HTTP server, with no native window.

### Build, ship and update
- **A build system you can see.** Projects use a [Taskfile](https://taskfile.dev)-based build that you can read, extend and debug, instead of one opaque command.
- **Cross-compilation**, with platform guides for building on and for every desktop target.
- **Code signing and installers.** Signing for every platform, NSIS installers and MSIX packages for Windows, DMG and app bundles for macOS, and AppImage, deb, rpm and Arch packages for Linux.
- **An in-app updater.** Pluggable providers (including GitHub Releases), cryptographic verification, an atomic swap with rollback, and a default UI you can theme or replace.
- **Single instance, file associations and custom URL protocols**, so your app opens the files and links it owns.
- **Obfuscated builds** and access to private macOS APIs when you need them.

### Developer experience
- **`wails3 setup`**, a guided setup wizard that checks and installs dependencies.
- **`wails3 dev`**, live development with frontend hot reload and Go rebuilds, which stops its background processes when your app exits.
- **Templates** for vanilla, React, Vue and Svelte, each in JavaScript or TypeScript, plus your own custom templates.
- **Testing support**, including end-to-end testing guides.
- **An MCP service** that lets AI tools drive and inspect a running Wails app.
- **A new documentation site** at [wails.io](https://wails.io/), built with [M-Press](https://m-press.me), with tutorials, guides and API reference translated into nine languages.

## Moving from v2

v3 is a new major version, so moving to it is a real port rather than a version bump. The main changes are:
- the application and window lifecycle;
- services in place of context-bound bindings;
- direct application and window APIs in place of the v2 runtime package;
- regenerated frontend bindings.

The [migration guide](https://wails.io/migration/v2-to-v3/) covers each one with before-and-after examples.

**Wails v2 isn't going away.** It continues to receive critical updates, and its documentation is now at [v2.wails.io](https://v2.wails.io/). The two CLIs (`wails` and `wails3`) are separate binaries with different module paths, so you can install both.

## Security

This release fixes every published Wails v3 security advisory except one: **GHSA-2p5f-2gx8-mh8m is only partly fixed.** Cross-origin child frames on Linux are now blocked, and the cases that remain open are listed under *Known limitations* below. If you are on `v3.0.0-beta.30` or earlier, please upgrade.

- Block reserved `wails:` messages from cross-origin child frames on Linux (GTK3 and GTK4) by injecting a guard script into all frames, mitigating GHSA-2p5f-2gx8-mh8m
- Update `github.com/klauspost/compress` to 1.18.7 (GHSA-259r-337f-4rfw) and `golang.org/x/image` to 0.45.0, and build with Go 1.26.10 to pick up the `crypto/x509` fix (GO-2026-6809)
- Update `vite` to 8.2.4 or later in the runtime package, setup wizard, examples and project templates (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Known limitations
- **Linux:** the child-frame guard cannot stop a *same-origin* child frame (for example a `srcdoc` or `about:blank` iframe) from reaching the native message handler through its parent, because WebKitGTK gives the handler no frame information. Cross-origin frames, and sandboxed frames without `allow-same-origin`, are blocked. Do not embed untrusted same-origin content in a Wails window.
- **Android (experimental):** the `window.wails` JavaScript bridge is exposed to every frame in the WebView. Do not embed untrusted third-party content in the app's WebView. A fix is planned for 3.0.1.
- **All platforms:** bound-method calls from child frames go over HTTP to `/wails/runtime`, so the guard does not restrict them. Treat any iframe you embed as having access to your bound Go methods.

## Changes since v3.0.0-beta.34

- Fix Windows apps crashing with `panic: SetTextColor failed` when painting the application menu bar on some Windows builds (#6311)
- Fix Windows apps crashing in `applyDarkMode` when `uxtheme.dll` was unloaded while its functions were still in use (#6358)
- Keep a staged update safe while a restart is pending, and validate an update before replacing the installed application, so a damaged update rolls back instead of breaking a working app (#6350, #6360)
- Fix `wails3 update cli` not finding v3 releases
- Project templates pin `@wailsio/runtime` to the version of the CLI that generated them

## Releases from here on

- **Weekly releases** bump the minor version (`v3.1.0`, `v3.2.0`, ...).
- **Nightly pre-releases** preview the next weekly (`v3.1.0-nightly.YYYYMMDD`).
- **Hotfixes** bump the patch version.

Go's `@latest` always resolves to the newest weekly or hotfix release, never a nightly. The `@wailsio/runtime` npm package always has the same version as the Wails release. Pin an exact version in production builds.

To propose a new feature or a change to public behaviour, open a draft [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md) pull request rather than a feature-request issue.

## Thank you

Wails v3 exists because of the people who tested incomplete builds, reported hard bugs, translated documentation, answered questions and contributed code: see the [credits page](https://wails.io/credits/). Special thanks to [Stffabi](https://github.com/stffabi), [Travis McLane](https://github.com/tmclane), [Atterpac](https://github.com/atterpac), [Lyimmi](https://github.com/Lyimmi) and [fbbdev](https://github.com/fbbdev); to Valerio Tesei for his years of support; to our sponsors; and to JetBrains, ZSA and [CodeRabbit](https://www.coderabbit.ai/oss).

If Wails is useful to you or your company, please consider [sponsoring the project](https://github.com/sponsors/leaanthony). 💙

**Full changelog:** [changelog](https://github.com/wailsapp/wails/blob/master/docs/mpress/content/changelog.md)
