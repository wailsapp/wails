# Changelog

All notable changes to the webview2 bindings (`v3/internal/webview2`).
Entries for SDK regenerations are produced by `webview2gen changelog`;
engineering notes are added by hand. The bindings are part of the v3
module and ship with v3 releases; historical `webview2/v1.x` tags belong
to the retired standalone `github.com/wailsapp/wails/webview2` module.

## SDK 1.0.3967.48 — generator overhaul (2026-06-10)

Complete rework of the IDL→Go generator and full regeneration of
`pkg/webview2`. Originally staged as a `/v2` release of the standalone
webview2 module; that release was superseded by folding the package into
the v3 module (`v3/internal/webview2`), where it versions with v3 itself.

### Fixed — calling convention / ABI

- **Derived interface vtbls were missing every inherited slot.** All 89
  derived interfaces (`ICoreWebView2_2`…`_28`, `Controller2-4`,
  `Settings2-9`, `Profile2-8`, `Environment2-15`, `Frame2-7`, …) dispatched
  every method to the wrong vtbl slot. Vtbls now embed their base
  interface's vtbl recursively, reproducing the flat COM layout.
- `[in] double` parameters were truncated (`uintptr(v)` turned `1.5` into
  `1`). They now pass the IEEE-754 bit pattern (12 methods, including
  `PutZoomFactor`), with the correct two-word encoding on 386.
- 8-byte by-value aggregates (`EventRegistrationToken` in all 65 `Remove*`
  event methods, `POINT`) were passed by address instead of by value — event
  handlers could never be unregistered. Now passed by value per the Win64
  ABI, with per-architecture encodings for arm64 and 386.
- `COREWEBVIEW2_COLOR` is packed by value into a single word
  (`PutDefaultBackgroundColor` previously received a stack address).
- `RECT` parameters are pointer-to-copy on amd64, a register pair on arm64
  and four stack words on 386 (previously wrong on arm64/386).
- `VARIANT` had a 16-byte layout; the win64 ABI size is 24 bytes. `Val` is
  now `[2]uintptr` (24/16 bytes on 64/32-bit).
- `[in] UINT64` splits into two stack words on 386.

### Fixed — API shape

- **QueryInterface helpers target the right object.** `Get<Interface>()`
  helpers were all emitted on `*ICoreWebView2`; obtaining
  `ICoreWebView2Controller2` etc. was impossible. Helpers now hang off the
  root of each interface's inheritance chain
  (`(*ICoreWebView2Controller).GetICoreWebView2Controller2()`, …).
- `SetCustomSchemeRegistrations` took `**ICoreWebView2CustomSchemeRegistration`
  and passed its address (a triple pointer). It now takes
  `[]*ICoreWebView2CustomSchemeRegistration`.
- `GetCustomSchemeRegistrations` marshalled the out-array into a struct
  value. It now returns `[]*ICoreWebView2CustomSchemeRegistration`, copied
  from the CoTaskMem array, which is then freed.
- `GetAllowedOrigins` returned `*string` holding a raw native pointer. It
  now returns `[]string`, with every string and the array CoTaskMem-freed.

### Changed — capabilities

- `InterfaceMinimumVersion` (map of SDK versions) is replaced by
  `InterfaceSupportTable` (`SDKVersion` + `MinRuntimeVersion` per
  interface). `SupportsInterface` now compares the installed runtime's
  browser version against a *runtime* minimum — the old comparison against
  SDK versions (1.0.x) was vacuously true for every real runtime.
- The table covers all 266 generated interfaces (plus release-notes-only
  interop interfaces); completeness is enforced at generation time and by a
  generated test. `Capability`, `AllCapabilities` and `HasCapability` were
  removed.

### Added

- Generated test suite: 267 `*_gen_test.go` files (1,019 tests) emitted by
  the generator — per-method marshalling round-trips against fake recording
  vtbls, QueryInterface IID checks, handler Invoke dispatch tests, and
  package-wide vtbl/struct-layout and capability-coverage assertions.
- `webview2gen changelog` — diffs two SDK IDLs into a changelog entry.
- Incremental generation: unchanged files are not rewritten; files the IDL
  no longer produces are removed.
- This changelog.

### Earlier history

The module previously tracked `github.com/wailsapp/go-webview2` (archived);
see that repository for pre-monorepo history. The hand-written
`pkg/edge`/`webviewloader` layers are unchanged by the generator work.
