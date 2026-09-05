//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2NavigationCompletedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetHttpStatusCode ComProc
}

type ICoreWebView2NavigationCompletedEventArgs2 struct {
	Vtbl *ICoreWebView2NavigationCompletedEventArgs2Vtbl
}

func (i *ICoreWebView2NavigationCompletedEventArgs2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2NavigationCompletedEventArgs2() *ICoreWebView2NavigationCompletedEventArgs2 {
	var result *ICoreWebView2NavigationCompletedEventArgs2

	iidICoreWebView2NavigationCompletedEventArgs2 := NewGUID("{fdf8b738-ee1e-4db2-a329-8d7d7b74d792}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NavigationCompletedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2NavigationCompletedEventArgs2) GetHttpStatusCode() (int, error) {

	var value int

	hr, _, _ := i.Vtbl.GetHttpStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
