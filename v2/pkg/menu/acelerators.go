package menu

// Modifier is actually a string
type Modifier string

const (
	// CmdOrCtrl represents Command on Mac and Control on other platforms
	CmdOrCtrl Modifier = "CmdOrCtrl"
	// OptionOrAlt represents Option on Mac and Alt on other platforms
	OptionOrAlt Modifier = "OptionOrAlt"
	// Shift represents the shift key on all systems
	Shift Modifier = "Shift"
	// Super represents Command on Mac and the Windows key on the other platforms
	Super Modifier = "Super"
	// Control represents the control key on all systems
	Control Modifier = "Control"
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

// CmdOrCtrlAccel creates a 'CmdOrCtrl' Accelerator
func CmdOrCtrlAccel(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{CmdOrCtrl},
	}
}

// OptionOrAltAccel creates a 'OptionOrAlt' Accelerator
func OptionOrAltAccel(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{OptionOrAlt},
	}
}

// ShiftAccel creates a 'Shift' Accelerator
func ShiftAccel(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{Shift},
	}
}

// ControlAccel creates a 'Control' Accelerator
func ControlAccel(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{Control},
	}
}

// SuperAccel creates a 'Super' Accelerator
func SuperAccel(key string) *Accelerator {
	return &Accelerator{
		Key:       key,
		Modifiers: []Modifier{Super},
	}
}

// ComboAccel creates an Accelerator with multiple Modifiers
func ComboAccel(key string, modifier1 Modifier, modifier2 Modifier, rest ...Modifier) *Accelerator {
	result := &Accelerator{
		Key:       key,
		Modifiers: []Modifier{modifier1, modifier2},
	}
	for _, extra := range rest {
		result.Modifiers = append(result.Modifiers, extra)
	}
	return result
}
