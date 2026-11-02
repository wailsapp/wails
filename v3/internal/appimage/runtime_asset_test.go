package appimage

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRuntimeAssetOverridesValidateTarget(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "runtime")
	if err := os.WriteFile(file, minimalELF("amd64", true), 0755); err != nil {
		t.Fatal(err)
	}
	options := Options{Root: root, RuntimePath: "runtime", Arch: "amd64"}
	got, err := options.PrepareRuntime(t.Context())
	if err != nil || got != file {
		t.Fatalf("runtime=%s err=%v", got, err)
	}
	for _, arch := range []string{"arm64", "unknown"} {
		options.Arch = arch
		if _, err := options.PrepareRuntime(t.Context()); err == nil {
			t.Fatalf("accepted %s runtime", arch)
		}
	}
}

func TestRuntimeDownloadRejectsUnpinnedBytes(t *testing.T) {
	options := Options{Root: t.TempDir(), Arch: "amd64", Client: &http.Client{Transport: launcherDownloadTransport(func(request *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(request.URL.Path, "/657374641") || request.Header.Get("Accept") != "application/octet-stream" {
			t.Errorf("unexpected runtime request: %s, %v", request.URL, request.Header)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("untrusted bytes"))}, nil
	})}}
	if _, err := options.PrepareRuntime(context.Background()); err == nil || !strings.Contains(err.Error(), "SHA256") {
		t.Fatalf("unexpected error: %v", err)
	}
	cache := filepath.Join(options.Root, ".wake/appimage-tools/amd64/runtime-"+runtimeAssets["amd64"].sha256)
	if _, err := os.Stat(cache); !os.IsNotExist(err) {
		t.Fatalf("invalid runtime retained: %v", err)
	}
}
