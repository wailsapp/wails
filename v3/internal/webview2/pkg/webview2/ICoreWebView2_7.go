//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_7Vtbl struct {
	IUnknownVtbl
	PrintToPdf ComProc
}

type ICoreWebView2_7 struct {
	Vtbl *ICoreWebView2_7Vtbl
}

func (i *ICoreWebView2_7) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_7() *ICoreWebView2_7 {
	var result *ICoreWebView2_7

	iidICoreWebView2_7 := NewGUID("{79c24d83-09a3-45ae-9418-487f32a58740}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_7)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_7) PrintToPdf(ResultFilePath string, printSettings *ICoreWebView2PrintSettings, handler *ICoreWebView2PrintToPdfCompletedHandler) error {

	// Convert string 'ResultFilePath' to *uint16
	_ResultFilePath, err := UTF16PtrFromString(ResultFilePath)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PrintToPdf.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_ResultFilePath)),
		uintptr(unsafe.Pointer(printSettings)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
