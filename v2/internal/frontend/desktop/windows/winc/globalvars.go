//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"syscall"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

// Private global variables.
var (
	gAppInstance        w32.HINSTANCE
	gControllerRegistry map[w32.HWND]Controller
	gRegisteredClasses  []string
)

// Public global variables.
var (
	GeneralWndprocCallBack = syscall.NewCallback(generalWndProc)
	DefaultFont            *Font
)
