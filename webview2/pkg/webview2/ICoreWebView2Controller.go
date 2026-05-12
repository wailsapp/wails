//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ControllerVtbl struct {
	IUnknownVtbl
	GetIsVisible ComProc
	PutIsVisible ComProc
	GetBounds ComProc
	PutBounds ComProc
	GetZoomFactor ComProc
	PutZoomFactor ComProc
	AddZoomFactorChanged ComProc
	RemoveZoomFactorChanged ComProc
	SetBoundsAndZoomFactor ComProc
	MoveFocus ComProc
	AddMoveFocusRequested ComProc
	RemoveMoveFocusRequested ComProc
	AddGotFocus ComProc
	RemoveGotFocus ComProc
	AddLostFocus ComProc
	RemoveLostFocus ComProc
	AddAcceleratorKeyPressed ComProc
	RemoveAcceleratorKeyPressed ComProc
	GetParentWindow ComProc
	PutParentWindow ComProc
	NotifyParentWindowPositionChanged ComProc
	Close ComProc
	GetCoreWebView2 ComProc
}

type ICoreWebView2Controller struct {
	Vtbl *ICoreWebView2ControllerVtbl
}

func (i *ICoreWebView2Controller) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Controller) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2Controller) GetIsVisible() (bool, error) {
	// Create int32 to hold bool result
	var _isVisible int32

	hr, _, err := i.Vtbl.GetIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isVisible)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    isVisible := _isVisible != 0
	return isVisible, err
}

func (i *ICoreWebView2Controller) PutIsVisible(isVisible bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _isVisible int32
	if isVisible {
		_isVisible = 1
	}

	hr, _, err := i.Vtbl.PutIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_isVisible),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) GetBounds() (RECT, error) {

	var bounds RECT

	hr, _, err := i.Vtbl.GetBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return RECT{}, syscall.Errno(hr)
	}
	return bounds, err
}

func (i *ICoreWebView2Controller) PutBounds(bounds RECT) error {


	hr, _, err := i.Vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) GetZoomFactor() (float64, error) {

	var zoomFactor float64

	hr, _, err := i.Vtbl.GetZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&zoomFactor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return zoomFactor, err
}

func (i *ICoreWebView2Controller) PutZoomFactor(zoomFactor float64) error {


	hr, _, err := i.Vtbl.PutZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(zoomFactor),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) AddZoomFactorChanged(eventHandler *ICoreWebView2ZoomFactorChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddZoomFactorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Controller) RemoveZoomFactorChanged(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveZoomFactorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) SetBoundsAndZoomFactor(bounds RECT, zoomFactor float64) error {


	hr, _, err := i.Vtbl.SetBoundsAndZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
		uintptr(zoomFactor),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) MoveFocus(reason COREWEBVIEW2_MOVE_FOCUS_REASON) error {


	hr, _, err := i.Vtbl.MoveFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(reason),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) AddMoveFocusRequested(eventHandler *ICoreWebView2MoveFocusRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddMoveFocusRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Controller) RemoveMoveFocusRequested(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveMoveFocusRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) AddGotFocus(eventHandler *ICoreWebView2FocusChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddGotFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Controller) RemoveGotFocus(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveGotFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) AddLostFocus(eventHandler *ICoreWebView2FocusChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddLostFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Controller) RemoveLostFocus(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveLostFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) AddAcceleratorKeyPressed(eventHandler *ICoreWebView2AcceleratorKeyPressedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Controller) RemoveAcceleratorKeyPressed(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) GetParentWindow() (HWND, error) {

	var parentWindow HWND

	hr, _, err := i.Vtbl.GetParentWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&parentWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return parentWindow, err
}

func (i *ICoreWebView2Controller) PutParentWindow(parentWindow HWND) error {


	hr, _, err := i.Vtbl.PutParentWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&parentWindow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) NotifyParentWindowPositionChanged() error {


	hr, _, err := i.Vtbl.NotifyParentWindowPositionChanged.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) Close() error {


	hr, _, err := i.Vtbl.Close.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Controller) GetCoreWebView2() (*ICoreWebView2, error) {

	var coreWebView2 *ICoreWebView2

	hr, _, err := i.Vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&coreWebView2)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return coreWebView2, err
}
