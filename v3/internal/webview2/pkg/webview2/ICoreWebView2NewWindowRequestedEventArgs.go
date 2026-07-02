//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2NewWindowRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetUri             ComProc
	PutNewWindow       ComProc
	GetNewWindow       ComProc
	PutHandled         ComProc
	GetHandled         ComProc
	GetIsUserInitiated ComProc
	GetDeferral        ComProc
	GetWindowFeatures  ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs struct {
	Vtbl *ICoreWebView2NewWindowRequestedEventArgsVtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetUri() (string, error) {
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

func (i *ICoreWebView2NewWindowRequestedEventArgs) PutNewWindow(newWindow *ICoreWebView2) error {

	hr, _, _ := i.Vtbl.PutNewWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(newWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetNewWindow() (*ICoreWebView2, error) {

	var newWindow *ICoreWebView2

	hr, _, _ := i.Vtbl.GetNewWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&newWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return newWindow, nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) PutHandled(handled bool) error {

	hr, _, _ := i.Vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&handled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetHandled() (bool, error) {
	// Create int32 to hold bool result
	var _handled int32

	hr, _, _ := i.Vtbl.GetHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_handled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	handled := _handled != 0
	return handled, nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetIsUserInitiated() (bool, error) {
	// Create int32 to hold bool result
	var _isUserInitiated int32

	hr, _, _ := i.Vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isUserInitiated)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isUserInitiated := _isUserInitiated != 0
	return isUserInitiated, nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

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

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetWindowFeatures() (*ICoreWebView2WindowFeatures, error) {

	var value *ICoreWebView2WindowFeatures

	hr, _, _ := i.Vtbl.GetWindowFeatures.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
