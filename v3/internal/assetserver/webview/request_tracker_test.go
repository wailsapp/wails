package webview

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"testing"
)

type trackerTestRequest struct {
	finalizerTestRequest
	headers http.Header
}

func (r *trackerTestRequest) Header() (http.Header, error) { return r.headers, nil }
func trackedTestRequest(token string) *trackerTestRequest {
	return &trackerTestRequest{headers: http.Header{RequestIDHeader: {token}, "Content-Type": {"application/json"}}}
}

func TestRequestTrackerIdenticalURLsAreIndependent(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	defer tracker.Close()
	a := trackedTestRequest(tracker.Start("a"))
	b := trackedTestRequest(tracker.Start("b"))
	first, second := tracker.Wrap(a), tracker.Wrap(b)
	defer first.Close()
	defer second.Close()
	tracker.Cancel("a")
	if requestContext(first).Err() != context.Canceled {
		t.Fatal("aborted request not cancelled")
	}
	if requestContext(second).Err() != nil {
		t.Fatal("unrelated same-URL request cancelled")
	}
	if a.headers.Get(RequestIDHeader) != "" || a.headers.Get("Content-Type") != "application/json" {
		t.Fatal("internal header leaked or application headers lost")
	}
}

func TestRequestTrackerAbortBeforeClaim(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	defer tracker.Close()
	token := tracker.Start("request")
	tracker.Cancel("request")
	if len(tracker.byID) != 0 || len(tracker.byToken) != 0 {
		t.Fatal("unclaimed abort leaked")
	}
	r := tracker.Wrap(trackedTestRequest(token))
	defer r.Close()
	if requestContext(r).Err() != context.Canceled {
		t.Fatal("late callback resurrected an aborted request")
	}
}

func TestRequestTrackerRedirectReusesNetworkID(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	defer tracker.Close()
	old := tracker.Wrap(trackedTestRequest(tracker.Start("redirect")))
	current := tracker.Wrap(trackedTestRequest(tracker.Start("redirect")))
	defer current.Close()
	old.Close()
	if requestContext(current).Err() != nil {
		t.Fatal("old response cleanup cancelled its successor")
	}
	tracker.Cancel("redirect")
	if requestContext(current).Err() != context.Canceled {
		t.Fatal("successor lost cancellation registration")
	}
}

func TestRequestTrackerWindowCloseAndLateRequests(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	unmarked := tracker.Wrap(trackedTestRequest(""))
	marked := tracker.Wrap(trackedTestRequest(tracker.Start("active")))
	tracker.Start("unclaimed")
	tracker.Close()
	defer unmarked.Close()
	defer marked.Close()
	for _, r := range []Request{unmarked, marked} {
		if requestContext(r).Err() != context.Canceled {
			t.Fatal("window close missed a request")
		}
	}
	if len(tracker.byID) != 0 || len(tracker.byToken) != 0 {
		t.Fatal("window close leaked requests")
	}
	if tracker.Start("late") != "" {
		t.Fatal("closed tracker accepted a request")
	}
	late := tracker.Wrap(trackedTestRequest(""))
	defer late.Close()
	if requestContext(late).Err() != context.Canceled {
		t.Fatal("late request escaped window cancellation")
	}
}

func TestRequestTrackerConcurrentCompletionAndAbort(t *testing.T) {
	for range 100 {
		tracker := NewRequestTracker(context.Background())
		raw := trackedTestRequest(tracker.Start("request"))
		r := tracker.Wrap(raw)
		var workers sync.WaitGroup
		for range 8 {
			workers.Go(func() { tracker.Cancel("request") })
			workers.Go(func() { r.Close() })
			workers.Go(func() { tracker.Close() })
		}
		workers.Wait()
		if raw.closeCount.Load() != 1 {
			t.Fatal("request closed more than once")
		}
		if len(tracker.byID) != 0 || len(tracker.byToken) != 0 {
			t.Fatal("completed request leaked")
		}
	}
}

func TestRequestTrackerRejectsDuplicateAndForeignMarkers(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	defer tracker.Close()
	token := tracker.Start("request")
	owner := tracker.Wrap(trackedTestRequest(token))
	defer owner.Close()
	for _, marker := range []string{token, "untrusted"} {
		r := tracker.Wrap(trackedTestRequest(marker))
		defer r.Close()
		if requestContext(r).Err() != context.Canceled {
			t.Fatal("invalid marker acquired a live context")
		}
	}
	if requestContext(owner).Err() != nil {
		t.Fatal("invalid marker cancelled another request")
	}
}

func TestRequestTrackerParentCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	tracker := NewRequestTracker(parent)
	defer tracker.Close()
	r := tracker.Wrap(trackedTestRequest(tracker.Start("active")))
	defer r.Close()
	cancel()
	if requestContext(r).Err() != context.Canceled || tracker.Start("late") != "" {
		t.Fatal("application shutdown did not cancel requests")
	}
}

func TestRequestTrackerKeepaliveRedirectAndFailedDelivery(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	defer tracker.Close()
	token := tracker.StartKeepalive("keepalive")
	redirect := tracker.Wrap(trackedTestRequest(token))
	redirect.Close()
	if tracker.Token("keepalive") != token {
		t.Fatal("redirect retired keepalive identity")
	}
	handler := tracker.Wrap(trackedTestRequest(token))
	tracker.FailKeepalive("keepalive")
	if requestContext(handler).Err() != nil {
		t.Fatal("page teardown cancelled running keepalive")
	}
	handler.Close()
	if tracker.Token("keepalive") != "" {
		t.Fatal("failed delivery leaked a completed keepalive")
	}
}

func TestRequestTrackerKeepaliveAbortAndFailureBeforeDispatch(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	defer tracker.Close()
	token := tracker.StartKeepalive("active")
	handler := tracker.Wrap(trackedTestRequest(token))
	defer handler.Close()
	tracker.Cancel("active")
	if requestContext(handler).Err() != context.Canceled {
		t.Fatal("explicit abort missed keepalive")
	}
	tracker.StartKeepalive("failed")
	tracker.FailKeepalive("failed")
	if tracker.Token("failed") != "" {
		t.Fatal("failed fetch leaked before dispatch")
	}
}

func TestRequestTrackerBoundsKeepaliveRegistrations(t *testing.T) {
	tracker := NewRequestTracker(context.Background())
	defer tracker.Close()
	for i := range maxKeepaliveRequests {
		if tracker.StartKeepalive(strconv.Itoa(i)) == "" {
			t.Fatal("registration rejected below capacity")
		}
	}
	if tracker.StartKeepalive("overflow") != "" {
		t.Fatal("unbounded keepalive registration")
	}
	old := tracker.Wrap(trackedTestRequest(tracker.Token("0")))
	tracker.Cancel("0")
	tracker.Cancel("unknown")
	if tracker.StartKeepalive("replacement") == "" {
		t.Fatal("cancellation did not release capacity")
	}
	old.Close()
	if tracker.keepalives != maxKeepaliveRequests {
		t.Fatal("old close released capacity twice")
	}
	tracker.Close()
	if tracker.keepalives != 0 {
		t.Fatal("window close retained registrations")
	}
}
