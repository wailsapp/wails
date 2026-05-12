//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2PermissionRequestedEventArgs3Vtbl struct {
	IUnknownVtbl
	GetSavesInProfile ComProc
	PutSavesInProfile ComProc
}

type ICoreWebView2PermissionRequestedEventArgs3 struct {
	Vtbl *ICoreWebView2PermissionRequestedEventArgs3Vtbl
}

func (i *ICoreWebView2PermissionRequestedEventArgs3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2PermissionRequestedEventArgs3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2PermissionRequestedEventArgs3() (*ICoreWebView2PermissionRequestedEventArgs3, error) {
	var result *ICoreWebView2PermissionRequestedEventArgs3

	iidICoreWebView2PermissionRequestedEventArgs3 := NewGUID("{e61670bc-3dce-4177-86d2-c629ae3cb6ac}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2PermissionRequestedEventArgs3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2PermissionRequestedEventArgs3) GetSavesInProfile() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetSavesInProfile.Call(
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

func (i *ICoreWebView2PermissionRequestedEventArgs3) PutSavesInProfile(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutSavesInProfile.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
