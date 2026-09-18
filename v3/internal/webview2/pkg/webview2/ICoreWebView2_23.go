//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_23Vtbl struct {
	IUnknownVtbl
	PostWebMessageAsJsonWithAdditionalObjects ComProc
}

type ICoreWebView2_23 struct {
	Vtbl *ICoreWebView2_23Vtbl
}

func (i *ICoreWebView2_23) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_23() *ICoreWebView2_23 {
	var result *ICoreWebView2_23

	iidICoreWebView2_23 := NewGUID("{508f0db5-90c4-5872-90a7-267a91377502}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_23)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_23) PostWebMessageAsJsonWithAdditionalObjects(webMessageAsJson string, additionalObjects *ICoreWebView2ObjectCollectionView) error {

	// Convert string 'webMessageAsJson' to *uint16
	_webMessageAsJson, err := UTF16PtrFromString(webMessageAsJson)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PostWebMessageAsJsonWithAdditionalObjects.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_webMessageAsJson)),
		uintptr(unsafe.Pointer(additionalObjects)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
