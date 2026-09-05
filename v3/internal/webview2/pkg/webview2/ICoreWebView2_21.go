//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_21Vtbl struct {
	IUnknownVtbl
	ExecuteScriptWithResult ComProc
}

type ICoreWebView2_21 struct {
	Vtbl *ICoreWebView2_21Vtbl
}

func (i *ICoreWebView2_21) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_21() *ICoreWebView2_21 {
	var result *ICoreWebView2_21

	iidICoreWebView2_21 := NewGUID("{c4980dea-587b-43b9-8143-3ef3bf552d95}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_21)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_21) ExecuteScriptWithResult(javaScript string, handler *ICoreWebView2ExecuteScriptWithResultCompletedHandler) error {

	// Convert string 'javaScript' to *uint16
	_javaScript, err := UTF16PtrFromString(javaScript)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.ExecuteScriptWithResult.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_javaScript)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
