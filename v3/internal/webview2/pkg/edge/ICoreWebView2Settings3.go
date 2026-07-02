//go:build windows

package edge

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Settings3Vtbl struct {
	_IUnknownVtbl
	GetIsScriptEnabled                  ComProc
	PutIsScriptEnabled                  ComProc
	GetIsWebMessageEnabled              ComProc
	PutIsWebMessageEnabled              ComProc
	GetAreDefaultScriptDialogsEnabled   ComProc
	PutAreDefaultScriptDialogsEnabled   ComProc
	GetIsStatusBarEnabled               ComProc
	PutIsStatusBarEnabled               ComProc
	GetAreDevToolsEnabled               ComProc
	PutAreDevToolsEnabled               ComProc
	GetAreDefaultContextMenusEnabled    ComProc
	PutAreDefaultContextMenusEnabled    ComProc
	GetAreHostObjectsAllowed            ComProc
	PutAreHostObjectsAllowed            ComProc
	GetIsZoomControlEnabled             ComProc
	PutIsZoomControlEnabled             ComProc
	GetIsBuiltInErrorPageEnabled        ComProc
	PutIsBuiltInErrorPageEnabled        ComProc
	GetUserAgent                        ComProc
	PutUserAgent                        ComProc
	GetAreBrowserAcceleratorKeysEnabled ComProc
	PutAreBrowserAcceleratorKeysEnabled ComProc
}

type ICoreWebView2Settings3 struct {
	Vtbl *ICoreWebView2Settings3Vtbl
}

func (i *ICoreWebView2Settings3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebViewSettings) GetICoreWebView2Settings3() *ICoreWebView2Settings3 {
	var result *ICoreWebView2Settings3

	iidICoreWebView2Settings3 := NewGUID("{fdb5ab74-af33-4854-84f0-0a631deb5eba}")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Settings3) GetAreBrowserAcceleratorKeysEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _areBrowserAcceleratorKeysEnabled int32

	hr, _, _ := i.Vtbl.GetAreBrowserAcceleratorKeysEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_areBrowserAcceleratorKeysEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	areBrowserAcceleratorKeysEnabled := _areBrowserAcceleratorKeysEnabled != 0
	return areBrowserAcceleratorKeysEnabled, nil
}

func (i *ICoreWebView2Settings3) PutAreBrowserAcceleratorKeysEnabled(areBrowserAcceleratorKeysEnabled bool) error {

	hr, _, _ := i.Vtbl.PutAreBrowserAcceleratorKeysEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&areBrowserAcceleratorKeysEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
