//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2SettingsVtbl struct {
	_IUnknownVtbl
	GetIsScriptEnabled                ComProc
	PutIsScriptEnabled                ComProc
	GetIsWebMessageEnabled            ComProc
	PutIsWebMessageEnabled            ComProc
	GetAreDefaultScriptDialogsEnabled ComProc
	PutAreDefaultScriptDialogsEnabled ComProc
	GetIsStatusBarEnabled             ComProc
	PutIsStatusBarEnabled             ComProc
	GetAreDevToolsEnabled             ComProc
	PutAreDevToolsEnabled             ComProc
	GetAreDefaultContextMenusEnabled  ComProc
	PutAreDefaultContextMenusEnabled  ComProc
	GetAreHostObjectsAllowed          ComProc
	PutAreHostObjectsAllowed          ComProc
	GetIsZoomControlEnabled           ComProc
	PutIsZoomControlEnabled           ComProc
	GetIsBuiltInErrorPageEnabled      ComProc
	PutIsBuiltInErrorPageEnabled      ComProc
}

type ICoreWebView2Settings struct {
	vtbl *_ICoreWebView2SettingsVtbl
}

func (i *ICoreWebView2Settings) AddRef() uintptr {
	return i.AddRef()
}

func (i *ICoreWebView2Settings) GetIsScriptEnabled() (bool, error) {
	var err error
	var isScriptEnabled bool
	_, _, err = i.vtbl.GetIsScriptEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isScriptEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return isScriptEnabled, nil
}

func (i *ICoreWebView2Settings) PutIsScriptEnabled(isScriptEnabled bool) error {
	var err error

	_, _, err = i.vtbl.PutIsScriptEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isScriptEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsWebMessageEnabled() (bool, error) {
	var err error
	var isWebMessageEnabled bool
	_, _, err = i.vtbl.GetIsWebMessageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isWebMessageEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return isWebMessageEnabled, nil
}

func (i *ICoreWebView2Settings) PutIsWebMessageEnabled(isWebMessageEnabled bool) error {
	var err error

	_, _, err = i.vtbl.PutIsWebMessageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isWebMessageEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreDefaultScriptDialogsEnabled() (bool, error) {
	var err error
	var areDefaultScriptDialogsEnabled bool
	_, _, err = i.vtbl.GetAreDefaultScriptDialogsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&areDefaultScriptDialogsEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return areDefaultScriptDialogsEnabled, nil
}

func (i *ICoreWebView2Settings) PutAreDefaultScriptDialogsEnabled(areDefaultScriptDialogsEnabled bool) error {
	var err error

	_, _, err = i.vtbl.PutAreDefaultScriptDialogsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(areDefaultScriptDialogsEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsStatusBarEnabled() (bool, error) {
	var err error
	var isStatusBarEnabled bool
	_, _, err = i.vtbl.GetIsStatusBarEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isStatusBarEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return isStatusBarEnabled, nil
}

func (i *ICoreWebView2Settings) PutIsStatusBarEnabled(isStatusBarEnabled bool) error {
	var err error

	_, _, err = i.vtbl.PutIsStatusBarEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(isStatusBarEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreDevToolsEnabled() (bool, error) {
	var err error
	var areDevToolsEnabled bool
	_, _, err = i.vtbl.GetAreDevToolsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&areDevToolsEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return areDevToolsEnabled, nil
}

func (i *ICoreWebView2Settings) PutAreDevToolsEnabled(areDevToolsEnabled bool) error {
	var err error
	_, _, err = i.vtbl.PutAreDevToolsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(areDevToolsEnabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreDefaultContextMenusEnabled() (bool, error) {
	var err error
	var enabled bool
	_, _, err = i.vtbl.GetAreDefaultContextMenusEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return enabled, nil
}

func (i *ICoreWebView2Settings) PutAreDefaultContextMenusEnabled(enabled bool) error {
	var err error
	_, _, err = i.vtbl.PutAreDefaultContextMenusEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetAreHostObjectsAllowed() (bool, error) {
	var err error
	var allowed bool
	_, _, err = i.vtbl.GetAreHostObjectsAllowed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&allowed)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return allowed, nil
}

func (i *ICoreWebView2Settings) PutAreHostObjectsAllowed(allowed bool) error {
	var err error

	_, _, err = i.vtbl.PutAreHostObjectsAllowed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(allowed)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsZoomControlEnabled() (bool, error) {
	var err error
	var enabled bool
	_, _, err = i.vtbl.GetIsZoomControlEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return enabled, nil
}

func (i *ICoreWebView2Settings) PutIsZoomControlEnabled(enabled bool) error {
	var err error

	_, _, err = i.vtbl.PutIsZoomControlEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Settings) GetIsBuiltInErrorPageEnabled() (bool, error) {
	var err error
	var enabled bool
	_, _, err = i.vtbl.GetIsBuiltInErrorPageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&enabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return enabled, nil
}

func (i *ICoreWebView2Settings) PutIsBuiltInErrorPageEnabled(enabled bool) error {
	var err error

	_, _, err = i.vtbl.PutIsBuiltInErrorPageEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(enabled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
