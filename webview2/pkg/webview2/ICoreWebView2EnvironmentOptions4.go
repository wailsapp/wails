//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2EnvironmentOptions4Vtbl struct {
	IUnknownVtbl
	GetCustomSchemeRegistrations ComProc
	SetCustomSchemeRegistrations ComProc
}

type ICoreWebView2EnvironmentOptions4 struct {
	Vtbl *ICoreWebView2EnvironmentOptions4Vtbl
}

func (i *ICoreWebView2EnvironmentOptions4) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2EnvironmentOptions4) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2EnvironmentOptions4) GetCustomSchemeRegistrations() (uint32, ICoreWebView2CustomSchemeRegistration, error) {

	var count uint32
	var schemeRegistrations ICoreWebView2CustomSchemeRegistration

	hr, _, _ := i.Vtbl.GetCustomSchemeRegistrations.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&schemeRegistrations)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, ICoreWebView2CustomSchemeRegistration{}, syscall.Errno(hr)
	}
	return count, schemeRegistrations, nil
}

func (i *ICoreWebView2EnvironmentOptions4) SetCustomSchemeRegistrations(count uint32, schemeRegistrations **ICoreWebView2CustomSchemeRegistration) error {


	hr, _, _ := i.Vtbl.SetCustomSchemeRegistrations.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(count),
		uintptr(unsafe.Pointer(&schemeRegistrations)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
