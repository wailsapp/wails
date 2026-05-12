//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2WebResourceRequestVtbl struct {
	IUnknownVtbl
	GetUri ComProc
	PutUri ComProc
	GetMethod ComProc
	PutMethod ComProc
	GetContent ComProc
	PutContent ComProc
	GetHeaders ComProc
}

type ICoreWebView2WebResourceRequest struct {
	Vtbl *ICoreWebView2WebResourceRequestVtbl
}

func (i *ICoreWebView2WebResourceRequest) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WebResourceRequest) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2WebResourceRequest) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16


	hr, _, err := i.Vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	uri := UTF16PtrToString(_uri)
	CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, err
}

func (i *ICoreWebView2WebResourceRequest) PutUri(uri string) error {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.PutUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2WebResourceRequest) GetMethod() (string, error) {
	// Create *uint16 to hold result
	var _method *uint16


	hr, _, err := i.Vtbl.GetMethod.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_method)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	method := UTF16PtrToString(_method)
	CoTaskMemFree(unsafe.Pointer(_method))
	return method, err
}

func (i *ICoreWebView2WebResourceRequest) PutMethod(method string) error {

	// Convert string 'method' to *uint16
	_method, err := UTF16PtrFromString(method)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.PutMethod.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_method)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2WebResourceRequest) GetContent() (*IStream, error) {

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

func (i *ICoreWebView2WebResourceRequest) PutContent(content *IStream) error {


	hr, _, err := i.Vtbl.PutContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(content)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2WebResourceRequest) GetHeaders() (*ICoreWebView2HttpRequestHeaders, error) {

	var headers *ICoreWebView2HttpRequestHeaders

	hr, _, err := i.Vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return headers, err
}
