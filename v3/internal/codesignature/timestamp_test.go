package codesignature

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func fixtureIdentity(t testing.TB) Identity {
	t.Helper()
	data, err := os.ReadFile("testdata/modern.pfx")
	if err != nil {
		t.Fatal(err)
	}
	id, err := DecodePFX(data, "fixture-only")
	if err != nil {
		t.Fatal(err)
	}
	return id
}
func timestampIdentity(t *testing.T) (Identity, *x509.CertPool) {
	return timestampIdentityWithUsage(t, x509.KeyUsageDigitalSignature)
}
func timestampIdentityWithUsage(t *testing.T, usage x509.KeyUsage) (Identity, *x509.CertPool) {
	t.Helper()
	id := fixtureIdentity(t)
	eku, err := asn1.Marshal([]asn1.ObjectIdentifier{{1, 3, 6, 1, 5, 5, 7, 3, 8}})
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{SerialNumber: big.NewInt(5), Subject: pkix.Name{CommonName: "Owned test timestamp authority"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), KeyUsage: usage, ExtraExtensions: []pkix.Extension{{Id: asn1.ObjectIdentifier{2, 5, 29, 37}, Critical: true, Value: eku}}}
	der, err := x509.CreateCertificate(rand.Reader, template, template, id.Signer.Public(), id.Signer)
	if err != nil {
		t.Fatal(err)
	}
	id.Certificate, err = x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	id.Chain = nil
	roots := x509.NewCertPool()
	roots.AddCert(id.Certificate)
	return id, roots
}
func timestampToken(t *testing.T, id Identity, digest []byte, nonce *big.Int, now time.Time, alter func(*signedData)) []byte {
	t.Helper()
	info := timestampInfo{Version: 1, Policy: asn1.ObjectIdentifier{1, 2, 3, 4}, Imprint: messageImprint{algorithm(oidSHA256), digest}, Serial: big.NewInt(123), Time: now.UTC().Truncate(time.Second), Nonce: nonce}
	body, err := asn1.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}
	contentDigest := sha256.Sum256(body)
	certDigest := sha256.Sum256(id.Certificate.Raw)
	ess := struct{ Certificates []struct{ Hash []byte } }{Certificates: []struct{ Hash []byte }{{certDigest[:]}}}
	attrs := []Attribute{}
	for _, a := range []struct {
		oid   asn1.ObjectIdentifier
		value any
	}{{oidContentType, oidTimestampInfo}, {oidMessageDigest, contentDigest[:]}, {asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 47}, ess}} {
		v, err := rawDER(a.value)
		if err != nil {
			t.Fatal(err)
		}
		attrs = append(attrs, Attribute{a.oid, []asn1.RawValue{v}})
	}
	encoded, err := asn1.MarshalWithParams(attrs, "set")
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(encoded)
	signature, err := id.Signer.Sign(rand.Reader, hash[:], crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	octets, err := asn1.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	data := signedData{Version: 3, Digests: []pkix.AlgorithmIdentifier{algorithm(oidSHA256)}, Content: contentInfo{oidTimestampInfo, explicitContent(octets)}, Certificates: asn1.RawValue{Class: 2, Tag: 0, IsCompound: true, Bytes: id.Certificate.Raw}, Signers: []signerInfo{{Version: 1, Signer: issuerSerial{asn1.RawValue{FullBytes: id.Certificate.RawIssuer}, id.Certificate.SerialNumber}, Digest: algorithm(oidSHA256), Signed: attrs, Algorithm: algorithm(oidRSA), Signature: signature}}}
	if alter != nil {
		alter(&data)
	}
	encoded, err = asn1.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	token, err := asn1.Marshal(contentInfo{oidSignedData, explicitContent(encoded)})
	if err != nil {
		t.Fatal(err)
	}
	return token
}
func TestTimestampVerificationRejectsForgedTokens(t *testing.T) {
	id, roots := timestampIdentity(t)
	digest := sha256.Sum256([]byte("signature"))
	nonce := big.NewInt(42)
	now := time.Now()
	valid := timestampToken(t, id, digest[:], nonce, now, nil)
	if err := verifyTimestamp(valid, digest[:], nonce, roots, now); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		token  []byte
		digest []byte
		nonce  *big.Int
		roots  *x509.CertPool
		now    time.Time
	}{
		{"wrong nonce", valid, digest[:], big.NewInt(43), roots, now}, {"wrong imprint", valid, make([]byte, 32), nonce, roots, now}, {"untrusted", valid, digest[:], nonce, x509.NewCertPool(), now}, {"stale", valid, digest[:], nonce, roots, now.Add(6 * time.Minute)}, {"fake ASN1", []byte{0x30, 0}, digest[:], nonce, roots, now}, {"bad signature", timestampToken(t, id, digest[:], nonce, now, func(d *signedData) { d.Signers[0].Signature[0] ^= 1 }), digest[:], nonce, roots, now}, {"missing ESS", timestampToken(t, id, digest[:], nonce, now, func(d *signedData) { d.Signers[0].Signed = d.Signers[0].Signed[:2] }), digest[:], nonce, roots, now}} {
		t.Run(test.name, func(t *testing.T) {
			if err := verifyTimestamp(test.token, test.digest, test.nonce, test.roots, test.now); err == nil {
				t.Fatal("invalid timestamp accepted")
			}
		})
	}
	codeID := fixtureIdentity(t)
	codeRoots := x509.NewCertPool()
	codeRoots.AddCert(codeID.Certificate)
	token := timestampToken(t, codeID, digest[:], nonce, now, nil)
	if err := verifyTimestamp(token, digest[:], nonce, codeRoots, now); err == nil {
		t.Fatal("code signing certificate accepted as TSA")
	}
}
func TestTimestampHTTP(t *testing.T) {
	id, roots := timestampIdentity(t)
	signature := []byte("owned signature fixture")
	expected := sha256.Sum256(signature)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.Header.Get("Content-Type") != "application/timestamp-query" {
			t.Error("wrong HTTP request")
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var query struct {
			Version int
			Imprint messageImprint
			Nonce   *big.Int
			Certs   bool
		}
		if err := decodeDER(data, &query); err != nil {
			t.Error(err)
			return
		}
		if !bytes.Equal(expected[:], query.Imprint.Digest) || !query.Certs {
			t.Error("wrong request imprint")
		}
		token := timestampToken(t, id, query.Imprint.Digest, query.Nonce, time.Now(), nil)
		reply := struct {
			Status struct{ Code int }
			Token  asn1.RawValue
		}{Token: asn1.RawValue{FullBytes: token}}
		data, err = asn1.Marshal(reply)
		if err != nil {
			t.Error(err)
			return
		}
		_, _ = w.Write(data)
	}))
	defer server.Close()
	if _, err := requestTimestamp(context.Background(), signature, Options{TimestampURL: server.URL, TimestampRoots: roots}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := requestTimestamp(ctx, signature, Options{TimestampURL: server.URL, TimestampRoots: roots}); err == nil {
		t.Fatal("canceled request accepted")
	}
}
func TestTimestampHTTPRejectsFakeAndRedirects(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   string
	}{{"fake success", 200, "not a timestamp"}, {"unavailable", 503, "retry"}, {"redirect", 302, ""}, {"oversized", 200, strings.Repeat("x", (4<<20)+1)}} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", "http://example.invalid/")
				w.WriteHeader(test.status)
				_, _ = io.WriteString(w, test.body)
			}))
			defer server.Close()
			if _, err := requestTimestamp(context.Background(), []byte("signature"), Options{TimestampURL: server.URL}); err == nil {
				t.Fatal("invalid timestamp accepted")
			}
		})
	}
}

func TestTimestampCertificateKeyUsage(t *testing.T) {
	digest := sha256.Sum256([]byte("signature"))
	nonce := big.NewInt(42)
	now := time.Now()
	for _, usage := range []x509.KeyUsage{x509.KeyUsageKeyEncipherment, x509.KeyUsageContentCommitment} {
		id, roots := timestampIdentityWithUsage(t, usage)
		token := timestampToken(t, id, digest[:], nonce, now, nil)
		err := verifyTimestamp(token, digest[:], nonce, roots, now)
		if (err == nil) != (usage == x509.KeyUsageContentCommitment) {
			t.Fatalf("key usage %v: %v", usage, err)
		}
	}
}
