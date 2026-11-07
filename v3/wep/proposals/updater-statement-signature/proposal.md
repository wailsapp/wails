# Wails Enhancement Proposal (WEP)

## Updater statement signature

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Damon Blais  
**Created**: 2026-11-07  
**Discussion**: none  
**Implementor**: Damon Blais  
**Target**: Wails v3

## Summary

Add an optional second signature to each artifact in an update manifest. It signs the release's version, platform, arch and digest together. A new `Config.RequireStatement` makes the updater install only a release whose statement verifies and whose version is newer than the running one. `wails3 updater statement` writes the signatures and `wails3 updater verify -require-statement` checks them.

## Motivation

A manifest artifact is signed over its file: `ed25519` signs the digest, `ed25519ph` signs the file's SHA-512. Nothing signed ties that file to the entry that points at it. The manifest's `version`, and each entry's `platform` and `arch`, are plain JSON.

Whoever controls the feed, or a cache or proxy in front of it, can therefore:

- serve an older, vulnerable build that is still validly signed, labelled as the newer release;
- serve another platform's or architecture's signed build to a client that asked for its own.

The signature verifies in both cases, because the file is genuine. Pinning a public key does not help.

## Detailed Design

Reference: commit `2694a7d01` on https://github.com/Rethunk-AI/wails/tree/feat/updater-statement-signature. That branch also carries two unrelated updater hardening commits that are not part of this proposal.

**Statement.** `updater.Statement(version, platform, arch, digestAlgo, digest)` returns the signed bytes:

```
wails-update-statement-v1
version=<version>
platform=<platform>
arch=<arch>
digest=<algo>:<hex digest>
```

The algorithm is `sha512` when the entry's `signatureAlgo` is `ed25519ph`, otherwise the entry's `digestAlgo`, defaulting to `sha256`.

**Manifest.** Each artifact gains an optional `statementSignature`: the base64 Ed25519 signature of that text. The endpoint provider decodes it into `Verification.StatementSignature` and fails on invalid base64. An entry with only a statement signature still produces a verification block.

**Signing.** `wails3 updater statement -manifest manifest.json -key <private key>` reads a finished manifest and signs a statement for every artifact. It is a separate step so it runs after any edit to an entry's platform or arch. It uses the same private key as `updater manifest -key` and fails if an entry has no usable digest.

**Checking in the app.** `Config.RequireStatement` (default false), when set:

1. a release with no verification block, or no statement signature, is refused;
2. the existing signature check runs first;
3. the statement is rebuilt from the release's version, **this Updater's** `Platform` and `Arch`, and the digest computed while downloading, and verified under `Config.PublicKey`, which must be Ed25519;
4. the release's version must be newer than `CurrentVersion`.

The platform and arch come from the running app, not from the manifest, so a mislabelled entry fails.

**Checking in CI.** `wails3 updater verify -require-statement` fails any artifact without a valid statement signature. Without the flag, a statement signature that is present is still verified.

## Non-Goals

- Rollback protection beyond "newer than the running version". No monotonic counter is stored.
- Signing the manifest as a whole, or its channel, size or URLs.
- ECDSA P-256 statement signatures.
- Changing how existing `signature` fields are produced or checked.

## Platform Considerations

Platform independent: Go code in `pkg/updater`, `pkg/updater/providers/endpoint` and the `wails3` CLI. It behaves the same on Windows, macOS, Linux and mobile wherever the updater runs. The platform and arch it binds are the strings the Updater is already configured with.

## Pros/Cons

Pros:

- Closes the downgrade and cross-platform relabelling attacks that a pinned key alone does not.
- Optional, no change to existing manifests or apps.
- Reuses the existing key and digest, with no new key to manage.

Cons:

- One more publishing step (`updater statement`) that must run after the manifest is final.
- Ed25519 only.
- Enabling it refuses every release published without a statement.

## Alternatives Considered

- **Sign the whole manifest.** Covers more, but the feed serves a manifest per platform and channel and tooling edits it after generation; per-entry signatures survive that. It remains compatible with this proposal.
- **Put version and target inside the signed digest payload.** Would change the meaning of existing `signature` values and break current clients.
- **Check version and platform only client-side.** The platform check already exists via entry selection, but the entry is what the feed controls, so it proves nothing.
- **Make `updater manifest` sign statements in the same pass.** Rejected because the manifest is commonly edited afterwards, which invalidates them.

## Backwards Compatibility

Additive. `statementSignature` is optional and older clients ignore the unknown field. `RequireStatement` defaults to false. An app turns it on once every release it should still install carries a statement. `updater verify` without the flag behaves as before except that a present statement is now checked.

## Security and Privacy

Improves release authenticity. The statement is signed with the key that already signs artifacts, so compromise of that key remains total. Domain separation comes from the `wails-update-statement-v1` prefix, so a statement cannot be confused with a file signature. No data is collected.

## Test Plan

Unit tests in the reference branch cover: a correctly signed release is accepted; a downgrade or relabelled version is refused; another platform's file is refused; a release without a statement is refused when required; and `updater statement` output passing `updater verify`. Run with `go test ./pkg/updater/... ./internal/commands`.

## Reference Implementation

https://github.com/Rethunk-AI/wails/tree/feat/updater-statement-signature (commit `2694a7d01`).

## Maintenance Plan

The Implementor maintains it. The `-v1` in the statement text allows a later format without ambiguity.

## Open Questions

- Should the manifest's `version` also be signed once per manifest rather than repeated per entry?
- Should `RequireStatement` become the default in a later release, with an opt-out?
- Is an Ed25519-only restriction acceptable, or should ECDSA keys be supported?
- Should rollback protection store the highest version seen?

## Conclusion

A signed file does not say which release or platform it belongs to. The statement signature adds that, optionally, with one extra signing step and a single config switch.
