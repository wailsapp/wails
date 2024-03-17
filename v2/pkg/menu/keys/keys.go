package keys

import (
	"fmt"
	"strings"
)

// Modifier is actually a string
type Modifier string

const (
	// CmdOrCtrlKey represents Command on Mac and Control on other platforms
	CmdOrCtrlKey Modifier = "cmdorctrl"
	// OptionOrAltKey represents Option on Mac and Alt on other platforms
	OptionOrAltKey Modifier = "optionoralt"
	// ShiftKey represents the shift key on all systems
	ShiftKey Modifier = "shift"
	// SuperKey represents Command on Mac and the Windows key on the other platforms
	// SuperKey Modifier = "super"
	// ControlKey represents the control key on all systems
	ControlKey Modifier = "ctrl"
)

var modifierMap = map[string]Modifier{
	"cmdorctrl":   CmdOrCtrlKey,
	"optionoralt": OptionOrAltKey,
	"shift":       ShiftKey,
	//"super":       SuperKey,
	"ctrl": ControlKey,
}

func parseModifier(text string) (*Modifier, error) {
	lowertext := strings.ToLower(text)
	result, valid := modifierMap[lowertext]
	if !valid {
		return nil, fmt.Errorf("'%s' is not a valid modifier", text)
	}

	return &result, nil
}

// Accelerator holds the keyboard shortcut for a menu item
type Accelerator struct {
	Key       string
	Modifiers []Modifier
}

// Key creates a standard key Accelerator
func Key(key string) *Accelerator {
	return &Accelerator{
		Key: strings.ToLower(key),
	}
}

// CmdOrCtrl creates a 'CmdOrCtrl' Accelerator
func CmdOrCtrl(key string) *Accelerator {
	return &Accelerator{
		Key:       strings.ToLower(key),
		Modifiers: []Modifier{CmdOrCtrlKey},
	}
}

// OptionOrAlt creates a 'OptionOrAlt' Accelerator
func OptionOrAlt(key string) *Accelerator {
	return &Accelerator{
		Key:       strings.ToLower(key),
		Modifiers: []Modifier{OptionOrAltKey},
	}
}

// Shift creates a 'Shift' Accelerator
func Shift(key string) *Accelerator {
	return &Accelerator{
		Key:       strings.ToLower(key),
		Modifiers: []Modifier{ShiftKey},
	}
}

// Control creates a 'Control' Accelerator
func Control(key string) *Accelerator {
	return &Accelerator{
		Key:       strings.ToLower(key),
		Modifiers: []Modifier{ControlKey},
	}
}

//
//// Super creates a 'Super' Accelerator
//func Super(key string) *Accelerator {
//	return &Accelerator{
//		Key:       strings.ToLower(key),
//		Modifiers: []Modifier{SuperKey},
//	}
//}

// Combo creates an Accelerator with multiple Modifiers
func Combo(key string, modifier1 Modifier, modifier2 Modifier, rest ...Modifier) *Accelerator {
	result := &Accelerator{
		Key:       key,
		Modifiers: []Modifier{modifier1, modifier2},
	}
	result.Modifiers = append(result.Modifiers, rest...)
	return result
}
