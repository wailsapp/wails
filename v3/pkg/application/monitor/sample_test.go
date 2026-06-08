package monitor

import (
	"bufio"
	"encoding/json"
	"net"
	"path/filepath"
	"testing"
	"time"
)

// TestReadProcOnHost is a smoke test for the OS probe. On Linux it must return
// a non-zero RSS for the running test binary; elsewhere it is allowed to be a
// no-op (the sampler degrades to runtime-only stats).
func TestReadProc(t *testing.T) {
	p, ok := readProc()
	if !ok {
		t.Skip("readProc unavailable on this platform")
	}
	if p.RSS == 0 {
		t.Errorf("expected non-zero RSS, got 0")
	}
}

// TestSamplerEmits starts a sink with a short sampler interval and asserts a
// MsgSample envelope arrives over the socket with sane runtime fields.
func TestSamplerEmits(t *testing.T) {
	sock := filepath.Join(t.TempDir(), "s.sock")
	s, err := Start(Config{SocketPath: sock})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Stop()

	// Replace the 1s production sampler with a fast one for the test.
	s.sampler.stop()
	s.sampler = startSampler(s, 20*time.Millisecond)

	conn, err := net.Dial("unix", sock)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		var env Envelope
		if err := json.Unmarshal(scanner.Bytes(), &env); err != nil {
			continue
		}
		if env.Type != MsgSample || env.Sample == nil {
			continue
		}
		if env.Sample.Goroutines <= 0 {
			t.Errorf("expected goroutines > 0, got %d", env.Sample.Goroutines)
		}
		if env.Sample.HeapAlloc == 0 {
			t.Errorf("expected non-zero heap alloc")
		}
		return // got a sample, done
	}
	t.Fatalf("no MsgSample received: %v", scanner.Err())
}
