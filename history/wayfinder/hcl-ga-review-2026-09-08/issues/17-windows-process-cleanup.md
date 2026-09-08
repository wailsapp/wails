# Clean Windows descendants after their wrapper exits

Type: task
Status: open
Blocked by: none
Label: ready-for-agent
Priority: P2

## Question

The Linux fix for [development process cleanup](15.md) cleans Unix process groups after parent exit. The Windows cleanup callback remains a no-op, and taskkill cannot reliably find the tree after its parent has gone away. A frontend wrapper may leave a server and port alive.

## Acceptance criteria

Own descendant lifetime independently of the parent PID, for example with an appropriately assigned Windows Job Object. Add native Windows regressions for early parent exit, interrupt-ignoring descendants, cancellation and port reuse. Preserve existing ordinary shutdown behavior.

## Comments

2026-09-08: Retained from CodeRabbit review. Linux process tests and Windows cross-compilation do not verify Windows descendant cleanup. Requires native Windows validation before closure.
