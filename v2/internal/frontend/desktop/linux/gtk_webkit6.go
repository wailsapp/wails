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
	item := gActionIdToMenuItem[actionName]

	prevActionId := radioActionName + "::" + prevId
	prevItem := gActionIdToMenuItem[prevActionId]

	prevItem.Checked = false
	item.Checked = true

	go item.Click(&menu.CallbackData{MenuItem: item})
}

//export handleMenuCheckItemClick
func handleMenuCheckItemClick(aName *C.char, checked C.int) {
	actionName := C.GoString(aName)
	item := gActionIdToMenuItem[actionName]

	item.Checked = int(checked) == 1

	go item.Click(&menu.CallbackData{MenuItem: item})
}

//export handleMenuItemClick
func handleMenuItemClick(aName *C.char) {
	actionName := C.GoString(aName)
	item := gActionIdToMenuItem[actionName]

	go item.Click(&menu.CallbackData{MenuItem: item})
}
