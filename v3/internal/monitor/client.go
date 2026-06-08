// Package monitor (internal) is the consumer side of the Wails v3 IPC monitor.
// It dials the unix socket exposed by a running app's monitor tap, decodes the
// NDJSON envelope stream (live traces), and can request on-demand snapshots.
package monitor

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"sync"

	mon "github.com/wailsapp/wails/v3/pkg/application/monitor"
)

// Re-exported wire types so the TUI depends on a single package.
type (
	Trace             = mon.Trace
	Snapshot          = mon.Snapshot
	AppInfo           = mon.AppInfo
	WindowInfo        = mon.WindowInfo
	BindingInfo       = mon.BindingInfo
	ParamInfo         = mon.ParamInfo
	EventListenerInfo = mon.EventListenerInfo
	DiscoveryEntry    = mon.DiscoveryEntry
	Sample            = mon.Sample
)

// List returns the live, attachable app instances.
func List() ([]DiscoveryEntry, error) { return mon.ListDiscovery() }

// Client is a live connection to an app's monitor socket.
type Client struct {
	conn    net.Conn
	traces  chan Trace
	snaps   chan *Snapshot
	samples chan Sample
	errc    chan error

	wmu sync.Mutex // serializes writes (describe requests)
}

// Connect dials the socket and starts decoding the envelope stream. Traces are
// delivered on Traces(); snapshot replies (from Describe) on Snapshots(). The
// reader stops when ctx is cancelled or the connection closes.
func Connect(ctx context.Context, sockPath string) (*Client, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "unix", sockPath)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:    conn,
		traces:  make(chan Trace, 256),
		snaps:   make(chan *Snapshot, 4),
		samples: make(chan Sample, 256),
		errc:    make(chan error, 1),
	}

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	go func() {
		defer close(c.traces)
		sc := bufio.NewScanner(conn)
		sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
		for sc.Scan() {
			line := sc.Bytes()
			if len(line) == 0 {
				continue
			}
			var env mon.Envelope
			if err := json.Unmarshal(line, &env); err != nil {
				continue
			}
			switch env.Type {
			case mon.MsgTrace:
				if env.Trace != nil {
					select {
					case c.traces <- *env.Trace:
					case <-ctx.Done():
						return
					}
				}
			case mon.MsgSnapshot:
				if env.Snapshot != nil {
					select {
					case c.snaps <- env.Snapshot:
					default: // drop if no one is waiting
					}
				}
			case mon.MsgSample:
				if env.Sample != nil {
					select {
					case c.samples <- *env.Sample:
					default: // drop if consumer is behind
					}
				}
			}
		}
		if err := sc.Err(); err != nil && ctx.Err() == nil {
			select {
			case c.errc <- err:
			default:
			}
		}
	}()

	return c, nil
}

// Traces is the live trace stream; closed when the connection ends.
func (c *Client) Traces() <-chan Trace { return c.traces }

// Snapshots delivers snapshot replies to Describe requests.
func (c *Client) Snapshots() <-chan *Snapshot { return c.snaps }

// Samples delivers periodic resource-usage samples.
func (c *Client) Samples() <-chan Sample { return c.samples }

// Errors delivers a terminal stream error, if any.
func (c *Client) Errors() <-chan error { return c.errc }

// Describe sends a snapshot request. The reply arrives on Snapshots().
func (c *Client) Describe() error {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	_, err := c.conn.Write([]byte(`{"req":"describe"}` + "\n"))
	return err
}

// Close closes the connection.
func (c *Client) Close() error { return c.conn.Close() }
