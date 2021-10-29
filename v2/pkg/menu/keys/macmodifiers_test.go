package keys

import "testing"

func TestToMacModifier(t *testing.T) {

	tests := []struct {
		name        string
		accelerator *Accelerator
		want        int
	}{
		// TODO: Add test cases.
		{"nil", nil, 0},
		{"empty", &Accelerator{}, 0},
		{"key", &Accelerator{Key: "p"}, 0},
		{"cmd", CmdOrCtrl(""), NSEventModifierFlagCommand},
		{"ctrl", Control(""), NSEventModifierFlagControl},
		{"shift", Shift(""), NSEventModifierFlagShift},
		{"option", OptionOrAlt(""), NSEventModifierFlagOption},
		{"cmd+ctrl", Combo("", CmdOrCtrlKey, ControlKey), NSEventModifierFlagCommand | NSEventModifierFlagControl},
		{"cmd+ctrl+shift", Combo("", CmdOrCtrlKey, ControlKey, ShiftKey), NSEventModifierFlagCommand | NSEventModifierFlagControl | NSEventModifierFlagShift},
		{"cmd+ctrl+shift+option", Combo("", CmdOrCtrlKey, ControlKey, ShiftKey, OptionOrAltKey), NSEventModifierFlagCommand | NSEventModifierFlagControl | NSEventModifierFlagShift | NSEventModifierFlagOption},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMacModifier(tt.accelerator); got != tt.want {
				t.Errorf("ToMacModifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
