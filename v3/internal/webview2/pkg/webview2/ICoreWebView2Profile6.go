//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Profile6Vtbl struct {
	ICoreWebView2Profile5Vtbl
	GetIsPasswordAutosaveEnabled ComProc
	PutIsPasswordAutosaveEnabled ComProc
	GetIsGeneralAutofillEnabled ComProc
	PutIsGeneralAutofillEnabled ComProc
}

type ICoreWebView2Profile6 struct {
	Vtbl *ICoreWebView2Profile6Vtbl
}

func (i *ICoreWebView2Profile6) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Profile6) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Profile6 queries the object for its ICoreWebView2Profile6 interface. The receiver
// is the root of ICoreWebView2Profile6's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Profile) GetICoreWebView2Profile6() (*ICoreWebView2Profile6, error) {
	var result *ICoreWebView2Profile6

	iidICoreWebView2Profile6 := NewGUID("{BD82FA6A-1D65-4C33-B2B4-0393020CC61B}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile6)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Profile6) GetIsPasswordAutosaveEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsPasswordAutosaveEnabled.Call(
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

func (i *ICoreWebView2Profile6) PutIsPasswordAutosaveEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutIsPasswordAutosaveEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Profile6) GetIsGeneralAutofillEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsGeneralAutofillEnabled.Call(
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

func (i *ICoreWebView2Profile6) PutIsGeneralAutofillEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutIsGeneralAutofillEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
