package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2HttpRequestHeadersVtbl struct {
	_IUnknownVtbl
	GetHeader    ComProc
	GetHeaders   ComProc
	Contains     ComProc
	SetHeader    ComProc
	RemoveHeader ComProc
	GetIterator  ComProc
}

type ICoreWebView2HttpRequestHeaders struct {
	vtbl *_ICoreWebView2HttpRequestHeadersVtbl
}

func (i *ICoreWebView2HttpRequestHeaders) Release() error {
	return i.vtbl.CallRelease(unsafe.Pointer(i))
}

// GetIterator returns an iterator over the collection of request headers. Make sure to call
// Release on the returned Object after finished using it.
func (i *ICoreWebView2HttpRequestHeaders) GetIterator() (*ICoreWebView2HttpHeadersCollectionIterator, error) {
	var headers *ICoreWebView2HttpHeadersCollectionIterator
	res, _, err := i.vtbl.GetIterator.Call(
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
