package devruntime

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestNotifyIsOptInAndSentOnce(t *testing.T) {
	once = sync.Once{}
	calls := make(chan string, 3)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls <- r.Header.Get("Authorization")
		w.WriteHeader(204)
	}))
	defer server.Close()
	t.Setenv(AddressEnv, server.URL)
	t.Setenv(TokenEnv, "")
	Notify()
	select {
	case <-calls:
		t.Fatal("notified without CLI token")
	case <-time.After(20 * time.Millisecond):
	}
	t.Setenv(TokenEnv, "launch-specific-token")
	Notify()
	Notify()
	select {
	case header := <-calls:
		if header != "Bearer launch-specific-token" {
			t.Fatal(header)
		}
	case <-time.After(time.Second):
		t.Fatal("no readiness notification")
	}
	select {
	case <-calls:
		t.Fatal("duplicate notification")
	case <-time.After(20 * time.Millisecond):
	}
}
