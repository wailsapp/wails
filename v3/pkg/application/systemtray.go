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
	setTemplateIcon(icon []byte)
	destroy()
	setDarkModeIcon(icon []byte)
}

type SystemTray struct {
	id           uint
	label        string
	icon         []byte
	darkModeIcon []byte
	iconPosition int

	leftButtonClickHandler        func()
	rightButtonClickHandler       func()
	leftButtonDoubleClickHandler  func()
	rightButtonDoubleClickHandler func()
	mouseEnterHandler             func()
	mouseLeaveHandler             func()

	// Platform specific implementation
	impl           systemTrayImpl
	menu           *Menu
	isTemplateIcon bool
}

func NewSystemTray(id uint) *SystemTray {
	return &SystemTray{
		id:           id,
		label:        "",
		iconPosition: NSImageLeading,
		icon:         DefaultApplicationIcon,
	}
}

func (s *SystemTray) SetLabel(label string) {
	if s.impl == nil {
		s.label = label
		return
	}
	invokeSync(func() {
		s.impl.setLabel(label)
	})
}

func (s *SystemTray) Label() string {
	return s.label
}

func (s *SystemTray) Run() {
	s.impl = newSystemTrayImpl(s)
	invokeSync(s.impl.run)
}

func (s *SystemTray) SetIcon(icon []byte) *SystemTray {
	if s.impl == nil {
		s.icon = icon
	} else {
		invokeSync(func() {
			s.impl.setIcon(icon)
		})
	}
	return s
}

func (s *SystemTray) SetDarkModeIcon(icon []byte) *SystemTray {
	if s.impl == nil {
		s.darkModeIcon = icon
	} else {
		invokeSync(func() {
			s.impl.setDarkModeIcon(icon)
		})
	}
	return s
}

func (s *SystemTray) SetMenu(menu *Menu) *SystemTray {
	if s.impl == nil {
		s.menu = menu
	} else {
		invokeSync(func() {
			s.impl.setMenu(menu)
		})
	}
	return s
}

func (s *SystemTray) SetIconPosition(iconPosition int) *SystemTray {
	if s.impl == nil {
		s.iconPosition = iconPosition
	} else {
		invokeSync(func() {
			s.impl.setIconPosition(iconPosition)
		})
	}
	return s
}

func (s *SystemTray) SetTemplateIcon(icon []byte) *SystemTray {
	if s.impl == nil {
		s.icon = icon
		s.isTemplateIcon = true
	} else {
		invokeSync(func() {
			s.impl.setTemplateIcon(icon)
		})
	}
	return s
}

func (s *SystemTray) Destroy() {
	if s.impl == nil {
		return
	}
	s.impl.destroy()
}

func (s *SystemTray) OnLeftButtonClick(handler func()) *SystemTray {
	s.leftButtonClickHandler = handler
	return s
}

func (s *SystemTray) OnRightButtonClick(handler func()) *SystemTray {
	s.rightButtonClickHandler = handler
	return s
}

func (s *SystemTray) OnLeftButtonDoubleClick(handler func()) *SystemTray {
	s.leftButtonDoubleClickHandler = handler
	return s
}

func (s *SystemTray) OnRightButtonDoubleClick(handler func()) *SystemTray {
	s.rightButtonDoubleClickHandler = handler
	return s
}

func (s *SystemTray) OnMouseEnter(handler func()) *SystemTray {
	s.mouseEnterHandler = handler
	return s
}

func (s *SystemTray) OnMouseLeave(handler func()) *SystemTray {
	s.mouseLeaveHandler = handler
	return s
}
