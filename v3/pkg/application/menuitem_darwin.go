//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "Cocoa/Cocoa.h"
#include "menuitem_darwin.h"
#include "application_darwin.h"

#define unicode(input) [NSString stringWithFormat:@"%C", input]

// Create menu item
void* newMenuItem(unsigned int menuItemID, char *label, bool disabled, char* tooltip, char* selector) {
    NSLog(@"newMenuItem: %s, %s\n", label, selector);
	MenuItem *menuItem = [MenuItem new];

    // Label
    menuItem.title = [NSString stringWithUTF8String:label];

	if( disabled ) {
		[menuItem setTarget:nil];
	} else {
		if (selector != NULL) {
			menuItem.action = NSSelectorFromString([NSString stringWithUTF8String:selector]);
			menuItem.target = nil; // Allow the action to be sent up the responder chain
		} else {
			menuItem.action = @selector(handleClick);
			menuItem.target = menuItem;
		}
	}
    menuItem.menuItemID = menuItemID;

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
	dispatch_async(dispatch_get_main_queue(), ^{
		MenuItem *menuItem = (MenuItem *)nsMenuItem;
		[menuItem setEnabled:!disabled];
		// remove target
		if( disabled ) {
			[menuItem setTarget:nil];
		} else {
			[menuItem setTarget:menuItem];
		}
	});
}

// set menu item hidden
void setMenuItemHidden(void* nsMenuItem, bool hidden) {
	dispatch_async(dispatch_get_main_queue(), ^{
		MenuItem *menuItem = (MenuItem *)nsMenuItem;
		[menuItem setHidden:hidden];
	});
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

NSString* translateKey(NSString* key) {

    // Guard against no accelerator key
    if( key == NULL ) {
        return @"";
    }

    if( [key isEqualToString:@"backspace"] ) {
        return unicode(0x0008);
    }
    if( [key isEqualToString:@"tab"] ) {
        return unicode(0x0009);
    }
    if( [key isEqualToString:@"return"] ) {
        return unicode(0x000d);
    }
    if( [key isEqualToString:@"enter"] ) {
        return unicode(0x000d);
    }
    if( [key isEqualToString:@"escape"] ) {
        return unicode(0x001b);
    }
    if( [key isEqualToString:@"left"] ) {
        return unicode(0xf702);
    }
    if( [key isEqualToString:@"right"] ) {
        return unicode(0xf703);
    }
    if( [key isEqualToString:@"up"] ) {
        return unicode(0xf700);
    }
    if( [key isEqualToString:@"down"] ) {
        return unicode(0xf701);
    }
    if( [key isEqualToString:@"space"] ) {
        return unicode(0x0020);
    }
    if( [key isEqualToString:@"delete"] ) {
        return unicode(0x007f);
    }
    if( [key isEqualToString:@"home"] ) {
        return unicode(0x2196);
    }
    if( [key isEqualToString:@"end"] ) {
        return unicode(0x2198);
    }
    if( [key isEqualToString:@"page up"] ) {
        return unicode(0x21de);
    }
    if( [key isEqualToString:@"page down"] ) {
        return unicode(0x21df);
    }
    if( [key isEqualToString:@"f1"] ) {
        return unicode(0xf704);
    }
    if( [key isEqualToString:@"f2"] ) {
        return unicode(0xf705);
    }
    if( [key isEqualToString:@"f3"] ) {
        return unicode(0xf706);
    }
    if( [key isEqualToString:@"f4"] ) {
        return unicode(0xf707);
    }
    if( [key isEqualToString:@"f5"] ) {
        return unicode(0xf708);
    }
    if( [key isEqualToString:@"f6"] ) {
        return unicode(0xf709);
    }
    if( [key isEqualToString:@"f7"] ) {
        return unicode(0xf70a);
    }
    if( [key isEqualToString:@"f8"] ) {
        return unicode(0xf70b);
    }
    if( [key isEqualToString:@"f9"] ) {
        return unicode(0xf70c);
    }
    if( [key isEqualToString:@"f10"] ) {
        return unicode(0xf70d);
    }
    if( [key isEqualToString:@"f11"] ) {
        return unicode(0xf70e);
    }
    if( [key isEqualToString:@"f12"] ) {
        return unicode(0xf70f);
    }
    if( [key isEqualToString:@"f13"] ) {
        return unicode(0xf710);
    }
    if( [key isEqualToString:@"f14"] ) {
        return unicode(0xf711);
    }
    if( [key isEqualToString:@"f15"] ) {
        return unicode(0xf712);
    }
    if( [key isEqualToString:@"f16"] ) {
        return unicode(0xf713);
    }
    if( [key isEqualToString:@"f17"] ) {
        return unicode(0xf714);
    }
    if( [key isEqualToString:@"f18"] ) {
        return unicode(0xf715);
    }
    if( [key isEqualToString:@"f19"] ) {
        return unicode(0xf716);
    }
    if( [key isEqualToString:@"f20"] ) {
        return unicode(0xf717);
    }
    if( [key isEqualToString:@"f21"] ) {
        return unicode(0xf718);
    }
    if( [key isEqualToString:@"f22"] ) {
        return unicode(0xf719);
    }
    if( [key isEqualToString:@"f23"] ) {
        return unicode(0xf71a);
    }
    if( [key isEqualToString:@"f24"] ) {
        return unicode(0xf71b);
    }
    if( [key isEqualToString:@"f25"] ) {
        return unicode(0xf71c);
    }
    if( [key isEqualToString:@"f26"] ) {
        return unicode(0xf71d);
    }
    if( [key isEqualToString:@"f27"] ) {
        return unicode(0xf71e);
    }
    if( [key isEqualToString:@"f28"] ) {
        return unicode(0xf71f);
    }
    if( [key isEqualToString:@"f29"] ) {
        return unicode(0xf720);
    }
    if( [key isEqualToString:@"f30"] ) {
        return unicode(0xf721);
    }
    if( [key isEqualToString:@"f31"] ) {
        return unicode(0xf722);
    }
    if( [key isEqualToString:@"f32"] ) {
        return unicode(0xf723);
    }
    if( [key isEqualToString:@"f33"] ) {
        return unicode(0xf724);
    }
    if( [key isEqualToString:@"f34"] ) {
        return unicode(0xf725);
    }
    if( [key isEqualToString:@"f35"] ) {
        return unicode(0xf726);
    }
    if( [key isEqualToString:@"numLock"] ) {
        return unicode(0xf739);
    }
    return key;
}

// Set the menuitem key equivalent
void setMenuItemKeyEquivalent(void* nsMenuItem, char *key, int modifier) {
	MenuItem *menuItem = (MenuItem *)nsMenuItem;
	NSString *nskey = [NSString stringWithUTF8String:key];
	menuItem.keyEquivalent = translateKey(nskey);
	menuItem.keyEquivalentModifierMask = modifier;
	free(key);
}

// Call the copy selector on the pasteboard
static void copyToPasteboard(char *text) {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	[pasteboard clearContents];
	[pasteboard setString:[NSString stringWithUTF8String:text] forType:NSPasteboardTypeString];
}

// Call the paste selector on the pasteboard
static char *pasteFromPasteboard(void) {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	NSString *text = [pasteboard stringForType:NSPasteboardTypeString];
	if( text == nil ) {
		return NULL;
	}
	return strdup([text UTF8String]);
}

void performSelectorOnMainThreadForFirstResponder(SEL selector) {
    NSWindow *activeWindow = [[NSApplication sharedApplication] keyWindow];
    if (activeWindow) {
		[activeWindow performSelectorOnMainThread:selector withObject:nil waitUntilDone:YES];
    }
}

void setMenuItemBitmap(void* nsMenuItem, unsigned char *bitmap, int length) {
	MenuItem *menuItem = (MenuItem *)nsMenuItem;
	NSImage *image = [[NSImage alloc] initWithData:[NSData dataWithBytes:bitmap length:length]];
	[menuItem setImage:image];
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

func (m macosMenuItem) setHidden(hidden bool) {
	C.setMenuItemHidden(m.nsMenuItem, C.bool(hidden))
}

func (m macosMenuItem) setBitmap(bitmap []byte) {
	C.setMenuItemBitmap(m.nsMenuItem, (*C.uchar)(&bitmap[0]), C.int(len(bitmap)))
}

func (m macosMenuItem) setAccelerator(accelerator *accelerator) {
	// Set the keyboard shortcut of the menu item
	var modifier C.int
	var key *C.char
	if accelerator != nil {
		modifier = C.int(toMacModifier(accelerator.Modifiers))
		key = C.CString(accelerator.Key)
	}

	// Convert the key to a string
	C.setMenuItemKeyEquivalent(m.nsMenuItem, key, modifier)
}

func newMenuItemImpl(item *MenuItem) *macosMenuItem {
	result := &macosMenuItem{
		menuItem: item,
	}

	selector := getSelectorForRole(item.role)
	if selector != nil {
		defer C.free(unsafe.Pointer(selector))
	}
	result.nsMenuItem = unsafe.Pointer(C.newMenuItem(
		C.uint(item.id),
		C.CString(item.label),
		C.bool(item.disabled),
		C.CString(item.tooltip),
		selector,
	))

	switch item.itemType {
	case checkbox, radio:
		C.setMenuItemChecked(result.nsMenuItem, C.bool(item.checked))
	}

	if item.accelerator != nil {
		result.setAccelerator(item.accelerator)
	}
	return result
}

func newSpeechMenu() *MenuItem {
	speechMenu := NewMenu()
	speechMenu.AddRole(StartSpeaking)
	speechMenu.AddRole(StopSpeaking)
	subMenu := NewSubMenuItem("Speech")
	subMenu.submenu = speechMenu
	return subMenu
}

func newHideMenuItem() *MenuItem {
	return NewMenuItem("Hide " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+h").
		SetRole(Hide)
}

func newHideOthersMenuItem() *MenuItem {
	return NewMenuItem("Hide Others").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+h").
		SetRole(HideOthers)
}

func newFrontMenuItem() *MenuItem {
	return NewMenuItem("Bring All to Front")
}

func newUnhideMenuItem() *MenuItem {
	return NewMenuItem("Show All")
}

func newUndoMenuItem() *MenuItem {
	return NewMenuItem("Undo").
		SetAccelerator("CmdOrCtrl+z")
}

// newRedoMenuItem creates a new menu item for redoing the last action
func newRedoMenuItem() *MenuItem {
	return NewMenuItem("Redo").
		SetAccelerator("CmdOrCtrl+Shift+z")
}

func newCutMenuItem() *MenuItem {
	return NewMenuItem("Cut").
		SetAccelerator("CmdOrCtrl+x")
}

func newCopyMenuItem() *MenuItem {
	return NewMenuItem("Copy").
		SetAccelerator("CmdOrCtrl+c")
}

func newPasteMenuItem() *MenuItem {
	return NewMenuItem("Paste").
		SetAccelerator("CmdOrCtrl+v")
}

func newPasteAndMatchStyleMenuItem() *MenuItem {
	return NewMenuItem("Paste and Match Style").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+v")
}

func newDeleteMenuItem() *MenuItem {
	return NewMenuItem("Delete").
		SetAccelerator("backspace")
}

func newQuitMenuItem() *MenuItem {
	return NewMenuItem("Quit " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+q").
		OnClick(func(ctx *Context) {
			globalApplication.Quit()
		})
}

func newSelectAllMenuItem() *MenuItem {
	return NewMenuItem("Select All").
		SetAccelerator("CmdOrCtrl+a")
}

func newAboutMenuItem() *MenuItem {
	return NewMenuItem("About " + globalApplication.options.Name).
		OnClick(func(ctx *Context) {
			globalApplication.ShowAboutDialog()
		})
}

func newCloseMenuItem() *MenuItem {
	return NewMenuItem("Close").
		SetAccelerator("CmdOrCtrl+w")
}

func newReloadMenuItem() *MenuItem {
	return NewMenuItem("Reload").
		SetAccelerator("CmdOrCtrl+r")
}

func newForceReloadMenuItem() *MenuItem {
	return NewMenuItem("Force Reload").
		SetAccelerator("CmdOrCtrl+Shift+r")
}

func newToggleFullscreenMenuItem() *MenuItem {
	return NewMenuItem("Toggle Full Screen").
		SetAccelerator("Ctrl+Command+F")
}

func newZoomResetMenuItem() *MenuItem {
	// reset zoom menu item
	return NewMenuItem("Actual Size").
		SetAccelerator("CmdOrCtrl+0")
}

func newZoomInMenuItem() *MenuItem {
	return NewMenuItem("Zoom In").
		SetAccelerator("CmdOrCtrl+plus")
}

func newZoomOutMenuItem() *MenuItem {
	return NewMenuItem("Zoom Out").
		SetAccelerator("CmdOrCtrl+-")
}

func newMinimizeMenuItem() *MenuItem {
	return NewMenuItem("Minimize").
		SetAccelerator("CmdOrCtrl+M")
}

func newZoomMenuItem() *MenuItem {
	return NewMenuItem("Zoom")
}

func newFullScreenMenuItem() *MenuItem {
	return NewMenuItem("Fullscreen")
}

func newPrintMenuItem() *MenuItem {
	return NewMenuItem("Print").
		SetAccelerator("CmdOrCtrl+p")
}

func newPageLayoutMenuItem() *MenuItem {
	return NewMenuItem("Page Setup...").
		SetAccelerator("CmdOrCtrl+Shift+p")
}

// -----

func newShowAllMenuItem() *MenuItem {
	return NewMenuItem("Show All")
}

func newBringAllToFrontMenuItem() *MenuItem {
	return NewMenuItem("Bring All to Front")
}

func newNewFileMenuItem() *MenuItem {
	return NewMenuItem("New File").
		SetAccelerator("CmdOrCtrl+n")
}

func newOpenMenuItem() *MenuItem {
	return NewMenuItem("Open...").
		SetAccelerator("CmdOrCtrl+o")
}

func newSaveMenuItem() *MenuItem {
	return NewMenuItem("Save").
		SetAccelerator("CmdOrCtrl+s")
}

func newSaveAsMenuItem() *MenuItem {
	return NewMenuItem("Save As...").
		SetAccelerator("CmdOrCtrl+Shift+s")
}

func newStartSpeakingMenuItem() *MenuItem {
	return NewMenuItem("Start Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+.")
}

func newStopSpeakingMenuItem() *MenuItem {
	return NewMenuItem("Stop Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+,")
}

func newRevertMenuItem() *MenuItem {
	return NewMenuItem("Revert").
		SetAccelerator("CmdOrCtrl+r")
}

func newFindMenuItem() *MenuItem {
	return NewMenuItem("Find...").
		SetAccelerator("CmdOrCtrl+f")
}

func newFindAndReplaceMenuItem() *MenuItem {
	return NewMenuItem("Find and Replace...").
		SetAccelerator("CmdOrCtrl+Shift+f")
}

func newFindNextMenuItem() *MenuItem {
	return NewMenuItem("Find Next").
		SetAccelerator("CmdOrCtrl+g")
}

func newFindPreviousMenuItem() *MenuItem {
	return NewMenuItem("Find Previous").
		SetAccelerator("CmdOrCtrl+Shift+g")
}

func newHelpMenuItem() *MenuItem {
	return NewMenuItem("Help").
		SetAccelerator("CmdOrCtrl+?")
}
