//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2DevToolsProtocolEventReceivedEventArgsVtbl struct {
	IUnknownVtbl
	GetParameterObjectAsJson ComProc
}

type ICoreWebView2DevToolsProtocolEventReceivedEventArgs struct {
	Vtbl *ICoreWebView2DevToolsProtocolEventReceivedEventArgsVtbl
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) GetParameterObjectAsJson() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetParameterObjectAsJson.Call(
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
