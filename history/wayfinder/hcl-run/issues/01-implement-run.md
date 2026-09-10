# Implement manifest-driven run defaults

Type: task
Status: resolved
Blocked by: none
Labels: needs-triage

## Question

Implement the [proposed reference](../../../../docs/src/content/docs/reference/wails-hcl-run.md)
for `wails3 run`, shared run settings and target overrides. Settle the supported-
platform declaration and mobile selection/lifecycle contracts before implementation.
Cover the reference acceptance criteria, regenerate schema documentation and
complete native platform checks before marking the feature available.

## Comments

2026-09-10: User requested documentation of the capability. Documentation is
complete; this ticket tracks implementation separately. Feature release placement
requires maintainer triage; no GA availability is claimed.

2026-09-10: Implementation requested, including every v3 example and a `go run .` fallback when no manifest exists.

## Answer

2026-09-10: Implemented manifest-driven run defaults, OS/architecture overrides,
mobile selection and the no-manifest `go run .` fallback. All 105 example entry
directories have valid local manifests and unique identities. Linux verification
completed with 100 launch smoke passes, five expected platform rejections and no
failures. Go tests, race checks, vet, Android APK assembly/signature checks and
CLI cross-compilation passed. See [verification](../verification.md) for commands,
evidence and the limits of these checks.

Native acceptance and interactive feature checks are explicitly retained in
[native acceptance](02-native-acceptance.md), following the user's offer to test
other platforms manually. This resolution does not claim GA readiness or completed
Windows/macOS/iOS/Android device execution.
