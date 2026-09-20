//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICoreWebView2Settings9Vtbl struct {
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
	GetIsPasswordAutosaveEnabled        ComProc
	PutIsPasswordAutosaveEnabled        ComProc
	GetIsGeneralAutofillEnabled         ComProc
	PutIsGeneralAutofillEnabled         ComProc
	GetIsPinchZoomEnabled               ComProc
	PutIsPinchZoomEnabled               ComProc
	GetIsSwipeNavigationEnabled         ComProc
	PutIsSwipeNavigationEnabled         ComProc
	GetHiddenPdfToolbarItems            ComProc
	PutHiddenPdfToolbarItems            ComProc
	GetIsReputationCheckingRequired     ComProc
	PutIsReputationCheckingRequired     ComProc
	GetIsNonClientRegionSupportEnabled  ComProc
	PutIsNonClientRegionSupportEnabled  ComProc
}

type ICoreWebView2Settings9 struct {
	Vtbl *ICoreWebView2Settings9Vtbl
}

func (i *ICoreWebView2Settings9) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebViewSettings) GetICoreWebView2Settings9() *ICoreWebView2Settings9 {
	var result *ICoreWebView2Settings9

	iidICoreWebView2Settings9 := NewGUID("{0528a73b-e92d-49f4-927a-e547dddaa37d}")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings9)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Settings9) GetIsNonClientRegionSupportEnabled() (bool, error) {
	var value int32

	hr, _, _ := i.Vtbl.GetIsNonClientRegionSupportEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	return value != 0, nil
}

func (i *ICoreWebView2Settings9) PutIsNonClientRegionSupportEnabled(value bool) error {
	hr, _, _ := i.Vtbl.PutIsNonClientRegionSupportEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
