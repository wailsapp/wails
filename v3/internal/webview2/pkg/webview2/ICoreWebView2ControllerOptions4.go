//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ControllerOptions4Vtbl struct {
	ICoreWebView2ControllerOptions3Vtbl
	GetAllowHostInputProcessing ComProc
	PutAllowHostInputProcessing ComProc
}

type ICoreWebView2ControllerOptions4 struct {
	Vtbl *ICoreWebView2ControllerOptions4Vtbl
}

func (i *ICoreWebView2ControllerOptions4) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ControllerOptions4) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2ControllerOptions4 queries the object for its ICoreWebView2ControllerOptions4 interface. The receiver
// is the root of ICoreWebView2ControllerOptions4's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2ControllerOptions) GetICoreWebView2ControllerOptions4() (*ICoreWebView2ControllerOptions4, error) {
	var result *ICoreWebView2ControllerOptions4

	iidICoreWebView2ControllerOptions4 := NewGUID("{21eb052f-ad39-555e-824a-c87b091d4d36}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ControllerOptions4)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2ControllerOptions4) GetAllowHostInputProcessing() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetAllowHostInputProcessing.Call(
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

func (i *ICoreWebView2ControllerOptions4) PutAllowHostInputProcessing(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutAllowHostInputProcessing.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
