//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2FindOptionsVtbl struct {
	IUnknownVtbl
	GetFindTerm ComProc
	PutFindTerm ComProc
	GetIsCaseSensitive ComProc
	PutIsCaseSensitive ComProc
	GetShouldHighlightAllMatches ComProc
	PutShouldHighlightAllMatches ComProc
	GetShouldMatchWord ComProc
	PutShouldMatchWord ComProc
	GetSuppressDefaultFindDialog ComProc
	PutSuppressDefaultFindDialog ComProc
}

type ICoreWebView2FindOptions struct {
	Vtbl *ICoreWebView2FindOptionsVtbl
}

func (i *ICoreWebView2FindOptions) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FindOptions) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2FindOptions) GetFindTerm() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetFindTerm.Call(
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

func (i *ICoreWebView2FindOptions) PutFindTerm(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutFindTerm.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2FindOptions) GetIsCaseSensitive() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsCaseSensitive.Call(
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

func (i *ICoreWebView2FindOptions) PutIsCaseSensitive(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutIsCaseSensitive.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2FindOptions) GetShouldHighlightAllMatches() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetShouldHighlightAllMatches.Call(
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

func (i *ICoreWebView2FindOptions) PutShouldHighlightAllMatches(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutShouldHighlightAllMatches.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2FindOptions) GetShouldMatchWord() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetShouldMatchWord.Call(
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

func (i *ICoreWebView2FindOptions) PutShouldMatchWord(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutShouldMatchWord.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2FindOptions) GetSuppressDefaultFindDialog() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetSuppressDefaultFindDialog.Call(
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

func (i *ICoreWebView2FindOptions) PutSuppressDefaultFindDialog(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutSuppressDefaultFindDialog.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
