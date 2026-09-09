# Document the HCL experimental opt-in

Type: task
Status: open
Blocked by: none
Label: ready-for-agent
Priority: P2
Release: nightly patch; excluded from feature-only v3.1.0

## Question

The primary HCL tutorial claims manifest presence alone selects HCL. WAILS_EXP_USE_WAKE is now required; real config check without it prints help instead of validating.

## Acceptance criteria

Document shell-specific opt-in before all HCL examples and manual acceptance commands, and verify the instructions from an ordinary unset environment.

## Comments

2026-09-10: Audited at `6c556f0b9`. See the [Ripwire audit](../../../ripwire-audit-2026-09-10/report.md) for evidence and limits.
