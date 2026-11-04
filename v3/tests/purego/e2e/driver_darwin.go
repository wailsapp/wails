//go:build darwin

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// The darwin driver posts NSEvents to this process's NSApplication with
// -[NSApplication postEvent:atStart:]. The probe runs inside the app, so the
// events go through the same queue as hardware input and reach the same
// responder chain, without accessibility or TCC permission (CGEventPost would
// need the Input Monitoring / Accessibility grant).
//
// Every AppKit call runs on the main thread. The probe calls the driver from a
// background goroutine while the main thread runs [NSApp run]; a call is
// scheduled with CFRunLoopPerformBlock in the common modes, so it also runs
// inside a modal session (NSAlert, NSOpenPanel) or menu tracking, and the
// caller waits on a channel, never the main thread.
//
// Key codes are the kVK_* constants of HIToolbox/Events.h, modifier flags and
// event types are NSEvent.h, and the function-key characters are the NSEvent.h
// NSUpArrowFunctionKey ... NSPageDownFunctionKey values.

const (
	nsEventTypeLeftMouseDown  = 1
	nsEventTypeLeftMouseUp    = 2
	nsEventTypeRightMouseDown = 3
	nsEventTypeRightMouseUp   = 4
	nsEventTypeKeyDown        = 10
	nsEventTypeKeyUp          = 11

	flagCommand  uint = 1 << 20
	flagShift    uint = 1 << 17
	flagControl  uint = 1 << 18
	flagOption   uint = 1 << 19
	flagFunction uint = 1 << 23 // set by AppKit on arrows, F-keys, home/end, page keys

	// mainThreadTimeout bounds a driver call. A call never waits forever for
	// a main thread that is not running its run loop.
	mainThreadTimeout = 20 * time.Second

	// clickHold is the gap between mouse down and mouse up, so the web view
	// sees a press and a release, as with a real click.
	clickHold = 50 * time.Millisecond
)

// keyStroke is one key press as AppKit reports it in a key NSEvent.
type keyStroke struct {
	code     uint16 // virtual key code (kVK_*)
	flags    uint   // NSEventModifierFlags
	chars    string // characters, with modifiers applied
	ignoring string // charactersIgnoringModifiers
}

// namedKey is a non-printable key such as Return or an arrow.
type namedKey struct {
	code     uint16
	chars    string
	function bool
}

// keyCodes are the virtual key codes of the printable keys on the ANSI US
// layout, by unshifted character (kVK_ANSI_*, kVK_Space, kVK_Grave).
var keyCodes = map[rune]uint16{
	'a': 0x00, 's': 0x01, 'd': 0x02, 'f': 0x03, 'h': 0x04, 'g': 0x05,
	'z': 0x06, 'x': 0x07, 'c': 0x08, 'v': 0x09, 'b': 0x0B, 'q': 0x0C,
	'w': 0x0D, 'e': 0x0E, 'r': 0x0F, 'y': 0x10, 't': 0x11,
	'1': 0x12, '2': 0x13, '3': 0x14, '4': 0x15, '6': 0x16, '5': 0x17,
	'=': 0x18, '9': 0x19, '7': 0x1A, '-': 0x1B, '8': 0x1C, '0': 0x1D,
	']': 0x1E, 'o': 0x1F, 'u': 0x20, '[': 0x21, 'i': 0x22, 'p': 0x23,
	'l': 0x25, 'j': 0x26, '\'': 0x27, 'k': 0x28, ';': 0x29, '\\': 0x2A,
	',': 0x2B, '/': 0x2C, 'n': 0x2D, 'm': 0x2E, '.': 0x2F, '`': 0x32,
	' ': 0x31,
}

// shiftedOf maps an unshifted key to the character it types with Shift on the
// ANSI US layout.
var shiftedOf = map[rune]rune{
	'1': '!', '2': '@', '3': '#', '4': '$', '5': '%', '6': '^', '7': '&',
	'8': '*', '9': '(', '0': ')', '-': '_', '=': '+', '[': '{', ']': '}',
	'\\': '|', ';': ':', '\'': '"', ',': '<', '.': '>', '/': '?', '`': '~',
}

// baseOfShifted is the inverse of shiftedOf.
var baseOfShifted = func() map[rune]rune {
	m := make(map[rune]rune, len(shiftedOf))
	for base, shifted := range shiftedOf {
		m[shifted] = base
	}
	return m
}()

// fnKeyCodes are kVK_F1 ... kVK_F12.
var fnKeyCodes = [12]uint16{0x7A, 0x78, 0x63, 0x76, 0x60, 0x61, 0x62, 0x64, 0x65, 0x6D, 0x67, 0x6F}

// namedKeys are the non-printable keys a chord can name, lower case.
var namedKeys = func() map[string]namedKey {
	m := map[string]namedKey{
		"return":    {0x24, "\r", false},
		"enter":     {0x24, "\r", false},
		"tab":       {0x30, "\t", false},
		"space":     {0x31, " ", false},
		"delete":    {0x33, "\x7f", false},
		"backspace": {0x33, "\x7f", false},
		"escape":    {0x35, "\x1b", false},
		"esc":       {0x35, "\x1b", false},
		"up":        {0x7E, "", true},
		"down":      {0x7D, "", true},
		"left":      {0x7B, "", true},
		"right":     {0x7C, "", true},
		"home":      {0x73, "", true},
		"end":       {0x77, "", true},
		"pageup":    {0x74, "", true},
		"pagedown":  {0x79, "", true},
	}
	for i, code := range fnKeyCodes {
		m[fmt.Sprintf("f%d", i+1)] = namedKey{code, string(rune(0xF704 + i)), true}
	}
	return m
}()

// modifierFlags is the NSEventModifierFlags for a chord's modifiers.
func modifierFlags(c chord) uint {
	var f uint
	if c.cmd {
		f |= flagCommand
	}
	if c.shift {
		f |= flagShift
	}
	if c.ctrl {
		f |= flagControl
	}
	if c.alt {
		f |= flagOption
	}
	return f
}

// chordStroke maps a parsed chord to the key event AppKit expects for it:
// the virtual key code, the modifier flags, and the characters the key types
// with those modifiers (Control turns a letter into its control character,
// Shift into the upper case or the shifted symbol).
func chordStroke(c chord) (keyStroke, error) {
	s := keyStroke{flags: modifierFlags(c)}
	if nk, ok := namedKeys[c.key]; ok {
		s.code = nk.code
		s.chars = nk.chars
		s.ignoring = nk.chars
		if nk.function {
			s.flags |= flagFunction
		}
		return s, nil
	}
	r, size := utf8.DecodeRuneInString(c.key)
	if size == 0 || size != len(c.key) || r == utf8.RuneError {
		return s, fmt.Errorf("chord %q: unknown key %q", c.key, c.key)
	}
	base := unicode.ToLower(r)
	code, ok := keyCodes[base]
	if !ok {
		return s, fmt.Errorf("chord key %q has no macOS virtual key code", c.key)
	}
	s.code = code
	s.ignoring = string(base)
	s.chars = string(base)
	switch {
	case c.ctrl && unicode.IsLetter(base):
		s.chars = string(base - 'a' + 1)
	case c.shift && unicode.IsLetter(base):
		s.chars = string(unicode.ToUpper(base))
	case c.shift && shiftedOf[base] != 0:
		s.chars = string(shiftedOf[base])
	}
	return s, nil
}

// typeStroke is the key press that types r. Upper case letters and shifted
// symbols are Shift plus the base key; characters with no key on the ANSI
// layout get key code 0 and carry the character alone, which text input
// (NSTextView, WKWebView, NSTextField) takes from the event's characters.
func typeStroke(r rune) keyStroke {
	base := r
	var flags uint
	switch {
	case unicode.IsUpper(r):
		base = unicode.ToLower(r)
		flags = flagShift
	case baseOfShifted[r] != 0:
		base = baseOfShifted[r]
		flags = flagShift
	}
	return keyStroke{code: keyCodes[base], flags: flags, chars: string(r), ignoring: string(base)}
}

// contentToWindow converts a point in web content coordinates (points, origin
// top left) to window base coordinates (origin bottom left) for a content view
// of contentHeight points.
func contentToWindow(x, y, contentHeight float64) (float64, float64) {
	return x, contentHeight - y
}

// ---------------------------------------------------------------------------
// Run loop and main-thread dispatch
// ---------------------------------------------------------------------------

var (
	runLoopOnce           sync.Once
	runLoopErr            error
	cfRunLoopGetMain      func() uintptr
	cfRunLoopPerformBlock func(loop, mode uintptr, block objc.Block)
	cfRunLoopWakeUp       func(loop uintptr)
	cfCommonModes         uintptr
)

// loadRunLoop binds the CoreFoundation run loop calls once.
func loadRunLoop() error {
	runLoopOnce.Do(func() {
		if _, err := purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit", purego.RTLD_NOW|purego.RTLD_GLOBAL); err != nil {
			runLoopErr = err
			return
		}
		cf, err := purego.Dlopen("/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation", purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			runLoopErr = err
			return
		}
		modes, err := purego.Dlsym(cf, "kCFRunLoopCommonModes")
		if err != nil {
			runLoopErr = err
			return
		}
		// modes is the address of the kCFRunLoopCommonModes variable; the
		// value of the variable is the CFStringRef to pass.
		cfCommonModes = **(**uintptr)(unsafe.Pointer(&modes))
		purego.RegisterLibFunc(&cfRunLoopGetMain, cf, "CFRunLoopGetMain")
		purego.RegisterLibFunc(&cfRunLoopPerformBlock, cf, "CFRunLoopPerformBlock")
		purego.RegisterLibFunc(&cfRunLoopWakeUp, cf, "CFRunLoopWakeUp")
	})
	return runLoopErr
}

// onMain runs fn on the main thread and returns its error. It schedules fn on
// the main run loop in the common modes and waits for it, so fn also runs
// while a modal session or menu tracking owns the main thread.
func onMain(fn func() error) error {
	if call[bool](class("NSThread"), "isMainThread") {
		return runSafe(fn)
	}
	done := make(chan error, 1)
	block := objc.NewBlock(func(objc.Block) { done <- runSafe(fn) })
	loop := cfRunLoopGetMain()
	cfRunLoopPerformBlock(loop, cfCommonModes, block)
	cfRunLoopWakeUp(loop)
	// CFRunLoopPerformBlock copies the block; drop our reference.
	block.Release()
	select {
	case err := <-done:
		return err
	case <-time.After(mainThreadTimeout):
		return fmt.Errorf("the main thread did not run the input event within %s (is the app running?)", mainThreadTimeout)
	}
}

// runSafe turns a panic inside fn into an error. A panic must not unwind
// through the Objective-C frames that called the block.
func runSafe(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("input driver panic on the main thread: %v", r)
		}
	}()
	return fn()
}

// ---------------------------------------------------------------------------
// Objective-C helpers. The probe is package main, so it cannot use the
// application package's unexported helpers; these mirror them.
// ---------------------------------------------------------------------------

var selCache sync.Map

func sel(name string) objc.SEL {
	if s, ok := selCache.Load(name); ok {
		return s.(objc.SEL)
	}
	s := objc.RegisterName(name)
	selCache.Store(name, s)
	return s
}

func class(name string) objc.ID { return objc.ID(objc.GetClass(name)) }

// send returns the object a message answers.
func send(o objc.ID, name string, args ...any) objc.ID {
	return o.Send(sel(name), args...)
}

// call sends a message and returns a typed result (BOOL, integers, floats,
// structs).
func call[T any](o objc.ID, name string, args ...any) T {
	return objc.Send[T](o, sel(name), args...)
}

func nsStr(s string) objc.ID {
	return objc.Send[objc.ID](class("NSString"), sel("stringWithUTF8String:"), s)
}

// goString converts an NSString to a Go string.
func goString(o objc.ID) string {
	if o == 0 {
		return ""
	}
	p := call[unsafe.Pointer](o, "UTF8String")
	if p == nil {
		return ""
	}
	var n int
	for *(*byte)(unsafe.Add(p, n)) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(p), n))
}

type nsPoint struct{ X, Y float64 }
type nsSize struct{ W, H float64 }
type nsRect struct {
	Origin nsPoint
	Size   nsSize
}

func nsApp() objc.ID { return objc.Send[objc.ID](class("NSApplication"), sel("sharedApplication")) }

// uptime is the timestamp for new events: seconds since boot, as NSEvent wants.
func uptime() float64 {
	return call[float64](send(class("NSProcessInfo"), "processInfo"), "systemUptime")
}

// findWindow returns the window titled windowTitle, preferring a visible one.
func findWindow(app objc.ID, windowTitle string) (objc.ID, error) {
	windows := send(app, "windows")
	n := call[uint](windows, "count")
	var hidden objc.ID
	for i := uint(0); i < n; i++ {
		w := windows.Send(sel("objectAtIndex:"), i)
		if goString(send(w, "title")) != windowTitle {
			continue
		}
		if call[bool](w, "isVisible") {
			return w, nil
		}
		if hidden == 0 {
			hidden = w
		}
	}
	if hidden != 0 {
		return hidden, nil
	}
	return 0, fmt.Errorf("no window titled %q", windowTitle)
}

// keyTarget is the window that receives key events. The key window comes
// first: a sheet (the Go To Folder sheet of an open panel) is key while its
// modal parent is the modal window. The modal and main windows are fallbacks.
func keyTarget(app objc.ID) (objc.ID, error) {
	for _, name := range []string{"keyWindow", "modalWindow", "mainWindow"} {
		if w := send(app, name); w != 0 {
			return w, nil
		}
	}
	return 0, errors.New("no key, modal or main window to send keys to")
}

// keyEvent builds a key event with +[NSEvent keyEventWithType:...], which sets
// both characters and charactersIgnoringModifiers (main-menu key equivalents
// match on the latter, so CGEvent-built events do not trigger menu items).
//
// purego limitation, Apple arm64: the method takes more integer arguments than
// fit in registers, so isARepeat (BOOL, 1 byte) and keyCode (unsigned short)
// are passed on the stack. Apple's ABI packs them at offsets 0 and 2; purego
// gives every stack argument its own 8-byte slot, so AppKit read key code 0
// for every key (a key binding saw cmd+a for cmd+k). Passing keyCode<<16 in the
// isARepeat slot puts the key code at offset 2; BOOL reads only the low byte,
// so isARepeat stays NO. On amd64 the System V ABI uses 8-byte stack slots and
// the arguments are passed normally.
func keyEvent(typ int, s keyStroke, windowNumber int) objc.ID {
	var repeatSlot uint64
	if runtime.GOARCH == "arm64" {
		repeatSlot = uint64(s.code) << 16
	}
	return objc.Send[objc.ID](class("NSEvent"),
		sel("keyEventWithType:location:modifierFlags:timestamp:windowNumber:context:characters:charactersIgnoringModifiers:isARepeat:keyCode:"),
		typ, nsPoint{}, s.flags, uptime(), windowNumber, objc.ID(0),
		nsStr(s.chars), nsStr(s.ignoring), repeatSlot, uint(s.code))
}

// postKeys posts a press and release for each stroke, to the key window.
func postKeys(strokes []keyStroke) error {
	return onMain(func() error {
		app := nsApp()
		w, err := keyTarget(app)
		if err != nil {
			return err
		}
		wn := call[int](w, "windowNumber")
		// With no key window (the application could not be activated),
		// NSApplication drops posted key-downs, except main-menu key
		// equivalents. Hand them to the window instead, which runs the same
		// performKeyEquivalent: and first-responder keyDown: path.
		direct := send(app, "keyWindow") == 0
		for _, s := range strokes {
			down := keyEvent(nsEventTypeKeyDown, s, wn)
			up := keyEvent(nsEventTypeKeyUp, s, wn)
			if down == 0 || up == 0 {
				return errors.New("AppKit refused to make a key event")
			}
			if direct {
				w.Send(sel("sendEvent:"), down)
				w.Send(sel("sendEvent:"), up)
				continue
			}
			app.Send(sel("postEvent:atStart:"), down, false)
			app.Send(sel("postEvent:atStart:"), up, false)
		}
		return nil
	})
}

// postMouse posts one mouse event at a content point of the titled window.
// pressure is 1 on press and 0 on release, as AppKit reports them.
func postMouse(windowTitle string, x, y, typ int, pressure float32) error {
	app := nsApp()
	w, err := findWindow(app, windowTitle)
	if err != nil {
		return err
	}
	bounds := call[nsRect](send(w, "contentView"), "bounds")
	wx, wy := contentToWindow(float64(x), float64(y), bounds.Size.H)
	ev := objc.Send[objc.ID](class("NSEvent"),
		sel("mouseEventWithType:location:modifierFlags:timestamp:windowNumber:context:eventNumber:clickCount:pressure:"),
		typ, nsPoint{wx, wy}, uint(0), uptime(), call[int](w, "windowNumber"), objc.ID(0),
		0, 1, pressure)
	if ev == 0 {
		return errors.New("AppKit refused to make a mouse event")
	}
	app.Send(sel("postEvent:atStart:"), ev, false)
	return nil
}

// ---------------------------------------------------------------------------
// The driver
// ---------------------------------------------------------------------------

type darwinDriver struct{}

func newDriver() (driver, error) {
	if err := loadRunLoop(); err != nil {
		return nil, fmt.Errorf("bind the run loop: %w", err)
	}
	return darwinDriver{}, nil
}

// focus activates the application and makes the titled window key, then waits
// until it really is key: a process started from a shell (as on the lab VMs)
// is not frontmost, and macOS 14+'s cooperative activate may decline.
func (darwinDriver) focus(windowTitle string) error {
	var target objc.ID
	err := onMain(func() error {
		app := nsApp()
		w, err := findWindow(app, windowTitle)
		if err != nil {
			return err
		}
		target = w
		// activateIgnoringOtherApps: is deprecated but still forces it; the
		// cooperative activate (macOS 14+) is asked too.
		if call[bool](app, "respondsToSelector:", sel("activateIgnoringOtherApps:")) {
			app.Send(sel("activateIgnoringOtherApps:"), true)
		}
		if call[bool](app, "respondsToSelector:", sel("activate")) {
			app.Send(sel("activate"))
		}
		w.Send(sel("makeKeyAndOrderFront:"), objc.ID(0))
		return nil
	})
	if err != nil {
		return err
	}
	for deadline := time.Now().Add(3 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		var key bool
		_ = onMain(func() error { key = call[bool](target, "isKeyWindow"); return nil })
		if key {
			return nil
		}
	}
	var active bool
	_ = onMain(func() error { active = call[bool](nsApp(), "isActive"); return nil })
	// Not fatal: key() delivers to the window directly when nothing is key.
	fmt.Fprintf(os.Stderr, "e2e: %q is not key after 3s (application active: %v)\n", windowTitle, active)
	return nil
}

func (darwinDriver) key(chordText string) error {
	c, err := parseChord(chordText)
	if err != nil {
		return err
	}
	s, err := chordStroke(c)
	if err != nil {
		return err
	}
	return postKeys([]keyStroke{s})
}

func (darwinDriver) typeText(text string) error {
	strokes := make([]keyStroke, 0, len(text))
	for _, r := range text {
		strokes = append(strokes, typeStroke(r))
	}
	return postKeys(strokes)
}

func (darwinDriver) click(windowTitle string, x, y int) error {
	return clickAt(windowTitle, x, y, nsEventTypeLeftMouseDown, nsEventTypeLeftMouseUp)
}

func (darwinDriver) rightClick(windowTitle string, x, y int) error {
	return clickAt(windowTitle, x, y, nsEventTypeRightMouseDown, nsEventTypeRightMouseUp)
}

// clickAt posts a press, waits clickHold on the calling goroutine (never on
// the main thread), then posts the release.
func clickAt(windowTitle string, x, y, downType, upType int) error {
	if err := onMain(func() error { return postMouse(windowTitle, x, y, downType, 1) }); err != nil {
		return err
	}
	time.Sleep(clickHold)
	return onMain(func() error { return postMouse(windowTitle, x, y, upType, 0) })
}

// cancelDialog cancels the open or save panel, as its Cancel button does. On
// macOS the panel's contents run in a separate Apple process, so posted key
// events cannot reach its fields; the panel window in this process can still
// be sent cancel:.
func (darwinDriver) cancelDialog() error {
	return onMain(func() error {
		// Wails shows the panel without a modal session, so look for a
		// visible NSSavePanel (NSOpenPanel's superclass) among the windows.
		panelClass := class("NSSavePanel")
		windows := send(nsApp(), "windows")
		for i := uint(0); i < call[uint](windows, "count"); i++ {
			w := windows.Send(sel("objectAtIndex:"), i)
			if call[bool](w, "isKindOfClass:", panelClass) && call[bool](w, "isVisible") {
				w.Send(sel("cancel:"), objc.ID(0))
				return nil
			}
		}
		return errors.New("no open or save panel is showing")
	})
}
