//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Profile5Vtbl struct {
	IUnknownVtbl
	GetCookieManager ComProc
}

type ICoreWebView2Profile5 struct {
	Vtbl *ICoreWebView2Profile5Vtbl
}

func (i *ICoreWebView2Profile5) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Profile5) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Profile5() (*ICoreWebView2Profile5, error) {
	var result *ICoreWebView2Profile5

	iidICoreWebView2Profile5 := NewGUID("{2ee5b76e-6e80-4df2-bcd3-d4ec3340a01b}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile5)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Profile5) GetCookieManager() (*ICoreWebView2CookieManager, error) {

	var value *ICoreWebView2CookieManager

	hr, _, _ := i.Vtbl.GetCookieManager.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
