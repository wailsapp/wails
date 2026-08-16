# Runtime-generated platform assets and customization escape hatches

Type: grilling
Status: resolved
Blocked by: 01

## Question

How should the manifest express metadata, file associations, protocols,
signing, installers, entitlements, and platform-specific packaging options so
that Wails can generate plist, NSIS, MSIX, nfpm, desktop, Xcode, iOS, and
Android inputs at build time? Where are structured fields insufficient, and
what stable user-owned template contract replaces hand-editing generated files?

## Comments

### 2026-08-16 — Accepted implementation default

- Generate all Wails-owned platform inputs into
  `.wails/build/<profile>/<target>/`; outputs are disposable and atomically
  replaced by their owning Node.
- Keep shared metadata, associations, protocols, entitlements, package-format
  options, and non-secret signing references typed in the Manifest. Exact
  Target values override Platform values.
- A user-owned template is an explicit path beside the feature that consumes
  it. Wails reads and fingerprints it but never edits it. Template variables
  are a documented, format-specific typed model rather than arbitrary process
  environment or internal Node state.
- Structured values win when supplied. A custom template owns the rendered
  file and receives the fully resolved structured model so users can replace
  presentation without replacing the Pipeline.

## Answer

The manifest-native pipeline now owns generated platform and package inputs.
It resolves shared project metadata plus Platform/Target values into an
immutable `AssetsSpec`, filters associations and protocols by their optional
`platforms` list, and atomically generates the selected target's assets under
`.wails/build/<profile>/<target>/assets`. Package-owned generated inputs live
beside them under `.wails/build/<profile>/<target>/package/<format>`; package
adapters never mutate the shared assets Artifact.

Every supported package format has a customization adapter: NSIS
`project.nsi`, MSIX `AppxManifest.xml`, macOS and iOS `Info.plist`, DMG JSON,
AppImage desktop files, nfpm YAML, and Android project directories. A custom
template is a project-relative user-owned file or directory. Wails fingerprints
but never edits it, renders `.tmpl` files with strict missing-field errors,
copies ordinary files and executable modes, rejects symlinks, and atomically
replaces the generated destination. The versioned model exposes only resolved
Project, Target (including capabilities), Package, Paths, Associations,
Protocols, and format-owned Options—not Pipeline Nodes or ambient environment.

Built-in format options remain typed: DMG resources and window dimensions and
AppImage categories are accepted without ejection; arbitrary options require a
custom template. Structured manifest values override template defaults where
both participate. Signing policy, certificate references, entitlements, and
notarization remain typed under `[signing.<platform>]`; secrets are not part of
the manifest or template model. Project-relative validation covers templates,
DMG resources, association icons, certificates, and entitlements.

The implementation also repaired the shared built-in MSIX manifest model and
added protocol propagation, which was exposed by executing the template in the
new adapter tests. Evidence:

- package renderer, Manifest, Planner, and command package tests pass;
- the same focused packages pass under the race detector and `go vet`;
- `internal/commands` cross-compiles for Windows amd64;
- a real Linux acceptance run builds the badge fixture and produces DEB, RPM,
  and Arch Linux packages from the new target-local workspaces.

Native Windows, macOS, iOS, Android, and signing execution belongs to the
separate build-performance acceptance matrix; the structural Planner matrix
already covers every supported target/format pair.
