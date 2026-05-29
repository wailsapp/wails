//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2SaveFileSecurityCheckStartingEventArgsVtbl struct {
	IUnknownVtbl
	GetCancelSave ComProc
	PutCancelSave ComProc
	GetDocumentOriginUri ComProc
	GetFileExtension ComProc
	GetFilePath ComProc
	GetSuppressDefaultPolicy ComProc
	PutSuppressDefaultPolicy ComProc
	GetDeferral ComProc
}

type ICoreWebView2SaveFileSecurityCheckStartingEventArgs struct {
	Vtbl *ICoreWebView2SaveFileSecurityCheckStartingEventArgsVtbl
}

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) GetCancelSave() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetCancelSave.Call(
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

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) PutCancelSave(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutCancelSave.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) GetDocumentOriginUri() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetDocumentOriginUri.Call(
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

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) GetFileExtension() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetFileExtension.Call(
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

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) GetFilePath() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetFilePath.Call(
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

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) GetSuppressDefaultPolicy() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetSuppressDefaultPolicy.Call(
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

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) PutSuppressDefaultPolicy(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutSuppressDefaultPolicy.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var value *ICoreWebView2Deferral

	hr, _, _ := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
