# Portable notarization client

This internal package implements Apple's documented Notary API with Go's standard
library. It creates submissions, uploads signed archives, reads their status, and
retrieves JSON logs. It does not sign software or staple tickets.

`New` takes an App Store Connect **team** key ID, issuer ID, and the PEM contents
of its downloaded PKCS#8 P-256 `.p8` private key. Individual API keys, Apple ID
passwords, and macOS keychain profiles are not inputs to this client. Existing
profile-based callers must keep their separate `notarytool` path when explicit
API credentials are absent. Callers should load the private key from a protected
file and keep its bytes out of configuration files, logs, and receipts.

The submission sequence is:

1. `Submit(ctx, archive)` hashes and streams a nonempty ZIP, DMG, or PKG file.
   It returns the submission ID even if upload fails after Apple creates it.
2. `Wait(ctx, id, interval)` polls status. Accepted, Invalid, and Rejected are
   terminal results; the caller must distinguish them. A context deadline bounds
   the overall wait. HTTP failures stop polling and expose a typed `HTTPError`.
3. `Log(ctx, id)` downloads the JSON report. Retrieve it even after acceptance,
   because it can contain warnings. An invalid submission is not successful just
   because its log was retrieved.

Each API request gets a five-minute ES256 JWT. Uploads use Apple's temporary S3
credentials with SHA-256 payload verification and AWS Signature Version 4. An
anonymous `HeadBucket` request obtains the bucket region, including when its
response is 403. The upload uses the S3 acceleration endpoint. Redirects and
cookies are disabled; neither the Apple bearer token nor S3 credentials are
forwarded to the temporary log URL. Response errors omit remote bodies and
signed URLs. API responses are limited to 1 MiB and log downloads to 16 MiB.
For Apple's specific missing-or-expired-agreements error, the client adds a fixed
instruction for the Account Holder to review agreements in the developer account.
It never displays the server's error title, detail, or supplied links; unknown
API errors and S3 errors retain the generic HTTP status message.

The client supports a single S3 PUT, capped at **5 GiB**. Larger archives are
rejected before creating a submission. Multipart and automatic POST/PUT retries
are not implemented. Upload errors therefore do not silently create replacement
submissions. API calls and log downloads have one-minute request deadlines; a
nil HTTP client selects a 30-minute HTTP timeout for uploads. Caller context
cancellation also interrupts hashing, HTTP operations, and polling waits.

Tests exercise the complete flow with an HTTP test server, verify ES256
signatures with `crypto/ecdsa`, and check the S3 signature against AWS's published
PUT test vector. They also cover cancellation, rejected submissions, malformed
responses, upload failure IDs, response limits, redirect refusal, and credential
separation. These tests do **not** validate a real Apple account or production
upload. Benchmarks for JWT signing and streamed archive hashing use no external
libraries.

Ticket retrieval and stapling belong to the separate
[`notaryticket`](../notaryticket) package, which follows Apple's observed ticket
protocol. Gatekeeper can find tickets online; stapling supports offline
verification. This client remains responsible for submission, status and logs.

## Live acceptance

On 2026-11-01, a stdlib Go arm64 application signed with Developer ID, hardened
runtime and a secure timestamp completed the live workflow from Linux: the
Wails archive writer created the ZIP, this client submitted and uploaded it,
polled Apple's `Accepted` result and retrieved the diagnostic log, and the Go
ticket package verified and stapled the returned ticket. The resulting bundle
passed macOS `codesign --verify --strict --deep`, `stapler validate` and
Gatekeeper assessment with `source=Notarized Developer ID`.

Code-signature generation for this acceptance test used native `codesign`.
The initial account-agreement rejection exercised the fixed guidance above.
A separate tiny compiler-fixture app was rejected by Apple despite passing
native signature checks; its signed bytes were unchanged by copying and ZIP
creation. Both outcomes are retained in the test evidence, and generated
signing fixtures are not treated as proof of notarization acceptance.

## Protocol sources

This implementation was written from the protocol descriptions, without copying
third-party implementations:

- [Apple: submitting software over the web](https://developer.apple.com/documentation/notaryapi/submitting-software-for-notarization-over-the-web)
- [Apple: API token generation](https://developer.apple.com/documentation/appstoreconnectapi/generating-tokens-for-api-requests)
- [Apple: submission status](https://developer.apple.com/documentation/notaryapi/get-submission-status)
- [Apple: submission logs](https://developer.apple.com/documentation/notaryapi/get-submission-log)
- [AWS: HeadBucket](https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadBucket.html)
- [AWS: single-chunk Signature Version 4, including test vectors](https://docs.aws.amazon.com/AmazonS3/latest/developerguide/sig-v4-header-based-auth.html)
- [AWS: uploading objects and single-PUT limit](https://docs.aws.amazon.com/AmazonS3/latest/userguide/upload-objects.html)
- [Apple: custom notarization workflow and stapling](https://developer.apple.com/documentation/security/customizing-the-notarization-workflow)
