//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment6Vtbl struct {
	ICoreWebView2Environment5Vtbl
	CreatePrintSettings ComProc
}

type ICoreWebView2Environment6 struct {
	Vtbl *ICoreWebView2Environment6Vtbl
}

func (i *ICoreWebView2Environment6) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment6) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Environment6 queries the object for its ICoreWebView2Environment6 interface. The receiver
// is the root of ICoreWebView2Environment6's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Environment) GetICoreWebView2Environment6() (*ICoreWebView2Environment6, error) {
	var result *ICoreWebView2Environment6

	iidICoreWebView2Environment6 := NewGUID("{e59ee362-acbd-4857-9a8e-d3644d9459a9}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment6)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment6) CreatePrintSettings() (*ICoreWebView2PrintSettings, error) {

	var value *ICoreWebView2PrintSettings

	hr, _, _ := i.Vtbl.CreatePrintSettings.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
