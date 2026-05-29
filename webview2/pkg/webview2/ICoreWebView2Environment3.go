//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment3Vtbl struct {
	IUnknownVtbl
	CreateCoreWebView2CompositionController ComProc
	CreateCoreWebView2PointerInfo ComProc
}

type ICoreWebView2Environment3 struct {
	Vtbl *ICoreWebView2Environment3Vtbl
}

func (i *ICoreWebView2Environment3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Environment3() (*ICoreWebView2Environment3, error) {
	var result *ICoreWebView2Environment3

	iidICoreWebView2Environment3 := NewGUID("{80a22ae3-be7c-4ce2-afe1-5a50056cdeeb}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment3) CreateCoreWebView2CompositionController(ParentWindow HWND, handler *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler) error {


	hr, _, _ := i.Vtbl.CreateCoreWebView2CompositionController.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(ParentWindow),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Environment3) CreateCoreWebView2PointerInfo() (*ICoreWebView2PointerInfo, error) {

	var value *ICoreWebView2PointerInfo

	hr, _, _ := i.Vtbl.CreateCoreWebView2PointerInfo.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
