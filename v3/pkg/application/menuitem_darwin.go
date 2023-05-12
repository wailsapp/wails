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
void* newMenuItem(unsigned int menuItemID, char *label, bool disabled, char* tooltip) {
    MenuItem *menuItem = [MenuItem new];

    // Label
    menuItem.title = [NSString stringWithUTF8String:label];

	if( disabled ) {
		[menuItem setTarget:nil];
	} else {
		[menuItem setTarget:menuItem];
	}
    menuItem.menuItemID = menuItemID;
    menuItem.action = @selector(handleClick);
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
        return unicode(0x001c);
    }
    if( [key isEqualToString:@"right"] ) {
        return unicode(0x001d);
    }
    if( [key isEqualToString:@"up"] ) {
        return unicode(0x001e);
    }
    if( [key isEqualToString:@"down"] ) {
        return unicode(0x001f);
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

// Call paste selector to paste text
static void paste(void) {
	[NSApp sendAction:@selector(paste:) to:nil from:nil];
}

// Call copy selector to copy text
static void copy(void) {
	[NSApp sendAction:@selector(copy:) to:nil from:nil];
}

// Call cut selector to cut text
static void cut(void) {
	[NSApp sendAction:@selector(cut:) to:nil from:nil];
}

// Call selectAll selector to select all text
static void selectAll(void) {
	[NSApp sendAction:@selector(selectAll:) to:nil from:nil];
}

// Call delete selector to delete text
static void delete(void) {
	[NSApp sendAction:@selector(delete:) to:nil from:nil];
}

// Call undo selector to undo text
static void undo(void) {
	[NSApp sendAction:@selector(undo:) to:nil from:nil];
}

// Call redo selector to redo text
static void redo(void) {
	[NSApp sendAction:@selector(redo:) to:nil from:nil];
}

// Call startSpeaking selector to start speaking text
static void startSpeaking(void) {
	[NSApp sendAction:@selector(startSpeaking:) to:nil from:nil];
}

// Call stopSpeaking selector to stop speaking text
static void stopSpeaking(void) {
	[NSApp sendAction:@selector(stopSpeaking:) to:nil from:nil];
}

static void pasteAndMatchStyle(void) {
	[NSApp sendAction:@selector(pasteAndMatchStyle:) to:nil from:nil];
}

static void hideApplication(void) {
    [[NSApplication sharedApplication] hide:nil];
}

// hideOthers hides all other applications
static void hideOthers(void) {
	[[NSApplication sharedApplication] hideOtherApplications:nil];
}

// showAll shows all hidden applications
static void showAll(void) {
	[[NSApplication sharedApplication] unhideAllApplications:nil];
}

*/
import "C"
import (
	"runtime"
	"unsafe"
)

type windowsMenuItem struct {
	menuItem *MenuItem

	nsMenuItem unsafe.Pointer
}

func (m windowsMenuItem) setTooltip(tooltip string) {
	C.setMenuItemTooltip(m.nsMenuItem, C.CString(tooltip))
}

func (m windowsMenuItem) setLabel(s string) {
	C.setMenuItemLabel(m.nsMenuItem, C.CString(s))
}

func (m windowsMenuItem) setDisabled(disabled bool) {
	C.setMenuItemDisabled(m.nsMenuItem, C.bool(disabled))
}

func (m windowsMenuItem) setChecked(checked bool) {
	C.setMenuItemChecked(m.nsMenuItem, C.bool(checked))
}

func (m windowsMenuItem) setHidden(hidden bool) {
	C.setMenuItemHidden(m.nsMenuItem, C.bool(hidden))
}

func (m windowsMenuItem) setAccelerator(accelerator *accelerator) {
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

func newMenuItemImpl(item *MenuItem) *windowsMenuItem {
	result := &windowsMenuItem{
		menuItem: item,
	}

	switch item.itemType {
	case text, checkbox, submenu, radio:
		result.nsMenuItem = unsafe.Pointer(C.newMenuItem(C.uint(item.id), C.CString(item.label), C.bool(item.disabled), C.CString(item.tooltip)))
		if item.itemType == checkbox || item.itemType == radio {
			C.setMenuItemChecked(result.nsMenuItem, C.bool(item.checked))
		}
		if item.accelerator != nil {
			result.setAccelerator(item.accelerator)
		}
	default:
		panic("WTF")
	}
	return result
}

func newSpeechMenu() *MenuItem {
	speechMenu := NewMenu()
	speechMenu.Add("Start Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+.").
		OnClick(func(ctx *Context) {
			C.startSpeaking()
		})
	speechMenu.Add("Stop Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+,").
		OnClick(func(ctx *Context) {
			C.stopSpeaking()
		})
	subMenu := newSubMenuItem("Speech")
	subMenu.submenu = speechMenu
	return subMenu
}

func newHideMenuItem() *MenuItem {
	return newMenuItem("Hide " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+h").
		OnClick(func(ctx *Context) {
			C.hideApplication()
		})
}

func newHideOthersMenuItem() *MenuItem {
	return newMenuItem("Hide Others").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+h").
		OnClick(func(ctx *Context) {
			C.hideOthers()
		})
}

func newUnhideMenuItem() *MenuItem {
	return newMenuItem("Show All").
		OnClick(func(ctx *Context) {
			C.showAll()
		})
}

func newUndoMenuItem() *MenuItem {
	return newMenuItem("Undo").
		SetAccelerator("CmdOrCtrl+z").
		OnClick(func(ctx *Context) {
			C.undo()
		})
}

// newRedoMenuItem creates a new menu item for redoing the last action
func newRedoMenuItem() *MenuItem {
	return newMenuItem("Redo").
		SetAccelerator("CmdOrCtrl+Shift+z").
		OnClick(func(ctx *Context) {
			C.redo()
		})
}

func newCutMenuItem() *MenuItem {
	return newMenuItem("Cut").
		SetAccelerator("CmdOrCtrl+x").
		OnClick(func(ctx *Context) {
			C.cut()
		})
}

func newCopyMenuItem() *MenuItem {
	return newMenuItem("Copy").
		SetAccelerator("CmdOrCtrl+c").
		OnClick(func(ctx *Context) {
			C.copy()
		})
}

func newPasteMenuItem() *MenuItem {
	return newMenuItem("Paste").
		SetAccelerator("CmdOrCtrl+v").
		OnClick(func(ctx *Context) {
			C.paste()
		})
}

func newPasteAndMatchStyleMenuItem() *MenuItem {
	return newMenuItem("Paste and Match Style").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+v").
		OnClick(func(ctx *Context) {
			C.pasteAndMatchStyle()
		})
}

func newDeleteMenuItem() *MenuItem {
	return newMenuItem("Delete").
		SetAccelerator("backspace").
		OnClick(func(ctx *Context) {
			C.delete()
		})
}

func newQuitMenuItem() *MenuItem {
	return newMenuItem("Quit " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+q").
		OnClick(func(ctx *Context) {
			globalApplication.Quit()
		})
}

func newSelectAllMenuItem() *MenuItem {
	return newMenuItem("Select All").
		SetAccelerator("CmdOrCtrl+a").
		OnClick(func(ctx *Context) {
			C.selectAll()
		})
}

func newAboutMenuItem() *MenuItem {
	return newMenuItem("About " + globalApplication.options.Name).
		OnClick(func(ctx *Context) {
			globalApplication.ShowAboutDialog()
		})
}

func newCloseMenuItem() *MenuItem {
	return newMenuItem("Close").
		SetAccelerator("CmdOrCtrl+w").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Close()
			}
		})
}

func newReloadMenuItem() *MenuItem {
	return newMenuItem("Reload").
		SetAccelerator("CmdOrCtrl+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Reload()
			}
		})
}

func newForceReloadMenuItem() *MenuItem {
	return newMenuItem("Force Reload").
		SetAccelerator("CmdOrCtrl+Shift+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ForceReload()
			}
		})
}

func newToggleFullscreenMenuItem() *MenuItem {
	result := newMenuItem("Toggle Full Screen").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ToggleFullscreen()
			}
		})
	if runtime.GOOS == "darwin" {
		result.SetAccelerator("Ctrl+Command+F")
	} else {
		result.SetAccelerator("F11")
	}
	return result
}

func newToggleDevToolsMenuItem() *MenuItem {
	return newMenuItem("Toggle Developer Tools").
		SetAccelerator("Alt+Command+I").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ToggleDevTools()
			}
		})
}

func newZoomResetMenuItem() *MenuItem {
	// reset zoom menu item
	return newMenuItem("Actual Size").
		SetAccelerator("CmdOrCtrl+0").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomReset()
			}
		})
}

func newZoomInMenuItem() *MenuItem {
	return newMenuItem("Zoom In").
		SetAccelerator("CmdOrCtrl+plus").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomIn()
			}
		})
}

func newZoomOutMenuItem() *MenuItem {
	return newMenuItem("Zoom Out").
		SetAccelerator("CmdOrCtrl+-").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomOut()
			}
		})
}

func newMinimizeMenuItem() *MenuItem {
	return newMenuItem("Minimize").
		SetAccelerator("CmdOrCtrl+M").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Minimise()
			}
		})
}

func newZoomMenuItem() *MenuItem {
	return newMenuItem("Zoom").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Zoom()
			}
		})
}

func newFullScreenMenuItem() *MenuItem {
	return newMenuItem("Fullscreen").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Fullscreen()
			}
		})
}
