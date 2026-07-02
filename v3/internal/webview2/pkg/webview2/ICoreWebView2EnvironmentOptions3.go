//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2EnvironmentOptions3Vtbl struct {
	IUnknownVtbl
	GetIsCustomCrashReportingEnabled ComProc
	PutIsCustomCrashReportingEnabled ComProc
}

type ICoreWebView2EnvironmentOptions3 struct {
	Vtbl *ICoreWebView2EnvironmentOptions3Vtbl
}

func (i *ICoreWebView2EnvironmentOptions3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2EnvironmentOptions3) GetIsCustomCrashReportingEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsCustomCrashReportingEnabled.Call(
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

func (i *ICoreWebView2EnvironmentOptions3) PutIsCustomCrashReportingEnabled(value bool) error {

	hr, _, _ := i.Vtbl.PutIsCustomCrashReportingEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
