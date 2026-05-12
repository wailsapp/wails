//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_21Vtbl struct {
	IUnknownVtbl
	ExecuteScriptWithResult ComProc
}

type ICoreWebView2_21 struct {
	Vtbl *ICoreWebView2_21Vtbl
}

func (i *ICoreWebView2_21) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_21) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2_21() (*ICoreWebView2_21, error) {
	var result *ICoreWebView2_21

	iidICoreWebView2_21 := NewGUID("{c4980dea-587b-43b9-8143-3ef3bf552d95}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_21)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_21) ExecuteScriptWithResult(javaScript string, handler *ICoreWebView2ExecuteScriptWithResultCompletedHandler) error {

	// Convert string 'javaScript' to *uint16
	_javaScript, err := UTF16PtrFromString(javaScript)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.ExecuteScriptWithResult.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_javaScript)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
