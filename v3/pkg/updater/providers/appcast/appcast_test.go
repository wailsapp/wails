package appcast_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/appcast"
)

const feedTemplate = `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle">
  <channel>
    <title>App Updates</title>
    <item>
      <title>1.0.0</title>
      <pubDate>Sat, 01 Mar 2026 10:00:00 +0000</pubDate>
      <sparkle:version>1.0.0</sparkle:version>
      <sparkle:shortVersionString>1.0.0</sparkle:shortVersionString>
      <enclosure url="%s/dl/old-darwin.dmg" length="100" type="application/octet-stream" sparkle:os="macos"/>
    </item>
    <item>
      <title>2.0.0</title>
      <description><![CDATA[<p>Big upgrade.</p>]]></description>
      <pubDate>Wed, 01 Apr 2026 10:00:00 +0000</pubDate>
      <sparkle:version>2.0.0</sparkle:version>
      <sparkle:shortVersionString>2.0.0</sparkle:shortVersionString>
      <sparkle:os>macos</sparkle:os>
      <enclosure url="%s/dl/new-darwin.dmg" length="12345" type="application/octet-stream" sparkle:edSignature="dGVzdC1zaWc"/>
    </item>
    <item>
      <title>2.0.0-windows</title>
      <pubDate>Wed, 01 Apr 2026 10:00:00 +0000</pubDate>
      <sparkle:version>2.0.0</sparkle:version>
      <sparkle:shortVersionString>2.0.0</sparkle:shortVersionString>
      <sparkle:os>windows</sparkle:os>
      <enclosure url="%s/dl/new-win.exe" length="22222" type="application/octet-stream"/>
    </item>
  </channel>
</rss>`

func TestCheck_PicksLatestMatchingPlatform(t *testing.T) {
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, feedTemplate, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL

	p, err := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml"})
	if err != nil {
		t.Fatal(err)
	}
	rel, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "darwin",
		Arch:           "arm64",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil {
		t.Fatal("expected release")
	}
	if rel.Version != "2.0.0" {
		t.Errorf("version: %q", rel.Version)
	}
	if rel.Artifact.Filename != "new-darwin.dmg" {
		t.Errorf("filename: %q", rel.Artifact.Filename)
	}
	if rel.Artifact.Size != 12345 {
		t.Errorf("size: %d", rel.Artifact.Size)
	}
	// Signature should be decoded from base64 raw-stdenc.
	if rel.Verification == nil {
		t.Fatal("expected verification populated from sparkle:edSignature")
	}
	if rel.Verification.SignatureAlgo != "ed25519" {
		t.Errorf("algo: %s", rel.Verification.SignatureAlgo)
	}
	expSig, _ := base64.RawStdEncoding.DecodeString("dGVzdC1zaWc")
	if !bytes.Equal(rel.Verification.Signature, expSig) {
		t.Errorf("signature: %x", rel.Verification.Signature)
	}
}

func TestCheck_PicksWindowsItem(t *testing.T) {
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, feedTemplate, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL
	p, _ := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml"})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0", Platform: "windows", Arch: "amd64"})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Artifact.Filename != "new-win.exe" {
		t.Fatalf("got %+v", rel)
	}
}

func TestCheck_UpToDate_NoNewerVersion(t *testing.T) {
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, feedTemplate, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL
	p, _ := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml"})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "2.0.0", Platform: "darwin", Arch: "arm64"})
	if err != nil {
		t.Fatal(err)
	}
	if rel != nil {
		t.Errorf("expected nil release, got %+v", rel)
	}
}

func TestCheck_NoItemsForPlatform(t *testing.T) {
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, feedTemplate, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL
	p, _ := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml"})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0", Platform: "linux", Arch: "amd64"})
	if err != nil {
		t.Fatal(err)
	}
	if rel != nil {
		t.Errorf("expected nil release for unsupported platform, got %+v", rel)
	}
}

func TestCheck_FeedNotFound_TreatedAsUpToDate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	p, _ := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml"})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err != nil {
		t.Fatal(err)
	}
	if rel != nil {
		t.Errorf("404 should map to nil, got %+v", rel)
	}
}

func TestCheck_FeedError_Surfaced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	p, _ := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml"})
	_, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("expected 500 error, got %v", err)
	}
}

func TestCheck_ChannelFilter(t *testing.T) {
	feed := `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle">
  <channel>
    <item><title>3.0.0</title>
      <sparkle:version>3.0.0</sparkle:version>
      <sparkle:shortVersionString>3.0.0</sparkle:shortVersionString>
      <sparkle:channel>beta</sparkle:channel>
      <enclosure url="EXAMPLE/3.dmg" length="1" type="x"/>
    </item>
    <item><title>2.5.0</title>
      <sparkle:version>2.5.0</sparkle:version>
      <sparkle:shortVersionString>2.5.0</sparkle:shortVersionString>
      <sparkle:channel>stable</sparkle:channel>
      <enclosure url="EXAMPLE/2.dmg" length="1" type="x"/>
    </item>
  </channel>
</rss>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, strings.ReplaceAll(feed, "EXAMPLE", ""))
	}))
	defer srv.Close()
	p, _ := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml", Channel: "stable"})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Version != "2.5.0" {
		t.Fatalf("channel filter ignored: %+v", rel)
	}
}

// An explicit Config.Channel must not match items that ship without a
// sparkle:channel — Sparkle's contract is that items opt *in* to a channel,
// so an unlabelled item should never satisfy a requested channel filter.
func TestCheck_ChannelFilter_RejectsUnchannelledItems(t *testing.T) {
	feed := `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle">
  <channel>
    <item><title>3.0.0</title>
      <sparkle:version>3.0.0</sparkle:version>
      <sparkle:shortVersionString>3.0.0</sparkle:shortVersionString>
      <enclosure url="EXAMPLE/3.dmg" length="1" type="x"/>
    </item>
    <item><title>2.5.0</title>
      <sparkle:version>2.5.0</sparkle:version>
      <sparkle:shortVersionString>2.5.0</sparkle:shortVersionString>
      <sparkle:channel>stable</sparkle:channel>
      <enclosure url="EXAMPLE/2.dmg" length="1" type="x"/>
    </item>
  </channel>
</rss>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, strings.ReplaceAll(feed, "EXAMPLE", ""))
	}))
	defer srv.Close()
	p, _ := appcast.New(appcast.Config{URL: srv.URL + "/appcast.xml", Channel: "stable"})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Version != "2.5.0" {
		t.Fatalf("channel filter should have skipped the unchannelled 3.0.0 item: %+v", rel)
	}
}

func TestDownload_StreamsAndReportsProgress(t *testing.T) {
	body := []byte("dmg-contents")
	var hits int32
	dl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.Write(body)
	}))
	defer dl.Close()

	p, _ := appcast.New(appcast.Config{URL: "irrelevant"})
	var buf bytes.Buffer
	var ticks int32
	err := p.Download(context.Background(), &updater.Release{
		Artifact: updater.Artifact{Filename: "x.dmg", Size: int64(len(body))},
		Metadata: map[string]any{"appcast.enclosure.url": dl.URL + "/asset.dmg"},
	}, &buf, func(_, _ int64) { atomic.AddInt32(&ticks, 1) })
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), body) {
		t.Errorf("body mismatch: %q", buf.String())
	}
	if atomic.LoadInt32(&hits) != 1 {
		t.Errorf("hits: %d", hits)
	}
	if atomic.LoadInt32(&ticks) == 0 {
		t.Error("expected progress callbacks")
	}
}

func TestNew_RequiresURL(t *testing.T) {
	if _, err := appcast.New(appcast.Config{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestProviderInterface(t *testing.T) {
	var _ updater.Provider = (*appcast.Provider)(nil)
}
