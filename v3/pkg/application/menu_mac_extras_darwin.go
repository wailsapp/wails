//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include <stdlib.h>
#include "menu_mac_extras_darwin.h"
*/
import "C"
import (
	"strings"
	"unsafe"
)

// macMenuPaletteEvent is a swatch selection from a palette menu.
type macMenuPaletteEvent struct {
	itemID uint
	index  int
}

var macMenuPaletteSelected = make(chan macMenuPaletteEvent, 5)

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macMenuPaletteSelected
			if item := getMenuItemByID(event.itemID); item != nil {
				item.handlePaletteSelection(event.index)
			}
		}
	})
}

//export processMenuPaletteSelection
func processMenuPaletteSelection(menuItemID C.uint, index C.int) {
	macMenuPaletteSelected <- macMenuPaletteEvent{itemID: uint(menuItemID), index: int(index)}
}

// HandleDockMenu is called by AppDelegate applicationDockMenu: on the main
// thread each time the Dock is about to show the app's menu. It returns
// the NSMenu owned by the Go Menu, or NULL for no menu.
//
//export HandleDockMenu
func HandleDockMenu() unsafe.Pointer {
	app := globalApplication
	if app == nil || app.Menu == nil {
		return nil
	}
	menu := app.Menu.resolveDockMenu()
	if menu == nil {
		return nil
	}
	menu.Update()
	impl, ok := menu.impl.(*macosMenu)
	if !ok || impl == nil {
		return nil
	}
	return impl.nsMenu
}

// newSectionHeaderImpl creates the native item for a sectionHeader entry.
func newSectionHeaderImpl(item *MenuItem) *macosMenuItem {
	return &macosMenuItem{
		menuItem:   item,
		nsMenuItem: C.newMenuItemSectionHeader(C.uint(item.id), C.CString(item.label)),
	}
}

// newPaletteImpl creates the native palette submenu item, or nil below
// macOS 14 (the native side logs the reason).
func newPaletteImpl(item *MenuItem) *macosMenuItem {
	count := len(item.paletteColours)
	var rgba *C.double
	if count > 0 {
		rgba = (*C.double)(C.malloc(C.size_t(count*4) * C.size_t(unsafe.Sizeof(C.double(0)))))
		defer C.free(unsafe.Pointer(rgba))
		values := unsafe.Slice(rgba, count*4)
		for i, colour := range item.paletteColours {
			values[i*4] = C.double(colour.Red) / 255
			values[i*4+1] = C.double(colour.Green) / 255
			values[i*4+2] = C.double(colour.Blue) / 255
			values[i*4+3] = C.double(colour.Alpha) / 255
		}
	}
	ptr := C.newMenuItemPalette(
		C.uint(item.id),
		C.CString(item.label),
		C.CString(strings.Join(item.paletteSymbols, "\n")),
		rgba,
		C.int(count),
		C.int(item.paletteSelected),
	)
	if ptr == nil {
		return nil
	}
	return &macosMenuItem{menuItem: item, nsMenuItem: ptr}
}

func (m macosMenuItem) setSymbol(name string) {
	C.setMenuItemSymbol(m.nsMenuItem, C.CString(name))
}

func (m macosMenuItem) setBadge(text string, count int, present bool) {
	C.setMenuItemBadge(m.nsMenuItem, C.CString(text), C.int(count), C.bool(present))
}

func (m macosMenuItem) setMixed(mixed bool) {
	C.setMenuItemMixed(m.nsMenuItem, C.bool(mixed))
}

func (m macosMenuItem) setAlternate(alternate bool) {
	C.setMenuItemAlternate(m.nsMenuItem, C.bool(alternate))
}

func (m macosMenuItem) setIndentationLevel(level int) {
	C.setMenuItemIndentationLevel(m.nsMenuItem, C.int(level))
}

// applyExtras pushes symbol, badge, mixed state, alternate flag and
// indentation from the model to a freshly created native item.
func (m macosMenuItem) applyExtras() {
	item := m.menuItem
	if item == nil || m.nsMenuItem == nil {
		return
	}
	if item.mixed {
		m.setMixed(true)
	}
	if item.alternate {
		m.setAlternate(true)
	}
	if item.indentationLevel > 0 {
		m.setIndentationLevel(item.indentationLevel)
	}
	if item.hasBadge {
		m.setBadge(item.badgeText, item.badgeCount, true)
	}
	if item.symbol != "" {
		m.setSymbol(item.symbol)
	}
}

// installOpenRecentMenu wires the native Open Recent delegate onto the
// NSMenu created for an OpenRecent role submenu.
func installOpenRecentMenu(nsMenu unsafe.Pointer) {
	C.installOpenRecentMenu(nsMenu)
}

// Recent documents via NSDocumentController (recentDocumentsImpl).

func (m *macosApp) addRecentDocument(path string) {
	InvokeSync(func() {
		C.noteRecentDocument(C.CString(path))
	})
}

func (m *macosApp) clearRecentDocuments() {
	InvokeSync(func() {
		C.clearRecentDocumentsList()
	})
}

func (m *macosApp) recentDocuments() []string {
	return InvokeSyncWithResult(func() []string {
		count := int(C.recentDocumentCount())
		result := make([]string, 0, count)
		for i := 0; i < count; i++ {
			cPath := C.recentDocumentPathAt(C.int(i))
			if cPath == nil {
				continue
			}
			result = append(result, C.GoString(cPath))
			C.free(unsafe.Pointer(cPath))
		}
		return result
	})
}
