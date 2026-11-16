package main

import (
	"io"
	"net/http"
)

// Binary round trips preserve arbitrary payloads without treating them as HTML.
func writeBinaryResponse(w http.ResponseWriter, body io.Reader) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = io.Copy(w, body)
}
