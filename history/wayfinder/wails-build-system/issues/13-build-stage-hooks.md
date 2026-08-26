# Build and stage hook contract

Type: task
Status: resolved
Blocked by: none
Label: ready-for-agent

## Question

Design and implement a bounded build/stage hook capability for projects that
need more than declarative configuration without turning `wails.hcl` into a
second general-purpose Taskfile language.

The following is an illustration of the desired readability only. It is not a
syntax, naming, stage, settings, execution, caching, or compatibility contract:

```hcl
stage "dependencies.prepare" {
  settings {
    go_modules = "download"
    frontend   = "install-if-needed"
  }
}
```

The design must derive hook points from real current Taskfile customisations
and `wails3 tool` calls, define typed and external-script boundaries, and cover
scope, ordering, working directory, environment, inputs, outputs, cache policy,
failure behaviour, cancellation, Plan rendering, migration and security. After
maintainer review of the contract, implement it with schema, planner, executor,
migration, documentation, correctness, coverage and performance tests.

## Comments

### 2026-08-26 — Added to the resumed goal

Do not treat the illustrative `dependencies.prepare` block as a specification.
The first action is evidence gathering and interface design; no syntax should
be frozen until the execution and cache semantics are understood.

## Answer

Implemented as a bounded lifecycle contract rather than a general stage API.
`hook "<phase>"` accepts one project-owned script for `before_build`,
`after_build`, `before_package`, `after_package`, `before_sign`, or
`after_sign`. Project and target scopes are fixed; package and sign hooks form
barriers around the complete requested format set while independent format
nodes remain concurrent.

Scripts execute directly with no inline shell or argument interpolation. Paths
and working directories are project-contained after symlink resolution. Stable
`WAILS_*` values override inherited environment, non-zero exits retain output
and exit status, and cancellation terminates the complete process group.

Hooks run every time unless `cache = true` declares complete `inputs` and
`outputs`. The script bytes and executable mode are implicit inputs. Multiple
outputs share one bounded non-root directory, and an output root may not contain
the script or declared inputs or contain symbolic links or special files.
Always-run hooks make their downstream operations non-cacheable; cacheable
hooks contribute their declared artifact identity instead. Conservative
argument-free Taskfile
lifecycle scripts migrate to this shape; arbitrary commands remain blockers.

The implementation includes strict schema/writer support, typed Plan nodes,
artifact/terminal-barrier separation, Linux and Windows launch adapters, Plan
rendering, cache restoration, migration, documentation, and focused validation,
topology, environment, failure, cancellation, cache, and migration tests.
The full focused suite and race detector pass; Windows amd64 and macOS arm64
production packages cross-compile. The unchanged no-hook graph passes the
seven-sample Linux gate at 105.388ms median with 3.45% MAD, zero executed and
six cached Nodes, and stable artifact bytes.
