//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Frame2Vtbl struct {
	IUnknownVtbl
	AddNavigationStarting ComProc
	RemoveNavigationStarting ComProc
	AddContentLoading ComProc
	RemoveContentLoading ComProc
	AddNavigationCompleted ComProc
	RemoveNavigationCompleted ComProc
	AddDOMContentLoaded ComProc
	RemoveDOMContentLoaded ComProc
	ExecuteScript ComProc
	PostWebMessageAsJson ComProc
	PostWebMessageAsString ComProc
	AddWebMessageReceived ComProc
	RemoveWebMessageReceived ComProc
}

type ICoreWebView2Frame2 struct {
	Vtbl *ICoreWebView2Frame2Vtbl
}

func (i *ICoreWebView2Frame2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Frame2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Frame2() (*ICoreWebView2Frame2, error) {
	var result *ICoreWebView2Frame2

	iidICoreWebView2Frame2 := NewGUID("{7a6a5834-d185-4dbf-b63f-4a9bc43107d4}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Frame2) AddNavigationStarting(eventHandler *ICoreWebView2FrameNavigationStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddNavigationStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Frame2) RemoveNavigationStarting(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveNavigationStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame2) AddContentLoading(eventHandler *ICoreWebView2FrameContentLoadingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddContentLoading.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Frame2) RemoveContentLoading(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveContentLoading.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame2) AddNavigationCompleted(eventHandler *ICoreWebView2FrameNavigationCompletedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Frame2) RemoveNavigationCompleted(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveNavigationCompleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame2) AddDOMContentLoaded(eventHandler *ICoreWebView2FrameDOMContentLoadedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddDOMContentLoaded.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Frame2) RemoveDOMContentLoaded(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveDOMContentLoaded.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame2) ExecuteScript(javaScript string, handler *ICoreWebView2ExecuteScriptCompletedHandler) error {

	// Convert string 'javaScript' to *uint16
	_javaScript, err := UTF16PtrFromString(javaScript)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.ExecuteScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_javaScript)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame2) PostWebMessageAsJson(webMessageAsJson string) error {

	// Convert string 'webMessageAsJson' to *uint16
	_webMessageAsJson, err := UTF16PtrFromString(webMessageAsJson)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.PostWebMessageAsJson.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_webMessageAsJson)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame2) PostWebMessageAsString(webMessageAsString string) error {

	// Convert string 'webMessageAsString' to *uint16
	_webMessageAsString, err := UTF16PtrFromString(webMessageAsString)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.PostWebMessageAsString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_webMessageAsString)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame2) AddWebMessageReceived(handler *ICoreWebView2FrameWebMessageReceivedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddWebMessageReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Frame2) RemoveWebMessageReceived(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveWebMessageReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
