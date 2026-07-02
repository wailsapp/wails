//go:build windows

package edge

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Controller3Vtbl struct {
	_IUnknownVtbl
	GetIsVisible                       ComProc
	PutIsVisible                       ComProc
	GetBounds                          ComProc
	PutBounds                          ComProc
	GetZoomFactor                      ComProc
	PutZoomFactor                      ComProc
	AddZoomFactorChanged               ComProc
	RemoveZoomFactorChanged            ComProc
	SetBoundsAndZoomFactor             ComProc
	MoveFocus                          ComProc
	AddMoveFocusRequested              ComProc
	RemoveMoveFocusRequested           ComProc
	AddGotFocus                        ComProc
	RemoveGotFocus                     ComProc
	AddLostFocus                       ComProc
	RemoveLostFocus                    ComProc
	AddAcceleratorKeyPressed           ComProc
	RemoveAcceleratorKeyPressed        ComProc
	GetParentWindow                    ComProc
	PutParentWindow                    ComProc
	NotifyParentWindowPositionChanged  ComProc
	Close                              ComProc
	GetCoreWebView2                    ComProc
	GetDefaultBackgroundColor          ComProc
	PutDefaultBackgroundColor          ComProc
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

func (i *ICoreWebView2Controller) GetICoreWebView2Controller3() *ICoreWebView2Controller3 {
	var result *ICoreWebView2Controller3

	iidICoreWebView2Controller3 := NewGUID("{f9614724-5d2b-41dc-aef7-73d62b51543b}")
	_, _, _ = i.vtbl.QueryInterface.Call(
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

	hr, _, _ := i.Vtbl.PutShouldDetectMonitorScaleChanges.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(value)),
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
