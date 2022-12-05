//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */
package winc

import (
	"runtime"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

var (
	// resource compilation tool assigns app.ico ID of 3
	// rsrc -manifest app.manifest -ico app.ico -o rsrc.syso
	AppIconID = 3
)

func init() {
	runtime.LockOSThread()

	gAppInstance = w32.GetModuleHandle("")
	if gAppInstance == 0 {
		panic("Error occurred in App.Init")
	}

	// Initialize the common controls
	var initCtrls w32.INITCOMMONCONTROLSEX
	initCtrls.DwSize = uint32(unsafe.Sizeof(initCtrls))
	initCtrls.DwICC =
		w32.ICC_LISTVIEW_CLASSES | w32.ICC_PROGRESS_CLASS | w32.ICC_TAB_CLASSES |
			w32.ICC_TREEVIEW_CLASSES | w32.ICC_BAR_CLASSES

	w32.InitCommonControlsEx(&initCtrls)
}

// SetAppIconID sets recource icon ID for the apps windows.
func SetAppIcon(appIconID int) {
	AppIconID = appIconID
}

func GetAppInstance() w32.HINSTANCE {
	return gAppInstance
}

func PreTranslateMessage(msg *w32.MSG) bool {
	// This functions is called by the MessageLoop. It processes the
	// keyboard accelerator keys and calls Controller.PreTranslateMessage for
	// keyboard and mouse events.

	processed := false

	if (msg.Message >= w32.WM_KEYFIRST && msg.Message <= w32.WM_KEYLAST) ||
		(msg.Message >= w32.WM_MOUSEFIRST && msg.Message <= w32.WM_MOUSELAST) {

		if msg.Hwnd != 0 {
			if controller := GetMsgHandler(msg.Hwnd); controller != nil {
				// Search the chain of parents for pretranslated messages.
				for p := controller; p != nil; p = p.Parent() {

					if processed = p.PreTranslateMessage(msg); processed {
						break
					}
				}
			}
		}
	}

	return processed
}

// RunMainLoop processes messages in main application loop.
func RunMainLoop() int {
	m := (*w32.MSG)(unsafe.Pointer(w32.GlobalAlloc(0, uint32(unsafe.Sizeof(w32.MSG{})))))
	defer w32.GlobalFree(w32.HGLOBAL(unsafe.Pointer(m)))

	for w32.GetMessage(m, 0, 0, 0) != 0 {

		if !PreTranslateMessage(m) {
			w32.TranslateMessage(m)
			w32.DispatchMessage(m)
		}
	}

	w32.GdiplusShutdown()
	return int(m.WParam)
}

// PostMessages processes recent messages. Sometimes helpful for instant window refresh.
func PostMessages() {
	m := (*w32.MSG)(unsafe.Pointer(w32.GlobalAlloc(0, uint32(unsafe.Sizeof(w32.MSG{})))))
	defer w32.GlobalFree(w32.HGLOBAL(unsafe.Pointer(m)))

	for i := 0; i < 10; i++ {
		if w32.GetMessage(m, 0, 0, 0) != 0 {
			if !PreTranslateMessage(m) {
				w32.TranslateMessage(m)
				w32.DispatchMessage(m)
			}
		}
	}
}

func Exit() {
	w32.PostQuitMessage(0)
}
