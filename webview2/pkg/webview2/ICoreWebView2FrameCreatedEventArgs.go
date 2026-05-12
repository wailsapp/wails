//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2FrameCreatedEventArgsVtbl struct {
	IUnknownVtbl
	GetFrame ComProc
}

type ICoreWebView2FrameCreatedEventArgs struct {
	Vtbl *ICoreWebView2FrameCreatedEventArgsVtbl
}

func (i *ICoreWebView2FrameCreatedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameCreatedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2FrameCreatedEventArgs) GetFrame() (*ICoreWebView2Frame, error) {

	var value *ICoreWebView2Frame

	hr, _, err := i.Vtbl.GetFrame.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
