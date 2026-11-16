//go:build darwin && purego && !ios && !server

package application

// CGO-free port of the animation, hover and close parts of WebviewNotchWindow
// (webview_notch_window_darwin.m): showAnimated:duration:, hideAnimated:
// duration:, mouseEntered:/mouseExited: (driven by the tracking area that
// configureNotchWebView installs) and the close override. The frame, timing and
// mode decisions are pure functions so they can be table-tested; the NSWindow
// and NSAnimationContext calls go through objc.

import (
	"sync"

	"github.com/ebitengine/purego/objc"
)

// AppKit tracking-area options (NSTrackingArea.h of the macOS 26 SDK).
const (
	nsTrackingMouseEnteredAndExited = 0x01  // NSTrackingMouseEnteredAndExited
	nsTrackingActiveAlways          = 0x80  // NSTrackingActiveAlways
	nsTrackingInVisibleRect         = 0x200 // NSTrackingInVisibleRect
)

// notchTrackingOptions are the tracking options of the notch hover area, as in
// configureWebView:.
const notchTrackingOptions = nsTrackingMouseEnteredAndExited | nsTrackingActiveAlways | nsTrackingInVisibleRect

// notchShowStartAlpha is the alpha a showing notch starts from; it settles to 1.
const notchShowStartAlpha = 0.96

// notchShowMode is what showAnimated:duration: does for a request.
type notchShowMode int

const (
	// notchShowAlreadyVisible: a show is already requested; the window is just
	// ordered front again.
	notchShowAlreadyVisible notchShowMode = iota
	// notchShowInstant: the frame is set and the window ordered front at once.
	notchShowInstant
	// notchShowSlide: the window slides in over the duration.
	notchShowSlide
)

// notchShowModeFor is the decision of showAnimated:duration: given whether a
// show is already requested.
func notchShowModeFor(requestedVisible, animated bool, duration float64) notchShowMode {
	switch {
	case requestedVisible:
		return notchShowAlreadyVisible
	case !animated || duration <= 0:
		return notchShowInstant
	default:
		return notchShowSlide
	}
}

// notchHideMode is what hideAnimated:duration: does for a request.
type notchHideMode int

const (
	// notchHideNoop: nothing is requested and nothing is on screen.
	notchHideNoop notchHideMode = iota
	// notchHideInstant: the window is ordered out at once.
	notchHideInstant
	// notchHideSlide: the window slides out over two thirds of the duration.
	notchHideSlide
)

// notchHideModeFor is the decision of hideAnimated:duration: given whether a
// show is requested and whether the window is visible.
func notchHideModeFor(requestedVisible, visible, animated bool, duration float64) notchHideMode {
	switch {
	case !requestedVisible && !visible:
		return notchHideNoop
	case !animated || duration <= 0:
		return notchHideInstant
	default:
		return notchHideSlide
	}
}

// notchHideDuration is the hide animation length: two thirds of the show
// duration, so dismissal is faster than the reveal.
func notchHideDuration(duration float64) float64 { return duration * (2.0 / 3.0) }

// notchAnimationPlan describes one slide animation: an optional frame and
// alpha set before it starts, then the end frame and alpha reached over
// Duration seconds along the cubic Bezier Timing (its two control points).
type notchAnimationPlan struct {
	// StartFrame and StartAlpha are set before the animation starts. The hide
	// plan leaves them zero: cgo starts it from the window's current frame and
	// alpha.
	StartFrame NSRect
	StartAlpha float64
	EndFrame   NSRect
	EndAlpha   float64
	Duration   float64
	Timing     [4]float64
}

// notchShowPlan is the slide-in: from the frame above the screen, at alpha 0.96
// to the target frame at alpha 1, over the show duration.
func notchShowPlan(target NSRect, duration float64) notchAnimationPlan {
	return notchAnimationPlan{
		StartFrame: notchHiddenFrame(target),
		StartAlpha: notchShowStartAlpha,
		EndFrame:   target,
		EndAlpha:   1.0,
		Duration:   duration,
		Timing:     [4]float64{0.22, 1.0, 0.36, 1.0},
	}
}

// notchHidePlan is the slide-out: from the target frame to the frame above the
// screen, fading to alpha 0.9, over two thirds of the duration.
func notchHidePlan(target NSRect, duration float64) notchAnimationPlan {
	return notchAnimationPlan{
		EndFrame: notchHiddenFrame(target),
		EndAlpha: 0.9,
		Duration: notchHideDuration(duration),
		Timing:   [4]float64{0.4, 0.0, 1.0, 1.0},
	}
}

// notchWindowState is the per-window state cgo keeps as WebviewNotchWindow
// properties: whether a show is requested, the generation that invalidates
// stale animation completions, the target screen and the webview to focus on
// hover. It is keyed by the NSWindow pointer.
type notchWindowState struct {
	requestedVisible bool
	generation       uint
	screenID         string
	webView          id
}

var notchWindowStates sync.Map // uintptr(NSWindow) -> *notchWindowState

// notchStateFor returns the state of win, creating an empty one if needed.
func notchStateFor(win id) *notchWindowState {
	state, _ := notchWindowStates.LoadOrStore(win.ptr(), &notchWindowState{})
	return state.(*notchWindowState)
}

// notchWindowMethods are the notch-specific overrides of WailsWebviewNotchPanel.
func notchWindowMethods() []objc.MethodDef {
	return []objc.MethodDef{
		{Cmd: sel_("showAnimated:duration:"), Fn: func(self objc.ID, cmd objc.SEL, animated bool, duration float64) {
			notchShowAnimated(id(self), animated, duration)
		}},
		{Cmd: sel_("hideAnimated:duration:"), Fn: func(self objc.ID, cmd objc.SEL, animated bool, duration float64) {
			notchHideAnimated(id(self), animated, duration)
		}},
		{Cmd: sel_("mouseEntered:"), Fn: func(self objc.ID, cmd objc.SEL, event objc.ID) {
			notchMouseEntered(id(self))
		}},
		// Keep the panel key until the user focuses another window: hover only
		// acquires focus, it does not dismiss.
		{Cmd: sel_("mouseExited:"), Fn: func(self objc.ID, cmd objc.SEL, event objc.ID) {}},
		{Cmd: sel_("close"), Fn: func(self objc.ID, cmd objc.SEL) {
			notchClose(id(self))
		}},
	}
}

// notchClose is the close override: it withdraws the show request and
// invalidates pending animations, then closes the window as the superclass
// does. The state entry is dropped, since the window is closed for good.
func notchClose(win id) {
	state := notchStateFor(win)
	state.requestedVisible = false
	state.generation++
	sendSuperFrom(notchWindowClass, win.raw(), "close")
	notchWindowStates.Delete(win.ptr())
}

// notchMouseEntered orders the notch front and makes it key so its webview
// takes keyboard input at once. The panel is non-activating, so the application
// stays inactive.
func notchMouseEntered(win id) {
	if !get[bool](win, "isVisible") {
		return
	}
	win.send("orderFrontRegardless")
	win.send("makeKeyWindow")
	if webView := notchStateFor(win).webView; !webView.isNil() {
		win.send("makeFirstResponder:", webView)
	}
}

// notchCurrentFrame returns the target frame of win on its configured screen.
func notchCurrentFrame(win id, state *notchWindowState) NSRect {
	return notchFrameForScreen(get[NSRect](win, "frame"), notchTargetScreen(state.screenID))
}

// notchShowAnimated is showAnimated:duration: of WebviewNotchWindow.
func notchShowAnimated(win id, animated bool, duration float64) {
	state := notchStateFor(win)
	mode := notchShowModeFor(state.requestedVisible, animated, duration)
	if mode == notchShowAlreadyVisible {
		win.send("orderFrontRegardless")
		return
	}

	state.requestedVisible = true
	state.generation++
	generation := state.generation
	target := notchCurrentFrame(win, state)
	if mode == notchShowInstant {
		win.send("setAlphaValue:", CGFloat(1.0))
		win.send("setFrame:display:", target, false)
		win.send("orderFrontRegardless")
		return
	}

	plan := notchShowPlan(target, duration)
	win.send("setFrame:display:", plan.StartFrame, false)
	win.send("setAlphaValue:", plan.StartAlpha)
	win.send("orderFrontRegardless")
	notchRunAnimation(win, plan, func() {
		if state.generation != generation || !state.requestedVisible {
			return
		}
		win.send("setFrame:display:", target, false)
		win.send("setAlphaValue:", CGFloat(1.0))
	})
}

// notchHideAnimated is hideAnimated:duration: of WebviewNotchWindow.
func notchHideAnimated(win id, animated bool, duration float64) {
	state := notchStateFor(win)
	mode := notchHideModeFor(state.requestedVisible, get[bool](win, "isVisible"), animated, duration)
	if mode == notchHideNoop {
		return
	}
	state.requestedVisible = false
	state.generation++
	generation := state.generation
	target := notchCurrentFrame(win, state)
	if mode == notchHideInstant {
		win.send("orderOut:", objc.ID(0))
		win.send("setFrame:display:", target, false)
		win.send("setAlphaValue:", CGFloat(1.0))
		return
	}

	plan := notchHidePlan(target, duration)
	notchRunAnimation(win, plan, func() {
		if state.generation != generation || state.requestedVisible {
			return
		}
		win.send("orderOut:", objc.ID(0))
		win.send("setFrame:display:", target, false)
		win.send("setAlphaValue:", CGFloat(1.0))
	})
}

// notchRunAnimation runs plan on win's animator inside an NSAnimationContext
// group and calls done when the group completes. The timing function takes
// four C floats, so the control points are passed as float32.
func notchRunAnimation(win id, plan notchAnimationPlan, done func()) {
	timing := class("CAMediaTimingFunction").send("functionWithControlPoints::::",
		float32(plan.Timing[0]), float32(plan.Timing[1]), float32(plan.Timing[2]), float32(plan.Timing[3]))
	animate := objc.NewBlock(func(_ objc.Block, context objc.ID) {
		ctx := id(context)
		ctx.send("setDuration:", plan.Duration)
		ctx.send("setTimingFunction:", timing)
		win.send("animator").send("setFrame:display:", plan.EndFrame, true)
		win.send("animator").send("setAlphaValue:", plan.EndAlpha)
	})
	completion := objc.NewBlock(func(objc.Block) { done() })
	class("NSAnimationContext").send("runAnimationGroup:completionHandler:", animate, completion)
	// AppKit copies both blocks it keeps; drop our creation references.
	animate.Release()
	completion.Release()
}

// installNotchTrackingArea adds the hover tracking area of configureWebView:
// to the content view root, with win as its owner. The content view retains
// the area, so the creation reference is dropped after adding it, as in cgo.
func installNotchTrackingArea(win, root id) {
	area := class("NSTrackingArea").send("alloc").send("initWithRect:options:owner:userInfo:",
		CGRect{}, uint(notchTrackingOptions), win, objc.ID(0))
	root.send("addTrackingArea:", area)
	area.send("release")
}

// showNotch and hideNotch send the animated selectors that cgo windowShow and
// windowHide send to a WebviewNotchWindow.
func (w *macosWebviewWindow) showNotch() {
	animated, duration := w.notchAnimation()
	runOnMain(func() { w.win().send("showAnimated:duration:", animated, duration) })
}

func (w *macosWebviewWindow) hideNotch() {
	animated, duration := w.notchAnimation()
	runOnMain(func() { w.win().send("hideAnimated:duration:", animated, duration) })
}
