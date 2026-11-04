//go:build server

package application

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestIsLoopbackBindHost(t *testing.T) {
	tests := []struct {
		bindHost string
		want     bool
	}{
		{"localhost", true},
		{"LocalHost", true},
		{"127.0.0.1", true},
		{"127.1.2.3", true},
		{"::1", true},
		{"[::1]", true},
		{"0.0.0.0", false},
		{"::", false},
		{"192.168.1.10", false},
		{"app.example", false},
	}
	for _, tt := range tests {
		if got := isLoopbackBindHost(tt.bindHost); got != tt.want {
			t.Errorf("isLoopbackBindHost(%q) = %v, want %v", tt.bindHost, got, tt.want)
		}
	}
}

func TestServerHostAllowed(t *testing.T) {
	tests := []struct {
		name        string
		requestHost string
		bindHost    string
		patterns    []string
		want        bool
	}{
		{name: "localhost", requestHost: "localhost:8080", bindHost: "localhost", want: true},
		{name: "localhost mixed case", requestHost: "LOCALHOST:8080", bindHost: "localhost", want: true},
		{name: "localhost without port", requestHost: "localhost", bindHost: "localhost", want: true},
		{name: "IPv4 loopback", requestHost: "127.0.0.1:8080", bindHost: "localhost", want: true},
		{name: "IPv6 loopback", requestHost: "[::1]:8080", bindHost: "localhost", want: true},
		{name: "IPv6 loopback without port", requestHost: "[::1]", bindHost: "localhost", want: true},
		{name: "empty host", requestHost: "", bindHost: "localhost", want: true},
		{name: "rebound name", requestHost: "evil.example:8080", bindHost: "localhost", want: false},
		{name: "rebound name without port", requestHost: "evil.example", bindHost: "127.0.0.1", want: false},
		{name: "localhost trailing dot", requestHost: "localhost.:8080", bindHost: "localhost", want: false},
		{name: "localhost subdomain", requestHost: "localhost.evil.example:8080", bindHost: "localhost", want: false},
		{name: "shorthand IP", requestHost: "127.1:8080", bindHost: "localhost", want: false},
		// "::1:8080" is itself a valid IPv6 literal, which cannot be rebound.
		{name: "unbracketed IPv6 literal", requestHost: "::1:8080", bindHost: "localhost", want: true},
		{
			name:        "pattern host",
			requestHost: "app.example",
			bindHost:    "localhost",
			patterns:    []string{"https://app.example"},
			want:        true,
		},
		{
			name:        "pattern host with port",
			requestHost: "app.example:8443",
			bindHost:    "localhost",
			patterns:    []string{"app.example"},
			want:        true,
		},
		{
			name:        "pattern wildcard",
			requestHost: "a.app.example",
			bindHost:    "localhost",
			patterns:    []string{"*.app.example"},
			want:        true,
		},
		{
			name:        "pattern mismatch",
			requestHost: "evil.example",
			bindHost:    "localhost",
			patterns:    []string{"https://app.example"},
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serverHostAllowed(tt.requestHost, tt.bindHost, tt.patterns); got != tt.want {
				t.Fatalf("serverHostAllowed(%q, %q, %q) = %v, want %v", tt.requestHost, tt.bindHost, tt.patterns, got, tt.want)
			}
		})
	}
}

// createHandler must apply the guard to every route, including /health,
// which never reaches the asset server.
func TestServerModeCreateHandlerAppliesHostGuard(t *testing.T) {
	a := newServerTestApp(t)
	h := &serverApp{app: a, bindHost: "localhost", broadcaster: NewWebSocketBroadcaster(a)}
	handler := h.createHandler()

	for host, want := range map[string]int{
		"evil.example:8080": http.StatusForbidden,
		"localhost:8080":    http.StatusOK,
	} {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Host = host
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Errorf("Host %q: status = %d, want %d", host, rec.Code, want)
		}
	}
}

// The rebinding attack sends a matching Origin and Host for the attacker's
// name, so the origin checks alone accept it; the Host guard must not.
func TestServerModeHostGuardRuntime(t *testing.T) {
	tests := []struct {
		name     string
		bindHost string
		server   ServerOptions
		host     string
		wantCode int
	}{
		{name: "rebound host rejected on localhost bind", bindHost: "localhost", host: "evil.example:8080", wantCode: http.StatusForbidden},
		{name: "rebound host rejected on IPv4 loopback bind", bindHost: "127.0.0.1", host: "evil.example:8080", wantCode: http.StatusForbidden},
		{name: "rebound host rejected on IPv6 loopback bind", bindHost: "::1", host: "evil.example:8080", wantCode: http.StatusForbidden},
		{
			name:     "rebound host rejected when all origins allowed",
			bindHost: "localhost",
			server:   ServerOptions{WebSocketAllowAllOrigins: true},
			host:     "evil.example:8080",
			wantCode: http.StatusForbidden,
		},
		{name: "localhost accepted", bindHost: "localhost", host: "localhost:8080", wantCode: http.StatusUnprocessableEntity},
		{name: "IPv4 loopback accepted", bindHost: "localhost", host: "127.0.0.1:8080", wantCode: http.StatusUnprocessableEntity},
		{name: "IPv6 loopback accepted", bindHost: "localhost", host: "[::1]:8080", wantCode: http.StatusUnprocessableEntity},
		{
			name:     "configured pattern host accepted",
			bindHost: "localhost",
			server:   ServerOptions{WebSocketOriginPatterns: []string{"https://app.example"}},
			host:     "app.example",
			wantCode: http.StatusUnprocessableEntity,
		},
		{name: "non-loopback bind unaffected", bindHost: "0.0.0.0", host: "app.example:8080", wantCode: http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newServerTestApp(t)
			a.options.Server = tt.server
			previous := globalApplication
			globalApplication = a
			t.Cleanup(func() { globalApplication = previous })

			next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				t.Fatal("runtime request fell through to the asset handler")
			})
			handler := requireAllowedHost(NewHTTPTransport().Handler()(next), a, tt.bindHost, tt.server)

			req := httptest.NewRequest(http.MethodPost, "/wails/runtime", strings.NewReader(`{}`))
			req.Host = tt.host
			req.Header.Set("Origin", "http://"+tt.host)
			req.Header.Set("Sec-Fetch-Site", "same-origin")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tt.wantCode, rec.Body.String())
			}
		})
	}
}

func TestServerModeHostGuardEventsWebSocket(t *testing.T) {
	tests := []struct {
		name       string
		bindHost   string
		server     ServerOptions
		host       string
		wantAccept bool
	}{
		{name: "rebound host rejected on loopback bind", bindHost: "localhost", host: "evil.example:8080"},
		{
			name:     "rebound host rejected when all origins allowed",
			bindHost: "127.0.0.1",
			server:   ServerOptions{WebSocketAllowAllOrigins: true},
			host:     "evil.example:8080",
		},
		{name: "localhost accepted", bindHost: "localhost", host: "localhost:8080", wantAccept: true},
		{name: "IPv4 loopback accepted", bindHost: "localhost", host: "127.0.0.1:8080", wantAccept: true},
		{name: "IPv6 loopback accepted", bindHost: "localhost", host: "[::1]:8080", wantAccept: true},
		{
			name:       "configured pattern host accepted",
			bindHost:   "localhost",
			server:     ServerOptions{WebSocketOriginPatterns: []string{"app.example"}},
			host:       "app.example",
			wantAccept: true,
		},
		{name: "non-loopback bind unaffected", bindHost: "0.0.0.0", host: "app.example:8080", wantAccept: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newServerTestApp(t)
			a.options.Server = tt.server
			srv := httptest.NewServer(requireAllowedHost(NewWebSocketBroadcaster(a), a, tt.bindHost, tt.server))
			defer srv.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			conn, resp, err := websocket.Dial(ctx, wsURL(srv.URL), &websocket.DialOptions{
				Host:       tt.host,
				HTTPHeader: http.Header{"Origin": []string{"http://" + tt.host}},
			})
			if tt.wantAccept {
				if err != nil {
					t.Fatalf("event WebSocket rejected: %v", err)
				}
				conn.CloseNow()
				return
			}
			if err == nil {
				conn.CloseNow()
				t.Fatal("rebound event WebSocket unexpectedly connected")
			}
			if resp == nil || resp.StatusCode != http.StatusForbidden {
				t.Fatalf("rejected response = %#v, want HTTP %d", resp, http.StatusForbidden)
			}
		})
	}
}
