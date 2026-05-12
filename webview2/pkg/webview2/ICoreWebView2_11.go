//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_11Vtbl struct {
	IUnknownVtbl
	CallDevToolsProtocolMethodForSession ComProc
	AddContextMenuRequested ComProc
	RemoveContextMenuRequested ComProc
}

type ICoreWebView2_11 struct {
	Vtbl *ICoreWebView2_11Vtbl
}

func (i *ICoreWebView2_11) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_11) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2_11() (*ICoreWebView2_11, error) {
	var result *ICoreWebView2_11

	iidICoreWebView2_11 := NewGUID("{0be78e56-c193-4051-b943-23b460c08bdb}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_11)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_11) CallDevToolsProtocolMethodForSession(sessionId string, methodName string, parametersAsJson string, handler *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) error {

	// Convert string 'sessionId' to *uint16
	_sessionId, err := UTF16PtrFromString(sessionId)
	if err != nil {
		return err
	}
	// Convert string 'methodName' to *uint16
	_methodName, err := UTF16PtrFromString(methodName)
	if err != nil {
		return err
	}
	// Convert string 'parametersAsJson' to *uint16
	_parametersAsJson, err := UTF16PtrFromString(parametersAsJson)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.CallDevToolsProtocolMethodForSession.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_sessionId)),
		uintptr(unsafe.Pointer(_methodName)),
		uintptr(unsafe.Pointer(_parametersAsJson)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2_11) AddContextMenuRequested(eventHandler *ICoreWebView2ContextMenuRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddContextMenuRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2_11) RemoveContextMenuRequested(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveContextMenuRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
