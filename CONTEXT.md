# Wails domain glossary

## Build system

### Manifest

The project's `wails.hcl`. It records project identity and user build intent,
not the implementation steps of the build. Every Wails project has a manifest
with its required project identity.

### Pipeline

The versioned build process supplied by Wails. A pipeline resolves manifest
intent into an executable build graph.

### Planner

The read-only part of Wake that resolves a Manifest and project facts into an
immutable Plan before execution begins.

### Plan

The immutable directed acyclic graph resolved for one Wails invocation. A Plan
may span multiple Targets while sharing Target-independent Nodes.

### Node

One schedulable and cacheable unit of work in a Plan. A Node owns typed intent,
dependencies, inputs, outputs, and produced Artifacts.

### Artifact

A declared Node output intended for another Node or surfaced to the user.

### Snapshot

The stable content identity of a Node's direct input set. A Snapshot represents
what the work consumes independently of where the project is stored.

### Action Key

The identity of one requested Node operation, derived from its typed intent,
direct input Snapshots, relevant tools and environment, and consumed Artifact
content.

### Artifact Store

The local content-addressed collection from which previously produced
Artifacts can be reused or restored.

### Receipt

Evidence that a stateful operation, such as frontend dependency installation,
completed successfully for a particular Action Key. A Receipt is not an
Artifact.

### Binding Model

The semantic services, methods, events, and data types exposed by the Go
application to its frontend. Changes outside this model do not alter generated
bindings.

### Step

The reporter's presentation of a running Node. A Step is not a user-defined
task or an execution primitive.

### Dev Session

The long-lived development controller that owns watchers and persistent
processes while requesting finite incremental Plans for rebuilds.

### Profile

A complete named production build request. A Profile selects concrete Targets,
formats, signing or notarisation intent, and bounded compiler policy. It is not
an inherited or sparse configuration overlay and cannot change project identity.

### Ejection

Materialization of the complete resolved reference Manifest, including current
Wails defaults, as the inactive `wails.ejected.hcl`. The snapshot records the
exact Wails version that generated it and never changes the active Manifest.

### Target

One concrete Platform and architecture pair produced by a Pipeline, such as
`windows/amd64`.

### Platform

An operating-system family supported by the build system, such as `darwin`,
`windows`, or `linux`. Platform configuration provides common values inherited
by its Targets.
