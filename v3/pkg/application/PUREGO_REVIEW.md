# Purego macOS backend — code review findings

Review of the CGO-free macOS backend (branch `worktree-mac-purego`, commits
`978afbf47..501f74989`) against the cgo backend as the behavioural reference.
Nine parallel review passes (foundation/FFI, application core, window/webview,
dialogs, menus/systray, asset server, leaf features, services, cross-cutting
build/ABI audit); all findings below were verified against the code, and the
FFI claims against the purego v0.10.1 sources.

Severity: **critical** = crash or broken core behaviour on common paths;
**major** = real bug, leak, or user-visible parity break; **minor** = edge
case, cosmetic parity, or robustness gap.

## Fix status (2026-07-02, commits `94a8391d3..b24f1c387`)

**Fixed:** all three criticals (C1 frameless close via direct `close`;
C2 `markAsDestroyed` in `destroy()`; C3 scheme-task stop race closed by
confining check+send atomically to the main thread with `dispatch_sync`) and
all majors M1–M14 (block `Release()` sweep everywhere, dialog map race +
non-modal windowless panels, file-input open panel, renderer-crash reload,
navigation/willClose/universal-link events, full geometry sweep incl.
`options.Screen`, construction parity + window-level drag destination,
registry hygiene + body-stream close on stop, notifications typeHint +
retained block captures, systray deferred `showMenu`, selector guards,
per-window teardown in `windowWillClose:`, foundation dlopen/superclass
panics). Most minors fixed in the same commits (show/focus, zoom semantics,
restore, isVisible, execJS guards, setTitle, min/max sizes, center,
maximise-button state, hidden else-branch, startDrag anchor, appInit pool,
terminate re-entrancy, post-run release, screen-ID signedness, go.mod,
pressure float, parseHexColor, doc comment).

**Leak follow-up (2026-07-03):** the "cgo has the same leak" category was
subsequently eliminated as well — menu `setBitmap` NSImage, the
`Menu.Update()` rebuild (old impls destroyed, submenus released, idempotent
`destroy()`), dialog icons, About/message NSAlert lifetime, app icon,
translucent-backdrop NSVisualEffectView, notifications content/payload, the
second-instance URL handler, and autorelease-pool wrapping for every
entry point that runs ObjC on a plain goroutine thread (clipboard,
`pkg/mac`, `isDarkMode`/`name`, all notifications APIs). The same leaks were
fixed in the **cgo** backend on branch `fix/darwin-cgo-leaks` (PR #5714).

**Left as-is (deliberate):** the single-instance URL-capture timeout start
point (ms-level); `goStringFromC` double copy and the vet `unsafeptr`
notices (benign purego idiom); `objc.Send[T]` per-call `RegisterFunc`
overhead (perf only); `nsString` invalid-UTF-8/NUL behaviour (cgo parity);
the pre-existing `menuitem_selectors_darwin.go` tag asymmetry.

**Validated:** both backends build (arm64 + amd64 cross), vet clean, and a
CGO-free lifecycle app (create titled → Close → recreate frameless → Close →
recreate → Quit) runs to a clean exit on a real macOS 26 desktop;
`examples/window` renders and serves assets.

---

## Critical

### C1. Frameless windows can never be closed
`webview_window_darwin_purego.go:790` — `close()` sends `performClose:`. cgo
closes directly (`[window close]` in `windowClose`). Frameless windows are
created with `borderless|resizable|miniaturizable` (`:417-418`) — no
`NSWindowStyleMaskClosable` — and AppKit's documented behaviour for
`performClose:` on a window with no close button is to beep and return.
Worse: the shared `WindowClosing` handler has already run `markAsDestroyed` +
`Window.Remove` before `impl.close`, so the NSWindow stays on screen but is
unreachable through Wails — a zombie window. Fix: mirror cgo — send `close`
directly (the `unconditionallyClose` flag already provides the should-close
gate for titled windows).

### C2. `destroy()` never calls `markAsDestroyed()` → use-after-free
`webview_window_darwin_purego.go:791-798` cleans the maps and sends `close`,
but unlike cgo (`webview_window_darwin.go:1687-1692`) never calls
`w.parent.markAsDestroyed()`. NSWindow defaults to `releasedWhenClosed=YES`,
so `close` deallocs the window; with `isDestroyed()` still false, every public
API guard is bypassed and a subsequent `SetTitle`/`Focus`/second `Destroy()`
messages freed memory.

### C3. Scheme-task stop race → uncatchable NSException
`responsewriter_darwin_purego.go:131-188`. Each write does
`if taskStopped(task) { return } ` then messages the task — a check-then-act
race against `webView:stopURLSchemeTask:` arriving on the main thread while
the asset handler goroutine is mid-response. If the send lands after WebKit
marks the task stopped, WebKit raises `NSException "This task has already
been stopped"` synchronously on the calling (Go) thread. cgo survives this
exact case with `@try/@catch` (`responsewriter_darwin.go:21-33`); pure Go has
no net, so the process aborts. Triggered by navigation/reload/close with
in-flight custom-scheme requests — load-dependent but will fire in
production. The registry narrows the window; it cannot close it (WebKit sets
its internal stopped flag before delivering the delegate callback). Robust
fixes: dispatch the check+send atomically to the main thread (stop is also
delivered there), or add a tiny ObjC trampoline class registered at runtime
whose IMP wraps the send in a native `@try/@catch`.

---

## Major

### M1. Every `objc.NewBlock` is leaked — unbounded growth on the hot path
purego's `NewBlock` returns a +1 heap block owned by the caller and pins the
Go closure in a global cache until refcount 0; `dispatch_async` /
`completionHandler:` copy and release only their own reference, so without an
explicit `block.Release()` nothing is ever freed. Zero `Release()` calls
exist in the codebase. Per-call leak sites:
- `mainthread_darwin_purego.go:38,68` (`dispatchOnMainThread`/`runOnMain`) —
  behind **every** `InvokeSync`/`InvokeAsync`, window op, `execJS`, and the
  50 ms drag-throttle tick;
- `dialogs_darwin_purego.go:241,472,569` (per sheet dialog);
- `webview_window_darwin_purego.go:1210` (per `attachModal`);
- `services/notifications/notifications_darwin_purego.go:682,699,862,919,951`
  (per notifications API call);
- `services/dock/dock_darwin_purego.go:140` (per dock call off-main);
- `systemtray_darwin_purego.go:237` (per tray create/destroy cycle).
This is heap growth (~100-200 B/call), not callback-slot exhaustion (the
NewCallback budget was audited separately and is safe). Fix is one line per
site: `block.Release()` immediately after the ObjC call that copies it
returns (dispatch_async and beginSheet… copy synchronously).

### M2. Data race on dialog response maps → fatal runtime crash
`dialogs_darwin_purego.go:451-461,557-559` — the panel completion `finish`
hops to a bare goroutine that reads/`delete`s `openFileResponses` /
`saveFileResponses` (plain maps, no mutex), while `show()` writes them on the
main thread. cgo is race-free because both sides run on the main thread.
Two dialogs in close succession can throw the unrecoverable
"concurrent map read and map write". Fix: do the map access on the main
thread like cgo (only the channel send needs to leave it), or add a mutex.

### M3. `<input type="file">` is broken — WKUIDelegate open panel missing
purego sets the UIDelegate (`webview_window_darwin_purego.go:465`) but
implements no `runOpenPanelWithParameters:` (cgo:
`webview_window_darwin.m:810-824` shows an NSOpenPanel). File-picker clicks
in the page do nothing.

### M4. No recovery from WebContent process crashes
`webViewWebContentProcessDidTerminate:` is not registered. cgo
(`webview_window_darwin.m:799-808`) emits the event **and reloads** — its
comment notes the renderer dies "notably after macOS sleep" and without the
handler the window is a blank unresponsive page forever.

### M5. Missing navigation/window events
- Only `didFinishNavigation` is registered; cgo also forwards
  `didStartProvisionalNavigation`, `didReceiveServerRedirectFor…`,
  `didCommitNavigation` (`webview_window_darwin.m:778-797`) — those
  `events.Mac.WebView*` events never fire.
- `windowWillClose:` is absent from the notification map
  (`webview_window_darwin_purego.go:93-112`) — `EventWindowWillClose` never
  fires (cgo: `.m:638-642`).
- Universal links: `application:continueUserActivity:restorationHandler:` is
  not implemented in the app delegate (cgo:
  `application_darwin_delegate.m:20-29`) — universal-link activations are
  silently dropped (`ApplicationLaunchedWithUrl` never fires for them).

### M6. Multi-monitor geometry is wrong
- `size()`/`width()`/`height()` return the **contentView** size; cgo returns
  the window **frame** size (`webview_window_darwin_purego.go:898-909` vs cgo
  `:623-637`). Since `setSize()` sets the frame, Set/Get is asymmetric and a
  `bounds()`→`setBounds()` round-trip shrinks the window by the titlebar
  height every time.
- `setPosition()`/`position()` compute against `[NSScreen mainScreen]` (the
  key window's screen); cgo uses the **primary** screen
  (`[[NSScreen screens] firstObject]`). Wrong coordinates whenever the window
  sits on a secondary display.
- `relativePosition()`/`setRelativePosition()` are aliased to the absolute
  variants (`:1037-1038`); cgo computes relative to the window's own screen
  (the #5408 fix). Broken on any non-primary screen.
- `centerOnScreen(screen)` ignores its argument (`:1039`) and `run()` ignores
  `options.Screen` (`:306-310`) — `WebviewWindow.SetScreen()` doesn't work
  and windows always open on the main screen.
- `setMinSize`/`setMaxSize` treat values as content size; cgo converts via
  `contentRectForFrameRect` and also sets `minSize`/`maxSize` (minor sizing
  skew, grouped here for the geometry sweep).

### M7. Window construction / drag parity
- The NSWindow subclass lacks cgo's init overrides (`.m:22-30`):
  `backgroundColor clearColor`, `opaque NO`, `movableByWindowBackground YES`,
  `alphaValue 1.0` — transparent/frameless rendering and
  background-drag behaviour differ.
- Frameless rounded corners missing (cgo sets `layer.cornerRadius = 8.0` on
  the content view; `webview_window_darwin.go:58-61`).
- cgo registers **every** window for `NSFilenamesPboardType` and implements
  window-level `draggingEntered/exited/performDragOperation`
  (`.m:201-208,254-304`), so `EventWindowFileDragging*` fire even with
  `EnableFileDrop=false`. purego installs its overlay only when
  `EnableFileDrop=true` (`:470-472`) — no drag events otherwise.

### M8. Stopped-task registry leaks and can kill unrelated future requests
`responsewriter_darwin_purego.go:117-127` + `request_darwin_purego.go:286-295`.
If `stopURLSchemeTask:` arrives after the request already finished and
`Close()` ran, `MarkTaskStopped` stores an entry keyed on a released pointer
that is never deleted: (a) slow unbounded growth; (b) when the heap reuses
that address for a later WKURLSchemeTask, the new request sees
`taskStopped()==true` for its whole lifetime and silently never loads.
Fix: `NewRequest` should clear any stale registry entry for its task pointer
on registration.

### M9. Notifications: dropped type hint + latent use-after-free
- `notifications_darwin_purego.go:753-756` passes the literal string
  `"UNNotificationAttachmentOptionsTypeHintKey"` as the options key — that's
  the *symbol name*; the framework constant's value is `"typeHint"` (verified
  at runtime). UTI type hints for attachments are silently ignored.
- `:882-975` — the category register/remove completion blocks capture
  **autoreleased** ObjC objects (`nsCategoryID`, `newCategory`) as raw
  uintptrs with no retain (cgo blocks retained them via `_Block_copy`).
  If the caller's autorelease pool drains before the async handler runs
  (e.g. main-thread caller + handler arriving after the 5 s timeout), the
  block messages freed memory. Retain before dispatch, release in the block.

### M10. Systray menu opens re-entrantly inside the click dispatch
`systemtray_darwin_purego.go:451-493` — cgo's `showMenu` is always
`dispatch_async` (next runloop turn); purego uses `dispatchOnMainThread`,
which runs **inline** when already on main. With the default tray config the
synthesized `mouseDown:` + nested menu-tracking runloop runs inside the real
right-click's `sendAction`, and `OpenMenu()` now blocks until the menu is
dismissed instead of returning after enqueueing. Fix: force async dispatch.

### M11. Unguarded selectors that abort the process
Per the port's own rule, a missing selector is an uncatchable NSException:
- `webview_window_darwin_dev_purego.go:16-23` — `openDevTools()` sends the
  **private** `_inspector` with no guard; cgo wraps it in
  `@available(macOS 12,*)` *and* `@try/@catch`. Private API can vanish in any
  WebKit update → add `respondsTo` guard (highest risk of this group).
- `webview_window_darwin_purego.go:1073` `printOperationWithPrintInfo:`
  (macOS 11+), `:1278` `setToolbarStyle:` (11+),
  `application_darwin_purego.go:279` `controlAccentColor` (10.14+),
  `screen_darwin_purego.go:93` `localizedName` (10.15+). All currently
  unreachable-in-practice because Go ≥1.23 binaries require macOS 11+, but
  each should be guarded or carry a floor comment so a backport doesn't trip.

### M12. Per-window object graph never released
`createWindow` (`webview_window_darwin_purego.go:428-473`) allocs the NSView,
WKWebViewConfiguration, user content controller, WKWebView and delegate and
never releases them; cgo autoreleases each and its `dealloc` removes the
`"external"` script message handler specifically to break the retain cycle.
purego never removes the handler, so the whole window/webview/delegate graph
survives every close. `setUseToolbar` also leaks an NSToolbar per call
(`:1269`). Additionally, `windowImplCache`/`windowDisableEscape`/
`dragViewToWindowID` are only cleaned in `destroy()`, not on the normal
`Close()` path — and `windowDisableEscape` is keyed by the freed NSWindow
*address*, so a future window at the same address silently inherits
`DisableEscapeExitsFullscreen`. **Caution when fixing:** today's leaks are
what keep the weak `NSWindow.delegate`/`navigationDelegate` references alive;
releasing the graph requires explicitly retaining the delegate.

### M13. Window-less open/save dialogs are app-modal
`dialogs_darwin_purego.go:467,564` use `runModal` when no window is attached;
cgo uses non-modal `beginWithCompletionHandler:` — `show()` returns
immediately and other windows stay interactive. Under purego the entire app
is modal for the duration.

### M14. Foundation robustness
- `darwin_purego_cocoa.go:61` — `Dlopen` failures are silently swallowed; if
  AppKit/WebKit fail to map, every `class()` is nil and all sends become
  silent no-ops (no window, no diagnostic). Panic/log for the big three;
  UniformTypeIdentifiers may stay tolerant.
- `darwin_purego_cocoa.go:269-276` — `registerDelegateClass` doesn't check
  `objc.GetClass(super) != 0`. A missing superclass silently creates a *root*
  class (that's how ObjC makes root classes), and the first `alloc` then dies
  in `doesNotRecognizeSelector:` — the intended fail-fast panic never fires.

---

## Minor

**Window behaviour**
- `show()` also calls `activateIgnoringOtherApps:` — every `Show()` steals
  app focus; cgo only does `makeKeyAndOrderFront:`. `focus()` always
  activates; cgo activates only when inactive and also calls `makeKeyWindow`.
- `zoom()` calls `zoomReset()`; cgo's `zoom()` is `[window zoom:nil]`
  (maximise toggle) — public `Zoom()` changes meaning. Webview zoom uses
  `pageZoom` vs cgo `magnification`; the `>=1.0` clamp and step sizes differ.
- `restore()` only deminiaturizes; cgo also exits fullscreen and un-zooms.
- `isVisible()` uses `NSWindow.isVisible`; cgo uses
  `occlusionState & NSWindowOcclusionStateVisible`.
- `execJS` lost cgo's `performingShutdown`/`isDestroyed`/nil guards and is
  now synchronous (`runOnMain` vs `InvokeAsync`) — JS is dispatched into
  closing windows and callers block.
- `setTitle` lost the `if !Frameless` guard; `setSize` lost `animate:YES`;
  `center()` uses `[NSWindow center]` (top-third bias) vs cgo's exact
  `visibleFrame` centring; `setMaximiseButtonState` skips
  `effectiveZoomButtonState(state, FullscreenButtonState)`; windows created
  `Hidden` never get `DisableShadow`/`AlwaysOnTop` applied on first key
  (cgo's `run()` else-branch was dropped).
- `startDrag` uses `NSApp.currentEvent` instead of cgo's retained original
  mouse-down event — JS-initiated drags anchor at the current cursor
  position (slight jump) and can no-op if another event intervened.

**Asset server**
- `stopURLSchemeTask:` doesn't close the request's `HTTPBodyStream` (cgo
  does, `.m:380-388`) — cancelled streamed uploads keep being consumed until
  the handler finishes.

**Dialogs**
- Empty-title About dialog shows NSAlert's default "Alert" text (cgo sets "").
- Open-panel filter delegate uses `hasDirectoryPath` (string-based); cgo uses
  `fileExistsAtPath:isDirectory:` — symlinks to directories are greyed out.

**Systray / menus**
- `pressure:` argument passed as `float64`; the parameter is a C `float`, so
  the callee reads ~0.0. Should be `float32(1.0)`.
- `parseHexColor` accepts colours without the leading `#`; cgo's sscanf
  applies no colour then.
- The event-monitor block (and captured tray struct) leaks per
  create/destroy cycle (`removeMonitor:` releases only AppKit's copy).

**Application**
- `appInit()` creates autoreleased objects before any pool exists (one-time
  "autorelease with no pool" leaks); cgo used compile-time constants. Wrap in
  `withAutoreleasePool`.
- `applicationShouldTerminate:` lacks cgo's `shuttingDown` re-entrancy guard
  (currently mitigated by `App.cleanup()`'s own idempotence).
- cgo releases the delegate and calls `abortModal` after `[NSApp run]`
  returns; purego doesn't.

**Leaf features**
- Screen.ID formats the display ID unsigned; cgo's `%d` prints it signed —
  persisted IDs differ when the top bit is set.
- Second-instance URL-capture timeout starts before `[NSApp run]` rather than
  in `applicationDidFinishLaunching` — marginally shorter effective window.
- Notifications availability check is class-existence (10.14+) where cgo
  requires macOS 11; on 10.14/10.15 the presentation options used are
  11-only values (dead code on the Go 1.25 floor, but divergent).

**Foundation / hygiene**
- `goStringFromC` uses the uintptr-arithmetic pattern vet flags (unsafeptr)
  and copies twice; purego's own `GoString` idiom copies once. 11 vet
  "possible misuse of unsafe.Pointer" notices across the purego files —
  benign, but CI running vet will need suppression.
- `nsString` returns nil for invalid UTF-8 / truncates at interior NUL
  (parity with cgo, but downstream sends silently no-op on nil).
- Autoreleased helpers leak when called from plain goroutine threads with no
  pool — document/enforce `withAutoreleasePool` at non-main call sites.
- `objc.Send[T]`/`get[T]` re-runs `RegisterFunc` (reflection) on every call —
  pure overhead on the hottest paths, incl. `isOnMainThread()`.
- `sendSuper` resolves the superclass dynamically (purego behaviour): if a
  delegate instance is ever KVO-observed, super-dispatch recurses infinitely.
  Document the constraint.
- `v3/go.mod` still marks `github.com/ebitengine/purego` as `// indirect`
  (`go mod tidy` would fix); package doc example in
  `darwin_purego_cocoa.go:14-15` references non-existent `alloc()`/`init()`
  methods.
- `menuitem_selectors_darwin.go` lacks `!ios` where its purego twin has it
  (pre-existing asymmetry, behaviour-preserving).

---

## Verified clean (audited, no findings)

- **Build tags:** every darwin cgo file carries `!purego`; all purego twins'
  qualifiers (`!ios`/`!server`/production/devtools) mirror correctly.
- **Builds:** `CGO_ENABLED=0 -tags purego` passes for `./pkg/...` and
  `./internal/...` on arm64 **and** amd64, plus `purego,production` /
  `purego,devtools` combos and `examples/window`; the default cgo build is
  unregressed. (Known unrelated `build_assets/ios` failure ignored.)
- **Callback budget:** ~1 direct `purego.NewCallback` + a few dozen
  `sync.Once`-guarded class registrations — far under purego's 2000-slot cap.
  Block signatures share ~6 slots total.
- **FFI/ABI:** string args (copied + kept alive), bool, unsafe.Pointer,
  float args in FP registers, CGRect/NSRange/NSOperatingSystemVersion
  struct returns (stret on amd64, x8/HFA on arm64) all verified correct in
  purego v0.10.1.
- **Parity verified exact:** key-code and menu-role/selector tables, modifier
  masks, screen geometry math (Y-flip, visibleFrame, scale), Carbon hot-key
  ABI + FourCharCodes, single-instance lock/notification flow, clipboard,
  autostart (SMAppService class-existence guard ≡ cgo's `@available(13,*)`),
  dock badge behaviour, UN framework constants (vs SDK headers), dialog
  message/button/callback semantics, file-drop overlay coordinate math,
  window-notification event IDs (all 53 reflected names resolve),
  `windowShouldClose` logic, script-message and URL-scheme request
  extraction (purego actually fixes two cgo bugs: dangling `UTF8String`
  after pool drain, and a potential empty-slice panic in stream reads).

## Suggested fix order

1. C1 + C2 (window close/destroy) — small, crash-class.
2. M1 block `Release()` sweep — mechanical, unblocks long-running apps.
3. C3 + M8 (scheme-task stop race + registry hygiene) — needs the
   main-thread-confined check or an ObjC trampoline with `@try/@catch`.
4. M2 dialog map race; M3/M4/M5 missing delegate methods (file upload,
   renderer-crash reload, events, universal links).
5. M6 geometry sweep (size/position/relative/screen targeting) — restores
   multi-monitor correctness.
6. Remainder in any order; minors opportunistically.
