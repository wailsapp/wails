# Portable Windows signatures

`SignPE` and `SignMSIX` sign with exported RSA identities on any Go host. They
use the shared `codesignature` package for authenticated PFX decoding, SHA-256
CMS signatures and verified RFC 3161 timestamps. No modules were added and no
third-party signing implementation was copied or adapted.

PE signing follows Microsoft's [PE format](https://learn.microsoft.com/en-us/windows/win32/debug/pe-format)
and [Authenticode PE specification](https://download.microsoft.com/download/9/c/5/9c5b2167-8017-4bae-9fde-d599bac8184a/Authenticode_PE.docx):
the checksum, certificate directory entry and certificate table are excluded
from the image digest. An aligned certificate table at the end of a previously
signed image is replaced; unsupported nonterminal tables fail explicitly.

MSIX signing uses the APPX SIP digest identifiers declared by Microsoft's SDK:
AXPC authenticates local file records, AXCD the central directory and unsigned
ZIP footer, AXCT content types, AXBM the block map, and optional AXCI the code
integrity catalog. Owned SignTool-generated fixtures independently established
these ranges and the PKCX/CMS encoding. The signature entry's ZIP creator version
must be 45, matching the ZIP64 footer, for native verification to accept it.
Existing local file records are copied literally, preserving local extra fields
that may differ from central-directory metadata. Only content types and the
signature footprint are changed. Payload hashes, block counts, compressed block
boundaries, local-header sizes and publisher/certificate RDNs are checked first.
The original code-integrity catalog, if present, is preserved and authenticated.

Input is copied into a private immutable snapshot before validation/signing and
closed before publication. This supports in-place replacement on Windows and
prevents input changes during a timestamp request from invalidating the output.
Failures remove temporary files and preserve the prior output.

Supported formats are PE32/PE32+ and individual MSIX/APPX packages smaller than
4 GiB. MSI, MSIX bundles, ZIP comments and packages with
65535 or more entries are rejected. ZIP64 local/central metadata and end-of-directory records produced by
native package tools are supported when individual sizes and offsets fit in
32 bits. Signature certificates currently require RSA 2048–8192 and SHA-256.
Certificate-store, hardware/cloud-provider and unsupported PFX profiles remain
platform signing-tool workflows; this package does not emulate those providers.

Validation uses retained standard-library tests for signing, re-signing,
transactions, immutable input, publisher parsing, package corruption, metadata
and fake timestamps. External native checks use Windows SDK SignTool plus
WinVerifyTrust with an in-memory exclusive root containing only an owned test
certificate. Valid PE/MSIX signatures pass; tampered counterparts fail with
TRUST_E_BAD_DIGEST. This test trust does not establish production publisher trust.
