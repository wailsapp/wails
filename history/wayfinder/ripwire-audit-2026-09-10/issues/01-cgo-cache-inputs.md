# Include CGo sources in compilation cache inputs

Type: task
Status: open
Blocked by: none
Label: ready-for-agent
Priority: P1
Release: nightly patch; excluded from feature-only v3.1.0

## Question

planner.go lines 316 and 834 omit .cxx and other valid CGo source types. The isolated compile-snapshot test fails after a .cxx-only edit; real Go builds consume the file.

## Acceptance criteria

Cover every compiler input type in root and local-module snapshots; prove native-code-only edits invalidate a warm build and update executable behavior.

## Comments

2026-09-10: Audited at `6c556f0b9`. See the [Ripwire audit](../../../ripwire-audit-2026-09-10/report.md) for evidence and limits.
