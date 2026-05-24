//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceRequestVtbl struct {
	_IUnknownVtbl
	GetUri     ComProc
	PutUri     ComProc
	GetMethod  ComProc
	PutMethod  ComProc
	GetContent ComProc
	PutContent ComProc
	GetHeaders ComProc
}

type ICoreWebView2WebResourceRequest struct {
	vtbl *_ICoreWebView2WebResourceRequestVtbl
}

func (i *ICoreWebView2WebResourceRequest) AddRef() uintptr {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2WebResourceRequest) Release() uintptr {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2WebResourceRequest) GetMethod() (string, error) {
	// Create *uint16 to hold result
	var _method *uint16
	hr, _, _ := i.vtbl.GetMethod.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_method)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}
	// Get result and cleanup
	uri := windows.UTF16PtrToString(_method)
	windows.CoTaskMemFree(unsafe.Pointer(_method))
	return uri, nil
}

func (i *ICoreWebView2WebResourceRequest) GetUri() (string, error) {
	
	// Create *uint16 to hold result
	var _uri *uint16
	hr, _, _ := i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	} // Get result and cleanup
	uri := windows.UTF16PtrToString(_uri)
	windows.CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

// GetContent returns the body of the request. Returns nil if there's no body. Make sure to call
// Release on the returned IStream after finished using it.
func (i *ICoreWebView2WebResourceRequest) GetContent() (*IStream, error) {
	var stream *IStream
	hr, _, _ := i.vtbl.GetContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&stream)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}
	return stream, nil
}

// GetHeaders returns the mutable HTTP request headers. Make sure to call
// Release on the returned Object after finished using it.
func (i *ICoreWebView2WebResourceRequest) GetHeaders() (*ICoreWebView2HttpRequestHeaders, error) {
	var headers *ICoreWebView2HttpRequestHeaders
	hr, _, _ := i.vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}
	return headers, nil
}
