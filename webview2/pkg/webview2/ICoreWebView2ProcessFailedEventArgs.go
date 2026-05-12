//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ProcessFailedEventArgsVtbl struct {
	IUnknownVtbl
	GetProcessFailedKind ComProc
}

type ICoreWebView2ProcessFailedEventArgs struct {
	Vtbl *ICoreWebView2ProcessFailedEventArgsVtbl
}

func (i *ICoreWebView2ProcessFailedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ProcessFailedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2ProcessFailedEventArgs) GetProcessFailedKind() (COREWEBVIEW2_PROCESS_FAILED_KIND, error) {

	var value COREWEBVIEW2_PROCESS_FAILED_KIND

	hr, _, err := i.Vtbl.GetProcessFailedKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}
