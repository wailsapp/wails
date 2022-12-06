//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"bytes"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type Key uint16

func (k Key) String() string {
	return key2string[k]
}

const (
	KeyLButton           Key = w32.VK_LBUTTON
	KeyRButton           Key = w32.VK_RBUTTON
	KeyCancel            Key = w32.VK_CANCEL
	KeyMButton           Key = w32.VK_MBUTTON
	KeyXButton1          Key = w32.VK_XBUTTON1
	KeyXButton2          Key = w32.VK_XBUTTON2
	KeyBack              Key = w32.VK_BACK
	KeyTab               Key = w32.VK_TAB
	KeyClear             Key = w32.VK_CLEAR
	KeyReturn            Key = w32.VK_RETURN
	KeyShift             Key = w32.VK_SHIFT
	KeyControl           Key = w32.VK_CONTROL
	KeyAlt               Key = w32.VK_MENU
	KeyMenu              Key = w32.VK_MENU
	KeyPause             Key = w32.VK_PAUSE
	KeyCapital           Key = w32.VK_CAPITAL
	KeyKana              Key = w32.VK_KANA
	KeyHangul            Key = w32.VK_HANGUL
	KeyJunja             Key = w32.VK_JUNJA
	KeyFinal             Key = w32.VK_FINAL
	KeyHanja             Key = w32.VK_HANJA
	KeyKanji             Key = w32.VK_KANJI
	KeyEscape            Key = w32.VK_ESCAPE
	KeyConvert           Key = w32.VK_CONVERT
	KeyNonconvert        Key = w32.VK_NONCONVERT
	KeyAccept            Key = w32.VK_ACCEPT
	KeyModeChange        Key = w32.VK_MODECHANGE
	KeySpace             Key = w32.VK_SPACE
	KeyPrior             Key = w32.VK_PRIOR
	KeyNext              Key = w32.VK_NEXT
	KeyEnd               Key = w32.VK_END
	KeyHome              Key = w32.VK_HOME
	KeyLeft              Key = w32.VK_LEFT
	KeyUp                Key = w32.VK_UP
	KeyRight             Key = w32.VK_RIGHT
	KeyDown              Key = w32.VK_DOWN
	KeySelect            Key = w32.VK_SELECT
	KeyPrint             Key = w32.VK_PRINT
	KeyExecute           Key = w32.VK_EXECUTE
	KeySnapshot          Key = w32.VK_SNAPSHOT
	KeyInsert            Key = w32.VK_INSERT
	KeyDelete            Key = w32.VK_DELETE
	KeyHelp              Key = w32.VK_HELP
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
	KeyLWIN              Key = w32.VK_LWIN
	KeyRWIN              Key = w32.VK_RWIN
	KeyApps              Key = w32.VK_APPS
	KeySleep             Key = w32.VK_SLEEP
	KeyNumpad0           Key = w32.VK_NUMPAD0
	KeyNumpad1           Key = w32.VK_NUMPAD1
	KeyNumpad2           Key = w32.VK_NUMPAD2
	KeyNumpad3           Key = w32.VK_NUMPAD3
	KeyNumpad4           Key = w32.VK_NUMPAD4
	KeyNumpad5           Key = w32.VK_NUMPAD5
	KeyNumpad6           Key = w32.VK_NUMPAD6
	KeyNumpad7           Key = w32.VK_NUMPAD7
	KeyNumpad8           Key = w32.VK_NUMPAD8
	KeyNumpad9           Key = w32.VK_NUMPAD9
	KeyMultiply          Key = w32.VK_MULTIPLY
	KeyAdd               Key = w32.VK_ADD
	KeySeparator         Key = w32.VK_SEPARATOR
	KeySubtract          Key = w32.VK_SUBTRACT
	KeyDecimal           Key = w32.VK_DECIMAL
	KeyDivide            Key = w32.VK_DIVIDE
	KeyF1                Key = w32.VK_F1
	KeyF2                Key = w32.VK_F2
	KeyF3                Key = w32.VK_F3
	KeyF4                Key = w32.VK_F4
	KeyF5                Key = w32.VK_F5
	KeyF6                Key = w32.VK_F6
	KeyF7                Key = w32.VK_F7
	KeyF8                Key = w32.VK_F8
	KeyF9                Key = w32.VK_F9
	KeyF10               Key = w32.VK_F10
	KeyF11               Key = w32.VK_F11
	KeyF12               Key = w32.VK_F12
	KeyF13               Key = w32.VK_F13
	KeyF14               Key = w32.VK_F14
	KeyF15               Key = w32.VK_F15
	KeyF16               Key = w32.VK_F16
	KeyF17               Key = w32.VK_F17
	KeyF18               Key = w32.VK_F18
	KeyF19               Key = w32.VK_F19
	KeyF20               Key = w32.VK_F20
	KeyF21               Key = w32.VK_F21
	KeyF22               Key = w32.VK_F22
	KeyF23               Key = w32.VK_F23
	KeyF24               Key = w32.VK_F24
	KeyNumlock           Key = w32.VK_NUMLOCK
	KeyScroll            Key = w32.VK_SCROLL
	KeyLShift            Key = w32.VK_LSHIFT
	KeyRShift            Key = w32.VK_RSHIFT
	KeyLControl          Key = w32.VK_LCONTROL
	KeyRControl          Key = w32.VK_RCONTROL
	KeyLAlt              Key = w32.VK_LMENU
	KeyLMenu             Key = w32.VK_LMENU
	KeyRAlt              Key = w32.VK_RMENU
	KeyRMenu             Key = w32.VK_RMENU
	KeyBrowserBack       Key = w32.VK_BROWSER_BACK
	KeyBrowserForward    Key = w32.VK_BROWSER_FORWARD
	KeyBrowserRefresh    Key = w32.VK_BROWSER_REFRESH
	KeyBrowserStop       Key = w32.VK_BROWSER_STOP
	KeyBrowserSearch     Key = w32.VK_BROWSER_SEARCH
	KeyBrowserFavorites  Key = w32.VK_BROWSER_FAVORITES
	KeyBrowserHome       Key = w32.VK_BROWSER_HOME
	KeyVolumeMute        Key = w32.VK_VOLUME_MUTE
	KeyVolumeDown        Key = w32.VK_VOLUME_DOWN
	KeyVolumeUp          Key = w32.VK_VOLUME_UP
	KeyMediaNextTrack    Key = w32.VK_MEDIA_NEXT_TRACK
	KeyMediaPrevTrack    Key = w32.VK_MEDIA_PREV_TRACK
	KeyMediaStop         Key = w32.VK_MEDIA_STOP
	KeyMediaPlayPause    Key = w32.VK_MEDIA_PLAY_PAUSE
	KeyLaunchMail        Key = w32.VK_LAUNCH_MAIL
	KeyLaunchMediaSelect Key = w32.VK_LAUNCH_MEDIA_SELECT
	KeyLaunchApp1        Key = w32.VK_LAUNCH_APP1
	KeyLaunchApp2        Key = w32.VK_LAUNCH_APP2
	KeyOEM1              Key = w32.VK_OEM_1
	KeyOEMPlus           Key = w32.VK_OEM_PLUS
	KeyOEMComma          Key = w32.VK_OEM_COMMA
	KeyOEMMinus          Key = w32.VK_OEM_MINUS
	KeyOEMPeriod         Key = w32.VK_OEM_PERIOD
	KeyOEM2              Key = w32.VK_OEM_2
	KeyOEM3              Key = w32.VK_OEM_3
	KeyOEM4              Key = w32.VK_OEM_4
	KeyOEM5              Key = w32.VK_OEM_5
	KeyOEM6              Key = w32.VK_OEM_6
	KeyOEM7              Key = w32.VK_OEM_7
	KeyOEM8              Key = w32.VK_OEM_8
	KeyOEM102            Key = w32.VK_OEM_102
	KeyProcessKey        Key = w32.VK_PROCESSKEY
	KeyPacket            Key = w32.VK_PACKET
	KeyAttn              Key = w32.VK_ATTN
	KeyCRSel             Key = w32.VK_CRSEL
	KeyEXSel             Key = w32.VK_EXSEL
	KeyErEOF             Key = w32.VK_EREOF
	KeyPlay              Key = w32.VK_PLAY
	KeyZoom              Key = w32.VK_ZOOM
	KeyNoName            Key = w32.VK_NONAME
	KeyPA1               Key = w32.VK_PA1
	KeyOEMClear          Key = w32.VK_OEM_CLEAR
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

func AltDown() bool {
	return w32.GetKeyState(int32(KeyAlt))>>15 != 0
}

func ControlDown() bool {
	return w32.GetKeyState(int32(KeyControl))>>15 != 0
}

func ShiftDown() bool {
	return w32.GetKeyState(int32(KeyShift))>>15 != 0
}
