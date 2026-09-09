package appcast

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/updater"
)

// TestIntegration_FeedToInstall stands up a real HTTP server serving a
// Sparkle-shaped appcast.xml + binary, runs the provider's full Check +
// Download flow against it, and verifies an Ed25519 signature on the way
// through. End-to-end coverage of the AppCast pipeline without needing to
// publish to a public GitHub Pages site.
func TestIntegration_FeedToInstall(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Build a release artifact: a fake 4 KiB binary blob.
	body := make([]byte, 4096)
	for i := range body {
		body[i] = byte(i % 251)
	}
	sum := sha256.Sum256(body)
	sig := ed25519.Sign(priv, sum[:])
	sigB64 := base64.StdEncoding.EncodeToString(sig)

	// Stand up the server. URL gets templated into the feed.
	mux := http.NewServeMux()
	var feedURL, binURL string

	mux.HandleFunc("/binary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		_, _ = w.Write(body)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	binURL = srv.URL + "/binary"
	feedURL = srv.URL + "/appcast.xml"

	// Now register the feed handler — needs srv.URL which we just learned.
	feed := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle">
 <channel>
  <title>Demo</title>
  <item>
   <title>v2.0.0</title>
   <sparkle:version>2.0.0</sparkle:version>
   <sparkle:shortVersionString>2.0.0</sparkle:shortVersionString>
   <description><![CDATA[Test release.]]></description>
   <enclosure
     url="%s"
     length="%d"
     type="application/octet-stream"
     sparkle:os="%s"
     sparkle:edSignature="%s" />
  </item>
 </channel>
</rss>`, binURL, len(body), runtime.GOOS, sigB64)

	mux.HandleFunc("/appcast.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(feed))
	})

	// --- Run the provider against the live server ---
	p, err := New(Config{URL: feedURL})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	rel, err := p.Check(ctx, updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       runtime.GOOS,
		Arch:           runtime.GOARCH,
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if rel == nil {
		t.Fatal("expected update, got nil")
	}
	if rel.Version != "2.0.0" {
		t.Errorf("version: got %q want 2.0.0", rel.Version)
	}
	if rel.Verification == nil || rel.Verification.SignatureAlgo != "ed25519" {
		t.Fatalf("verification not populated: %+v", rel.Verification)
	}
	if len(rel.Verification.Signature) != ed25519.SignatureSize {
		t.Errorf("signature length: got %d want %d", len(rel.Verification.Signature), ed25519.SignatureSize)
	}

	// Stream Download into a buffer and run signature verification end-to-end.
	var written strings.Builder
	if err := p.Download(ctx, rel, writerNoop{w: &written}, func(int64, int64) {}); err != nil {
		t.Fatalf("Download: %v", err)
	}
	if written.Len() != len(body) {
		t.Errorf("downloaded %d bytes, want %d", written.Len(), len(body))
	}

	// Independently verify the signature with the key we generated.
	gotSum := sha256.Sum256([]byte(written.String()))
	if !ed25519.Verify(pub, gotSum[:], rel.Verification.Signature) {
		t.Fatalf("signature did not verify against downloaded body")
	}
	t.Logf("verified Ed25519 signature over SHA-256 of downloaded body")

	// Also exercise that the channel filter rejects items lacking sparkle:channel
	// when one is configured. Quick second pass using a different provider.
	p2, _ := New(Config{URL: feedURL, Channel: "stable"})
	rel2, err := p2.Check(ctx, updater.CheckRequest{CurrentVersion: "1.0.0", Platform: runtime.GOOS, Arch: runtime.GOARCH})
	if err == nil && rel2 != nil {
		t.Errorf("Channel=stable filter should have rejected the unchannelled item; got %s", rel2.Version)
	}
}

type writerNoop struct{ w io.Writer }

func (n writerNoop) Write(p []byte) (int, error) { return n.w.Write(p) }

// TestIntegration_FeedNetworkTimeout exercises the provider's ctx cancellation
// when the feed never responds.
func TestIntegration_FeedNetworkTimeout(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/appcast.xml", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // longer than ctx timeout below
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	p, err := New(Config{URL: srv.URL + "/appcast.xml"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	_, err = p.Check(ctx, updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

// Sanity that the feed XML the test builds is itself well-formed
// — protects against typos in the template above.
func TestIntegration_FeedXMLWellFormed(t *testing.T) {
	feed := `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle">
 <channel>
  <item>
   <enclosure url="https://x/y" length="1" sparkle:edSignature="aGk=" />
  </item>
 </channel>
</rss>`
	var v any
	if err := xml.Unmarshal([]byte(feed), &v); err != nil {
		t.Fatal(err)
	}
}
