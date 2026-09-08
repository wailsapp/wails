# macOS verification after native acceptance fixes

8 September 2026, `codex/hcl-build-system`. Follow-up to `platform-fixes.md` and [build-system issue #6021](https://github.com/wailsapp/wails/issues/6021).

The changes accompanying this record resolve the macOS failures found on base commit `906a9aa2776f800137c8a476028059582e1b33c8`. The native checks below ran against the resulting source tree and its freshly built CLI. The disposable applications used a Go module replacement and locally built frontend runtime from this checkout.

## Corrections

- Native signing passes `darwin.notarization.credential` to notarytool. A regression with distinct signing and notarization credentials failed before the change and passes afterward.
- Canonical explicit package-script commands, such as the migrated `frontend.dev = ["npm", "run", "dev"]`, receive the same Vite host/port/strict-port arguments as compiled defaults. Non-Vite commands and arbitrary custom argv retain their existing launch semantics.
- Failure panels preserve wrapped native diagnostics even when the error has a known exit code and no separately captured output. Long diagnostic lines wrap instead of losing their actionable suffix.
- A complete repeat of the frontend timeout confirmed that its error is already printed after frontend cleanup. The earlier incomplete capture was not evidence of a separate CLI error-printing defect; no top-level CLI change was required.

## Environment

macOS 26.4.1 (25E253), arm64; Xcode 26.6 (17F113); macOS SDK 26.5; Apple clang 21.0.0; Go 1.26.2; Node 20.20.2; npm 10.8.2. Signing used the available Developer ID Application identity and a validated notarytool keychain profile. Secrets were entered through the native secure prompt and were not placed in HCL or logs.

Signed builds used a non-synced local temporary directory. This host's synced Documents folder adds Finder metadata to generated app directories, which codesign rejects. No machine-wide sync or signing settings were changed.

## Native results

| Check | Result |
|---|---|
| arm64, amd64 and universal app/DMG matrix | Passed, including artifact receipts and zero-work unsigned warm reruns |
| Custom complete Info.plist | Preserved in the native artifacts |
| Universal standalone `.app` signing/notarization/stapling | Passed with only the documented notarization credential field |
| Standalone app validation | `stapler validate`, strict/deep `codesign`, and Gatekeeper execution assessment passed; source `Notarized Developer ID` |
| Universal DMG signing/notarization/stapling | Passed without duplicated signing credentials |
| DMG validation | Ticket validation, Gatekeeper primary-signature assessment and strict image/enclosed-app signature checks passed |
| Repeated signed/notarized DMG | Passed; unsigned packaging reused while signing and notarization reran |
| Signed standalone app and mounted DMG launch | Passed frontend-to-Go request/response round trips on arm64 |
| Intel build execution (community-supported) | Passed request/response round trip under Rosetta; native Intel hardware execution remains untested and is non-blocking |
| Existing migrated Vite app | Passed initial launch and frontend configuration restart with no manual host argument |
| Fresh migrate / activate / build / dev | Passed with generated HCL unchanged |
| Development isolation and cleanup | Production app file hashes unchanged; dev app under `.wails/dev/`; SIGINT exit 0; frontend/backend gone and listening port reusable |
| Invalid signing identity | Exit 1 with native `no identity found` diagnostic in the failure panel |
| Missing and nonexistent notarization profiles | Exit 1 with actionable credential diagnostics |
| Failed signing publication | Last successful signed artifact and receipt remained byte-identical after each negative check |

Accepted Apple submissions include standalone app `40a8ac70-c23d-4098-ac3b-6ce9b383c294`, DMG `0cbac68f-7ad3-4166-9ae2-f46c46181c5d`, and repeated DMG `72b32bb8-322d-4910-b67e-643734edf687`.

The launch probe calls Go from JavaScript, checks the returned nonce, and makes a second Go call confirming that response. Thus `HCL_BINDING_ROUNDTRIP` in the captured logs verifies both directions of the bridge.

## Automated verification

```sh
go test -race ./internal/wake/... ./internal/commands ./internal/report/... ./cmd/wails3
go vet ./internal/wake/... ./internal/commands ./internal/report/... ./cmd/wails3
GOOS=windows GOARCH=amd64 go test -c ./internal/commands
```

All passed. The focused credential-routing, migrated-command and failure-panel regressions failed on the original implementation and passed after the fixes. Windows test compilation was checked because the development launcher is shared; native Windows execution is not claimed.

CodeRabbit reviewed the implementation and tests with no findings. The installed CLI no longer accepts `--plain`; its supported equivalent, `coderabbit review --uncommitted`, was used after attempting the repository-prescribed command.

## Acceptance boundary

macOS HCL acceptance is complete for Apple Silicon. Per the project owner's decision on 8 September 2026, Intel macOS is community-supported: native Intel hardware confirmation is non-blocking. Intel and universal build targets remain available, with builds and Intel execution under Rosetta verified above; native Intel hardware execution is not claimed.

This record does not close Windows or iOS acceptance; in particular, the previously observed iOS icon assembly defect is outside these macOS corrections. No production release or upload was published.
