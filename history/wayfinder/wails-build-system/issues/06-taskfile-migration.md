# Taskfile AST migration and unsupported customization policy

Type: grilling
Status: open
Blocked by: 01

## Question

What generated Taskfile patterns and user modifications can be translated into
manifest fields, profiles, hooks, or template references? What must be reported
for manual migration? Define dry-run output, machine-readable diagnostics,
write/remove behavior, behavior when `wails.toml` already exists, and the
preconditions for removing legacy files when unsupported Taskfile logic
remains. Define a machine-readable migration completion state, any explicit
cutover marker, and provenance-aware removal that cannot delete modified,
user-owned, or incompletely represented files. Define the safety rule that
inline shell blocks are never silently converted into generated scripts.

## Comments

### 2026-08-16 — Accepted implementation default

- `wails3 migrate` always analyses first and produces a human summary plus a
  stable JSON report under `.wails/migration-report.json`; `--json` writes the
  report to stdout and `--dry-run` performs no project writes.
- Migration recognizes shipped Taskfiles/config by structure, not only by
  byte equality. Known variables and task differences map to typed Manifest
  fields; script-file commands may map to hooks; inline shell and arbitrary
  task logic are manual diagnostics and are never copied into generated files.
- A completed migration produces the same Manifest shape as a native project:
  no migration metadata. An incomplete migration records only
  `[wake.migration] complete = false` so legacy execution remains active while
  the Manifest is inspected and edited. Completion provenance, source digests,
  and diagnostics are internal workflow state and live exclusively in
  `.wails/migration-report.json`.
- Existing `wails.toml` is never overwritten unless it already carries the
  same migration provenance and the update is a safe merge. Conflicts stop
  before writing.
- `--remove-old-files` removes only files whose current digests match the
  analyzed, fully represented sources in the migration report. Any unsupported
  or subsequently modified file remains, and removal failure does not
  invalidate the Manifest.

### 2026-08-16 — MVP implementation evidence

The AST migration recognizes historical stock files, translates project
metadata/package manager/associations/protocols, and now maps a conservatively
recognized single project-relative script task (`before-build`, `after-build`,
package/sign equivalents) to the corresponding typed hook. Commands with
arguments, interpolation, or shell operators remain manual diagnostics.
