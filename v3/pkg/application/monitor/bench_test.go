package monitor

import (
	"bufio"
	"encoding/json"
	"net"
	"path/filepath"
	"sync/atomic"
	"testing"
)

// sampleTrace mirrors a realistic binding-call trace: a method name, a window,
// and a small JSON argument payload (already a RawMessage on the hot path).
func sampleTrace() Trace {
	return Trace{
		Kind:       "call",
		Dir:        "in",
		CallID:     "V1StGXR8_Z5jdHi6B-myT",
		ObjectName: "call",
		Method:     "main.GreetService.Greet",
		Window:     "main",
		Args:       json.RawMessage(`["Alice",42,{"verbose":true}]`),
	}
}

// BenchmarkEmitDisabled is the production-default path: the monitor is off, so
// Emit must do nothing but an atomic load + nil check. This is the cost paid by
// every shipped app on every IPC message.
func BenchmarkEmitDisabled(b *testing.B) {
	sink.Store(nil) // ensure disabled
	t := sampleTrace()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Emit(t)
	}
}

// BenchmarkEnabledGate measures the call-site guard (`if monitor.Enabled()`)
// that wraps the hot-path taps when disabled.
func BenchmarkEnabledGate(b *testing.B) {
	sink.Store(nil)
	b.ReportAllocs()
	b.ResetTimer()
	var on bool
	for i := 0; i < b.N; i++ {
		on = Enabled()
	}
	_ = on
}

// BenchmarkEmitEnabledNoClients isolates the marshal + ring-store cost when the
// monitor is running but no consumer is attached (the common dev case: app up,
// TUI not yet launched).
func BenchmarkEmitEnabledNoClients(b *testing.B) {
	s := &Sink{
		ringSize:     4096,
		clientBuffer: 1024,
		clients:      make(map[*client]struct{}),
	}
	sink.Store(s)
	defer sink.Store(nil)

	t := sampleTrace()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Emit(t)
	}
}

// BenchmarkEmitEnabledWithClient measures the full hot path with one attached,
// continuously draining client (marshal + ring-store + non-blocking send).
func BenchmarkEmitEnabledWithClient(b *testing.B) {
	c := &client{ch: make(chan []byte, 1024)}
	s := &Sink{
		ringSize:     4096,
		clientBuffer: 1024,
		clients:      map[*client]struct{}{c: {}},
	}
	sink.Store(s)
	defer sink.Store(nil)

	done := make(chan struct{})
	var drained atomic.Uint64
	go func() {
		for range c.ch {
			drained.Add(1)
		}
		close(done)
	}()

	t := sampleTrace()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Emit(t)
	}
	b.StopTimer()
	close(c.ch)
	<-done
}

// BenchmarkEmitEndToEnd measures Emit reaching a real client over a unix socket
// (marshal + ring + chan + conn.Write), the true wire cost per traced message.
func BenchmarkEmitEndToEnd(b *testing.B) {
	sock := filepath.Join(b.TempDir(), "b.sock")
	s, err := Start(Config{SocketPath: sock})
	if err != nil {
		b.Fatalf("Start: %v", err)
	}
	defer s.Stop()
	s.sampler.stop() // exclude the 1 Hz sampler from the measurement

	conn, err := net.Dial("unix", sock)
	if err != nil {
		b.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		r := bufio.NewReader(conn)
		buf := make([]byte, 32*1024)
		for {
			if _, err := r.Read(buf); err != nil {
				close(done)
				return
			}
		}
	}()

	t := sampleTrace()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Emit(t)
	}
	b.StopTimer()
}

// BenchmarkReadProc measures one OS resource probe (the per-tick sampler cost,
// minus ReadMemStats). Skips on platforms without a /proc implementation.
func BenchmarkReadProc(b *testing.B) {
	if _, ok := readProc(); !ok {
		b.Skip("readProc unavailable on this platform")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = readProc()
	}
}
