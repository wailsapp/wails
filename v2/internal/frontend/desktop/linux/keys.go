//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"

*/
import "C"
import (
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

var namedKeysToGTK = map[string]C.guint{
	"backspace": C.guint(0xff08),
	"tab":       C.guint(0xff09),
	"return":    C.guint(0xff0d),
	"enter":     C.guint(0xff0d),
	"escape":    C.guint(0xff1b),
	"left":      C.guint(0xff51),
	"right":     C.guint(0xff53),
	"up":        C.guint(0xff52),
	"down":      C.guint(0xff54),
	"space":     C.guint(0xff80),
	"delete":    C.guint(0xff9f),
	"home":      C.guint(0xff95),
	"end":       C.guint(0xff9c),
	"page up":   C.guint(0xff9a),
	"page down": C.guint(0xff9b),
	"f1":        C.guint(0xffbe),
	"f2":        C.guint(0xffbf),
	"f3":        C.guint(0xffc0),
	"f4":        C.guint(0xffc1),
	"f5":        C.guint(0xffc2),
	"f6":        C.guint(0xffc3),
	"f7":        C.guint(0xffc4),
	"f8":        C.guint(0xffc5),
	"f9":        C.guint(0xffc6),
	"f10":       C.guint(0xffc7),
	"f11":       C.guint(0xffc8),
	"f12":       C.guint(0xffc9),
	"f13":       C.guint(0xffca),
	"f14":       C.guint(0xffcb),
	"f15":       C.guint(0xffcc),
	"f16":       C.guint(0xffcd),
	"f17":       C.guint(0xffce),
	"f18":       C.guint(0xffcf),
	"f19":       C.guint(0xffd0),
	"f20":       C.guint(0xffd1),
	"f21":       C.guint(0xffd2),
	"f22":       C.guint(0xffd3),
	"f23":       C.guint(0xffd4),
	"f24":       C.guint(0xffd5),
	"f25":       C.guint(0xffd6),
	"f26":       C.guint(0xffd7),
	"f27":       C.guint(0xffd8),
	"f28":       C.guint(0xffd9),
	"f29":       C.guint(0xffda),
	"f30":       C.guint(0xffdb),
	"f31":       C.guint(0xffdc),
	"f32":       C.guint(0xffdd),
	"f33":       C.guint(0xffde),
	"f34":       C.guint(0xffdf),
	"f35":       C.guint(0xffe0),
	"numlock":   C.guint(0xff7f),
}

func acceleratorToGTK(accelerator *keys.Accelerator) (C.guint, C.GdkModifierType) {
	key := parseKey(accelerator.Key)
	mods := parseModifiers(accelerator.Modifiers)
	return key, mods
}

func parseKey(key string) C.guint {
	var result C.guint
	result, found := namedKeysToGTK[key]
	if found {
		return result
	}
	// Check for unknown namedkeys
	// Check if we only have a single character
	if len(key) != 1 {
		return C.guint(0)
	}
	keyval := rune(key[0])
	return C.gdk_unicode_to_keyval(C.guint(keyval))
}

func parseModifiers(modifiers []keys.Modifier) C.GdkModifierType {

	var result C.GdkModifierType

	for _, modifier := range modifiers {
		switch modifier {
		case keys.ShiftKey:
			result |= C.GDK_SHIFT_MASK
		case keys.ControlKey, keys.CmdOrCtrlKey:
			result |= C.GDK_CONTROL_MASK
		case keys.OptionOrAltKey:
			result |= C.GDK_MOD1_MASK
		}
	}
	return result
}
