//go:build darwin && purego && !ios && !server

package application

// CGO-free screen enumeration.
//
// This mirrors the behaviour of the cgo screen_darwin.go by driving NSScreen
// directly through the Objective-C runtime helpers in darwin_purego_cocoa.go.
// The numeric conventions (Y-down normalisation against the primary screen and
// the point->device-pixel pre-multiplication consumed by applyDPIScaling in
// screenmanager.go) are reproduced exactly so the Screen values match the cgo
// backend.

import (
	"fmt"
	"sync"

	"github.com/ebitengine/purego"
)

// CGDisplayRotation lives in CoreGraphics and is a plain C function rather than
// an Objective-C method, so it is resolved lazily via dlsym.
var (
	cgOnce            sync.Once
	cgDisplayRotation func(displayID uint32) float64
)

func loadCoreGraphics() {
	cgOnce.Do(func() {
		handle, err := purego.Dlopen(frameworkCoreGfx, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			return
		}
		purego.RegisterLibFunc(&cgDisplayRotation, handle, "CGDisplayRotation")
	})
}

// nsScreens returns the current [NSScreen screens] array as an id.
func nsScreens() id {
	return class("NSScreen").send("screens")
}

// processScreen converts a single NSScreen into the shared *Screen value,
// applying the same coordinate normalisation and point->pixel scaling as the
// cgo processScreen + cScreenToScreen pair.
func processScreen(screen id, isPrimary bool) *Screen {
	// scaleFactor is stored as float32 in the cgo struct; round-trip through
	// float32 so the derived scaling matches bit-for-bit.
	scaleFactor := float32(get[CGFloat](screen, "backingScaleFactor"))
	sf := float64(scaleFactor)

	// NSScreen's native space is Y-up with (0,0) at the bottom-left of the
	// primary screen. Normalise to Y-down with (0,0) at the top-left so Bounds
	// matches windowGetPosition/windowSetPosition and the cross-platform
	// conventions.
	primaryScreen := nsScreens().send("firstObject")
	if primaryScreen.isNil() {
		primaryScreen = class("NSScreen").send("mainScreen")
	}
	primaryHeight := get[NSRect](primaryScreen, "frame").Size.Height

	frame := get[NSRect](screen, "frame")
	height := int(frame.Size.Height)
	width := int(frame.Size.Width)
	x := int(frame.Origin.X)
	y := int(primaryHeight - frame.Origin.Y - frame.Size.Height)

	workArea := get[NSRect](screen, "visibleFrame")
	wHeight := int(workArea.Size.Height)
	wWidth := int(workArea.Size.Width)
	wX := int(workArea.Origin.X)
	wY := int(primaryHeight - workArea.Origin.Y - workArea.Size.Height)

	// adapted from https://stackoverflow.com/a/1237490/4188138
	screenDictionary := screen.send("deviceDescription")
	screenID := screenDictionary.send("objectForKey:", nsString("NSScreenNumber"))
	displayID := get[uint32](screenID, "unsignedIntValue")
	// cgo renders the CGDirectDisplayID through printf %d (signed); keep the
	// same rendering so persisted Screen.IDs match across backends.
	idStr := fmt.Sprintf("%d", int32(displayID))

	// Physical monitor size (device pixels).
	sizeValue := screenDictionary.send("objectForKey:", nsString("NSDeviceSize"))
	physicalSize := get[NSSize](sizeValue, "sizeValue")
	pWidth := int(physicalSize.Width)
	pHeight := int(physicalSize.Height)

	loadCoreGraphics()
	var rotation float64
	if cgDisplayRotation != nil {
		rotation = cgDisplayRotation(displayID)
	}

	// localizedName is macOS 10.15+ (cgo guards with @available); a missed
	// selector is an uncatchable NSException.
	var name string
	if respondsTo(screen, "localizedName") {
		name = screen.send("localizedName").string()
	}

	// Pre-multiply the point values by backingScaleFactor so that
	// applyDPIScaling's later division by ScaleFactor lands back on the
	// original point values (see screen_darwin.go / #5556 note).
	toPhysical := func(points int) int { return int(float64(points) * sf) }

	return &Screen{
		X: toPhysical(x),
		Y: toPhysical(y),
		Size: Size{
			Width:  pWidth,
			Height: pHeight,
		},
		Bounds: Rect{
			X:      toPhysical(x),
			Y:      toPhysical(y),
			Height: toPhysical(height),
			Width:  toPhysical(width),
		},
		PhysicalBounds: Rect{
			X:      toPhysical(x),
			Y:      toPhysical(y),
			Height: toPhysical(height),
			Width:  toPhysical(width),
		},
		WorkArea: Rect{
			X:      toPhysical(wX),
			Y:      toPhysical(wY),
			Height: toPhysical(wHeight),
			Width:  toPhysical(wWidth),
		},
		PhysicalWorkArea: Rect{
			X:      toPhysical(wX),
			Y:      toPhysical(wY),
			Height: toPhysical(wHeight),
			Width:  toPhysical(wWidth),
		},
		ScaleFactor: scaleFactor,
		ID:          idStr,
		Name:        name,
		IsPrimary:   isPrimary,
		Rotation:    float32(rotation),
	}
}

// allScreens enumerates the attached screens and converts them to Go values.
// The first screen in the [NSScreen screens] snapshot is treated as primary,
// matching the cgo getAllScreens.
func allScreens() []*Screen {
	var screens []*Screen
	withAutoreleasePool(func() {
		arr := nsScreens()
		count := int(get[uint](arr, "count"))
		screens = make([]*Screen, count)
		for i := 0; i < count; i++ {
			screen := arr.send("objectAtIndex:", uint(i))
			screens[i] = processScreen(screen, i == 0)
		}
	})
	return screens
}

func (m *macosApp) processAndCacheScreens() error {
	return m.parent.Screen.LayoutScreens(allScreens())
}

func (m *macosApp) getPrimaryScreen() (*Screen, error) {
	if m.parent.Screen.GetPrimary() == nil {
		if err := m.processAndCacheScreens(); err != nil {
			return nil, err
		}
	}
	return m.parent.Screen.GetPrimary(), nil
}

func (m *macosApp) getScreens() ([]*Screen, error) {
	if len(m.parent.Screen.GetAll()) == 0 {
		if err := m.processAndCacheScreens(); err != nil {
			return nil, err
		}
	}
	return m.parent.Screen.GetAll(), nil
}

func getScreenForWindow(window *macosWebviewWindow) (*Screen, error) {
	var screen *Screen
	withAutoreleasePool(func() {
		nsScreen := id(uintptr(window.nsWindow)).send("screen")
		screen = processScreen(nsScreen, false)
	})
	return screen, nil
}

func getScreenForSystray(systray *macosSystemTray) (*Screen, error) {
	// Resolve the status item's window (statusItem.button.window) and then the
	// screen that window is displayed on. https://stackoverflow.com/a/5875019
	var screen *Screen
	withAutoreleasePool(func() {
		statusItem := id(uintptr(systray.nsStatusItem))
		window := statusItem.send("button").send("window")
		screen = processScreen(window.send("screen"), false)
	})
	return screen, nil
}
