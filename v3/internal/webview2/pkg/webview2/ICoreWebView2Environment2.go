//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment2Vtbl struct {
	IUnknownVtbl
	CreateWebResourceRequest ComProc
}

type ICoreWebView2Environment2 struct {
	Vtbl *ICoreWebView2Environment2Vtbl
}

func (i *ICoreWebView2Environment2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment2() *ICoreWebView2Environment2 {
	var result *ICoreWebView2Environment2

	iidICoreWebView2Environment2 := NewGUID("{41f3632b-5ef4-404f-ad82-2d606c5a9a21}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment2) CreateWebResourceRequest(uri string, Method string, postData *IStream, Headers string) (*ICoreWebView2WebResourceRequest, error) {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return nil, err
	}
	// Convert string 'Method' to *uint16
	_Method, err := UTF16PtrFromString(Method)
	if err != nil {
		return nil, err
	}
	// Convert string 'Headers' to *uint16
	_Headers, err := UTF16PtrFromString(Headers)
	if err != nil {
		return nil, err
	}
	var value *ICoreWebView2WebResourceRequest

	hr, _, _ := i.Vtbl.CreateWebResourceRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(unsafe.Pointer(_Method)),
		uintptr(unsafe.Pointer(postData)),
		uintptr(unsafe.Pointer(_Headers)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
