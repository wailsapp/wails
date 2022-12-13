package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "Cocoa/Cocoa.h"
#include "menuitem.h"

// Create menu item
void* newMenuItem(unsigned int menuItemID, char *label, bool disabled, char* tooltip) {
    MenuItem *menuItem = [MenuItem new];

    // Label
    menuItem.title = [NSString stringWithUTF8String:label];

    // Process callback
    menuItem.menuItemID = menuItemID;
    menuItem.action = @selector(handleClick);
    menuItem.target = menuItem;
	menuItem.enabled = !disabled;

	// Tooltip
	if( tooltip != NULL ) {
		menuItem.toolTip = [NSString stringWithUTF8String:tooltip];
		free(tooltip);
	}

	// Set the tag
	[menuItem setTag:menuItemID];

	return (void*)menuItem;
}

// set menu item label
void setMenuItemLabel(void* nsMenuItem, char *label) {
	MenuItem *menuItem = (MenuItem *)nsMenuItem;
	menuItem.title = [NSString stringWithUTF8String:label];
}

// set menu item disabled
void setMenuItemDisabled(void* nsMenuItem, bool disabled) {
	MenuItem *menuItem = (MenuItem *)nsMenuItem;
	menuItem.enabled = !disabled;
}

// set menu item tooltip
void setMenuItemTooltip(void* nsMenuItem, char *tooltip) {
	MenuItem *menuItem = (MenuItem *)nsMenuItem;
	menuItem.toolTip = [NSString stringWithUTF8String:tooltip];
}

// Check menu item
void setMenuItemChecked(void* nsMenuItem, bool checked) {
	MenuItem *menuItem = (MenuItem *)nsMenuItem;
	menuItem.state = checked ? NSControlStateValueOn : NSControlStateValueOff;
}

*/
import "C"
import (
	"unsafe"
)

type macosMenuItem struct {
	menuItem *MenuItem

	nsMenuItem unsafe.Pointer
}

func (m macosMenuItem) setTooltip(tooltip string) {
	C.setMenuItemTooltip(m.nsMenuItem, C.CString(tooltip))
}

func (m macosMenuItem) setLabel(s string) {
	C.setMenuItemLabel(m.nsMenuItem, C.CString(s))
}

func (m macosMenuItem) setDisabled(disabled bool) {
	C.setMenuItemDisabled(m.nsMenuItem, C.bool(disabled))
}

func (m macosMenuItem) setChecked(checked bool) {
	C.setMenuItemChecked(m.nsMenuItem, C.bool(checked))
}

func newMenuItemImpl(item *MenuItem) *macosMenuItem {
	result := &macosMenuItem{
		menuItem: item,
	}
	switch item.itemType {
	case text, checkbox:
		result.nsMenuItem = unsafe.Pointer(C.newMenuItem(C.uint(item.id), C.CString(item.label), C.bool(item.disabled), C.CString(item.tooltip)))
		if item.itemType == checkbox {
			C.setMenuItemChecked(result.nsMenuItem, C.bool(item.checked))
		}
	default:
		panic("WTF")
	}
	return result
}
