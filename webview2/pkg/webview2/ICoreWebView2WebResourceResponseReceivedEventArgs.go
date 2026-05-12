//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2WebResourceResponseReceivedEventArgsVtbl struct {
	IUnknownVtbl
	GetRequest ComProc
	GetResponse ComProc
}

type ICoreWebView2WebResourceResponseReceivedEventArgs struct {
	Vtbl *ICoreWebView2WebResourceResponseReceivedEventArgsVtbl
}

func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest, error) {

	var value *ICoreWebView2WebResourceRequest

	hr, _, err := i.Vtbl.GetRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) GetResponse() (*ICoreWebView2WebResourceResponseView, error) {

	var value *ICoreWebView2WebResourceResponseView

	hr, _, err := i.Vtbl.GetResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
