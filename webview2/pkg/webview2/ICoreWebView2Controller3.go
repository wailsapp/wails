//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Controller3Vtbl struct {
	IUnknownVtbl
	GetRasterizationScale ComProc
	PutRasterizationScale ComProc
	GetShouldDetectMonitorScaleChanges ComProc
	PutShouldDetectMonitorScaleChanges ComProc
	AddRasterizationScaleChanged ComProc
	RemoveRasterizationScaleChanged ComProc
	GetBoundsMode ComProc
	PutBoundsMode ComProc
}

type ICoreWebView2Controller3 struct {
	Vtbl *ICoreWebView2Controller3Vtbl
}

func (i *ICoreWebView2Controller3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Controller3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Controller3() (*ICoreWebView2Controller3, error) {
	var result *ICoreWebView2Controller3

	iidICoreWebView2Controller3 := NewGUID("{f9614724-5d2b-41dc-aef7-73d62b51543b}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Controller3) GetRasterizationScale() (float64, error) {

	var scale float64

	hr, _, err := i.Vtbl.GetRasterizationScale.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&scale)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return scale, err
}

func (i *ICoreWebView2Controller3) PutRasterizationScale(scale float64) error {


	hr, _, err := i.Vtbl.PutRasterizationScale.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(scale),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller3) GetShouldDetectMonitorScaleChanges() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetShouldDetectMonitorScaleChanges.Call(
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

func (i *ICoreWebView2Controller3) PutShouldDetectMonitorScaleChanges(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutShouldDetectMonitorScaleChanges.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller3) AddRasterizationScaleChanged(eventHandler *ICoreWebView2RasterizationScaleChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddRasterizationScaleChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Controller3) RemoveRasterizationScaleChanged(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveRasterizationScaleChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller3) GetBoundsMode() (COREWEBVIEW2_BOUNDS_MODE, error) {

	var boundsMode COREWEBVIEW2_BOUNDS_MODE

	hr, _, err := i.Vtbl.GetBoundsMode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&boundsMode)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return boundsMode, err
}

func (i *ICoreWebView2Controller3) PutBoundsMode(boundsMode COREWEBVIEW2_BOUNDS_MODE) error {


	hr, _, err := i.Vtbl.PutBoundsMode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boundsMode),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
