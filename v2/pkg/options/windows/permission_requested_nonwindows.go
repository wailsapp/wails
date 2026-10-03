//go:build !windows

package windows

// PermissionRequestedCallback handles WebView2 permission requests.
type PermissionRequestedCallback func(uri string, kind uint32, uriErr error) uint32
