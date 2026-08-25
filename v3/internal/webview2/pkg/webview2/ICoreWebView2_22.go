//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_22Vtbl struct {
	IUnknownVtbl
	AddWebResourceRequestedFilterWithRequestSourceKinds    ComProc
	RemoveWebResourceRequestedFilterWithRequestSourceKinds ComProc
}

type ICoreWebView2_22 struct {
	Vtbl *ICoreWebView2_22Vtbl
}

func (i *ICoreWebView2_22) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_22() *ICoreWebView2_22 {
	var result *ICoreWebView2_22

	iidICoreWebView2_22 := NewGUID("{db75dfc7-a857-4632-a398-6969dde26c0a}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_22)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_22) AddWebResourceRequestedFilterWithRequestSourceKinds(uri string, ResourceContext COREWEBVIEW2_WEB_RESOURCE_CONTEXT, requestSourceKinds COREWEBVIEW2_WEB_RESOURCE_REQUEST_SOURCE_KINDS) error {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.AddWebResourceRequestedFilterWithRequestSourceKinds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(ResourceContext),
		uintptr(requestSourceKinds),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_22) RemoveWebResourceRequestedFilterWithRequestSourceKinds(uri string, ResourceContext COREWEBVIEW2_WEB_RESOURCE_CONTEXT, requestSourceKinds COREWEBVIEW2_WEB_RESOURCE_REQUEST_SOURCE_KINDS) error {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.RemoveWebResourceRequestedFilterWithRequestSourceKinds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(ResourceContext),
		uintptr(requestSourceKinds),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
