//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NonClientRegionChangedEventArgsVtbl struct {
	IUnknownVtbl
	GetRegionKind ComProc
}

type ICoreWebView2NonClientRegionChangedEventArgs struct {
	Vtbl *ICoreWebView2NonClientRegionChangedEventArgsVtbl
}

func (i *ICoreWebView2NonClientRegionChangedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NonClientRegionChangedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2NonClientRegionChangedEventArgs) GetRegionKind() (COREWEBVIEW2_NON_CLIENT_REGION_KIND, error) {

	var value COREWEBVIEW2_NON_CLIENT_REGION_KIND

	hr, _, err := i.Vtbl.GetRegionKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}
