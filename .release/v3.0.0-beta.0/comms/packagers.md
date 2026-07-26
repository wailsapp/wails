<!-- Downstream packagers and integrators. Send before the freeze ends, not on
     launch day: their lag is what users experience as a broken release.
     Date and version are placeholders until the go/no-go fixes them. -->

# Heads-up: Wails v3 beta

Send to: Homebrew and distro package maintainers, plugin and template authors,
documentation translators, and anyone whose work breaks when ours changes.

## What is changing

**Version scheme.** The CLI moves from the `v3.0.0-alpha2.N` series to
`v3.0.0-beta.N`. A beta sorts above every alpha under semantic versioning, so
`@latest` resolves to it.

**`@wailsio/runtime` now matches the CLI version.** Until now the npm package
counted its own `alpha.N` series independently, so a CLI and its runtime
reported unrelated versions. From this release the package version is derived
from the same file the CLI embeds. If you pin the runtime, expect the version
string to jump series rather than increment.

**Release artifacts now ship with checksums and provenance.** Each release
carries a `SHA256SUMS` file generated in the same job that publishes the
binaries, plus a signed build-provenance attestation. If you verify downloads,
you can now do so properly.

**Artifact names are unchanged:** `wails3-linux-amd64`, `wails3-linux-arm64`,
`wails3-windows-amd64.exe`, `wails3-darwin-amd64`, `wails3-darwin-arm64`,
`wails3-darwin-universal`.

**macOS signing.** Whether the darwin binaries are signed and notarized in this
release depends on the signing credentials being in place. If they are not, the
release will say so explicitly rather than shipping silently unsigned binaries.

**The repository no longer ships a `go.work`.** If your build tooling assumed a
workspace at the repository root, it does not need one now.

## What we need from you

- Test the release candidate against your package before the date below.
- Tell us if any of the above breaks your build, while there is still time.
- Note that this is a beta: the API is stable and it receives security fixes,
  but it is not the stable 3.0.

## Dates

- Release candidate available: TBC
- Deadline for telling us it broke: TBC
- Release date: TBC

## Contact

Named contact for release week: TBC
