//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment8Vtbl struct {
	IUnknownVtbl
	AddProcessInfosChanged ComProc
	RemoveProcessInfosChanged ComProc
	GetProcessInfos ComProc
}

type ICoreWebView2Environment8 struct {
	Vtbl *ICoreWebView2Environment8Vtbl
}

func (i *ICoreWebView2Environment8) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment8) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Environment8() (*ICoreWebView2Environment8, error) {
	var result *ICoreWebView2Environment8

	iidICoreWebView2Environment8 := NewGUID("{d6eb91dd-c3d2-45e5-bd29-6dc2bc4de9cf}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment8)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment8) AddProcessInfosChanged(eventHandler *ICoreWebView2ProcessInfosChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddProcessInfosChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Environment8) RemoveProcessInfosChanged(token EventRegistrationToken) error {


	hr, _, _ := i.Vtbl.RemoveProcessInfosChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Environment8) GetProcessInfos() (*ICoreWebView2ProcessInfoCollection, error) {

	var value *ICoreWebView2ProcessInfoCollection

	hr, _, _ := i.Vtbl.GetProcessInfos.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
