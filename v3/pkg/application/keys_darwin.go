//go:build darwin

package application

const (
	NSEventModifierFlagShift   = 1 << 17 // Set if Shift key is pressed.
	NSEventModifierFlagControl = 1 << 18 // Set if Control key is pressed.
	NSEventModifierFlagOption  = 1 << 19 // Set if Option or Alternate key is pressed.
	NSEventModifierFlagCommand = 1 << 20 // Set if Command key is pressed.
)

// macModifierMap maps accelerator modifiers to macOS modifiers.
var macModifierMap = map[modifier]int{
	CmdOrCtrlKey:   NSEventModifierFlagCommand,
	ControlKey:     NSEventModifierFlagControl,
	OptionOrAltKey: NSEventModifierFlagOption,
	ShiftKey:       NSEventModifierFlagShift,
	SuperKey:       NSEventModifierFlagCommand,
}

// toMacModifier converts the accelerator to a macOS modifier.
func toMacModifier(modifiers []modifier) int {
	result := 0
	for _, modifier := range modifiers {
		result |= macModifierMap[modifier]
	}
	return result
}
