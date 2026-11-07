# Wails Enhancement Proposal (WEP)

## Title

Native file drag-out

| Field | Value |
| --- | --- |
| WEP Number | |
| Status | Draft |
| Author | Philipp Mieden (@dreadl0ck) |
| Created | 2026-11-07 |
| Discussion | [#3158](https://github.com/wailsapp/wails/issues/3158) |
| Implementor | Philipp Mieden (@dreadl0ck) |
| Target | Wails v3 |

## Summary

Let an app drag existing local files into Finder, Explorer, a Linux file manager,
or another application. The receiver gets files, not text containing their paths.

## Motivation

AudioVault needs to drag audio samples into a DAW. Wails v3.1.0 handles incoming
file drops, but needs a local patch for outgoing native file drags. #3158 asks
for the same capability in an image-conversion app.

## Detailed Design

Add an opt-in window option and two Go methods:

```go
application.WebviewWindowOptions{EnableFileDragOut: true}

window.ArmFileDragOut(paths []string, imagePath string) error
window.DisarmFileDragOut()
```

The Go backend authorises the files and arms the next pointer drag. Paths must
be absolute and exist. The operation offers copies only; it never deletes the
originals. `imagePath` is an optional preview hint, with a standard icon fallback.

Completion dispatches `wails:file-drag-out-ended` to the page with
`detail.performed`: true for a completed drop, false for cancellation.
Disarming clears a pending drag; Escape cancels an active native session.
Ordinary clicks must keep working while armed.

## Non-Goals

Moving or deleting source files, virtual files generated during a drop, clipboard
changes, mobile drag-out, and browser/server mode.

## Platform Considerations

| Platform | Approach |
| --- | --- |
| macOS | AppKit `NSDraggingSession` with file URLs. Keep the existing mouse gesture intact. |
| Windows | OLE drag-and-drop with a Unicode `CF_HDROP` file list. Support files from different directories and volumes. |
| Linux | GTK3/GTK4 native drag sources offering `text/uri-list`. Preserve clicks and release all drag state on cancellation. |
| Mobile/server | Return an unsupported-operation error; do not silently succeed. |

The desktop API should work with the normal CGO backend and the experimental
purego backend. Backend differences must not change the file list or completion
result.

## Pros/Cons

One API covers a common desktop workflow. The cost is maintaining three native
implementations and testing their input handling, not just compilation.

## Alternatives Considered

HTML5 dragging alone did not reliably export multiple local files in the
WebKitGTK 2.52.6 prototype. Asking each app to maintain an AppKit/OLE/GTK patch
duplicates the same platform work.

## Backwards Compatibility

`EnableFileDragOut` defaults to false. Existing file-drop destinations and HTML5
drags keep their current behaviour when it is disabled.

## Security and Privacy

Arming is Go-only. The page must not be able to grant itself access to arbitrary
file paths. Apps should resolve selected record IDs to authorised files in their
own backend. A user drop gives the receiving application access to those files.

## Test Plan

Require native transfer tests before merging an implementation:

| Check | Required result |
| --- | --- |
| Multiple files, different folders, spaces and Unicode | Receiver gets the exact file list. |
| Escape, then another drag | Cancellation reports false; the next drag works. |
| Plain click while armed | The webview still receives the click. |
| Disarm, window close and disabled option | No stale drag source or input grab remains. |
| File validation | Reject empty, relative, missing and NUL-containing paths. |

The prototype has manual macOS confirmation in AudioVault. Native tests pass on
Windows amd64 and Ubuntu 24.04 with GTK3/GTK4, each using CGO and purego. They
cover ordinary clicks, multi-file transfer, Unicode paths, cancellation and
reuse. Windows arm64 is compile-checked only.
[CI results](https://github.com/dreadl0ck/wails/actions/runs/42426040522).

## Reference Implementation

[Prototype and native test probes](https://github.com/dreadl0ck/wails/tree/cff11ab6c22449c91fd4c013bba94ba3461a6de2).
This PR contains only the proposal; the prototype is separate.

## Maintenance Plan

Keep the Go API and runtime completion handler together. Run native file-list,
click, cancellation and reuse tests on Windows and both Linux GTK backends when
the implementation changes.

## Conclusion

Add a small, opt-in Go API for dragging real files out of a desktop window,
with the same result on macOS, Windows and Linux.
