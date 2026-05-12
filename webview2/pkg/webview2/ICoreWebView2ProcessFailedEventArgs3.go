//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ProcessFailedEventArgs3Vtbl struct {
	IUnknownVtbl
	GetFailureSourceModulePath ComProc
}

type ICoreWebView2ProcessFailedEventArgs3 struct {
	Vtbl *ICoreWebView2ProcessFailedEventArgs3Vtbl
}

func (i *ICoreWebView2ProcessFailedEventArgs3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ProcessFailedEventArgs3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2ProcessFailedEventArgs3() (*ICoreWebView2ProcessFailedEventArgs3, error) {
	var result *ICoreWebView2ProcessFailedEventArgs3

	iidICoreWebView2ProcessFailedEventArgs3 := NewGUID("{ab667428-094d-5fd1-b480-8b4c0fdbdf2f}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ProcessFailedEventArgs3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2ProcessFailedEventArgs3) GetFailureSourceModulePath() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetFailureSourceModulePath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}
