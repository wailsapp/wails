# Build and stage hook contract

Type: task
Status: open
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

Pending.
