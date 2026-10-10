//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_22Vtbl struct {
	ICoreWebView2_21Vtbl
	AddWebResourceRequestedFilterWithRequestSourceKinds ComProc
	RemoveWebResourceRequestedFilterWithRequestSourceKinds ComProc
}

type ICoreWebView2_22 struct {
	Vtbl *ICoreWebView2_22Vtbl
}

func (i *ICoreWebView2_22) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_22) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_22 queries the object for its ICoreWebView2_22 interface. The receiver
// is the root of ICoreWebView2_22's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_22() (*ICoreWebView2_22, error) {
	var result *ICoreWebView2_22

	iidICoreWebView2_22 := NewGUID("{db75dfc7-a857-4632-a398-6969dde26c0a}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_22)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
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
