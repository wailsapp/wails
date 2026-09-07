//go:build windows

package windows

import "github.com/wailsapp/go-webview2/pkg/edge"

// PermissionRequestedCallback handles WebView2 permission requests.
type PermissionRequestedCallback func(uri string, kind edge.CoreWebView2PermissionKind, uriErr error) edge.CoreWebView2PermissionState
