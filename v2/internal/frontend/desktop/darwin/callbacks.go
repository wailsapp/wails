//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "Application.h"

#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"strconv"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

func (f *Frontend) handleCallback(menuItemID uint) error {
	menuItem := getMenuItemForID(menuItemID)
	if menuItem == nil {
		return errors.New("unknown menuItem ID: " + strconv.Itoa(int(menuItemID)))
	}

	wailsMenuItem := menuItem.wailsMenuItem
	if wailsMenuItem.Type == menu.CheckboxType {
		wailsMenuItem.Checked = !wailsMenuItem.Checked
		C.UpdateMenuItem(menuItem.nsmenuitem, bool2Cint(wailsMenuItem.Checked))
	}
	if wailsMenuItem.Type == menu.RadioType {
		// Ignore if we clicked the item that is already checked
		if !wailsMenuItem.Checked {
			for _, item := range menuItem.radioGroupMembers {
				if item.wailsMenuItem.Checked {
					item.wailsMenuItem.Checked = false
					C.UpdateMenuItem(item.nsmenuitem, C.int(0))
				}
			}
			wailsMenuItem.Checked = true
			C.UpdateMenuItem(menuItem.nsmenuitem, C.int(1))
		}
	}
	if wailsMenuItem.Click != nil {
		go wailsMenuItem.Click(&menu.CallbackData{MenuItem: wailsMenuItem})
	}
	return nil
}
