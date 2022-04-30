package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "Application.h"
#import "WailsContext.h"

#include <stdlib.h>
*/

import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

func (f *Frontend) TrayMenuAdd(trayMenu *menu.TrayMenu) {
	if f.applicationDidFinishLaunching == false {
		f.trayMenusBuffer = append(f.trayMenusBuffer, trayMenu)
		return
	}
	nsTrayMenu := f.mainWindow.TrayMenuAdd(trayMenu)
	f.trayMenus[trayMenu] = nsTrayMenu
}

type NSTrayMenu struct {
	context      unsafe.Pointer
	nsStatusItem unsafe.Pointer // NSStatusItem
}

func (w *Window) TrayMenuAdd(trayMenu *menu.TrayMenu) *NSTrayMenu {
	return NewNSTrayMenu(w.context, trayMenu)
}
