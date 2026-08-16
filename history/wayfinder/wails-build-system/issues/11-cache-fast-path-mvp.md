# Manifest-to-binary cache fast-path MVP

Type: prototype
Status: resolved
Blocked by: 01, 02, 03

## Question

Does a narrow host-desktop vertical slice of the manifest-driven Pipeline prove
that typed planning and content-addressed caching can make Wails no-op and
incremental builds materially faster without exposing build-system machinery?
Build a reviewable prototype from a minimal `wails.toml` through
`GenerateBindings`, `BuildFrontend`, and `CompileApplication`, including the
filesystem Snapshot fast path, project Action Index, local Artifact Store,
Artifact-digest propagation, and existing Wake reporting. Capture cold, warm
no-op, implementation-only Go edit, binding-shape edit, and frontend-only edit
timings. Stop for hands-on review before broadening Targets, package managers,
packaging, signing, hooks, or migration.

## Comments

### 2026-08-16 — Executable prototype ready for review

- Prototype branch: `codex/wake-cache-mvp` (local only; never pushed).
- Entry point: `WAILS_WAKE_MVP=1 wails3 build` selects the prototype Pipeline
  when a minimal `wails.toml` exists; the explicit opt-in prevents throwaway
  code from changing existing Taskfile project behavior.
- Prototype code:
  [`v3/internal/wake/mvpprototype/`](../../../../v3/internal/wake/mvpprototype/)
  and
  [`v3/internal/commands/wake_mvp_prototype.go`](../../../../v3/internal/commands/wake_mvp_prototype.go).
- Vertical slice: npm Receipt, in-process TypeScript interface bindings,
  frontend production build, host Go compile, BLAKE3 Snapshots and Action Keys,
  project Action Index, machine-local Artifact bytes, missing-output restore,
  and Pulse reporting.
- An isolated cache-cold run with a cold Go build cache took **58.85s**:
  dependency installation 1.7s, bindings/package analysis 25.2s, frontend
  0.886s, and compilation 30.9s. The external Go/type-analysis work dominates;
  the planner and cache are not on the cold critical path.
- The immediate no-op run took **0.07s wall time** and **6ms reported build
  time**, with all four Nodes cached. The equivalent current Taskfile Wake
  no-op on the same prepared fixture took **0.12s wall time** and **43ms
  reported build time** with six cached tasks.
- With underlying tool caches warm, an implementation-only Go edit took 1.93s
  (bindings analysis 0.933s, compile 0.859s, frontend cached); a binding-shape
  edit took 2.84s (bindings 0.965s, frontend 0.873s, compile 0.854s); and a
  frontend-only edit took 1.86s (frontend 0.908s, compile 0.874s).
- Removing the produced binary restored it from the Artifact Store in the next
  all-cached 0.07s invocation rather than recompiling it.
- Failed bindings and compile attempts created no successful Action record.
  The prototype also confirmed that a build reports required `go mod tidy`
  remediation without running it automatically.
- The clearest remaining speed opportunity is the approximately one-second
  binding analysis paid by implementation-only Go edits. Artifact-digest
  propagation correctly prevents that work from cascading into the frontend,
  but a persistent semantic Binding Model index would be needed to avoid the
  analysis itself.
- Awaiting hands-on review before this prototype ticket is resolved or its
  ideas are lifted into production Wake code.

## Answer

The prototype is accepted as proof of the cache fast path, not as the complete
build-system replacement. It proved that a minimal Manifest can drive typed
work without Taskfiles, that downstream Nodes can consume Artifact content
digests, and that a correct no-op build can complete in roughly 60–70ms while
missing outputs are restored instead of rebuilt.

The production implementation will keep these contracts and replace the
prototype package. The remaining system—resolved configuration, the complete
Planner and scheduler, generated platform assets, hooks, Dev Session,
migration, and CLI cutover—must be reviewed together as one solution. The
largest measured incremental-build opportunity is persistent binding/package
analysis, because an implementation-only Go edit still pays about one second
for analysis even though byte-identical bindings correctly stop downstream
frontend invalidation.
