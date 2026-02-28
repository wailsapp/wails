//go:build !android || !cgo || server

package application

// AndroidSetScrollEnabled is a no-op on non-Android platforms.
func AndroidSetScrollEnabled(enabled bool) {}

// AndroidSetBounceEnabled is a no-op on non-Android platforms.
func AndroidSetBounceEnabled(enabled bool) {}

// AndroidSetScrollIndicatorsEnabled is a no-op on non-Android platforms.
func AndroidSetScrollIndicatorsEnabled(enabled bool) {}

// AndroidSetBackForwardGesturesEnabled is a no-op on non-Android platforms.
func AndroidSetBackForwardGesturesEnabled(enabled bool) {}

// AndroidSetLinkPreviewEnabled is a no-op on non-Android platforms.
func AndroidSetLinkPreviewEnabled(enabled bool) {}

// AndroidSetCustomUserAgent is a no-op on non-Android platforms.
func AndroidSetCustomUserAgent(ua string) {}
