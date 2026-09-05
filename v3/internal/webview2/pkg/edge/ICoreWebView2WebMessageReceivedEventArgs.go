//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type iCoreWebView2WebMessageReceivedEventArgsVtbl struct {
	_IUnknownVtbl
	GetSource                ComProc
	GetWebMessageAsJSON      ComProc
	TryGetWebMessageAsString ComProc
	GetAdditionalObjects     ComProc
}

type ICoreWebView2WebMessageReceivedEventArgs struct {
	vtbl *iCoreWebView2WebMessageReceivedEventArgsVtbl
}

func (i *ICoreWebView2WebMessageReceivedEventArgs) GetSource() (string, error) {
	var _source *uint16
	hr, _, _ := i.vtbl.GetSource.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_source)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}
	source := windows.UTF16PtrToString(_source)
	windows.CoTaskMemFree(unsafe.Pointer(_source))
	return source, nil
}

func (i *ICoreWebView2WebMessageReceivedEventArgs) GetAdditionalObjects() (*ICoreWebView2ObjectCollectionView, error) {
	var value *ICoreWebView2ObjectCollectionView

	hr, _, _ := i.vtbl.GetAdditionalObjects.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2WebMessageReceivedEventArgs) TryGetWebMessageAsString() (string, error) {
	var u16msg *uint16

	hr, _, _ := i.vtbl.TryGetWebMessageAsString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&u16msg)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(u16msg))

	msg := windows.UTF16PtrToString(u16msg)

	return msg, nil
}
