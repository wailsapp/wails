//go:build linux

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

// The Linux driver injects input through the XTEST extension of the X
// server, loaded with purego: no cgo, no xdotool. It needs DISPLAY and the
// XTEST extension (Xvfb has both). Windows are found by title on the X tree.
// Focus is set explicitly, so no window manager is needed.

const (
	libX11  = "libX11.so.6"
	libXtst = "libXtst.so.6"

	revertToParent = 2 // RevertToParent
	currentTime    = 0 // CurrentTime
	buttonLeft     = 1
	buttonRight    = 3

	// eventDelay separates XTEST events, a realistic gap between keystrokes.
	// (The GTK3 key handler used to drop presses within 50 ms of each other;
	// fixed 2026-10-28, and the probe passes at 20 ms too.)
	eventDelay = 80 * time.Millisecond

	// menuBarHeight is the GTK menu bar inside the e2e window, in pixels. GTK
	// draws it in the window's own X window and gives the page no X window of
	// its own, so the page origin is the window origin plus this height. It
	// was measured on the container's GTK3 (the page background starts at
	// window y=27). Other themes or fonts may differ: set E2E_MENUBAR_HEIGHT.
	menuBarHeight = 27
)

// xlib holds the Xlib and XTEST entry points the driver calls.
type xlib struct {
	dpy  uintptr
	root uint64

	openDisplay       func(name *byte) uintptr
	stringToKeysym    func(s *byte) uint64
	keysymToKeycode   func(dpy uintptr, ks uint64) uint8
	keycodeToKeysym   func(dpy uintptr, kc uint8, index int32) uint64
	defaultRootWindow func(dpy uintptr) uint64
	queryTree         func(dpy uintptr, w uint64, root, parent *uint64, children *unsafe.Pointer, n *uint32) int32
	fetchName         func(dpy uintptr, w uint64, name *unsafe.Pointer) int32
	translateCoords   func(dpy uintptr, src, dst uint64, sx, sy int32, dx, dy *int32, child *uint64) int32
	raiseWindow       func(dpy uintptr, w uint64) int32
	setInputFocus     func(dpy uintptr, w uint64, revert int32, t uint64) int32
	free              func(p unsafe.Pointer) int32
	flush             func(dpy uintptr) int32
	sync              func(dpy uintptr, discard int32) int32

	fakeKey    func(dpy uintptr, kc uint32, press int32, delay uint64) int32
	fakeButton func(dpy uintptr, button uint32, press int32, delay uint64) int32
	fakeMotion func(dpy uintptr, screen, x, y int32, delay uint64) int32
}

// loadXlib dlopens libX11 and libXtst and binds every symbol the driver
// uses. A missing library or symbol is an error that names it.
func loadXlib() (*xlib, error) {
	x11, err := purego.Dlopen(libX11, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, fmt.Errorf("cannot load %s (install libX11): %w", libX11, err)
	}
	xtst, err := purego.Dlopen(libXtst, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, fmt.Errorf("cannot load %s (install libXtst): %w", libXtst, err)
	}
	x := &xlib{}
	binds := []struct {
		h    uintptr
		lib  string
		fptr any
		name string
	}{
		{x11, libX11, &x.openDisplay, "XOpenDisplay"},
		{x11, libX11, &x.stringToKeysym, "XStringToKeysym"},
		{x11, libX11, &x.keysymToKeycode, "XKeysymToKeycode"},
		{x11, libX11, &x.keycodeToKeysym, "XKeycodeToKeysym"},
		{x11, libX11, &x.defaultRootWindow, "XDefaultRootWindow"},
		{x11, libX11, &x.queryTree, "XQueryTree"},
		{x11, libX11, &x.fetchName, "XFetchName"},
		{x11, libX11, &x.translateCoords, "XTranslateCoordinates"},
		{x11, libX11, &x.raiseWindow, "XRaiseWindow"},
		{x11, libX11, &x.setInputFocus, "XSetInputFocus"},
		{x11, libX11, &x.free, "XFree"},
		{x11, libX11, &x.flush, "XFlush"},
		{x11, libX11, &x.sync, "XSync"},
		{xtst, libXtst, &x.fakeKey, "XTestFakeKeyEvent"},
		{xtst, libXtst, &x.fakeButton, "XTestFakeButtonEvent"},
		{xtst, libXtst, &x.fakeMotion, "XTestFakeMotionEvent"},
	}
	for _, b := range binds {
		addr, err := purego.Dlsym(b.h, b.name)
		if err != nil || addr == 0 {
			return nil, fmt.Errorf("%s: symbol %s not found", b.lib, b.name)
		}
		purego.RegisterFunc(b.fptr, addr)
	}
	return x, nil
}

// open connects to the display named by DISPLAY.
func (x *xlib) open() error {
	// A NULL name makes XOpenDisplay use $DISPLAY.
	x.dpy = x.openDisplay(nil)
	if x.dpy == 0 {
		return fmt.Errorf("cannot open the X display (DISPLAY=%q): run under xvfb-run", os.Getenv("DISPLAY"))
	}
	x.root = x.defaultRootWindow(x.dpy)
	return nil
}

// cstr returns a NUL-terminated copy of s.
func cstr(s string) *byte {
	b := append([]byte(s), 0)
	return &b[0]
}

// keysym resolves an X keysym name to its value, or 0 if unknown.
func (x *xlib) keysym(name string) uint64 {
	return x.stringToKeysym(cstr(name))
}

// children lists the child windows of w, bottom to top.
func (x *xlib) children(w uint64) []uint64 {
	var root, parent uint64
	var list unsafe.Pointer
	var n uint32
	if x.queryTree(x.dpy, w, &root, &parent, &list, &n) == 0 || list == nil {
		return nil
	}
	defer x.free(list)
	return append([]uint64(nil), unsafe.Slice((*uint64)(list), n)...)
}

// title returns the window's WM_NAME, which GTK sets from the window title.
func (x *xlib) title(w uint64) string {
	var name unsafe.Pointer
	if x.fetchName(x.dpy, w, &name) == 0 || name == nil {
		return ""
	}
	defer x.free(name)
	n := 0
	for *(*byte)(unsafe.Add(name, n)) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(name), n))
}

// find returns the window whose title equals title, searching the whole
// window tree depth first. The error lists the titles it saw.
func (x *xlib) find(title string) (uint64, error) {
	var seen []string
	var walk func(w uint64) uint64
	walk = func(w uint64) uint64 {
		kids := x.children(w)
		for _, c := range kids {
			t := x.title(c)
			if t == title {
				return c
			}
			if t != "" {
				seen = append(seen, t)
			}
		}
		for _, c := range kids {
			if r := walk(c); r != 0 {
				return r
			}
		}
		return 0
	}
	if w := walk(x.root); w != 0 {
		return w, nil
	}
	return 0, fmt.Errorf("no X window titled %q (titles seen: %q)", title, seen)
}

// focusWindow raises w and gives it keyboard focus.
func (x *xlib) focusWindow(w uint64) {
	x.raiseWindow(x.dpy, w)
	x.setInputFocus(x.dpy, w, revertToParent, currentTime)
	x.sync(x.dpy, 0)
	time.Sleep(eventDelay)
}

// pageOrigin returns the root-window coordinates of the page's top-left
// corner in the window top: the window origin plus the menu bar height.
func (x *xlib) pageOrigin(top uint64) (int32, int32) {
	var dx, dy int32
	var child uint64
	x.translateCoords(x.dpy, top, x.root, 0, 0, &dx, &dy, &child)
	return dx, dy + menuBarOffset()
}

// menuBarOffset is menuBarHeight, or E2E_MENUBAR_HEIGHT when it is set.
func menuBarOffset() int32 {
	if v, err := strconv.Atoi(os.Getenv("E2E_MENUBAR_HEIGHT")); err == nil && v >= 0 {
		return int32(v)
	}
	return menuBarHeight
}

// sendKey presses and releases one keycode with the modifiers held.
func (x *xlib) sendKey(kc uint32, mods []uint32) {
	for _, m := range mods {
		x.keyEvent(m, true)
	}
	x.keyEvent(kc, true)
	x.keyEvent(kc, false)
	for i := len(mods) - 1; i >= 0; i-- {
		x.keyEvent(mods[i], false)
	}
	x.sync(x.dpy, 0)
}

// keyEvent sends one key event and flushes it: Xlib buffers requests, so
// an unflushed event would reach the server with the next batch.
func (x *xlib) keyEvent(kc uint32, press bool) {
	x.fakeKey(x.dpy, kc, boolInt(press), 0)
	x.flush(x.dpy)
	time.Sleep(eventDelay)
}

func boolInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// keycodeFor resolves a keysym name to a keycode, and reports whether Shift
// is needed to type it on this keymap. It fails when the keymap has no key
// for the keysym.
func (x *xlib) keycodeFor(name string) (kc uint32, shift bool, err error) {
	ks := x.keysym(name)
	if ks == 0 {
		return 0, false, fmt.Errorf("the X server has no keysym %q", name)
	}
	code := x.keysymToKeycode(x.dpy, ks)
	if code == 0 {
		return 0, false, fmt.Errorf("the keymap has no key for keysym %q", name)
	}
	shift, ok := needsShift(x.keycodeToKeysym(x.dpy, code, 0), x.keycodeToKeysym(x.dpy, code, 1), ks)
	if !ok {
		return 0, false, fmt.Errorf("keycode %d cannot produce keysym %q", code, name)
	}
	return uint32(code), shift, nil
}

// xDriver implements driver over one X display.
type xDriver struct{ x *xlib }

// newDriver loads XTEST and connects to the display.
func newDriver() (driver, error) {
	x, err := loadXlib()
	if err != nil {
		return nil, err
	}
	if err := x.open(); err != nil {
		return nil, err
	}
	return &xDriver{x}, nil
}

func (d *xDriver) focus(title string) error {
	w, err := d.x.find(title)
	if err != nil {
		return err
	}
	d.x.focusWindow(w)
	return nil
}

func (d *xDriver) key(chordText string) error {
	c, err := parseChord(chordText)
	if err != nil {
		return err
	}
	name, err := keysymForKey(c.key)
	if err != nil {
		return err
	}
	kc, _, err := d.x.keycodeFor(name)
	if err != nil {
		return err
	}
	var mods []uint32
	for _, m := range []struct {
		on   bool
		name string
	}{{c.ctrl, "ctrl"}, {c.shift, "shift"}, {c.alt, "alt"}, {c.cmd, "cmd"}} {
		if !m.on {
			continue
		}
		mc, _, err := d.x.keycodeFor(modifierKeysym(m.name))
		if err != nil {
			return err
		}
		mods = append(mods, mc)
	}
	d.x.sendKey(kc, mods)
	return nil
}

func (d *xDriver) typeText(text string) error {
	shiftKC, _, err := d.x.keycodeFor(modifierKeysym("shift"))
	if err != nil {
		return err
	}
	for _, r := range text {
		name, ok := keysymForRune(r)
		if !ok {
			return fmt.Errorf("cannot type %q: no X keysym for it", r)
		}
		kc, shift, err := d.x.keycodeFor(name)
		if err != nil {
			return fmt.Errorf("cannot type %q: %w", r, err)
		}
		var mods []uint32
		if shift {
			mods = append(mods, shiftKC)
		}
		d.x.sendKey(kc, mods)
	}
	return nil
}

func (d *xDriver) click(title string, x, y int) error { return d.mouse(title, x, y, buttonLeft) }

func (d *xDriver) rightClick(title string, x, y int) error {
	return d.mouse(title, x, y, buttonRight)
}

// mouse moves the pointer to page point (x, y) of the window titled title,
// then presses and releases button. The window is focused first.
func (d *xDriver) mouse(title string, x, y int, button uint32) error {
	w, err := d.x.find(title)
	if err != nil {
		return err
	}
	d.x.focusWindow(w)
	ox, oy := d.x.pageOrigin(w)
	d.x.fakeMotion(d.x.dpy, 0, ox+int32(x), oy+int32(y), 0)
	d.x.flush(d.x.dpy)
	time.Sleep(eventDelay)
	d.x.fakeButton(d.x.dpy, button, 1, 0)
	d.x.flush(d.x.dpy)
	time.Sleep(eventDelay)
	d.x.fakeButton(d.x.dpy, button, 0, 0)
	d.x.sync(d.x.dpy, 0)
	time.Sleep(eventDelay)
	return nil
}
