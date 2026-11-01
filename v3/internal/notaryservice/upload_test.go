package notaryservice

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSignatureMatchesAWSReferencePUT(t *testing.T) {
	// Test data and expected digest from AWS's published SigV4 protocol example:
	// https://docs.aws.amazon.com/AmazonS3/latest/developerguide/sig-v4-header-based-auth.html
	request, err := http.NewRequest(http.MethodPut, "https://examplebucket.s3.amazonaws.com/test%24file.text", strings.NewReader("Welcome to Amazon S3."))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Date", "Fri, 24 May 2013 00:00:00 GMT")
	request.Header.Set("X-Amz-Storage-Class", "REDUCED_REDUNDANCY")
	digest := sha256.Sum256([]byte("Welcome to Amazon S3."))
	credentials := uploadCredentials{AccessKey: "AKIAIOSFODNN7EXAMPLE", SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"}
	signUpload(request, credentials, "us-east-1", hex.EncodeToString(digest[:]), time.Date(2013, 5, 24, 0, 0, 0, 0, time.UTC))
	expected := "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request,SignedHeaders=date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class,Signature=98ad721746da40c64f1a55b78f14c238d841ea1380cd77a1b5971af0ece108bd"
	if actual := request.Header.Get("Authorization"); actual != expected {
		t.Fatalf("SigV4 mismatch\ngot %s\nwant %s", actual, expected)
	}
}

func TestS3PathEncodingPreservesObjectIdentity(t *testing.T) {
	for input, expected := range map[string]string{"/a//b/../c": "/a//b/../c", "/a b+c%/é": "/a%20b%2Bc%25/%C3%A9", "/A_z-0.~": "/A_z-0.~", "/a?b#c": "/a%3Fb%23c"} {
		if actual := escapeS3Path(input); actual != expected {
			t.Errorf("%q = %q, want %q", input, actual, expected)
		}
	}
}

func BenchmarkArchiveSHA256(b *testing.B) {
	path := filepath.Join(b.TempDir(), "archive.zip")
	payload := make([]byte, 4<<20)
	for i := range payload {
		payload[i] = byte(i)
	}
	if err := os.WriteFile(path, payload, 0600); err != nil {
		b.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	b.SetBytes(int64(len(payload)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := file.Seek(0, 0); err != nil {
			b.Fatal(err)
		}
		if _, err := hashArchive(context.Background(), file); err != nil {
			b.Fatal(err)
		}
	}
}
