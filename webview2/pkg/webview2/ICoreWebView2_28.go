//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_28Vtbl struct {
	IUnknownVtbl
	GetFind ComProc
}

type ICoreWebView2_28 struct {
	Vtbl *ICoreWebView2_28Vtbl
}

func (i *ICoreWebView2_28) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_28) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2_28() (*ICoreWebView2_28, error) {
	var result *ICoreWebView2_28

	iidICoreWebView2_28 := NewGUID("{62e50381-5bf5-51a8-aae0-f20a3a9c8a90}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_28)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_28) GetFind() (*ICoreWebView2Find, error) {

	var value *ICoreWebView2Find

	hr, _, _ := i.Vtbl.GetFind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
