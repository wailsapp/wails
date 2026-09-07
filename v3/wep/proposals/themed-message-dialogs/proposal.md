# Wails Enhancement Proposal (WEP)

## Themed Windows message dialogs

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Julian Storer  
**Created**: 2026-08-04  
**Discussion**: [optional link to any prior discussion of the idea]  
**Implementor**: Julian Storer  
**Target**: Wails v3

## Summary

Render Windows message dialogs with the modern comctl32 v6 `TaskDialogIndirect`
API instead of the legacy Win32 `MessageBox`. This makes dialogs render themed
to match the OS, and lets them honour the custom buttons a dialog was built
with. Falls back to the current `MessageBox` when the modern API is unavailable.

## Motivation

Two problems with the current `MessageBox` path:

1. **Appearance.** `MessageBox` always draws in the classic, unthemed ("Windows
   Classic") style regardless of the OS visual style, so Wails dialogs look
   dated next to the rest of the system.
2. **Custom buttons don't work.** `MessageBox` can only show the fixed
   Yes/No/OK/Cancel captions. A dialog built with custom button labels — e.g.
   `AddButton("Discard")` — returns a standard result string that never matches
   the caller's labels, so the registered `OnClick` handler never fires. This is
   effectively a bug: the public custom-button API silently does nothing on
   Windows.

`TaskDialogIndirect` is the API Windows itself uses for modern dialogs today. It
renders themed, supports a bold main instruction over the message body, and
supports arbitrary custom buttons with real click results.

## Detailed Design

Message dialog display gains a modern path with a graceful fallback:

- `showModern()` builds and shows a `TaskDialogIndirect` dialog. On **any**
  failure — most importantly when the process has no Common-Controls v6
  side-by-side assembly, so the API is unavailable — it returns `false`.
- `showLegacy()` is the existing unchanged `MessageBox` path, used whenever
  `showModern()` returns `false`.

So the change is additive: environments that can render a task dialog do;
everything else behaves exactly as today.

### Custom buttons

Custom buttons are assigned stable IDs and mapped back to their `Button`, so the
correct `OnClick` handler runs and the result matches the caller's label. The
`IsDefault`/`IsCancel` flags and Esc / title-bar-close are honoured.

### w32 binding

Adds a `pkg/w32` binding for `TaskDialogIndirect`. `TASKDIALOGCONFIG` and
`TASKDIALOG_BUTTON` are `#pragma pack(1)`, so both are serialized by hand into
byte buffers at their exact packed offsets rather than as naturally aligned Go
structs; pointer-sized fields are written at native width so the layout stays
correct on 386 as well as amd64/arm64, and buried pointers are held live across
the call with `runtime.KeepAlive`. The config-to-task-dialog translation is
factored into a pure `buildTaskDialogConfig` function so it can be unit-tested
(button ID assignment, default/cancel mapping, per-type common buttons and
icons, app-icon fallback).

## Non-Goals

- No change to the public dialog API (`AddButton`, `Question`, message/title,
  icons). This changes how existing dialogs are rendered and dispatched, not how
  they are constructed.
- Not a custom/HTML dialog system; that is a separate concern.
- No change to macOS or Linux dialogs.

## Platform Considerations

- **Windows**: modern themed rendering when Common-Controls v6 is available;
  automatic fallback to `MessageBox` otherwise (e.g. a process without the v6
  side-by-side manifest). Packed-struct handling is 32/64-bit safe.
- **macOS / Linux**: unaffected; this is a Windows-only rendering change.

## Pros/Cons

### Pros

- Dialogs render themed and native, matching the rest of Windows.
- Custom buttons and their `OnClick` handlers work on Windows for the first
  time — fixes a real gap in the existing API.
- Degrades gracefully; no hard dependency on comctl32 v6.
- Core translation is pure and unit-tested.

### Cons

- Adds a non-trivial `w32` binding with hand-packed structs to maintain.
- Dialog appearance changes for existing apps (see Backwards Compatibility).

## Alternatives Considered

- **Keep `MessageBox`, only fix custom-button mapping.** Not possible in
  general: `MessageBox` cannot render arbitrary button captions, so custom
  buttons fundamentally require a different API.
- **Depend on comctl32 v6 unconditionally** and drop the `MessageBox` path.
  Rejected: a process without the v6 assembly would lose dialogs entirely; the
  fallback keeps Wails robust.
- **A managed struct via cgo/syscall marshalling** instead of hand-packed
  buffers. Rejected: `TASKDIALOGCONFIG` is `#pragma pack(1)` with embedded
  pointers, so explicit packed serialization is the reliable cross-arch choice.

## Backwards Compatibility

The public API is unchanged. Two observable changes for existing apps:

- Dialogs look different (themed instead of classic) where comctl32 v6 is
  available. This is a visual change, not an API change.
- Custom-button dialogs now return the caller's label and fire the registered
  handler, where previously they returned a fixed string and fired nothing.
  Code that depended on the old (broken) result strings would be affected, but
  that path never did what the API promised.

Where the modern API is unavailable, behaviour is byte-for-byte the current
`MessageBox` path.

## Security and Privacy

None. `TaskDialogIndirect` is a standard user32/comctl32 UI call with no new
capability, permission, or data-handling implication.

## Test Plan

- Unit tests cover `buildTaskDialogConfig`: button ID assignment, default/cancel
  mapping, per-type common buttons and icons, and app-icon fallback. (Present in
  the reference implementation.)
- Manual: standard Info/Warning/Question/Error dialogs render themed; a dialog
  with custom buttons returns the right label and fires the right handler;
  Esc / title-bar-close map to the cancel button.
- Fallback: confirm a process without the v6 assembly still shows a working
  `MessageBox`.

## Reference Implementation

A working reference implementation exists, including the `pkg/w32`
`TaskDialogIndirect` binding, the modern/legacy dialog paths, and unit tests for
the config builder. The implementation PR will be linked here once opened.

## Maintenance Plan

The main ongoing surface is the hand-packed `TASKDIALOGCONFIG`/`TASKDIALOG_BUTTON`
serialization, which is pinned to a stable Win32 ABI and covered by unit tests,
so drift is unlikely and would be caught. Maintained alongside the rest of the
Windows dialog code.

## Conclusion

Switching to `TaskDialogIndirect` modernises Windows dialog appearance and makes
custom buttons actually work, while preserving the existing `MessageBox` path as
a safe fallback and leaving the public API untouched.
