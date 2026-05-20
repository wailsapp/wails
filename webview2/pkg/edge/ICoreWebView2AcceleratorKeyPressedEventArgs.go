//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2AcceleratorKeyPressedEventArgsVtbl struct {
	_IUnknownVtbl
	GetKeyEventKind      ComProc
	GetVirtualKey        ComProc
	GetKeyEventLParam    ComProc
	GetPhysicalKeyStatus ComProc
	GetHandled           ComProc
	PutHandled           ComProc
}

type ICoreWebView2AcceleratorKeyPressedEventArgs struct {
	vtbl *_ICoreWebView2AcceleratorKeyPressedEventArgsVtbl
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) GetKeyEventKind() (COREWEBVIEW2_KEY_EVENT_KIND, error) {
	var keyEventKind COREWEBVIEW2_KEY_EVENT_KIND
	hr, _, _ := i.vtbl.GetKeyEventKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&keyEventKind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, windows.Errno(hr)
	}
	return keyEventKind, nil
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) GetVirtualKey() (uint, error) {
	var virtualKey uint
	hr, _, _ := i.vtbl.GetVirtualKey.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&virtualKey)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, windows.Errno(hr)
	}
	return virtualKey, nil
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) GetPhysicalKeyStatus() (COREWEBVIEW2_PHYSICAL_KEY_STATUS, error) {
	var physicalKeyStatus internal_COREWEBVIEW2_PHYSICAL_KEY_STATUS
	hr, _, _ := i.vtbl.GetPhysicalKeyStatus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&physicalKeyStatus)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return COREWEBVIEW2_PHYSICAL_KEY_STATUS{}, windows.Errno(hr)
	}
	return COREWEBVIEW2_PHYSICAL_KEY_STATUS{
		RepeatCount:   physicalKeyStatus.RepeatCount,
		ScanCode:      physicalKeyStatus.ScanCode,
		IsExtendedKey: physicalKeyStatus.IsExtendedKey != 0,
		IsMenuKeyDown: physicalKeyStatus.IsMenuKeyDown != 0,
		WasKeyDown:    physicalKeyStatus.WasKeyDown != 0,
		IsKeyReleased: physicalKeyStatus.IsKeyReleased != 0,
	}, nil
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs) PutHandled(handled bool) error {
	hr, _, _ := i.vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(handled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}
