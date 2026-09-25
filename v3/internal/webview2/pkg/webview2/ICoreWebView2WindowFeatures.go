//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2WindowFeaturesVtbl struct {
	IUnknownVtbl
	GetHasPosition ComProc
	GetHasSize ComProc
	GetLeft ComProc
	GetTop ComProc
	GetHeight ComProc
	GetWidth ComProc
	GetShouldDisplayMenuBar ComProc
	GetShouldDisplayStatus ComProc
	GetShouldDisplayToolbar ComProc
	GetShouldDisplayScrollBars ComProc
}

type ICoreWebView2WindowFeatures struct {
	Vtbl *ICoreWebView2WindowFeaturesVtbl
}

func (i *ICoreWebView2WindowFeatures) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WindowFeatures) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2WindowFeatures) GetHasPosition() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetHasPosition.Call(
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

func (i *ICoreWebView2WindowFeatures) GetHasSize() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetHasSize.Call(
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

func (i *ICoreWebView2WindowFeatures) GetLeft() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetLeft.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2WindowFeatures) GetTop() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetTop.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2WindowFeatures) GetHeight() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetHeight.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2WindowFeatures) GetWidth() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetWidth.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2WindowFeatures) GetShouldDisplayMenuBar() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetShouldDisplayMenuBar.Call(
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

func (i *ICoreWebView2WindowFeatures) GetShouldDisplayStatus() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetShouldDisplayStatus.Call(
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

func (i *ICoreWebView2WindowFeatures) GetShouldDisplayToolbar() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetShouldDisplayToolbar.Call(
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

func (i *ICoreWebView2WindowFeatures) GetShouldDisplayScrollBars() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetShouldDisplayScrollBars.Call(
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
