//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2SettingsVtbl struct {
	IUnknownVtbl
	GetIsScriptEnabled ComProc
	PutIsScriptEnabled ComProc
	GetIsWebMessageEnabled ComProc
	PutIsWebMessageEnabled ComProc
	GetAreDefaultScriptDialogsEnabled ComProc
	PutAreDefaultScriptDialogsEnabled ComProc
	GetIsStatusBarEnabled ComProc
	PutIsStatusBarEnabled ComProc
	GetAreDevToolsEnabled ComProc
	PutAreDevToolsEnabled ComProc
	GetAreDefaultContextMenusEnabled ComProc
	PutAreDefaultContextMenusEnabled ComProc
	GetAreHostObjectsAllowed ComProc
	PutAreHostObjectsAllowed ComProc
	GetIsZoomControlEnabled ComProc
	PutIsZoomControlEnabled ComProc
	GetIsBuiltInErrorPageEnabled ComProc
	PutIsBuiltInErrorPageEnabled ComProc
}

type ICoreWebView2Settings struct {
	Vtbl *ICoreWebView2SettingsVtbl
}

func (i *ICoreWebView2Settings) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Settings) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2Settings) GetIsScriptEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _isScriptEnabled int32

	hr, _, _ := i.Vtbl.GetIsScriptEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isScriptEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    isScriptEnabled := _isScriptEnabled != 0
	return isScriptEnabled, nil
}

func (i *ICoreWebView2Settings) PutIsScriptEnabled(isScriptEnabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _isScriptEnabled int32
	if isScriptEnabled {
		_isScriptEnabled = 1
	}

	hr, _, _ := i.Vtbl.PutIsScriptEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_isScriptEnabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsWebMessageEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _isWebMessageEnabled int32

	hr, _, _ := i.Vtbl.GetIsWebMessageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isWebMessageEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    isWebMessageEnabled := _isWebMessageEnabled != 0
	return isWebMessageEnabled, nil
}

func (i *ICoreWebView2Settings) PutIsWebMessageEnabled(isWebMessageEnabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _isWebMessageEnabled int32
	if isWebMessageEnabled {
		_isWebMessageEnabled = 1
	}

	hr, _, _ := i.Vtbl.PutIsWebMessageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_isWebMessageEnabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreDefaultScriptDialogsEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _areDefaultScriptDialogsEnabled int32

	hr, _, _ := i.Vtbl.GetAreDefaultScriptDialogsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_areDefaultScriptDialogsEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    areDefaultScriptDialogsEnabled := _areDefaultScriptDialogsEnabled != 0
	return areDefaultScriptDialogsEnabled, nil
}

func (i *ICoreWebView2Settings) PutAreDefaultScriptDialogsEnabled(areDefaultScriptDialogsEnabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _areDefaultScriptDialogsEnabled int32
	if areDefaultScriptDialogsEnabled {
		_areDefaultScriptDialogsEnabled = 1
	}

	hr, _, _ := i.Vtbl.PutAreDefaultScriptDialogsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_areDefaultScriptDialogsEnabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsStatusBarEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _isStatusBarEnabled int32

	hr, _, _ := i.Vtbl.GetIsStatusBarEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isStatusBarEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    isStatusBarEnabled := _isStatusBarEnabled != 0
	return isStatusBarEnabled, nil
}

func (i *ICoreWebView2Settings) PutIsStatusBarEnabled(isStatusBarEnabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _isStatusBarEnabled int32
	if isStatusBarEnabled {
		_isStatusBarEnabled = 1
	}

	hr, _, _ := i.Vtbl.PutIsStatusBarEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_isStatusBarEnabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreDevToolsEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _areDevToolsEnabled int32

	hr, _, _ := i.Vtbl.GetAreDevToolsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_areDevToolsEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    areDevToolsEnabled := _areDevToolsEnabled != 0
	return areDevToolsEnabled, nil
}

func (i *ICoreWebView2Settings) PutAreDevToolsEnabled(areDevToolsEnabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _areDevToolsEnabled int32
	if areDevToolsEnabled {
		_areDevToolsEnabled = 1
	}

	hr, _, _ := i.Vtbl.PutAreDevToolsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_areDevToolsEnabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreDefaultContextMenusEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _enabled int32

	hr, _, _ := i.Vtbl.GetAreDefaultContextMenusEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    enabled := _enabled != 0
	return enabled, nil
}

func (i *ICoreWebView2Settings) PutAreDefaultContextMenusEnabled(enabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _enabled int32
	if enabled {
		_enabled = 1
	}

	hr, _, _ := i.Vtbl.PutAreDefaultContextMenusEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_enabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreHostObjectsAllowed() (bool, error) {
	// Create int32 to hold bool result
	var _allowed int32

	hr, _, _ := i.Vtbl.GetAreHostObjectsAllowed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_allowed)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    allowed := _allowed != 0
	return allowed, nil
}

func (i *ICoreWebView2Settings) PutAreHostObjectsAllowed(allowed bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _allowed int32
	if allowed {
		_allowed = 1
	}

	hr, _, _ := i.Vtbl.PutAreHostObjectsAllowed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_allowed),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsZoomControlEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _enabled int32

	hr, _, _ := i.Vtbl.GetIsZoomControlEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    enabled := _enabled != 0
	return enabled, nil
}

func (i *ICoreWebView2Settings) PutIsZoomControlEnabled(enabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _enabled int32
	if enabled {
		_enabled = 1
	}

	hr, _, _ := i.Vtbl.PutIsZoomControlEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_enabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsBuiltInErrorPageEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _enabled int32

	hr, _, _ := i.Vtbl.GetIsBuiltInErrorPageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    enabled := _enabled != 0
	return enabled, nil
}

func (i *ICoreWebView2Settings) PutIsBuiltInErrorPageEnabled(enabled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _enabled int32
	if enabled {
		_enabled = 1
	}

	hr, _, _ := i.Vtbl.PutIsBuiltInErrorPageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_enabled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
