//go:build android && cgo && !server

package application

// AndroidSetScrollEnabled enables or disables WebView scrolling on Android.
func AndroidSetScrollEnabled(enabled bool) { androidSetScrollEnabled(enabled) }

// AndroidSetBounceEnabled enables or disables overscroll bounce on Android.
func AndroidSetBounceEnabled(enabled bool) { androidSetBounceEnabled(enabled) }

// AndroidSetScrollIndicatorsEnabled enables or disables scroll indicators on Android.
func AndroidSetScrollIndicatorsEnabled(enabled bool) {
	androidSetScrollIndicatorsEnabled(enabled)
}

// AndroidSetBackForwardGesturesEnabled enables or disables swipe back/forward gestures.
func AndroidSetBackForwardGesturesEnabled(enabled bool) {
	androidSetBackForwardGesturesEnabled(enabled)
}

// AndroidSetLinkPreviewEnabled enables or disables link preview behavior on Android.
func AndroidSetLinkPreviewEnabled(enabled bool) { androidSetLinkPreviewEnabled(enabled) }

// AndroidSetCustomUserAgent sets the custom user agent for the Android WebView.
func AndroidSetCustomUserAgent(ua string) { androidSetCustomUserAgent(ua) }
