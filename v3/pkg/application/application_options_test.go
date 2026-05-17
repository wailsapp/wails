package application

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActivationPolicy_Constants(t *testing.T) {
	if ActivationPolicyRegular != 0 {
		t.Error("ActivationPolicyRegular should be 0")
	}
	if ActivationPolicyAccessory != 1 {
		t.Error("ActivationPolicyAccessory should be 1")
	}
	if ActivationPolicyProhibited != 2 {
		t.Error("ActivationPolicyProhibited should be 2")
	}
}

func TestNativeTabIcon_Constants(t *testing.T) {
	tests := []struct {
		name     string
		icon     NativeTabIcon
		expected string
	}{
		{"NativeTabIconNone", NativeTabIconNone, ""},
		{"NativeTabIconHouse", NativeTabIconHouse, "house"},
		{"NativeTabIconGear", NativeTabIconGear, "gear"},
		{"NativeTabIconStar", NativeTabIconStar, "star"},
		{"NativeTabIconPerson", NativeTabIconPerson, "person"},
		{"NativeTabIconBell", NativeTabIconBell, "bell"},
		{"NativeTabIconMagnify", NativeTabIconMagnify, "magnifyingglass"},
		{"NativeTabIconList", NativeTabIconList, "list.bullet"},
		{"NativeTabIconFolder", NativeTabIconFolder, "folder"},
	}

	for _, tt := range tests {
		if string(tt.icon) != tt.expected {
			t.Errorf("%s = %q, want %q", tt.name, string(tt.icon), tt.expected)
		}
	}
}

func TestChainMiddleware_Empty(t *testing.T) {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("base"))
	})

	chained := ChainMiddleware()
	handler := chained(baseHandler)

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

	chained := ChainMiddleware(middleware)
	handler := chained(baseHandler)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if len(callOrder) != 2 {
		t.Errorf("Expected 2 calls, got %d", len(callOrder))
	}
	if callOrder[0] != "middleware" {
		t.Errorf("First call should be middleware, got %s", callOrder[0])
	}
	if callOrder[1] != "base" {
		t.Errorf("Second call should be base, got %s", callOrder[1])
	}
}

func TestChainMiddleware_Multiple(t *testing.T) {
	callOrder := []string{}

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "middleware1")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "middleware2")
			next.ServeHTTP(w, r)
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "middleware3")
			next.ServeHTTP(w, r)
		})
	}

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "base")
		w.WriteHeader(http.StatusOK)
	})

	chained := ChainMiddleware(middleware1, middleware2, middleware3)
	handler := chained(baseHandler)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := []string{"middleware1", "middleware2", "middleware3", "base"}
	if len(callOrder) != len(expected) {
		t.Errorf("Expected %d calls, got %d", len(expected), len(callOrder))
	}
	for i, v := range expected {
		if i < len(callOrder) && callOrder[i] != v {
			t.Errorf("Call %d: expected %s, got %s", i, v, callOrder[i])
		}
	}
}

func TestOptions_Defaults(t *testing.T) {
	opts := Options{}

	if opts.Name != "" {
		t.Error("Name should default to empty string")
	}
	if opts.Description != "" {
		t.Error("Description should default to empty string")
	}
	if opts.Icon != nil {
		t.Error("Icon should default to nil")
	}
	if opts.Logger != nil {
		t.Error("Logger should default to nil")
	}
	if opts.DisableDefaultSignalHandler != false {
		t.Error("DisableDefaultSignalHandler should default to false")
	}
}

func TestMacOptions_Defaults(t *testing.T) {
	opts := MacOptions{}

	if opts.ActivationPolicy != ActivationPolicyRegular {
		t.Error("ActivationPolicy should default to ActivationPolicyRegular")
	}
	if opts.ApplicationShouldTerminateAfterLastWindowClosed != false {
		t.Error("ApplicationShouldTerminateAfterLastWindowClosed should default to false")
	}
}

func TestWindowsOptions_Defaults(t *testing.T) {
	opts := WindowsOptions{}

	if opts.WndClass != "" {
		t.Error("WndClass should default to empty string")
	}
	if opts.DisableQuitOnLastWindowClosed != false {
		t.Error("DisableQuitOnLastWindowClosed should default to false")
	}
	if opts.WebviewUserDataPath != "" {
		t.Error("WebviewUserDataPath should default to empty string")
	}
	if opts.WebviewBrowserPath != "" {
		t.Error("WebviewBrowserPath should default to empty string")
	}
}

func TestLinuxOptions_Defaults(t *testing.T) {
	opts := LinuxOptions{}

	if opts.DisableQuitOnLastWindowClosed != false {
		t.Error("DisableQuitOnLastWindowClosed should default to false")
	}
	if opts.ProgramName != "" {
		t.Error("ProgramName should default to empty string")
	}
}

func TestIOSOptions_Defaults(t *testing.T) {
	opts := IOSOptions{}

	if opts.DisableInputAccessoryView != false {
		t.Error("DisableInputAccessoryView should default to false")
	}
	if opts.DisableScroll != false {
		t.Error("DisableScroll should default to false")
	}
	if opts.DisableBounce != false {
		t.Error("DisableBounce should default to false")
	}
	if opts.EnableBackForwardNavigationGestures != false {
		t.Error("EnableBackForwardNavigationGestures should default to false")
	}
	if opts.EnableNativeTabs != false {
		t.Error("EnableNativeTabs should default to false")
	}
}

func TestAndroidOptions_Defaults(t *testing.T) {
	opts := AndroidOptions{}

	if opts.DisableScroll != false {
		t.Error("DisableScroll should default to false")
	}
	if opts.DisableOverscroll != false {
		t.Error("DisableOverscroll should default to false")
	}
	if opts.EnableZoom != false {
		t.Error("EnableZoom should default to false")
	}
	if opts.DisableHardwareAcceleration != false {
		t.Error("DisableHardwareAcceleration should default to false")
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
	if opts.DisableLogging != false {
		t.Error("DisableLogging should default to false")
	}
}

func TestNativeTabItem_Fields(t *testing.T) {
	item := NativeTabItem{
		Title:       "Home",
		SystemImage: NativeTabIconHouse,
	}

	if item.Title != "Home" {
		t.Errorf("Title = %q, want %q", item.Title, "Home")
	}
	if item.SystemImage != NativeTabIconHouse {
		t.Errorf("SystemImage = %q, want %q", item.SystemImage, NativeTabIconHouse)
	}
}

func TestMiddleware_ShortCircuit(t *testing.T) {
	shortCircuit := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			// Don't call next
		})
	}

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("base"))
	})

	chained := ChainMiddleware(shortCircuit)
	handler := chained(baseHandler)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusForbidden)
	}
	if rec.Body.String() != "forbidden" {
		t.Errorf("Body = %q, want %q", rec.Body.String(), "forbidden")
	}
}
