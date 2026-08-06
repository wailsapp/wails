package application

import (
	"errors"
	"fmt"
	"unsafe"
)

// MacScrollEdgeEffectStyle is AppKit's preferred scroll-edge treatment for a
// titlebar or split-view-item accessory controller. AppKit controls the exact
// blur radius, fade depth, and adaptation; these values select only its
// documented automatic, soft, and hard policies.
type MacScrollEdgeEffectStyle int

const (
	MacScrollEdgeEffectStyleAutomatic MacScrollEdgeEffectStyle = iota
	MacScrollEdgeEffectStyleSoft
	MacScrollEdgeEffectStyleHard
)

func validMacScrollEdgeEffectStyle(style MacScrollEdgeEffectStyle) bool {
	return style >= MacScrollEdgeEffectStyleAutomatic && style <= MacScrollEdgeEffectStyleHard
}

var (
	// ErrMacAccessoryControllerRequired means a nil native controller was
	// supplied to WrapMacAccessoryViewController.
	ErrMacAccessoryControllerRequired = errors.New("a native macOS accessory view controller is required")
	// ErrMacAccessoryControllerType means the native object is neither an
	// NSTitlebarAccessoryViewController nor an
	// NSSplitViewItemAccessoryViewController.
	ErrMacAccessoryControllerType = errors.New("native object is not a supported macOS accessory view controller")
	// ErrMacScrollEdgeEffectStyleUnavailable means the requested explicit style
	// requires macOS 26.1 or newer. Automatic remains a valid no-op fallback.
	ErrMacScrollEdgeEffectStyleUnavailable = errors.New("preferred scroll-edge effect styles require macOS 26.1 or newer")
	// ErrMacAccessoryControllerUnsupported means native AppKit accessory
	// controllers are unavailable on the current platform.
	ErrMacAccessoryControllerUnsupported = errors.New("macOS accessory view controllers are unavailable on this platform")
)

// MacAccessoryViewController is a non-owning, type-checked wrapper around an
// NSTitlebarAccessoryViewController or
// NSSplitViewItemAccessoryViewController. It exists for native integrations
// that construct or receive an AppKit accessory controller and need to use
// the shared scroll-edge style API without private selectors.
//
// The native owner (normally NSWindow or NSSplitViewItem) must keep the
// controller alive for as long as this wrapper is used.
type MacAccessoryViewController struct {
	native unsafe.Pointer
	kind   MacAccessoryViewControllerKind
}

// MacAccessoryViewControllerKind identifies the AppKit controller class held
// by a MacAccessoryViewController.
type MacAccessoryViewControllerKind uint8

const (
	MacAccessoryViewControllerKindUnknown MacAccessoryViewControllerKind = iota
	MacAccessoryViewControllerKindTitlebar
	MacAccessoryViewControllerKindSplitItem
)

// WrapMacAccessoryViewController validates and wraps an existing native
// NSTitlebarAccessoryViewController or
// NSSplitViewItemAccessoryViewController. The wrapper does not retain or
// release the native object.
func WrapMacAccessoryViewController(native unsafe.Pointer) (*MacAccessoryViewController, error) {
	if native == nil {
		return nil, ErrMacAccessoryControllerRequired
	}
	kind, err := macAccessoryControllerKind(native)
	if err != nil {
		return nil, err
	}
	if kind != MacAccessoryViewControllerKindTitlebar && kind != MacAccessoryViewControllerKindSplitItem {
		return nil, ErrMacAccessoryControllerType
	}
	return &MacAccessoryViewController{native: native, kind: kind}, nil
}

// Kind reports whether the wrapped native object is a titlebar or split-item
// accessory view controller.
func (c *MacAccessoryViewController) Kind() MacAccessoryViewControllerKind {
	if c == nil {
		return MacAccessoryViewControllerKindUnknown
	}
	return c.kind
}

// NativeController returns the wrapped AppKit controller pointer. The pointer
// is borrowed and has the same lifetime as its native owner.
func (c *MacAccessoryViewController) NativeController() unsafe.Pointer {
	if c == nil {
		return nil
	}
	return c.native
}

// SupportsPreferredScrollEdgeEffectStyle reports whether this controller can
// currently apply Apple's preferredScrollEdgeEffectStyle property. The API is
// available for both supported accessory controller classes on macOS 26.1+.
func (c *MacAccessoryViewController) SupportsPreferredScrollEdgeEffectStyle() bool {
	return c != nil && c.native != nil && macAccessoryControllerSupportsScrollEdgeEffectStyle(c.native)
}

// SetPreferredScrollEdgeEffectStyle sets AppKit's preferred effect for content
// scrolling behind this accessory. Automatic, Soft, and Hard map directly to
// NSScrollEdgeEffectStyle. The system still owns the precise fade and blur.
//
// On macOS versions before 26.1, Automatic succeeds as a no-op because it is
// the platform default; explicit Soft or Hard returns
// ErrMacScrollEdgeEffectStyleUnavailable.
func (c *MacAccessoryViewController) SetPreferredScrollEdgeEffectStyle(style MacScrollEdgeEffectStyle) error {
	if c == nil || c.native == nil {
		return ErrMacAccessoryControllerRequired
	}
	if !validMacScrollEdgeEffectStyle(style) {
		return fmt.Errorf("unknown macOS scroll-edge effect style %d", style)
	}
	return macAccessoryControllerSetScrollEdgeEffectStyle(c.native, style)
}

// PreferredScrollEdgeEffectStyle returns the controller's current AppKit
// preference. Before macOS 26.1 it returns Automatic, the system behavior.
func (c *MacAccessoryViewController) PreferredScrollEdgeEffectStyle() (MacScrollEdgeEffectStyle, error) {
	if c == nil || c.native == nil {
		return MacScrollEdgeEffectStyleAutomatic, ErrMacAccessoryControllerRequired
	}
	return macAccessoryControllerScrollEdgeEffectStyle(c.native)
}
