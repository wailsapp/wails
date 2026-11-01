package codesignature

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestDarwinTimestampSystemTrustRejectsUntrusted(t *testing.T) {
	id, roots := timestampIdentity(t)
	now := time.Now()
	digest := sha256.Sum256([]byte("platform trust regression"))
	nonce := big.NewInt(42)
	token := timestampToken(t, id, digest[:], nonce, now, nil)
	if err := verifyTimestamp(token, digest[:], nonce, nil, now); err == nil {
		t.Fatal("platform trust accepted a self-signed test certificate")
	}
	if err := verifyTimestamp(token, digest[:], nonce, roots, now); err != nil {
		t.Fatalf("explicit test roots rejected: %v", err)
	}
	if err := verifyTimestampChain(id.Certificate, []*x509.Certificate{id.Certificate}, roots, id.Certificate.NotAfter.Add(time.Second)); err == nil {
		t.Fatal("expired certificate accepted with explicit roots")
	}
}

// Opt-in acceptance exercises the platform root store against a real RFC 3161
// authority. Normal test runs stay offline and require no external credentials.
func TestDarwinTimestampLiveSystemTrust(t *testing.T) {
	url := os.Getenv("WAILS_TEST_TIMESTAMP_URL")
	if url == "" {
		t.Skip("set WAILS_TEST_TIMESTAMP_URL for live system-trust acceptance")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := requestTimestamp(ctx, []byte("Wails owned timestamp trust acceptance"), Options{TimestampURL: url}); err != nil {
		t.Fatal(err)
	}
}
