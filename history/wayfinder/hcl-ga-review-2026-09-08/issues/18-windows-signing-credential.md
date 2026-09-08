# Pass the declared Windows signing credential to the signer

Type: task
Status: resolved
Blocked by: none
Label: ready-for-agent
Priority: P1

## Question

The HCL signing credential names a password environment variable, but the handler passes it only as a macOS keychain-profile field. Windows signing receives an empty password and falls back to unrelated global credentials. The real HCL-to-handler regression fails with expected test-pfx-password, actual empty string.

## Acceptance criteria

Resolve the declared Windows password at execution time and pass it to the signer. Reject an explicitly named missing credential before signing. Keep secrets out of plans and diagnostics. Verify protected PFX signing on Windows.

## Comments

2026-09-08: Reproduced on Linux using the real HCL plan and signer boundary. Native verification remains with the platform owner.

## Answer

2026-09-08: Fixed credential resolution and missing-credential failure; real HCL signer-boundary regression passes. Protected PFX native verification remains in matching-host acceptance.
