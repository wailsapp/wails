//go:build server && !windows

package application

import (
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

// waitForListener returns once addr accepts connections; run registers its signal channel before listening.
func waitForListener(t *testing.T, addr string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if conn, err := net.Dial("tcp", addr); err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server did not listen on %s", addr)
}

// The default signal handler is off, so only serverApp.run's own shutdown path can run the hooks.
func TestServerMode_SignalRunsShutdownHooks(t *testing.T) {
	resetGlobalApp()

	ran := make(chan struct{}, 1)
	app := New(Options{
		Name:                        "Test Server",
		Server:                      ServerOptions{Host: "127.0.0.1", Port: 18089},
		DisableDefaultSignalHandler: true,
		OnShutdown:                  func() { ran <- struct{}{} },
	})
	errCh := make(chan error, 1)
	go func() { errCh <- app.Run() }()
	waitForListener(t, "127.0.0.1:18089")

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("app.Run() returned error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for shutdown")
	}
	select {
	case <-ran:
	default:
		t.Fatal("OnShutdown did not run after SIGTERM")
	}
}
