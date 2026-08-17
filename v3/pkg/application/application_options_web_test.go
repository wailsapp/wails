//go:build !wails_native

package application

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChainMiddleware_Empty(t *testing.T) {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("base"))
	})

	handler := ChainMiddleware()(baseHandler)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "base" {
		t.Errorf("Body = %q, want %q", rec.Body.String(), "base")
	}
}

func TestChainMiddleware_Single(t *testing.T) {
	callOrder := []string{}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "middleware")
			next.ServeHTTP(w, r)
		})
	}
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "base")
		w.WriteHeader(http.StatusOK)
	})

	handler := ChainMiddleware(middleware)(baseHandler)
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	expected := []string{"middleware", "base"}
	for i, value := range expected {
		if i >= len(callOrder) || callOrder[i] != value {
			t.Fatalf("Call order = %v, want %v", callOrder, expected)
		}
	}
}

func TestChainMiddleware_Multiple(t *testing.T) {
	callOrder := []string{}
	middleware := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callOrder = append(callOrder, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "base")
		w.WriteHeader(http.StatusOK)
	})

	handler := ChainMiddleware(middleware("middleware1"), middleware("middleware2"), middleware("middleware3"))(baseHandler)
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	expected := []string{"middleware1", "middleware2", "middleware3", "base"}
	for i, value := range expected {
		if i >= len(callOrder) || callOrder[i] != value {
			t.Fatalf("Call order = %v, want %v", callOrder, expected)
		}
	}
}

func TestAssetOptions_Defaults(t *testing.T) {
	opts := AssetOptions{}
	if opts.Handler != nil {
		t.Error("Handler should default to nil")
	}
	if opts.Middleware != nil {
		t.Error("Middleware should default to nil")
	}
	if opts.DisableLogging {
		t.Error("DisableLogging should default to false")
	}
}

func TestMiddleware_ShortCircuit(t *testing.T) {
	shortCircuit := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("forbidden"))
		})
	}
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("base"))
	})

	handler := ChainMiddleware(shortCircuit)(baseHandler)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusForbidden)
	}
	if rec.Body.String() != "forbidden" {
		t.Errorf("Body = %q, want %q", rec.Body.String(), "forbidden")
	}
}
