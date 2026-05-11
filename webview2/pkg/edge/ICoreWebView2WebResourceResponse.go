//go:build windows

package edge

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/webview2/internal/w32"
	"golang.org/x/sys/windows"
)

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

func (i *ICoreWebView2WebResourceResponse) AddRef() error {
	_, _, err := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	if err != nil && !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}

	return nil
}

func (i *ICoreWebView2WebResourceResponse) Release() error {
	_, _, err := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	if err != nil && !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}

	return nil
}

// GetHeaders returns the mutable HTTP request headers. Make sure to call
// Release on the returned Object after finished using it.
func (i *ICoreWebView2WebResourceResponse) GetHeaders() (*ICoreWebView2HttpResponseHeaders, error) {
	var headers *ICoreWebView2HttpResponseHeaders
	hr, _, _ := i.vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	if headers == nil {

		return nil, fmt.Errorf("unknown error")
	}
	return headers, nil
}

func (i *ICoreWebView2WebResourceResponse) PutStatusCode(statusCode int) error {
	hr, _, _ := i.vtbl.PutStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(statusCode),
	)

	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}

	// Convert string 'reasonPhrase' to *uint16
	_reasonPhrase, err := UTF16PtrFromString(http.StatusText(statusCode))
	if err != nil {
		return err
	}

	hr, _, _ = i.vtbl.PutReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_reasonPhrase)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2WebResourceResponse) PutContent(content *IStream) error {
	hr, _, _ := i.vtbl.PutContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(content)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}

	return nil
}

func (i *ICoreWebView2WebResourceResponse) PutByteContent(content []byte) error {
	var stream *IStream
	if len(content) > 0 {
		// Create stream for response
		str, err := w32.SHCreateMemStream(content)
		if err != nil {
			return err
		}
		stream = (*IStream)(unsafe.Pointer(str))
		defer stream.Release()
	}

	return i.PutContent(stream)
}
