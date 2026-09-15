//go:build windows && !server

package application

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

type cancellationCall struct{ session, method, params string }
type cancellationTestCDP struct {
	calls   []cancellationCall
	events  map[string]func(string, string, error)
	fail    string
	removed int
}

func (f *cancellationTestCDP) CallDevTools(method, params string, done func(string, error)) error {
	return f.CallDevToolsForSession("", method, params, done)
}
func (f *cancellationTestCDP) CallDevToolsForSession(session, method, params string, done func(string, error)) error {
	f.calls = append(f.calls, cancellationCall{session, method, params})
	if method == f.fail {
		done("", errors.New("unsupported"))
	} else {
		done("{}", nil)
	}
	return nil
}
func (f *cancellationTestCDP) OnDevTools(event string, receive func(string, string, error)) (func() error, error) {
	f.events[event] = receive
	return func() error { f.removed++; return nil }, nil
}
func cancellationFixture(t *testing.T) (*windowsRequestCancellation, *cancellationTestCDP) {
	t.Helper()
	previous := globalApplication
	globalApplication = &App{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	fake := &cancellationTestCDP{events: make(map[string]func(string, string, error))}
	c := &windowsRequestCancellation{chromium: fake, tracker: webview.NewRequestTracker(context.Background())}
	t.Cleanup(func() { c.close(); globalApplication = previous })
	if err := c.subscribe(); err != nil {
		t.Fatal(err)
	}
	return c, fake
}
func emitCancellation(t *testing.T, c *windowsRequestCancellation, session, event string, payload any) {
	t.Helper()
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.event(session, event, string(data)); err != nil {
		t.Fatal(err)
	}
}
func keepaliveState(t *testing.T, c *windowsRequestCancellation, id, state string) {
	t.Helper()
	payload, _ := json.Marshal(map[string]string{"ID": id, "State": state})
	emitCancellation(t, c, "", "Runtime.bindingCalled", map[string]string{"name": keepaliveBinding, "payload": string(payload)})
}

type cancellationTestRequest struct {
	webview.Request
	header http.Header
}

func (r *cancellationTestRequest) Header() (http.Header, error) { return r.header, nil }
func (r *cancellationTestRequest) Close() error                 { return nil }
func pauseCancellation(t *testing.T, c *windowsRequestCancellation, f *cancellationTestCDP, id string, extra map[string]string) webview.Request {
	t.Helper()
	headers := map[string]string{"Accept": "application/json"}
	for k, v := range extra {
		if k != "fragment" {
			headers[k] = v
		}
	}
	emitCancellation(t, c, "", "Fetch.requestPaused", map[string]any{"requestId": "pause-" + id, "networkId": id, "request": map[string]any{"headers": headers, "urlFragment": extra["fragment"]}})
	var continued struct{ Headers []requestHeader }
	if err := json.Unmarshal([]byte(f.calls[len(f.calls)-1].params), &continued); err != nil {
		t.Fatal(err)
	}
	header := make(http.Header)
	for _, h := range continued.Headers {
		header.Add(h.Name, h.Value)
	}
	if header.Get("fragment") != "" {
		t.Fatal("fragment leaked into headers")
	}
	r := c.tracker.Wrap(&cancellationTestRequest{header: header})
	t.Cleanup(func() { r.Close() })
	return r
}
func isCancelled(r webview.Request) bool {
	ctx, _ := webview.RequestContext(r)
	return ctx.Err() == context.Canceled
}

func TestWindowsCancellationWorkerUsesRootInterceptionIdentity(t *testing.T) {
	c, f := cancellationFixture(t)
	emitCancellation(t, c, "worker", "Network.requestWillBeSent", map[string]string{"requestId": "42.1"})
	worker := pauseCancellation(t, c, f, "42.1", nil)
	control := pauseCancellation(t, c, f, "42.2", nil)
	emitCancellation(t, c, "worker", "Network.loadingFailed", map[string]string{"requestId": "42.1"})
	if !isCancelled(worker) || isCancelled(control) {
		t.Fatal("worker abort was lost or cancelled another request")
	}
}

func TestWindowsCancellationWorkerDetach(t *testing.T) {
	c, f := cancellationFixture(t)
	emitCancellation(t, c, "worker", "Network.requestWillBeSent", map[string]string{"requestId": "42.1"})
	worker := pauseCancellation(t, c, f, "42.1", nil)
	control := pauseCancellation(t, c, f, "42.2", nil)
	emitCancellation(t, c, "", "Target.detachedFromTarget", map[string]string{"sessionId": "worker"})
	if !isCancelled(worker) || isCancelled(control) {
		t.Fatal("detachment cancelled the wrong requests")
	}
}

func TestWindowsCancellationKeepaliveSurvivesNavigation(t *testing.T) {
	c, f := cancellationFixture(t)
	id := "11111111-1111-1111-1111-111111111111"
	keepaliveState(t, c, id, "start")
	redirect := pauseCancellation(t, c, f, "42.1", map[string]string{"fragment": "#" + keepaliveFragment + id})
	redirect.Close()
	handler := pauseCancellation(t, c, f, "42.1", map[string]string{"fragment": "#" + keepaliveFragment + id})
	emitCancellation(t, c, "", "Network.loadingFailed", map[string]string{"requestId": "42.1"})
	keepaliveState(t, c, id, "failed")
	if isCancelled(handler) {
		t.Fatal("navigation cancelled keepalive handler")
	}
	keepaliveState(t, c, id, "abort")
	if !isCancelled(handler) {
		t.Fatal("explicit abort did not cancel keepalive")
	}
}

func TestWindowsCancellationKeepaliveAbortBeforeInterception(t *testing.T) {
	c, f := cancellationFixture(t)
	id := "11111111-1111-1111-1111-111111111111"
	keepaliveState(t, c, id, "start")
	keepaliveState(t, c, id, "abort")
	r := pauseCancellation(t, c, f, "42.1", map[string]string{"fragment": "#" + keepaliveFragment + id})
	if !isCancelled(r) {
		t.Fatal("late interception resurrected keepalive")
	}
}

func TestWindowsCancellationSetupAndFailedWorkerResume(t *testing.T) {
	c, f := cancellationFixture(t)
	var setupErr error
	c.enableSession("", func(err error) { setupErr = err })
	if setupErr != nil {
		t.Fatal(setupErr)
	}
	foundFetch := false
	for _, call := range f.calls {
		if call.method == "Fetch.enable" {
			foundFetch = true
		}
	}
	if !foundFetch {
		t.Fatal("root interception missing")
	}
	f.calls = nil
	f.fail = "Network.enable"
	emitCancellation(t, c, "", "Target.attachedToTarget", map[string]string{"sessionId": "worker"})
	last := f.calls[len(f.calls)-1]
	if last.session != "worker" || last.method != "Runtime.runIfWaitingForDebugger" {
		t.Fatal("failed setup left worker paused")
	}
}

func TestWindowsCancellationCloseReleasesRequestsAndSubscriptions(t *testing.T) {
	c, f := cancellationFixture(t)
	r := pauseCancellation(t, c, f, "42.1", nil)
	c.close()
	c.close()
	if !isCancelled(r) || f.removed != len(f.events) {
		t.Fatal("close leaked request or event callback")
	}
	calls := len(f.calls)
	for _, receive := range f.events {
		receive("", "invalid late event", nil)
	}
	if len(f.calls) != calls {
		t.Fatal("late event accessed closed WebView")
	}
}

func TestWindowsCancellationKeepaliveContextLossReleasesAfterHandler(t *testing.T) {
	c, f := cancellationFixture(t)
	id := "11111111-1111-1111-1111-111111111111"
	keepaliveState(t, c, id, "start")
	handler := pauseCancellation(t, c, f, "42.1", map[string]string{"fragment": "#" + keepaliveFragment + id})
	emitCancellation(t, c, "", "Runtime.executionContextsCleared", struct{}{})
	if isCancelled(handler) {
		t.Fatal("context loss cancelled a running keepalive")
	}
	handler.Close()
	if c.tracker.Token("keepalive:"+id) != "" || len(c.keepaliveOwners) != 0 {
		t.Fatal("context loss leaked keepalive registration")
	}
}
