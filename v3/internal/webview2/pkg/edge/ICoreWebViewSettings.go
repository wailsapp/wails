//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// ICoreWebviewSettings is the merged settings class

type _ICoreWebViewSettingsVtbl struct {
	_IUnknownVtbl
	GetIsScriptEnabled                  ComProc
	PutIsScriptEnabled                  ComProc
	GetIsWebMessageEnabled              ComProc
	PutIsWebMessageEnabled              ComProc
	GetAreDefaultScriptDialogsEnabled   ComProc
	PutAreDefaultScriptDialogsEnabled   ComProc
	GetIsStatusBarEnabled               ComProc
	PutIsStatusBarEnabled               ComProc
	GetAreDevToolsEnabled               ComProc
	PutAreDevToolsEnabled               ComProc
	GetAreDefaultContextMenusEnabled    ComProc
	PutAreDefaultContextMenusEnabled    ComProc
	GetAreHostObjectsAllowed            ComProc
	PutAreHostObjectsAllowed            ComProc
	GetIsZoomControlEnabled             ComProc
	PutIsZoomControlEnabled             ComProc
	GetIsBuiltInErrorPageEnabled        ComProc
	PutIsBuiltInErrorPageEnabled        ComProc
	GetUserAgent                        ComProc // ICoreWebView2Settings2: SDK 1.0.864.35
	PutUserAgent                        ComProc
	GetAreBrowserAcceleratorKeysEnabled ComProc // ICoreWebView2Settings3: SDK 1.0.864.35
	PutAreBrowserAcceleratorKeysEnabled ComProc
	GetIsPasswordAutosaveEnabled        ComProc // ICoreWebView2Settings4: SDK 1.0.902.49
	PutIsPasswordAutosaveEnabled        ComProc
	GetIsGeneralAutofillEnabled         ComProc
	PutIsGeneralAutofillEnabled         ComProc
	GetIsPinchZoomEnabled               ComProc // ICoreWebView2Settings5: SDK 1.0.902.49
	PutIsPinchZoomEnabled               ComProc
	GetIsSwipeNavigationEnabled         ComProc // ICoreWebView2Settings6: SDK 1.0.992.28
	PutIsSwipeNavigationEnabled         ComProc
}

type ICoreWebViewSettings struct {
	vtbl *_ICoreWebViewSettingsVtbl
}

func (i *ICoreWebViewSettings) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebViewSettings) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebViewSettings) GetIsScriptEnabled() (bool, error) {
	var isScriptEnabled bool
	hr, _, _ := i.vtbl.GetIsScriptEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isScriptEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return isScriptEnabled, nil
}

func (i *ICoreWebViewSettings) PutIsScriptEnabled(isScriptEnabled bool) error {
	hr, _, _ := i.vtbl.PutIsScriptEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isScriptEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetIsWebMessageEnabled() (bool, error) {
	var isWebMessageEnabled bool
	hr, _, _ := i.vtbl.GetIsWebMessageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isWebMessageEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return isWebMessageEnabled, nil
}

func (i *ICoreWebViewSettings) PutIsWebMessageEnabled(isWebMessageEnabled bool) error {
	hr, _, _ := i.vtbl.PutIsWebMessageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isWebMessageEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetAreDefaultScriptDialogsEnabled() (bool, error) {
	var areDefaultScriptDialogsEnabled bool
	hr, _, _ := i.vtbl.GetAreDefaultScriptDialogsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&areDefaultScriptDialogsEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return areDefaultScriptDialogsEnabled, nil
}

func (i *ICoreWebViewSettings) PutAreDefaultScriptDialogsEnabled(areDefaultScriptDialogsEnabled bool) error {
	hr, _, _ := i.vtbl.PutAreDefaultScriptDialogsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(areDefaultScriptDialogsEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetIsStatusBarEnabled() (bool, error) {
	var isStatusBarEnabled bool
	hr, _, _ := i.vtbl.GetIsStatusBarEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isStatusBarEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return isStatusBarEnabled, nil
}

func (i *ICoreWebViewSettings) PutIsStatusBarEnabled(isStatusBarEnabled bool) error {
	hr, _, _ := i.vtbl.PutIsStatusBarEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isStatusBarEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetAreDevToolsEnabled() (bool, error) {
	var areDevToolsEnabled bool
	hr, _, _ := i.vtbl.GetAreDevToolsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&areDevToolsEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return areDevToolsEnabled, nil
}

func (i *ICoreWebViewSettings) PutAreDevToolsEnabled(areDevToolsEnabled bool) error {
	hr, _, _ := i.vtbl.PutAreDevToolsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(areDevToolsEnabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetAreDefaultContextMenusEnabled() (bool, error) {
	var enabled bool
	hr, _, _ := i.vtbl.GetAreDefaultContextMenusEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return enabled, nil
}

func (i *ICoreWebViewSettings) PutAreDefaultContextMenusEnabled(enabled bool) error {
	hr, _, _ := i.vtbl.PutAreDefaultContextMenusEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetAreHostObjectsAllowed() (bool, error) {
	var allowed bool
	hr, _, _ := i.vtbl.GetAreHostObjectsAllowed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&allowed)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return allowed, nil
}

func (i *ICoreWebViewSettings) PutAreHostObjectsAllowed(allowed bool) error {
	hr, _, _ := i.vtbl.PutAreHostObjectsAllowed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(allowed)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetIsZoomControlEnabled() (bool, error) {
	var enabled bool
	hr, _, _ := i.vtbl.GetIsZoomControlEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return enabled, nil
}

func (i *ICoreWebViewSettings) PutIsZoomControlEnabled(enabled bool) error {
	hr, _, _ := i.vtbl.PutIsZoomControlEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetIsBuiltInErrorPageEnabled() (bool, error) {
	var enabled bool
	hr, _, _ := i.vtbl.GetIsBuiltInErrorPageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return enabled, nil
}

func (i *ICoreWebViewSettings) PutIsBuiltInErrorPageEnabled(enabled bool) error {
	hr, _, _ := i.vtbl.PutIsBuiltInErrorPageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetUserAgent() (string, error) {
	// Create *uint16 to hold result
	var _userAgent *uint16
	hr, _, _ := i.vtbl.GetUserAgent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_userAgent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	} // Get result and cleanup
	userAgent := windows.UTF16PtrToString(_userAgent)
	windows.CoTaskMemFree(unsafe.Pointer(_userAgent))
	return userAgent, nil
}

func (i *ICoreWebViewSettings) PutUserAgent(userAgent string) error {
	
	// Convert string 'userAgent' to *uint16
	_userAgent, err := windows.UTF16PtrFromString(userAgent)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.PutUserAgent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_userAgent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetAreBrowserAcceleratorKeysEnabled() (bool, error) {
	var enabled bool
	hr, _, _ := i.vtbl.GetAreBrowserAcceleratorKeysEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return enabled, nil
}

func (i *ICoreWebViewSettings) PutAreBrowserAcceleratorKeysEnabled(enabled bool) error {
	hr, _, _ := i.vtbl.PutAreBrowserAcceleratorKeysEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetIsPinchZoomEnabled() (bool, error) {
	var enabled bool
	hr, _, _ := i.vtbl.GetIsPinchZoomEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return enabled, nil
}

func (i *ICoreWebViewSettings) PutIsPinchZoomEnabled(enabled bool) error {
	hr, _, _ := i.vtbl.PutIsPinchZoomEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}

func (i *ICoreWebViewSettings) GetIsSwipeNavigationEnabled() (bool, error) {
	var enabled bool
	hr, _, _ := i.vtbl.GetIsSwipeNavigationEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return enabled, nil
}

func (i *ICoreWebViewSettings) PutIsSwipeNavigationEnabled(enabled bool) error {
	hr, _, _ := i.vtbl.PutIsSwipeNavigationEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}
