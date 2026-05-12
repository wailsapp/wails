//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2FrameVtbl struct {
	IUnknownVtbl
	GetName ComProc
	AddNameChanged ComProc
	RemoveNameChanged ComProc
	AddHostObjectToScriptWithOrigins ComProc
	RemoveHostObjectFromScript ComProc
	AddDestroyed ComProc
	RemoveDestroyed ComProc
	IsDestroyed ComProc
}

type ICoreWebView2Frame struct {
	Vtbl *ICoreWebView2FrameVtbl
}

func (i *ICoreWebView2Frame) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Frame) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2Frame) GetName() (string, error) {
	// Create *uint16 to hold result
	var _name *uint16


	hr, _, err := i.Vtbl.GetName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_name)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	name := UTF16PtrToString(_name)
	CoTaskMemFree(unsafe.Pointer(_name))
	return name, err
}

func (i *ICoreWebView2Frame) AddNameChanged(eventHandler *ICoreWebView2FrameNameChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddNameChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Frame) RemoveNameChanged(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveNameChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame) AddHostObjectToScriptWithOrigins(name string, object *VARIANT, originsCount uint32, origins []string) error {

	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}	// Convert []string 'origins' to **uint16 (LPCWSTR* / LPWSTR*)
	_originsptrs := make([]*uint16, len(origins))
	for _i, _s := range origins {
		_p, err := UTF16PtrFromString(_s)
		if err != nil {
			return err
		}
		_originsptrs[_i] = _p
	}
	var _origins **uint16
	if len(_originsptrs) > 0 {
		_origins = &_originsptrs[0]
	}


	hr, _, err := i.Vtbl.AddHostObjectToScriptWithOrigins.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(object)),
		uintptr(originsCount),
		uintptr(unsafe.Pointer(_origins)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame) RemoveHostObjectFromScript(name string) error {

	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.RemoveHostObjectFromScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame) AddDestroyed(eventHandler *ICoreWebView2FrameDestroyedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddDestroyed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2Frame) RemoveDestroyed(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveDestroyed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Frame) IsDestroyed() (bool, error) {
	// Create int32 to hold bool result
	var _destroyed int32

	hr, _, err := i.Vtbl.IsDestroyed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_destroyed)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    destroyed := _destroyed != 0
	return destroyed, err
}
