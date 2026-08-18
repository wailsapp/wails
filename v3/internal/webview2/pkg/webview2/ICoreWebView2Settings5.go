//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Settings5Vtbl struct {
	ICoreWebView2Settings4Vtbl
	GetIsPinchZoomEnabled ComProc
	PutIsPinchZoomEnabled ComProc
}

type ICoreWebView2Settings5 struct {
	Vtbl *ICoreWebView2Settings5Vtbl
}

func (i *ICoreWebView2Settings5) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Settings5) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Settings5 queries the object for its ICoreWebView2Settings5 interface. The receiver
// is the root of ICoreWebView2Settings5's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Settings) GetICoreWebView2Settings5() (*ICoreWebView2Settings5, error) {
	var result *ICoreWebView2Settings5

	iidICoreWebView2Settings5 := NewGUID("{183e7052-1d03-43a0-ab99-98e043b66b39}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings5)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Settings5) GetIsPinchZoomEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsPinchZoomEnabled.Call(
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

func (i *ICoreWebView2Settings5) PutIsPinchZoomEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutIsPinchZoomEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
