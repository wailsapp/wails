//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2DOMContentLoadedEventArgsVtbl struct {
	IUnknownVtbl
	GetNavigationId ComProc
}

type ICoreWebView2DOMContentLoadedEventArgs struct {
	Vtbl *ICoreWebView2DOMContentLoadedEventArgsVtbl
}

func (i *ICoreWebView2DOMContentLoadedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2DOMContentLoadedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2DOMContentLoadedEventArgs) GetNavigationId() (uint64, error) {

	var value uint64

	hr, _, err := i.Vtbl.GetNavigationId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}
