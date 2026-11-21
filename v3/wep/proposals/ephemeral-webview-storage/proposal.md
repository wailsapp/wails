# Wails Enhancement Proposal (WEP)

## Ephemeral webview storage

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Julian Storer  
**Created**: 2026-10-25  
**Discussion**: [#2255](https://github.com/wailsapp/wails/issues/2255) (related symptom), [#6258](https://github.com/wailsapp/wails/pull/6258) (open, overlapping Windows options)  
**Implementor**: Julian Storer  
**Target**: Wails v3

## Summary

Add `Options.EphemeralWebviewStorage bool`, which stops every webview in the
application persisting website data (HTTP cache, cookies, local storage,
IndexedDB and the rest of the platform's website data store). On macOS and
Linux all windows share one in-memory store for the life of the process; on
Windows each WebView2 controller is created in InPrivate mode. The default,
`false`, is today's behaviour.

## Motivation

Webviews persist website data per origin, and nothing in Wails prunes it. An
app served from a local HTTP server on a random port gets a new origin each
run, leaving stores no later run can read — gigabytes over time, somewhere the
user does not know about.

In [#2255](https://github.com/wailsapp/wails/issues/2255) ("Wails spams AppData
folder") each run left a new ~100 MB `EBWebView` folder. The cause differed
(the folder is named after the binary, which `wails dev` renamed each reload;
fixed in #2103, though a 2024 comment notes it recurs when a built binary is
renamed) and is not fixed here, since WebView2 still creates the folder in
InPrivate mode — but it is the same unbounded growth.

Kiosk tools, viewers of untrusted or one-off content, and test harnesses also
want nothing to outlive the process, and can only delete directories they do
not own, whose layout varies by platform and WebView2 version.

## Detailed Design

### API

```go
type Options struct {
    // ...

    // EphemeralWebviewStorage stops webviews persisting website data — HTTP
    // cache, cookies, local storage, IndexedDB — so nothing a page stores
    // outlives the application.
    //
    // macOS and Linux use one non-persistent store shared by every window for
    // the life of the process. Windows creates each WebView2 controller in
    // InPrivate mode; WebView2 still creates its user data folder (see
    // Windows.WebviewUserDataPath), and window creation fails on a runtime
    // too old to support InPrivate. Ignored on iOS and Android.
    EphemeralWebviewStorage bool
}
```

Read when each window's webview is created; not a runtime toggle.

### Semantics

- **Scope:** one store shared by all windows for the process lifetime, as
  today's persistent store is; nothing is shared with any other run.
- **Covered:** the platform's website data store — HTTP cache, cookies,
  `localStorage`, `sessionStorage`, IndexedDB, Cache Storage, service
  workers. Not files the application writes, nor engine data outside that
  store.
- **Guarantee:** "not persisted by the webview"; data lives at most as long as
  the process, with no claim about when the engine frees it.

### Implementation

- **macOS:** `windowNew` sets `WKWebViewConfiguration.websiteDataStore` to one
  retained `[WKWebsiteDataStore nonPersistentDataStore]` (each call returns a
  new empty store); it runs on the main thread, so no lock.
- **WebKitGTK 6:** `"network-session"` is set to one
  `webkit_network_session_new_ephemeral()`, created on first use and never
  released; unset, the default session is used.
- **WebKitGTK 4.1:** `"web-context"` is set to one
  `webkit_web_context_new_ephemeral()`, likewise. Sharing it also keeps the
  `wails` URI scheme, registered once on the first context, working in every
  window.
- **Windows:** `edge.Chromium` gains `InPrivate bool`. When set, controllers are
  created via `ICoreWebView2Environment10.CreateCoreWebView2ControllerOptions`
  with `put_IsInPrivateModeEnabled(TRUE)`, then
  `CreateCoreWebView2ControllerWithOptions` or
  `CreateCoreWebView2CompositionControllerWithOptions`, at every creation site:
  windowed, composition hosting (`Windows.WebView2CompositionHosting`), the
  composition-to-windowed fallback, and the rebuild after a browser process
  failure; the options and `Environment10` are released after. Without
  `Environment10`, creation returns an error wrapping
  `edge.UnsupportedCapabilityError` through the existing error path.

**COM bindings** for both interfaces go in `internal/webview2/pkg/edge`,
beside those `Chromium` uses. The generated ones in
`internal/webview2/pkg/webview2` are unused and wrong: the HWND is passed to
`CreateCoreWebView2ControllerWithOptions` by pointer, and
`GetICoreWebView2Environment10` is defined on `ICoreWebView2`, not
`ICoreWebView2Environment` (#6309 fixed only the inherited slots). The `edge`
vtable is `iCoreWebView2Environment3Vtbl` plus the Environment4–9 slots in
`WebView2.idl` order, putting Environment10's three methods at slots 20–22,
pinned by a unit test.

### Impact on existing functionality

None when unset: the Linux helpers take a `NULL` session/context, and the
Windows windowed path calls the same `CreateCoreWebView2Controller` through a
helper.

## Non-Goals

- Clearing data persisted earlier or by other builds.
- Removing the WebView2 user data folder or other engine data; the
  folder-per-binary-name behaviour from #2255 is unchanged.
- Named or custom-location persistent profiles (WebView2 `ProfileName`, macOS
  14+ `WKWebsiteDataStore(forIdentifier:)`, WebKitNetworkSession with custom
  directories); room is left for them.
- Per-window ephemeral storage; encryption of data at rest; iOS and Android.

## Platform Considerations

| Platform | Mechanism | Scope | Limitations |
|---|---|---|---|
| macOS | `WKWebsiteDataStore.nonPersistentDataStore` | One store, all windows, process lifetime | None known; on every macOS version Wails v3 supports |
| Linux (WebKitGTK 6) | `webkit_network_session_new_ephemeral()` | One session, all windows, process lifetime | — |
| Linux (WebKitGTK 4.1) | `webkit_web_context_new_ephemeral()` | One context, all windows, process lifetime | — |
| Windows | WebView2 InPrivate (`ICoreWebView2ControllerOptions.IsInPrivateModeEnabled`) | Per WebView2 environment (user data folder) | Runtime ≥ 101.0.1210.39; user data folder still created |
| iOS, Android | — | — | Ignored |

**Windows.**

- *Minimum runtime:* SDK 1.0.1210.39, which made both interfaces stable,
  requires Runtime **101.0.1210.39** (May 2022)
  ([release notes](https://learn.microsoft.com/en-us/microsoft-edge/webview2/release-notes/sdk/1-0-1210-39),
  [ICoreWebView2Environment10](https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environment10)).
  Evergreen is long past it; failures mainly concern pinned fixed-version
  runtimes (`Windows.WebviewBrowserPath`) and machines with updates blocked.
- *User data folder:* the profile "will be created on disk or opened ... no
  matter InPrivate mode is enabled or not" (Microsoft). So
  `Windows.WebviewUserDataPath` (default `%APPDATA%\<binary name>`) is still
  created, honoured, and may hold engine data; only browsing data stays off
  the record.
- *Open PR #6258* adds per-window `WindowsWindow.WebviewUserDataPath` and
  `AdditionalBrowserArgs` (a per-window environment). `InPrivate` is per
  controller, so they compose; windows with different user data paths get
  different ephemeral sessions, as with persistent storage today (to be
  documented). Both touch a line or two of `webview_window_windows.go`; the
  second to merge rebases trivially.
- *`--incognito`* via `AdditionalBrowserArgs` is undocumented for InPrivate;
  the controller option is the supported API.

**Mobile:** iOS could use `nonPersistentDataStore`; Android has no ephemeral
store (storage is process-global; AndroidX profiles persist), so it would need
clearing. Both are left out rather than half-done, and documented as ignored.

## Pros/Cons

### Pros

- Stops unbounded data growth for apps whose origin changes per run.
- Uses each platform's first-party ephemeral store; deletes nothing Wails does
  not own.
- Additive and small: one `Chromium` field, two COM bindings, a few lines per
  platform.

### Cons

- Weaker on Windows (user data folder still created); must be documented.
- Hard failure on WebView2 runtimes older than 101.0.1210.39.
- Application-wide only; not on iOS or Android.

## Alternatives Considered

- **Struct/enum mode** (`Persistent | Ephemeral`, later `Named`): a two-value
  enum is a bool; a named profile is a separate axis (WebView2 allows a named
  InPrivate profile; macOS 14+ `dataStoreForIdentifier:` and custom-directory
  WebKitNetworkSessions persist). A later `Options.WebviewProfile string` can
  sit beside the bool. If preferred, only where the flag is read changes.
- **Per-window option** (`WebviewWindowOptions.EphemeralStorage`): would split
  cookies and `localStorage` between windows of one origin. Possible later as
  an override sharing one store among opting-in windows (WebKitGTK 4.1 would
  need the `wails` scheme on that context too); fits beside #6258.
- **Fall back to persistent with a warning** on old runtimes: would silently
  persist everything on the least-tested machines. If preferred, a small change
  in `withInPrivateOptions` (log via the existing error callback, call the
  non-options function), or a `Windows` field could choose.
- **Delete the store on exit/startup:** engine-owned, version-specific paths;
  useless after a crash; macOS's default store cannot be safely identified.
- **Fix and use `pkg/webview2`:** unused, and its environment type differs from
  `Chromium`'s, forcing conversions at each call; cleanup is separate.
- **Do nothing / document `--incognito` or cleanup:** nothing for macOS and
  Linux; Windows on an undocumented flag.

## Backwards Compatibility

Fully compatible: the zero value is current behaviour, no existing option
changes meaning, and the `edge` additions are internal.

## Security and Privacy

- Website data is "not persisted by the webview" — nothing more. It neither
  encrypts nor clears data persisted by earlier runs or builds, nor stops
  writes elsewhere (downloads, app files, OS caches).
- On Windows the user data folder is still created and may hold engine data;
  docs must not claim nothing is left on disk or any stronger guarantee.
- Isolation between windows is unchanged (one shared store).
- On an old WebView2 runtime it fails closed (no window), not open.

## Test Plan

- **Unit:** `TestOptions_Defaults` (zero value `false`);
  `TestEnvironment10VtableSlots`, `TestControllerOptionsVtableSlots`
  (Environment10 at 20–22, `put_IsInPrivateModeEnabled` at 6).
- **Windows:** `ICoreWebView2_13.get_Profile` →
  `get_IsInPrivateModeEnabled` is `TRUE` for windowed, composition-hosted and
  post-crash controllers, `FALSE` when unset; two windows share
  `localStorage`; a runtime older than 101.0.1210.39 fails with the
  unsupported-capability error.
- **macOS:** `websiteDataStore.persistent` is false and the store is shared;
  a cookie and `localStorage` key cross windows.
- **Linux:** `webkit_network_session_is_ephemeral()` /
  `webkit_web_context_is_ephemeral()` is true and shared by two windows; the
  `wails` scheme loads in a second window.
- **Manual:** `localStorage`/IndexedDB data is gone after relaunch when set,
  kept when unset.

## Reference Implementation

https://github.com/juggler-ai/wails/commit/14701a0ec00b3d270f8a57831103ca66b2eb84db
— option, `edge` bindings and vtable tests, all desktop paths, and docs in
`reference/application.md`. It changes only if review picks the struct, the
per-window override, or the fallback.

## Maintenance Plan

Stable, non-deprecated first-party APIs throughout. The care points are the
hand-written vtables (guarded by the slot tests) and routing any new WebView2
controller creation site through the InPrivate helper. Maintained by the
implementor with the other webview options.

## Conclusion

A supported one-line way to stop webviews accumulating data no run will read,
using each platform's own ephemeral store: off by default, additive, and
claiming only that data is not persisted by the webview.
