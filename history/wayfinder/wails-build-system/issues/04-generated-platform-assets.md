# Runtime-generated platform assets and customization escape hatches

Type: grilling
Status: open
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
