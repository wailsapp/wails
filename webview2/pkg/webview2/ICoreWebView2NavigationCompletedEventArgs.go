//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NavigationCompletedEventArgsVtbl struct {
	IUnknownVtbl
	GetIsSuccess ComProc
	GetWebErrorStatus ComProc
	GetNavigationId ComProc
}

type ICoreWebView2NavigationCompletedEventArgs struct {
	Vtbl *ICoreWebView2NavigationCompletedEventArgsVtbl
}

func (i *ICoreWebView2NavigationCompletedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NavigationCompletedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2NavigationCompletedEventArgs) GetIsSuccess() (bool, error) {
	// Create int32 to hold bool result
	var _isSuccess int32

	hr, _, err := i.Vtbl.GetIsSuccess.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isSuccess)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    isSuccess := _isSuccess != 0
	return isSuccess, err
}

func (i *ICoreWebView2NavigationCompletedEventArgs) GetWebErrorStatus() (COREWEBVIEW2_WEB_ERROR_STATUS, error) {

	var webErrorStatus COREWEBVIEW2_WEB_ERROR_STATUS

	hr, _, err := i.Vtbl.GetWebErrorStatus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&webErrorStatus)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return webErrorStatus, err
}

func (i *ICoreWebView2NavigationCompletedEventArgs) GetNavigationId() (uint64, error) {

	var navigationId uint64

	hr, _, err := i.Vtbl.GetNavigationId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&navigationId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return navigationId, err
}
