//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package win32

import (
	"bytes"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"strings"
	"unsafe"
)

type Key uint16

func (k Key) String() string {
	return key2string[k]
}

// Virtual key codes
const (
	VK_LBUTTON             = 1
	VK_RBUTTON             = 2
	VK_CANCEL              = 3
	VK_MBUTTON             = 4
	VK_XBUTTON1            = 5
	VK_XBUTTON2            = 6
	VK_BACK                = 8
	VK_TAB                 = 9
	VK_CLEAR               = 12
	VK_RETURN              = 13
	VK_SHIFT               = 16
	VK_CONTROL             = 17
	VK_MENU                = 18
	VK_PAUSE               = 19
	VK_CAPITAL             = 20
	VK_KANA                = 0x15
	VK_HANGEUL             = 0x15
	VK_HANGUL              = 0x15
	VK_JUNJA               = 0x17
	VK_FINAL               = 0x18
	VK_HANJA               = 0x19
	VK_KANJI               = 0x19
	VK_ESCAPE              = 0x1B
	VK_CONVERT             = 0x1C
	VK_NONCONVERT          = 0x1D
	VK_ACCEPT              = 0x1E
	VK_MODECHANGE          = 0x1F
	VK_SPACE               = 32
	VK_PRIOR               = 33
	VK_NEXT                = 34
	VK_END                 = 35
	VK_HOME                = 36
	VK_LEFT                = 37
	VK_UP                  = 38
	VK_RIGHT               = 39
	VK_DOWN                = 40
	VK_SELECT              = 41
	VK_PRINT               = 42
	VK_EXECUTE             = 43
	VK_SNAPSHOT            = 44
	VK_INSERT              = 45
	VK_DELETE              = 46
	VK_HELP                = 47
	VK_LWIN                = 0x5B
	VK_RWIN                = 0x5C
	VK_APPS                = 0x5D
	VK_SLEEP               = 0x5F
	VK_NUMPAD0             = 0x60
	VK_NUMPAD1             = 0x61
	VK_NUMPAD2             = 0x62
	VK_NUMPAD3             = 0x63
	VK_NUMPAD4             = 0x64
	VK_NUMPAD5             = 0x65
	VK_NUMPAD6             = 0x66
	VK_NUMPAD7             = 0x67
	VK_NUMPAD8             = 0x68
	VK_NUMPAD9             = 0x69
	VK_MULTIPLY            = 0x6A
	VK_ADD                 = 0x6B
	VK_SEPARATOR           = 0x6C
	VK_SUBTRACT            = 0x6D
	VK_DECIMAL             = 0x6E
	VK_DIVIDE              = 0x6F
	VK_F1                  = 0x70
	VK_F2                  = 0x71
	VK_F3                  = 0x72
	VK_F4                  = 0x73
	VK_F5                  = 0x74
	VK_F6                  = 0x75
	VK_F7                  = 0x76
	VK_F8                  = 0x77
	VK_F9                  = 0x78
	VK_F10                 = 0x79
	VK_F11                 = 0x7A
	VK_F12                 = 0x7B
	VK_F13                 = 0x7C
	VK_F14                 = 0x7D
	VK_F15                 = 0x7E
	VK_F16                 = 0x7F
	VK_F17                 = 0x80
	VK_F18                 = 0x81
	VK_F19                 = 0x82
	VK_F20                 = 0x83
	VK_F21                 = 0x84
	VK_F22                 = 0x85
	VK_F23                 = 0x86
	VK_F24                 = 0x87
	VK_NUMLOCK             = 0x90
	VK_SCROLL              = 0x91
	VK_LSHIFT              = 0xA0
	VK_RSHIFT              = 0xA1
	VK_LCONTROL            = 0xA2
	VK_RCONTROL            = 0xA3
	VK_LMENU               = 0xA4
	VK_RMENU               = 0xA5
	VK_BROWSER_BACK        = 0xA6
	VK_BROWSER_FORWARD     = 0xA7
	VK_BROWSER_REFRESH     = 0xA8
	VK_BROWSER_STOP        = 0xA9
	VK_BROWSER_SEARCH      = 0xAA
	VK_BROWSER_FAVORITES   = 0xAB
	VK_BROWSER_HOME        = 0xAC
	VK_VOLUME_MUTE         = 0xAD
	VK_VOLUME_DOWN         = 0xAE
	VK_VOLUME_UP           = 0xAF
	VK_MEDIA_NEXT_TRACK    = 0xB0
	VK_MEDIA_PREV_TRACK    = 0xB1
	VK_MEDIA_STOP          = 0xB2
	VK_MEDIA_PLAY_PAUSE    = 0xB3
	VK_LAUNCH_MAIL         = 0xB4
	VK_LAUNCH_MEDIA_SELECT = 0xB5
	VK_LAUNCH_APP1         = 0xB6
	VK_LAUNCH_APP2         = 0xB7
	VK_OEM_1               = 0xBA
	VK_OEM_PLUS            = 0xBB
	VK_OEM_COMMA           = 0xBC
	VK_OEM_MINUS           = 0xBD
	VK_OEM_PERIOD          = 0xBE
	VK_OEM_2               = 0xBF
	VK_OEM_3               = 0xC0
	VK_OEM_4               = 0xDB
	VK_OEM_5               = 0xDC
	VK_OEM_6               = 0xDD
	VK_OEM_7               = 0xDE
	VK_OEM_8               = 0xDF
	VK_OEM_102             = 0xE2
	VK_PROCESSKEY          = 0xE5
	VK_PACKET              = 0xE7
	VK_ATTN                = 0xF6
	VK_CRSEL               = 0xF7
	VK_EXSEL               = 0xF8
	VK_EREOF               = 0xF9
	VK_PLAY                = 0xFA
	VK_ZOOM                = 0xFB
	VK_NONAME              = 0xFC
	VK_PA1                 = 0xFD
	VK_OEM_CLEAR           = 0xFE
)

const (
	KeyLButton           Key = VK_LBUTTON
	KeyRButton           Key = VK_RBUTTON
	KeyCancel            Key = VK_CANCEL
	KeyMButton           Key = VK_MBUTTON
	KeyXButton1          Key = VK_XBUTTON1
	KeyXButton2          Key = VK_XBUTTON2
	KeyBack              Key = VK_BACK
	KeyTab               Key = VK_TAB
	KeyClear             Key = VK_CLEAR
	KeyReturn            Key = VK_RETURN
	KeyShift             Key = VK_SHIFT
	KeyControl           Key = VK_CONTROL
	KeyAlt               Key = VK_MENU
	KeyMenu              Key = VK_MENU
	KeyPause             Key = VK_PAUSE
	KeyCapital           Key = VK_CAPITAL
	KeyKana              Key = VK_KANA
	KeyHangul            Key = VK_HANGUL
	KeyJunja             Key = VK_JUNJA
	KeyFinal             Key = VK_FINAL
	KeyHanja             Key = VK_HANJA
	KeyKanji             Key = VK_KANJI
	KeyEscape            Key = VK_ESCAPE
	KeyConvert           Key = VK_CONVERT
	KeyNonconvert        Key = VK_NONCONVERT
	KeyAccept            Key = VK_ACCEPT
	KeyModeChange        Key = VK_MODECHANGE
	KeySpace             Key = VK_SPACE
	KeyPrior             Key = VK_PRIOR
	KeyNext              Key = VK_NEXT
	KeyEnd               Key = VK_END
	KeyHome              Key = VK_HOME
	KeyLeft              Key = VK_LEFT
	KeyUp                Key = VK_UP
	KeyRight             Key = VK_RIGHT
	KeyDown              Key = VK_DOWN
	KeySelect            Key = VK_SELECT
	KeyPrint             Key = VK_PRINT
	KeyExecute           Key = VK_EXECUTE
	KeySnapshot          Key = VK_SNAPSHOT
	KeyInsert            Key = VK_INSERT
	KeyDelete            Key = VK_DELETE
	KeyHelp              Key = VK_HELP
	Key0                 Key = 0x30
	Key1                 Key = 0x31
	Key2                 Key = 0x32
	Key3                 Key = 0x33
	Key4                 Key = 0x34
	Key5                 Key = 0x35
	Key6                 Key = 0x36
	Key7                 Key = 0x37
	Key8                 Key = 0x38
	Key9                 Key = 0x39
	KeyA                 Key = 0x41
	KeyB                 Key = 0x42
	KeyC                 Key = 0x43
	KeyD                 Key = 0x44
	KeyE                 Key = 0x45
	KeyF                 Key = 0x46
	KeyG                 Key = 0x47
	KeyH                 Key = 0x48
	KeyI                 Key = 0x49
	KeyJ                 Key = 0x4A
	KeyK                 Key = 0x4B
	KeyL                 Key = 0x4C
	KeyM                 Key = 0x4D
	KeyN                 Key = 0x4E
	KeyO                 Key = 0x4F
	KeyP                 Key = 0x50
	KeyQ                 Key = 0x51
	KeyR                 Key = 0x52
	KeyS                 Key = 0x53
	KeyT                 Key = 0x54
	KeyU                 Key = 0x55
	KeyV                 Key = 0x56
	KeyW                 Key = 0x57
	KeyX                 Key = 0x58
	KeyY                 Key = 0x59
	KeyZ                 Key = 0x5A
	KeyLWIN              Key = VK_LWIN
	KeyRWIN              Key = VK_RWIN
	KeyApps              Key = VK_APPS
	KeySleep             Key = VK_SLEEP
	KeyNumpad0           Key = VK_NUMPAD0
	KeyNumpad1           Key = VK_NUMPAD1
	KeyNumpad2           Key = VK_NUMPAD2
	KeyNumpad3           Key = VK_NUMPAD3
	KeyNumpad4           Key = VK_NUMPAD4
	KeyNumpad5           Key = VK_NUMPAD5
	KeyNumpad6           Key = VK_NUMPAD6
	KeyNumpad7           Key = VK_NUMPAD7
	KeyNumpad8           Key = VK_NUMPAD8
	KeyNumpad9           Key = VK_NUMPAD9
	KeyMultiply          Key = VK_MULTIPLY
	KeyAdd               Key = VK_ADD
	KeySeparator         Key = VK_SEPARATOR
	KeySubtract          Key = VK_SUBTRACT
	KeyDecimal           Key = VK_DECIMAL
	KeyDivide            Key = VK_DIVIDE
	KeyF1                Key = VK_F1
	KeyF2                Key = VK_F2
	KeyF3                Key = VK_F3
	KeyF4                Key = VK_F4
	KeyF5                Key = VK_F5
	KeyF6                Key = VK_F6
	KeyF7                Key = VK_F7
	KeyF8                Key = VK_F8
	KeyF9                Key = VK_F9
	KeyF10               Key = VK_F10
	KeyF11               Key = VK_F11
	KeyF12               Key = VK_F12
	KeyF13               Key = VK_F13
	KeyF14               Key = VK_F14
	KeyF15               Key = VK_F15
	KeyF16               Key = VK_F16
	KeyF17               Key = VK_F17
	KeyF18               Key = VK_F18
	KeyF19               Key = VK_F19
	KeyF20               Key = VK_F20
	KeyF21               Key = VK_F21
	KeyF22               Key = VK_F22
	KeyF23               Key = VK_F23
	KeyF24               Key = VK_F24
	KeyNumlock           Key = VK_NUMLOCK
	KeyScroll            Key = VK_SCROLL
	KeyLShift            Key = VK_LSHIFT
	KeyRShift            Key = VK_RSHIFT
	KeyLControl          Key = VK_LCONTROL
	KeyRControl          Key = VK_RCONTROL
	KeyLAlt              Key = VK_LMENU
	KeyLMenu             Key = VK_LMENU
	KeyRAlt              Key = VK_RMENU
	KeyRMenu             Key = VK_RMENU
	KeyBrowserBack       Key = VK_BROWSER_BACK
	KeyBrowserForward    Key = VK_BROWSER_FORWARD
	KeyBrowserRefresh    Key = VK_BROWSER_REFRESH
	KeyBrowserStop       Key = VK_BROWSER_STOP
	KeyBrowserSearch     Key = VK_BROWSER_SEARCH
	KeyBrowserFavorites  Key = VK_BROWSER_FAVORITES
	KeyBrowserHome       Key = VK_BROWSER_HOME
	KeyVolumeMute        Key = VK_VOLUME_MUTE
	KeyVolumeDown        Key = VK_VOLUME_DOWN
	KeyVolumeUp          Key = VK_VOLUME_UP
	KeyMediaNextTrack    Key = VK_MEDIA_NEXT_TRACK
	KeyMediaPrevTrack    Key = VK_MEDIA_PREV_TRACK
	KeyMediaStop         Key = VK_MEDIA_STOP
	KeyMediaPlayPause    Key = VK_MEDIA_PLAY_PAUSE
	KeyLaunchMail        Key = VK_LAUNCH_MAIL
	KeyLaunchMediaSelect Key = VK_LAUNCH_MEDIA_SELECT
	KeyLaunchApp1        Key = VK_LAUNCH_APP1
	KeyLaunchApp2        Key = VK_LAUNCH_APP2
	KeyOEM1              Key = VK_OEM_1
	KeyOEMPlus           Key = VK_OEM_PLUS
	KeyOEMComma          Key = VK_OEM_COMMA
	KeyOEMMinus          Key = VK_OEM_MINUS
	KeyOEMPeriod         Key = VK_OEM_PERIOD
	KeyOEM2              Key = VK_OEM_2
	KeyOEM3              Key = VK_OEM_3
	KeyOEM4              Key = VK_OEM_4
	KeyOEM5              Key = VK_OEM_5
	KeyOEM6              Key = VK_OEM_6
	KeyOEM7              Key = VK_OEM_7
	KeyOEM8              Key = VK_OEM_8
	KeyOEM102            Key = VK_OEM_102
	KeyProcessKey        Key = VK_PROCESSKEY
	KeyPacket            Key = VK_PACKET
	KeyAttn              Key = VK_ATTN
	KeyCRSel             Key = VK_CRSEL
	KeyEXSel             Key = VK_EXSEL
	KeyErEOF             Key = VK_EREOF
	KeyPlay              Key = VK_PLAY
	KeyZoom              Key = VK_ZOOM
	KeyNoName            Key = VK_NONAME
	KeyPA1               Key = VK_PA1
	KeyOEMClear          Key = VK_OEM_CLEAR
)

var key2string = map[Key]string{
	KeyLButton:           "LButton",
	KeyRButton:           "RButton",
	KeyCancel:            "Cancel",
	KeyMButton:           "MButton",
	KeyXButton1:          "XButton1",
	KeyXButton2:          "XButton2",
	KeyBack:              "Back",
	KeyTab:               "Tab",
	KeyClear:             "Clear",
	KeyReturn:            "Return",
	KeyShift:             "Shift",
	KeyControl:           "Control",
	KeyAlt:               "Alt / Menu",
	KeyPause:             "Pause",
	KeyCapital:           "Capital",
	KeyKana:              "Kana / Hangul",
	KeyJunja:             "Junja",
	KeyFinal:             "Final",
	KeyHanja:             "Hanja / Kanji",
	KeyEscape:            "Escape",
	KeyConvert:           "Convert",
	KeyNonconvert:        "Nonconvert",
	KeyAccept:            "Accept",
	KeyModeChange:        "ModeChange",
	KeySpace:             "Space",
	KeyPrior:             "Prior",
	KeyNext:              "Next",
	KeyEnd:               "End",
	KeyHome:              "Home",
	KeyLeft:              "Left",
	KeyUp:                "Up",
	KeyRight:             "Right",
	KeyDown:              "Down",
	KeySelect:            "Select",
	KeyPrint:             "Print",
	KeyExecute:           "Execute",
	KeySnapshot:          "Snapshot",
	KeyInsert:            "Insert",
	KeyDelete:            "Delete",
	KeyHelp:              "Help",
	Key0:                 "0",
	Key1:                 "1",
	Key2:                 "2",
	Key3:                 "3",
	Key4:                 "4",
	Key5:                 "5",
	Key6:                 "6",
	Key7:                 "7",
	Key8:                 "8",
	Key9:                 "9",
	KeyA:                 "A",
	KeyB:                 "B",
	KeyC:                 "C",
	KeyD:                 "D",
	KeyE:                 "E",
	KeyF:                 "F",
	KeyG:                 "G",
	KeyH:                 "H",
	KeyI:                 "I",
	KeyJ:                 "J",
	KeyK:                 "K",
	KeyL:                 "L",
	KeyM:                 "M",
	KeyN:                 "N",
	KeyO:                 "O",
	KeyP:                 "P",
	KeyQ:                 "Q",
	KeyR:                 "R",
	KeyS:                 "S",
	KeyT:                 "T",
	KeyU:                 "U",
	KeyV:                 "V",
	KeyW:                 "W",
	KeyX:                 "X",
	KeyY:                 "Y",
	KeyZ:                 "Z",
	KeyLWIN:              "LWIN",
	KeyRWIN:              "RWIN",
	KeyApps:              "Apps",
	KeySleep:             "Sleep",
	KeyNumpad0:           "Numpad0",
	KeyNumpad1:           "Numpad1",
	KeyNumpad2:           "Numpad2",
	KeyNumpad3:           "Numpad3",
	KeyNumpad4:           "Numpad4",
	KeyNumpad5:           "Numpad5",
	KeyNumpad6:           "Numpad6",
	KeyNumpad7:           "Numpad7",
	KeyNumpad8:           "Numpad8",
	KeyNumpad9:           "Numpad9",
	KeyMultiply:          "Multiply",
	KeyAdd:               "Add",
	KeySeparator:         "Separator",
	KeySubtract:          "Subtract",
	KeyDecimal:           "Decimal",
	KeyDivide:            "Divide",
	KeyF1:                "F1",
	KeyF2:                "F2",
	KeyF3:                "F3",
	KeyF4:                "F4",
	KeyF5:                "F5",
	KeyF6:                "F6",
	KeyF7:                "F7",
	KeyF8:                "F8",
	KeyF9:                "F9",
	KeyF10:               "F10",
	KeyF11:               "F11",
	KeyF12:               "F12",
	KeyF13:               "F13",
	KeyF14:               "F14",
	KeyF15:               "F15",
	KeyF16:               "F16",
	KeyF17:               "F17",
	KeyF18:               "F18",
	KeyF19:               "F19",
	KeyF20:               "F20",
	KeyF21:               "F21",
	KeyF22:               "F22",
	KeyF23:               "F23",
	KeyF24:               "F24",
	KeyNumlock:           "Numlock",
	KeyScroll:            "Scroll",
	KeyLShift:            "LShift",
	KeyRShift:            "RShift",
	KeyLControl:          "LControl",
	KeyRControl:          "RControl",
	KeyLMenu:             "LMenu",
	KeyRMenu:             "RMenu",
	KeyBrowserBack:       "BrowserBack",
	KeyBrowserForward:    "BrowserForward",
	KeyBrowserRefresh:    "BrowserRefresh",
	KeyBrowserStop:       "BrowserStop",
	KeyBrowserSearch:     "BrowserSearch",
	KeyBrowserFavorites:  "BrowserFavorites",
	KeyBrowserHome:       "BrowserHome",
	KeyVolumeMute:        "VolumeMute",
	KeyVolumeDown:        "VolumeDown",
	KeyVolumeUp:          "VolumeUp",
	KeyMediaNextTrack:    "MediaNextTrack",
	KeyMediaPrevTrack:    "MediaPrevTrack",
	KeyMediaStop:         "MediaStop",
	KeyMediaPlayPause:    "MediaPlayPause",
	KeyLaunchMail:        "LaunchMail",
	KeyLaunchMediaSelect: "LaunchMediaSelect",
	KeyLaunchApp1:        "LaunchApp1",
	KeyLaunchApp2:        "LaunchApp2",
	KeyOEM1:              "OEM1",
	KeyOEMPlus:           "OEMPlus",
	KeyOEMComma:          "OEMComma",
	KeyOEMMinus:          "OEMMinus",
	KeyOEMPeriod:         "OEMPeriod",
	KeyOEM2:              "OEM2",
	KeyOEM3:              "OEM3",
	KeyOEM4:              "OEM4",
	KeyOEM5:              "OEM5",
	KeyOEM6:              "OEM6",
	KeyOEM7:              "OEM7",
	KeyOEM8:              "OEM8",
	KeyOEM102:            "OEM102",
	KeyProcessKey:        "ProcessKey",
	KeyPacket:            "Packet",
	KeyAttn:              "Attn",
	KeyCRSel:             "CRSel",
	KeyEXSel:             "EXSel",
	KeyErEOF:             "ErEOF",
	KeyPlay:              "Play",
	KeyZoom:              "Zoom",
	KeyNoName:            "NoName",
	KeyPA1:               "PA1",
	KeyOEMClear:          "OEMClear",
}

type Modifiers byte

func (m Modifiers) String() string {
	return modifiers2string[m]
}

var modifiers2string = map[Modifiers]string{
	ModShift:                       "Shift",
	ModControl:                     "Ctrl",
	ModControl | ModShift:          "Ctrl+Shift",
	ModAlt:                         "Alt",
	ModAlt | ModShift:              "Alt+Shift",
	ModAlt | ModControl | ModShift: "Alt+Ctrl+Shift",
}

const (
	ModShift Modifiers = 1 << iota
	ModControl
	ModAlt
)

func ModifiersDown() Modifiers {
	var m Modifiers

	if ShiftDown() {
		m |= ModShift
	}
	if ControlDown() {
		m |= ModControl
	}
	if AltDown() {
		m |= ModAlt
	}

	return m
}

type Shortcut struct {
	Modifiers Modifiers
	Key       Key
}

func (s Shortcut) String() string {
	m := s.Modifiers.String()
	if m == "" {
		return s.Key.String()
	}

	b := new(bytes.Buffer)

	b.WriteString(m)
	b.WriteRune('+')
	b.WriteString(s.Key.String())

	return b.String()
}

func GetKeyState(nVirtKey int32) int16 {
	ret, _, _ := procGetKeyState.Call(
		uintptr(nVirtKey),
	)

	return int16(ret)
}

func AltDown() bool {
	return GetKeyState(int32(KeyAlt))>>15 != 0
}

func ControlDown() bool {
	return GetKeyState(int32(KeyControl))>>15 != 0
}

func ShiftDown() bool {
	return GetKeyState(int32(KeyShift))>>15 != 0
}

var ModifierMap = map[keys.Modifier]Modifiers{
	keys.ShiftKey:       ModShift,
	keys.ControlKey:     ModControl,
	keys.OptionOrAltKey: ModAlt,
	keys.CmdOrCtrlKey:   ModControl,
}

var NoShortcut = Shortcut{}

func AcceleratorToShortcut(accelerator *keys.Accelerator) Shortcut {

	if accelerator == nil {
		return NoShortcut
	}
	inKey := strings.ToUpper(accelerator.Key)
	key, exists := KeyMap[inKey]
	if !exists {
		return NoShortcut
	}
	var modifiers Modifiers
	if _, exists := shiftMap[inKey]; exists {
		modifiers = ModShift
	}
	for _, mod := range accelerator.Modifiers {
		modifiers |= ModifierMap[mod]
	}
	return Shortcut{
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

var KeyMap = map[string]Key{
	"0":   Key0,
	"1":   Key1,
	"2":   Key2,
	"3":   Key3,
	"4":   Key4,
	"5":   Key5,
	"6":   Key6,
	"7":   Key7,
	"8":   Key8,
	"9":   Key9,
	"A":   KeyA,
	"B":   KeyB,
	"C":   KeyC,
	"D":   KeyD,
	"E":   KeyE,
	"F":   KeyF,
	"G":   KeyG,
	"H":   KeyH,
	"I":   KeyI,
	"J":   KeyJ,
	"K":   KeyK,
	"L":   KeyL,
	"M":   KeyM,
	"N":   KeyN,
	"O":   KeyO,
	"P":   KeyP,
	"Q":   KeyQ,
	"R":   KeyR,
	"S":   KeyS,
	"T":   KeyT,
	"U":   KeyU,
	"V":   KeyV,
	"W":   KeyW,
	"X":   KeyX,
	"Y":   KeyY,
	"Z":   KeyZ,
	"F1":  KeyF1,
	"F2":  KeyF2,
	"F3":  KeyF3,
	"F4":  KeyF4,
	"F5":  KeyF5,
	"F6":  KeyF6,
	"F7":  KeyF7,
	"F8":  KeyF8,
	"F9":  KeyF9,
	"F10": KeyF10,
	"F11": KeyF11,
	"F12": KeyF12,
	"F13": KeyF13,
	"F14": KeyF14,
	"F15": KeyF15,
	"F16": KeyF16,
	"F17": KeyF17,
	"F18": KeyF18,
	"F19": KeyF19,
	"F20": KeyF20,
	"F21": KeyF21,
	"F22": KeyF22,
	"F23": KeyF23,
	"F24": KeyF24,

	"`":    KeyOEM3,
	",":    KeyOEMComma,
	".":    KeyOEMPeriod,
	"/":    KeyOEM2,
	";":    KeyOEM1,
	"'":    KeyOEM7,
	"[":    KeyOEM4,
	"]":    KeyOEM6,
	`\`:    KeyOEM5,
	"~":    KeyOEM3,
	")":    Key0,
	"!":    Key1,
	"@":    Key2,
	"#":    Key3,
	"$":    Key4,
	"%":    Key5,
	"^":    Key6,
	"&":    Key7,
	"*":    Key8,
	"(":    Key9,
	"_":    KeyOEMMinus,
	"PLUS": KeyOEMPlus,
	"<":    KeyOEMComma,
	">":    KeyOEMPeriod,
	"?":    KeyOEM2,
	":":    KeyOEM1,
	`"`:    KeyOEM7,
	"{":    KeyOEM4,
	"}":    KeyOEM6,
	"|":    KeyOEM5,

	"SPACE":              KeySpace,
	"TAB":                KeyTab,
	"CAPSLOCK":           KeyCapital,
	"NUMLOCK":            KeyNumlock,
	"SCROLLLOCK":         KeyScroll,
	"BACKSPACE":          KeyBack,
	"DELETE":             KeyDelete,
	"INSERT":             KeyInsert,
	"RETURN":             KeyReturn,
	"ENTER":              KeyReturn,
	"UP":                 KeyUp,
	"DOWN":               KeyDown,
	"LEFT":               KeyLeft,
	"RIGHT":              KeyRight,
	"HOME":               KeyHome,
	"END":                KeyEnd,
	"PAGEUP":             KeyPrior,
	"PAGEDOWN":           KeyNext,
	"ESCAPE":             KeyEscape,
	"ESC":                KeyEscape,
	"VOLUMEUP":           KeyVolumeUp,
	"VOLUMEDOWN":         KeyVolumeDown,
	"VOLUMEMUTE":         KeyVolumeMute,
	"MEDIANEXTTRACK":     KeyMediaNextTrack,
	"MEDIAPREVIOUSTRACK": KeyMediaPrevTrack,
	"MEDIASTOP":          KeyMediaStop,
	"MEDIAPLAYPAUSE":     KeyMediaPlayPause,
	"PRINTSCREEN":        KeyPrint,
	"NUM0":               KeyNumpad0,
	"NUM1":               KeyNumpad1,
	"NUM2":               KeyNumpad2,
	"NUM3":               KeyNumpad3,
	"NUM4":               KeyNumpad4,
	"NUM5":               KeyNumpad5,
	"NUM6":               KeyNumpad6,
	"NUM7":               KeyNumpad7,
	"NUM8":               KeyNumpad8,
	"NUM9":               KeyNumpad9,
	"nummult":            KeyMultiply,
	"numadd":             KeyAdd,
	"numsub":             KeySubtract,
	"numdec":             KeyDecimal,
	"numdiv":             KeyDivide,
}

type Accelerator struct {
	Virtual byte
	Key     uint16
	Cmd     uint16
}

func CreateAcceleratorTable(acc []Accelerator) uintptr {
	if len(acc) == 0 {
		return 0
	}
	ret, _, _ := procCreateAcceleratorTable.Call(
		uintptr(unsafe.Pointer(&acc[0])),
		uintptr(len(acc)),
	)
	return ret
}

func TranslateAccelerator(hwnd HWND, hAccTable uintptr, lpMsg *MSG) bool {
	ret, _, _ := procTranslateAccelerator.Call(
		uintptr(hwnd),
		hAccTable,
		uintptr(unsafe.Pointer(lpMsg)),
	)
	return ret != 0
}
