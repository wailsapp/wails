//go:build windows
// +build windows

package windows

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"strings"
)

var ModifierMap = map[keys.Modifier]winc.Modifiers{
	keys.ShiftKey:       winc.ModShift,
	keys.ControlKey:     winc.ModControl,
	keys.OptionOrAltKey: winc.ModAlt,
	keys.CmdOrCtrlKey:   winc.ModControl,
}

func acceleratorToWincShortcut(accelerator *keys.Accelerator) winc.Shortcut {

	if accelerator == nil {
		return winc.NoShortcut
	}
	inKey := strings.ToUpper(accelerator.Key)
	key, exists := keyMap[inKey]
	if !exists {
		return winc.NoShortcut
	}
	var modifiers winc.Modifiers
	if _, exists := shiftMap[inKey]; exists {
		modifiers = winc.ModShift
	}
	for _, mod := range accelerator.Modifiers {
		modifiers |= ModifierMap[mod]
	}
	return winc.Shortcut{
		Modifiers: modifiers,
		Key:       key,
	}
}

var shiftMap = map[string]struct{}{
	"~":    {},
	")":    {},
	"!":    {},
	"@":    {},
	"#":    {},
	"$":    {},
	"%":    {},
	"^":    {},
	"&":    {},
	"*":    {},
	"(":    {},
	"_":    {},
	"PLUS": {},
	"<":    {},
	">":    {},
	"?":    {},
	":":    {},
	`"`:    {},
	"{":    {},
	"}":    {},
	"|":    {},
}

var keyMap = map[string]winc.Key{
	"0":   winc.Key0,
	"1":   winc.Key1,
	"2":   winc.Key2,
	"3":   winc.Key3,
	"4":   winc.Key4,
	"5":   winc.Key5,
	"6":   winc.Key6,
	"7":   winc.Key7,
	"8":   winc.Key8,
	"9":   winc.Key9,
	"A":   winc.KeyA,
	"B":   winc.KeyB,
	"C":   winc.KeyC,
	"D":   winc.KeyD,
	"E":   winc.KeyE,
	"F":   winc.KeyF,
	"G":   winc.KeyG,
	"H":   winc.KeyH,
	"I":   winc.KeyI,
	"J":   winc.KeyJ,
	"K":   winc.KeyK,
	"L":   winc.KeyL,
	"M":   winc.KeyM,
	"N":   winc.KeyN,
	"O":   winc.KeyO,
	"P":   winc.KeyP,
	"Q":   winc.KeyQ,
	"R":   winc.KeyR,
	"S":   winc.KeyS,
	"T":   winc.KeyT,
	"U":   winc.KeyU,
	"V":   winc.KeyV,
	"W":   winc.KeyW,
	"X":   winc.KeyX,
	"Y":   winc.KeyY,
	"Z":   winc.KeyZ,
	"F1":  winc.KeyF1,
	"F2":  winc.KeyF2,
	"F3":  winc.KeyF3,
	"F4":  winc.KeyF4,
	"F5":  winc.KeyF5,
	"F6":  winc.KeyF6,
	"F7":  winc.KeyF7,
	"F8":  winc.KeyF8,
	"F9":  winc.KeyF9,
	"F10": winc.KeyF10,
	"F11": winc.KeyF11,
	"F12": winc.KeyF12,
	"F13": winc.KeyF13,
	"F14": winc.KeyF14,
	"F15": winc.KeyF15,
	"F16": winc.KeyF16,
	"F17": winc.KeyF17,
	"F18": winc.KeyF18,
	"F19": winc.KeyF19,
	"F20": winc.KeyF20,
	"F21": winc.KeyF21,
	"F22": winc.KeyF22,
	"F23": winc.KeyF23,
	"F24": winc.KeyF24,

	"`": winc.KeyOEM3,
	",": winc.KeyOEMComma,
	".": winc.KeyOEMPeriod,
	"/": winc.KeyOEM2,
	";": winc.KeyOEM1,
	"'": winc.KeyOEM7,
	"[": winc.KeyOEM4,
	"]": winc.KeyOEM6,
	`\`: winc.KeyOEM5,

	"~":    winc.KeyOEM3, //
	")":    winc.Key0,
	"!":    winc.Key1,
	"@":    winc.Key2,
	"#":    winc.Key3,
	"$":    winc.Key4,
	"%":    winc.Key5,
	"^":    winc.Key6,
	"&":    winc.Key7,
	"*":    winc.Key8,
	"(":    winc.Key9,
	"_":    winc.KeyOEMMinus,
	"PLUS": winc.KeyOEMPlus,
	"<":    winc.KeyOEMComma,
	">":    winc.KeyOEMPeriod,
	"?":    winc.KeyOEM2,
	":":    winc.KeyOEM1,
	`"`:    winc.KeyOEM7,
	"{":    winc.KeyOEM4,
	"}":    winc.KeyOEM6,
	"|":    winc.KeyOEM5,

	"SPACE":              winc.KeySpace,
	"TAB":                winc.KeyTab,
	"CAPSLOCK":           winc.KeyCapital,
	"NUMLOCK":            winc.KeyNumlock,
	"SCROLLLOCK":         winc.KeyScroll,
	"BACKSPACE":          winc.KeyBack,
	"DELETE":             winc.KeyDelete,
	"INSERT":             winc.KeyInsert,
	"RETURN":             winc.KeyReturn,
	"ENTER":              winc.KeyReturn,
	"UP":                 winc.KeyUp,
	"DOWN":               winc.KeyDown,
	"LEFT":               winc.KeyLeft,
	"RIGHT":              winc.KeyRight,
	"HOME":               winc.KeyHome,
	"END":                winc.KeyEnd,
	"PAGEUP":             winc.KeyPrior,
	"PAGEDOWN":           winc.KeyNext,
	"ESCAPE":             winc.KeyEscape,
	"ESC":                winc.KeyEscape,
	"VOLUMEUP":           winc.KeyVolumeUp,
	"VOLUMEDOWN":         winc.KeyVolumeDown,
	"VOLUMEMUTE":         winc.KeyVolumeMute,
	"MEDIANEXTTRACK":     winc.KeyMediaNextTrack,
	"MEDIAPREVIOUSTRACK": winc.KeyMediaPrevTrack,
	"MEDIASTOP":          winc.KeyMediaStop,
	"MEDIAPLAYPAUSE":     winc.KeyMediaPlayPause,
	"PRINTSCREEN":        winc.KeyPrint,
	"NUM0":               winc.KeyNumpad0,
	"NUM1":               winc.KeyNumpad1,
	"NUM2":               winc.KeyNumpad2,
	"NUM3":               winc.KeyNumpad3,
	"NUM4":               winc.KeyNumpad4,
	"NUM5":               winc.KeyNumpad5,
	"NUM6":               winc.KeyNumpad6,
	"NUM7":               winc.KeyNumpad7,
	"NUM8":               winc.KeyNumpad8,
	"NUM9":               winc.KeyNumpad9,
	"nummult":            winc.KeyMultiply,
	"numadd":             winc.KeyAdd,
	"numsub":             winc.KeySubtract,
	"numdec":             winc.KeyDecimal,
	"numdiv":             winc.KeyDivide,
}
