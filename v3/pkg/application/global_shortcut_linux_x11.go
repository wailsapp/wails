//go:build linux && cgo && !android && !server && !purego

package application

/*
#cgo pkg-config: x11
#include <X11/Xlib.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <sys/select.h>

// A grab error flag, set by the X error handler. Because every Xlib call this
// file makes happens on a single goroutine (the event loop), no locking is
// needed around it.
static int gGrabError = 0;

static int grabErrorHandler(Display *d, XErrorEvent *e) {
    gGrabError = 1;
    return 0;
}

static Display *gsOpenDisplay(void) {
    Display *d = XOpenDisplay(NULL);
    if (d != NULL) {
        XSetErrorHandler(grabErrorHandler);
    }
    return d;
}

static void gsCloseDisplay(Display *d) {
    if (d != NULL) {
        XCloseDisplay(d);
    }
}

// gsKeycodeForName resolves an X keysym name (e.g. "a", "F1", "Return") to a
// hardware keycode for this display. Returns 0 if unknown.
static unsigned int gsKeycodeForName(Display *d, const char *name) {
    KeySym ks = XStringToKeysym(name);
    if (ks == NoSymbol) {
        return 0;
    }
    return (unsigned int)XKeysymToKeycode(d, ks);
}

// The lock modifiers (CapsLock, NumLock) alter the event state, so each shortcut
// must be grabbed for every combination of them or it will not fire while a
// lock is engaged.
static const unsigned int gsLockMasks[4] = {0, LockMask, Mod2Mask, LockMask | Mod2Mask};

// gsGrabKey grabs keycode+modmask (and every lock-mask variant) on the root
// window. Returns 0 on success, -1 if the grab was refused (BadAccess), which
// happens when another client already holds the combination.
static int gsGrabKey(Display *d, unsigned int keycode, unsigned int modmask) {
    Window root = DefaultRootWindow(d);
    gGrabError = 0;
    for (int i = 0; i < 4; i++) {
        XGrabKey(d, keycode, modmask | gsLockMasks[i], root, False, GrabModeAsync, GrabModeAsync);
    }
    XSync(d, False);
    return gGrabError ? -1 : 0;
}

static void gsUngrabKey(Display *d, unsigned int keycode, unsigned int modmask) {
    Window root = DefaultRootWindow(d);
    for (int i = 0; i < 4; i++) {
        XUngrabKey(d, keycode, modmask | gsLockMasks[i], root);
    }
    XSync(d, False);
}

// gsWaitForEvent blocks until either an X KeyPress arrives or wakeFd becomes
// readable. Returns:
//   1  -> a KeyPress; *keycode and *state are filled in
//   0  -> woken via wakeFd (caller should service its request queue)
//  -1  -> the connection was lost
static int gsWaitForEvent(Display *d, int wakeFd, unsigned int *keycode, unsigned int *state) {
    int xfd = ConnectionNumber(d);
    for (;;) {
        while (XPending(d) > 0) {
            XEvent ev;
            XNextEvent(d, &ev);
            if (ev.type == KeyPress) {
                *keycode = ev.xkey.keycode;
                *state = ev.xkey.state;
                return 1;
            }
        }
        fd_set fds;
        FD_ZERO(&fds);
        FD_SET(xfd, &fds);
        FD_SET(wakeFd, &fds);
        int maxfd = xfd > wakeFd ? xfd : wakeFd;
        int r = select(maxfd + 1, &fds, NULL, NULL, NULL);
        if (r < 0) {
            if (errno == EINTR) {
                continue;
            }
            return -1;
        }
        if (FD_ISSET(wakeFd, &fds)) {
            char buf[64];
            while (read(wakeFd, buf, sizeof(buf)) > 0) {
            }
            return 0;
        }
    }
}
*/
import "C"

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

// X11 keyboard state mask bits (from X.h) that we treat as significant
// modifiers. LockMask (CapsLock) and Mod2Mask (NumLock) are deliberately
// excluded so that shortcuts fire regardless of those locks.
const (
	x11ShiftMask   = 1 << 0 // ShiftMask
	x11ControlMask = 1 << 2 // ControlMask
	x11Mod1Mask    = 1 << 3 // Mod1Mask  (Alt)
	x11Mod4Mask    = 1 << 6 // Mod4Mask  (Super)
)

const x11SignificantMask = x11ShiftMask | x11ControlMask | x11Mod1Mask | x11Mod4Mask

// x11Binding records what was grabbed so it can be matched against incoming
// events and ungrabbed later.
type x11Binding struct {
	keycode uint
	modMask uint
}

// x11GlobalShortcuts implements globalShortcutImpl on X11 using XGrabKey. All
// Xlib calls are funnelled onto a single event-loop goroutine; register and
// unregister hand work to it over opCh and wake it through a self-pipe. This
// keeps Xlib single-threaded (no XInitThreads) while still letting the loop
// block in select().
type x11GlobalShortcuts struct {
	manager *GlobalShortcutManager

	display *C.Display
	wakeR   int
	wakeW   int
	opCh    chan func()

	mu       sync.RWMutex
	bindings map[int]x11Binding // id -> grabbed keycode/mask
	match    map[x11Binding]int // keycode/mask -> id (for event lookup)
	startErr error
}

func newX11GlobalShortcuts(manager *GlobalShortcutManager) globalShortcutImpl {
	g := &x11GlobalShortcuts{
		manager:  manager,
		opCh:     make(chan func(), 16),
		bindings: make(map[int]x11Binding),
		match:    make(map[x11Binding]int),
	}

	g.display = C.gsOpenDisplay()
	if g.display == nil {
		g.startErr = fmt.Errorf("could not open an X11 display (global shortcuts via X require an X11 session)")
		return g
	}

	fds := make([]int, 2)
	if err := syscall.Pipe(fds); err != nil {
		C.gsCloseDisplay(g.display)
		g.display = nil
		g.startErr = fmt.Errorf("could not create wake pipe: %w", err)
		return g
	}
	g.wakeR, g.wakeW = fds[0], fds[1]
	syscall.SetNonblock(g.wakeR, true)
	syscall.SetNonblock(g.wakeW, true)

	go g.eventLoop()
	return g
}

// x11KeysymNames maps Wails accelerator key names (already lower-cased by
// parseAccelerator) to X keysym names accepted by XStringToKeysym. Letters and
// digits map to themselves.
var x11KeysymNames = map[string]string{
	"a": "a", "b": "b", "c": "c", "d": "d", "e": "e", "f": "f", "g": "g",
	"h": "h", "i": "i", "j": "j", "k": "k", "l": "l", "m": "m", "n": "n",
	"o": "o", "p": "p", "q": "q", "r": "r", "s": "s", "t": "t", "u": "u",
	"v": "v", "w": "w", "x": "x", "y": "y", "z": "z",
	"0": "0", "1": "1", "2": "2", "3": "3", "4": "4",
	"5": "5", "6": "6", "7": "7", "8": "8", "9": "9",
	// Punctuation
	";": "semicolon", "=": "equal", ",": "comma", "-": "minus", ".": "period",
	"/": "slash", "`": "grave", "[": "bracketleft", "\\": "backslash",
	"]": "bracketright", "'": "apostrophe", "+": "plus",
	// Named keys
	"backspace": "BackSpace",
	"tab":       "Tab",
	"return":    "Return",
	"enter":     "Return",
	"escape":    "Escape",
	"space":     "space",
	"page up":   "Prior",
	"page down": "Next",
	"end":       "End",
	"home":      "Home",
	"left":      "Left",
	"up":        "Up",
	"right":     "Right",
	"down":      "Down",
	"delete":    "Delete",
	"numlock":   "Num_Lock",
	// Function keys
	"f1": "F1", "f2": "F2", "f3": "F3", "f4": "F4", "f5": "F5", "f6": "F6",
	"f7": "F7", "f8": "F8", "f9": "F9", "f10": "F10", "f11": "F11", "f12": "F12",
	"f13": "F13", "f14": "F14", "f15": "F15", "f16": "F16", "f17": "F17",
	"f18": "F18", "f19": "F19", "f20": "F20", "f21": "F21", "f22": "F22",
	"f23": "F23", "f24": "F24",
}

// waitForEvent blocks until a key press (1, with its keycode and state), a
// wake-up (0), or a display failure (-1) occurs.
func (g *x11GlobalShortcuts) waitForEvent() (r int, keycode, state uint32) {
	var kc, st C.uint
	r = int(C.gsWaitForEvent(g.display, C.int(g.wakeR), &kc, &st))
	return r, uint32(kc), uint32(st)
}

// keycodeForName resolves an X keysym name to a hardware keycode, or 0 if unknown.
func (g *x11GlobalShortcuts) keycodeForName(name string) uint {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return uint(C.gsKeycodeForName(g.display, cname))
}

// grabKey grabs keycode+modmask on the root window. It returns 0 on success.
func (g *x11GlobalShortcuts) grabKey(keycode, modmask uint32) int {
	return int(C.gsGrabKey(g.display, C.uint(keycode), C.uint(modmask)))
}

// ungrabKey releases keycode+modmask on the root window.
func (g *x11GlobalShortcuts) ungrabKey(keycode, modmask uint32) {
	C.gsUngrabKey(g.display, C.uint(keycode), C.uint(modmask))
}
