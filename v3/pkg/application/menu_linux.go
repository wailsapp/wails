//go:build linux

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>

void handleClick(void*);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

var (
	gtkSignalHandlers   map[*C.GtkWidget]C.gulong
	gtkSignalToMenuItem map[*C.GtkWidget]*MenuItem
)

func init() {
	gtkSignalHandlers = map[*C.GtkWidget]C.gulong{}
	gtkSignalToMenuItem = map[*C.GtkWidget]*MenuItem{}
}

//export handleClick
func handleClick(idPtr unsafe.Pointer) {
	id := (*C.GtkWidget)(idPtr)
	item, ok := gtkSignalToMenuItem[id]
	if !ok {
		return
	}

	//impl := (item.impl).(*linuxMenuItem)

	switch item.itemType {
	case text, checkbox:
		processMenuItemClick(C.uint(item.id))
	case radio:
		menuItem := (item.impl).(*linuxMenuItem)
		if menuItem.isChecked() {
			processMenuItemClick(C.uint(item.id))
		}
	default:
		fmt.Println("handleClick", item.itemType, item.id)
	}
}

type linuxMenu struct {
	menu   *Menu
	native unsafe.Pointer
}

func newMenuImpl(menu *Menu) *linuxMenu {
	result := &linuxMenu{
		menu:   menu,
		native: unsafe.Pointer(C.gtk_menu_bar_new()),
	}
	return result
}

func (m *linuxMenu) update() {
	//	fmt.Println("linuxMenu.update()")
	// if m.native != nil {
	// 	C.gtk_widget_destroy((*C.GtkWidget)(m.native))
	// 	m.native = unsafe.Pointer(C.gtk_menu_new())
	// }
	m.processMenu(m.menu)
}

func (m *linuxMenu) processMenu(menu *Menu) {
	if menu.impl == nil {
		menu.impl = &linuxMenu{
			menu:   menu,
			native: unsafe.Pointer(C.gtk_menu_new()),
		}
	}
	var currentRadioGroup *C.GSList

	for _, item := range menu.items {
		// drop the group if we have run out of radio items
		if item.itemType != radio {
			currentRadioGroup = nil
		}

		switch item.itemType {
		case submenu:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			m.processMenu(item.submenu)
			m.addSubMenuToItem(item.submenu, item)
			m.addMenuItem(menu, item)
		case text, checkbox:
			menuItem := newMenuItemImpl(item)
			item.impl = menuItem
			m.addMenuItem(menu, item)
		case radio:
			menuItem := newRadioItemImpl(item, currentRadioGroup)
			item.impl = menuItem
			m.addMenuItem(menu, item)
			currentRadioGroup = C.gtk_radio_menu_item_get_group((*C.GtkRadioMenuItem)(menuItem.native))
		case separator:
			m.addMenuSeparator(menu)
		}

	}

	for _, item := range menu.items {
		if item.callback != nil {
			m.attachHandler(item)
		}
	}

}

func (m *linuxMenu) attachHandler(item *MenuItem) {
	signal := C.CString("activate")
	defer C.free(unsafe.Pointer(signal))

	impl := (item.impl).(*linuxMenuItem)
	widget := impl.native
	flags := C.GConnectFlags(0)
	handlerId := C.g_signal_connect_object(
		C.gpointer(widget),
		signal,
		C.GCallback(C.handleClick),
		C.gpointer(widget),
		flags)

	id := (*C.GtkWidget)(widget)
	gtkSignalToMenuItem[id] = item
	gtkSignalHandlers[id] = handlerId
	impl.handlerId = handlerId
}

func (m *linuxMenu) addSubMenuToItem(menu *Menu, item *MenuItem) {
	if menu.impl == nil {
		menu.impl = &linuxMenu{
			menu:   menu,
			native: unsafe.Pointer(C.gtk_menu_new()),
		}
	}

	C.gtk_menu_item_set_submenu(
		(*C.GtkMenuItem)((item.impl).(*linuxMenuItem).native),
		(*C.GtkWidget)((menu.impl).(*linuxMenu).native))

	if item.role == ServicesMenu {
		// FIXME: what does this mean?
	}
}

func (m *linuxMenu) addMenuItem(parent *Menu, menu *MenuItem) {
	//	fmt.Println("addMenuIteam", fmt.Sprintf("%+v", parent), fmt.Sprintf("%+v", menu))
	C.gtk_menu_shell_append(
		(*C.GtkMenuShell)((parent.impl).(*linuxMenu).native),
		(*C.GtkWidget)((menu.impl).(*linuxMenuItem).native),
	)
	/*
		C.gtk_menu_item_set_submenu(
			(*C.struct__GtkMenuItem)((menu.impl).(*linuxMenuItem).native),
			(*C.struct__GtkWidget)((parent.impl).(*linuxMenu).native),
		)
	*/
}

func (m *linuxMenu) addMenuSeparator(menu *Menu) {
	//	fmt.Println("addMenuSeparator", fmt.Sprintf("%+v", menu))
	sep := C.gtk_separator_menu_item_new()
	native := (menu.impl).(*linuxMenu).native
	C.gtk_menu_shell_append((*C.GtkMenuShell)(native), sep)
}

func (m *linuxMenu) addServicesMenu(menu *Menu) {
	fmt.Println("addServicesMenu - not implemented")
	//C.addServicesMenu(unsafe.Pointer(menu.impl.(*linuxMenu).nsMenu))
}

func (l *linuxMenu) createMenu(name string, items []*MenuItem) *Menu {
	impl := newMenuImpl(&Menu{label: name})
	menu := &Menu{
		label: name,
		items: items,
		impl:  impl,
	}
	impl.menu = menu
	return menu
}
