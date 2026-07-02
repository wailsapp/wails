//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Controller3Vtbl struct {
	IUnknownVtbl
	GetRasterizationScale              ComProc
	PutRasterizationScale              ComProc
	GetShouldDetectMonitorScaleChanges ComProc
	PutShouldDetectMonitorScaleChanges ComProc
	AddRasterizationScaleChanged       ComProc
	RemoveRasterizationScaleChanged    ComProc
	GetBoundsMode                      ComProc
	PutBoundsMode                      ComProc
}

type ICoreWebView2Controller3 struct {
	Vtbl *ICoreWebView2Controller3Vtbl
}

func (i *ICoreWebView2Controller3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Controller3() *ICoreWebView2Controller3 {
	var result *ICoreWebView2Controller3

	iidICoreWebView2Controller3 := NewGUID("{f9614724-5d2b-41dc-aef7-73d62b51543b}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Controller3) GetRasterizationScale() (float64, error) {

	var scale float64

	hr, _, _ := i.Vtbl.GetRasterizationScale.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&scale)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return scale, nil
}

func (i *ICoreWebView2Controller3) PutRasterizationScale(scale float64) error {

	hr, _, _ := i.Vtbl.PutRasterizationScale.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&scale)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller3) GetShouldDetectMonitorScaleChanges() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetShouldDetectMonitorScaleChanges.Call(
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

func (i *ICoreWebView2Controller3) PutShouldDetectMonitorScaleChanges(value bool) error {
	var intValue uintptr
	if value {
		intValue = 1
	} else {
		intValue = 0
	}

	hr, _, _ := i.Vtbl.PutShouldDetectMonitorScaleChanges.Call(
		uintptr(unsafe.Pointer(i)),
		intValue,
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller3) AddRasterizationScaleChanged(eventHandler *ICoreWebView2RasterizationScaleChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddRasterizationScaleChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Controller3) RemoveRasterizationScaleChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveRasterizationScaleChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller3) GetBoundsMode() (COREWEBVIEW2_BOUNDS_MODE, error) {

	var boundsMode COREWEBVIEW2_BOUNDS_MODE

	hr, _, _ := i.Vtbl.GetBoundsMode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&boundsMode)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return boundsMode, nil
}

func (i *ICoreWebView2Controller3) PutBoundsMode(boundsMode COREWEBVIEW2_BOUNDS_MODE) error {

	hr, _, _ := i.Vtbl.PutBoundsMode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boundsMode),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
