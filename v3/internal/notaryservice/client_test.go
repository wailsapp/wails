package notaryservice

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const testID = "2efe2717-52ef-43a5-96dc-0797e4ca1041"

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func testCredentials(t testing.TB) (Credentials, *ecdsa.PrivateKey) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return Credentials{KeyID: "TESTKEY", IssuerID: "test-team", PrivateKeyPEM: pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})}, key
}

func localClient(t *testing.T, handler http.HandlerFunc) (*Client, *ecdsa.PrivateKey) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	address, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	t.Cleanup(transport.CloseIdleConnections)
	client := &http.Client{Transport: transportFunc(func(request *http.Request) (*http.Response, error) {
		copyRequest := request.Clone(request.Context())
		copyRequest.URL.Scheme = address.Scheme
		copyRequest.URL.Host = address.Host
		copyRequest.Host = request.URL.Host
		return transport.RoundTrip(copyRequest)
	})}
	credentials, key := testCredentials(t)
	result, err := New(credentials, client)
	if err != nil {
		t.Fatal(err)
	}
	return result, key
}

func verifyToken(t *testing.T, token string, key *ecdsa.PrivateKey, now time.Time) {
	t.Helper()
	pieces := strings.Split(token, ".")
	if len(pieces) != 3 {
		t.Errorf("invalid JWT shape")
		return
	}
	var header struct{ Alg, Kid, Typ string }
	var payload struct {
		Iss, Aud string
		Iat, Exp int64
	}
	data, err := base64.RawURLEncoding.DecodeString(pieces[0])
	if err != nil {
		t.Error(err)
		return
	}
	if err = json.Unmarshal(data, &header); err != nil {
		t.Error(err)
		return
	}
	data, err = base64.RawURLEncoding.DecodeString(pieces[1])
	if err != nil {
		t.Error(err)
		return
	}
	if err = json.Unmarshal(data, &payload); err != nil {
		t.Error(err)
		return
	}
	if header.Alg != "ES256" || header.Kid != "TESTKEY" || header.Typ != "JWT" || payload.Iss != "test-team" || payload.Aud != "appstoreconnect-v1" || payload.Iat > now.Unix() || payload.Exp <= now.Unix() || payload.Exp-payload.Iat > 1200 {
		t.Errorf("incorrect JWT claims: %+v %+v", header, payload)
	}
	signature, err := base64.RawURLEncoding.DecodeString(pieces[2])
	if err != nil || len(signature) != 64 {
		t.Errorf("JWT signature must be 64 raw bytes: %v", err)
		return
	}
	digest := sha256.Sum256([]byte(pieces[0] + "." + pieces[1]))
	if !ecdsa.Verify(&key.PublicKey, digest[:], new(big.Int).SetBytes(signature[:32]), new(big.Int).SetBytes(signature[32:])) {
		t.Error("JWT signature did not verify")
	}
}

func TestTokenSignatureAndClaims(t *testing.T) {
	credentials, key := testCredentials(t)
	client, err := New(credentials, nil)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(1700000000, 0)
	token, err := client.token(now)
	if err != nil {
		t.Fatal(err)
	}
	verifyToken(t, token, key, now)
}

func TestCredentialsRejectMalformedAndWrongCurve(t *testing.T) {
	credentials, _ := testCredentials(t)
	for _, field := range []string{"key", "issuer", "pem", "curve"} {
		t.Run(field, func(t *testing.T) {
			value := credentials
			switch field {
			case "key":
				value.KeyID = ""
			case "issuer":
				value.IssuerID = ""
			case "pem":
				value.PrivateKeyPEM = []byte("secret-malformed-key")
			case "curve":
				key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
				if err != nil {
					t.Fatal(err)
				}
				der, err := x509.MarshalPKCS8PrivateKey(key)
				if err != nil {
					t.Fatal(err)
				}
				value.PrivateKeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
			}
			if _, err := New(value, nil); err == nil || strings.Contains(err.Error(), "secret-malformed-key") {
				t.Fatalf("unsafe or absent credential validation: %v", err)
			}
		})
	}
}

func uploadResponse(w http.ResponseWriter) {
	_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": testID, "attributes": uploadCredentials{AccessKey: "test-access", SecretKey: "test-secret", SessionToken: "test-session", Bucket: "test-bucket", Object: "folder//app +%é.zip"}}})
}

func TestSubmitWaitAndLog(t *testing.T) {
	payload := []byte("signed application fixture")
	digest := sha256.Sum256(payload)
	var statusCalls atomic.Int32
	var uploaded atomic.Bool
	var key *ecdsa.PrivateKey
	client, generatedKey := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "appstoreconnect.apple.com" {
			verifyToken(t, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "), key, time.Now())
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/notary/v2/submissions":
			var body map[string]string
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body["submissionName"] != "app.zip" || body["sha256"] != hex.EncodeToString(digest[:]) {
				t.Errorf("wrong submission: %#v", body)
			}
			uploadResponse(w)
		case r.Method == http.MethodHead && r.Host == "test-bucket.s3.amazonaws.com":
			if r.Header.Get("Authorization") != "" {
				t.Error("Apple authorization leaked to bucket discovery")
			}
			w.Header().Set("X-Amz-Bucket-Region", "eu-west-1")
			w.WriteHeader(http.StatusForbidden)
		case r.Method == http.MethodPut && r.Host == "test-bucket.s3-accelerate.amazonaws.com":
			data, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			if string(data) != string(payload) || r.ContentLength != int64(len(payload)) {
				t.Error("upload bytes or length differ")
			}
			if r.URL.EscapedPath() != "/folder//app%20%2B%25%C3%A9.zip" {
				t.Errorf("upload path changed: %s", r.URL.EscapedPath())
			}
			if !strings.Contains(r.Header.Get("Authorization"), "/eu-west-1/s3/aws4_request") || r.Header.Get("X-Amz-Security-Token") != "test-session" || r.Header.Get("X-Amz-Content-Sha256") != hex.EncodeToString(digest[:]) {
				t.Error("upload authorization or digest missing")
			}
			uploaded.Store(true)
		case r.URL.Path == "/notary/v2/submissions/"+testID:
			if !uploaded.Load() {
				t.Error("polled before upload")
			}
			status := "In Progress"
			if statusCalls.Add(1) > 1 {
				status = "Accepted"
			}
			fmt.Fprintf(w, `{"data":{"id":%q,"attributes":{"status":%q,"name":"app.zip"}}}`, testID, status)
		case r.URL.Path == "/notary/v2/submissions/"+testID+"/logs":
			fmt.Fprintf(w, `{"data":{"id":%q,"attributes":{"developerLogUrl":"https://logs.example/report?signature=private"}}}`, testID)
		case r.Host == "logs.example":
			if r.Header.Get("Authorization") != "" || r.Header.Get("X-Amz-Security-Token") != "" {
				t.Error("credentials leaked to log download")
			}
			fmt.Fprint(w, `{"status":"Accepted","issues":[{"severity":"warning"}]}`)
		default:
			t.Errorf("unexpected request %s %s %s", r.Method, r.Host, r.URL)
			w.WriteHeader(404)
		}
	})
	key = generatedKey
	path := filepath.Join(t.TempDir(), "app.zip")
	if err := os.WriteFile(path, payload, 0600); err != nil {
		t.Fatal(err)
	}
	result, err := client.Submit(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != testID {
		t.Fatalf("submission ID = %q", result.ID)
	}
	result, err = client.Wait(context.Background(), result.ID, time.Millisecond)
	if err != nil || result.Status != "Accepted" || statusCalls.Load() != 2 {
		t.Fatalf("poll result %+v %v", result, err)
	}
	log, err := client.Log(context.Background(), result.ID)
	if err != nil || !strings.Contains(string(log), "warning") {
		t.Fatalf("missing log warnings: %s %v", log, err)
	}
}

func TestUploadFailureRetainsSubmissionID(t *testing.T) {
	client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			uploadResponse(w)
		case http.MethodHead:
			w.Header().Set("X-Amz-Bucket-Region", "us-west-2")
			w.WriteHeader(403)
		default:
			w.WriteHeader(403)
			fmt.Fprint(w, "test-secret test-session")
		}
	})
	path := filepath.Join(t.TempDir(), "app.dmg")
	if err := os.WriteFile(path, []byte("archive"), 0600); err != nil {
		t.Fatal(err)
	}
	result, err := client.Submit(context.Background(), path)
	var responseError *HTTPError
	if result.ID != testID || !errors.As(err, &responseError) || responseError.StatusCode != 403 || strings.Contains(err.Error(), "test-secret") {
		t.Fatalf("upload failure lost ID or exposed credentials: %+v %v", result, err)
	}
}

func TestOversizedArchiveRejectedBeforePOST(t *testing.T) {
	var calls atomic.Int32
	client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(500) })
	path := filepath.Join(t.TempDir(), "large.zip")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	err = file.Truncate(MaxUploadBytes + 1)
	closeErr := file.Close()
	if err != nil {
		t.Skipf("sparse files unavailable: %v", err)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	_, err = client.Submit(context.Background(), path)
	if err == nil || !strings.Contains(err.Error(), "5 GiB") || calls.Load() != 0 {
		t.Fatalf("oversize must fail before HTTP: %v, %d requests", err, calls.Load())
	}
}

func TestWaitCancellationAndRejection(t *testing.T) {
	for _, status := range []string{"In Progress", "Invalid", "Rejected", "unexpected"} {
		t.Run(status, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, `{"data":{"id":%q,"attributes":{"status":%q}}}`, testID, status)
			})
			if status == "In Progress" {
				time.AfterFunc(20*time.Millisecond, cancel)
			}
			result, err := client.Wait(ctx, testID, time.Hour)
			if status == "In Progress" {
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("not canceled: %v", err)
				}
				return
			}
			if status == "unexpected" {
				if err == nil {
					t.Fatal("unknown status accepted")
				}
				return
			}
			if err != nil || result.Status != status {
				t.Fatalf("rejected result unavailable for log retrieval: %+v %v", result, err)
			}
		})
	}
}

func TestRedirectAndCredentialErrorsAreNotFollowedOrLeaked(t *testing.T) {
	var calls atomic.Int32
	client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Location", "https://elsewhere.example/?private=token")
		w.WriteHeader(307)
	})
	_, err := client.Status(context.Background(), testID)
	var responseError *HTTPError
	if !errors.As(err, &responseError) || responseError.StatusCode != 307 || calls.Load() != 1 || strings.Contains(err.Error(), "private") {
		t.Fatalf("unsafe redirect handling: %v calls%d", err, calls.Load())
	}
}

func TestStatusHTTPErrorAndMalformedResponse(t *testing.T) {
	for _, body := range []string{"malformed", `{"data":{"id":"wrong","attributes":{"status":"Accepted"}}}`, strings.Repeat("x", maxResponseBytes+1)} {
		t.Run(fmt.Sprint(len(body)), func(t *testing.T) {
			client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, body) })
			if _, err := client.Status(context.Background(), testID); err == nil {
				t.Fatal("invalid response accepted")
			}
		})
	}
	client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(429) })
	_, err := client.Status(context.Background(), testID)
	var responseError *HTTPError
	if !errors.As(err, &responseError) || responseError.StatusCode != 429 {
		t.Fatalf("missing HTTP status: %v", err)
	}
}

func TestAgreementErrorUsesOnlyFixedLocalGuidance(t *testing.T) {
	const code = "FORBIDDEN.REQUIRED_AGREEMENTS_MISSING_OR_EXPIRED"
	const guidance = "notary API: HTTP 403: required Apple agreements are missing or expired; the Account Holder must review the agreements at https://developer.apple.com/account/ before retrying"
	for _, test := range []struct {
		name, operation, code string
		status                int
		alter                 func(string) string
		wantGuidance          bool
	}{
		{name: "recognized", operation: "notary API", code: code, status: 403, wantGuidance: true},
		{name: "unknown", operation: "notary API", code: "FORBIDDEN.UNKNOWN secret-code", status: 403},
		{name: "partial match", operation: "notary API", code: code + ".secret", status: 403},
		{name: "wrong status", operation: "notary API", code: code, status: 500},
		{name: "upload", operation: "notarization upload", code: code, status: 403},
		{name: "log download", operation: "notarization log", code: code, status: 403},
		{name: "malformed", operation: "notary API", code: code, status: 403, alter: func(body string) string { return body + "secret trailing content" }},
		{name: "oversized", operation: "notary API", code: code, status: 403, alter: func(body string) string { return body + strings.Repeat(" ", maxResponseBytes) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) {
				// A hostile body may echo the actual bearer token or suggest a
				// credential-stealing URL. None of these fields may be retained.
				body, err := json.Marshal(map[string]any{"errors": []map[string]string{{
					"code": test.code, "title": "secret-title", "detail": r.Header.Get("Authorization") + " secret-key",
					"url": "https://attacker.invalid/?credential=secret-token",
				}}})
				if err != nil {
					t.Error(err)
					return
				}
				data := string(body)
				if test.alter != nil {
					data = test.alter(data)
				}
				w.WriteHeader(test.status)
				fmt.Fprint(w, data)
			})
			var err error
			if test.operation == "notary API" {
				_, err = client.Status(context.Background(), testID)
			} else {
				request, requestError := http.NewRequest(http.MethodGet, "https://download.invalid/artifact", nil)
				if requestError != nil {
					t.Fatal(requestError)
				}
				_, err = client.response(request, test.operation, maxResponseBytes)
			}
			var responseError *HTTPError
			if !errors.As(err, &responseError) || responseError.StatusCode != test.status || responseError.Operation != test.operation {
				t.Fatalf("HTTP error identity lost: %v", err)
			}
			want := fmt.Sprintf("%s: HTTP %d", test.operation, test.status)
			if test.wantGuidance {
				want = guidance
			}
			if err.Error() != want {
				t.Fatalf("unexpected error text: %q", err.Error())
			}
			if strings.Contains(fmt.Sprintf("%+v", err), "secret") || strings.Contains(fmt.Sprintf("%+v", responseError), "Bearer") {
				t.Fatal("remote error content retained")
			}
		})
	}
}

func TestInvalidIDAndCanceledSubmissionDoNotSendRequests(t *testing.T) {
	var calls atomic.Int32
	client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1) })
	for _, id := range []string{"../logs", strings.Repeat("a", 36), ""} {
		if _, err := client.Status(context.Background(), id); err == nil {
			t.Errorf("accepted invalid ID %q", id)
		}
	}
	path := filepath.Join(t.TempDir(), "app.pkg")
	if err := os.WriteFile(path, []byte("package"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Submit(ctx, path); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("made %d invalid requests", calls.Load())
	}
}

func BenchmarkTokenES256(b *testing.B) {
	credentials, _ := testCredentials(b)
	client, err := New(credentials, nil)
	if err != nil {
		b.Fatal(err)
	}
	now := time.Unix(1700000000, 0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := client.token(now); err != nil {
			b.Fatal(err)
		}
	}
}

func TestTransportFailureOmitsSignedURL(t *testing.T) {
	credentials, _ := testCredentials(t)
	address := "https://logs.example/report?signature=private"
	client, err := New(credentials, &http.Client{Transport: transportFunc(func(*http.Request) (*http.Response, error) {
		return nil, &url.Error{Op: "Get", URL: address, Err: context.DeadlineExceeded}
	})})
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.response(request, "notarization log", maxLogBytes)
	if !errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "private") || strings.Contains(err.Error(), "logs.example") {
		t.Fatalf("transport error leaked signed URL or lost cause: %v", err)
	}
}

func TestStatusRequestCancellation(t *testing.T) {
	client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() })
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := client.Status(ctx, testID); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("request cancellation lost: %v", err)
	}
}

func TestLogRejectsInsecureLocation(t *testing.T) {
	var calls atomic.Int32
	client, _ := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		fmt.Fprintf(w, `{"data":{"id":%q,"attributes":{"developerLogUrl":"http://logs.example/report"}}}`, testID)
	})
	if _, err := client.Log(context.Background(), testID); err == nil || calls.Load() != 1 {
		t.Fatalf("insecure log URL accepted: %v, calls %d", err, calls.Load())
	}
}
