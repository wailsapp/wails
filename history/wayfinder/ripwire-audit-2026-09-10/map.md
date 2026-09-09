# Ripwire audit follow-up

Audited `6c556f0b9` on 10 September 2026. [Report and evidence](../../ripwire-audit-2026-09-10/report.md). Open findings are recorded under `issues/`, ordered by severity. These are bug fixes for nightly patches, outside v3.1.0.

## Decisions so far

Ripwire and its 16 Codex usage skills are installed. Structural risk counts are kept separate from the three actionable review findings. All three findings have now been fixed and regression-tested on Linux.

- [Native compilation cache inputs](issues/01-cgo-cache-inputs.md): track all Go-supported native extensions in project and local-module sources.
- [Unchanged mobile backend](issues/02-mobile-unchanged-restart.md): compare installable content against the last successful launch; preserve retry and rollback behavior.
- [HCL opt-in instructions](issues/03-hcl-opt-in-docs.md): document the prerequisite in the tutorial and manual acceptance commands.

Native Android/iOS session acceptance remains a matching-host verification step in the [platform handoff](../wails-build-system/multi-platform-testing-handoff.md).
