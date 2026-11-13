# Wails Enhancement Proposal (WEP)

## Themed Windows message dialogs

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Julian Storer  
**Created**: 2026-08-04  
**Discussion**: Builds on [#6308](https://github.com/wailsapp/wails/pull/6308) (merged; fixes [#5794](https://github.com/wailsapp/wails/issues/5794)). Related: [#5193](https://github.com/wailsapp/wails/pull/5193), [#6250](https://github.com/wailsapp/wails/pull/6250)  
**Implementor**: Julian Storer  
**Target**: Wails v3

## Summary

#6308 shows Windows Question dialogs that have custom buttons with
`TaskDialogIndirect`. Every other message dialog — Info, Warning, Error, and
Questions without custom buttons — still goes through `MessageBox`, which draws
the classic, unthemed box. This WEP extends #6308's TaskDialog path to every
message dialog, so they all render in the current visual style, with the title
as a bold heading over the message. `MessageBox` stays as the fallback.

## Motivation

`MessageBox` ignores the OS visual style, so most Wails dialogs look dated next
to the rest of Windows — and, since #6308, inconsistent with the custom-button
Questions in the same app, which are already themed. The TaskDialog machinery
(binding, packed encoder, button mapping, topmost callback, fallback) is now in
master; using it for the remaining dialogs is a small change on top of it.

## Detailed Design

All of this sits inside #6308's `showTaskDialog` in
`v3/pkg/application/dialogs_windows.go`; no new binding or code path.

- **Remove the gate.** `show()` calls `showTaskDialog` for every message dialog,
  not only `QuestionDialogType` with buttons.
- **Standard buttons.** A dialog with no buttons of its own gets the standard
  set for its type through `dwCommonButtons` (`TDCBF_*`): OK, or Yes/No for a
  Question — the same buttons `calculateMessageDialogFlags` gives `MessageBox`.
  Their results are the same `IDOK`/`IDYES`/`IDNO` values. Dialogs with custom
  buttons are unchanged from #6308 (IDs from 1000, default/cancel mapping).
- **Standard icons.** `pszMainIcon` gets `TD_INFORMATION_ICON`,
  `TD_WARNING_ICON` or `TD_ERROR_ICON` by type. TaskDialog has no question icon,
  and `MessageBox` shows none for Questions today (`MB_YESNO` only), so a
  Question gets no icon. The app icon (resource 3) still takes precedence when
  requested and present, as in #6308.
- **Main instruction.** The title is passed as `pszMainInstruction` (bold
  heading) as well as `pszWindowTitle`.
- **Esc / close box.** `TDF_ALLOW_DIALOG_CANCELLATION` is set exactly when
  Esc and the close box would dismiss the `MessageBox` the dialog replaces:
  when it has a cancel button (as in #6308), or when it is a standard OK
  dialog (Info, Warning or Error with no buttons of its own). A standard
  Yes/No question still requires an answer, as `MB_YESNO` does, and so does a
  dialog with custom buttons and no cancel button.
- **Unchanged from #6308:** `HWND_TOPMOST` on `TDN_CREATED` (preserving
  `MB_SYSTEMMODAL`), positioning relative to the parent, and the `MessageBox`
  fallback with caption/position button mapping.

The per-type choice is a pure `taskDialogStandard(DialogType)` helper beside
`planTaskDialog`; `encodeTaskDialogConfig` gains the `dwCommonButtons` and
`pszMainInstruction` fields it currently writes as zero.

### Fallback policy

`TaskDialogAvailable()` is checked first, so the expected "no Common-Controls v6
manifest" case never calls the API. If `TaskDialogIndirect` is called and
returns any failure HRESULT, the dialog falls back to `MessageBox`, as #6308
does. This is deliberate: a message dialog carries information or a decision
the user has to see, and failing to show it is worse than showing it unthemed;
the fallback still fires the right callbacks via `messageBoxButtonFor`.

What this WEP changes is how loudly that is reported. `E_INVALIDARG` can only
mean a malformed `TASKDIALOGCONFIG` — a Wails bug, not an environment problem —
so it is logged as an error asking for a report; other failures (e.g.
`E_OUTOFMEMORY`) stay a warning. Encoding bugs are meant to be caught before
that, by the layout test (see Test Plan).

## Non-Goals

- No public API change (`MessageDialog`, `AddButton`, `SetDefaultButton`,
  `SetCancelButton`, icons).
- Rendering arbitrary `Icon` bytes (needs `TDF_USE_HICON_MAIN`; separate).
- Custom/HTML dialogs; macOS and Linux dialogs.
- Changes to the v2 dialog code (#6250's area).

## Platform Considerations

- **Windows**: themed dialogs whenever comctl32 v6 is loaded, i.e. the app has
  the Common-Controls v6 manifest (debug builds now get it by default). Without it,
  behaviour is the current `MessageBox` path. Packed layout is 386/amd64/arm64
  safe (#6308's encoder).
- **macOS / Linux / mobile**: unaffected.

## Pros/Cons

### Pros

- Every Windows message dialog is themed and consistent with #6308's.
- Small: reuses #6308's binding, encoder, mapping and fallback.
- Same buttons, results and callbacks as `MessageBox` for standard dialogs.

### Cons

- Visible change for every existing app's dialogs (see below).
- The title appears twice: window caption and heading.

## Alternatives Considered

- **Leave #6308's gate in place.** Keeps two different looks in one app, and
  standard dialogs stay classic.
- **Title only as the window caption** (no main instruction). Closer to
  `MessageBox`, but loses the heading that makes task dialogs readable; the
  caption is small and often truncated.
- **Fall back only for specific HRESULTs.** Rejected above: any failure to show
  the dialog is worse than an unthemed one.
- **#5193** (custom-button callbacks via TaskDialog) is largely superseded by
  #6308 and does not cover standard dialogs. **#6250** (app icon via
  `MessageBoxIndirectW`) is v2's `MessageBox` path; v3 already has that
  (`MessageBoxWithIcon`), and the TaskDialog path here uses the same icon
  resource.

## Backwards Compatibility

No API change. Observable differences, where comctl32 v6 is loaded:

- Info/Warning/Error dialogs, and Questions without custom buttons, render as
  themed task dialogs with the title as a bold heading.
- Info/Warning/Error dialogs **with** custom buttons show those buttons (as
  Questions already do since #6308), where `MessageBox` showed only OK and
  matched callbacks by caption.
- Such a dialog without a cancel button no longer closes on Esc.
- Standard buttons, their localised captions and their results are unchanged.

Without comctl32 v6, or on fallback, behaviour is exactly today's.

## Security and Privacy

None. Standard comctl32 UI call; no new capability or data handling.

## Test Plan

Unit (platform-independent, run on every CI OS):

- `TestPlanTaskDialogMapsButtons`, `TestPlanTaskDialogNoCancelButton`,
  `TestMessageBoxButtonFor` — #6308's mapping, unchanged.
- `TestEncodeTaskDialogConfigLayout` — byte offsets and `cbSize` for 64-bit
  (160 bytes; amd64 and arm64 share it) and 386 (96 bytes), extended to
  `dwCommonButtons` and `pszMainInstruction`, plus the `TASKDIALOG_BUTTON`
  layout.
- `TestTaskDialogStandard` — buttons and icon per dialog type.
- `TestTaskDialogAllowsCancel` — the cancellation rule above, for each standard
  type and for custom buttons with and without a cancel button.
- `TestTaskDialogFailureIsBug` — only `E_INVALIDARG` is reported as an error.

Manual, Windows 10 and 11:

- Each type, with and without custom buttons, renders themed with the right
  icon, buttons and default.
- Esc / close box on a standard OK dialog (closes it), on a standard Yes/No
  question (ignored), on a custom dialog with a cancel button (fires its
  callback), and on one without (ignored).
- App-icon dialog (About) with and without the embedded resource.
- A binary without the v6 manifest shows `MessageBox` with working callbacks.

## Reference Implementation

[`9ad6f13066670174b4bdf27b496f8b856a581404`](https://github.com/juggler-ai/wails/commit/9ad6f13066670174b4bdf27b496f8b856a581404)
on the juggler-ai/wails branch
[`feat/windows-themed-taskdialog`](https://github.com/juggler-ai/wails/tree/feat/windows-themed-taskdialog):
one commit on top of current `master`, which includes #6308
(`dialogs_windows.go`, `dialogs_taskdialog.go`, `dialogs_taskdialog_test.go`).
It includes the error-level logging for `E_INVALIDARG` described above.

## Maintenance Plan

The added surface is two encoder fields and a four-case helper, both
unit-tested; the rest is #6308's code, maintained with the Windows dialogs.

## Conclusion

#6308 did the hard part. Routing the remaining message dialogs through it makes
every Windows dialog themed and consistent, keeps `MessageBox` as a safe
fallback, and leaves the public API untouched.
