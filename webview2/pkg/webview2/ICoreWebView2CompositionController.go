//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2CompositionControllerVtbl struct {
	IUnknownVtbl
	GetRootVisualTarget ComProc
	PutRootVisualTarget ComProc
	SendMouseInput ComProc
	SendPointerInput ComProc
	GetCursor ComProc
	GetSystemCursorId ComProc
	AddCursorChanged ComProc
	RemoveCursorChanged ComProc
}

type ICoreWebView2CompositionController struct {
	Vtbl *ICoreWebView2CompositionControllerVtbl
}

func (i *ICoreWebView2CompositionController) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CompositionController) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2CompositionController) GetRootVisualTarget() (*IUnknown, error) {

	var target *IUnknown

	hr, _, err := i.Vtbl.GetRootVisualTarget.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&target)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return target, err
}

func (i *ICoreWebView2CompositionController) PutRootVisualTarget(target *IUnknown) error {


	hr, _, err := i.Vtbl.PutRootVisualTarget.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(target)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2CompositionController) SendMouseInput(eventKind COREWEBVIEW2_MOUSE_EVENT_KIND, virtualKeys COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS, mouseData uint32, point POINT) error {


	hr, _, err := i.Vtbl.SendMouseInput.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(eventKind),
		uintptr(virtualKeys),
		uintptr(mouseData),
		uintptr(unsafe.Pointer(&point)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2CompositionController) SendPointerInput(eventKind COREWEBVIEW2_POINTER_EVENT_KIND, pointerInfo *ICoreWebView2PointerInfo) error {


	hr, _, err := i.Vtbl.SendPointerInput.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(eventKind),
		uintptr(unsafe.Pointer(pointerInfo)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2CompositionController) GetCursor() (HCURSOR, error) {

	var cursor HCURSOR

	hr, _, err := i.Vtbl.GetCursor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&cursor)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return cursor, err
}

func (i *ICoreWebView2CompositionController) GetSystemCursorId() (uint32, error) {

	var systemCursorId uint32

	hr, _, err := i.Vtbl.GetSystemCursorId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&systemCursorId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return systemCursorId, err
}

func (i *ICoreWebView2CompositionController) AddCursorChanged(eventHandler *ICoreWebView2CursorChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddCursorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2CompositionController) RemoveCursorChanged(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveCursorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
