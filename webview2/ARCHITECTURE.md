# wails/webview2 — Architecture

This document records the design decisions behind the v2 WebView2 binding
generator. It is the reference for "why was this done this way?" questions.

## Goals

1. **Correct.** Generated COM bindings must respect the Windows ABI —
   no caller-stack pointer writes, no missing release methods, no truncated
   reference counts.
2. **Reproducible.** Given a frozen IDL + release-notes file, two runs of
   `webview2gen full` on different machines must produce byte-identical
   output. CI gates on this via `webview2gen verify`.
3. **Capability-aware.** The Wails runtime sees a stable Go API, but at
   runtime degrades gracefully when the installed WebView2 lacks an
   interface (we ship for a wide range of runtime versions).
4. **Single command refresh.** Refreshing the bindings for a new SDK
   release must be one CLI invocation, not a sequence of hand edits.

## High-level pipeline

```
                 Microsoft NuGet                          MicrosoftDocs/edge-developer
                       │                                            │
                       │ .nupkg (zip of headers + WebView2.idl)     │ index.md
                       ▼                                            ▼
        ┌──────────────────────────┐                ┌─────────────────────────┐
        │ internal/idl.Fetcher     │                │ internal/notes.Fetch    │
        │ extracts WebView2.idl    │                │ scrapes release notes   │
        │ caches in scripts/*.idl  │                │                         │
        └──────────────┬───────────┘                └────────────┬────────────┘
                       │                                         │
                       ▼                                         ▼
        ┌──────────────────────────┐                ┌─────────────────────────┐
        │ generator.ParseIDL       │                │ notes.Parse             │
        │  participle/v2 grammar   │                │ extract per-release     │
        │  types/typemap (17 pat.) │                │ interface mentions      │
        │  emits []GeneratedFile   │                │                         │
        └──────────────┬───────────┘                └────────────┬────────────┘
                       │                                         │
                       ▼                                         ▼
       pkg/webview2/*.go (306 files)            pkg/webview2/capabilities.go
```

Each box is one Go package under `webview2/scripts/`; `webview2gen` (the
CLI) wires them together.

## Design decisions

### Template-driven emitter (kept)

The generator emits Go source by composing `text/template` files in
`scripts/generator/types/templates/`. The spec considered moving to a
fully programmatic emitter (`fmt.Fprintf` + `bytes.Buffer`), but the
existing template surface is small (12 files), debuggable, and already
covers all the patterns we need. The 6 bug fixes were all expressible as
template edits; a rewrite would have been a much larger churn for the
same correctness outcome.

**Trade-off.** Templates are weakly typed and easy to break with whitespace
changes. Mitigated by the snapshot-style `testfiles/` fixtures: every
generator change re-runs the existing fixtures and diffs them.

### File-per-interface output

`pkg/webview2/` has 306 files because each interface and each enum lands
in its own file. Diffs from an SDK refresh stay readable, and parallel
editing of related but separate interfaces doesn't cause merge churn.

**Trade-off.** `ls pkg/webview2 | wc -l` is intimidating. Mitigated by
`doc.go` (the only hand-written file) carrying the "do not edit" notice
and the regeneration recipe.

### IDL inlined upfront

The fetcher pulls `WebView2.idl` directly out of the NuGet `.nupkg`. The
file is self-contained — it inlines the COM/OLE types it depends on, so
we do not need a preprocessor that resolves `import "objidl.idl"` etc.
If a future SDK switches to multi-file IDL, `internal/idl` is where the
inliner will live.

### Capabilities from release notes + IDL inventory

Each entry in `InterfaceSupportTable` pairs the SDK release that
introduced an interface with that release's *minimum WebView2 Runtime
version* — the value `SupportsInterface` compares against, because what a
running app can query is the installed runtime's browser version, not an
SDK version. Both values come from scraping the Markdown release notes
(`MicrosoftDocs/edge-developer`): the IDL itself records neither.

The scrape is regex-based and walks the notes oldest-first so the
earliest stable mention wins. Coverage is made total by construction: the
`capabilities` command parses the current IDL for the full interface
inventory; interfaces the notes never mention but that exist in the
oldest cached IDL predate the release-notes archive and become baseline
(always-supported) entries, and anything absent from both fails
generation loudly. A generated test additionally asserts every generated
interface has a table entry.

When Microsoft changes the notes format we will see test failures (and
generation failures from the inventory check) and update the regex; the
table is committed at `pkg/webview2/capabilities.go` so a broken scrape
can be fixed without forcing every consumer to wait.

**Trade-off.** A maintenance burden (every few months Microsoft may
restructure the notes). Mitigated by the test suite covering the parser
against the cached `test.md` snapshot, plus the loud inventory check.

### `pkg/edge/` is kept, not deleted

The legacy `pkg/edge/` package is what the Wails v3 runtime actually
imports today. Removing it would break the runtime in the same PR that
ships the new generator. The migration plan:

1. This branch lands the generator + tested `pkg/webview2/`.
2. A follow-up PR refactors `v3/pkg/application/...` (chromium.go) to
   import `pkg/webview2`.
3. After that lands, `pkg/edge/` is moved or removed.

`pkg/edge/` is therefore actively used and must not be touched by
generator runs. The generator only writes to `pkg/webview2/`.

### Two go modules in one repo

`webview2/go.mod` is the runtime module (`github.com/wailsapp/wails/webview2`).
`webview2/scripts/go.mod` is a dev-time module (`module updater`, a
legacy name kept to avoid a noisy rename). The split:

- Lets `scripts/` depend on developer-only packages (NuGet fetcher, HTTP
  client, participle/v2) without forcing those deps on every consumer.
- Lets `scripts/` use a separate, sometimes-older Go version (1.20) than
  the runtime module (1.24) without conflict.

If we ever want a single module, the rename is straightforward but
out-of-scope for this branch.

### Rejected approaches

- **Hand-maintain the bindings.** ~109 out-pointer methods to audit, and
  one new SDK per six weeks. Untenable. The whole point of v2 is
  automation.
- **`go-ole` or other COM libraries.** Adds a third-party dependency
  that we then need to track for the same set of fixes. The thin
  template-based wrapper is simpler to audit.
- **Monolithic `webview2.go`.** A single file would be 30k+ lines.
  Reviewing diffs for a new SDK release would be impractical.

## Future considerations

- **Generics for type-safe vtable calls.** Today every vtable invocation
  is `uintptr(unsafe.Pointer(...))`. Go generics could express the call
  shape per parameter type. Out of scope for v2; tracked as a follow-up.
- **Automatic version detection at startup.** `webviewloader` already
  reports the runtime version. We could populate a process-wide cache
  so `SupportsInterface` calls don't need to re-pass the string.
- **Cross-platform compile checks.** `pkg/webview2/` is `//go:build
  windows` so non-Windows CI can't catch regressions. Could add a
  `linux` shim package that satisfies the same surface using build tags.
- **Windows VM CI.** The generator unit tests run on Linux. Smoke tests
  for the COM bindings need a real WebView2 runtime — tracked as
  WAI-297; not in this branch.
