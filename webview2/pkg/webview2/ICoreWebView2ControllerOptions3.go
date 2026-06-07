//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ControllerOptions3Vtbl struct {
	IUnknownVtbl
	GetDefaultBackgroundColor ComProc
	PutDefaultBackgroundColor ComProc
}

type ICoreWebView2ControllerOptions3 struct {
	Vtbl *ICoreWebView2ControllerOptions3Vtbl
}

func (i *ICoreWebView2ControllerOptions3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ControllerOptions3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2ControllerOptions3() (*ICoreWebView2ControllerOptions3, error) {
	var result *ICoreWebView2ControllerOptions3

	iidICoreWebView2ControllerOptions3 := NewGUID("{b32b191a-8998-57ca-b7cb-e04617e4ce4a}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ControllerOptions3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2ControllerOptions3) GetDefaultBackgroundColor() (COREWEBVIEW2_COLOR, error) {

	var value COREWEBVIEW2_COLOR

	hr, _, _ := i.Vtbl.GetDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return COREWEBVIEW2_COLOR{}, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ControllerOptions3) PutDefaultBackgroundColor(value COREWEBVIEW2_COLOR) error {


	hr, _, _ := i.Vtbl.PutDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
