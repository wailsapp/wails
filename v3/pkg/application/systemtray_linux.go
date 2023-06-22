//go:build linux

package application

type linuxSystemTray struct {
	id    uint
	label string
	icon  []byte
	menu  *Menu

	iconPosition   int
	isTemplateIcon bool
}

func (s *linuxSystemTray) setIconPosition(position int) {
	s.iconPosition = position
}

func (s *linuxSystemTray) setMenu(menu *Menu) {
	s.menu = menu
}

func (s *linuxSystemTray) run() {
	globalApplication.dispatchOnMainThread(func() {
		// if s.nsStatusItem != nil {
		// 	Fatal("System tray '%d' already running", s.id)
		// }
		//		s.nsStatusItem = unsafe.Pointer(C.systemTrayNew())
		if s.label != "" {
			//			C.systemTraySetLabel(s.nsStatusItem, C.CString(s.label))
		}
		if s.icon != nil {
			//		s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&s.icon[0]), C.int(len(s.icon))))
			//			C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
		}
		if s.menu != nil {
			s.menu.Update()
			// Convert impl to macosMenu object
			//			s.nsMenu = (s.menu.impl).(*macosMenu).nsMenu
			//			C.systemTraySetMenu(s.nsStatusItem, s.nsMenu)
		}

	})
}

func (s *linuxSystemTray) setIcon(icon []byte) {
	s.icon = icon
	globalApplication.dispatchOnMainThread(func() {
		//		s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&icon[0]), C.int(len(icon))))
		//		C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
	})
}

func (s *linuxSystemTray) setDarkModeIcon(icon []byte) {
	s.icon = icon
	globalApplication.dispatchOnMainThread(func() {
		//		s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&icon[0]), C.int(len(icon))))
		//		C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
	})
}

func (s *linuxSystemTray) setTemplateIcon(icon []byte) {
	s.icon = icon
	s.isTemplateIcon = true
	globalApplication.dispatchOnMainThread(func() {
		//		s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&icon[0]), C.int(len(icon))))
		//		C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
	})
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &linuxSystemTray{
		id:             s.id,
		label:          s.label,
		icon:           s.icon,
		menu:           s.menu,
		iconPosition:   s.iconPosition,
		isTemplateIcon: s.isTemplateIcon,
	}
}

func (s *linuxSystemTray) setLabel(label string) {
	s.label = label
	//	C.systemTraySetLabel(s.nsStatusItem, C.CString(label))
}

func (s *linuxSystemTray) destroy() {
	// Remove the status item from the status bar and its associated menu
	//	C.systemTrayDestroy(s.nsStatusItem)
}
