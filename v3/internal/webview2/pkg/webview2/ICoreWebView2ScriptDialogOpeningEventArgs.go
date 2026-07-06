//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ScriptDialogOpeningEventArgsVtbl struct {
	IUnknownVtbl
	GetUri         ComProc
	GetKind        ComProc
	GetMessage     ComProc
	Accept         ComProc
	GetDefaultText ComProc
	GetResultText  ComProc
	PutResultText  ComProc
	GetDeferral    ComProc
}

type ICoreWebView2ScriptDialogOpeningEventArgs struct {
	Vtbl *ICoreWebView2ScriptDialogOpeningEventArgsVtbl
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16

	hr, _, _ := i.Vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	uri := UTF16PtrToString(_uri)
	CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) GetKind() (COREWEBVIEW2_SCRIPT_DIALOG_KIND, error) {

	var kind COREWEBVIEW2_SCRIPT_DIALOG_KIND

	hr, _, _ := i.Vtbl.GetKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&kind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return kind, nil
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) GetMessage() (string, error) {
	// Create *uint16 to hold result
	var _message *uint16

	hr, _, _ := i.Vtbl.GetMessage.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_message)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	message := UTF16PtrToString(_message)
	CoTaskMemFree(unsafe.Pointer(_message))
	return message, nil
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) Accept() error {

	hr, _, _ := i.Vtbl.Accept.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) GetDefaultText() (string, error) {
	// Create *uint16 to hold result
	var _defaultText *uint16

	hr, _, _ := i.Vtbl.GetDefaultText.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_defaultText)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	defaultText := UTF16PtrToString(_defaultText)
	CoTaskMemFree(unsafe.Pointer(_defaultText))
	return defaultText, nil
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) GetResultText() (string, error) {
	// Create *uint16 to hold result
	var _resultText *uint16

	hr, _, _ := i.Vtbl.GetResultText.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_resultText)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	resultText := UTF16PtrToString(_resultText)
	CoTaskMemFree(unsafe.Pointer(_resultText))
	return resultText, nil
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) PutResultText(resultText string) error {

	// Convert string 'resultText' to *uint16
	_resultText, err := UTF16PtrFromString(resultText)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutResultText.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_resultText)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ScriptDialogOpeningEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, _ := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, nil
}
