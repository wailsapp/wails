//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2DevToolsProtocolEventReceivedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetSessionId ComProc
}

type ICoreWebView2DevToolsProtocolEventReceivedEventArgs2 struct {
	Vtbl *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2Vtbl
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2DevToolsProtocolEventReceivedEventArgs2() *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2 {
	var result *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2

	iidICoreWebView2DevToolsProtocolEventReceivedEventArgs2 := NewGUID("{2dc4959d-1494-4393-95ba-bea4cb9ebd1b}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2DevToolsProtocolEventReceivedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2) GetSessionId() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetSessionId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}
