//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ControllerOptions2Vtbl struct {
	ICoreWebView2ControllerOptionsVtbl
	GetScriptLocale ComProc
	PutScriptLocale ComProc
}

type ICoreWebView2ControllerOptions2 struct {
	Vtbl *ICoreWebView2ControllerOptions2Vtbl
}

func (i *ICoreWebView2ControllerOptions2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ControllerOptions2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2ControllerOptions2 queries the object for its ICoreWebView2ControllerOptions2 interface. The receiver
// is the root of ICoreWebView2ControllerOptions2's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2ControllerOptions) GetICoreWebView2ControllerOptions2() (*ICoreWebView2ControllerOptions2, error) {
	var result *ICoreWebView2ControllerOptions2

	iidICoreWebView2ControllerOptions2 := NewGUID("{06c991d8-9e7e-11ed-a8fc-0242ac120002}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ControllerOptions2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2ControllerOptions2) GetScriptLocale() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetScriptLocale.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2ControllerOptions2) PutScriptLocale(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutScriptLocale.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
