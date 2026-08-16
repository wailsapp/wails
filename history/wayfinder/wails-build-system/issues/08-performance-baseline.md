# Current build performance baseline

Type: task
Status: open
Blocked by: none

## Question

Measure reproducible cold, warm no-op, frontend-only, binding-relevant Go,
ordinary Go, multi-Target, packaging, and development-restart timings across
representative Wails projects using the current Task and experimental Wake
paths. Record wall time, executed work, critical path, CPU time, and variance so
the later performance acceptance decision is grounded in facts rather than an
arbitrary percentage.

## Comments

### 2026-08-16 — Accepted implementation default

- Use the existing badge project and at least one larger real application.
  Capture cache-cold, tool-cache-warm cold, no-op, ordinary Go edit,
  binding-shape edit, frontend-only edit, missing-output restore, package, and
  Dev restart scenarios.
- Record median and range over five warm samples, wall time, Step durations,
  cache outcomes, and critical-path work. Existing Task CLI, Taskfile Wake, and
  Manifest Wake run from equivalent disposable fixtures.
- The already captured MVP figures seed this baseline; broader platform and
  package measurements are implementation verification rather than another
  upfront decision gate.

### 2026-08-16 — MVP evidence

- Disposable badge: 2.88s cold build with warm tool caches; 70–80ms warm
  no-op, all four Nodes cached.
- Service method-body edit: 0.91s wall; only Go compile executed while binding
  generation and frontend work remained cached.
- DEB package: 2.80–2.81s after rebuilding the CLI; 70–90ms warm no-op with all
  six Nodes cached. The final run included the coarse-timestamp cache safety
  window and packaging over an existing output. Broader five-sample and
  native-platform measurements remain.
