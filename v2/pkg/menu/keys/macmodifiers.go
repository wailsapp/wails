package keys

const (
	NSEventModifierFlagShift   = 1 << 17 // Set if Shift key is pressed.
	NSEventModifierFlagControl = 1 << 18 // Set if Control key is pressed.
	NSEventModifierFlagOption  = 1 << 19 // Set if Option or Alternate key is pressed.
	NSEventModifierFlagCommand = 1 << 20 // Set if Command key is pressed.
)

var macModifierMap = map[Modifier]int{
	CmdOrCtrlKey:   NSEventModifierFlagCommand,
	ControlKey:     NSEventModifierFlagControl,
	OptionOrAltKey: NSEventModifierFlagOption,
	ShiftKey:       NSEventModifierFlagShift,
}

func ToMacModifier(accelerator *Accelerator) int {
	if accelerator == nil {
		return 0
	}
	result := 0
	for _, modifier := range accelerator.Modifiers {
		result |= macModifierMap[modifier]
	}
	return result
}
