package updater

import "testing"

// A signature without a declared algorithm must fail verification rather
// than silently downgrading to digest-only mode.
func TestRunVerification_SignatureWithoutAlgoFailsClosed(t *testing.T) {
	v := &Verification{Signature: []byte("sig-bytes")}
	if err := runVerification(nil, v, []byte("key")); err == nil {
		t.Fatal("expected error for signature without signatureAlgo, got nil")
	}
}

// With a trust anchor configured, a release must carry a signature — a
// digest-only block or no Verification at all is the downgrade a compromised
// feed would attempt, so both must fail closed.
func TestRunVerification_PublicKeyRequiresSignature(t *testing.T) {
	key := []byte("pinned-key")
	digest := []byte{0xAA, 0xBB, 0xCC, 0xDD}

	cases := []struct {
		name string
		v    *Verification
	}{
		{"nil-verification", nil},
		{"empty-verification", &Verification{}},
		{"digest-only", &Verification{DigestAlgo: "sha256", Digest: digest}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := runVerification(digest, tc.v, key); err == nil {
				t.Fatalf("expected fail-closed error with PublicKey set for %s, got nil", tc.name)
			}
		})
	}
}

// Without a key nothing is pinned, so digest-only and no-verification releases
// still install (digest comparison still applies when a digest is present).
func TestRunVerification_NoPublicKeyPreservesDigestOnly(t *testing.T) {
	digest := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	if err := runVerification(digest, &Verification{DigestAlgo: "sha256", Digest: digest}, nil); err != nil {
		t.Fatalf("digest-only release without a key should verify, got %v", err)
	}
	if err := runVerification(digest, nil, nil); err != nil {
		t.Fatalf("release without verification and without a key should verify, got %v", err)
	}
	// A mismatching digest is still rejected.
	if err := runVerification(digest, &Verification{DigestAlgo: "sha256", Digest: []byte{0x00}}, nil); err == nil {
		t.Fatal("expected digest mismatch error, got nil")
	}
}
