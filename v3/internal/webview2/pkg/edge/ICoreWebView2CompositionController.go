//go:build windows

package edge

import (
	"syscall"
	"unsafe"
)

type ICoreWebView2CompositionControllerVtbl struct {
	_IUnknownVtbl
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
	ret, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2CompositionController) Release() uintptr {
	ret, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2CompositionController) GetICoreWebView2Controller() *ICoreWebView2Controller {
	var result *ICoreWebView2Controller

	iidICoreWebView2Controller := NewGUID("{4D00C0D1-9434-4EB6-8078-8697A560334F}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2CompositionController) PutRootVisualTarget(target *IUnknown) error {
	hr, _, _ := i.Vtbl.PutRootVisualTarget.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(target)),
	)
	if int32(hr) < 0 {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CompositionController) SendMouseInput(eventKind COREWEBVIEW2_MOUSE_EVENT_KIND, virtualKeys COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS, mouseData uint32, point POINT) error {
	hr, _, _ := i.Vtbl.SendMouseInput.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(eventKind),
		uintptr(virtualKeys),
		uintptr(mouseData),
		point.uintptr(),
	)
	if int32(hr) < 0 {
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
	if int32(hr) < 0 {
		return 0, syscall.Errno(hr)
	}
	return cursor, nil
}

func (i *ICoreWebView2CompositionController) GetSystemCursorId() (uint32, error) {
	var cursorID uint32
	hr, _, _ := i.Vtbl.GetSystemCursorId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&cursorID)),
	)
	if int32(hr) < 0 {
		return 0, syscall.Errno(hr)
	}
	return cursorID, nil
}

func (i *ICoreWebView2CompositionController) AddCursorChanged(eventHandler *iCoreWebView2CursorChangedEventHandler, token *_EventRegistrationToken) error {
	hr, _, _ := i.Vtbl.AddCursorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(token)),
	)
	if int32(hr) < 0 {
		return syscall.Errno(hr)
	}
	return nil
}
