# Wails domain glossary

## Build system

### Manifest

The project's `wails.toml`. It records project identity and user build intent,
not the implementation steps of the build. Every Wails project has a manifest
with its required project identity.

### Pipeline

The versioned build process supplied by Wails. A pipeline resolves manifest
intent into an executable build graph.

### Profile

A named set of build-policy overrides, such as debug or release. Profiles are
sparse overlays over the manifest's effective top-level configuration until
explicitly ejected as a complete frozen snapshot. A Profile controls build
policy and cannot change project identity.

### Ejection

Materialization of an effective base or named Profile into a complete frozen
snapshot, including every Target reachable by that Profile. The snapshot
records which Wails version ejected it and does not inherit later defaults.

### Target

One concrete Platform and architecture pair produced by a Pipeline, such as
`windows/amd64`.

### Platform

An operating-system family supported by the build system, such as `darwin`,
`windows`, or `linux`. Platform configuration provides common values inherited
by its Targets.
