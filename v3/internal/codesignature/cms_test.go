package codesignature

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"testing"
)

func TestDetachedCMS(t *testing.T) {
	id := fixtureIdentity(t)
	content := []byte("owned code directory")
	encoded, err := SignDetached(context.Background(), id, content, nil, Options{})
	if err != nil {
		t.Fatal(err)
	}
	var outer contentInfo
	if err := decodeDER(encoded, &outer); err != nil {
		t.Fatal(err)
	}
	var data signedData
	if err := decodeDER(outer.Content.Bytes, &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Content.Content.FullBytes) != 0 || !data.Content.Type.Equal(oidData) || len(data.Signers) != 1 {
		t.Fatal("CMS is not detached id-data")
	}
	signer := data.Signers[0]
	attrs, err := asn1.MarshalWithParams(signer.Signed, "set")
	if err != nil {
		t.Fatal(err)
	}
	if err := id.Certificate.CheckSignature(x509.SHA256WithRSA, attrs, signer.Signature); err != nil {
		t.Fatal(err)
	}
	want := sha256.Sum256(content)
	found := false
	for _, attr := range signer.Signed {
		if attr.Type.Equal(oidMessageDigest) {
			var got []byte
			if err := decodeDER(attr.Values[0].FullBytes, &got); err != nil {
				t.Fatal(err)
			}
			if string(got) != string(want[:]) {
				t.Fatal("incorrect content digest")
			}
			found = true
		}
	}
	if !found {
		t.Fatal("missing content digest")
	}
	attrs[5] ^= 1
	if err := id.Certificate.CheckSignature(x509.SHA256WithRSA, attrs, signer.Signature); err == nil {
		t.Fatal("tampered attributes accepted")
	}
}
func TestDuplicateCMSAttributes(t *testing.T) {
	id := fixtureIdentity(t)
	value, _ := rawDER(oidData)
	_, err := SignDetached(context.Background(), id, []byte("data"), []Attribute{{oidContentType, []asn1.RawValue{value}}}, Options{})
	if err == nil {
		t.Fatal("duplicate mandatory attribute accepted")
	}
}
func BenchmarkSignDetached(b *testing.B) {
	id := fixtureIdentity(b)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := SignDetached(context.Background(), id, []byte("owned code directory"), nil, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}
