//go:build windows

/*
 * Copyright (C) 2019 Tad Vizbaras. All Rights Reserved.
 * Copyright (C) 2010-2012 The W32 Authors. All Rights Reserved.
 */

package w32

import (
	"syscall"
	"unsafe"
)

// LISTVIEW parts
const (
	LVP_LISTITEM         = 1
	LVP_LISTGROUP        = 2
	LVP_LISTDETAIL       = 3
	LVP_LISTSORTEDDETAIL = 4
	LVP_EMPTYTEXT        = 5
	LVP_GROUPHEADER      = 6
	LVP_GROUPHEADERLINE  = 7
	LVP_EXPANDBUTTON     = 8
	LVP_COLLAPSEBUTTON   = 9
	LVP_COLUMNDETAIL     = 10
)

// LVP_LISTITEM states
const (
	LISS_NORMAL           = 1
	LISS_HOT              = 2
	LISS_SELECTED         = 3
	LISS_DISABLED         = 4
	LISS_SELECTEDNOTFOCUS = 5
	LISS_HOTSELECTED      = 6
)

// TREEVIEW parts
const (
	TVP_TREEITEM = 1
	TVP_GLYPH    = 2
	TVP_BRANCH   = 3
	TVP_HOTGLYPH = 4
)

// TVP_TREEITEM states
const (
	TREIS_NORMAL           = 1
	TREIS_HOT              = 2
	TREIS_SELECTED         = 3
	TREIS_DISABLED         = 4
	TREIS_SELECTEDNOTFOCUS = 5
	TREIS_HOTSELECTED      = 6
)

type HTHEME HANDLE

var (
	// Library
	libuxtheme uintptr

	// Functions
	closeThemeData      uintptr
	drawThemeBackground uintptr
	drawThemeText       uintptr
	getThemeTextExtent  uintptr
	openThemeData       uintptr
	setWindowTheme      uintptr
)

func init() {
	// Library
	libuxtheme = MustLoadLibrary("uxtheme.dll")

	// Functions
	closeThemeData = MustGetProcAddress(libuxtheme, "CloseThemeData")
	drawThemeBackground = MustGetProcAddress(libuxtheme, "DrawThemeBackground")
	drawThemeText = MustGetProcAddress(libuxtheme, "DrawThemeText")
	getThemeTextExtent = MustGetProcAddress(libuxtheme, "GetThemeTextExtent")
	openThemeData = MustGetProcAddress(libuxtheme, "OpenThemeData")
	setWindowTheme = MustGetProcAddress(libuxtheme, "SetWindowTheme")
}

func CloseThemeData(hTheme HTHEME) HRESULT {
	ret, _, _ := syscall.Syscall(closeThemeData, 1,
		uintptr(hTheme),
		0,
		0)

	return HRESULT(ret)
}

func DrawThemeBackground(hTheme HTHEME, hdc HDC, iPartId, iStateId int32, pRect, pClipRect *RECT) HRESULT {
	ret, _, _ := syscall.Syscall6(drawThemeBackground, 6,
		uintptr(hTheme),
		uintptr(hdc),
		uintptr(iPartId),
		uintptr(iStateId),
		uintptr(unsafe.Pointer(pRect)),
		uintptr(unsafe.Pointer(pClipRect)))

	return HRESULT(ret)
}

func DrawThemeText(hTheme HTHEME, hdc HDC, iPartId, iStateId int32, pszText *uint16, iCharCount int32, dwTextFlags, dwTextFlags2 uint32, pRect *RECT) HRESULT {
	ret, _, _ := syscall.Syscall9(drawThemeText, 9,
		uintptr(hTheme),
		uintptr(hdc),
		uintptr(iPartId),
		uintptr(iStateId),
		uintptr(unsafe.Pointer(pszText)),
		uintptr(iCharCount),
		uintptr(dwTextFlags),
		uintptr(dwTextFlags2),
		uintptr(unsafe.Pointer(pRect)))

	return HRESULT(ret)
}

func GetThemeTextExtent(hTheme HTHEME, hdc HDC, iPartId, iStateId int32, pszText *uint16, iCharCount int32, dwTextFlags uint32, pBoundingRect, pExtentRect *RECT) HRESULT {
	ret, _, _ := syscall.Syscall9(getThemeTextExtent, 9,
		uintptr(hTheme),
		uintptr(hdc),
		uintptr(iPartId),
		uintptr(iStateId),
		uintptr(unsafe.Pointer(pszText)),
		uintptr(iCharCount),
		uintptr(dwTextFlags),
		uintptr(unsafe.Pointer(pBoundingRect)),
		uintptr(unsafe.Pointer(pExtentRect)))

	return HRESULT(ret)
}

func OpenThemeData(hwnd HWND, pszClassList *uint16) HTHEME {
	ret, _, _ := syscall.Syscall(openThemeData, 2,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pszClassList)),
		0)

	return HTHEME(ret)
}

func SetWindowTheme(hwnd HWND, pszSubAppName, pszSubIdList *uint16) HRESULT {
	ret, _, _ := syscall.Syscall(setWindowTheme, 3,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pszSubAppName)),
		uintptr(unsafe.Pointer(pszSubIdList)))

	return HRESULT(ret)
}
