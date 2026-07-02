//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2NavigationStartingEventArgsVtbl struct {
	IUnknownVtbl
	GetUri             ComProc
	GetIsUserInitiated ComProc
	GetIsRedirected    ComProc
	GetRequestHeaders  ComProc
	GetCancel          ComProc
	PutCancel          ComProc
	GetNavigationId    ComProc
}

type ICoreWebView2NavigationStartingEventArgs struct {
	Vtbl *ICoreWebView2NavigationStartingEventArgsVtbl
}

func (i *ICoreWebView2NavigationStartingEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2NavigationStartingEventArgs) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16

	hr, _, _ := i.Vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	uri := UTF16PtrToString(_uri)
	CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

func (i *ICoreWebView2NavigationStartingEventArgs) GetIsUserInitiated() (bool, error) {
	// Create int32 to hold bool result
	var _isUserInitiated int32

	hr, _, _ := i.Vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isUserInitiated)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isUserInitiated := _isUserInitiated != 0
	return isUserInitiated, nil
}

func (i *ICoreWebView2NavigationStartingEventArgs) GetIsRedirected() (bool, error) {
	// Create int32 to hold bool result
	var _isRedirected int32

	hr, _, _ := i.Vtbl.GetIsRedirected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isRedirected)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isRedirected := _isRedirected != 0
	return isRedirected, nil
}

func (i *ICoreWebView2NavigationStartingEventArgs) GetRequestHeaders() (*ICoreWebView2HttpRequestHeaders, error) {

	var requestHeaders *ICoreWebView2HttpRequestHeaders

	hr, _, _ := i.Vtbl.GetRequestHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&requestHeaders)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return requestHeaders, nil
}

func (i *ICoreWebView2NavigationStartingEventArgs) GetCancel() (bool, error) {
	// Create int32 to hold bool result
	var _cancel int32

	hr, _, _ := i.Vtbl.GetCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_cancel)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	cancel := _cancel != 0
	return cancel, nil
}

func (i *ICoreWebView2NavigationStartingEventArgs) PutCancel(cancel bool) error {

	hr, _, _ := i.Vtbl.PutCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&cancel)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2NavigationStartingEventArgs) GetNavigationId() (uint64, error) {

	var navigationId uint64

	hr, _, _ := i.Vtbl.GetNavigationId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&navigationId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return navigationId, nil
}
