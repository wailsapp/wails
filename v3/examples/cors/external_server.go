// +build ignore

package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed frontend/*
var frontendFiles embed.FS

func main() {
	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get source file location")
	}
	scriptDir := filepath.Dir(filename)

	// Check if certificates exist
	certFile := filepath.Join(scriptDir, "certs", "server.crt")
	keyFile := filepath.Join(scriptDir, "certs", "server.key")

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatalf("Certificate file not found at %s. Please run: go run generate_certs.go", certFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Fatalf("Key file not found at %s. Please run: go run generate_certs.go", keyFile)
	}

	// Create a file server for the frontend
	mux := http.NewServeMux()

	// Serve frontend files
	mux.Handle("/", http.FileServer(http.FS(frontendFiles)))

	// Add some debug endpoints
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","server":"external","cors":"testing"}`))
	})

	// Load certificates
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load certificates: %v", err)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Create HTTPS server
	server := &http.Server{
		Addr:      ":3000",
		Handler:   loggingMiddleware(mux),
		TLSConfig: tlsConfig,
	}

	fmt.Println("🚀 External HTTPS server starting...")
	fmt.Println("📍 URL: https://app-local.wails-awesome.io:3000")
	fmt.Println("📍 Alternative: https://localhost:3000")
	fmt.Println("")
	fmt.Println("⚠️  Make sure you have added to your hosts file:")
	fmt.Println("    127.0.0.1    app-local.wails-awesome.io")
	fmt.Println("")
	fmt.Println("🔒 Make sure you have trusted the CA certificate (see generate_certs.go output)")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")

	// Start HTTPS server
	log.Fatal(server.ListenAndServeTLS("", ""))
}

// loggingMiddleware logs all incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, r.Header.Get("Origin"))

		// Add CORS headers for debugging (the Wails app will handle the actual CORS)
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("X-Debug-Origin", origin)
		}

		next.ServeHTTP(w, r)
	})
}