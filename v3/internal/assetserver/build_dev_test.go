//go:build !production

package assetserver

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
)

// A dial that failed because the host had no ephemeral port left must not be treated as
// transient: retrying prolongs an outage that is already affecting every process on the machine.
// That is the current behaviour by omission, and this guards it, since the obvious "fix" for the
// EADDRNOTAVAIL errors this bug produces is to add them to isConnectionError.
func TestIsConnectionErrorIgnoresAddressExhaustion(t *testing.T) {
	tests := map[string]struct {
		err  error
		want bool
	}{
		"nil": {nil, false},
		"refused": {
			errors.New("dial tcp4 127.0.0.1:9245: connect: connection refused"), true},
		"reset": {
			errors.New("read tcp 127.0.0.1:9245: connection reset by peer"), true},
		"windows refused": {
			errors.New("dial tcp4 127.0.0.1:9245: connectex: No connection could be made " +
				"because the target machine actively refused it."), true},
		"darwin EADDRNOTAVAIL": {
			errors.New("dial tcp4 127.0.0.1:9245: connect: can't assign requested address"), false},
		"linux EADDRNOTAVAIL": {
			errors.New("dial tcp4 127.0.0.1:9245: connect: cannot assign requested address"), false},
		// Windows binds a local port before calling ConnectEx, so exhaustion surfaces from the
		// bind rather than as one of the "connectex" errors above.
		"windows bind WSAEADDRINUSE": {
			errors.New("dial tcp4 127.0.0.1:9245: bind: Only one usage of each socket " +
				"address (protocol/network address/port) is normally permitted."), false},
		"windows bind WSAEADDRNOTAVAIL": {
			errors.New("dial tcp4 127.0.0.1:9245: bind: The requested address is not " +
				"valid in its context."), false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isConnectionError(tt.err); got != tt.want {
				t.Errorf("isConnectionError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// The dev proxy must reuse its upstream connections. With Go's default of two idle connections
// per host, the concurrent module requests a dev page makes each open and close their own TCP
// connection, and every one of those leaves a TIME_WAIT on an ephemeral port for 2*MSL. On a
// large frontend that drains the host's entire ephemeral range and every process on the machine
// starts failing to make outbound connections.
func TestDevProxyReusesUpstreamConnections(t *testing.T) {
	const (
		concurrency = 64
		requests    = 1000
	)

	var accepted atomic.Int64
	backend := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			io.WriteString(w, "ok")
		}))
	backend.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			accepted.Add(1)
		}
	}
	backend.Start()
	defer backend.Close()

	t.Setenv("FRONTEND_DEVSERVER_URL", backend.URL)

	front := httptest.NewServer(NewAssetFileServer(nil))
	defer front.Close()

	client := front.Client()
	// Keep the test's own client from churning connections the way the bug under test does.
	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.MaxIdleConns = concurrency
		transport.MaxIdleConnsPerHost = concurrency
	}

	work := make(chan struct{}, requests)
	for i := 0; i < requests; i++ {
		work <- struct{}{}
	}
	close(work)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range work {
				resp, err := client.Get(front.URL + "/module.js")
				if err != nil {
					t.Errorf("request failed: %v", err)
					return
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
					return
				}
			}
		}()
	}
	wg.Wait()

	// With reuse the proxy needs roughly one upstream connection per concurrent in-flight
	// request rather than one per request. The ceiling is deliberately generous so scheduler
	// noise cannot make this flaky: unpatched this count is close to `requests`, patched it is
	// close to `concurrency`.
	t.Logf("proxy opened %d upstream connections for %d requests", accepted.Load(), requests)
	if limit := int64(concurrency * 2); accepted.Load() > limit {
		t.Fatalf("proxy opened %d upstream connections for %d requests (limit %d) - "+
			"connections are not being reused; check MaxIdleConnsPerHost",
			accepted.Load(), requests, limit)
	}
}
