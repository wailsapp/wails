//go:build windows

package windows

import (
	"errors"
	"testing"

	"github.com/wailsapp/go-webview2/pkg/edge"
	windowsoptions "github.com/wailsapp/wails/v2/pkg/options/windows"
)

func TestFrontendWiresPermissionCallbackInsteadOfGlobalAllow(t *testing.T) {
	wantURI := "http://wails.localhost/console/"
	wantKind := edge.CoreWebView2PermissionKindMicrophone
	wantURIErr := errors.New("GetURI failed")
	wantState := edge.CoreWebView2PermissionStateDeny

	callbackCalled := false
	var permissionRequested windowsoptions.PermissionRequestedCallback = func(uri string, kind edge.CoreWebView2PermissionKind, uriErr error) edge.CoreWebView2PermissionState {
		callbackCalled = true
		if uri != wantURI {
			t.Errorf("permission callback URI = %q, want %q", uri, wantURI)
		}
		if kind != wantKind {
			t.Errorf("permission callback kind = %v, want %v", kind, wantKind)
		}
		if uriErr != wantURIErr {
			t.Errorf("permission callback URI error = %v, want %v", uriErr, wantURIErr)
		}
		return wantState
	}

	var installed windowsoptions.PermissionRequestedCallback
	globalAllowCalls := 0
	configurePermissionRequested(
		&windowsoptions.Options{PermissionRequested: permissionRequested},
		func(callback windowsoptions.PermissionRequestedCallback) {
			installed = callback
		},
		func(state edge.CoreWebView2PermissionState) {
			globalAllowCalls++
			if state != edge.CoreWebView2PermissionStateAllow {
				t.Errorf("global permission = %v, want Allow", state)
			}
		},
	)

	if installed == nil {
		t.Fatal("permission callback was not installed")
	}
	if got := installed(wantURI, wantKind, wantURIErr); got != wantState {
		t.Errorf("installed permission callback state = %v, want %v", got, wantState)
	}
	if !callbackCalled {
		t.Error("installed permission callback was not called")
	}
	if globalAllowCalls != 0 {
		t.Errorf("global Allow calls = %d, want 0", globalAllowCalls)
	}

	for _, test := range []struct {
		name           string
		windowsOptions *windowsoptions.Options
	}{
		{name: "nil Windows options"},
		{name: "nil callback", windowsOptions: &windowsoptions.Options{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			callbackInstalled := false
			var globalPermissions []edge.CoreWebView2PermissionState
			configurePermissionRequested(
				test.windowsOptions,
				func(windowsoptions.PermissionRequestedCallback) {
					callbackInstalled = true
				},
				func(state edge.CoreWebView2PermissionState) {
					globalPermissions = append(globalPermissions, state)
				},
			)

			if callbackInstalled {
				t.Error("permission callback was installed")
			}
			if len(globalPermissions) != 1 || globalPermissions[0] != edge.CoreWebView2PermissionStateAllow {
				t.Errorf("global permissions = %v, want [Allow]", globalPermissions)
			}
		})
	}
}
