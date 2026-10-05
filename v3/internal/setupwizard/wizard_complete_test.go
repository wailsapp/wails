package setupwizard

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCompleteWaitsForBackgroundBuild(t *testing.T) {
	w := New()
	w.origin = "http://127.0.0.1:12345"
	w.dockerMu.Lock()
	w.dockerStatus.PullStatus = "pulling"
	w.buildWg.Add(1)
	w.dockerMu.Unlock()

	request := httptest.NewRequest(http.MethodPost, "/api/complete", nil)
	request.Header.Set("Origin", w.origin)
	response := httptest.NewRecorder()
	w.handleComplete(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("completion returned %d: %s", response.Code, response.Body.String())
	}
	if !w.stopping {
		t.Fatal("completion did not prevent new background builds")
	}
	select {
	case <-w.done:
		t.Fatal("wizard stopped before background build finished")
	default:
	}
	w.buildWg.Done()
	select {
	case <-w.done:
	case <-time.After(time.Second):
		t.Fatal("wizard did not stop after background build finished")
	}
}

func TestCompleteRejectsGetAndCrossOrigin(t *testing.T) {
	w := New()
	w.origin = "http://127.0.0.1:12345"
	for _, tc := range []struct {
		method string
		origin string
		want   int
	}{
		{http.MethodGet, w.origin, http.StatusMethodNotAllowed},
		{http.MethodPost, "https://example.com", http.StatusForbidden},
	} {
		request := httptest.NewRequest(tc.method, "/api/complete", nil)
		request.Header.Set("Origin", tc.origin)
		response := httptest.NewRecorder()
		w.handleComplete(response, request)
		if response.Code != tc.want {
			t.Errorf("%s from %s returned %d, want %d", tc.method, tc.origin, response.Code, tc.want)
		}
		select {
		case <-w.done:
			t.Fatal("rejected completion stopped the wizard")
		default:
		}
	}
}
