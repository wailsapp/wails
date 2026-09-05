//go:build server

package application

// Native embedded webviews are unavailable in server deployments.
func newPanelImpl(_ *WebviewPanel) webviewPanelImpl { return nil }
