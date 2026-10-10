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


func (i *ICoreWebView2EnvironmentOptions4) GetCustomSchemeRegistrations() (uint32, []*ICoreWebView2CustomSchemeRegistration, error) {

	var count uint32
	// COM writes the address of a CoTaskMem-allocated array of *ICoreWebView2CustomSchemeRegistration
	var _schemeRegistrations **ICoreWebView2CustomSchemeRegistration

	hr, _, _ := i.Vtbl.GetCustomSchemeRegistrations.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&_schemeRegistrations)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, nil, syscall.Errno(hr)
	}
	// Copy the COM-allocated array into a Go slice, then free the array itself.
	// The elements are interface pointers whose references are owned by the caller.
	schemeRegistrations := make([]*ICoreWebView2CustomSchemeRegistration, count)
	if _schemeRegistrations != nil {
		copy(schemeRegistrations, unsafe.Slice(_schemeRegistrations, count))
		CoTaskMemFree(unsafe.Pointer(_schemeRegistrations))
	}
	return count, schemeRegistrations, nil
}

func (i *ICoreWebView2EnvironmentOptions4) SetCustomSchemeRegistrations(count uint32, schemeRegistrations []*ICoreWebView2CustomSchemeRegistration) error {
	// Convert []*ICoreWebView2CustomSchemeRegistration 'schemeRegistrations' to a pointer to its first element (T**)
	var _schemeRegistrations **ICoreWebView2CustomSchemeRegistration
	if len(schemeRegistrations) > 0 {
		_schemeRegistrations = &schemeRegistrations[0]
	}


	hr, _, _ := i.Vtbl.SetCustomSchemeRegistrations.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(count),
		uintptr(unsafe.Pointer(_schemeRegistrations)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
