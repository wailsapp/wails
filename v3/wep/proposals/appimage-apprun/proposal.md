# Wails Enhancement Proposal (WEP)

## AppImage: use a local AppRun

**WEP Number**:  
**Status**: Draft  
**Author**: Damon Blais (@Albinogeek)  
**Created**: 2026-10-09  
**Discussion**:  
**Implementor**: Damon Blais (@Albinogeek); the implementation is written and on branch `pr/appimage-apprun`  
**Target**: Wails v3

## Summary

Add `-apprun <path>` (env `WAILS_APPIMAGE_APPRUN`) to `wails3 generate appimage`. The file is copied into the AppDir as `AppRun` instead of downloading AppImageKit's `AppRun-<arch>` from GitHub at build time. The flag wins over the environment variable; with neither set, the download is unchanged.

## Motivation

`generate appimage` always downloads `AppRun-<arch>` from the AppImageKit `continuous` release. A build with no network at that point (distro packaging, hermetic or locked-down CI) cannot produce an AppImage, and an app that needs its own AppRun (environment setup, a working directory, a wrapper) has no way to supply one short of patching Wails. The `continuous` tag also moves, so two builds of the same commit can bundle different AppRun binaries.

## Detailed Design

A new field on `GenerateAppImageOptions`:

```go
AppRun string `description:"Path to an AppRun file to use instead of downloading one (env WAILS_APPIMAGE_APPRUN)"`
```

Resolution is a small helper shared with the `-linuxdeploy` proposal: the flag value, else the environment variable, made absolute; empty means download. When a path is supplied the download of `AppRun-<arch>` is skipped and the file is copied into the AppDir as `AppRun` and made executable, exactly where the downloaded file would have gone. Nothing else in the command changes.

Documentation: a row in the `generate appimage` flag table of the CLI guide and an `Added` entry in the unreleased changelog.

## Non-Goals

Does not change how `linuxdeploy` is obtained (see the `-linuxdeploy` WEP), does not verify or sign the supplied file, and does not add AppRun templating.

## Platform Considerations

Linux only: `wails3 generate appimage` produces Linux AppImages and is not available on Windows or macOS, so there is no behaviour to define there. The option works on every architecture the command already supports.

## Pros/Cons

Pros: offline and reproducible builds; custom entry points; no behaviour change by default. Cons: one more option to document and keep working; the caller owns the file's correctness.

## Alternatives Considered

Vendoring AppRun into Wails (adds a binary to the repository and a refresh burden); pinning the download to a release tag (helps reproducibility but not offline builds or custom AppRun); a template hook (more machinery than the need justifies).

## Backwards Compatibility

Additive. With the option unset the command behaves exactly as before.

## Security and Privacy

The supplied file becomes the AppImage entry point and runs with the user's privileges, so the caller is trusting a file of their choosing instead of one fetched over HTTPS from GitHub. That is the intent of the option: it removes a network fetch of an unpinned binary from the build.

## Test Plan

A unit test of the shared path helper: flag over environment over empty, and the result is absolute. Manual check: build an AppImage on a host with GitHub blocked, supplying a local AppRun, and confirm the AppDir contains that file.

## Reference Implementation

Branch `pr/appimage-apprun` in the author's fork; it will be linked from the PR once the WEP is accepted.

## Maintenance Plan

The author maintains the option. It is a few lines in `v3/internal/commands/appimage.go` and follows the existing `generate appimage` structure, so it needs no separate upkeep.

## Conclusion

A small, additive option for `wails3 generate appimage` (-apprun <path>) that makes AppImage builds possible and repeatable in environments the command currently cannot serve.
