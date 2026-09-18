//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_2Vtbl struct {
	IUnknownVtbl
	AddWebResourceResponseReceived    ComProc
	RemoveWebResourceResponseReceived ComProc
	NavigateWithWebResourceRequest    ComProc
	AddDOMContentLoaded               ComProc
	RemoveDOMContentLoaded            ComProc
	GetCookieManager                  ComProc
	GetEnvironment                    ComProc
}

type ICoreWebView2_2 struct {
	Vtbl *ICoreWebView2_2Vtbl
}

func (i *ICoreWebView2_2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_2() *ICoreWebView2_2 {
	var result *ICoreWebView2_2

	iidICoreWebView2_2 := NewGUID("{9E8F0CF8-E670-4B5E-B2BC-73E061E3184C}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_2) AddWebResourceResponseReceived(eventHandler *ICoreWebView2WebResourceResponseReceivedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddWebResourceResponseReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_2) RemoveWebResourceResponseReceived(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveWebResourceResponseReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_2) NavigateWithWebResourceRequest(request *ICoreWebView2WebResourceRequest) error {

	hr, _, _ := i.Vtbl.NavigateWithWebResourceRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(request)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_2) AddDOMContentLoaded(eventHandler *ICoreWebView2DOMContentLoadedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddDOMContentLoaded.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_2) RemoveDOMContentLoaded(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveDOMContentLoaded.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_2) GetCookieManager() (*ICoreWebView2CookieManager, error) {

	var cookieManager *ICoreWebView2CookieManager

	hr, _, _ := i.Vtbl.GetCookieManager.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&cookieManager)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return cookieManager, nil
}

func (i *ICoreWebView2_2) GetEnvironment() (*ICoreWebView2Environment, error) {

	var environment *ICoreWebView2Environment

	hr, _, _ := i.Vtbl.GetEnvironment.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&environment)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return environment, nil
}
