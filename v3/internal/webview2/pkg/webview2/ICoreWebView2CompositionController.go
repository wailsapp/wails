//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2CompositionControllerVtbl struct {
	IUnknownVtbl
	GetRootVisualTarget ComProc
	PutRootVisualTarget ComProc
	SendMouseInput      ComProc
	SendPointerInput    ComProc
	GetCursor           ComProc
	GetSystemCursorId   ComProc
	AddCursorChanged    ComProc
	RemoveCursorChanged ComProc
}

type ICoreWebView2CompositionController struct {
	Vtbl *ICoreWebView2CompositionControllerVtbl
}

func (i *ICoreWebView2CompositionController) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2CompositionController) GetRootVisualTarget() (*IUnknown, error) {

	var target *IUnknown

	hr, _, _ := i.Vtbl.GetRootVisualTarget.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&target)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return target, nil
}

func (i *ICoreWebView2CompositionController) PutRootVisualTarget(target *IUnknown) error {

	hr, _, _ := i.Vtbl.PutRootVisualTarget.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(target)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CompositionController) SendMouseInput(eventKind COREWEBVIEW2_MOUSE_EVENT_KIND, virtualKeys COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS, mouseData uint32, point POINT) error {

	hr, _, _ := i.Vtbl.SendMouseInput.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(eventKind),
		uintptr(virtualKeys),
		uintptr(unsafe.Pointer(&mouseData)),
		uintptr(unsafe.Pointer(&point)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CompositionController) SendPointerInput(eventKind COREWEBVIEW2_POINTER_EVENT_KIND, pointerInfo *ICoreWebView2PointerInfo) error {

	hr, _, _ := i.Vtbl.SendPointerInput.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(eventKind),
		uintptr(unsafe.Pointer(pointerInfo)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CompositionController) GetCursor() (HCURSOR, error) {

	var cursor HCURSOR

	hr, _, _ := i.Vtbl.GetCursor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&cursor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return cursor, nil
}

func (i *ICoreWebView2CompositionController) GetSystemCursorId() (uint32, error) {

	var systemCursorId uint32

	hr, _, _ := i.Vtbl.GetSystemCursorId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&systemCursorId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return systemCursorId, nil
}

func (i *ICoreWebView2CompositionController) AddCursorChanged(eventHandler *ICoreWebView2CursorChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddCursorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2CompositionController) RemoveCursorChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveCursorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
