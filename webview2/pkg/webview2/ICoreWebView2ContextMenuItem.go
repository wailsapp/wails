//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ContextMenuItemVtbl struct {
	IUnknownVtbl
	GetName ComProc
	GetLabel ComProc
	GetCommandId ComProc
	GetShortcutKeyDescription ComProc
	GetIcon ComProc
	GetKind ComProc
	PutIsEnabled ComProc
	GetIsEnabled ComProc
	PutIsChecked ComProc
	GetIsChecked ComProc
	GetChildren ComProc
	AddCustomItemSelected ComProc
	RemoveCustomItemSelected ComProc
}

type ICoreWebView2ContextMenuItem struct {
	Vtbl *ICoreWebView2ContextMenuItemVtbl
}

func (i *ICoreWebView2ContextMenuItem) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ContextMenuItem) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2ContextMenuItem) GetName() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) GetLabel() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetLabel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) GetCommandId() (int32, error) {

	var value int32

	hr, _, err := i.Vtbl.GetCommandId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) GetShortcutKeyDescription() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetShortcutKeyDescription.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) GetIcon() (*IStream, error) {

	var value *IStream

	hr, _, err := i.Vtbl.GetIcon.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) GetKind() (COREWEBVIEW2_CONTEXT_MENU_ITEM_KIND, error) {

	var value COREWEBVIEW2_CONTEXT_MENU_ITEM_KIND

	hr, _, err := i.Vtbl.GetKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) PutIsEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutIsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2ContextMenuItem) GetIsEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetIsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) PutIsChecked(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutIsChecked.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2ContextMenuItem) GetIsChecked() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetIsChecked.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) GetChildren() (*ICoreWebView2ContextMenuItemCollection, error) {

	var value *ICoreWebView2ContextMenuItemCollection

	hr, _, err := i.Vtbl.GetChildren.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuItem) AddCustomItemSelected(eventHandler *ICoreWebView2CustomItemSelectedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddCustomItemSelected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2ContextMenuItem) RemoveCustomItemSelected(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveCustomItemSelected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
