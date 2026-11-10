//go:build server && !windows

package application

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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

// freePort is a port nothing listens on, so the readiness check cannot be answered by another process.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// The default signal handler is off, so only serverApp.run's own shutdown path can run the hooks.
func TestServerMode_SignalRunsShutdownHooks(t *testing.T) {
	resetGlobalApp()

	// While any channel is registered SIGTERM is delivered rather than exiting, so a signal that misses the server
	// fails this test instead of ending the test binary.
	guard := make(chan os.Signal, 1)
	signal.Notify(guard, syscall.SIGTERM)
	defer signal.Stop(guard)

	port := freePort(t)
	ran := make(chan struct{}, 1)
	app := New(Options{
		Name:                        "Test Server",
		Server:                      ServerOptions{Host: "127.0.0.1", Port: port},
		DisableDefaultSignalHandler: true,
		OnShutdown:                  func() { ran <- struct{}{} },
	})
	errCh := make(chan error, 1)
	go func() { errCh <- app.Run() }()
	waitForListener(t, fmt.Sprintf("127.0.0.1:%d", port))
	select {
	case err := <-errCh:
		t.Fatalf("app.Run() returned before the signal: %v", err)
	default:
	}

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
