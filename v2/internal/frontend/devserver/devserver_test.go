//go:build dev

package devserver

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/ipcauth"
	"github.com/wailsapp/wails/v2/internal/logger"
)

const echoPrefix = "echo:"

// stubDispatcher echoes each message back, so the IPC round-trip can be
// observed without a real bindings dispatcher.
type stubDispatcher struct{}

func (stubDispatcher) ProcessMessage(message string, _ frontend.Frontend) (string, error) {
	return echoPrefix + message, nil
}

// nullFrontend satisfies frontend.Frontend for the reload path. Only
// WindowReload is exercised by these tests; the embedded nil interface covers
// the rest, which are never called.
type nullFrontend struct {
	frontend.Frontend
}

func (nullFrontend) WindowReload() {}

// newTestServer builds a DevWebServer served over httptest and returns it with
// the server and its bind port. devServerAddr is set to the httptest listener
// address before registerRoutes runs, so the host guard's allowlist is keyed on
// the real bind address (a loopback IP).
func newTestServer(t *testing.T) (*DevWebServer, *httptest.Server, string) {
	t.Helper()
	d := &DevWebServer{
		server:           echo.New(),
		logger:           logger.New(nil),
		dispatcher:       stubDispatcher{},
		Frontend:         nullFrontend{},
		websocketClients: make(map[*websocket.Conn]*sync.Mutex),
	}
	d.server.HideBanner = true
	d.server.HidePort = true

	ts := httptest.NewServer(d.server)
	t.Cleanup(ts.Close)

	d.devServerAddr = ts.Listener.Addr().String()
	d.registerRoutes()

	_, port, err := net.SplitHostPort(d.devServerAddr)
	if err != nil {
		t.Fatalf("splitting test server addr %q: %v", d.devServerAddr, err)
	}
	return d, ts, port
}

func getWithHost(t *testing.T, url, host string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("building request: %v", err)
	}
	if host != "" {
		req.Host = host
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request to %s (host %q): %v", url, host, err)
	}
	return resp
}

func hasCapabilityCookie(resp *http.Response, port string) bool {
	for _, ck := range resp.Cookies() {
		if ck.Name == ipcauth.CookieName+"_"+port {
			return true
		}
	}
	return false
}

func statusOf(resp *http.Response) int {
	if resp == nil {
		return -1
	}
	return resp.StatusCode
}

func wsURL(ts *httptest.Server) string {
	return "ws" + strings.TrimPrefix(ts.URL, "http") + "/wails/ipc"
}

func TestReloadRouteHostGuard(t *testing.T) {
	_, ts, port := newTestServer(t)
	reloadURL := ts.URL + "/wails/reload"

	t.Run("allowed host sets the capability cookie", func(t *testing.T) {
		resp := getWithHost(t, reloadURL, "")
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
		}
		if !hasCapabilityCookie(resp, port) {
			t.Errorf("an allowed response should carry a %s_%s cookie", ipcauth.CookieName, port)
		}
	})

	t.Run("localhost host is allowed against a loopback-IP bind", func(t *testing.T) {
		resp := getWithHost(t, reloadURL, "localhost:"+port)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
		}
	})

	t.Run("forged host is rejected and gets no cookie", func(t *testing.T) {
		resp := getWithHost(t, reloadURL, "attacker.example:"+port)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
		if hasCapabilityCookie(resp, port) {
			t.Errorf("a rejected request must not receive a %s_%s cookie", ipcauth.CookieName, port)
		}
	})

	t.Run("wrong port is rejected", func(t *testing.T) {
		resp := getWithHost(t, reloadURL, "127.0.0.1:1")
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
	})
}

func TestIPCWebSocketGate(t *testing.T) {
	_, ts, port := newTestServer(t)
	sameOriginHeader := ts.URL // http://127.0.0.1:<port>, equal to the request Host

	withCookie := func() http.Header {
		h := http.Header{}
		h.Set("Cookie", ipcauth.CookieName+"_"+port+"="+ipcauth.Token())
		return h
	}
	expectRejected := func(t *testing.T, conn *websocket.Conn, resp *http.Response, err error) {
		t.Helper()
		if err == nil {
			conn.Close()
			t.Fatal("expected the upgrade to be rejected")
		}
		if statusOf(resp) != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", statusOf(resp), http.StatusForbidden)
		}
	}

	t.Run("missing cookie is rejected", func(t *testing.T) {
		h := http.Header{}
		h.Set("Origin", sameOriginHeader)
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), h)
		expectRejected(t, conn, resp, err)
	})

	t.Run("valid cookie and same origin upgrade and echo", func(t *testing.T) {
		h := withCookie()
		h.Set("Origin", sameOriginHeader)
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), h)
		if err != nil {
			t.Fatalf("upgrade failed: %v (status %d)", err, statusOf(resp))
		}
		defer conn.Close()
		if resp.StatusCode != http.StatusSwitchingProtocols {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusSwitchingProtocols)
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
			t.Fatalf("write: %v", err)
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if string(msg) != echoPrefix+"ping" {
			t.Errorf("got %q, want %q", msg, echoPrefix+"ping")
		}
	})

	// Browsers do not scope cookies by port, so another dev server's cookie may
	// be sent along; it must not open this server's IPC.
	t.Run("another port's cookie is rejected", func(t *testing.T) {
		h := http.Header{}
		h.Set("Cookie", ipcauth.CookieName+"_1="+ipcauth.Token())
		h.Set("Origin", sameOriginHeader)
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), h)
		expectRejected(t, conn, resp, err)
	})

	t.Run("unscoped cookie name is rejected", func(t *testing.T) {
		h := http.Header{}
		h.Set("Cookie", ipcauth.CookieName+"="+ipcauth.Token())
		h.Set("Origin", sameOriginHeader)
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), h)
		expectRejected(t, conn, resp, err)
	})

	t.Run("foreign origin is rejected", func(t *testing.T) {
		h := withCookie()
		h.Set("Origin", "http://attacker.example")
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), h)
		expectRejected(t, conn, resp, err)
	})

	t.Run("missing origin is rejected", func(t *testing.T) {
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), withCookie())
		expectRejected(t, conn, resp, err)
	})

	// The DNS-rebinding regression: a rebound page presents Host and Origin that
	// agree (both the attacker's name) and a valid capability cookie the browser
	// attached first-party. Only the Host allowlist stands between it and the
	// dispatcher — so this must be rejected, and it must be the guard doing it.
	t.Run("rebound host is rejected despite matching origin and a valid cookie", func(t *testing.T) {
		forged := "attacker.example:" + port
		h := withCookie()
		h.Set("Host", forged)
		h.Set("Origin", "http://"+forged)
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), h)
		expectRejected(t, conn, resp, err)
	})
}

func TestCapabilityCookieName(t *testing.T) {
	tests := []struct {
		addr string
		want string
	}{
		{"localhost:34115", ipcauth.CookieName + "_34115"},
		{"0.0.0.0:5000", ipcauth.CookieName + "_5000"},
		{"[::1]:34115", ipcauth.CookieName + "_34115"},
		{":34115", ipcauth.CookieName + "_34115"},
		{"", ipcauth.CookieName},
	}
	for _, tt := range tests {
		if got := capabilityCookieName(tt.addr); got != tt.want {
			t.Errorf("capabilityCookieName(%q) = %q, want %q", tt.addr, got, tt.want)
		}
	}
}

// TestIPCAuthRoute covers the no-op endpoint the dev runtime fetches before
// each IPC (re)connect to refresh its capability cookie.
func TestIPCAuthRoute(t *testing.T) {
	_, ts, port := newTestServer(t)
	authURL := ts.URL + "/wails/ipcauth"

	t.Run("allowed host gets 204 and a port-scoped cookie", func(t *testing.T) {
		resp := getWithHost(t, authURL, "")
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
		}
		var cookie *http.Cookie
		for _, ck := range resp.Cookies() {
			if ck.Name == ipcauth.CookieName+"_"+port {
				cookie = ck
			}
		}
		if cookie == nil {
			t.Fatalf("response should carry a %s_%s cookie", ipcauth.CookieName, port)
		}
		if !ipcauth.Valid(cookie.Value) {
			t.Errorf("cookie value is not the current capability")
		}
		if !cookie.HttpOnly || cookie.SameSite != http.SameSiteStrictMode || cookie.Path != "/" {
			t.Errorf("cookie attributes = HttpOnly %v, SameSite %v, Path %q; want true, Strict, /", cookie.HttpOnly, cookie.SameSite, cookie.Path)
		}
	})

	t.Run("forged host is rejected and gets no cookie", func(t *testing.T) {
		resp := getWithHost(t, authURL, "attacker.example:"+port)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
		if hasCapabilityCookie(resp, port) {
			t.Errorf("a rejected request must not receive a %s_%s cookie", ipcauth.CookieName, port)
		}
	})

	t.Run("refreshed cookie opens the IPC socket", func(t *testing.T) {
		resp := getWithHost(t, authURL, "")
		resp.Body.Close()
		h := http.Header{}
		for _, ck := range resp.Cookies() {
			h.Add("Cookie", ck.Name+"="+ck.Value)
		}
		h.Set("Origin", ts.URL)
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL(ts), h)
		if err != nil {
			t.Fatalf("upgrade failed: %v (status %d)", err, statusOf(resp))
		}
		conn.Close()
	})
}

// TestAllowedHostsEnvVar proves the escape hatch: a host that the fixed rules
// reject is accepted once named in WAILS_DEV_ALLOWED_HOSTS. The env var is read
// when routes are registered, so it must be set before newTestServer.
func TestAllowedHostsEnvVar(t *testing.T) {
	t.Setenv(allowedHostsEnvVar, "attacker.example")
	_, ts, port := newTestServer(t)

	resp := getWithHost(t, ts.URL+"/wails/reload", "attacker.example:"+port)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("with %s set, status = %d, want %d", allowedHostsEnvVar, resp.StatusCode, http.StatusNoContent)
	}
}
