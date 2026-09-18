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
