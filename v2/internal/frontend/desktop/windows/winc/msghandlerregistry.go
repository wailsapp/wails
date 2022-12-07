//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

func RegMsgHandler(controller Controller) {
	gControllerRegistry[controller.Handle()] = controller
}

func UnRegMsgHandler(hwnd w32.HWND) {
	delete(gControllerRegistry, hwnd)
}

func GetMsgHandler(hwnd w32.HWND) Controller {
	if controller, isExists := gControllerRegistry[hwnd]; isExists {
		return controller
	}

	return nil
}
