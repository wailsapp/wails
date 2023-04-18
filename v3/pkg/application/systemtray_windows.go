//go:build windows

package application

type windowsSystemTray struct {
	id   uint
	icon []byte
	menu *Menu
}

func (s *windowsSystemTray) setIconPosition(position int) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) setMenu(menu *Menu) {
	s.menu = menu
}

func (s *windowsSystemTray) run() {
	globalApplication.dispatchOnMainThread(func() {
		//if s.nsStatusItem != nil {
		//	Fatal("System tray '%d' already running", s.id)
		//}
		//s.nsStatusItem = unsafe.Pointer(C.systemTrayNew())
		//if s.label != "" {
		//	C.systemTraySetLabel(s.nsStatusItem, C.CString(s.label))
		//}
		//if s.icon != nil {
		//	s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&s.icon[0]), C.int(len(s.icon))))
		//	C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
		//}
		//if s.menu != nil {
		//	s.menu.Update()
		//	// Convert impl to macosMenu object
		//	s.nsMenu = (s.menu.impl).(*macosMenu).nsMenu
		//	C.systemTraySetMenu(s.nsStatusItem, s.nsMenu)
		//}
		panic("implement me")
	})
}

func (s *windowsSystemTray) setIcon(icon []byte) {
	s.icon = icon
	globalApplication.dispatchOnMainThread(func() {
		//s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&icon[0]), C.int(len(icon))))
		//C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
	})
}

func (s *windowsSystemTray) setTemplateIcon(icon []byte) {
	// Unsupported - do nothing
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &windowsSystemTray{
		id:   s.id,
		icon: s.icon,
		menu: s.menu,
	}
}

func (s *windowsSystemTray) setLabel(label string) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) destroy() {
	panic("implement me")
}
