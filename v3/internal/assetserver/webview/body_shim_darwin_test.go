//go:build darwin && !ios

package webview

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestUnshimHeaderRestoresTheOriginalRequest(t *testing.T) {
	tests := []struct {
		name     string
		header   map[string]string
		shimmed  bool
		wantType string // "" = no Content-Type
	}{
		{"not shimmed", map[string]string{"Content-Type": "text/plain"}, false, "text/plain"},
		{"binary with no type", map[string]string{bodyEncodingHeader: "base64", bodyTypeHeader: "", "Content-Type": "text/plain", "Content-Length": "8"}, true, ""},
		{"blob with a type", map[string]string{bodyEncodingHeader: "base64", bodyTypeHeader: "image/png", "Content-Type": "text/plain"}, true, "image/png"},
		{"form data", map[string]string{bodyEncodingHeader: "base64", bodyTypeHeader: "multipart/form-data; boundary=x"}, true, "multipart/form-data; boundary=x"},
		{"other encoding ignored", map[string]string{bodyEncodingHeader: "gzip", "Content-Type": "text/plain"}, false, "text/plain"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			for k, v := range tt.header {
				h.Set(k, v)
			}
			if got := unshimHeader(h); got != tt.shimmed {
				t.Fatalf("unshimHeader = %v, want %v", got, tt.shimmed)
			}
			if got := h.Get("Content-Type"); got != tt.wantType {
				t.Errorf("Content-Type = %q, want %q", got, tt.wantType)
			}
			if tt.shimmed {
				for _, k := range []string{bodyEncodingHeader, bodyTypeHeader, "Content-Length"} {
					if h.Get(k) != "" {
						t.Errorf("%s left in the header", k)
					}
				}
			}
		})
	}
}

func TestUnshimBodyDecodes(t *testing.T) {
	body := unshimBody(io.NopCloser(strings.NewReader("AAGA/wo="))) // 0 1 128 255 10
	got, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0, 1, 128, 255, 10}; string(got) != string(want) {
		t.Fatalf("decoded %v, want %v", got, want)
	}
	if unshimBody(http.NoBody) != http.NoBody {
		t.Error("an empty body must stay http.NoBody")
	}
}
