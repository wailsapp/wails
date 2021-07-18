package ffenestri

type callbackData struct {
	menuID   string
	menuType menuType
}

var callbacks = map[uint16]map[uint8]callbackData{}

func addMenuCallback(key uint16, modifiers uint8, menuID string, menutype menuType) {

	if callbacks[key] == nil {
		callbacks[key] = make(map[uint8]callbackData)
	}
	callbacks[key][modifiers] = callbackData{
		menuID:   menuID,
		menuType: menutype,
	}
}

func resetCallbacks() {
	callbacks = map[uint16]map[uint8]callbackData{}
}

func getCallbackForKeyPress(key uint16, modifiers uint8) (string, menuType) {
	if callbacks[key] == nil {
		return "", ""
	}
	result := callbacks[key][modifiers]
	return result.menuID, result.menuType
}

func calculateKeycode(key string) uint16 {
	return keymap[key]
}

// TODO: Complete this list
var keymap = map[string]uint16{
	"a": 0x41,
	"b": 0x42,
	"c": 0x43,
	"d": 0x44,
	"e": 0x45,
	"f": 0x46,
	"g": 0x47,
	"h": 0x48,
	"i": 0x49,
	"j": 0x4A,
	"k": 0x4B,
	"l": 0x4C,
	"m": 0x4D,
	"n": 0x4E,
	"o": 0x4F,
	"p": 0x50,
	"q": 0x51,
	"r": 0x52,
	"s": 0x53,
	"t": 0x54,
	"u": 0x55,
	"v": 0x56,
	"w": 0x57,
	"x": 0x58,
	"y": 0x59,
	"z": 0x5A,
}
