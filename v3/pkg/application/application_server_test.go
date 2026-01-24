//go:build server

package application

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"
)

// resetGlobalApp resets the global application state for testing
func resetGlobalApp() {
	globalApplication = nil
}

func TestServerMode_HealthEndpoint(t *testing.T) {
	resetGlobalApp()

	// Create a server mode app (server mode is enabled via build tag)
	app := New(Options{
		Name: "Test Server",
		Server: ServerOptions{
			Host: "127.0.0.1",
			Port: 18081, // Use specific port for this test
		},
		Assets: AssetOptions{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
	})

	// Start app in background
	var wg sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		defer wg.Done()
		errCh <- app.Run()
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Test health endpoint
	resp, err := http.Get("http://127.0.0.1:18081/health")
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Shutdown
	app.Quit()

	// Wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("app.Run() returned error: %v", err)
		}
	case <-ctx.Done():
		t.Error("timeout waiting for app shutdown")
	}
}

func TestServerMode_AssetServing(t *testing.T) {
	resetGlobalApp()

	testContent := "Hello from server mode!"

	app := New(Options{
		Name: "Test Assets",
		Server: ServerOptions{
			Host: "127.0.0.1",
			Port: 18082, // Use specific port for this test
		},
		Assets: AssetOptions{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(testContent))
			}),
		},
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	time.Sleep(200 * time.Millisecond)

	// Test asset serving
	resp, err := http.Get("http://127.0.0.1:18082/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	app.Quit()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-errCh:
	case <-ctx.Done():
		t.Error("timeout waiting for app shutdown")
	}
}

func TestServerMode_Defaults(t *testing.T) {
	resetGlobalApp()

	app := New(Options{
		Name: "Test Defaults",
		Server: ServerOptions{
			Port: 18083, // Use specific port to avoid conflicts
		},
		Assets: AssetOptions{
			task: [build:docker] docker build -t server-example:latest -f examples/server/Dockerfile .	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		},
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	time.Sleep(200 * time.Millisecond)

	// Should be listening on localhost:18083
	resp, err := http.Get("http://localhost:18083/health")
	if err != nil {
		t.Fatalf("request to address failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	app.Quit()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-errCh:
	case <-ctx.Done():
		t.Error("timeout waiting for app shutdown")
	}
}
