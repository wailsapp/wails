//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Frame4Vtbl struct {
	IUnknownVtbl
	PostSharedBufferToScript ComProc
}

type ICoreWebView2Frame4 struct {
	Vtbl *ICoreWebView2Frame4Vtbl
}

func (i *ICoreWebView2Frame4) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Frame4() *ICoreWebView2Frame4 {
	var result *ICoreWebView2Frame4

	iidICoreWebView2Frame4 := NewGUID("{188782dc-92aa-4732-ab3c-fcc59f6f68b9}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame4)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Frame4) PostSharedBufferToScript(sharedBuffer *ICoreWebView2SharedBuffer, access COREWEBVIEW2_SHARED_BUFFER_ACCESS, additionalDataAsJson string) error {

	// Convert string 'additionalDataAsJson' to *uint16
	_additionalDataAsJson, err := UTF16PtrFromString(additionalDataAsJson)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PostSharedBufferToScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(sharedBuffer)),
		uintptr(access),
		uintptr(unsafe.Pointer(_additionalDataAsJson)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
