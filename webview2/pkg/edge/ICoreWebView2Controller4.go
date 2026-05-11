//go:build windows

package edge

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Controller4Vtbl struct {
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
	GetAllowExternalDrop               ComProc
	PutAllowExternalDrop               ComProc
}

type ICoreWebView2Controller4 struct {
	Vtbl *ICoreWebView2Controller4Vtbl
}

func (i *ICoreWebView2Controller4) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2Controller) GetICoreWebView2Controller4() *ICoreWebView2Controller4 {
	var result *ICoreWebView2Controller4

	iidICoreWebView2Controller4 := NewGUID("{97d418d5-a426-4e49-a151-e1a10f327d9e}")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller4)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Controller4) GetAllowExternalDrop() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetAllowExternalDrop.Call(
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

func (i *ICoreWebView2Controller4) PutAllowExternalDrop(value bool) error {

	hr, _, _ := i.Vtbl.PutAllowExternalDrop.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
