//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2SourceChangedEventArgsVtbl struct {
	IUnknownVtbl
	GetIsNewDocument ComProc
}

type ICoreWebView2SourceChangedEventArgs struct {
	Vtbl *ICoreWebView2SourceChangedEventArgsVtbl
}

func (i *ICoreWebView2SourceChangedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2SourceChangedEventArgs) GetIsNewDocument() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsNewDocument.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}
