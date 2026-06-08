package monitor

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	mon "github.com/wailsapp/wails/v3/pkg/application/monitor"
)

// TestClientEndToEnd exercises the real producer Sink and the consumer Client:
// a trace emitted by the app is decoded by the client, and a Describe request
// returns the registered snapshot.
func TestClientEndToEnd(t *testing.T) {
	sock := filepath.Join(t.TempDir(), "e2e.sock")
	s, err := mon.Start(mon.Config{SocketPath: sock, RingSize: 16})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer s.Stop()

	mon.SetDescribeFunc(func() *mon.Snapshot {
		return &mon.Snapshot{
			App:      mon.AppInfo{Name: "e2e", PID: 99},
			Windows:  []mon.WindowInfo{{ID: 1, Name: "main", Width: 800, Height: 600, Focused: true}},
			Bindings: []mon.BindingInfo{{FQN: "main.Greet.Greet", Name: "Greet"}},
		}
	})
	defer mon.SetDescribeFunc(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	c, err := Connect(ctx, sock)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer c.Close()

	mon.Emit(mon.Trace{Kind: "call", Dir: "in", CallID: "x1", Method: "main.Greet.Greet"})

	select {
	case tr := <-c.Traces():
		if tr.CallID != "x1" || tr.Method != "main.Greet.Greet" {
			t.Fatalf("unexpected trace: %+v", tr)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for trace")
	}

	if err := c.Describe(); err != nil {
		t.Fatalf("Describe: %v", err)
	}
	select {
	case snap := <-c.Snapshots():
		if snap == nil || snap.App.Name != "e2e" {
			t.Fatalf("bad snapshot: %+v", snap)
		}
		if len(snap.Windows) != 1 || snap.Windows[0].Name != "main" {
			t.Fatalf("bad windows: %+v", snap.Windows)
		}
		if len(snap.Bindings) != 1 || snap.Bindings[0].Name != "Greet" {
			t.Fatalf("bad bindings: %+v", snap.Bindings)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for snapshot")
	}
}
