//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// The args of both NavigationStarting and FrameNavigationStarting.
type _ICoreWebView2NavigationStartingEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri             ComProc
	GetIsUserInitiated ComProc
	GetIsRedirected    ComProc
	GetRequestHeaders  ComProc
	GetCancel          ComProc
	PutCancel          ComProc
	GetNavigationId    ComProc
}

type ICoreWebView2NavigationStartingEventArgs struct {
	vtbl *_ICoreWebView2NavigationStartingEventArgsVtbl
}

func (i *ICoreWebView2NavigationStartingEventArgs) GetUri() (string, error) {
	var _uri *uint16
	hr, _, _ := i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}
	uri := windows.UTF16PtrToString(_uri)
	windows.CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

// PutCancel cancels the navigation when cancel is true. NavigationStarting
// takes no deferral, so it must be called before the event handler returns.
func (i *ICoreWebView2NavigationStartingEventArgs) PutCancel(cancel bool) error {
	hr, _, _ := i.vtbl.PutCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(cancel)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}
