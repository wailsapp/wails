package application

import "testing"

func TestShouldNavigateWithoutPolicyAllows(t *testing.T) {
	window := &WebviewWindow{}
	if !window.shouldNavigate(NavigationRequest{URL: "https://example.com/"}) {
		t.Error("a window without OnBeforeNavigate should allow every navigation")
	}
}

func TestShouldNavigatePassesRequestAndWindow(t *testing.T) {
	var gotWindow Window
	var gotRequest NavigationRequest
	window := &WebviewWindow{
		options: WebviewWindowOptions{
			OnBeforeNavigate: func(w Window, request NavigationRequest) bool {
				gotWindow = w
				gotRequest = request
				return false
			},
		},
	}

	request := NavigationRequest{URL: "https://example.com/login", IsMainFrame: true}
	if window.shouldNavigate(request) {
		t.Error("shouldNavigate should return the policy's refusal")
	}
	if gotWindow != Window(window) {
		t.Error("the policy should be given the window being navigated")
	}
	if gotRequest != request {
		t.Errorf("the policy got %+v, want %+v", gotRequest, request)
	}
}

func TestInterceptsNavigation(t *testing.T) {
	withNavigationWindow(t, 1, func(Window, NavigationRequest) bool { return true })

	if !windowInterceptsNavigationByID(1) {
		t.Error("a window with OnBeforeNavigate should intercept navigation")
	}
	if windowInterceptsNavigationByID(2) {
		t.Error("an unknown window should not intercept navigation")
	}

	withNavigationWindow(t, 1, nil)
	if windowInterceptsNavigationByID(1) {
		t.Error("a window without OnBeforeNavigate should not intercept navigation")
	}
}

func TestDecideNavigation(t *testing.T) {
	sameOrigin := func(_ Window, request NavigationRequest) bool {
		return !request.IsMainFrame || request.URL == "wails://localhost/"
	}

	tests := []struct {
		name     string
		windowID uint
		policy   func(Window, NavigationRequest) bool
		request  NavigationRequest
		want     bool
	}{
		{
			name:     "allowed by the policy",
			windowID: 1,
			policy:   sameOrigin,
			request:  NavigationRequest{URL: "wails://localhost/", IsMainFrame: true},
			want:     true,
		},
		{
			name:     "refused by the policy",
			windowID: 1,
			policy:   sameOrigin,
			request:  NavigationRequest{URL: "https://example.com/", IsMainFrame: true},
			want:     false,
		},
		{
			name:     "an iframe is reported as such",
			windowID: 1,
			policy:   sameOrigin,
			request:  NavigationRequest{URL: "https://example.com/", IsMainFrame: false},
			want:     true,
		},
		{
			// The window may close while the webview is waiting on an answer.
			name:     "a window that has gone away is refused",
			windowID: 2,
			policy:   func(Window, NavigationRequest) bool { return true },
			request:  NavigationRequest{URL: "wails://localhost/", IsMainFrame: true},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withNavigationWindow(t, 1, tt.policy)

			if got := decideNavigation(tt.windowID, tt.request); got != tt.want {
				t.Errorf("decideNavigation(%+v) = %t, want %t", tt.request, got, tt.want)
			}
		})
	}
}

// A policy that panics must still give the webview an answer, and the answer
// must be a refusal.
func TestDecideNavigationRefusesWhenPolicyPanics(t *testing.T) {
	app := withNavigationWindow(t, 1, func(Window, NavigationRequest) bool {
		panic("policy failed")
	})
	panics := make(chan *PanicDetails, 1)
	app.options.PanicHandler = func(details *PanicDetails) { panics <- details }

	if decideNavigation(1, NavigationRequest{URL: "wails://localhost/", IsMainFrame: true}) {
		t.Error("a panicking policy should refuse the navigation")
	}
	select {
	case <-panics:
	default:
		t.Error("the panic should reach the application's PanicHandler")
	}
}

// withNavigationWindow installs an application holding one window with the
// given OnBeforeNavigate, and restores the previous one when the test ends.
func withNavigationWindow(t *testing.T, windowID uint, policy func(Window, NavigationRequest) bool) *App {
	t.Helper()

	previous := globalApplication

	app := &App{
		windows: map[uint]Window{
			windowID: &WebviewWindow{
				id:      windowID,
				options: WebviewWindowOptions{OnBeforeNavigate: policy},
			},
		},
	}
	app.Window = newWindowManager(app)

	globalApplication = app
	t.Cleanup(func() { globalApplication = previous })
	return app
}
