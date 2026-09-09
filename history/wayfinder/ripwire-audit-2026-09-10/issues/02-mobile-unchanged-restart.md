# Preserve unchanged mobile development applications

Type: task
Status: open
Blocked by: none
Label: ready-for-agent
Priority: P2
Release: nightly patch; excluded from feature-only v3.1.0

## Question

manifest_dev_mobile.go line 119 makes BackendChanged always true. Non-artifact edits therefore stop and reinstall an unchanged Android/iOS app.

## Acceptance criteria

Compare installable artifact identity and retain live app state for cached/non-artifact rebuilds. Still replace the app after executable or packaging changes; add controller and native acceptance coverage.

## Comments

2026-09-10: Audited at `6c556f0b9`. See the [Ripwire audit](../../../ripwire-audit-2026-09-10/report.md) for evidence and limits.
