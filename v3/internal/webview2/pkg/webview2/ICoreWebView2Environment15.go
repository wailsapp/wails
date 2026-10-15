//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment15Vtbl struct {
	ICoreWebView2Environment14Vtbl
	CreateFindOptions ComProc
}

type ICoreWebView2Environment15 struct {
	Vtbl *ICoreWebView2Environment15Vtbl
}

func (i *ICoreWebView2Environment15) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment15) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Environment15 queries the object for its ICoreWebView2Environment15 interface. The receiver
// is the root of ICoreWebView2Environment15's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Environment) GetICoreWebView2Environment15() (*ICoreWebView2Environment15, error) {
	var result *ICoreWebView2Environment15

	iidICoreWebView2Environment15 := NewGUID("{2ac5ebfb-e654-5961-a667-7971885c7b27}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment15)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment15) CreateFindOptions() (*ICoreWebView2FindOptions, error) {

	var value *ICoreWebView2FindOptions

	hr, _, _ := i.Vtbl.CreateFindOptions.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
