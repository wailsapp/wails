# CLI compatibility and migration cutover

Type: grilling
Status: resolved
Blocked by: 01

## Question

What should `wails3 init`, `build`, `dev`, `package`, `sign`, `eject`,
`migrate`, `config check`, and `config show` do in manifest projects and
pre-migration Taskfile projects? Define detection, error messages, release
staging, precedence and conflict behavior when `wails.toml` and Taskfiles
coexist, how migration completion or an explicit cutover marker affects that
precedence, the exact legacy fallback condition, and the point at which
generated Taskfiles stop being supported.

## Comments

### 2026-08-16 — Accepted implementation default

- New projects contain `wails.toml` and no Taskfiles. Manifest-only projects
  use the built-in Pipeline for build, package, sign, and dev.
- A project containing both systems uses the hidden migration report as the
  cutover decision. An incomplete report leaves build verbs on the legacy path
  and prints a concise migration status hint; a complete report uses the
  Manifest. Coexistence without a report is rejected as ambiguous.
- `config check`, `config show`, `eject`, and `migrate` are Manifest-aware
  commands and never dispatch through Task. `eject` requires a Manifest;
  `migrate` requires legacy inputs.
- Legacy Taskfile execution remains a compatibility adapter for projects with
  no active Manifest. `wails3 task` is hidden/deprecated during the migration
  cycle and is not generated or documented for new projects.
- Unknown positional `KEY=value` build arguments are rejected in Manifest
  projects with a mapping hint; supported target/profile/options use explicit
  flags. Legacy projects retain their current passthrough behavior.

## Answer

Manifest activation is resolved in one command-routing seam. Manifest-only
projects use the built-in pipeline; Taskfile-only projects use legacy
execution; an incomplete private migration report deliberately falls back; and
a complete report activates the Manifest even during cleanup. Both systems
without a report fail with an actionable ambiguity error, while an incomplete
report without a Taskfile fails as an invalid migration. Config/eject commands
remain Manifest-native, and completed migrations normally remove the source of
future ambiguity by retiring the Taskfiles. `migrate --complete` is the only
existing-manifest migration path and never rewrites that manifest.
