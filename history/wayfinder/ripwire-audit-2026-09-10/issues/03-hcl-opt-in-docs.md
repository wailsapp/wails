# Document the HCL experimental opt-in

Type: task
Status: resolved
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

## Answer

2026-09-10: The HCL tutorial and manual acceptance instructions now set WAILS_EXP_USE_WAKE before any experimental command, for Bash/Zsh and PowerShell. They explain presence-based enabling, unset behavior and an explicit manifest-validation acceptance check. A freshly built CLI validates the documented fixture with the variable set and follows legacy command routing when it is unset.
