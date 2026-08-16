# Wake's typed build graph and execution boundary

Type: grilling
Status: resolved
Blocked by: none

## Question

Which Wails build concepts become first-class typed nodes, what inputs and
outputs does each node own, and which parts of the current Wake Taskfile parser,
resolver, cache, executor, and reporter can be reused versus replaced? Define
the graph contract before implementation begins.

## Comments

### 2026-08-15 — Grilling round 1

- The canonical graph terms are Pipeline, Plan, Node, Artifact, and Step. A
  Step is only the reporter's presentation of a running Node; Task and Command
  are not part of the new engine's domain language.
- Planning is read-only and completes before execution. The Planner resolves
  the Manifest, inspects the project and toolchains, expands Targets and
  formats, and produces an immutable Plan. Execution never adds Nodes or
  dependencies dynamically.
- One invocation produces one multi-Target Plan. Target-independent Nodes such
  as binding generation and frontend builds are structurally deduplicated.
- Development uses a long-lived Dev Session outside the finite DAG. It owns
  watchers and persistent processes and requests finite incremental Plans for
  rebuilds.
- Wake becomes the typed Planner and execution engine. Reporting, Pulse,
  output capture, process execution, and generalized DAG algorithms are reused.
  Taskfile AST parsing, templates, variables, namespaces, overrides, fallback,
  command-string routing, the recursive Task executor, and its cache remain
  isolated in legacy migration/compatibility code rather than shaping the new
  engine.
- Initial typed Node kinds are frontend dependency installation, binding and
  icon generation, frontend build, application compilation, platform-asset
  generation, bundle assembly, packaging, signing, and user-owned hook
  execution. Planning and input discovery are not Nodes.
- Normal builds no longer run `go mod tidy` or otherwise rewrite `go.mod` and
  `go.sum`. Module validation reports a remediation command; source mutation
  belongs to an explicit maintenance command.

### 2026-08-15 — Grilling round 2

- The user delegated remaining internal graph decisions to engineering
  judgment. Wall-clock speed is the overriding design criterion; internal
  mechanics may iterate as benchmark evidence improves.
- Nodes are immutable declarative data. Typed handlers execute Node Kinds;
  Nodes do not contain closures, shell commands, or execution methods.
- Human-readable structural Node Keys are separate from content/cache
  fingerprints. Keys drive diagnostics and structural deduplication; content
  changes invalidate fingerprints without obscuring the Plan.
- Node scopes are Project, Target, and Package. Profiles contribute resolved
  configuration rather than creating another execution scope. Hooks inherit
  the scope of their phase.
- Every concrete output has exactly one owning Node. Inputs refer to project
  sources or declared dependency outputs. Artifacts are the user-visible
  subset of outputs, and hidden mutable communication between Nodes is
  forbidden.
- A failed Node blocks descendants while independent branches continue. The
  Plan aggregates failures; explicit user cancellation cancels all work.

## Answer

Wake becomes a static typed build engine optimized for minimum wall-clock time.
The performance strategy is to eliminate orchestration subprocesses, plan the
whole invocation once, deduplicate shared work structurally, prune cached
subgraphs before execution, and schedule the remaining critical path with
resource-aware parallelism.

### Plan contract

The Planner is read-only. It resolves the Manifest, Profile, requested Targets
and formats, project layout, frontend conventions, source topology, and
toolchain capabilities before execution. It emits one immutable Plan for the
invocation. Execution cannot add Nodes, dependencies, or outputs.

Conceptually, the core model is:

```go
type Plan struct {
    Roots []NodeKey
    Nodes map[NodeKey]Node
}

type Node struct {
    Key          NodeKey
    Kind         NodeKind
    Scope        Scope
    Dependencies []NodeKey
    Spec         NodeSpec
    Inputs       []InputRef
    Outputs      []OutputSpec
    Policy       ExecutionPolicy
}
```

`NodeSpec` is a closed typed union with one concrete spec per Node Kind. It is
data, not an executable interface. A handler selected by Kind validates and
executes the spec. The implementation avoids reflection and runtime plugin
dispatch on the hot path.

A Node Key is stable, structural, and readable, for example:

```text
frontend:bindings
frontend:build
target:darwin/arm64:compile
package:windows/amd64:msix
```

The Key is not a cache digest. The cache fingerprint separately covers the
typed Spec, resolved configuration, tool identities, relevant environment,
source inputs, and dependency results. The Planner interns Nodes by structural
Key; encountering the same Key with a different Spec is a planning error.

### Scope and ownership

Nodes have one of three scopes:

- **Project** — work shared by all Targets, such as frontend dependency
  installation, bindings, and a target-independent frontend build.
- **Target** — work for one Platform/architecture pair, such as compilation,
  platform resources, and bundle assembly.
- **Package** — work for one Target and package format, such as MSIX, DMG, deb,
  or AAB creation and signing.

Profiles are already resolved before planning and are not execution scopes.
Hooks inherit the scope of the phase they attach to.

Every concrete output path has exactly one owning Node. Inputs are either
project sources or typed references to dependency outputs. Planning rejects
output collisions, missing producers, scope-invalid dependencies, cycles,
unknown handlers, and outputs escaping allowed roots. Nodes cannot communicate
through undeclared mutable process state. An Artifact is an Output promoted to
the user-facing result summary.

### Built-in Node kinds

The initial closed set is:

```text
InstallFrontendDependencies
GenerateBindings
GenerateIcons
BuildFrontend
CompileApplication
GeneratePlatformAssets
AssembleBundle
PackageArtifact
SignArtifact
RunHook
```

One Kind may appear multiple times with different typed Specs. For example,
platform assets may be split by independent output group and package formats
always receive separate Package Nodes. This keeps expensive branches parallel
and independently cacheable without exposing micro-tasks to users.

Manifest resolution, project and toolchain inspection, Target expansion,
capability checks, and input discovery belong to planning. They are not Nodes.
Normal Plans never include `go mod tidy`; build execution does not mutate module
source files.

### Execution and scheduling

The scheduler is the only component that traverses dependencies. It uses a
non-recursive, bounded worker pool over DAG in-degrees; handlers never execute
other Nodes. This replaces the current combination of an outer DAG walk and a
recursive Task executor.

Ready work is prioritized by estimated remaining critical-path duration.
Historical Node timings refine estimates between runs. Handlers declare
internal resource claims—CPU weight, memory pressure, and exclusive tool or
credential locks—so parallel Go, Xcode, Gradle, packaging, and signing work can
be tuned from benchmark evidence without changing the Plan contract. The
scheduler aims for maximum useful concurrency rather than maximum goroutine
count.

Cache evaluation occurs before a Node enters the ready queue. A valid hit
restores or verifies declared outputs and can prune an otherwise-unneeded
dependency subtree. Cache misses on independent branches execute concurrently.
The exact fingerprint, content-store, and invalidation rules belong to
“Automatic input discovery and cache identity.”

On failure, descendants are marked blocked while unrelated Target and Package
branches continue. The final result aggregates failures and successful
Artifacts. User cancellation propagates through context cancellation to every
running process and prevents new work from starting. Cache records and
Artifacts are committed only after successful Node completion.

Handlers receive an engine-managed scratch location and publish declared
outputs through atomic promotion where the platform tool permits it. Failed or
cancelled handlers cannot make partial outputs cache-valid.

### Fast execution boundary

Handlers call Wails functionality in-process whenever possible: bindings,
icons, manifest generation, metadata generation, bundling helpers, and
reporting do not spawn `wails3` subprocesses. Unavoidable external tools are
started directly with typed argument vectors and explicit environments. The
new path never reconstructs or reparses shell command strings. `RunHook` is the
only handler allowed to invoke a user-owned script.

The existing `internal/report` contract, Pulse renderer, structured producer
events, output capture, process cancellation, and platform detection are
retained and generalized from Task names to Node Keys. DAG cycle/topological
algorithms are retained in concept but rewritten over typed Nodes. Typed
process option structures may be reused after removing command-string routing.

The Taskfile AST, parser, include/template/variable resolver, namespace and
override behavior, fallback, recursive executor, and mtime Task cache remain an
isolated legacy compatibility and migration subsystem. They do not constrain
the typed Plan or scheduler.

### Development

`wails3 dev` owns a long-lived in-process Dev Session rather than representing
watchers and servers as DAG Nodes. The session keeps the resolved project,
Plan structure, file-to-Node invalidation index, timing history, and process
handles warm. Frontend HMR remains with the frontend server; relevant backend
changes request the smallest finite incremental Plan and restart only the
application process. Bursty file events are coalesced, and stale rebuilds may
be cancelled when newer changes supersede them.

This gives build, package, and sign finite Plans while allowing development to
reuse the same Planner, scheduler, cache, handlers, and reporter without paying
full startup and discovery cost on every edit.
