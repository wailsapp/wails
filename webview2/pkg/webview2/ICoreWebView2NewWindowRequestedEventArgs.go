//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NewWindowRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetUri ComProc
	PutNewWindow ComProc
	GetNewWindow ComProc
	PutHandled ComProc
	GetHandled ComProc
	GetIsUserInitiated ComProc
	GetDeferral ComProc
	GetWindowFeatures ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs struct {
	Vtbl *ICoreWebView2NewWindowRequestedEventArgsVtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2NewWindowRequestedEventArgs) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16


	hr, _, err := i.Vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	uri := UTF16PtrToString(_uri)
	CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, err
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) PutNewWindow(newWindow *ICoreWebView2) error {


	hr, _, err := i.Vtbl.PutNewWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(newWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetNewWindow() (*ICoreWebView2, error) {

	var newWindow *ICoreWebView2

	hr, _, err := i.Vtbl.GetNewWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&newWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return newWindow, err
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) PutHandled(handled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _handled int32
	if handled {
		_handled = 1
	}

	hr, _, err := i.Vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_handled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetHandled() (bool, error) {
	// Create int32 to hold bool result
	var _handled int32

	hr, _, err := i.Vtbl.GetHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_handled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    handled := _handled != 0
	return handled, err
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetIsUserInitiated() (bool, error) {
	// Create int32 to hold bool result
	var _isUserInitiated int32

	hr, _, err := i.Vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isUserInitiated)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    isUserInitiated := _isUserInitiated != 0
	return isUserInitiated, err
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, err := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, err
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetWindowFeatures() (*ICoreWebView2WindowFeatures, error) {

	var value *ICoreWebView2WindowFeatures

	hr, _, err := i.Vtbl.GetWindowFeatures.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
