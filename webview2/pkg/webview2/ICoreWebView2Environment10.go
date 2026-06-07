//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment10Vtbl struct {
	IUnknownVtbl
	CreateCoreWebView2ControllerOptions ComProc
	CreateCoreWebView2ControllerWithOptions ComProc
	CreateCoreWebView2CompositionControllerWithOptions ComProc
}

type ICoreWebView2Environment10 struct {
	Vtbl *ICoreWebView2Environment10Vtbl
}

func (i *ICoreWebView2Environment10) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment10) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Environment10() (*ICoreWebView2Environment10, error) {
	var result *ICoreWebView2Environment10

	iidICoreWebView2Environment10 := NewGUID("{ee0eb9df-6f12-46ce-b53f-3f47b9c928e0}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment10)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment10) CreateCoreWebView2ControllerOptions() (*ICoreWebView2ControllerOptions, error) {

	var value *ICoreWebView2ControllerOptions

	hr, _, _ := i.Vtbl.CreateCoreWebView2ControllerOptions.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Environment10) CreateCoreWebView2ControllerWithOptions(ParentWindow HWND, options *ICoreWebView2ControllerOptions, handler *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler) error {


	hr, _, _ := i.Vtbl.CreateCoreWebView2ControllerWithOptions.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(ParentWindow),
		uintptr(unsafe.Pointer(options)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Environment10) CreateCoreWebView2CompositionControllerWithOptions(ParentWindow HWND, options *ICoreWebView2ControllerOptions, handler *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler) error {


	hr, _, _ := i.Vtbl.CreateCoreWebView2CompositionControllerWithOptions.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(ParentWindow),
		uintptr(unsafe.Pointer(options)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
