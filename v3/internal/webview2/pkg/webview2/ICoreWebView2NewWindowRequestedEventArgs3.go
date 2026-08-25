//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2NewWindowRequestedEventArgs3Vtbl struct {
	IUnknownVtbl
	GetOriginalSourceFrameInfo ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs3 struct {
	Vtbl *ICoreWebView2NewWindowRequestedEventArgs3Vtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2NewWindowRequestedEventArgs3() *ICoreWebView2NewWindowRequestedEventArgs3 {
	var result *ICoreWebView2NewWindowRequestedEventArgs3

	iidICoreWebView2NewWindowRequestedEventArgs3 := NewGUID("{842bed3c-6ad6-4dd9-b938-28c96667ad66}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NewWindowRequestedEventArgs3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2NewWindowRequestedEventArgs3) GetOriginalSourceFrameInfo() (*ICoreWebView2FrameInfo, error) {

	var value *ICoreWebView2FrameInfo

	hr, _, _ := i.Vtbl.GetOriginalSourceFrameInfo.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
