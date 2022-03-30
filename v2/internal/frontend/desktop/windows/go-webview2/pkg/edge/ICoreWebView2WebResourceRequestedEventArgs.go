package edge

import (
	"golang.org/x/sys/windows"
	"unsafe"
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

func (i *ICoreWebView2WebResourceRequestedEventArgs) AddRef() uintptr {
	return i.AddRef()
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) PutResponse(response *ICoreWebView2WebResourceResponse) error {
	var err error

	_, _, err = i.vtbl.PutResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(response)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest, error) {
	var err error
	var request *ICoreWebView2WebResourceRequest
	_, _, err = i.vtbl.GetRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&request)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return request, nil
}
