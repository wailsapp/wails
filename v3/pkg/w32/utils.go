//go:build windows

/*
 * Copyright (C) 2019 Tad Vizbaras. All Rights Reserved.
 * Copyright (C) 2010-2012 The W32 Authors. All Rights Reserved.
 */

package w32

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func MustLoadLibrary(name string) uintptr {
	lib, err := syscall.LoadLibrary(name)
	if err != nil {
		panic(err)
	}

	return uintptr(lib)
}

func MustGetProcAddress(lib uintptr, name string) uintptr {
	addr, err := syscall.GetProcAddress(syscall.Handle(lib), name)
	if err != nil {
		panic(err)
	}

	return uintptr(addr)
}

func SUCCEEDED(hr HRESULT) bool {
	return hr >= 0
}

func FAILED(hr HRESULT) bool {
	return hr < 0
}

func LOWORD(dw uint32) uint16 {
	return uint16(dw)
}

func HIWORD(dw uint32) uint16 {
	return uint16(dw >> 16 & 0xffff)
}

func MAKELONG(lo, hi uint16) uint32 {
	return uint32(uint32(lo) | ((uint32(hi)) << 16))
}

func BoolToBOOL(value bool) BOOL {
	if value {
		return 1
	}

	return 0
}

func UTF16PtrToString(cstr *uint16) string {
	if cstr != nil {
		us := make([]uint16, 0, 256)
		for p := uintptr(unsafe.Pointer(cstr)); ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				return string(utf16.Decode(us))
			}
			us = append(us, u)
		}
	}

	return ""
}

func ComAddRef(unknown *IUnknown) int32 {
	ret, _, _ := syscall.SyscallN(uintptr(unknown.Vtbl.AddRef),
		uintptr(unsafe.Pointer(unknown)),
		0,
		0)
	return int32(ret)
}

func ComRelease(unknown *IUnknown) int32 {
	ret, _, _ := syscall.SyscallN(uintptr(unknown.Vtbl.Release),
		uintptr(unsafe.Pointer(unknown)),
		0,
		0)
	return int32(ret)
}

func ComQueryInterface(unknown *IUnknown, id *GUID) *IDispatch {
	var disp *IDispatch
	hr, _, _ := syscall.SyscallN(uintptr(unknown.Vtbl.QueryInterface),
		uintptr(unsafe.Pointer(unknown)),
		uintptr(unsafe.Pointer(id)),
		uintptr(unsafe.Pointer(&disp)))
	if hr != 0 {
		panic("Invoke QieryInterface error.")
	}
	return disp
}

func ComGetIDsOfName(disp *IDispatch, names []string) []int32 {
	wnames := make([]*uint16, len(names))
	dispid := make([]int32, len(names))
	for i := 0; i < len(names); i++ {
		wnames[i] = syscall.StringToUTF16Ptr(names[i])
	}
	hr, _, _ := syscall.SyscallN(disp.lpVtbl.pGetIDsOfNames,
		uintptr(unsafe.Pointer(disp)),
		uintptr(unsafe.Pointer(IID_NULL)),
		uintptr(unsafe.Pointer(&wnames[0])),
		uintptr(len(names)),
		uintptr(GetUserDefaultLCID()),
		uintptr(unsafe.Pointer(&dispid[0])))
	if hr != 0 {
		panic("Invoke GetIDsOfName error.")
	}
	return dispid
}

func ComInvoke(disp *IDispatch, dispid int32, dispatch int16, params ...interface{}) (result *VARIANT) {
	var dispparams DISPPARAMS

	if dispatch&DISPATCH_PROPERTYPUT != 0 {
		dispnames := [1]int32{DISPID_PROPERTYPUT}
		dispparams.RgdispidNamedArgs = uintptr(unsafe.Pointer(&dispnames[0]))
		dispparams.CNamedArgs = 1
	}
	var vargs []VARIANT
	if len(params) > 0 {
		vargs = make([]VARIANT, len(params))
		for i, v := range params {
			//n := len(params)-i-1
			n := len(params) - i - 1
			VariantInit(&vargs[n])
			switch v.(type) {
			case bool:
				if v.(bool) {
					vargs[n] = VARIANT{VT_BOOL, 0, 0, 0, 0xffff}
				} else {
					vargs[n] = VARIANT{VT_BOOL, 0, 0, 0, 0}
				}
			case *bool:
				vargs[n] = VARIANT{VT_BOOL | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*bool))))}
			case byte:
				vargs[n] = VARIANT{VT_I1, 0, 0, 0, int64(v.(byte))}
			case *byte:
				vargs[n] = VARIANT{VT_I1 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*byte))))}
			case int16:
				vargs[n] = VARIANT{VT_I2, 0, 0, 0, int64(v.(int16))}
			case *int16:
				vargs[n] = VARIANT{VT_I2 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*int16))))}
			case uint16:
				vargs[n] = VARIANT{VT_UI2, 0, 0, 0, int64(v.(int16))}
			case *uint16:
				vargs[n] = VARIANT{VT_UI2 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*uint16))))}
			case int, int32:
				vargs[n] = VARIANT{VT_UI4, 0, 0, 0, int64(v.(int))}
			case *int, *int32:
				vargs[n] = VARIANT{VT_I4 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*int))))}
			case uint, uint32:
				vargs[n] = VARIANT{VT_UI4, 0, 0, 0, int64(v.(uint))}
			case *uint, *uint32:
				vargs[n] = VARIANT{VT_UI4 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*uint))))}
			case int64:
				vargs[n] = VARIANT{VT_I8, 0, 0, 0, v.(int64)}
			case *int64:
				vargs[n] = VARIANT{VT_I8 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*int64))))}
			case uint64:
				vargs[n] = VARIANT{VT_UI8, 0, 0, 0, int64(v.(uint64))}
			case *uint64:
				vargs[n] = VARIANT{VT_UI8 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*uint64))))}
			case float32:
				vargs[n] = VARIANT{VT_R4, 0, 0, 0, int64(v.(float32))}
			case *float32:
				vargs[n] = VARIANT{VT_R4 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*float32))))}
			case float64:
				vargs[n] = VARIANT{VT_R8, 0, 0, 0, int64(v.(float64))}
			case *float64:
				vargs[n] = VARIANT{VT_R8 | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*float64))))}
			case string:
				vargs[n] = VARIANT{VT_BSTR, 0, 0, 0, int64(uintptr(unsafe.Pointer(SysAllocString(v.(string)))))}
			case *string:
				vargs[n] = VARIANT{VT_BSTR | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*string))))}
			case *IDispatch:
				vargs[n] = VARIANT{VT_DISPATCH, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*IDispatch))))}
			case **IDispatch:
				vargs[n] = VARIANT{VT_DISPATCH | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(**IDispatch))))}
			case nil:
				vargs[n] = VARIANT{VT_NULL, 0, 0, 0, 0}
			case *VARIANT:
				vargs[n] = VARIANT{VT_VARIANT | VT_BYREF, 0, 0, 0, int64(uintptr(unsafe.Pointer(v.(*VARIANT))))}
			default:
				panic("unknown type")
			}
		}
		dispparams.Rgvarg = uintptr(unsafe.Pointer(&vargs[0]))
		dispparams.CArgs = uint32(len(params))
	}

	var ret VARIANT
	var excepInfo EXCEPINFO
	VariantInit(&ret)
	hr, _, _ := syscall.SyscallN(disp.lpVtbl.pInvoke,
		uintptr(unsafe.Pointer(disp)),
		uintptr(dispid),
		uintptr(unsafe.Pointer(IID_NULL)),
		uintptr(GetUserDefaultLCID()),
		uintptr(dispatch),
		uintptr(unsafe.Pointer(&dispparams)),
		uintptr(unsafe.Pointer(&ret)),
		uintptr(unsafe.Pointer(&excepInfo)),
		0)
	if hr != 0 {
		if excepInfo.BstrDescription != nil {
			bs := UTF16PtrToString(excepInfo.BstrDescription)
			panic(bs)
		}
	}
	for _, varg := range vargs {
		if varg.VT == VT_BSTR && varg.Val != 0 {
			SysFreeString(((*int16)(unsafe.Pointer(uintptr(varg.Val)))))
		}
	}
	result = &ret
	return
}

func WMMessageToString(msg uintptr) string {
	// Convert windows message to string
	switch msg {
	case WM_APP:
		return "WM_APP"
	case WM_ACTIVATE:
		return "WM_ACTIVATE"
	case WM_ACTIVATEAPP:
		return "WM_ACTIVATEAPP"
	case WM_AFXFIRST:
		return "WM_AFXFIRST"
	case WM_AFXLAST:
		return "WM_AFXLAST"
	case WM_ASKCBFORMATNAME:
		return "WM_ASKCBFORMATNAME"
	case WM_CANCELJOURNAL:
		return "WM_CANCELJOURNAL"
	case WM_CANCELMODE:
		return "WM_CANCELMODE"
	case WM_CAPTURECHANGED:
		return "WM_CAPTURECHANGED"
	case WM_CHANGECBCHAIN:
		return "WM_CHANGECBCHAIN"
	case WM_CHAR:
		return "WM_CHAR"
	case WM_CHARTOITEM:
		return "WM_CHARTOITEM"
	case WM_CHILDACTIVATE:
		return "WM_CHILDACTIVATE"
	case WM_CLEAR:
		return "WM_CLEAR"
	case WM_CLOSE:
		return "WM_CLOSE"
	case WM_COMMAND:
		return "WM_COMMAND"
	case WM_COMMNOTIFY /* OBSOLETE */ :
		return "WM_COMMNOTIFY"
	case WM_COMPACTING:
		return "WM_COMPACTING"
	case WM_COMPAREITEM:
		return "WM_COMPAREITEM"
	case WM_CONTEXTMENU:
		return "WM_CONTEXTMENU"
	case WM_COPY:
		return "WM_COPY"
	case WM_COPYDATA:
		return "WM_COPYDATA"
	case WM_CREATE:
		return "WM_CREATE"
	case WM_CTLCOLORBTN:
		return "WM_CTLCOLORBTN"
	case WM_CTLCOLORDLG:
		return "WM_CTLCOLORDLG"
	case WM_CTLCOLOREDIT:
		return "WM_CTLCOLOREDIT"
	case WM_CTLCOLORLISTBOX:
		return "WM_CTLCOLORLISTBOX"
	case WM_CTLCOLORMSGBOX:
		return "WM_CTLCOLORMSGBOX"
	case WM_CTLCOLORSCROLLBAR:
		return "WM_CTLCOLORSCROLLBAR"
	case WM_CTLCOLORSTATIC:
		return "WM_CTLCOLORSTATIC"
	case WM_CUT:
		return "WM_CUT"
	case WM_DEADCHAR:
		return "WM_DEADCHAR"
	case WM_DELETEITEM:
		return "WM_DELETEITEM"
	case WM_DESTROY:
		return "WM_DESTROY"
	case WM_DESTROYCLIPBOARD:
		return "WM_DESTROYCLIPBOARD"
	case WM_DEVICECHANGE:
		return "WM_DEVICECHANGE"
	case WM_DEVMODECHANGE:
		return "WM_DEVMODECHANGE"
	case WM_DISPLAYCHANGE:
		return "WM_DISPLAYCHANGE"
	case WM_DRAWCLIPBOARD:
		return "WM_DRAWCLIPBOARD"
	case WM_DRAWITEM:
		return "WM_DRAWITEM"
	case WM_DROPFILES:
		return "WM_DROPFILES"
	case WM_ENABLE:
		return "WM_ENABLE"
	case WM_ENDSESSION:
		return "WM_ENDSESSION"
	case WM_ENTERIDLE:
		return "WM_ENTERIDLE"
	case WM_ENTERMENULOOP:
		return "WM_ENTERMENULOOP"
	case WM_ENTERSIZEMOVE:
		return "WM_ENTERSIZEMOVE"
	case WM_ERASEBKGND:
		return "WM_ERASEBKGND"
	case WM_EXITMENULOOP:
		return "WM_EXITMENULOOP"
	case WM_EXITSIZEMOVE:
		return "WM_EXITSIZEMOVE"
	case WM_FONTCHANGE:
		return "WM_FONTCHANGE"
	case WM_GETDLGCODE:
		return "WM_GETDLGCODE"
	case WM_GETFONT:
		return "WM_GETFONT"
	case WM_GETHOTKEY:
		return "WM_GETHOTKEY"
	case WM_GETICON:
		return "WM_GETICON"
	case WM_GETMINMAXINFO:
		return "WM_GETMINMAXINFO"
	case WM_GETTEXT:
		return "WM_GETTEXT"
	case WM_GETTEXTLENGTH:
		return "WM_GETTEXTLENGTH"
	case WM_HANDHELDFIRST:
		return "WM_HANDHELDFIRST"
	case WM_HANDHELDLAST:
		return "WM_HANDHELDLAST"
	case WM_HELP:
		return "WM_HELP"
	case WM_HOTKEY:
		return "WM_HOTKEY"
	case WM_HSCROLL:
		return "WM_HSCROLL"
	case WM_HSCROLLCLIPBOARD:
		return "WM_HSCROLLCLIPBOARD"
	case WM_ICONERASEBKGND:
		return "WM_ICONERASEBKGND"
	case WM_INITDIALOG:
		return "WM_INITDIALOG"
	case WM_INITMENU:
		return "WM_INITMENU"
	case WM_INITMENUPOPUP:
		return "WM_INITMENUPOPUP"
	case WM_INPUT:
		return "WM_INPUT"
	case WM_INPUTLANGCHANGE:
		return "WM_INPUTLANGCHANGE"
	case WM_INPUTLANGCHANGEREQUEST:
		return "WM_INPUTLANGCHANGEREQUEST"
	case WM_KEYDOWN:
		return "WM_KEYDOWN"
	case WM_KEYUP:
		return "WM_KEYUP"
	case WM_KILLFOCUS:
		return "WM_KILLFOCUS"
	case WM_MDIACTIVATE:
		return "WM_MDIACTIVATE"
	case WM_MDICASCADE:
		return "WM_MDICASCADE"
	case WM_MDICREATE:
		return "WM_MDICREATE"
	case WM_MDIDESTROY:
		return "WM_MDIDESTROY"
	case WM_MDIGETACTIVE:
		return "WM_MDIGETACTIVE"
	case WM_MDIICONARRANGE:
		return "WM_MDIICONARRANGE"
	case WM_MDIMAXIMIZE:
		return "WM_MDIMAXIMIZE"
	case WM_MDINEXT:
		return "WM_MDINEXT"
	case WM_MDIREFRESHMENU:
		return "WM_MDIREFRESHMENU"
	case WM_MDIRESTORE:
		return "WM_MDIRESTORE"
	case WM_MDISETMENU:
		return "WM_MDISETMENU"
	case WM_MDITILE:
		return "WM_MDITILE"
	case WM_MEASUREITEM:
		return "WM_MEASUREITEM"
	case WM_GETOBJECT:
		return "WM_GETOBJECT"
	case WM_CHANGEUISTATE:
		return "WM_CHANGEUISTATE"
	case WM_UPDATEUISTATE:
		return "WM_UPDATEUISTATE"
	case WM_QUERYUISTATE:
		return "WM_QUERYUISTATE"
	case WM_UNINITMENUPOPUP:
		return "WM_UNINITMENUPOPUP"
	case WM_MENURBUTTONUP:
		return "WM_MENURBUTTONUP"
	case WM_MENUCOMMAND:
		return "WM_MENUCOMMAND"
	case WM_MENUGETOBJECT:
		return "WM_MENUGETOBJECT"
	case WM_MENUDRAG:
		return "WM_MENUDRAG"
	case WM_APPCOMMAND:
		return "WM_APPCOMMAND"
	case WM_MENUCHAR:
		return "WM_MENUCHAR"
	case WM_MENUSELECT:
		return "WM_MENUSELECT"
	case WM_MOVE:
		return "WM_MOVE"
	case WM_MOVING:
		return "WM_MOVING"
	case WM_NCACTIVATE:
		return "WM_NCACTIVATE"
	case WM_NCCALCSIZE:
		return "WM_NCCALCSIZE"
	case WM_NCCREATE:
		return "WM_NCCREATE"
	case WM_NCDESTROY:
		return "WM_NCDESTROY"
	case WM_NCHITTEST:
		return "WM_NCHITTEST"
	case WM_NCLBUTTONDBLCLK:
		return "WM_NCLBUTTONDBLCLK"
	case WM_NCLBUTTONDOWN:
		return "WM_NCLBUTTONDOWN"
	case WM_NCLBUTTONUP:
		return "WM_NCLBUTTONUP"
	case WM_NCMBUTTONDBLCLK:
		return "WM_NCMBUTTONDBLCLK"
	case WM_NCMBUTTONDOWN:
		return "WM_NCMBUTTONDOWN"
	case WM_NCMBUTTONUP:
		return "WM_NCMBUTTONUP"
	case WM_NCXBUTTONDOWN:
		return "WM_NCXBUTTONDOWN"
	case WM_NCXBUTTONUP:
		return "WM_NCXBUTTONUP"
	case WM_NCXBUTTONDBLCLK:
		return "WM_NCXBUTTONDBLCLK"
	case WM_NCMOUSEHOVER:
		return "WM_NCMOUSEHOVER"
	case WM_NCMOUSELEAVE:
		return "WM_NCMOUSELEAVE"
	case WM_NCMOUSEMOVE:
		return "WM_NCMOUSEMOVE"
	case WM_NCPAINT:
		return "WM_NCPAINT"
	case WM_NCRBUTTONDBLCLK:
		return "WM_NCRBUTTONDBLCLK"
	case WM_NCRBUTTONDOWN:
		return "WM_NCRBUTTONDOWN"
	case WM_NCRBUTTONUP:
		return "WM_NCRBUTTONUP"
	case WM_NEXTDLGCTL:
		return "WM_NEXTDLGCTL"
	case WM_NEXTMENU:
		return "WM_NEXTMENU"
	case WM_NOTIFY:
		return "WM_NOTIFY"
	case WM_NOTIFYFORMAT:
		return "WM_NOTIFYFORMAT"
	case WM_NULL:
		return "WM_NULL"
	case WM_PAINT:
		return "WM_PAINT"
	case WM_PAINTCLIPBOARD:
		return "WM_PAINTCLIPBOARD"
	case WM_PAINTICON:
		return "WM_PAINTICON"
	case WM_PALETTECHANGED:
		return "WM_PALETTECHANGED"
	case WM_PALETTEISCHANGING:
		return "WM_PALETTEISCHANGING"
	case WM_PARENTNOTIFY:
		return "WM_PARENTNOTIFY"
	case WM_PASTE:
		return "WM_PASTE"
	case WM_PENWINFIRST:
		return "WM_PENWINFIRST"
	case WM_PENWINLAST:
		return "WM_PENWINLAST"
	case WM_POWER:
		return "WM_POWER"
	case WM_PRINT:
		return "WM_PRINT"
	case WM_PRINTCLIENT:
		return "WM_PRINTCLIENT"
	case WM_QUERYDRAGICON:
		return "WM_QUERYDRAGICON"
	case WM_QUERYENDSESSION:
		return "WM_QUERYENDSESSION"
	case WM_QUERYNEWPALETTE:
		return "WM_QUERYNEWPALETTE"
	case WM_QUERYOPEN:
		return "WM_QUERYOPEN"
	case WM_QUEUESYNC:
		return "WM_QUEUESYNC"
	case WM_QUIT:
		return "WM_QUIT"
	case WM_RENDERALLFORMATS:
		return "WM_RENDERALLFORMATS"
	case WM_RENDERFORMAT:
		return "WM_RENDERFORMAT"
	case WM_SETCURSOR:
		return "WM_SETCURSOR"
	case WM_SETFOCUS:
		return "WM_SETFOCUS"
	case WM_SETFONT:
		return "WM_SETFONT"
	case WM_SETHOTKEY:
		return "WM_SETHOTKEY"
	case WM_SETICON:
		return "WM_SETICON"
	case WM_SETREDRAW:
		return "WM_SETREDRAW"
	case WM_SETTEXT:
		return "WM_SETTEXT"
	case WM_SETTINGCHANGE:
		return "WM_SETTINGCHANGE"
	case WM_SHOWWINDOW:
		return "WM_SHOWWINDOW"
	case WM_SIZE:
		return "WM_SIZE"
	case WM_SIZECLIPBOARD:
		return "WM_SIZECLIPBOARD"
	case WM_SIZING:
		return "WM_SIZING"
	case WM_SPOOLERSTATUS:
		return "WM_SPOOLERSTATUS"
	case WM_STYLECHANGED:
		return "WM_STYLECHANGED"
	case WM_STYLECHANGING:
		return "WM_STYLECHANGING"
	case WM_SYSCHAR:
		return "WM_SYSCHAR"
	case WM_SYSCOLORCHANGE:
		return "WM_SYSCOLORCHANGE"
	case WM_SYSCOMMAND:
		return "WM_SYSCOMMAND"
	case WM_SYSDEADCHAR:
		return "WM_SYSDEADCHAR"
	case WM_SYSKEYDOWN:
		return "WM_SYSKEYDOWN"
	case WM_SYSKEYUP:
		return "WM_SYSKEYUP"
	case WM_TCARD:
		return "WM_TCARD"
	case WM_THEMECHANGED:
		return "WM_THEMECHANGED"
	case WM_TIMECHANGE:
		return "WM_TIMECHANGE"
	case WM_TIMER:
		return "WM_TIMER"
	case WM_UNDO:
		return "WM_UNDO"
	case WM_USER:
		return "WM_USER"
	case WM_USERCHANGED:
		return "WM_USERCHANGED"
	case WM_VKEYTOITEM:
		return "WM_VKEYTOITEM"
	case WM_VSCROLL:
		return "WM_VSCROLL"
	case WM_VSCROLLCLIPBOARD:
		return "WM_VSCROLLCLIPBOARD"
	case WM_WINDOWPOSCHANGED:
		return "WM_WINDOWPOSCHANGED"
	case WM_WINDOWPOSCHANGING:
		return "WM_WINDOWPOSCHANGING"
	case WM_KEYLAST:
		return "WM_KEYLAST"
	case WM_SYNCPAINT:
		return "WM_SYNCPAINT"
	case WM_MOUSEACTIVATE:
		return "WM_MOUSEACTIVATE"
	case WM_MOUSEMOVE:
		return "WM_MOUSEMOVE"
	case WM_LBUTTONDOWN:
		return "WM_LBUTTONDOWN"
	case WM_LBUTTONUP:
		return "WM_LBUTTONUP"
	case WM_LBUTTONDBLCLK:
		return "WM_LBUTTONDBLCLK"
	case WM_RBUTTONDOWN:
		return "WM_RBUTTONDOWN"
	case WM_RBUTTONUP:
		return "WM_RBUTTONUP"
	case WM_RBUTTONDBLCLK:
		return "WM_RBUTTONDBLCLK"
	case WM_MBUTTONDOWN:
		return "WM_MBUTTONDOWN"
	case WM_MBUTTONUP:
		return "WM_MBUTTONUP"
	case WM_MBUTTONDBLCLK:
		return "WM_MBUTTONDBLCLK"
	case WM_MOUSEWHEEL:
		return "WM_MOUSEWHEEL"
	case WM_XBUTTONDOWN:
		return "WM_XBUTTONDOWN"
	case WM_XBUTTONUP:
		return "WM_XBUTTONUP"
	case WM_MOUSELAST:
		return "WM_MOUSELAST"
	case WM_MOUSEHOVER:
		return "WM_MOUSEHOVER"
	case WM_MOUSELEAVE:
		return "WM_MOUSELEAVE"
	case WM_CLIPBOARDUPDATE:
		return "WM_CLIPBOARDUPDATE"
	default:
		return fmt.Sprintf("0x%08x", msg)
	}
}
