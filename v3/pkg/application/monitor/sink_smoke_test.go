package monitor

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"path/filepath"
	"testing"
	"time"
)

// TestSinkRoundTrip verifies a trace emitted into the sink is delivered to a
// connected client as an NDJSON envelope, and that the replay ring works.
func TestSinkRoundTrip(t *testing.T) {
	sock := filepath.Join(t.TempDir(), "test.sock")
	s, err := Start(Config{SocketPath: sock, RingSize: 8})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Stop()

	if !Enabled() {
		t.Fatal("Enabled() should be true after Start")
	}

	// Emit one trace BEFORE connecting → must be replayed from the ring.
	Emit(Trace{Kind: "call", Dir: "in", CallID: "c1", Method: "main.Greet"})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", sock)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)

	// First line: the replayed call.
	env := readEnvelope(t, sc)
	if env.Type != MsgTrace || env.Trace == nil {
		t.Fatalf("expected trace envelope, got %+v", env)
	}
	if env.Trace.CallID != "c1" || env.Trace.Kind != "call" {
		t.Fatalf("replay mismatch: %+v", env.Trace)
	}
	if env.Trace.Seq == 0 || env.Trace.Time.IsZero() {
		t.Fatalf("seq/time not stamped: %+v", env.Trace)
	}

	// Emit a live result → must stream.
	Emit(Trace{Kind: "result", Dir: "in", CallID: "c1", DurationMS: 1.5})
	env = readEnvelope(t, sc)
	if env.Trace == nil || env.Trace.Kind != "result" || env.Trace.CallID != "c1" {
		t.Fatalf("live mismatch: %+v", env)
	}
}

// TestDescribeRequest verifies a describe request returns a snapshot envelope.
func TestDescribeRequest(t *testing.T) {
	sock := filepath.Join(t.TempDir(), "desc.sock")
	s, err := Start(Config{SocketPath: sock, RingSize: 8})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Stop()

	SetDescribeFunc(func() *Snapshot {
		return &Snapshot{App: AppInfo{Name: "testapp", PID: 1234}}
	})
	defer SetDescribeFunc(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", sock)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(`{"req":"describe"}` + "\n")); err != nil {
		t.Fatalf("write request: %v", err)
	}

	sc := bufio.NewScanner(conn)
	for {
		env := readEnvelope(t, sc)
		if env.Type == MsgSnapshot {
			if env.Snapshot == nil || env.Snapshot.App.Name != "testapp" {
				t.Fatalf("bad snapshot: %+v", env.Snapshot)
			}
			return
		}
		// ignore any trace envelopes that arrive first
	}
}

func TestDisabledEmitIsNoop(t *testing.T) {
	Emit(Trace{Kind: "event", Method: "noop"})
}

func readEnvelope(t *testing.T, sc *bufio.Scanner) Envelope {
	t.Helper()
	if !sc.Scan() {
		t.Fatalf("scan failed: %v", sc.Err())
	}
	var e Envelope
	if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal %q: %v", sc.Text(), err)
	}
	return e
}
