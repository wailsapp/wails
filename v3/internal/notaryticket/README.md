# Portable notarization tickets

This is an original implementation of the ticket framing observed in Apple's
`stapler -v` output and independent public signed artifacts. No third-party
implementation was copied or adapted; no Go dependency was added. It reuses
Wails' `applesign` inspection and bundle copying, and `stagefile` staging.

Apple documents the notarization workflow and that stapler uses CloudKit, but
**does not publish a supported ticket-format or ticket-delivery API contract**.
This implementation accepts only the observed `s8ch` version 1 envelope and
positive `g8tk` version 2 payload, with SHA256 CodeDirectory hashes. Unknown
versions, flags, algorithms, oversized responses and malformed records fail.

`Digests` reads the primary SHA256 CodeDirectory for every executable
architecture. `Client.Fetch` requests those records from Apple's HTTPS service,
without credentials or redirects. `Verify` checks bounded lengths, the Apple
software-ticket certificate marker, digital-signature usage, a chain anchored
only at Apple's published Root CA G3 at the signed issuance time, and the
ECDSA/SHA256 signature over the complete envelope/certificate/payload prefix.
Every requested code hash must occur in each returned ticket. `Ticket` keeps
its authenticated bytes private.

`Staple` makes a fresh output and preserves the source. The caller publishes
that output after success. For canonical `.app` bundles it writes
`Contents/CodeResources`, separate from the signed resource seal in
`Contents/_CodeSignature/CodeResources`. For signed UDIF DMGs it replaces the
ticket component in the signature SuperBlob while preserving CodeDirectory and
CMS bytes, updating only the unsigned container and signature length. Existing
tickets can be replaced. Existing output paths are never overwritten. DMG
publication uses a same-volume hard link; filesystems without hard links fail
explicitly. Bundle copying preserves relative symlinks; Windows hosts need
permission to create symlinks for bundles that contain them.

This is **not** a replacement for Gatekeeper assessment. It does not verify the
executable's page hashes, resource seal or Developer ID signature, and does not
query Apple's current certificate/ticket revocation policy. A signed ticket
authenticates the listed CodeDirectories, not every byte in an arbitrary input
file. Use immutable signed outputs that the notary service accepted. Standalone
executables, unsigned/ticket-only DMGs, installer packages and other bundle
layouts are outside this interface.

## Independent acceptance and retained fixtures

On 2026-11-01, this code fetched and authenticated live Apple tickets and stapled
copies on Linux. On macOS, both the resulting Rectangle 2.0.4 universal app and
signed LuLu 4.5.3 DMG passed `codesign --verify --strict` (with `--deep` for the
app) and `xcrun stapler validate`. The app initially lacked a ticket; the signed
DMG already had one, and restapling produced identical bytes. Installed
applications were never modified. These probes demonstrate those artifacts,
not every Apple format.

Retained `.ticket` fixtures are public signed metadata, not executable source
or credentials: a ticket from a disposable copy of Tailscale's public app, its
updated CloudKit ticket (with two zero signature-padding bytes), and LuLu's
public DMG ticket. Tests need no external tool or third-party test library.
They cover corruption, unrelated hashes, bounds, response substitution,
immutable inputs, existing outputs, cancellation and DMG component preservation.
`BenchmarkVerifyTicket` retains a small CPU/allocation benchmark.

Public trust anchor source:
https://www.apple.com/certificateauthority/AppleRootCA-G3.cer

Its SHA256 is
`63343abfb89a6a03ebb57e9b3f5fa7be7c4f5c756f3017b3a8c488c3653e9179`.

Protocol references (format observations only):

- https://developer.apple.com/documentation/security/customizing-the-notarization-workflow
- https://developer.apple.com/forums/thread/720093 (ticket format is undocumented)
- Apple's Security `notarization.cpp`: bundle attachment location
- Apple's Security `codedirectory.h` and `diskimagerep.cpp`: DMG ticket slot and
  representation-header treatment
- Native `stapler -v` requests/responses and OpenSSL verification of public
  ticket signatures; no implementation source from other signers was used.
