# Wails Enhancement Proposal (WEP)

## AppImage: prepare the AppDir without packing

**WEP Number**:  
**Status**: Draft  
**Author**: Damon Blais (@Albinogeek)  
**Created**: 2026-10-09  
**Discussion**:  
**Implementor**: Damon Blais (@Albinogeek); the implementation is written and on branch `pr/appimage-nopack`  
**Target**: Wails v3

## Summary

Add `-nopack` to `wails3 generate appimage`. linuxdeploy still bundles the AppDir, but `--output appimage` is omitted; the command logs the AppDir path and returns without producing an AppImage.

## Motivation

Today the command prepares the AppDir and packs it in one step. Callers that need to inspect or post-process the AppDir first (patch libraries, add files, sign, or pack with their own `appimagetool` or a different runtime) have to reproduce the whole preparation themselves or race the command.

## Detailed Design

A new boolean on `GenerateAppImageOptions`:

```go
NoPack bool `description:"Prepare the AppDir but do not pack it into an AppImage"`
```

The linuxdeploy command line is built without `--output appimage` when the flag is set. After linuxdeploy finishes, the command logs `AppDir prepared: <path>` and returns before the step that moves the AppImage into the output directory. Without the flag the command line and flow are unchanged.

Documentation: a row in the `generate appimage` flag table of the CLI guide and an `Added` entry in the unreleased changelog.

## Non-Goals

Does not pack the AppDir with another tool, does not choose where the AppDir lives, and does not change the contents of the AppDir.

## Platform Considerations

Linux only: `wails3 generate appimage` produces Linux AppImages and is not available on Windows or macOS, so there is no behaviour to define there. The option works on every architecture the command already supports.

## Pros/Cons

Pros: lets callers own the packing step; small and isolated. Cons: one more flag; the command can now succeed without producing the artifact its name suggests.

## Alternatives Considered

A separate `prepare` subcommand (more surface for one flag's worth of behaviour); documenting the AppDir layout so callers rebuild it themselves (duplicates logic that drifts); a hook script run between bundling and packing (more machinery, and a shell-execution surface).

## Backwards Compatibility

Additive. With the flag unset the command behaves exactly as before.

## Security and Privacy

No new capability: the flag only skips the final packing step. The AppDir is left on disk in the build directory, as it already is during a normal run.

## Test Plan

A unit test that the command line gains `--output appimage` only without the flag. Manual check: run with `-nopack` and confirm the AppDir exists, no AppImage is written to the output directory, and the AppDir path is logged.

## Reference Implementation

Branch `pr/appimage-nopack` in the author's fork; it will be linked from the PR once the WEP is accepted.

## Maintenance Plan

The author maintains the option. It is a few lines in `v3/internal/commands/appimage.go` and follows the existing `generate appimage` structure, so it needs no separate upkeep.

## Conclusion

A small, additive option for `wails3 generate appimage` (-nopack) that makes AppImage builds possible and repeatable in environments the command currently cannot serve.
