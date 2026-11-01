# Portable signing primitives

This package is an original implementation using Go's standard cryptography
library. On macOS, system timestamp trust uses Security.framework through the
existing purego dependency. It adds no module dependencies. Protocol references
are RFC 7292 (PKCS#12), RFC 8018
(PBES2/PBKDF2), RFC 5652 (CMS), RFC 3161 and RFC 5035 (timestamp certificates).

`DecodePFX` requires a verified password MAC and exactly one private key with a
matching certificate. Supported encryption is PBES2/PBKDF2/AES-CBC and legacy
SHA-1/three-key 3DES. RC2, unauthenticated containers, multiple private keys and
other algorithms fail explicitly. Re-export an unsupported legacy PFX with
AES encryption; platform signing tools remain available for other profiles.
Input size, bag counts, password size and cumulative derivation work are bounded.

Signing currently requires RSA 2048–8192, SHA-256, a matching valid certificate
and code-signing usage. `SignDetached` signs id-data and accepts format-specific
signed attributes. `SignPKCS7` provides legacy content-octet hashing for
Authenticode. Neither function establishes public trust in a signing identity.

RFC 3161 responses are bounded, cancellation-aware and verified before use:
nonce, signature imprint, current generation time, CMS signature, signed content
digest, ESS certificate binding, exclusive critical timestamping EKU and an
X.509 chain to system or explicitly supplied roots. Redirects are disabled.
Timestamp authorities using unsupported signature profiles fail explicitly.
macOS system trust uses Apple's timestamp policy, rather than the TLS policy
used by Go's system verifier. Verification uses the token's generation time,
the platform trust store and supplied intermediates, without additional network
fetches. Explicitly supplied roots continue to use Go's verifier on every OS.

The retained PFX fixtures were generated for these tests, using an original Go
RSA/certificate generator and OpenSSL's PFX writer. They contain no real signing
credentials. Their password is `fixture-only`; subject is
`CN=Wails Portable Signing Test Only`, serial 20261101, validity 2020–2120.
The modern fixture uses OpenSSL's PBES2/AES defaults and SHA-256 MAC; the legacy
3DES fixture explicitly selects `PBE-SHA1-3DES` for both certificates and private keys.
The RC2 fixture uses OpenSSL `-legacy` defaults to prove that the unsupported
encryption sentinel is returned only after MAC verification; bad passwords and
tampered data cannot select a native fallback.
Test and benchmark helpers require only the standard library. Darwin tests also
exercise the production Security.framework adapter.

## macOS timestamp acceptance

On macOS 27.0.1 with Go 1.26.10, the original system verifier rejected a live
DigiCert SHA-256 timestamp responder as "not standards compliant": the system
chain verifier selected an SSL policy for a timestamp-only certificate. The
timestamp-policy adapter accepts the same authority using normal platform roots.
The native suite also rejects an untrusted generated timestamp identity and
retains explicit-root, nonce, imprint, signature, ESS and key-usage checks.

An owned nested application signed by this package with the generated fixture
identity and a live DigiCert timestamp passed `codesign --verify --strict --deep`;
`codesign -dvv` displayed the timestamp. This validates the timestamp and Apple
signature format, not public trust in the generated code-signing identity or
notarization. Run the optional live test on macOS with
`WAILS_TEST_TIMESTAMP_URL=http://timestamp.digicert.com go test ./internal/codesignature`.
