<!-- Downstream packagers and integrators. Send before the freeze ends, not on
     launch day: their lag is what users experience as a broken release.
     Date and version are placeholders until the go/no-go fixes them. -->

# Heads-up: Wails v3 beta

Send to: Homebrew and distro package maintainers, plugin and template authors,
documentation translators, and anyone whose work breaks when ours changes.

## What is changing

**Version scheme.** The CLI moves from the `v3.0.0-alpha2.N` series to
`v3.0.0-beta.N`. A beta sorts above every alpha under semantic versioning. The
npm publish workflow assigns the `latest` dist-tag explicitly; after publishing,
verify that `npm view @wailsio/runtime dist-tags` points `latest` to the new
`3.0.0-beta.N` version rather than a `3.0.0-alpha2.*` version.

**`@wailsio/runtime` now matches the CLI version.** Until now the npm package
counted its own `alpha.N` series independently, so a CLI and its runtime
reported unrelated versions. From this release the package version is derived
from the same file the CLI embeds. If you pin the runtime, expect the version
string to jump series rather than increment.

**The current beta is tag-only.** Wails does not publish prebuilt v3 CLI
binaries, checksums or platform signatures under the current policy. Users
install the CLI from the Go module channel:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Packagers should build from the published source tag and must not wait for
GitHub release assets. If binary distribution returns later, this notice and
the release ledger will be updated before that policy is used.

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
