package updater

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"testing"
	"time"
)

// positiveControlKey derives a synthetic public key so the pinned key is
// non-empty without embedding any fixture literal.
func positiveControlKey(t *testing.T) []byte {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("setup: generate key: %v", err)
	}
	enc, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("setup: marshal public key: %v", err)
	}
	return enc
}

// TestWorkflowReproUnsignedReleaseInstallsWithPublicKey pins the reported
// fail-open behaviour: with Config.PublicKey set, a release that carries no
// signature (digest only) or no Verification at all must currently be
// accepted by runVerification.
func TestWorkflowReproUnsignedReleaseInstallsWithPublicKey(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfgKey := positiveControlKey(t)
	if len(cfgKey) == 0 {
		t.Fatal("setup: empty control public key")
	}

	computed := []byte{0x01, 0x02, 0x03, 0x04}

	digestOnly := &Verification{Digest: computed}
	errDigest := runVerification(computed, digestOnly, cfgKey)
	t.Logf("REPRO_OBSERVATION: digest-only err=%v keyLen=%d", errDigest, len(cfgKey))

	none := &Verification{}
	errNone := runVerification(computed, none, cfgKey)
	t.Logf("REPRO_OBSERVATION: no-verification err=%v", errNone)

	nilV := runVerification(computed, nil, cfgKey)
	t.Logf("REPRO_OBSERVATION: nil-verification err=%v", nilV)

	if errDigest == nil {
		t.Errorf("REPRO_ASSERTION: digest-only release accepted with PublicKey set, err=%v keyLen=%d", errDigest, len(cfgKey))
	}
	if errNone == nil {
		t.Errorf("REPRO_ASSERTION: release with no verification accepted with PublicKey set, err=%v", errNone)
	}
}
