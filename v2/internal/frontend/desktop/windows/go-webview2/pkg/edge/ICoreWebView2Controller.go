//go:build windows

package edge

import (
	"math"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/internal/w32"
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
	return i.AddRef()
}

func (i *ICoreWebView2Controller) GetBounds() (*w32.Rect, error) {
	var err error
	var bounds w32.Rect
	_, _, err = i.vtbl.GetBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return &bounds, nil
}

func (i *ICoreWebView2Controller) PutBounds(bounds w32.Rect) error {
	var err error

	_, _, err = i.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) MoveFocus(reason COREWEBVIEW2_MOVE_FOCUS_REASON) error {
	var err error

	_, _, err = i.vtbl.MoveFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(reason),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) AddAcceleratorKeyPressed(eventHandler *ICoreWebView2AcceleratorKeyPressedEventHandler, token *_EventRegistrationToken) error {
	var err error
	_, _, err = i.vtbl.AddAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) PutIsVisible(isVisible bool) error {
	var err error

	_, _, err = i.vtbl.PutIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isVisible)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
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
	var err error
	_, _, err = i.vtbl.NotifyParentWindowPositionChanged.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) PutZoomFactor(zoomFactor float64) error {
	var err error
	_, _, err = i.vtbl.PutZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(math.Float64bits(zoomFactor)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Controller) GetZoomFactor() (float64, error) {
	var err error
	var zoomFactorUint64 uint64
	_, _, err = i.vtbl.GetZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&zoomFactorUint64)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0.0, err
	}
	return math.Float64frombits(zoomFactorUint64), nil
}
