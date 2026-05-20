//go:build windows

package edge

import (
	"math"
	"unsafe"

	"github.com/wailsapp/wails/webview2/internal/w32"
	"golang.org/x/sys/windows"
)

type _ICoreWebView2ControllerVtbl struct {
	_IUnknownVtbl
	GetIsVisible                      ComProc
	PutIsVisible                      ComProc
	GetBounds                         ComProc
	PutBounds                         ComProc
	GetZoomFactor                     ComProc
	PutZoomFactor                     ComProc
	AddZoomFactorChanged              ComProc
	RemoveZoomFactorChanged           ComProc
	SetBoundsAndZoomFactor            ComProc
	MoveFocus                         ComProc
	AddMoveFocusRequested             ComProc
	RemoveMoveFocusRequested          ComProc
	AddGotFocus                       ComProc
	RemoveGotFocus                    ComProc
	AddLostFocus                      ComProc
	RemoveLostFocus                   ComProc
	AddAcceleratorKeyPressed          ComProc
	RemoveAcceleratorKeyPressed       ComProc
	GetParentWindow                   ComProc
	PutParentWindow                   ComProc
	NotifyParentWindowPositionChanged ComProc
	Close                             ComProc
	GetCoreWebView2                   ComProc
}

type ICoreWebView2Controller struct {
	vtbl *_ICoreWebView2ControllerVtbl
}

func (i *ICoreWebView2Controller) AddRef() uintptr {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2Controller) Release() uintptr {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2Controller) GetCoreWebView2() (*ICoreWebView2, error) {
	var wv2Ptr *ICoreWebView2
	hr, _, _ := i.vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&wv2Ptr)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}

	return wv2Ptr, nil
}

func (i *ICoreWebView2Controller) GetBounds() (*w32.Rect, error) {
	var bounds w32.Rect
	hr, _, _ := i.vtbl.GetBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}
	return &bounds, nil
}

func (i *ICoreWebView2Controller) PutBounds(bounds w32.Rect) error {
	hr, _, _ := i.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) MoveFocus(reason COREWEBVIEW2_MOVE_FOCUS_REASON) error {

	hr, _, _ := i.vtbl.MoveFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(reason),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) AddAcceleratorKeyPressed(eventHandler *ICoreWebView2AcceleratorKeyPressedEventHandler, token *_EventRegistrationToken) error {

	hr, _, _ := i.vtbl.AddAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) PutIsVisible(isVisible bool) error {

	hr, _, _ := i.vtbl.PutIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isVisible)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) GetICoreWebView2Controller2() *ICoreWebView2Controller2 {

	var result *ICoreWebView2Controller2

	iidICoreWebView2Controller2 := NewGUID("{c979903e-d4ca-4228-92eb-47ee3fa96eab}")
	i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Controller) NotifyParentWindowPositionChanged() error {

	hr, _, _ := i.vtbl.NotifyParentWindowPositionChanged.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) PutZoomFactor(zoomFactor float64) error {

	hr, _, _ := i.vtbl.PutZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(math.Float64bits(zoomFactor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) GetZoomFactor() (float64, error) {

	var zoomFactorUint64 uint64
	hr, _, _ := i.vtbl.GetZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&zoomFactorUint64)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, windows.Errno(hr)
	}
	return math.Float64frombits(zoomFactorUint64), nil
}
