package main

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestBinaryResponseDoesNotRenderHTML(t *testing.T) {
	for _, payload := range [][]byte{
		[]byte(`<html><script>alert("echo")</script></html>`),
		{0, 1, 128, 255, 10},
	} {
		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "text/html")
		writeBinaryResponse(response, bytes.NewReader(payload))
		if got := response.Header().Get("Content-Type"); got != "application/octet-stream" {
			t.Errorf("unsafe content type %q", got)
		}
		if got := response.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Errorf("missing sniff protection: %q", got)
		}
		if !bytes.Equal(response.Body.Bytes(), payload) {
			t.Errorf("binary echo changed: %x", response.Body.Bytes())
		}
	}
}
