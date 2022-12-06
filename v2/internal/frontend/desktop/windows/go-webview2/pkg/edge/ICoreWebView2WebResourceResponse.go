//go:build windows

package edge

import "unsafe"

type _ICoreWebView2WebResourceResponseVtbl struct {
	_IUnknownVtbl
	GetContent      ComProc
	PutContent      ComProc
	GetHeaders      ComProc
	GetStatusCode   ComProc
	PutStatusCode   ComProc
	GetReasonPhrase ComProc
	PutReasonPhrase ComProc
}

type ICoreWebView2WebResourceResponse struct {
	vtbl *_ICoreWebView2WebResourceResponseVtbl
}

func (i *ICoreWebView2WebResourceResponse) AddRef() uintptr {
	return i.AddRef()
}

func (i *ICoreWebView2WebResourceResponse) Release() error {
	return i.vtbl.CallRelease(unsafe.Pointer(i))
}
