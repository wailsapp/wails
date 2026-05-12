//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2DevToolsProtocolEventReceiverVtbl struct {
	IUnknownVtbl
	AddDevToolsProtocolEventReceived ComProc
	RemoveDevToolsProtocolEventReceived ComProc
}

type ICoreWebView2DevToolsProtocolEventReceiver struct {
	Vtbl *ICoreWebView2DevToolsProtocolEventReceiverVtbl
}

func (i *ICoreWebView2DevToolsProtocolEventReceiver) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2DevToolsProtocolEventReceiver) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2DevToolsProtocolEventReceiver) AddDevToolsProtocolEventReceived(eventHandler *ICoreWebView2DevToolsProtocolEventReceivedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddDevToolsProtocolEventReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2DevToolsProtocolEventReceiver) RemoveDevToolsProtocolEventReceived(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveDevToolsProtocolEventReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
