//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2WebResourceRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetRequest         ComProc
	GetResponse        ComProc
	PutResponse        ComProc
	GetDeferral        ComProc
	GetResourceContext ComProc
}

type ICoreWebView2WebResourceRequestedEventArgs struct {
	Vtbl *ICoreWebView2WebResourceRequestedEventArgsVtbl
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest, error) {

	var request *ICoreWebView2WebResourceRequest

	hr, _, _ := i.Vtbl.GetRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&request)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return request, nil
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) GetResponse() (*ICoreWebView2WebResourceResponse, error) {

	var response *ICoreWebView2WebResourceResponse

	hr, _, _ := i.Vtbl.GetResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&response)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return response, nil
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) PutResponse(response *ICoreWebView2WebResourceResponse) error {

	hr, _, _ := i.Vtbl.PutResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(response)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, _ := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, nil
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) GetResourceContext() (COREWEBVIEW2_WEB_RESOURCE_CONTEXT, error) {

	var context COREWEBVIEW2_WEB_RESOURCE_CONTEXT

	hr, _, _ := i.Vtbl.GetResourceContext.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&context)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return context, nil
}
