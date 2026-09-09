# Fixes for the three Ripwire audit findings

Implemented on `codex/hcl-build-system` after the audit at `6c556f0b9`.

- Native-source inputs now use one complete extension list for root compilation,
  bindings and local-module snapshots. Tests cover 19 native extensions in both
  project and replacement-module sources.
- Mobile development compares the final installable artifact with the version
  that last reached readiness. Uncached iOS bundles are fingerprinted directly;
  failed replacements leave the previous identity active for rollback/retry.
- The HCL tutorial and manual platform handoff now explain the required
  `WAILS_EXP_USE_WAKE` opt-in for Bash/Zsh and PowerShell, including unset behavior.

## Verification on Linux

Passed from `v3/`:

```sh
go test ./internal/dev/... ./internal/devruntime/... ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/features ./internal/report/...
go test -race ./internal/dev/... ./internal/devruntime/... ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/features ./internal/report/...
go vet ./internal/dev/... ./internal/devruntime/... ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/features ./internal/report/...
go build -o /tmp/wails3-three-fixes ./cmd/wails3
```

The native snapshot tests failed before the fix and passed afterward. The real
CGo executable test demonstrates warm-cache reuse followed by recompilation
from output `1` to `2` after editing only `helper.cxx`, both in the project and
in a local replacement module. A second unchanged build hits the cache again.
The original `reproduce-cxx.py` audit proof also passes now.

Mobile regression tests cover unchanged APK/app contents after a debounce edit,
packaging-only changes, failed readiness and retries, rollback, older restored
artifacts, reusable APK digests and missing installable outputs. The focused
native and mobile regression tests were also run under the race detector.

A freshly built CLI reported the disposable fixture's manifest as valid with
`WAILS_EXP_USE_WAKE=1`; without the variable it followed legacy command routing
and did not validate HCL. Ripwire's focused panel scanned the modified planner
and mobile adapter successfully; existing complexity findings remain outside
these correctness fixes.

Native Android/iOS device execution was not repeated. The updated
[manual handoff](../wayfinder/wails-build-system/multi-platform-testing-handoff.md)
contains app-state preservation, changed-package deployment, retry and signed
iOS acceptance steps. The three implementation tickets are resolved in the
[audit map](../wayfinder/ripwire-audit-2026-09-10/map.md).

CodeRabbit completed the pre-commit review with no code findings. Its sole
comment questioned the resolution date; 10 September is the session date in
Australia/Sydney. The installed CLI rejects `--plain`, so the supported
`coderabbit review --uncommitted` equivalent was used after attempting it.
