package application

type IconPosition int

const (
	NSImageNone = iota
	NSImageOnly
	NSImageLeft
	NSImageRight
	NSImageBelow
	NSImageAbove
	NSImageOverlaps
	NSImageLeading
	NSImageTrailing
)

type systemTrayImpl interface {
	setLabel(label string)
	run()
	setIcon(icon []byte)
	setMenu(menu *Menu)
	setIconPosition(position int)
}

type SystemTray struct {
	id           uint
	label        string
	icon         []byte
	iconPosition int

	// Platform specific implementation
	impl systemTrayImpl
	menu *Menu
}

func NewSystemTray(id uint) *SystemTray {
	return &SystemTray{
		id:           id,
		label:        "",
		iconPosition: NSImageLeading,
	}
}

func (s *SystemTray) SetLabel(label string) {
	if s.impl == nil {
		s.label = label
		return
	}
	s.impl.setLabel(label)
}

func (s *SystemTray) Label() string {
	return s.label
}

func (s *SystemTray) Run() {
	s.impl = newSystemTrayImpl(s)
	s.impl.run()
}

func (s *SystemTray) SetIcon(icon []byte) *SystemTray {
	if s.impl == nil {
		s.icon = icon
	} else {
		s.impl.setIcon(icon)
	}
	return s
}

func (s *SystemTray) SetMenu(menu *Menu) *SystemTray {
	if s.impl == nil {
		s.menu = menu
	} else {
		s.impl.setMenu(menu)
	}
	return s
}

func (s *SystemTray) SetIconPosition(iconPosition int) *SystemTray {
	if s.impl == nil {
		s.iconPosition = iconPosition
	} else {
		s.impl.setIconPosition(iconPosition)
	}
	return s
}
