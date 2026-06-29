//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2WebResourceResponseVtbl struct {
	IUnknownVtbl
	GetContent      ComProc
	PutContent      ComProc
	GetHeaders      ComProc
	GetStatusCode   ComProc
	PutStatusCode   ComProc
	GetReasonPhrase ComProc
	PutReasonPhrase ComProc
}

type ICoreWebView2WebResourceResponse struct {
	Vtbl *ICoreWebView2WebResourceResponseVtbl
}

func (i *ICoreWebView2WebResourceResponse) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2WebResourceResponse) GetContent() (*IStream, error) {

	var content *IStream

	hr, _, _ := i.Vtbl.GetContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&content)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return content, nil
}

func (i *ICoreWebView2WebResourceResponse) PutContent(content *IStream) error {

	hr, _, _ := i.Vtbl.PutContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(content)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2WebResourceResponse) GetHeaders() (*ICoreWebView2HttpResponseHeaders, error) {

	var headers *ICoreWebView2HttpResponseHeaders

	hr, _, _ := i.Vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return headers, nil
}

func (i *ICoreWebView2WebResourceResponse) GetStatusCode() (int, error) {

	var statusCode int

	hr, _, _ := i.Vtbl.GetStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(statusCode),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return statusCode, nil
}

func (i *ICoreWebView2WebResourceResponse) PutStatusCode(statusCode int) error {

	hr, _, _ := i.Vtbl.PutStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(statusCode),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2WebResourceResponse) GetReasonPhrase() (string, error) {
	// Create *uint16 to hold result
	var _reasonPhrase *uint16

	hr, _, _ := i.Vtbl.GetReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_reasonPhrase)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	reasonPhrase := UTF16PtrToString(_reasonPhrase)
	CoTaskMemFree(unsafe.Pointer(_reasonPhrase))
	return reasonPhrase, nil
}

func (i *ICoreWebView2WebResourceResponse) PutReasonPhrase(reasonPhrase string) error {

	// Convert string 'reasonPhrase' to *uint16
	_reasonPhrase, err := UTF16PtrFromString(reasonPhrase)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_reasonPhrase)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
