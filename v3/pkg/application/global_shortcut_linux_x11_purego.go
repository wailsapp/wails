//go:build linux && purego && !android && !server

package application

// CGO-free twin of global_shortcut_linux_x11.go. It implements X11 global
// shortcuts via XGrabKey on a dedicated Display connection with its own event
// loop, loading libX11 at runtime through purego instead of linking it.
//
// This file is fully self-contained (like the cgo original): it must not
// reference identifiers from linux_purego_lib.go / linux_purego_callbacks.go /
// linux_purego.go, because those are only built under !gtk3 while this file
// builds for every purego Linux configuration.

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
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

// Xlib constants used below (from X.h).
const (
	gsX11LockMask      = 1 << 1 // LockMask (CapsLock)
	gsX11Mod2Mask      = 1 << 4 // Mod2Mask (NumLock)
	gsX11KeyPress      = 2      // KeyPress event type
	gsX11GrabModeAsync = 1      // GrabModeAsync
	gsX11False         = 0      // False
	gsX11NoSymbol      = 0      // NoSymbol
)

// The lock modifiers (CapsLock, NumLock) alter the event state, so each shortcut
// must be grabbed for every combination of them or it will not fire while a
// lock is engaged.
var gsX11LockMasks = [4]uint32{0, gsX11LockMask, gsX11Mod2Mask, gsX11LockMask | gsX11Mod2Mask}

// XKeyEvent field offsets inside the 192-byte XEvent union on 64-bit Linux:
// type int32 @0, ..., state uint32 @80, keycode uint32 @84.
const (
	gsX11EventTypeOffset    = 0
	gsX11KeyEventStateOff   = 80
	gsX11KeyEventKeycodeOff = 84
)

// gsX11Funcs holds the libX11 entry points this backend needs, bound at
// runtime via purego. All Xlib pointer types (Display*, Window, KeySym,
// XErrorHandler) are uintptr.
type gsX11Funcs struct {
	openDisplay       func(name uintptr) uintptr                                                                                                 // Display *XOpenDisplay(char*)
	closeDisplay      func(display uintptr) int32                                                                                                // int XCloseDisplay(Display*)
	setErrorHandler   func(handler uintptr) uintptr                                                                                              // XErrorHandler XSetErrorHandler(XErrorHandler)
	stringToKeysym    func(name string) uintptr                                                                                                  // KeySym XStringToKeysym(char*)
	keysymToKeycode   func(display uintptr, keysym uintptr) uint8                                                                                // KeyCode XKeysymToKeycode(Display*, KeySym)
	defaultRootWindow func(display uintptr) uintptr                                                                                              // Window XDefaultRootWindow(Display*)
	grabKey           func(display uintptr, keycode int32, modifiers uint32, window uintptr, ownerEvents, pointerMode, keyboardMode int32) int32 // int XGrabKey(...)
	ungrabKey         func(display uintptr, keycode int32, modifiers uint32, window uintptr) int32                                               // int XUngrabKey(...)
	sync              func(display uintptr, discard int32) int32                                                                                 // int XSync(Display*, Bool)
	pending           func(display uintptr) int32                                                                                                // int XPending(Display*)
	nextEvent         func(display uintptr, event uintptr) int32                                                                                 // int XNextEvent(Display*, XEvent*)
	connectionNumber  func(display uintptr) int32                                                                                                // int XConnectionNumber(Display*)
}

var (
	gsX11Once sync.Once
	gsX11     *gsX11Funcs
	gsX11Err  error
)

// gsX11GrabError is the grab error flag, set by the X error handler installed
// on our dedicated Display connection. Because every Xlib call this file makes
// is serialized onto a single goroutine (the event loop), plain reads/writes
// would suffice; atomics are used out of caution since the handler runs inside
// Xlib.
var gsX11GrabError int32

// gsX11ErrorHandler is the X error handler callback, created exactly once as a
// package var (purego callback slots are a finite resource). It mirrors the
// cgo grabErrorHandler: record that an error happened and swallow it.
var gsX11ErrorHandler = purego.NewCallback(func(display, event uintptr) uintptr {
	atomic.StoreInt32(&gsX11GrabError, 1)
	return 0
})

// gsX11Load loads libX11 and binds the required symbols, once.
func gsX11Load() (*gsX11Funcs, error) {
	gsX11Once.Do(func() {
		handle, err := purego.Dlopen("libX11.so.6", purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			var err2 error
			handle, err2 = purego.Dlopen("libX11.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
			if err2 != nil {
				gsX11Err = fmt.Errorf("libX11 could not be loaded: %w", err)
				return
			}
		}
		fns := &gsX11Funcs{}
		bind := func(fptr any, name string) {
			if gsX11Err != nil {
				return
			}
			sym, err := purego.Dlsym(handle, name)
			if err != nil || sym == 0 {
				gsX11Err = fmt.Errorf("libX11 is missing symbol %s: %w", name, err)
				return
			}
			purego.RegisterFunc(fptr, sym)
		}
		bind(&fns.openDisplay, "XOpenDisplay")
		bind(&fns.closeDisplay, "XCloseDisplay")
		bind(&fns.setErrorHandler, "XSetErrorHandler")
		bind(&fns.stringToKeysym, "XStringToKeysym")
		bind(&fns.keysymToKeycode, "XKeysymToKeycode")
		bind(&fns.defaultRootWindow, "XDefaultRootWindow")
		bind(&fns.grabKey, "XGrabKey")
		bind(&fns.ungrabKey, "XUngrabKey")
		bind(&fns.sync, "XSync")
		bind(&fns.pending, "XPending")
		bind(&fns.nextEvent, "XNextEvent")
		bind(&fns.connectionNumber, "XConnectionNumber")
		if gsX11Err == nil {
			gsX11 = fns
		}
	})
	return gsX11, gsX11Err
}

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

	lib     *gsX11Funcs
	display uintptr
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

	lib, err := gsX11Load()
	if err != nil {
		g.startErr = fmt.Errorf("X11 global shortcuts unavailable: %w", err)
		return g
	}
	g.lib = lib

	g.display = lib.openDisplay(0)
	if g.display == 0 {
		g.startErr = fmt.Errorf("could not open an X11 display (global shortcuts via X require an X11 session)")
		return g
	}
	lib.setErrorHandler(gsX11ErrorHandler)

	fds := make([]int, 2)
	if err := syscall.Pipe(fds); err != nil {
		lib.closeDisplay(g.display)
		g.display = 0
		g.startErr = fmt.Errorf("could not create wake pipe: %w", err)
		return g
	}
	g.wakeR, g.wakeW = fds[0], fds[1]
	syscall.SetNonblock(g.wakeR, true)
	syscall.SetNonblock(g.wakeW, true)

	go g.eventLoop()
	return g
}

func (g *x11GlobalShortcuts) wake() {
	var b [1]byte
	_, _ = syscall.Write(g.wakeW, b[:])
}

// run executes fn on the event-loop goroutine and waits for it to complete.
func (g *x11GlobalShortcuts) run(fn func()) {
	done := make(chan struct{})
	g.opCh <- func() {
		fn()
		close(done)
	}
	g.wake()
	<-done
}

// gsX11FdSet / gsX11FdIsSet are FD_SET / FD_ISSET for syscall.FdSet, whose
// Bits array elements are 64-bit words on linux/amd64 and linux/arm64.
func gsX11FdSet(set *syscall.FdSet, fd int) {
	set.Bits[fd/64] |= 1 << (uint(fd) % 64)
}

func gsX11FdIsSet(set *syscall.FdSet, fd int) bool {
	return set.Bits[fd/64]&(1<<(uint(fd)%64)) != 0
}

// waitForEvent blocks until either an X KeyPress arrives or the wake pipe
// becomes readable. Returns:
//
//	1  -> a KeyPress; keycode and state are filled in
//	0  -> woken via the pipe (caller should service its request queue)
//	-1 -> the connection was lost
//
// This is the Go twin of the cgo gsWaitForEvent helper.
func (g *x11GlobalShortcuts) waitForEvent() (r int, keycode, state uint32) {
	xfd := int(g.lib.connectionNumber(g.display))
	for {
		for g.lib.pending(g.display) > 0 {
			// XEvent is a 192-byte union; use an 8-byte-aligned buffer.
			var ev [24]uint64
			g.lib.nextEvent(g.display, uintptr(unsafe.Pointer(&ev[0])))
			p := unsafe.Pointer(&ev[0])
			if *(*int32)(unsafe.Add(p, gsX11EventTypeOffset)) == gsX11KeyPress {
				keycode = *(*uint32)(unsafe.Add(p, gsX11KeyEventKeycodeOff))
				state = *(*uint32)(unsafe.Add(p, gsX11KeyEventStateOff))
				return 1, keycode, state
			}
		}
		var fds syscall.FdSet
		gsX11FdSet(&fds, xfd)
		gsX11FdSet(&fds, g.wakeR)
		maxfd := xfd
		if g.wakeR > maxfd {
			maxfd = g.wakeR
		}
		_, err := syscall.Select(maxfd+1, &fds, nil, nil, nil)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			return -1, 0, 0
		}
		if gsX11FdIsSet(&fds, g.wakeR) {
			var buf [64]byte
			for {
				n, rerr := syscall.Read(g.wakeR, buf[:])
				if n <= 0 || rerr != nil {
					break
				}
			}
			return 0, 0, 0
		}
	}
}

func (g *x11GlobalShortcuts) eventLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for {
		r, keycode, state := g.waitForEvent()
		switch r {
		case 0:
			g.drainOps()
		case 1:
			b := x11Binding{keycode: uint(keycode), modMask: uint(state) & x11SignificantMask}
			g.mu.RLock()
			id, ok := g.match[b]
			g.mu.RUnlock()
			if ok {
				g.manager.dispatch(id)
			}
		default:
			return
		}
	}
}

func (g *x11GlobalShortcuts) drainOps() {
	for {
		select {
		case op := <-g.opCh:
			op()
		default:
			return
		}
	}
}

func (g *x11GlobalShortcuts) modMask(accel *accelerator) uint {
	var mask uint
	for _, m := range accel.Modifiers {
		switch m {
		case CmdOrCtrlKey, ControlKey:
			mask |= x11ControlMask
		case OptionOrAltKey:
			mask |= x11Mod1Mask
		case ShiftKey:
			mask |= x11ShiftMask
		case SuperKey:
			mask |= x11Mod4Mask
		}
	}
	return mask
}

// keycodeForName resolves an X keysym name (e.g. "a", "F1", "Return") to a
// hardware keycode for this display. Returns 0 if unknown. Must be called on
// the event-loop goroutine.
func (g *x11GlobalShortcuts) keycodeForName(name string) uint {
	ks := g.lib.stringToKeysym(name)
	if ks == gsX11NoSymbol {
		return 0
	}
	return uint(g.lib.keysymToKeycode(g.display, ks))
}

// grabKey grabs keycode+modmask (and every lock-mask variant) on the root
// window. Returns 0 on success, -1 if the grab was refused (BadAccess), which
// happens when another client already holds the combination. Must be called on
// the event-loop goroutine.
func (g *x11GlobalShortcuts) grabKey(keycode, modmask uint32) int {
	root := g.lib.defaultRootWindow(g.display)
	atomic.StoreInt32(&gsX11GrabError, 0)
	for _, lm := range gsX11LockMasks {
		g.lib.grabKey(g.display, int32(keycode), modmask|lm, root, gsX11False, gsX11GrabModeAsync, gsX11GrabModeAsync)
	}
	g.lib.sync(g.display, gsX11False)
	if atomic.LoadInt32(&gsX11GrabError) != 0 {
		return -1
	}
	return 0
}

// ungrabKey releases keycode+modmask (and every lock-mask variant) on the root
// window. Must be called on the event-loop goroutine.
func (g *x11GlobalShortcuts) ungrabKey(keycode, modmask uint32) {
	root := g.lib.defaultRootWindow(g.display)
	for _, lm := range gsX11LockMasks {
		g.lib.ungrabKey(g.display, int32(keycode), modmask|lm, root)
	}
	g.lib.sync(g.display, gsX11False)
}

func (g *x11GlobalShortcuts) register(id int, accel *accelerator) error {
	if g.startErr != nil {
		return g.startErr
	}
	keysymName, ok := x11KeysymNames[accel.Key]
	if !ok {
		return fmt.Errorf("key %q is not supported as a global shortcut", accel.Key)
	}
	modMask := g.modMask(accel)

	var regErr error
	var binding x11Binding
	g.run(func() {
		keycode := g.keycodeForName(keysymName)
		if keycode == 0 {
			regErr = fmt.Errorf("key %q has no keycode on this keyboard", accel.Key)
			return
		}
		if g.grabKey(uint32(keycode), uint32(modMask)) != 0 {
			regErr = fmt.Errorf("the shortcut is already registered (possibly by another application)")
			return
		}
		binding = x11Binding{keycode: keycode, modMask: modMask}
	})
	if regErr != nil {
		return regErr
	}

	g.mu.Lock()
	g.bindings[id] = binding
	g.match[binding] = id
	g.mu.Unlock()
	return nil
}

func (g *x11GlobalShortcuts) unregister(id int) error {
	g.mu.Lock()
	binding, ok := g.bindings[id]
	if ok {
		delete(g.bindings, id)
		delete(g.match, binding)
	}
	g.mu.Unlock()
	if !ok {
		return nil
	}
	g.run(func() {
		g.ungrabKey(uint32(binding.keycode), uint32(binding.modMask))
	})
	return nil
}

func (g *x11GlobalShortcuts) unregisterAll() error {
	g.mu.Lock()
	bindings := g.bindings
	g.bindings = make(map[int]x11Binding)
	g.match = make(map[x11Binding]int)
	g.mu.Unlock()
	if len(bindings) == 0 {
		return nil
	}
	g.run(func() {
		for _, b := range bindings {
			g.ungrabKey(uint32(b.keycode), uint32(b.modMask))
		}
	})
	return nil
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
