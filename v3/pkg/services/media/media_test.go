package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"sync/atomic"
	"testing"
	"testing/fstest"
)

type testConn struct {
	ctx     context.Context
	request []byte
	frames  [][]byte
	onSend  func([]byte) error
}

func (c *testConn) Context() context.Context { return c.ctx }
func (c *testConn) Receive() ([]byte, error) { return c.request, nil }
func (c *testConn) Send(b []byte) error {
	if c.onSend != nil {
		if err := c.onSend(b); err != nil {
			return err
		}
	}
	c.frames = append(c.frames, b)
	return nil
}
func newConn(name string, limit int64) *testConn {
	b, _ := json.Marshal(request{Version: 1, Name: name, MaxBytes: limit})
	return &testConn{ctx: context.Background(), request: b}
}
func header(t *testing.T, c *testConn) metadata {
	t.Helper()
	if len(c.frames) == 0 {
		t.Fatal("missing metadata")
	}
	var m metadata
	if err := json.Unmarshal(c.frames[0], &m); err != nil {
		t.Fatal(err)
	}
	return m
}

type observedFS struct {
	fs.FS
	read   atomic.Int64
	closed atomic.Int64
}
type observedFile struct {
	fs.File
	owner *observedFS
}

func (f *observedFS) Open(name string) (fs.File, error) {
	v, err := f.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return &observedFile{v, f}, nil
}
func (f *observedFile) Read(b []byte) (int, error) {
	n, err := f.File.Read(b)
	f.owner.read.Add(int64(n))
	return n, err
}
func (f *observedFile) Close() error { f.owner.closed.Add(1); return f.File.Close() }

func TestTransferChunksPreserveBytes(t *testing.T) {
	data := make([]byte, 3*chunkBytes+7)
	for i := range data {
		data[i] = byte(i / 37)
	}
	files := &observedFS{FS: fstest.MapFS{"clip.wav": {Data: data}}}
	c := newConn("clip.wav", int64(len(data)))
	transfer(c, files, int64(len(data)))
	m := header(t, c)
	if m.Version != 1 || m.Error != "" || m.Size != int64(len(data)) || m.Type == "" {
		t.Fatal(m)
	}
	var got []byte
	for _, chunk := range c.frames[1:] {
		if len(chunk) > chunkBytes {
			t.Fatalf("oversized frame %d", len(chunk))
		}
		got = append(got, chunk...)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("queued frames were corrupted")
	}
	if files.closed.Load() != 1 {
		t.Fatal("file not closed")
	}
}
func TestLimitsRejectBeforeReading(t *testing.T) {
	for _, tc := range []struct{ client, server int64 }{{3, 8}, {8, 3}} {
		files := &observedFS{FS: fstest.MapFS{"clip.wav": {Data: []byte("1234")}}}
		c := newConn("clip.wav", tc.client)
		transfer(c, files, tc.server)
		m := header(t, c)
		if !m.Limit || m.Error == "" || len(c.frames) != 1 {
			t.Fatal(m)
		}
		if files.read.Load() != 0 || files.closed.Load() != 1 {
			t.Fatal("oversized file read or left open")
		}
	}
}
func TestInvalidRequests(t *testing.T) {
	files := fstest.MapFS{"directory/clip.wav": {Data: []byte("data")}}
	for _, name := range []string{".", "..", "../secret", "/secret", `dir\secret`, "file:secret", "missing", "directory"} {
		c := newConn(name, 100)
		transfer(c, files, 100)
		if header(t, c).Error == "" || len(c.frames) != 1 {
			t.Fatalf("accepted %s", name)
		}
	}
	for _, frame := range [][]byte{[]byte("not json"), bytes.Repeat([]byte("x"), 4097), []byte(`{"version":2,"name":"x","maxBytes":3}`), []byte(`{"version":1,"name":"x","maxBytes":0}`)} {
		c := &testConn{ctx: context.Background(), request: frame}
		transfer(c, files, 100)
		if header(t, c).Error == "" {
			t.Fatal("invalid request accepted")
		}
	}
	if _, err := NewHandler(nil, 100); err == nil {
		t.Fatal("nil FS accepted")
	}
	if _, err := NewHandler(files, 0); err == nil {
		t.Fatal("zero limit accepted")
	}
}
func TestCancellationAndSendFailureStopReading(t *testing.T) {
	for _, fail := range []bool{false, true} {
		files := &observedFS{FS: fstest.MapFS{"clip.wav": {Data: make([]byte, 4*chunkBytes)}}}
		c := newConn("clip.wav", 4*chunkBytes)
		ctx, cancel := context.WithCancel(context.Background())
		c.ctx = ctx
		c.onSend = func(b []byte) error {
			if len(b) == chunkBytes {
				if fail {
					return errors.New("closed")
				}
				cancel()
			}
			return nil
		}
		transfer(c, files, 4*chunkBytes)
		cancel()
		if files.read.Load() != chunkBytes || files.closed.Load() != 1 {
			t.Fatalf("continued reading or leaked file: read=%d closed=%d", files.read.Load(), files.closed.Load())
		}
	}
}
