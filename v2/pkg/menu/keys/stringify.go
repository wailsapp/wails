package keys

import (
	"strings"

	"github.com/leaanthony/slicer"
)

var modifierStringMap = map[string]map[Modifier]string{
	"windows": {
		CmdOrCtrlKey:   "Ctrl",
		ControlKey:     "Ctrl",
		OptionOrAltKey: "Alt",
		ShiftKey:       "Shift",
		// SuperKey:       "Win",
	},
	"darwin": {
		CmdOrCtrlKey:   "Cmd",
		ControlKey:     "Ctrl",
		OptionOrAltKey: "Option",
		ShiftKey:       "Shift",
		// SuperKey:       "Cmd",
	},
	"linux": {
		CmdOrCtrlKey:   "Ctrl",
		ControlKey:     "Ctrl",
		OptionOrAltKey: "Alt",
		ShiftKey:       "Shift",
		// SuperKey:       "Super",
	},
}

func Stringify(accelerator *Accelerator, platform string) string {
	result := slicer.String()
	for _, modifier := range accelerator.Modifiers {
		result.Add(modifierStringMap[platform][modifier])
	}
	result.Deduplicate()
	result.Add(strings.ToUpper(accelerator.Key))
	return result.Join("+")
}
