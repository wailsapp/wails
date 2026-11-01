# Portable signing primitives

This package is an original implementation using Go's standard library. It adds
no module dependencies. Protocol references are RFC 7292 (PKCS#12), RFC 8018
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

The retained PFX fixtures were generated for these tests, using an original Go
RSA/certificate generator and OpenSSL's PFX writer. They contain no real signing
credentials. Their password is `fixture-only`; subject is
`CN=Wails Portable Signing Test Only`, serial 20261101, validity 2020–2120.
The modern fixture uses OpenSSL's PBES2/AES defaults and SHA-256 MAC; the legacy
fixture explicitly selects `PBE-SHA1-3DES` for both certificates and private keys.
Tests and benchmarks require only the standard library.
