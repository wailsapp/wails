//go:build windows

package edge

import (
	"syscall"
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
	return i.AddRef()
}

func (i *ICoreWebView2WebResourceRequest) GetMethod() (string, error) {
	// Create *uint16 to hold result
	var _method *uint16
	res, _, err := i.vtbl.GetMethod.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_method)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	if windows.Handle(res) != windows.S_OK {
		return "", syscall.Errno(res)
	}
	// Get result and cleanup
	uri := windows.UTF16PtrToString(_method)
	windows.CoTaskMemFree(unsafe.Pointer(_method))
	return uri, nil
}

func (i *ICoreWebView2WebResourceRequest) GetUri() (string, error) {
	var err error
	// Create *uint16 to hold result
	var _uri *uint16
	_, _, err = i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	} // Get result and cleanup
	uri := windows.UTF16PtrToString(_uri)
	windows.CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

// GetContent returns the body of the request. Returns nil if there's no body. Make sure to call
// Release on the returned IStream after finished using it.
func (i *ICoreWebView2WebResourceRequest) GetContent() (*IStream, error) {
	var stream *IStream
	res, _, err := i.vtbl.GetContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&stream)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	if windows.Handle(res) != windows.S_OK {
		return nil, syscall.Errno(res)
	}
	return stream, nil
}

// GetHeaders returns the mutable HTTP request headers. Make sure to call
// Release on the returned Object after finished using it.
func (i *ICoreWebView2WebResourceRequest) GetHeaders() (*ICoreWebView2HttpRequestHeaders, error) {
	var headers *ICoreWebView2HttpRequestHeaders
	res, _, err := i.vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	if windows.Handle(res) != windows.S_OK {
		return nil, syscall.Errno(res)
	}
	return headers, nil
}

func (i *ICoreWebView2WebResourceRequest) Release() error {
	return i.vtbl.CallRelease(unsafe.Pointer(i))
}
