# Wails Enhancement Proposal (WEP)

## AppImage: use a local linuxdeploy

**WEP Number**:  
**Status**: Draft  
**Author**: Damon Blais (@Albinogeek)  
**Created**: 2026-10-09  
**Discussion**:  
**Implementor**: Damon Blais (@Albinogeek); the implementation is written and on branch `pr/appimage-linuxdeploy`  
**Target**: Wails v3

## Summary

Add `-linuxdeploy <path>` (env `WAILS_APPIMAGE_LINUXDEPLOY`) to `wails3 generate appimage`. The given executable is run instead of downloading `linuxdeploy-<arch>.AppImage` into the build directory. The flag wins over the environment variable; with neither set, the download is unchanged.

## Motivation

`generate appimage` downloads `linuxdeploy-<arch>.AppImage` from the linuxdeploy `continuous` release the first time it runs. Offline or hermetic builds cannot do that, and because `continuous` moves, a build cannot pin a known-good linuxdeploy: a regression upstream breaks Wails builds with no change on the Wails side.

## Detailed Design

A new field on `GenerateAppImageOptions`:

```go
LinuxDeploy string `description:"Path to a linuxdeploy executable to use instead of downloading one (env WAILS_APPIMAGE_LINUXDEPLOY)"`
```

A helper resolves the value: the flag, else the environment variable, made absolute; empty means download. When a path is supplied the existing download-if-missing step is skipped, the file is marked executable, and it is the executable the command runs (still with `--appimage-extract-and-run`, so a linuxdeploy AppImage works without FUSE). The helper is the same one used by the `-apprun` proposal.

Documentation: a row in the `generate appimage` flag table of the CLI guide and an `Added` entry in the unreleased changelog.

## Non-Goals

Does not change how AppRun is obtained (see the `-apprun` WEP), does not manage linuxdeploy plugins, and does not verify the supplied binary.

## Platform Considerations

Linux only: `wails3 generate appimage` produces Linux AppImages and is not available on Windows or macOS, so there is no behaviour to define there. The option works on every architecture the command already supports.

## Pros/Cons

Pros: offline and pinned builds; no behaviour change by default. Cons: one more option to document; the caller owns compatibility between their linuxdeploy and Wails's plugin flags.

## Alternatives Considered

Pinning the download URL to a release (helps reproducibility, not offline builds); looking linuxdeploy up on PATH (implicit, and a distro linuxdeploy may be a different program); a config-file setting (no other `generate` command reads one).

## Backwards Compatibility

Additive. With the option unset the command behaves exactly as before.

## Security and Privacy

The supplied executable is run with the user's privileges. The option replaces an unpinned download with a file the caller chose, which is the point of it.

## Test Plan

A unit test of the shared path helper: flag over environment over empty, result absolute. Manual check: run the command with a locally downloaded linuxdeploy and the network blocked, and confirm no download is attempted and the AppImage is produced.

## Reference Implementation

Branch `pr/appimage-linuxdeploy` in the author's fork; it will be linked from the PR once the WEP is accepted.

## Maintenance Plan

The author maintains the option. It is a few lines in `v3/internal/commands/appimage.go` and follows the existing `generate appimage` structure, so it needs no separate upkeep.

## Conclusion

A small, additive option for `wails3 generate appimage` (-linuxdeploy <path>) that makes AppImage builds possible and repeatable in environments the command currently cannot serve.
