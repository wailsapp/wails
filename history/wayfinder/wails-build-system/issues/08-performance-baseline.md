# Current build performance baseline

Type: task
Status: resolved
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
- DEB package: 2.76–2.84s after rebuilding the CLI; 70–90ms warm no-op with all
  six Nodes cached. The final run included the coarse-timestamp cache safety
  window and packaging over an existing output. Broader five-sample and
  native-platform measurements remain.

## Answer

The reproducible baseline, environment, commands, semantic work matrix, and
measurement caveats are recorded in
[`performance-baseline.md`](../performance-baseline.md). A checked-in Go
benchmark package now measures complete CLI processes, strips terminal control
sequences, records wall/CPU/reported graph duration and ran/cached counts,
computes medians and ranges, reads and writes atomic JSON baselines, and applies
absolute plus relative median budgets. The CLI wrapper is
`v3/scripts/benchmark-manifest-build.go`; in-process Planner and Executor
benchmarks isolate orchestration overhead.

Five-sample results on the Ryzen 3950X Linux host establish:

- manifest badge no-op: 87.6ms median, 77.9–89.8ms, all four Nodes cached;
- larger dock no-op: 74.2ms median, 72.5–83.1ms, all four Nodes cached;
- Taskfile Wake no-op: 91.2ms median; external Task: 228.5ms median;
- DEB and AppImage no-op package: 92.0ms and 84.7ms median, respectively;
- method-body edit: 904.6ms, only compile ran;
- binding-shape edit: 2,816.2ms median, with generated-output identity safely
  skipping frontend work where different Go types produced identical
  TypeScript;
- frontend comment with identical bundle: 965.1ms, only frontend ran;
- frontend output change: 1,787.9ms, frontend and compile ran;
- missing binary restoration: 85.3ms, no handlers ran and restored bytes
  matched the removed Artifact;
- cache-cold dock build with warm tool caches: 3,184ms; first AppImage package:
  21,173ms including linuxdeploy/tool work;
- live Dev method-body rebuild: 870ms, compile only, followed by a clean
  backend replacement and child-process shutdown.

The three-target Planner median is 687µs/op; warm fake-handler Executor medians
are 2.26ms for one Target and 3.28ms for three Targets. Native multi-platform
tool timing is deliberately left to host acceptance rather than represented by
cross-compilation from Linux.

The checked baseline is
`v3/internal/wake/benchmark/testdata/badge-noop-linux-amd64.json` at 83.332ms.
The rebuilt CLI passes the sub-100ms absolute gate and the allowed 20% median
regression. This resolves the factual baseline; the permanent CI/release
workload and enforcement policy remain in the dependent performance-acceptance
ticket.
