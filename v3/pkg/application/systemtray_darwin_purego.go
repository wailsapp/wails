//go:build darwin && purego && !ios && !server

// Package application - CGO-free macOS system tray backend.
//
// This is the purego counterpart of systemtray_darwin.go / systemtray_darwin.m.
// It drives NSStatusBar / NSStatusItem / NSMenu directly through the
// Objective-C runtime helpers in darwin_purego_cocoa.go instead of compiling
// Objective-C through cgo.
//
// The behaviour mirrors the cgo backend exactly:
//
//   - a variable-length NSStatusItem is created on the system status bar,
//   - a per-tray controller object receives the button's target/action click,
//   - a local NSEvent monitor (installed as an Objective-C block) fires BEFORE
//     the button processes the mouse-down so the framework can, when
//     appropriate, hand the click to native menu tracking (proper highlight,
//     no app activation),
//   - programmatic OpenMenu() synthesizes a mouse-down to enter native menu
//     tracking.
package application

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego/objc"
)

// ---------------------------------------------------------------------------
// AppKit constants (values taken from the AppKit headers so we don't need to
// dlsym the exported symbols).
// ---------------------------------------------------------------------------

const (
	nsVariableStatusItemLength = -1.0 // NSVariableStatusItemLength

	nsEventMaskLeftMouseDown  = 1 << 1 // NSEventMaskLeftMouseDown  (== 2)
	nsEventMaskRightMouseDown = 1 << 3 // NSEventMaskRightMouseDown (== 8)

	nsEventTypeLeftMouseDown = 1 // NSEventTypeLeftMouseDown

	nsPopUpMenuWindowLevel = 101 // NSPopUpMenuWindowLevel

	// NSAttributedString attribute-name string values. These are the literal
	// NSString values of NSForegroundColorAttributeName /
	// NSBackgroundColorAttributeName, stable across macOS releases.
	nsForegroundColorAttributeName = "NSColor"
	nsBackgroundColorAttributeName = "NSBackgroundColor"
)

// ---------------------------------------------------------------------------
// macosSystemTray - field names identical to the cgo implementation so other
// files (e.g. screen_darwin_purego.go's getScreenForSystray) link against it.
// Native handles are stored as unsafe.Pointer, matching the cgo struct and the
// convention used by the other purego files; they are converted to `id` at use
// sites via idFromPtr / ptrFromID.
// ---------------------------------------------------------------------------

type macosSystemTray struct {
	id    uint
	label string
	icon  []byte
	menu  *Menu

	nsStatusItem      unsafe.Pointer
	nsImage           unsafe.Pointer
	nsMenu            unsafe.Pointer
	iconPosition      IconPosition
	isTemplateIcon    bool
	parent            *SystemTray
	lastClickedScreen unsafe.Pointer

	// purego-only bookkeeping (no cgo equivalent field). nsController is the
	// per-tray controller object that owns the target/action + menu delegate;
	// eventMonitor is the local NSEvent monitor object returned by AppKit.
	nsController unsafe.Pointer
	eventMonitor unsafe.Pointer
}

// button mirrors the cgo enum: the raw NSEventType values for the two mouse
// buttons we handle.
type button int

const (
	leftButtonDown  button = 1
	rightButtonDown button = 3
)

// system tray map (keyed by Wails system-tray id), mirroring the cgo backend.
var systemTrayMap = make(map[uint]*macosSystemTray)

// ---------------------------------------------------------------------------
// Controller class + click callbacks
// ---------------------------------------------------------------------------

var (
	statusItemControllerOnce  sync.Once
	statusItemControllerClass id

	systrayControllerMu  sync.Mutex
	systrayControllerMap = map[uintptr]*macosSystemTray{}
)

// statusItemControllerClassRef lazily registers the controller class used as
// the NSStatusItem target and NSMenu delegate. A single class backs every tray;
// per-instance state is resolved by looking the controller object up in
// systrayControllerMap.
func statusItemControllerClassRef() id {
	statusItemControllerOnce.Do(func() {
		statusItemControllerClass = registerDelegateClass(
			"WailsStatusItemController_purego", "NSObject", nil,
			[]objc.MethodDef{
				{Cmd: sel_("statusItemClicked:"), Fn: statusItemClickedIMP},
				{Cmd: sel_("menuDidClose:"), Fn: menuDidCloseIMP},
			},
		)
	})
	return statusItemControllerClass
}

func lookupSystrayController(self objc.ID) *macosSystemTray {
	systrayControllerMu.Lock()
	defer systrayControllerMu.Unlock()
	return systrayControllerMap[uintptr(self)]
}

// statusItemClickedIMP backs -[WailsStatusItemController statusItemClicked:].
// It reads the current NSEvent's type (left/right mouse-down) and forwards it,
// mirroring the cgo statusItemClicked:/systrayClickCallback path.
func statusItemClickedIMP(self objc.ID, cmd objc.SEL, sender objc.ID) {
	s := lookupSystrayController(self)
	if s == nil {
		return
	}
	event := class("NSApplication").send("sharedApplication").send("currentEvent")
	evType := int(get[uint](event, "type"))
	systrayClickCallback(s.id, evType)
}

// menuDidCloseIMP backs the NSMenuDelegate menuDidClose:. It detaches the menu
// from the status item so subsequent clicks invoke the action handler again.
func menuDidCloseIMP(self objc.ID, cmd objc.SEL, menu objc.ID) {
	s := lookupSystrayController(self)
	if s == nil {
		return
	}
	if s.nsStatusItem != nil {
		idFromPtr(s.nsStatusItem).send("setMenu:", id(0))
	}
	id(menu).send("setDelegate:", id(0))
}

// systrayClickCallback is the pure-Go analog of the cgo //export of the same
// name: it dispatches a processed click to the addressed system tray.
func systrayClickCallback(trayID uint, buttonID int) {
	systemTray := systemTrayMap[trayID]
	if systemTray == nil {
		globalApplication.error("system tray not found: %v", trayID)
		return
	}
	systemTray.processClick(button(buttonID))
}

// systrayPreClickCallback is the pure-Go analog of the cgo //export: it is
// called from the local NSEvent monitor BEFORE the button processes the
// mouse-down. It returns 1 when the framework should show the menu via native
// tracking, or 0 to let the action handler fire for custom click behaviour.
func systrayPreClickCallback(trayID uint, buttonID int) int {
	systemTray := systemTrayMap[trayID]
	if systemTray == nil || systemTray.nsMenu == nil {
		return 0
	}
	b := button(buttonID)
	switch b {
	case leftButtonDown:
		if systemTray.parent.clickHandler == nil &&
			systemTray.parent.attachedWindow.Window == nil {
			return 1
		}
	case rightButtonDown:
		if systemTray.parent.rightClickHandler == nil {
			// Hide the attached window before the menu appears.
			if systemTray.parent.attachedWindow.Window != nil &&
				systemTray.parent.attachedWindow.Window.IsVisible() {
				systemTray.parent.attachedWindow.Window.Hide()
			}
			return 1
		}
	}
	return 0
}

// ---------------------------------------------------------------------------
// Construction
// ---------------------------------------------------------------------------

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	result := &macosSystemTray{
		parent:         s,
		id:             s.id,
		label:          s.label,
		icon:           s.icon,
		menu:           s.menu,
		iconPosition:   s.iconPosition,
		isTemplateIcon: s.isTemplateIcon,
	}
	systemTrayMap[s.id] = result
	return result
}

// createStatusItem builds the NSStatusItem, its controller (target/action) and
// installs the pre-click local event monitor. It is the purego analog of the
// cgo systemTrayNew and must run on the main thread.
func (s *macosSystemTray) createStatusItem() id {
	controller := statusItemControllerClassRef().send("alloc").send("init")
	s.nsController = ptrFromID(controller)

	systrayControllerMu.Lock()
	systrayControllerMap[uintptr(controller.ptr())] = s
	systrayControllerMu.Unlock()

	statusBar := class("NSStatusBar").send("systemStatusBar")
	statusItem := statusBar.send("statusItemWithLength:", float64(nsVariableStatusItemLength)).send("retain")

	statusItem.send("setTarget:", controller)
	statusItem.send("setAction:", sel_("statusItemClicked:"))

	button := statusItem.send("button")
	button.send("sendActionOn:", uint(nsEventMaskLeftMouseDown|nsEventMaskRightMouseDown))

	// Install a local event monitor that fires BEFORE the button processes the
	// mouse-down. When the pre-click callback says "show menu" we temporarily
	// set statusItem.menu so the button enters native menu tracking.
	handler := objc.NewBlock(func(_ objc.Block, event objc.ID) objc.ID {
		ev := id(event)
		if ev.send("window").ptr() != button.send("window").ptr() {
			return event
		}
		action := systrayPreClickCallback(s.id, int(get[uint](ev, "type")))
		if action == 1 && s.nsMenu != nil {
			menu := idFromPtr(s.nsMenu)
			menu.send("setDelegate:", controller)
			statusItem.send("setMenu:", menu)
		}
		return event
	})
	monitor := class("NSEvent").send("addLocalMonitorForEventsMatchingMask:handler:",
		uint(nsEventMaskLeftMouseDown|nsEventMaskRightMouseDown), objc.ID(handler))
	s.eventMonitor = ptrFromID(monitor)

	return statusItem
}

func (s *macosSystemTray) run() {
	globalApplication.dispatchOnMainThread(func() {
		if s.nsStatusItem != nil {
			Fatal("System tray '%d' already running", s.id)
		}
		s.nsStatusItem = ptrFromID(s.createStatusItem())

		if s.label != "" {
			s.setLabel(s.label)
		}
		if s.icon != nil {
			s.applyIcon(s.icon)
		}
		if s.menu != nil {
			s.menu.Update()
			// Convert impl to macosMenu object; s.nsMenu doubles as the cached
			// menu read by the pre-click event monitor.
			s.nsMenu = (s.menu.impl).(*macosMenu).nsMenu
		}
	})
}

// ---------------------------------------------------------------------------
// Label handling (pure Go / objc - no C strings)
// ---------------------------------------------------------------------------

func (s *macosSystemTray) setLabel(label string) {
	s.label = label
	if s.nsStatusItem == nil {
		return
	}
	button := idFromPtr(s.nsStatusItem).send("button")

	if !hasANSICodes(label) {
		button.send("setTitle:", nsString(label))
		return
	}
	parts, err := SystemTrayLabelParser(label)
	if err != nil || len(parts) == 0 {
		button.send("setTitle:", nsString(label))
		return
	}

	l, fg, bg := partToStrings(parts[0])
	attr := createAttributedString(l, fg, bg)
	for _, p := range parts[1:] {
		l, fg, bg = partToStrings(p)
		attr = appendAttributedString(attr, l, fg, bg)
	}
	idFromPtr(s.nsStatusItem).send("setAttributedTitle:", attr)
}

func hasANSICodes(s string) bool {
	return strings.Contains(s, "\033[")
}

// partToStrings is the pure-Go analog of the cgo partToCStrings: it returns the
// text/foreground/background components of a label part as plain Go strings
// (no C allocation, no manual freeing).
func partToStrings(p SystemTrayLabelPart) (label, fg, bg string) {
	return p.Text, p.FgColor, p.BgColor
}

// createAttributedString builds an autoreleased NSMutableAttributedString with
// optional foreground/background colours, mirroring the cgo createAttributedString.
func createAttributedString(title, fg, bg string) id {
	dict := class("NSMutableDictionary").send("dictionary")
	if c := parseHexColor(fg); !c.isNil() {
		dict.send("setObject:forKey:", c, nsString(nsForegroundColorAttributeName))
	}
	if c := parseHexColor(bg); !c.isNil() {
		dict.send("setObject:forKey:", c, nsString(nsBackgroundColorAttributeName))
	}
	return class("NSMutableAttributedString").
		send("alloc").
		send("initWithString:attributes:", nsString(title), dict).
		send("autorelease")
}

// appendAttributedString appends a new styled run to current and returns the
// (possibly new) combined string, mirroring the cgo appendAttributedString.
func appendAttributedString(current id, title, fg, bg string) id {
	newString := createAttributedString(title, fg, bg)
	if !current.isNil() {
		current.send("appendAttributedString:", newString)
		return current
	}
	return newString
}

// parseHexColor parses "#rrggbb" or "#rrggbbaa" into an NSColor. It returns a
// nil id when the string is empty or contains no parseable colour component,
// matching the cgo behaviour (missing components default to 255).
func parseHexColor(hex string) id {
	if hex == "" {
		return id(0)
	}
	h := strings.TrimPrefix(hex, "#")
	// default white, fully opaque (r=g=b=a=255)
	comps := [4]float64{255, 255, 255, 255}
	parsed := 0
	for i := 0; i < 4 && len(h) >= (i+1)*2; i++ {
		v, err := strconv.ParseUint(h[i*2:i*2+2], 16, 16)
		if err != nil {
			break
		}
		comps[i] = float64(v)
		parsed++
	}
	if parsed == 0 {
		return id(0)
	}
	return class("NSColor").send("colorWithCalibratedRed:green:blue:alpha:",
		comps[0]/255.0, comps[1]/255.0, comps[2]/255.0, comps[3]/255.0)
}

// ---------------------------------------------------------------------------
// Icon handling
// ---------------------------------------------------------------------------

// imageFromBytes builds an NSImage from raw image bytes, mirroring the cgo
// imageFromBytes helper.
func imageFromBytes(b []byte) id {
	return class("NSImage").send("alloc").send("initWithData:", nsData(b))
}

// applyIcon renders the icon onto the status item button. Must run on the main
// thread.
func (s *macosSystemTray) applyIcon(icon []byte) {
	if s.nsStatusItem == nil || len(icon) == 0 {
		return
	}
	image := imageFromBytes(icon)
	s.nsImage = ptrFromID(image)

	thickness := get[CGFloat](class("NSStatusBar").send("systemStatusBar"), "thickness")
	image.send("setSize:", CGSize{Width: thickness, Height: thickness})
	if s.isTemplateIcon {
		image.send("setTemplate:", true)
	}
	button := idFromPtr(s.nsStatusItem).send("button")
	button.send("setImage:", image.send("autorelease"))
	button.send("setImagePosition:", int(s.iconPosition))
}

func (s *macosSystemTray) setIcon(icon []byte) {
	s.icon = icon
	globalApplication.dispatchOnMainThread(func() {
		s.applyIcon(icon)
	})
}

func (s *macosSystemTray) setDarkModeIcon(icon []byte) {
	s.setIcon(icon)
}

func (s *macosSystemTray) setTemplateIcon(icon []byte) {
	s.icon = icon
	s.isTemplateIcon = true
	globalApplication.dispatchOnMainThread(func() {
		s.applyIcon(icon)
	})
}

func (s *macosSystemTray) setIconPosition(position IconPosition) {
	s.iconPosition = position
}

func (s *macosSystemTray) setTooltip(tooltip string) {
	// Tooltips not supported on macOS
}

// ---------------------------------------------------------------------------
// Menu
// ---------------------------------------------------------------------------

func (s *macosSystemTray) setMenu(menu *Menu) {
	s.menu = menu
	if s.nsStatusItem != nil && menu != nil {
		menu.Update()
		s.nsMenu = (menu.impl).(*macosMenu).nsMenu
	}
}

func (s *macosSystemTray) openMenu() {
	if s.nsMenu == nil {
		return
	}
	s.showMenu()
}

// showMenu programmatically enters native menu tracking by synthesizing a
// mouse-down at the button centre, mirroring the cgo showMenu. Click-triggered
// menus are handled by the pre-click event monitor instead.
func (s *macosSystemTray) showMenu() {
	globalApplication.dispatchOnMainThread(func() {
		if s.nsStatusItem == nil || s.nsMenu == nil {
			return
		}
		statusItem := idFromPtr(s.nsStatusItem)
		menu := idFromPtr(s.nsMenu)
		controller := idFromPtr(s.nsController)
		button := statusItem.send("button")

		// Temporarily assign the menu for native tracking.
		menu.send("setDelegate:", controller)
		statusItem.send("setMenu:", menu)

		// Synthesize a mouse-down at the button centre.
		bounds := get[NSRect](button, "bounds")
		frame := get[NSRect](button, "convertRect:toView:", bounds, objc.ID(0))
		loc := NSPoint{
			X: frame.Origin.X + frame.Size.Width/2,
			Y: frame.Origin.Y + frame.Size.Height/2,
		}
		uptime := get[float64](class("NSProcessInfo").send("processInfo"), "systemUptime")
		windowNumber := get[int](button.send("window"), "windowNumber")

		event := class("NSEvent").send(
			"mouseEventWithType:location:modifierFlags:timestamp:windowNumber:context:eventNumber:clickCount:pressure:",
			uint(nsEventTypeLeftMouseDown),
			loc,
			uint(0),
			uptime,
			windowNumber,
			objc.ID(0),
			int(0),
			int(1),
			float64(1.0),
		)
		button.send("mouseDown:", event)

		// Menu dismissed - restore custom click handling.
		statusItem.send("setMenu:", id(0))
		menu.send("setDelegate:", id(0))
	})
}

// ---------------------------------------------------------------------------
// Geometry
// ---------------------------------------------------------------------------

func (s *macosSystemTray) getScreen() (*Screen, error) {
	if s.lastClickedScreen != nil {
		frame := get[NSRect](idFromPtr(s.lastClickedScreen), "frame")
		result := &Screen{
			Bounds: Rect{
				X:      int(frame.Origin.X),
				Y:      int(frame.Origin.Y),
				Width:  int(frame.Size.Width),
				Height: int(frame.Size.Height),
			},
		}
		return result, nil
	}
	return nil, errors.New("no screen available")
}

func (s *macosSystemTray) bounds() (*Rect, error) {
	if s.nsStatusItem == nil {
		return nil, errors.New("system tray not running")
	}
	statusItem := idFromPtr(s.nsStatusItem)
	button := statusItem.send("button")

	// Find the screen the mouse is currently on.
	mouseLocation := get[NSPoint](class("NSEvent"), "mouseLocation")
	var screen id
	screens := class("NSScreen").send("screens")
	count := int(get[uint](screens, "count"))
	for i := 0; i < count; i++ {
		candidate := screens.send("objectAtIndex:", uint(i))
		frame := get[NSRect](candidate, "frame")
		if pointInRect(mouseLocation, frame) {
			screen = candidate
			break
		}
	}
	if screen.isNil() {
		screen = class("NSScreen").send("mainScreen")
	}
	// Store the screen for use in positionWindow / getScreen.
	s.lastClickedScreen = ptrFromID(screen)

	// Button frame in screen coordinates.
	buttonFrame := get[NSRect](button, "frame")
	buttonFrameScreen := get[NSRect](button.send("window"), "convertRectToScreen:", buttonFrame)

	result := &Rect{
		X:      int(buttonFrameScreen.Origin.X),
		Y:      int(buttonFrameScreen.Origin.Y),
		Width:  int(buttonFrameScreen.Size.Width),
		Height: int(buttonFrameScreen.Size.Height),
	}
	return result, nil
}

// pointInRect is the Go reimplementation of NSPointInRect.
func pointInRect(p NSPoint, r NSRect) bool {
	return p.X >= r.Origin.X && p.X < r.Origin.X+r.Size.Width &&
		p.Y >= r.Origin.Y && p.Y < r.Origin.Y+r.Size.Height
}

func (s *macosSystemTray) positionWindow(window Window, offset int) error {
	nativeWindow := window.NativeWindow()
	if nativeWindow == nil {
		return errors.New("window native implementation unavailable")
	}
	win := idFromPtr(nativeWindow)

	button := idFromPtr(s.nsStatusItem).send("button")
	frame := get[NSRect](button.send("window"), "convertRectToScreen:", get[NSRect](button, "frame"))

	screen := button.send("window").send("screen")
	if screen.isNil() {
		screen = class("NSScreen").send("mainScreen")
	}
	scaleFactor := get[CGFloat](screen, "backingScaleFactor")

	windowFrame := get[NSRect](win, "frame")
	screenFrame := get[NSRect](screen, "frame")
	visibleFrame := get[NSRect](screen, "visibleFrame")

	// Horizontal position (centred under the status item), clamped to screen.
	windowX := frame.Origin.X + (frame.Size.Width-windowFrame.Size.Width)/2
	if windowX+windowFrame.Size.Width > screenFrame.Origin.X+screenFrame.Size.Width {
		windowX = screenFrame.Origin.X + screenFrame.Size.Width - windowFrame.Size.Width
	}
	if windowX < screenFrame.Origin.X {
		windowX = screenFrame.Origin.X
	}

	// Vertical position.
	scaledOffset := float64(offset) * scaleFactor
	windowY := visibleFrame.Origin.Y + visibleFrame.Size.Height - windowFrame.Size.Height - scaledOffset

	windowFrame.Origin.X = windowX
	windowFrame.Origin.Y = windowY
	win.send("setFrame:display:animate:", windowFrame, true, false)
	win.send("setLevel:", int(nsPopUpMenuWindowLevel))
	win.send("orderFrontRegardless")

	return nil
}

// ---------------------------------------------------------------------------
// Visibility / lifecycle
// ---------------------------------------------------------------------------

func (s *macosSystemTray) Show() {
	if s.nsStatusItem == nil {
		return
	}
	globalApplication.dispatchOnMainThread(func() {
		idFromPtr(s.nsStatusItem).send("setVisible:", true)
	})
}

func (s *macosSystemTray) Hide() {
	if s.nsStatusItem == nil {
		return
	}
	globalApplication.dispatchOnMainThread(func() {
		idFromPtr(s.nsStatusItem).send("setVisible:", false)
	})
}

func (s *macosSystemTray) destroy() {
	if s.nsStatusItem == nil {
		return
	}
	globalApplication.dispatchOnMainThread(func() {
		statusItem := idFromPtr(s.nsStatusItem)

		if s.eventMonitor != nil {
			class("NSEvent").send("removeMonitor:", idFromPtr(s.eventMonitor))
			s.eventMonitor = nil
		}
		class("NSStatusBar").send("systemStatusBar").send("removeStatusItem:", statusItem)

		if s.nsController != nil {
			systrayControllerMu.Lock()
			delete(systrayControllerMap, uintptr(idFromPtr(s.nsController).ptr()))
			systrayControllerMu.Unlock()
			idFromPtr(s.nsController).send("release")
			s.nsController = nil
		}
		statusItem.send("release")
		s.nsStatusItem = nil
	})
}

// ---------------------------------------------------------------------------
// Click processing (identical logic to the cgo backend)
// ---------------------------------------------------------------------------

func (s *macosSystemTray) processClick(b button) {
	switch b {
	case leftButtonDown:
		if s.parent.clickHandler != nil {
			s.parent.clickHandler()
			return
		}
		if s.parent.attachedWindow.Window != nil {
			s.parent.defaultClickHandler()
			return
		}
		if s.menu != nil {
			s.showMenu()
		}
	case rightButtonDown:
		if s.parent.rightClickHandler != nil {
			s.parent.rightClickHandler()
			return
		}
		if s.menu != nil {
			if s.parent.attachedWindow.Window != nil {
				s.parent.attachedWindow.Window.Hide()
			}
			s.showMenu()
			return
		}
		if s.parent.attachedWindow.Window != nil {
			s.parent.defaultClickHandler()
		}
	}
}
