//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Controller4Vtbl struct {
	ICoreWebView2Controller3Vtbl
	GetAllowExternalDrop ComProc
	PutAllowExternalDrop ComProc
}

type ICoreWebView2Controller4 struct {
	Vtbl *ICoreWebView2Controller4Vtbl
}

func (i *ICoreWebView2Controller4) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Controller4) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Controller4 queries the object for its ICoreWebView2Controller4 interface. The receiver
// is the root of ICoreWebView2Controller4's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Controller) GetICoreWebView2Controller4() (*ICoreWebView2Controller4, error) {
	var result *ICoreWebView2Controller4

	iidICoreWebView2Controller4 := NewGUID("{97d418d5-a426-4e49-a151-e1a10f327d9e}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller4)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Controller4) GetAllowExternalDrop() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetAllowExternalDrop.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, nil
}

func (i *ICoreWebView2Controller4) PutAllowExternalDrop(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutAllowExternalDrop.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
