//go:build darwin && purego && !ios && !server

package application

// CGO-free port of the static part of WebviewNotchWindow
// (webview_notch_window_darwin.m): the shaped background view, target-screen
// selection, camera-housing placement and configureWebView.
//
// The slide animation (showAnimated:/hideAnimated:), the hover tracking area
// (mouseEntered:/mouseExited:) and the close override are in
// webview_notch_animation_darwin_purego.go.

import (
	"fmt"
	"math"
	"sync"

	"github.com/ebitengine/purego/objc"
)

// Keep these shape constants synchronized with NotchWingInset and
// NotchBottomCornerRadius in webview_notch_window_darwin.m.
const (
	notchBackgroundWing       = 20.0
	notchBackgroundBodyRadius = 28.0
)

// notchPathSegment is one step of the notch background path. The first
// segment is always a move; Curve selects a curve over a straight line.
type notchPathSegment struct {
	Curve  bool
	C1, C2 CGPoint // control points, used when Curve is set
	End    CGPoint
}

// notchBackgroundOutline returns the segments of the background path for a
// view of the given size, in drawing order, starting at the top-left corner.
// It is the path of WailsNotchBackgroundView's drawRect: in
// webview_notch_window_darwin.m, in AppKit's bottom-left origin.
func notchBackgroundOutline(width, height float64) []notchPathSegment {
	const (
		wing   = notchBackgroundWing
		radius = notchBackgroundBodyRadius
		body   = radius // bottom of the straight flanks
	)
	return []notchPathSegment{
		{End: CGPoint{X: 0, Y: height}},
		{End: CGPoint{X: width, Y: height}},
		{Curve: true,
			C1:  CGPoint{X: width - wing*0.58, Y: height},
			C2:  CGPoint{X: width - wing, Y: height - wing*0.42},
			End: CGPoint{X: width - wing, Y: height - wing}},
		{End: CGPoint{X: width - wing, Y: body}},
		{Curve: true,
			C1:  CGPoint{X: width - wing, Y: radius * 0.45},
			C2:  CGPoint{X: width - wing - radius*0.45, Y: 0},
			End: CGPoint{X: width - wing - radius, Y: 0}},
		{End: CGPoint{X: wing + radius, Y: 0}},
		{Curve: true,
			C1:  CGPoint{X: wing + radius*0.45, Y: 0},
			C2:  CGPoint{X: wing, Y: radius * 0.45},
			End: CGPoint{X: wing, Y: body}},
		{End: CGPoint{X: wing, Y: height - wing}},
		{Curve: true,
			C1:  CGPoint{X: wing, Y: height - wing*0.42},
			C2:  CGPoint{X: wing * 0.58, Y: height},
			End: CGPoint{X: 0, Y: height}},
	}
}

// notchCameraHousingCentre returns the horizontal centre of the camera housing
// in screen coordinates, from the screen's auxiliary top areas and safe-area
// top inset. found is false when the screen has no housing: no safe-area inset,
// areas unavailable (macOS before 12), an empty area, or no gap between them.
// It is the decision of cameraHousingCentreForScreen:found: without the screen.
func notchCameraHousingCentre(leftArea, rightArea NSRect, safeAreaTop float64, areasAvailable bool) (float64, bool) {
	if !areasAvailable || !(safeAreaTop > 0) {
		return 0, false
	}
	if nsIsEmptyRect(leftArea) || nsIsEmptyRect(rightArea) {
		return 0, false
	}
	left := leftArea.Origin.X + leftArea.Size.Width // NSMaxX(leftArea)
	right := rightArea.Origin.X                     // NSMinX(rightArea)
	if right <= left {
		return 0, false
	}
	return left + (right-left)/2, true
}

func nsIsEmptyRect(r NSRect) bool {
	return r.Size.Width <= 0 || r.Size.Height <= 0
}

// notchTargetFrame returns frame moved onto its target: horizontally centred on
// the housing (rounded to whole points like the native round()) and flush with
// the top of the reference rect, which is the screen frame when there is a
// housing and the visible frame otherwise. Without a housing the centre is the
// screen frame's. frame keeps its size.
func notchTargetFrame(frame, screenFrame, visibleFrame NSRect, housingCentre float64, hasHousing bool) NSRect {
	reference := visibleFrame
	centre := screenFrame.Origin.X + screenFrame.Size.Width/2 // NSMidX(screen.frame)
	if hasHousing {
		reference = screenFrame
		centre = housingCentre
	}
	frame.Origin.X = math.Round(centre - frame.Size.Width/2)
	frame.Origin.Y = reference.Origin.Y + reference.Size.Height - frame.Size.Height // NSMaxY - height
	return frame
}

// notchHiddenFrame returns the frame the slide-out animation moves the window
// to: the target moved up by its own height, fully above the top of the screen.
func notchHiddenFrame(target NSRect) NSRect {
	hidden := target
	hidden.Origin.Y = target.Origin.Y + target.Size.Height // NSMaxY(target)
	return hidden
}

// notchTargetScreen returns the screen whose ID matches screenID, else the
// primary screen (the first in [NSScreen screens]), else the main screen.
// IDs are formatted as processScreen formats Screen.ID.
func notchTargetScreen(screenID string) id {
	screens := class("NSScreen").send("screens")
	if screenID != "" {
		count := get[uint](screens, "count")
		for i := uint(0); i < count; i++ {
			screen := screens.send("objectAtIndex:", i)
			if notchScreenDeviceID(screen) == screenID {
				return screen
			}
		}
	}
	if primary := screens.send("firstObject"); !primary.isNil() {
		return primary
	}
	return class("NSScreen").send("mainScreen")
}

// notchScreenDeviceID is the NSScreenNumber of a screen formatted as Screen.ID.
func notchScreenDeviceID(screen id) string {
	number := screen.send("deviceDescription").send("objectForKey:", nsString("NSScreenNumber"))
	if number.isNil() {
		return ""
	}
	return fmt.Sprintf("%d", int32(get[uint32](number, "unsignedIntValue")))
}

// notchScreenHousingCentre reads the screen's camera housing. The auxiliary
// areas and safeAreaInsets are macOS 12+, so each is guarded by macOSAtLeast
// and respondsTo; cgo guards the same calls with @available(macOS 12.0, *).
func notchScreenHousingCentre(screen id) (float64, bool) {
	areasAvailable := macOSAtLeast(12, 0) &&
		respondsTo(screen, "auxiliaryTopLeftArea") &&
		respondsTo(screen, "auxiliaryTopRightArea") &&
		respondsTo(screen, "safeAreaInsets")
	if !areasAvailable {
		return notchCameraHousingCentre(NSRect{}, NSRect{}, 0, false)
	}
	return notchCameraHousingCentre(
		get[NSRect](screen, "auxiliaryTopLeftArea"),
		get[NSRect](screen, "auxiliaryTopRightArea"),
		get[NSEdgeInsets](screen, "safeAreaInsets").Top,
		true)
}

// notchFrameForScreen returns the target frame of a notch window of the given
// frame on screen. A nil screen leaves the frame unchanged.
func notchFrameForScreen(frame NSRect, screen id) NSRect {
	if screen.isNil() {
		return frame
	}
	centre, hasHousing := notchScreenHousingCentre(screen)
	return notchTargetFrame(frame,
		get[NSRect](screen, "frame"), get[NSRect](screen, "visibleFrame"),
		centre, hasHousing)
}

// notchBackgroundBezierPath returns the background outline as an NSBezierPath
// for a view of the given size, closed like drawRect:'s path.
func notchBackgroundBezierPath(width, height float64) id {
	segments := notchBackgroundOutline(width, height)
	path := class("NSBezierPath").send("bezierPath")
	for i, s := range segments {
		switch {
		case i == 0:
			path.send("moveToPoint:", s.End)
		case s.Curve:
			path.send("curveToPoint:controlPoint1:controlPoint2:", s.End, s.C1, s.C2)
		default:
			path.send("lineToPoint:", s.End)
		}
	}
	path.send("closePath")
	return path
}

var (
	registerNotchBackgroundOnce sync.Once
	notchBackgroundClass        id
	registerNotchWindowOnce     sync.Once
	notchWindowClass            id
)

// registerNotchBackgroundClass registers WailsNotchBackgroundView, an NSView
// that fills the notch outline with black in drawRect:.
func registerNotchBackgroundClass() id {
	registerNotchBackgroundOnce.Do(func() {
		methods := []objc.MethodDef{
			{Cmd: sel_("isOpaque"), Fn: func(self objc.ID, cmd objc.SEL) bool { return false }},
			{Cmd: sel_("drawRect:"), Fn: func(self objc.ID, cmd objc.SEL, dirtyRect NSRect) {
				// The default NSView drawRect: draws nothing, so the super call
				// cgo makes is omitted.
				bounds := get[NSRect](id(self), "bounds")
				class("NSColor").send("blackColor").send("setFill")
				notchBackgroundBezierPath(bounds.Size.Width, bounds.Size.Height).send("fill")
			}},
		}
		notchBackgroundClass = registerDelegateClass("WailsNotchBackgroundView", "NSView", nil, methods)
	})
	return notchBackgroundClass
}

// registerNotchWindowClass registers WailsWebviewNotchPanel, the NSPanel
// subclass of notch windows. It overrides performWindowDragWithEvent: (notch
// windows stay anchored to the camera housing: web drag regions and the Wails
// drag bridge may not move the panel) and the animation, hover and close
// methods of notchWindowMethods.
func registerNotchWindowClass() id {
	registerNotchWindowOnce.Do(func() {
		registerWebviewPanelClass() // the superclass must exist first
		methods := append([]objc.MethodDef{
			{Cmd: sel_("performWindowDragWithEvent:"), Fn: func(self objc.ID, cmd objc.SEL, event objc.ID) {}},
		}, notchWindowMethods()...)
		notchWindowClass = registerDelegateClass("WailsWebviewNotchPanel", "WailsWebviewNSPanel", nil, methods)
	})
	return notchWindowClass
}

// configureNotchWebView is the purego configureWebView:contentWidth:contentHeight:
// targetScreenID:. It disables the native drag paths (the high-level API owns
// placement), installs the shaped background behind the webview, sizes the
// webview to the content, makes it transparent and moves the window onto its
// target. win's content view must already hold webView.
func configureNotchWebView(win, webView id, cfg *notchWindowConfig) {
	notchWindowStates.Store(win.ptr(), &notchWindowState{screenID: cfg.screenID, webView: webView})
	win.send("setMovable:", false)
	win.send("setMovableByWindowBackground:", false)

	const (
		belowWindow = -1     // NSWindowBelow
		autoWidth   = 1 << 1 // NSViewWidthSizable
		autoHeight  = 1 << 4 // NSViewHeightSizable
	)
	root := win.send("contentView")
	background := registerNotchBackgroundClass().send("alloc").
		send("initWithFrame:", get[NSRect](root, "bounds"))
	background.send("setAutoresizingMask:", uint(autoWidth|autoHeight))
	// addSubview retains the background; drop the creation reference.
	root.send("addSubview:positioned:relativeTo:", background, belowWindow, webView)
	background.send("release")

	// Let web content reach the top edge of the display; the camera housing
	// masks the centre and the usable areas stay on either side.
	webView.send("setFrame:", rect(notchWindowHorizontalInset, notchWindowBottomInset,
		CGFloat(cfg.contentWidth), CGFloat(cfg.contentHeight)))
	webView.send("setAutoresizingMask:", uint(0)) // NSViewNotSizable
	wailsPrivateSetWebviewTransparent(webView)
	installNotchTrackingArea(win, root)

	win.send("setFrame:display:", notchFrameForScreen(get[NSRect](win, "frame"), notchTargetScreen(cfg.screenID)), false)
}
