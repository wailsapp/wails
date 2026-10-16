//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type iCoreWebView2Environment3Vtbl struct {
	_IUnknownVtbl
	CreateCoreWebView2Controller            ComProc
	CreateWebResourceResponse               ComProc
	GetBrowserVersionString                 ComProc
	AddNewBrowserVersionAvailable           ComProc
	RemoveNewBrowserVersionAvailable        ComProc
	CreateWebResourceRequest                ComProc
	CreateCoreWebView2CompositionController ComProc
	CreateCoreWebView2PointerInfo           ComProc
}

type ICoreWebView2Environment3 struct {
	vtbl *iCoreWebView2Environment3Vtbl
}

func (e *ICoreWebView2Environment3) AddRef() uintptr {
	ret, _, _ := e.vtbl.AddRef.Call(uintptr(unsafe.Pointer(e)))

	return ret
}

func (e *ICoreWebView2Environment3) Release() uintptr {
	ret, _, _ := e.vtbl.Release.Call(uintptr(unsafe.Pointer(e)))

	return ret
}

func (e *ICoreWebView2Environment) GetICoreWebView2Environment3() *ICoreWebView2Environment3 {
	var result *ICoreWebView2Environment3

	iidICoreWebView2Environment3 := NewGUID("{80a22ae3-be7c-4ce2-afe1-5a50056cdeeb}")
	_, _, _ = e.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(e)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *ICoreWebView2Environment3) CreateCoreWebView2CompositionController(parentWindow uintptr, handler *iCoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler) error {
	hr, _, _ := e.vtbl.CreateCoreWebView2CompositionController.Call(
		uintptr(unsafe.Pointer(e)),
		parentWindow,
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}

	return nil
}
