//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2WebMessageReceivedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetAdditionalObjects ComProc
}

type ICoreWebView2WebMessageReceivedEventArgs2 struct {
	Vtbl *ICoreWebView2WebMessageReceivedEventArgs2Vtbl
}

func (i *ICoreWebView2WebMessageReceivedEventArgs2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WebMessageReceivedEventArgs2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2WebMessageReceivedEventArgs2() (*ICoreWebView2WebMessageReceivedEventArgs2, error) {
	var result *ICoreWebView2WebMessageReceivedEventArgs2

	iidICoreWebView2WebMessageReceivedEventArgs2 := NewGUID("{06fc7ab7-c90c-4297-9389-33ca01cf6d5e}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2WebMessageReceivedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2WebMessageReceivedEventArgs2) GetAdditionalObjects() (*ICoreWebView2ObjectCollectionView, error) {

	var value *ICoreWebView2ObjectCollectionView

	hr, _, err := i.Vtbl.GetAdditionalObjects.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
