//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_16Vtbl struct {
	IUnknownVtbl
	Print            ComProc
	ShowPrintUI      ComProc
	PrintToPdfStream ComProc
}

type ICoreWebView2_16 struct {
	Vtbl *ICoreWebView2_16Vtbl
}

func (i *ICoreWebView2_16) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_16() *ICoreWebView2_16 {
	var result *ICoreWebView2_16

	iidICoreWebView2_16 := NewGUID("{0EB34DC9-9F91-41E1-8639-95CD5943906B}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_16)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_16) Print(printSettings *ICoreWebView2PrintSettings, handler *ICoreWebView2PrintCompletedHandler) error {

	hr, _, _ := i.Vtbl.Print.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(printSettings)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_16) ShowPrintUI(printDialogKind COREWEBVIEW2_PRINT_DIALOG_KIND) error {

	hr, _, _ := i.Vtbl.ShowPrintUI.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(printDialogKind),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_16) PrintToPdfStream(printSettings *ICoreWebView2PrintSettings, handler *ICoreWebView2PrintToPdfStreamCompletedHandler) error {

	hr, _, _ := i.Vtbl.PrintToPdfStream.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(printSettings)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
