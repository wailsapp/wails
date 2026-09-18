//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2WebMessageReceivedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetAdditionalObjects ComProc
}

type ICoreWebView2WebMessageReceivedEventArgs2 struct {
	Vtbl *ICoreWebView2WebMessageReceivedEventArgs2Vtbl
}

func (i *ICoreWebView2WebMessageReceivedEventArgs2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2WebMessageReceivedEventArgs2() *ICoreWebView2WebMessageReceivedEventArgs2 {
	var result *ICoreWebView2WebMessageReceivedEventArgs2

	iidICoreWebView2WebMessageReceivedEventArgs2 := NewGUID("{06fc7ab7-c90c-4297-9389-33ca01cf6d5e}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2WebMessageReceivedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2WebMessageReceivedEventArgs2) GetAdditionalObjects() (*ICoreWebView2ObjectCollectionView, error) {

	var value *ICoreWebView2ObjectCollectionView

	hr, _, _ := i.Vtbl.GetAdditionalObjects.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
