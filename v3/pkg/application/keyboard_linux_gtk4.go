//go:build linux && !android && !gtk3 && !server

package application

// gtk4NamedKeys maps the named keys of keys.go to GDK keysyms for GTK4
// accelerators. It is shared by the cgo and purego GTK4 backends.
//
// The values are X11 keysymdef.h main-keyboard keysyms. The keypad keysyms
// (XK_KP_Space 0xff80, XK_KP_Home 0xff95, XK_KP_End 0xff9c, XK_KP_Prior 0xff9a,
// XK_KP_Next 0xff9b, XK_KP_Delete 0xff9f) are different keys from the main
// keyboard's and must not be used here: an accelerator on them never matches
// the key the user presses.
var gtk4NamedKeys = map[string]uint32{
	"backspace":      0xff08, // XK_BackSpace
	"tab":            0xff09, // XK_Tab
	"return":         0xff0d, // XK_Return
	"enter":          0xff0d, // alias of return
	"escape":         0xff1b, // XK_Escape
	"left":           0xff51, // XK_Left
	"up":             0xff52, // XK_Up
	"right":          0xff53, // XK_Right
	"down":           0xff54, // XK_Down
	"page up":        0xff55, // XK_Page_Up
	"page down":      0xff56, // XK_Page_Down
	"end":            0xff57, // XK_End
	"home":           0xff50, // XK_Home
	"space":          0x0020, // XK_space
	"delete":         0xffff, // XK_Delete (the main-keyboard key, not KP_Delete)
	"forward delete": 0xffff, // XK_Delete
	"f1":             0xffbe, // XK_F1
	"f2":             0xffbf,
	"f3":             0xffc0,
	"f4":             0xffc1,
	"f5":             0xffc2,
	"f6":             0xffc3,
	"f7":             0xffc4,
	"f8":             0xffc5,
	"f9":             0xffc6,
	"f10":            0xffc7,
	"f11":            0xffc8,
	"f12":            0xffc9,
	"f13":            0xffca,
	"f14":            0xffcb,
	"f15":            0xffcc,
	"f16":            0xffcd,
	"f17":            0xffce,
	"f18":            0xffcf,
	"f19":            0xffd0,
	"f20":            0xffd1,
	"f21":            0xffd2,
	"f22":            0xffd3,
	"f23":            0xffd4,
	"f24":            0xffd5,
	"f25":            0xffd6,
	"f26":            0xffd7,
	"f27":            0xffd8,
	"f28":            0xffd9,
	"f29":            0xffda,
	"f30":            0xffdb,
	"f31":            0xffdc,
	"f32":            0xffdd,
	"f33":            0xffde,
	"f34":            0xffdf,
	"f35":            0xffe0, // XK_F35
	"numlock":        0xff7f, // XK_Num_Lock
}
