//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ContextMenuItemVtbl struct {
	IUnknownVtbl
	GetName                   ComProc
	GetLabel                  ComProc
	GetCommandId              ComProc
	GetShortcutKeyDescription ComProc
	GetIcon                   ComProc
	GetKind                   ComProc
	PutIsEnabled              ComProc
	GetIsEnabled              ComProc
	PutIsChecked              ComProc
	GetIsChecked              ComProc
	GetChildren               ComProc
	AddCustomItemSelected     ComProc
	RemoveCustomItemSelected  ComProc
}

type ICoreWebView2ContextMenuItem struct {
	Vtbl *ICoreWebView2ContextMenuItemVtbl
}

func (i *ICoreWebView2ContextMenuItem) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ContextMenuItem) GetName() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) GetLabel() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetLabel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) GetCommandId() (int32, error) {

	var value int32

	hr, _, _ := i.Vtbl.GetCommandId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) GetShortcutKeyDescription() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetShortcutKeyDescription.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) GetIcon() (*IStream, error) {

	var value *IStream

	hr, _, _ := i.Vtbl.GetIcon.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) GetKind() (COREWEBVIEW2_CONTEXT_MENU_ITEM_KIND, error) {

	var value COREWEBVIEW2_CONTEXT_MENU_ITEM_KIND

	hr, _, _ := i.Vtbl.GetKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) PutIsEnabled(value bool) error {

	hr, _, _ := i.Vtbl.PutIsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ContextMenuItem) GetIsEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) PutIsChecked(value bool) error {

	hr, _, _ := i.Vtbl.PutIsChecked.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ContextMenuItem) GetIsChecked() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsChecked.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) GetChildren() (*ICoreWebView2ContextMenuItemCollection, error) {

	var value *ICoreWebView2ContextMenuItemCollection

	hr, _, _ := i.Vtbl.GetChildren.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ContextMenuItem) AddCustomItemSelected(eventHandler *ICoreWebView2CustomItemSelectedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddCustomItemSelected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2ContextMenuItem) RemoveCustomItemSelected(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveCustomItemSelected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
