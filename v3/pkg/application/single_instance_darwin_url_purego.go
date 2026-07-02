//go:build darwin && purego && !ios && !server

package application

// CGO-free re-implementation of single_instance_darwin_url.go.
//
// A force-launched second instance ("open -n URL") receives its URL as a
// kAEGetURL Apple Event which LaunchServices only delivers once the process has
// "finished launching". We therefore briefly spin an NSApplication run loop so
// the event can be dispatched, capture the URL, then stop the loop. A
// performSelector:afterDelay: timer guarantees the loop is bounded even when no
// URL event ever arrives.

import (
	"sync"

	"github.com/ebitengine/purego/objc"
)

// launchURLCaptureTimeout is the maximum time the second instance will wait
// for a kAEGetURL Apple Event before assuming no URL was launched.
const launchURLCaptureTimeout = 0.3

// Carbon / AppleEvent four-char-codes used by the URL capture handler.
const (
	kInternetEventClass = 0x4755524C // 'GURL'
	kAEGetURL           = 0x4755524C // 'GURL'
	keyDirectObject     = 0x2D2D2D2D // '----'
)

// NSApplicationActivationPolicyProhibited — the short-lived second instance
// must not get a Dock icon.
const nsApplicationActivationPolicyProhibited = 2

// NSEventTypeApplicationDefined — used for the synthetic wake-up event that
// makes [NSApp run] return promptly after [NSApp stop:].
const nsEventTypeApplicationDefined = 15

var (
	urlCaptureOnce  sync.Once
	urlCaptureClass id

	captureMu         sync.Mutex
	capturedLaunchURL string
)

// captureLaunchURL briefly runs an NSApplication event loop so that
// LaunchServices can deliver any pending kAEGetURL Apple Event (e.g. when this
// process was force-launched via "open -n URL"). Returns the URL string, or ""
// if none arrived within the timeout.
func captureLaunchURL() string {
	var result string
	withAutoreleasePool(func() {
		ensureURLCaptureClass()

		app := class("NSApplication").send("sharedApplication")
		// Run without a dock icon — this is a short-lived second instance.
		app.send("setActivationPolicy:", nsApplicationActivationPolicyProhibited)

		handler := urlCaptureClass.send("alloc").send("init")

		// Register the URL handler BEFORE [NSApp run] calls finishLaunching so
		// we catch the event the moment it is dispatched.
		mgr := class("NSAppleEventManager").send("sharedAppleEventManager")
		mgr.send("setEventHandler:andSelector:forEventClass:andEventID:",
			handler,
			sel_("handleGetURLEvent:withReplyEvent:"),
			uint32(kInternetEventClass),
			uint32(kAEGetURL),
		)

		// Safety-net timeout so the run loop never blocks indefinitely if no
		// URL event arrives. Scheduled on the main run loop that [NSApp run]
		// drives.
		handler.send("performSelector:withObject:afterDelay:",
			sel_("wailsCaptureTimeout:"),
			id(0),
			float64(launchURLCaptureTimeout),
		)

		captureMu.Lock()
		capturedLaunchURL = ""
		captureMu.Unlock()

		// [NSApp run] calls finishLaunching (signalling to LaunchServices that
		// this process is ready) and then enters the event loop. The run loop is
		// stopped by either the URL handler (early, on success) or the timeout.
		app.send("run")

		mgr.send("removeEventHandlerForEventClass:andEventID:",
			uint32(kInternetEventClass), uint32(kAEGetURL))

		captureMu.Lock()
		result = capturedLaunchURL
		captureMu.Unlock()
	})
	return result
}

// ensureURLCaptureClass registers the delegate class exactly once. The class
// exposes two instance methods: the Apple Event handler and the timeout target.
func ensureURLCaptureClass() {
	urlCaptureOnce.Do(func() {
		urlCaptureClass = registerDelegateClass(
			"WailsURLCaptureHandler_purego",
			"NSObject",
			nil,
			[]objc.MethodDef{
				{
					Cmd: sel_("handleGetURLEvent:withReplyEvent:"),
					Fn:  urlCaptureHandleGetURL,
				},
				{
					Cmd: sel_("wailsCaptureTimeout:"),
					Fn:  urlCaptureTimeout,
				},
			},
		)
	})
}

// urlCaptureHandleGetURL implements
// -handleGetURLEvent:withReplyEvent: — it extracts the URL string from the
// Apple Event descriptor and stops the run loop.
func urlCaptureHandleGetURL(self objc.ID, _ objc.SEL, event objc.ID, _ objc.ID) {
	desc := id(event).send("paramDescriptorForKeyword:", uint32(keyDirectObject))
	urlStr := desc.send("stringValue")
	if !urlStr.isNil() {
		captureMu.Lock()
		if capturedLaunchURL == "" {
			capturedLaunchURL = urlStr.string()
		}
		captureMu.Unlock()
	}
	stopAppEventLoop()
}

// urlCaptureTimeout implements -wailsCaptureTimeout: — the safety-net timer that
// stops the run loop if no URL event arrived.
func urlCaptureTimeout(_ objc.ID, _ objc.SEL, _ objc.ID) {
	stopAppEventLoop()
}

// stopAppEventLoop stops [NSApp run] by requesting a stop and posting a
// synthetic event so the run loop wakes and returns immediately.
func stopAppEventLoop() {
	app := class("NSApplication").send("sharedApplication")
	app.send("stop:", id(0))
	evt := class("NSEvent").send(
		"otherEventWithType:location:modifierFlags:timestamp:windowNumber:context:subtype:data1:data2:",
		uint(nsEventTypeApplicationDefined), // type
		CGPoint{X: 0, Y: 0},                 // location (NSZeroPoint)
		uint(0),                             // modifierFlags
		float64(0),                          // timestamp (NSTimeInterval)
		int(0),                              // windowNumber (NSInteger)
		id(0),                               // context
		int16(0),                            // subtype (short)
		int(0),                              // data1 (NSInteger)
		int(0),                              // data2 (NSInteger)
	)
	app.send("postEvent:atStart:", evt, false)
}
