# Windows HCL acceptance on win-node1

8 September 2026 — `codex/hcl-build-system`, following the Windows checks in `platform-fixes.md` (historical `hcl-ga-fixes` references interpreted as `hcl-build-system`). Related project tracker: [wailsapp/wails#6021](https://github.com/wailsapp/wails/issues/6021).

## Result

Windows x64 acceptance is confirmed on win-node1. Windows ARM64 EXE, NSIS and MSIX generation, metadata and cache behavior are confirmed from the x64 host; execution on native ARM64 hardware remains untested.

The branch started at `906a9aa2776f800137c8a476028059582e1b33c8`, with the preceding macOS fixes at `e8188f78d5b217321433a914e93d27b3e5e4deea`. The Windows changes were tested in the isolated checkout `C:\Users\leaan\hcl-acceptance-20260908\wails`, with a disposable app replacing its Wails Go module and frontend runtime with that checkout. SHA-256 comparison confirmed all changed source files matched the local worktree. The proof nonce contains `e8188f78d`; it identifies the fixture, not the final source revision.

Host: Windows 11 Pro x64, build 26200; Go 1.26.2; Node 24.15.0/npm 11.12.1; Git 2.54; MinGW GCC 15.2.0; NSIS; Windows SDK 10.0.26100.7705, using the SDK's x64 MakeAppx and SignTool; installed WebView2 152.0.4191.66. The SDK was installed for these checks.

## Fixes

- Accept stock CRLF Taskfiles during migration classification while preserving byte-exact source receipts and rejecting substantive edits.
- Reject Windows rooted and drive-relative paths where HCL requires project-relative paths.
- Link generated Windows resources through a real generated Go package: Go's object packer does not honor a virtual `.syso` overlay. Keep generated files out of the application source directory and reject collisions with the overlaid import filename.
- Preserve HCL product metadata when generating build assets. Include both numeric product version and string file version so Windows exposes product details correctly.
- Stage the sibling Windows assets and output directory required by NSIS templates.
- Generate correctly sized MSIX artwork from the application icon instead of transparent 1×1 placeholders.
- Allow native Windows Go builds across x64/ARM64 for the Wails Windows backend.
- Publish signed artifacts with usable extensions, such as `badge.signed.exe` and `badge.signed.msix`.
- Preserve native SignTool failures and redact passwords rather than hiding them behind a fallback signer error.
- Use fsnotify on Windows to avoid the previous watcher backend's race/checkptr crash.
- Launch npm's JavaScript entry point with Node on Windows, avoiding `cmd.exe` quoting failures under `Program Files`.
- Serialize Windows backend replacement to avoid killing the replacement's shared WebView2 browser when the old job closes. Keep the last working executable in memory and restore it if replacement startup fails; browser profile data is preserved. Microsoft documents the shared profile/process relationship in its [WebView2 process model](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/process-model).
- Always watch `wails.hcl`, including in migrated projects whose source watch list contains only `*.go`.
- Correct native test fixtures for Windows paths, line endings, file modes, process termination and temporary-directory cleanup. Unix shell fake-tool fixtures are explicitly limited to Unix; native Windows process, signing, resource and packaging paths were exercised separately.

## Verification

| Check | Evidence/result |
|---|---|
| Fresh stock Windows checkout migration and activation | 14 source files, zero diagnostics; legacy Taskfiles preserved |
| Clean x64 and ARM64 EXE/NSIS/MSIX builds | Passed with no certificate file or credential variable available |
| Warm unsigned builds | Every matrix rerun executed zero handlers; artifact receipts verified |
| PE resources | Machine 8664/AA64, file version 3.7.2, product name HCL Windows Acceptance |
| MSIX contents | x64/arm64 identity, version 3.7.2.0, correct 150×150, 44×44, 310×150, 620×300 and 50×50 artwork |
| Unsigned outputs | EXEs report NotSigned; MSIX has no AppxSignature.p7x |
| Protected-PFX signing | EXE and installer profiles built twice; signing reran while reusable stages were cached |
| Independent signature verification | SignTool `/pa /v` passed for signed EXE, NSIS and MSIX |
| Invalid credentials | Wrong and missing passwords rejected; no password in diagnostics; prior published artifact and receipt hashes preserved |
| Signed x64 EXE | Launched and completed frontend-to-Go request/response |
| NSIS | Silent installation, installed app Go binding, correct product details, silent removal |
| MSIX | Add-AppxPackage, activation through AppsFolder, Go binding, Remove-AppxPackage |
| Development from console | Initial and post-frontend-restart Go calls, console-interrupt exit 0, reusable port, unchanged production EXE |
| Nested Windows Job Object | Same development checks while the CLI inherited an outer job; tests Wails' nested job ownership |
| Backend replacement | Edited Go code, confirmed the new revision through the live WebView after the old backend stopped, and verified captured descendant handles exited |
| Failed Windows backend startup | Native regression restores and runs the saved executable after an invalid replacement |
| Repeated process/signing regressions | `platform-fixes.md` Windows regression selection passed 20 times |
| Native Windows race suite and vet | Passed for wake, commands, report and CLI packages; final race execution used `-count=1` |
| macOS regression suite and vet | Full suite passed; focused lifecycle race tests and vet passed again after the final Windows replacement fix |
| Minimum Go version | Go 1.25.0 actually compiled and linked x64 and ARM64 resource fixtures, including version metadata |
| CodeRabbit | Final review: no new findings. Its previous generated-import-path concern was contradicted by the actual Go 1.25 compiler regression |

The development harness makes recurring frontend-to-Go calls and requires a new confirmation after the recorded frontend process restart. It does not assume Vite will always reload an otherwise unchanged page. The nested-job check uses a native outer Job Object; a particular IDE's terminal UI was not operated. Console interruption used Ctrl+Break, handled by Go as an interrupt.

Signing and GUI checks ran in the logged-in `multica` desktop session through temporary scheduled tasks. The SSH session's user crypto context could not access the PFX private key, and session 0 could not initialize a WebView controller. The same HCL signing and application paths worked in the interactive desktop session.

## Reproduction

From `v3/`, with Git Bash on PATH for the existing legacy Taskfile tests:

```powershell
go test -race -count=1 -timeout=8m ./internal/wake/... ./internal/commands ./internal/report/... ./cmd/wails3/...
go vet ./internal/wake/... ./internal/commands ./internal/report/... ./cmd/wails3/...
go test ./internal/commands -run 'TestManifestProcess|TestHCLNativeSigningSettingsReachWindowsAndDarwinArtifacts|TestAndroidDeployUsesPlatformApplicationID' -count=20
go run scripts/verify-manifest-build-system.go -wails <built-cli.exe> -project <unsigned-fixture> -targets windows/amd64,windows/arm64
```

The final Windows race run disabled test-result caching with `-count=1`; an earlier cached run stalled in the Go driver and was stopped. Build-cache acceptance was tested separately and passed.

## Scope and cleanup

Signing used a disposable password-protected code-signing PFX and matching publisher, trusted only for this acceptance run. This verifies the HCL signing mechanism and installation flow, not a production publisher certificate, SmartScreen reputation, Microsoft Store submission or timestamp service.

The disposable certificate was removed from LocalMachine Root and TrustedPeople. The test PFX, private key and password files were deleted. Both installers were removed, all five temporary scheduled tasks were unregistered, no acceptance frontend/backend processes remained, and port 9346 was reusable. The SDK and isolated source/artifact directories remain available.

No changes were pushed or published. Intel macOS is recorded as community-supported and non-blocking per the project owner's instruction. This Windows work does not close the separate iOS icon/device acceptance work.
