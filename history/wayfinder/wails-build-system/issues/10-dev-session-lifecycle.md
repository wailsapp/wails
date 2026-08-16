# Dev Session invalidation and process lifecycle

Type: grilling
Status: resolved
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

## Answer

`wails3 dev` is implemented as one long-lived Dev Session around ordinary
finite production Plans. Watch bursts are debounced into generations; starting
a newer generation cancels stale execution while preserving `context.Canceled`
through the executor. Pulse clears that generation without a failed-build
summary. Shutdown cancels and waits for in-flight work before terminating and
reaping the frontend, backend, and their complete process groups.

The frontend and backend are readiness-gated persistent processes rather than
Nodes. A compile cache hit/restoration leaves the current backend running;
executed compile or `after_build` work stages a replacement, verifies startup,
then swaps it. Failed Plans and replacement processes leave the healthy app
running. An already-rendered Wake error is not repeated below its failure
panel; Dev adds only the useful current-app status.

Manifest changes are transactional session updates. Frontend command,
directory, manager, and port changes stage a listening replacement frontend
before starting the backend that consumes it. Port changes can keep the old
pair healthy until the new pair is ready. Candidate
`FRONTEND_DEVSERVER_URL`/`WAILS_VITE_PORT` values are isolated to the build and
child processes, and both server and backend use an explicit IPv4 loopback
address to prevent `localhost` IPv4/IPv6 disagreement.

Watch sets honor root and nested `.gitignore` files, reload ignore changes
without building, register newly created directories, scan them for inputs
created before registration, and retain the old set if replacement fails.
Watch/exclude/debounce policy is re-resolved from the Manifest without turning
persistent processes into graph Nodes.

Unit, race, vet, and Windows cross-compile gates pass across report, Wake,
pipeline, commands, and CLI packages. A live Linux badge session verified a
6ms all-cached event with no backend restart, Go-only rebuild/restart,
transactional 9245→9246 frontend/backend replacement, immediate `.gitignore`
reload, framed failure while the old app remained healthy, and clean Ctrl+C
process-tree shutdown.
