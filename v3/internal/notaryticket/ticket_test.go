package notaryticket

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"
)

func fixture(t testing.TB) ([]byte, []Digest) {
	t.Helper()
	data, err := os.ReadFile("testdata/tailscale.ticket")
	if err != nil {
		t.Fatal(err)
	}
	value, err := hex.DecodeString("043baa3056b56e8b21cfaa9d522d563727754d1d")
	if err != nil {
		t.Fatal(err)
	}
	digest := Digest{Algorithm: 2}
	copy(digest.Value[:], value)
	return data, []Digest{digest}
}

func TestAppleTicket(t *testing.T) {
	data, hashes := fixture(t)
	ticket, err := Verify(data, hashes)
	if err != nil {
		t.Fatal(err)
	}
	copyBytes := ticket.Bytes()
	copyBytes[0] ^= 1
	if !bytes.Equal(ticket.Bytes(), data) {
		t.Fatal("ticket bytes not isolated")
	}
	for _, offset := range []int{0, 4, 8, 12, 30, 1536, 1548, 1552, 1561, len(data) - 1} {
		altered := bytes.Clone(data)
		altered[offset] ^= 1
		if _, err := Verify(altered, hashes); err == nil {
			t.Fatalf("accepted tampering at %d", offset)
		}
	}
	wrong := append([]Digest(nil), hashes...)
	wrong[0].Value[0] ^= 1
	if _, err := Verify(data, wrong); err == nil {
		t.Fatal("accepted unrelated ticket")
	}
	if _, err := Verify(append(data, 0), hashes); err == nil {
		t.Fatal("accepted trailing bytes")
	}
	for _, size := range []int{0, 15, 16, 1536, 1980} {
		if _, err := Verify(data[:size], hashes); err == nil {
			t.Fatal("accepted truncated ticket")
		}
	}
}

func TestAdditionalAppleEnvelopes(t *testing.T) {
	for _, test := range []struct{ name, hash string }{{"tailscale-padded.ticket", "043baa3056b56e8b21cfaa9d522d563727754d1d"}, {"lulu.ticket", "313c4bb9251b99983605844b97fc092430fd8778"}} {
		data, err := os.ReadFile("testdata/" + test.name)
		if err != nil {
			t.Fatal(err)
		}
		value, _ := hex.DecodeString(test.hash)
		digest := Digest{Algorithm: 2}
		copy(digest.Value[:], value)
		if _, err := Verify(data, []Digest{digest}); err != nil {
			t.Fatalf("%s: %v", test.name, err)
		}
	}
}

func TestRejectSelfIssuedTicketCertificate(t *testing.T) {
	data, hashes := fixture(t)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	certificate := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "Owned fake ticket authority"}, NotBefore: time.Unix(0, 0), NotAfter: time.Now().Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature, ExtraExtensions: []pkix.Extension{{Id: asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 1, 30}, Value: []byte{5, 0}}}}
	der, err := x509.CreateCertificate(rand.Reader, certificate, certificate, key.Public(), key)
	if err != nil {
		t.Fatal(err)
	}
	certs, err := asn1.Marshal([]asn1.RawValue{{FullBytes: der}})
	if err != nil {
		t.Fatal(err)
	}
	le := binary.LittleEndian
	start := 16 + int(le.Uint32(data[8:]))
	payload := data[start : start+int(le.Uint32(data[12:]))]
	forged := bytes.Clone(data[:16])
	le.PutUint32(forged[8:], uint32(len(certs)))
	forged = append(forged, certs...)
	forged = append(forged, payload...)
	digest := sha256.Sum256(forged)
	signature, err := ecdsa.SignASN1(rand.Reader, key, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	for len(signature)%4 != 0 {
		signature = append(signature, 0)
	}
	forged = append(forged, signature...)
	if _, err := Verify(forged, hashes); err == nil {
		t.Fatal("accepted attacker self-issued certificate despite valid signature")
	}
}

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetch(t *testing.T) {
	data, hashes := fixture(t)
	record := map[string]any{"recordName": "2/2/043baa3056b56e8b21cfaa9d522d563727754d1d", "recordType": "DeveloperIDTicket", "deleted": false, "fields": map[string]any{"signedTicket": map[string]any{"type": "BYTES", "value": data}}}
	for _, mode := range []string{"valid", "wrong-record", "deleted", "bad-ticket", "redirect", "oversize"} {
		t.Run(mode, func(t *testing.T) {
			body, _ := json.Marshal(map[string]any{"records": []any{record}})
			switch mode {
			case "wrong-record":
				body = bytes.ReplaceAll(body, []byte("043baa30"), []byte("143baa30"))
			case "deleted":
				body = bytes.ReplaceAll(body, []byte(`"deleted":false`), []byte(`"deleted":true`))
			case "bad-ticket":
				body = bytes.ReplaceAll(body, []byte("czhja"), []byte("azhjY"))
			case "oversize":
				body = bytes.Repeat([]byte(" "), 8<<20+1)
			}
			client := Client{HTTPClient: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() != lookupURL || r.Method != "POST" || r.Header.Get("Authorization") != "" {
					t.Fatal("unexpected request")
				}
				status := 200
				if mode == "redirect" {
					status = 302
				}
				return &http.Response{StatusCode: status, Body: io.NopCloser(bytes.NewReader(body)), Header: make(http.Header), Request: r}, nil
			})}}
			_, err := client.Fetch(context.Background(), hashes)
			if (err == nil) != (mode == "valid") {
				t.Fatalf("mode %s: %v", mode, err)
			}
		})
	}
}

func BenchmarkVerifyTicket(b *testing.B) {
	data, hashes := fixture(b)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := Verify(data, hashes); err != nil {
			b.Fatal(err)
		}
	}
}

func TestFetchMultipleArchitectures(t *testing.T) {
	data, first := fixture(t)
	verified, err := Verify(data, first)
	if err != nil {
		t.Fatal(err)
	}
	var second Digest
	for hash := range verified.hashes {
		if hash != first[0] {
			second = hash
			break
		}
	}
	for _, mode := range []string{"valid", "missing", "duplicate", "partial-coverage"} {
		t.Run(mode, func(t *testing.T) {
			expected := []Digest{first[0], second}
			if mode == "partial-coverage" {
				expected[1].Value[0] ^= 1
			}
			records := make([]map[string]any, 0, 2)
			for _, hash := range expected {
				records = append(records, map[string]any{"recordName": fmt.Sprintf("2/2/%x", hash.Value), "recordType": "DeveloperIDTicket", "fields": map[string]any{"signedTicket": map[string]any{"type": "BYTES", "value": data}}})
			}
			if mode == "missing" {
				records = records[:1]
			}
			if mode == "duplicate" {
				records[1] = records[0]
			}
			body, err := json.Marshal(map[string]any{"records": records})
			if err != nil {
				t.Fatal(err)
			}
			client := Client{HTTPClient: &http.Client{Transport: roundTrip(func(request *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(body)), Header: make(http.Header)}, nil
			})}}
			_, err = client.Fetch(t.Context(), expected)
			if (err == nil) != (mode == "valid") {
				t.Fatalf("%s: %v", mode, err)
			}
		})
	}
}
