# Development and migration reproduction results

Reviewed commit: 3452def9d. Linux host. These are intentionally failing expected-contract regression tests; they were run temporarily in the commands package and then moved out of the active Go source tree.

Copy dev-migration-regressions_test.go to v3/internal/commands/review_hcl_regression_test.go, run from v3:

    go test ./internal/commands -run '^TestReview' -count=1

Remove only that copied file afterward. The test file requires a Unix host; its process-cleanup case kills the child it creates after observing the bug.

Observed failures:

- TestReviewDevUsesDiscoveredRoot: frontend receives the nested working directory, rather than the discovered parent manifest root.
- TestReviewMigrationDoesNotLoseTaskfileEnvironment: migration says complete, diagnostics are empty, and CGO_CFLAGS is missing from generated HCL.
- TestReviewMigrationDoesNotLosePort: migration says complete; internal Port=9357 but encoded HCL drops it.
- TestReviewFrontendDevReceivesEnvironment: real frontend child sees an empty variable instead of configured-value.
- TestReviewStopCleansChildrenAfterParentExit: child remains alive after stop on its exited parent.

Command result: FAIL, package github.com/wailsapp/wails/v3/internal/commands, 0.102s.

The original scoped suites and race/vet gates passed before introducing these temporary reproductions. Their success does not establish correctness of the uncovered cases.
