# Native package signing

`Sign` adds signatures to existing DEB and RPM packages using the existing
ProtonMail OpenPGP dependency. It does not execute a package manager, GnuPG or
another signing program. It streams package contents and preserves the existing
DEB members or RPM immutable header and compressed payload byte for byte.

The caller supplies separate input and output paths and owns transactional
publication. An error can leave a partial output; the input is never overwritten.
The key must be an exported armored or binary private key. Encrypted keys use the
supplied password. A multi-identity export requires a full fingerprint or long
key ID selector. Expired, revoked, public-only and unavailable signing keys are
rejected through the crypto library's signing-key selection.

DEB signatures use the dpkg-sig v4 clear-signed manifest, including the format's
MD5/SHA1 member digests. The OpenPGP signature uses SHA-256. Existing signatures
are preserved and included in the new signature; repeated roles receive a numeric
member suffix. This is package-level dpkg-sig compatibility, not apt repository
metadata signing or an assertion that apt automatically verifies these signatures.

RPM signing supports the v4 immutable-header format and OpenPGP v4 RSA, DSA,
ECDSA and legacy EdDSA keys. Both header and header-plus-payload signatures use
SHA-256. Existing header/payload digests are checked before adding signatures.
Package signatures are replaced while unrelated signature-header fields survive.
Header allocations are bounded; malformed indexes, strings, lengths, and regions
are rejected. No unpacking, recompression or package metadata regeneration occurs.

The `archlinux` format writes a binary detached OpenPGP signature to the supplied
output path. It shares key parsing and selection with DEB and RPM.

A bare GnuPG key ID is not a key file and cannot be resolved without access to a
GnuPG keyring. Callers retaining key-ID configurations must keep their legacy
external-tool path or ask users to export the private key. Missing or unreadable
key files must not silently fall back to a different identity.

## Validation

Unit tests use only standard-library test helpers and Wails-owned fixtures.
The encrypted fixture password is `package-test-password`; **these committed keys
are public test material and must never be used for releases**. The RSA and
Ed25519 keys were generated specifically for these tests with GnuPG. The unsigned
packages were generated through the already-installed nfpm API from an original
one-line shell payload, not copied from another project.

```sh
go test -race ./internal/packagesign
go test ./internal/packagesign -run '^$' -bench BenchmarkSignPackage -benchmem
WAILS_TEST_PACKAGE_SIGNING=1 go test ./internal/packagesign -run TestIndependentPackageReaders -v
```

The optional acceptance test requires GnuPG, dpkg-sig and RPM, with a writable
temporary RPM database. These are independent test verifiers, not runtime
dependencies. Acceptance on 2026-11-01 used dpkg-sig 0.13.1+nmu4, RPM 4.18.2 and
GnuPG 2.4.4: encrypted/plain RSA and encrypted Ed25519 signatures, repeated native
signing, and native signing after official dpkg-sig/rpmsign signing all verified.
Both official readers rejected modified payloads.

Three 200 ms benchmark samples on Linux/amd64, Go 1.26.10, Ryzen 3950X and
GOMAXPROCS=4 measured median signing times for the small fixtures:

| Key | DEB | RPM |
| --- | ---: | ---: |
| Encrypted RSA | 48.216 ms | 50.023 ms |
| Plain RSA | 2.525 ms | 3.895 ms |

These measure complete key parsing/unlocking and file signing, not large-package
throughput or a comparison against external tools. A separate CPU profile spent
about 97% in the encrypted key's iterated S2K derivation; SHA-NI hashing was the
largest leaf. Existing library/standard-library acceleration is used. No custom
assembly or weakened key derivation was justified.

## Format references and implementation provenance

The archive/signature-envelope code is original Wails code written against the
published formats and checked with independent command-line readers. No external
implementation was copied or adapted. Existing dependencies were inspected for
public API availability: nfpm and rpmpack expose creation-time signing callbacks,
but no existing-artifact rewrite API. No module was added.

- [dpkg-sig signature format](https://manpages.debian.org/unstable/dpkg-sig/dpkg-sig.1.en.html#SIGNATURE_FORMAT)
- [RPM v4 envelope](https://rpm.org/docs/4.20.x/manual/format_v4.html)
- [RPM header structure and immutable regions](https://rpm.org/docs/4.20.x/manual/format_header.html)
- [RPM digest and signature tags](https://rpm.org/docs/4.20.x/manual/tags.html)
