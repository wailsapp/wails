//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2WebResourceResponseVtbl struct {
	IUnknownVtbl
	GetContent ComProc
	PutContent ComProc
	GetHeaders ComProc
	GetStatusCode ComProc
	PutStatusCode ComProc
	GetReasonPhrase ComProc
	PutReasonPhrase ComProc
}

type ICoreWebView2WebResourceResponse struct {
	Vtbl *ICoreWebView2WebResourceResponseVtbl
}

func (i *ICoreWebView2WebResourceResponse) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WebResourceResponse) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2WebResourceResponse) GetContent() (*IStream, error) {

	var content *IStream

	hr, _, err := i.Vtbl.GetContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&content)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return content, err
}

func (i *ICoreWebView2WebResourceResponse) PutContent(content *IStream) error {


	hr, _, err := i.Vtbl.PutContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(content)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2WebResourceResponse) GetHeaders() (*ICoreWebView2HttpResponseHeaders, error) {

	var headers *ICoreWebView2HttpResponseHeaders

	hr, _, err := i.Vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return headers, err
}

func (i *ICoreWebView2WebResourceResponse) GetStatusCode() (int, error) {

	var statusCode int

	hr, _, err := i.Vtbl.GetStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&statusCode)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return statusCode, err
}

func (i *ICoreWebView2WebResourceResponse) PutStatusCode(statusCode int) error {


	hr, _, err := i.Vtbl.PutStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(statusCode),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2WebResourceResponse) GetReasonPhrase() (string, error) {
	// Create *uint16 to hold result
	var _reasonPhrase *uint16


	hr, _, err := i.Vtbl.GetReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_reasonPhrase)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	reasonPhrase := UTF16PtrToString(_reasonPhrase)
	CoTaskMemFree(unsafe.Pointer(_reasonPhrase))
	return reasonPhrase, err
}

func (i *ICoreWebView2WebResourceResponse) PutReasonPhrase(reasonPhrase string) error {

	// Convert string 'reasonPhrase' to *uint16
	_reasonPhrase, err := UTF16PtrFromString(reasonPhrase)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.PutReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_reasonPhrase)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
