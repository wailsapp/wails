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

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.SendMouseInput.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(eventKind),
			uintptr(virtualKeys),
			uintptr(mouseData),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.SendMouseInput.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(eventKind),
			uintptr(virtualKeys),
			uintptr(mouseData),
			uintptr(*(*uint64)(unsafe.Pointer(&point))),
		)
	}
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

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemoveCursorChanged.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemoveCursorChanged.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
