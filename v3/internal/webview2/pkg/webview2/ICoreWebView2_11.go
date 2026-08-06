//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_11Vtbl struct {
	ICoreWebView2_10Vtbl
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


// GetICoreWebView2_11 queries the object for its ICoreWebView2_11 interface. The receiver
// is the root of ICoreWebView2_11's inheritance chain — the object that actually
// implements it.
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

	hr, _, _ := i.Vtbl.CallDevToolsProtocolMethodForSession.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_sessionId)),
		uintptr(unsafe.Pointer(_methodName)),
		uintptr(unsafe.Pointer(_parametersAsJson)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_11) AddContextMenuRequested(eventHandler *ICoreWebView2ContextMenuRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddContextMenuRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_11) RemoveContextMenuRequested(token EventRegistrationToken) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemoveContextMenuRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemoveContextMenuRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
