# CLI compatibility and migration cutover

Type: grilling
Status: open
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
- A project containing both systems uses the Manifest when there is no pending
  migration marker. `[wake.migration].complete = false` leaves build verbs on
  the legacy path and prints a concise migration status hint. Older preview
  manifests with `complete = true` remain compatible.
- `config check`, `config show`, `eject`, and `migrate` are Manifest-aware
  commands and never dispatch through Task. `eject` requires a Manifest;
  `migrate` requires legacy inputs.
- Legacy Taskfile execution remains a compatibility adapter for projects with
  no active Manifest. `wails3 task` is hidden/deprecated during the migration
  cycle and is not generated or documented for new projects.
- Unknown positional `KEY=value` build arguments are rejected in Manifest
  projects with a mapping hint; supported target/profile/options use explicit
  flags. Legacy projects retain their current passthrough behavior.
