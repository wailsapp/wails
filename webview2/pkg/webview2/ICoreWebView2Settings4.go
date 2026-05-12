//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Settings4Vtbl struct {
	IUnknownVtbl
	GetIsPasswordAutosaveEnabled ComProc
	PutIsPasswordAutosaveEnabled ComProc
	GetIsGeneralAutofillEnabled ComProc
	PutIsGeneralAutofillEnabled ComProc
}

type ICoreWebView2Settings4 struct {
	Vtbl *ICoreWebView2Settings4Vtbl
}

func (i *ICoreWebView2Settings4) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Settings4) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Settings4() (*ICoreWebView2Settings4, error) {
	var result *ICoreWebView2Settings4

	iidICoreWebView2Settings4 := NewGUID("{cb56846c-4168-4d53-b04f-03b6d6796ff2}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings4)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Settings4) GetIsPasswordAutosaveEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetIsPasswordAutosaveEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, err
}

func (i *ICoreWebView2Settings4) PutIsPasswordAutosaveEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutIsPasswordAutosaveEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Settings4) GetIsGeneralAutofillEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetIsGeneralAutofillEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, err
}

func (i *ICoreWebView2Settings4) PutIsGeneralAutofillEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutIsGeneralAutofillEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
