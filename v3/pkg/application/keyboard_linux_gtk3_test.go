//go:build linux && !android && gtk3 && !server

package application

import "testing"

func TestGTK3AcceleratorKeyval(t *testing.T) {
	tests := []struct {
		key  string
		want uint32
	}{
		{"m", 'm'},
		{"M", 'm'}, // canonical lower case
		{"z", 'z'},
		{"+", 0x2b},
		{"-", 0x2d},
		{".", 0x2e},
		{"0", '0'},
		{"return", 0xff0d},
		{"Enter", 0xff0d},
		{"escape", 0xff1b},
		{"backspace", 0xff08},
		{"delete", 0xffff},         // the PC Delete key, as on GTK4 and Windows
		{"forward delete", 0xffff}, // GDK_KEY_Delete
		{"home", 0xff50},           // main keyboard, not KP_Home (0xff95)
		{"end", 0xff57},
		{"page up", 0xff55},
		{"page down", 0xff56},
		{"space", 0x20},
		{"f1", 0xffbe},
		{"f10", 0xffc7},
		{"f35", 0xffe0},
		{"f36", 0},
		{"f0", 0},
		{"fx", 0},
		{"", 0},
		{"ab", 0},
		{"unknownkey", 0},
		{"é", 0xe9},       // Latin-1 code point
		{"€", 0x010020ac}, // above 0xff: GDK adds 0x01000000
	}
	for _, tt := range tests {
		if got := gtk3AcceleratorKeyval(tt.key); got != tt.want {
			t.Errorf("gtk3AcceleratorKeyval(%q) = %#x, want %#x", tt.key, got, tt.want)
		}
	}
}

func TestGTK3AcceleratorModifiers(t *testing.T) {
	tests := []struct {
		name string
		mods []modifier
		want uint32
	}{
		{"none", nil, 0},
		{"shift", []modifier{ShiftKey}, 1 << 0},
		{"control", []modifier{ControlKey}, 1 << 2},
		{"cmdorctrl is control on linux", []modifier{CmdOrCtrlKey}, 1 << 2},
		{"alt is MOD1", []modifier{OptionOrAltKey}, 1 << 3},
		{"super", []modifier{SuperKey}, 1 << 26},
		{"ctrl+shift", []modifier{CmdOrCtrlKey, ShiftKey}, 1<<2 | 1<<0},
		{"all", []modifier{ControlKey, ShiftKey, OptionOrAltKey, SuperKey}, 1<<0 | 1<<2 | 1<<3 | 1<<26},
	}
	for _, tt := range tests {
		if got := gtk3AcceleratorModifiers(tt.mods); got != tt.want {
			t.Errorf("%s: gtk3AcceleratorModifiers = %#x, want %#x", tt.name, got, tt.want)
		}
	}
}

func TestDialogResponseIndex(t *testing.T) {
	yes := &Button{Label: "Yes"}
	no := &Button{Label: "No", IsCancel: true}
	noCancel := []*Button{yes, {Label: "Later"}}
	withCancel := []*Button{yes, no}
	tests := []struct {
		name     string
		response int
		buttons  []*Button
		want     int
	}{
		{"button response is its index", 0, withCancel, 0},
		{"second button", 1, withCancel, 1},
		{"no response (-1) selects nothing", -1, withCancel, -1},
		{"accept (-3) selects nothing", -3, withCancel, -1},
		{"escape or close selects the cancel button", gtk3ResponseDeleteEvent, withCancel, 1},
		{"escape or close without a cancel button does nothing", gtk3ResponseDeleteEvent, noCancel, -1},
		{"escape or close without buttons does nothing", gtk3ResponseDeleteEvent, nil, -1},
		{"nil button entries are skipped", gtk3ResponseDeleteEvent, []*Button{nil, no}, 1},
	}
	for _, tt := range tests {
		if got := dialogResponseIndex(tt.response, tt.buttons); got != tt.want {
			t.Errorf("%s: dialogResponseIndex(%d) = %d, want %d", tt.name, tt.response, got, tt.want)
		}
	}
}

// keyEvent is one key press (down) or release (up) of a hardware keycode.
type keyEvent struct {
	keycode uint
	down    bool
	repeat  bool // for a press: whether it is expected to be a repeat
}

func TestKeyRepeatFilter(t *testing.T) {
	const (
		kK     = 40 // hardware keycodes
		kJ     = 44
		kZ     = 52
		kCtrl  = 37
		kShift = 50
	)
	down := func(k uint, repeat bool) keyEvent { return keyEvent{k, true, repeat} }
	up := func(k uint) keyEvent { return keyEvent{keycode: k} }
	tests := []struct {
		name   string
		events []keyEvent
	}{
		{"modifiers then key, each new", []keyEvent{
			down(kCtrl, false), down(kShift, false), down(kK, false),
		}},
		{"fast double tap of the same key is two presses", []keyEvent{
			down(kK, false), up(kK), down(kK, false),
		}},
		{"release then press after a long pause is new", []keyEvent{
			down(kK, false), up(kK), down(kK, false), up(kK), down(kK, false),
		}},
		{"auto-repeat with no release is a repeat, however slow", []keyEvent{
			down(kK, false), down(kK, true), down(kK, true), down(kK, true),
		}},
		{"auto-repeat ends at release; the next press is new", []keyEvent{
			down(kK, false), down(kK, true), up(kK), down(kK, false),
		}},
		{"another key pressed while K is held is new, K still repeats", []keyEvent{
			down(kK, false), down(kJ, false), down(kK, true), down(kJ, true),
		}},
		{"releasing one key does not release another", []keyEvent{
			down(kK, false), down(kJ, false), up(kJ), down(kK, true),
		}},
		{"a release of a key never pressed is harmless", []keyEvent{
			up(kZ), down(kK, false),
		}},
	}
	for _, tt := range tests {
		var f keyRepeatFilter
		for i, e := range tt.events {
			if !e.down {
				f.release(e.keycode)
				continue
			}
			if got := f.press(e.keycode); got != e.repeat {
				t.Errorf("%s: event %d (press keycode %d) repeat = %v, want %v",
					tt.name, i, e.keycode, got, e.repeat)
			}
		}
	}
}

func TestKeyRepeatFilterResetOnFocusLoss(t *testing.T) {
	// K pressed, the window loses focus, K is released elsewhere: when the
	// window has focus again, the next K press is a new press, not a repeat.
	var f keyRepeatFilter
	f.press(40)
	f.reset()
	if f.press(40) {
		t.Fatal("a key released while the window was unfocused is still held")
	}
}
