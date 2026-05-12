//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_26Vtbl struct {
	IUnknownVtbl
	AddSaveFileSecurityCheckStarting ComProc
	RemoveSaveFileSecurityCheckStarting ComProc
}

type ICoreWebView2_26 struct {
	Vtbl *ICoreWebView2_26Vtbl
}

func (i *ICoreWebView2_26) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_26) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2_26() (*ICoreWebView2_26, error) {
	var result *ICoreWebView2_26

	iidICoreWebView2_26 := NewGUID("{806268b8-f897-5685-88e5-c45fca0b1a48}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_26)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_26) AddSaveFileSecurityCheckStarting(eventHandler *ICoreWebView2SaveFileSecurityCheckStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddSaveFileSecurityCheckStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2_26) RemoveSaveFileSecurityCheckStarting(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveSaveFileSecurityCheckStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
