//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Profile3Vtbl struct {
	IUnknownVtbl
	GetPreferredTrackingPreventionLevel ComProc
	PutPreferredTrackingPreventionLevel ComProc
}

type ICoreWebView2Profile3 struct {
	Vtbl *ICoreWebView2Profile3Vtbl
}

func (i *ICoreWebView2Profile3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Profile3() *ICoreWebView2Profile3 {
	var result *ICoreWebView2Profile3

	iidICoreWebView2Profile3 := NewGUID("{b188e659-5685-4e05-bdba-fc640e0f1992}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Profile3) GetPreferredTrackingPreventionLevel() (COREWEBVIEW2_TRACKING_PREVENTION_LEVEL, error) {

	var value COREWEBVIEW2_TRACKING_PREVENTION_LEVEL

	hr, _, _ := i.Vtbl.GetPreferredTrackingPreventionLevel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Profile3) PutPreferredTrackingPreventionLevel(value COREWEBVIEW2_TRACKING_PREVENTION_LEVEL) error {

	hr, _, _ := i.Vtbl.PutPreferredTrackingPreventionLevel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
