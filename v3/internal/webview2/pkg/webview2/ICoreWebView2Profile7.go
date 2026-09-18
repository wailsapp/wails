//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Profile7Vtbl struct {
	IUnknownVtbl
	AddBrowserExtension  ComProc
	GetBrowserExtensions ComProc
}

type ICoreWebView2Profile7 struct {
	Vtbl *ICoreWebView2Profile7Vtbl
}

func (i *ICoreWebView2Profile7) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Profile7() *ICoreWebView2Profile7 {
	var result *ICoreWebView2Profile7

	iidICoreWebView2Profile7 := NewGUID("{7b4c7906-a1aa-4cb4-b723-db09f813d541}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile7)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Profile7) AddBrowserExtension(extensionFolderPath string, handler *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler) error {

	// Convert string 'extensionFolderPath' to *uint16
	_extensionFolderPath, err := UTF16PtrFromString(extensionFolderPath)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.AddBrowserExtension.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_extensionFolderPath)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Profile7) GetBrowserExtensions(handler *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler) error {

	hr, _, _ := i.Vtbl.GetBrowserExtensions.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
