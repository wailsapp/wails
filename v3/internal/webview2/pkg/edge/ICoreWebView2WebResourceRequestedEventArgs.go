//go:build windows

package edge

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetRequest         ComProc
	GetResponse        ComProc
	PutResponse        ComProc
	GetDeferral        ComProc
	GetResourceContext ComProc
}

type ICoreWebView2WebResourceRequestedEventArgs struct {
	vtbl *_ICoreWebView2WebResourceRequestedEventArgsVtbl
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) PutResponse(response *ICoreWebView2WebResourceResponse) error {
	

	hr, _, _ := i.vtbl.PutResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(response)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest, error) {
	
	var request *ICoreWebView2WebResourceRequest
	hr, _, _ := i.vtbl.GetRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&request)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}
	return request, nil
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {
	var deferral *ICoreWebView2Deferral

	hr, _, _ := i.vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}

	if deferral == nil {
		return nil, fmt.Errorf("unknown error")
	}

	return deferral, nil
}
