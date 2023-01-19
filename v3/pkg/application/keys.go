package application

import (
	"fmt"
	"strconv"
	"strings"
)

// modifier is actually a string
type modifier int

const (
	// CmdOrCtrlKey represents Command on Mac and Control on other platforms
	CmdOrCtrlKey modifier = 0 << iota
	// OptionOrAltKey represents Option on Mac and Alt on other platforms
	OptionOrAltKey modifier = 1 << iota
	// ShiftKey represents the shift key on all systems
	ShiftKey modifier = 2 << iota
	// SuperKey represents Command on Mac and the Windows key on the other platforms
	SuperKey modifier = 3 << iota
	// ControlKey represents the control key on all systems
	ControlKey modifier = 4 << iota
)

var modifierMap = map[string]modifier{
	"cmdorctrl":   CmdOrCtrlKey,
	"cmd":         CmdOrCtrlKey,
	"command":     CmdOrCtrlKey,
	"ctrl":        CmdOrCtrlKey,
	"optionoralt": OptionOrAltKey,
	"alt":         OptionOrAltKey,
	"option":      OptionOrAltKey,
	"shift":       ShiftKey,
	"super":       SuperKey,
}

// accelerator holds the keyboard shortcut for a menu item
type accelerator struct {
	Key       string
	Modifiers []modifier
}

var namedKeys = map[string]struct{}{
	"backspace": {},
	"tab":       {},
	"return":    {},
	"enter":     {},
	"escape":    {},
	"left":      {},
	"right":     {},
	"up":        {},
	"down":      {},
	"space":     {},
	"delete":    {},
	"home":      {},
	"end":       {},
	"page up":   {},
	"page down": {},
	"f1":        {},
	"f2":        {},
	"f3":        {},
	"f4":        {},
	"f5":        {},
	"f6":        {},
	"f7":        {},
	"f8":        {},
	"f9":        {},
	"f10":       {},
	"f11":       {},
	"f12":       {},
	"f13":       {},
	"f14":       {},
	"f15":       {},
	"f16":       {},
	"f17":       {},
	"f18":       {},
	"f19":       {},
	"f20":       {},
	"f21":       {},
	"f22":       {},
	"f23":       {},
	"f24":       {},
	"f25":       {},
	"f26":       {},
	"f27":       {},
	"f28":       {},
	"f29":       {},
	"f30":       {},
	"f31":       {},
	"f32":       {},
	"f33":       {},
	"f34":       {},
	"f35":       {},
	"numlock":   {},
}

func parseKey(key string) (string, bool) {

	// Lowercase!
	key = strings.ToLower(key)

	// Check special case
	if key == "plus" {
		return "+", true
	}

	// Handle named keys
	_, namedKey := namedKeys[key]
	if namedKey {
		return key, true
	}

	// Check we only have a single character
	if len(key) != 1 {
		return "", false
	}

	runeKey := rune(key[0])

	// This may be too inclusive
	if strconv.IsPrint(runeKey) {
		return key, true
	}

	return "", false

}

// parseAccelerator parses a string into an accelerator
func parseAccelerator(shortcut string) (*accelerator, error) {

	var result accelerator

	// Split the shortcut by +
	components := strings.Split(shortcut, "+")

	// If we only have one it should be a key
	// We require components
	if len(components) == 0 {
		return nil, fmt.Errorf("no components given to validateComponents")
	}

	modifiers := map[modifier]struct{}{}

	// Check components
	for index, component := range components {

		// If last component
		if index == len(components)-1 {
			processedKey, validKey := parseKey(component)
			if !validKey {
				return nil, fmt.Errorf("'%s' is not a valid key", component)
			}
			result.Key = processedKey
			continue
		}

		// Not last component - needs to be modifier
		lowercaseComponent := strings.ToLower(component)
		thisModifier, valid := modifierMap[lowercaseComponent]
		if !valid {
			return nil, fmt.Errorf("'%s' is not a valid modifier", component)
		}
		// Save this data
		modifiers[thisModifier] = struct{}{}
	}
	// return the keys as a slice
	for thisModifier := range modifiers {
		result.Modifiers = append(result.Modifiers, thisModifier)
	}
	return &result, nil
}
