//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NavigationStartingEventArgs2Vtbl struct {
	ICoreWebView2NavigationStartingEventArgsVtbl
	GetAdditionalAllowedFrameAncestors ComProc
	PutAdditionalAllowedFrameAncestors ComProc
}

type ICoreWebView2NavigationStartingEventArgs2 struct {
	Vtbl *ICoreWebView2NavigationStartingEventArgs2Vtbl
}

func (i *ICoreWebView2NavigationStartingEventArgs2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NavigationStartingEventArgs2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2NavigationStartingEventArgs2 queries the object for its ICoreWebView2NavigationStartingEventArgs2 interface. The receiver
// is the root of ICoreWebView2NavigationStartingEventArgs2's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2NavigationStartingEventArgs) GetICoreWebView2NavigationStartingEventArgs2() (*ICoreWebView2NavigationStartingEventArgs2, error) {
	var result *ICoreWebView2NavigationStartingEventArgs2

	iidICoreWebView2NavigationStartingEventArgs2 := NewGUID("{9086be93-91aa-472d-a7e0-579f2ba006ad}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NavigationStartingEventArgs2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2NavigationStartingEventArgs2) GetAdditionalAllowedFrameAncestors() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetAdditionalAllowedFrameAncestors.Call(
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

func (i *ICoreWebView2NavigationStartingEventArgs2) PutAdditionalAllowedFrameAncestors(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutAdditionalAllowedFrameAncestors.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
