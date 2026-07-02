//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ControllerVtbl struct {
	IUnknownVtbl
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
	Vtbl *ICoreWebView2ControllerVtbl
}

func (i *ICoreWebView2Controller) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2Controller) GetIsVisible() (bool, error) {
	// Create int32 to hold bool result
	var _isVisible int32

	hr, _, _ := i.Vtbl.GetIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isVisible)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isVisible := _isVisible != 0
	return isVisible, nil
}

func (i *ICoreWebView2Controller) PutIsVisible(isVisible bool) error {

	hr, _, _ := i.Vtbl.PutIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isVisible)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) GetBounds() (RECT, error) {

	var bounds RECT

	hr, _, _ := i.Vtbl.GetBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return RECT{}, syscall.Errno(hr)
	}
	return bounds, nil
}

func (i *ICoreWebView2Controller) PutBounds(bounds RECT) error {

	hr, _, _ := i.Vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) GetZoomFactor() (float64, error) {

	var zoomFactor float64

	hr, _, _ := i.Vtbl.GetZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&zoomFactor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return zoomFactor, nil
}

func (i *ICoreWebView2Controller) PutZoomFactor(zoomFactor float64) error {

	hr, _, _ := i.Vtbl.PutZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&zoomFactor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) AddZoomFactorChanged(eventHandler *ICoreWebView2ZoomFactorChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddZoomFactorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Controller) RemoveZoomFactorChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveZoomFactorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) SetBoundsAndZoomFactor(bounds RECT, zoomFactor float64) error {

	hr, _, _ := i.Vtbl.SetBoundsAndZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
		uintptr(unsafe.Pointer(&zoomFactor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) MoveFocus(reason COREWEBVIEW2_MOVE_FOCUS_REASON) error {

	hr, _, _ := i.Vtbl.MoveFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(reason),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) AddMoveFocusRequested(eventHandler *ICoreWebView2MoveFocusRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddMoveFocusRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Controller) RemoveMoveFocusRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveMoveFocusRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) AddGotFocus(eventHandler *ICoreWebView2FocusChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddGotFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Controller) RemoveGotFocus(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveGotFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) AddLostFocus(eventHandler *ICoreWebView2FocusChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddLostFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Controller) RemoveLostFocus(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveLostFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) AddAcceleratorKeyPressed(eventHandler *ICoreWebView2AcceleratorKeyPressedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Controller) RemoveAcceleratorKeyPressed(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) GetParentWindow() (HWND, error) {

	var parentWindow HWND

	hr, _, _ := i.Vtbl.GetParentWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&parentWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return parentWindow, nil
}

func (i *ICoreWebView2Controller) PutParentWindow(parentWindow HWND) error {

	hr, _, _ := i.Vtbl.PutParentWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&parentWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) NotifyParentWindowPositionChanged() error {

	hr, _, _ := i.Vtbl.NotifyParentWindowPositionChanged.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) Close() error {

	hr, _, _ := i.Vtbl.Close.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Controller) GetCoreWebView2() (*ICoreWebView2, error) {

	var coreWebView2 *ICoreWebView2

	hr, _, _ := i.Vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&coreWebView2)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return coreWebView2, nil
}
