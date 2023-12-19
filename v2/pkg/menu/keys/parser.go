package keys

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/leaanthony/slicer"
)

var namedKeys = slicer.String([]string{"backspace", "tab", "return", "enter", "escape", "left", "right", "up", "down", "space", "delete", "home", "end", "page up", "page down", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17", "f18", "f19", "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27", "f28", "f29", "f30", "f31", "f32", "f33", "f34", "f35", "numlock"})

func parseKey(key string) (string, bool) {
	// Lowercase!
	key = strings.ToLower(key)

	// Check special case
	if key == "plus" {
		return "+", true
	}

	// Handle named keys
	if namedKeys.Contains(key) {
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

func Parse(shortcut string) (*Accelerator, error) {
	var result Accelerator

	// Split the shortcut by +
	components := strings.Split(shortcut, "+")

	// If we only have one it should be a key
	// We require components
	if len(components) == 0 {
		return nil, fmt.Errorf("no components given to validateComponents")
	}

	// Keep track of modifiers we have processed
	var modifiersProcessed slicer.StringSlicer

	// Check components
	for index, component := range components {

		// If last component
		if index == len(components)-1 {
			processedkey, validKey := parseKey(component)
			if !validKey {
				return nil, fmt.Errorf("'%s' is not a valid key", component)
			}
			result.Key = processedkey
			continue
		}

		// Not last component - needs to be modifier
		lowercaseComponent := strings.ToLower(component)
		thisModifier, valid := modifierMap[lowercaseComponent]
		if !valid {
			return nil, fmt.Errorf("'%s' is not a valid modifier", component)
		}
		// Needs to be unique
		if modifiersProcessed.Contains(lowercaseComponent) {
			return nil, fmt.Errorf("Modifier '%s' is defined twice for shortcut: %s", component, shortcut)
		}

		// Save this data
		result.Modifiers = append(result.Modifiers, thisModifier)
		modifiersProcessed.Add(lowercaseComponent)
	}

	return &result, nil
}
