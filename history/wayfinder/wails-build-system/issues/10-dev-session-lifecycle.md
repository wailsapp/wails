# Dev Session invalidation and process lifecycle

Type: grilling
Status: open
Blocked by: 03

## Question

How does the long-lived Dev Session map file changes to affected Nodes,
coordinate frontend HMR with backend rebuilds, coalesce bursts, cancel stale
Plans, restart or preserve the application process, report readiness and
failures, and recover when watched files or configuration change? Define the
fast path and lifecycle contract without turning persistent processes into DAG
Nodes.

## Comments

### 2026-08-16 — Accepted implementation default

- Dev is a long-lived Session owning one frontend process, filesystem watcher,
  current backend process, and generation counter. Each coalesced change burst
  requests a finite incremental Plan from the same production Planner.
- Frontend-source changes are left to the frontend server/HMR; Go, Manifest,
  binding, embed, or generated-input changes invalidate the corresponding
  Nodes. Binding-byte changes notify the frontend path; byte-identical output
  does not restart it.
- A newer generation cancels stale planning/execution. The currently healthy
  app stays alive until a replacement binary builds successfully, then is
  terminated gracefully and atomically replaced. A failed rebuild leaves it
  running and reports the failure.
- Manifest changes re-resolve the Session. Changes to frontend command, port,
  target, or watch policy restart the affected persistent process; cache and
  CLI changes do not turn those processes into graph Nodes.
- Readiness is explicit: frontend listening, backend built, backend started.
  Shutdown cancels work, terminates child process groups, flushes reporting,
  and leaves no detached children.

### 2026-08-16 — MVP implementation evidence

The manifest Dev Session coalesces watcher bursts and assigns each asynchronous
rebuild a generation. A newer burst cancels the previous context; serialized
execution prevents reporter/environment overlap. Successful manifest reloads
replace the backend only after a new process starts, re-register watch policy,
and restart the frontend when its manager, command, directory, or effective
port changes. Failed builds leave the current backend running.
