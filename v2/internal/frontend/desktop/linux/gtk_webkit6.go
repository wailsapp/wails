//go:build linux && webkit_6
// +build linux,webkit_6

package linux

import "C"
import (
	"github.com/wailsapp/wails/v2/pkg/menu"
)

//export handleMenuRadioItemClick
func handleMenuRadioItemClick(rName *C.char, prev *C.char, curr *C.char) {
	radioActionName := C.GoString(rName)
	prevId := C.GoString(prev)
	itemId := C.GoString(curr)

	actionName := radioActionName + "::" + itemId
	it, ok := gActionIdToMenuItem.Load(actionName)
	if !ok {
		return
	}

	item := it.(*menu.MenuItem)

	prevActionId := radioActionName + "::" + prevId
	prevIt, ok := gActionIdToMenuItem.Load(prevActionId)
	if !ok {
		return
	}

	prevItem := prevIt.(*menu.MenuItem)

	prevItem.Checked = false
	item.Checked = true

	go item.Click(&menu.CallbackData{MenuItem: item})
}

//export handleMenuCheckItemClick
func handleMenuCheckItemClick(aName *C.char, checked C.int) {
	actionName := C.GoString(aName)
	it, ok := gActionIdToMenuItem.Load(actionName)
	if !ok {
		return
	}

	item := it.(*menu.MenuItem)

	item.Checked = int(checked) == 1

	go item.Click(&menu.CallbackData{MenuItem: item})
}

//export handleMenuItemClick
func handleMenuItemClick(aName *C.char) {
	actionName := C.GoString(aName)
	it, ok := gActionIdToMenuItem.Load(actionName)
	if !ok {
		return
	}

	item := it.(*menu.MenuItem)

	go item.Click(&menu.CallbackData{MenuItem: item})
}
