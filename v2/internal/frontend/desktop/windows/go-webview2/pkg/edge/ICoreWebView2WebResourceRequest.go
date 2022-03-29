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
	return i.AddRef()
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

func (i *ICoreWebView2WebResourceRequest) Release() error {
	return i.vtbl.CallRelease(unsafe.Pointer(i))
}
