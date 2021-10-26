package keys

import (
	"strconv"
	"testing"
)

func TestStringify(t *testing.T) {

	const Windows = "windows"
	const Mac = "darwin"
	const Linux = "linux"
	tests := []struct {
		arg      *Accelerator
		want     string
		platform string
	}{
		// Single Keys
		{Key("a"), "A", Windows},
		{Key(""), "", Windows},
		{Key("?"), "?", Windows},
		{Key("a"), "A", Mac},
		{Key(""), "", Mac},
		{Key("?"), "?", Mac},
		{Key("a"), "A", Linux},
		{Key(""), "", Linux},
		{Key("?"), "?", Linux},

		// Single modifier
		{Control("a"), "Ctrl+A", Windows},
		{Control("a"), "Ctrl+A", Mac},
		{Control("a"), "Ctrl+A", Linux},
		{CmdOrCtrl("a"), "Ctrl+A", Windows},
		{CmdOrCtrl("a"), "Cmd+A", Mac},
		{CmdOrCtrl("a"), "Ctrl+A", Linux},
		{Shift("a"), "Shift+A", Windows},
		{Shift("a"), "Shift+A", Mac},
		{Shift("a"), "Shift+A", Linux},
		{OptionOrAlt("a"), "Alt+A", Windows},
		{OptionOrAlt("a"), "Option+A", Mac},
		{OptionOrAlt("a"), "Alt+A", Linux},
		//{Super("a"), "Win+A", Windows},
		//{Super("a"), "Cmd+A", Mac},
		//{Super("a"), "Super+A", Linux},

		// Dual Combo non duplicate
		{Combo("a", ControlKey, OptionOrAltKey), "Ctrl+Alt+A", Windows},
		{Combo("a", ControlKey, OptionOrAltKey), "Ctrl+Option+A", Mac},
		{Combo("a", ControlKey, OptionOrAltKey), "Ctrl+Alt+A", Linux},
		{Combo("a", CmdOrCtrlKey, OptionOrAltKey), "Ctrl+Alt+A", Windows},
		{Combo("a", CmdOrCtrlKey, OptionOrAltKey), "Cmd+Option+A", Mac},
		{Combo("a", CmdOrCtrlKey, OptionOrAltKey), "Ctrl+Alt+A", Linux},
		{Combo("a", ShiftKey, OptionOrAltKey), "Shift+Alt+A", Windows},
		{Combo("a", ShiftKey, OptionOrAltKey), "Shift+Option+A", Mac},
		{Combo("a", ShiftKey, OptionOrAltKey), "Shift+Alt+A", Linux},
		//{Combo("a", SuperKey, OptionOrAltKey), "Win+Alt+A", Windows},
		//{Combo("a", SuperKey, OptionOrAltKey), "Cmd+Option+A", Mac},
		//{Combo("a", SuperKey, OptionOrAltKey), "Super+Alt+A", Linux},

		// Combo duplicate
		{Combo("a", OptionOrAltKey, OptionOrAltKey), "Alt+A", Windows},
		{Combo("a", OptionOrAltKey, OptionOrAltKey), "Option+A", Mac},
		{Combo("a", OptionOrAltKey, OptionOrAltKey), "Alt+A", Linux},
		//{Combo("a", OptionOrAltKey, SuperKey, OptionOrAltKey), "Alt+Win+A", Windows},
		//{Combo("a", OptionOrAltKey, SuperKey, OptionOrAltKey), "Option+Cmd+A", Mac},
		//{Combo("a", OptionOrAltKey, SuperKey, OptionOrAltKey), "Alt+Super+A", Linux},
	}
	for index, tt := range tests {
		t.Run(strconv.Itoa(index), func(t *testing.T) {
			if got := Stringify(tt.arg, tt.platform); got != tt.want {
				t.Errorf("Stringify() = %v, want %v", got, tt.want)
			}
		})
	}
}
