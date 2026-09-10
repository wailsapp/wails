# HCL run configuration

## Destination

Provide local run defaults and platform overrides through `wails.hcl`, with a
`go run .` fallback when no manifest is discoverable.

## Decisions so far

[Implementation](issues/01-implement-run.md) is complete and tested on Linux.
Run builds use the production pipeline; runtime arguments and environment remain
outside compile cache identity. OS-wide and architecture-specific settings share
one manifest. Android deploys an APK; iOS uses the native Apple launch tools.

[Verification](verification.md) records the 105-example Linux sweep, automated
regressions, Android APK checks and platform plans. [Review](review.md) records
findings and fixes. Native acceptance remains separate from implementation and
must be completed before claiming release readiness on the remaining platforms.

[Development application arguments](issues/03-dev-arguments.md) now preserve
runtime configuration through rebuilds and rollback, with independent dev
defaults and platform overrides. [Verification](dev-arguments-verification.md)
records real Linux Wails acceptance and the remaining native checks.
