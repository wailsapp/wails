package menu

// Modifier is actually a string
type Modifier string

const (
	// CmdOrCtrlKey represents Command on Mac and Control on other platforms
	CmdOrCtrlKey Modifier = "CmdOrCtrl"
	// OptionOrAltKey represents Option on Mac and Alt on other platforms
	OptionOrAltKey Modifier = "OptionOrAlt"
	// ShiftKey represents the shift key on all systems
	ShiftKey Modifier = "Shift"
	// SuperKey represents Command on Mac and the Windows key on the other platforms
	SuperKey Modifier = "Super"
	// ControlKey represents the control key on all systems
	ControlKey Modifier = "Control"
)

// Accelerator holds the keyboard shortcut for a menu item
type Accelerator struct {
	Key       string
	Modifiers []Modifier
}

// Accel creates a standard key Accelerator
func Accel(key string) *Accelerator {
	return &Accelerator{
		Key: key,
	}
}

// CmdOrCtrl creates a 'CmdOrCtrl' Accelerator
func CmdOrCtrl(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{CmdOrCtrlKey},
	}
}

// OptionOrAlt creates a 'OptionOrAlt' Accelerator
func OptionOrAlt(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{OptionOrAltKey},
	}
}

// Shift creates a 'Shift' Accelerator
func Shift(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{ShiftKey},
	}
}

// Control creates a 'Control' Accelerator
func Control(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{ControlKey},
	}
}

// Super creates a 'Super' Accelerator
func Super(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{SuperKey},
	}
}

// Combo creates an Accelerator with multiple Modifiers
func Combo(key string, modifier1 Modifier, modifier2 Modifier, rest ...Modifier) *Accelerator {
	result := &Accelerator{
		Key:       key,
		Modifiers: []Modifier{modifier1, modifier2},
	}
	for _, extra := range rest {
		result.Modifiers = append(result.Modifiers, extra)
	}
	return result
}
