//go:build windows

package edge

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Settings2Vtbl struct {
	_IUnknownVtbl
	GetIsScriptEnabled                ComProc
	PutIsScriptEnabled                ComProc
	GetIsWebMessageEnabled            ComProc
	PutIsWebMessageEnabled            ComProc
	GetAreDefaultScriptDialogsEnabled ComProc
	PutAreDefaultScriptDialogsEnabled ComProc
	GetIsStatusBarEnabled             ComProc
	PutIsStatusBarEnabled             ComProc
	GetAreDevToolsEnabled             ComProc
	PutAreDevToolsEnabled             ComProc
	GetAreDefaultContextMenusEnabled  ComProc
	PutAreDefaultContextMenusEnabled  ComProc
	GetAreHostObjectsAllowed          ComProc
	PutAreHostObjectsAllowed          ComProc
	GetIsZoomControlEnabled           ComProc
	PutIsZoomControlEnabled           ComProc
	GetIsBuiltInErrorPageEnabled      ComProc
	PutIsBuiltInErrorPageEnabled      ComProc
	GetUserAgent                      ComProc
	PutUserAgent                      ComProc
}

type ICoreWebView2Settings2 struct {
	Vtbl *ICoreWebView2Settings2Vtbl
}

func (i *ICoreWebView2Settings2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebViewSettings) GetICoreWebView2Settings2() *ICoreWebView2Settings2 {
	var result *ICoreWebView2Settings2

	iidICoreWebView2Settings2 := NewGUID("{ee9a0f68-f46c-4e32-ac23-ef8cac224d2a}")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Settings2) GetUserAgent() (string, error) {
	// Create *uint16 to hold result
	var _userAgent *uint16

	hr, _, _ := i.Vtbl.GetUserAgent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_userAgent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	userAgent := UTF16PtrToString(_userAgent)
	CoTaskMemFree(unsafe.Pointer(_userAgent))
	return userAgent, nil
}

func (i *ICoreWebView2Settings2) PutUserAgent(userAgent string) error {

	// Convert string 'userAgent' to *uint16
	_userAgent, err := UTF16PtrFromString(userAgent)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutUserAgent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_userAgent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
