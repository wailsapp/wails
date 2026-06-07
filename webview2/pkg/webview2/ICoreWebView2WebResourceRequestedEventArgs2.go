//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2WebResourceRequestedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetRequestedSourceKind ComProc
}

type ICoreWebView2WebResourceRequestedEventArgs2 struct {
	Vtbl *ICoreWebView2WebResourceRequestedEventArgs2Vtbl
}

func (i *ICoreWebView2WebResourceRequestedEventArgs2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WebResourceRequestedEventArgs2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2WebResourceRequestedEventArgs2() (*ICoreWebView2WebResourceRequestedEventArgs2, error) {
	var result *ICoreWebView2WebResourceRequestedEventArgs2

	iidICoreWebView2WebResourceRequestedEventArgs2 := NewGUID("{9c562c24-b219-4d7f-92f6-b187fbbadd56}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2WebResourceRequestedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2WebResourceRequestedEventArgs2) GetRequestedSourceKind() (COREWEBVIEW2_WEB_RESOURCE_REQUEST_SOURCE_KINDS, error) {

	var value COREWEBVIEW2_WEB_RESOURCE_REQUEST_SOURCE_KINDS

	hr, _, _ := i.Vtbl.GetRequestedSourceKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
