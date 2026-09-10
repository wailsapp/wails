# Preserve development application arguments across restarts

Type: task
Status: resolved
Blocked by: none
Labels: needs-triage

## Question

Implement the runtime-configuration scenario from PR #6069 in Wake: development
argument defaults, platform overrides, CLI passthrough and an `--appargs`
compatibility spelling. Preserve arguments across rebuilds and rollback, restart
for argument-only edits without recompiling, and use platform-neutral resolution
with explicit adapter restrictions.

## Comments

2026-09-10: User requested implementation after reviewing the PR's scenario.
The implementation keeps development arguments separate from production run
arguments. Desktop and iOS launchers accept vectors; Android rejects nonempty
vectors before device effects. Legacy Taskfile command rewriting remains outside
this Wake-specific change.

## Answer

2026-09-10: Implemented and verified the shared argument contract, platform
adapters, literal CLI passthrough and compatibility spelling. Real Linux Wails
acceptance, normal/race tests, vet and Windows/macOS cross-compilation passed.
See [verification](../dev-arguments-verification.md) for commands, limitations and
native follow-up checks. Production run defaults remain separate.
