//go:build !mcp

package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestDisabledServiceIsInert(t *testing.T) {
	service := New()
	if service.Enabled() {
		t.Fatal("Enabled() must be false without the mcp build tag")
	}
	if err := service.ServiceStartup(context.Background(), application.ServiceOptions{}); err != nil {
		t.Fatalf("stub startup failed: %v", err)
	}
	if err := service.ServiceShutdown(); err != nil {
		t.Fatalf("stub shutdown failed: %v", err)
	}

	recorder := httptest.NewRecorder()
	service.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected 404 from disabled service, got %d", recorder.Code)
	}
}
