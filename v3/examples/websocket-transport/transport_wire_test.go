package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type wirePeer struct {
	net.Conn
	reader *bufio.Reader
}

// A small RFC 6455 peer keeps protocol checks independent of the server's library.
func openPeer(tb testing.TB) (*WebSocketTransport, *wirePeer) {
	tb.Helper()
	transport := NewWebSocketTransport(":0")
	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { defer close(done); transport.handleWebSocket(w, r) }))
	transport.server = server.Config
	conn, err := net.Dial("tcp", server.Listener.Addr().String())
	if err != nil {
		server.Close()
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		_ = conn.Close()
		server.Close()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			tb.Error("handler did not finish")
		}
	})
	if err := conn.SetDeadline(time.Now().Add(30 * time.Second)); err != nil {
		tb.Fatal(err)
	}
	_, err = fmt.Fprintf(conn, "GET /wails/ws HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Version: 13\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nOrigin: http://example.invalid\r\n\r\n", server.Listener.Addr())
	if err != nil {
		tb.Fatal(err)
	}
	peer := &wirePeer{Conn: conn, reader: bufio.NewReader(conn)}
	response, err := http.ReadResponse(peer.reader, nil)
	if err != nil {
		tb.Fatal(err)
	}
	if response.StatusCode != http.StatusSwitchingProtocols || response.Header.Get("Sec-WebSocket-Accept") != "s3pPLMBiTxaQ9kYGzzhZRbK+xOo=" {
		tb.Fatalf("handshake rejected: %s", response.Status)
	}
	// A pong proves the handler has registered the client and entered its read loop.
	if err := peer.writeFrame(9, []byte("ready")); err != nil {
		tb.Fatal(err)
	}
	opcode, payload, err := peer.readFrame()
	if err != nil || opcode != 10 || string(payload) != "ready" {
		tb.Fatalf("ping/pong: opcode=%d, payload=%q, error=%v", opcode, payload, err)
	}
	return transport, peer
}

func (p *wirePeer) writeFrame(opcode byte, payload []byte) error {
	header := []byte{0x80 | opcode, 0x80}
	switch {
	case len(payload) < 126:
		header[1] |= byte(len(payload))
	case len(payload) <= 65535:
		header[1] |= 126
		header = binary.BigEndian.AppendUint16(header, uint16(len(payload)))
	default:
		header[1] |= 127
		header = binary.BigEndian.AppendUint64(header, uint64(len(payload)))
	}
	mask := [4]byte{1, 2, 3, 4}
	header = append(header, mask[:]...)
	data := append(header, payload...)
	for i := range payload {
		data[len(header)+i] ^= mask[i%4]
	}
	_, err := p.Write(data)
	return err
}

func (p *wirePeer) readFrame() (byte, []byte, error) {
	var header [2]byte
	if _, err := io.ReadFull(p.reader, header[:]); err != nil {
		return 0, nil, err
	}
	n := uint64(header[1] & 127)
	if n == 126 {
		var size [2]byte
		if _, err := io.ReadFull(p.reader, size[:]); err != nil {
			return 0, nil, err
		}
		n = uint64(binary.BigEndian.Uint16(size[:]))
	}
	if n == 127 {
		var size [8]byte
		if _, err := io.ReadFull(p.reader, size[:]); err != nil {
			return 0, nil, err
		}
		n = binary.BigEndian.Uint64(size[:])
	}
	if header[0]&0x80 == 0 || header[1]&0x80 != 0 || n > 2<<20 {
		return 0, nil, fmt.Errorf("unexpected server frame")
	}
	payload := make([]byte, int(n))
	_, err := io.ReadFull(p.reader, payload)
	return header[0] & 15, payload, err
}

func TestTransportJSONMessages(t *testing.T) {
	transport, peer := openPeer(t)
	transport.DispatchWailsEvent(&application.CustomEvent{Name: "example", Data: "event data"})
	opcode, payload, err := peer.readFrame()
	if err != nil || opcode != 1 {
		t.Fatalf("event frame: %d, %v", opcode, err)
	}
	var message WebSocketMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if message.Type != "event" || message.Event == nil || message.Event.Name != "example" {
		t.Fatalf("event=%+v", message)
	}
	transport.mu.RLock()
	var responses chan *WebSocketMessage
	for _, ch := range transport.clients {
		responses = ch
	}
	transport.mu.RUnlock()
	transport.sendResponse(context.Background(), responses, "request-42", map[string]string{"result": "ok"}, nil)
	_, payload, err = peer.readFrame()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if message.Type != "response" || message.ID != "request-42" || message.Response == nil || message.Response.StatusCode != 200 {
		t.Fatalf("response=%+v", message)
	}
	transport.sendResponse(context.Background(), responses, "failed-42", nil, errors.New("request failed"))
	_, payload, err = peer.readFrame()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatal(err)
	}
	if message.ID != "failed-42" || message.Response == nil || message.Response.StatusCode != 422 || message.Response.Data != "request failed" {
		t.Fatalf("error response=%+v", message)
	}
}

func TestTransportLargeTextAndBinaryJSON(t *testing.T) {
	_, peer := openPeer(t)
	payload := []byte(`{"type":"ignored","padding":"` + strings.Repeat("x", 64<<10) + `"}`)
	for _, opcode := range []byte{1, 2} {
		if err := peer.writeFrame(opcode, payload); err != nil {
			t.Fatal(err)
		}
		if err := peer.writeFrame(9, []byte("after-large")); err != nil {
			t.Fatal(err)
		}
		kind, body, err := peer.readFrame()
		if err != nil || kind != 10 || string(body) != "after-large" {
			t.Fatalf("large message disconnected: kind=%d, payload=%q, error=%v", kind, body, err)
		}
	}
}

func BenchmarkTransportBroadcast(b *testing.B) {
	transport, peer := openPeer(b)
	event := &application.CustomEvent{Name: "bench", Data: strings.Repeat("x", 1024)}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		transport.DispatchWailsEvent(event)
		kind, _, err := peer.readFrame()
		if err != nil || kind != 1 {
			b.Fatalf("broadcast: frame=%d, error=%v", kind, err)
		}
	}
}

func TestTransportStopClosesClients(t *testing.T) {
	transport, peer := openPeer(t)
	if err := transport.Stop(); err != nil {
		t.Fatal(err)
	}
	if _, _, err := peer.readFrame(); err == nil {
		t.Error("client connection remained open")
	}
}
