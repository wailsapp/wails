//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2DevToolsProtocolEventReceivedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetSessionId ComProc
}

type ICoreWebView2DevToolsProtocolEventReceivedEventArgs2 struct {
	Vtbl *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2Vtbl
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2DevToolsProtocolEventReceivedEventArgs2() (*ICoreWebView2DevToolsProtocolEventReceivedEventArgs2, error) {
	var result *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2

	iidICoreWebView2DevToolsProtocolEventReceivedEventArgs2 := NewGUID("{2dc4959d-1494-4393-95ba-bea4cb9ebd1b}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2DevToolsProtocolEventReceivedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs2) GetSessionId() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetSessionId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}
