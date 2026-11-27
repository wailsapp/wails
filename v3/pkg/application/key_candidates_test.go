package application

import (
	"slices"
	"testing"
)

func TestPunctuationAcceleratorCandidates(t *testing.T) {
	plus := mustAccelerator(t, "CmdOrCtrl+plus")
	minus := mustAccelerator(t, "CmdOrCtrl+-")
	equal := mustAccelerator(t, "Ctrl+=")
	for _, test := range []struct {
		name               string
		mods               []modifier
		unshifted, shifted rune
		want               []string
		matches            []string
	}{
		{"Ctrl+Shift+= is Zoom In's Ctrl++", []modifier{ControlKey, ShiftKey}, '=', '+',
			[]string{"Ctrl+Shift+=", "Ctrl++", "Ctrl+Shift+OEM_PLUS"}, []string{plus}},
		{"Ctrl+= without Shift", []modifier{ControlKey}, '=', '+',
			[]string{"Ctrl+=", "Ctrl+OEM_PLUS"}, []string{equal}},
		{"Ctrl+- is Zoom Out", []modifier{ControlKey}, '-', '_',
			[]string{"Ctrl+-", "Ctrl+OEM_PLUS"}, []string{minus}},
		{"a key with no character keeps the legacy name", []modifier{ControlKey}, 0, 0,
			[]string{"Ctrl+OEM_PLUS"}, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := punctuationAcceleratorCandidates(test.mods, test.unshifted, test.shifted, "oem_plus")
			if !slices.Equal(got, test.want) {
				t.Errorf("candidates = %q, want %q", got, test.want)
			}
			for _, m := range test.matches {
				if !slices.Contains(got, m) {
					t.Errorf("no candidate matches the accelerator %q", m)
				}
			}
		})
	}
}

func mustAccelerator(t *testing.T, s string) string {
	t.Helper()
	a, err := parseAccelerator(s)
	if err != nil {
		t.Fatal(err)
	}
	if a.Modifiers != nil {
		for i, m := range a.Modifiers {
			if m == CmdOrCtrlKey {
				a.Modifiers[i] = ControlKey
			}
		}
	}
	return a.String()
}
