//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment13Vtbl struct {
	IUnknownVtbl
	GetProcessExtendedInfos ComProc
}

type ICoreWebView2Environment13 struct {
	Vtbl *ICoreWebView2Environment13Vtbl
}

func (i *ICoreWebView2Environment13) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment13) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Environment13() (*ICoreWebView2Environment13, error) {
	var result *ICoreWebView2Environment13

	iidICoreWebView2Environment13 := NewGUID("{af641f58-72b2-11ee-b962-0242ac120002}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment13)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment13) GetProcessExtendedInfos(handler *ICoreWebView2GetProcessExtendedInfosCompletedHandler) error {


	hr, _, err := i.Vtbl.GetProcessExtendedInfos.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
