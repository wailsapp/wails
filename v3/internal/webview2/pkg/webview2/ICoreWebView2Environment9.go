//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment9Vtbl struct {
	IUnknownVtbl
	CreateContextMenuItem ComProc
}

type ICoreWebView2Environment9 struct {
	Vtbl *ICoreWebView2Environment9Vtbl
}

func (i *ICoreWebView2Environment9) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment9() *ICoreWebView2Environment9 {
	var result *ICoreWebView2Environment9

	iidICoreWebView2Environment9 := NewGUID("{f06f41bf-4b5a-49d8-b9f6-fa16cd29f274}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment9)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment9) CreateContextMenuItem(Label string, iconStream *IStream, Kind COREWEBVIEW2_CONTEXT_MENU_ITEM_KIND) (*ICoreWebView2ContextMenuItem, error) {

	// Convert string 'Label' to *uint16
	_Label, err := UTF16PtrFromString(Label)
	if err != nil {
		return nil, err
	}
	var value *ICoreWebView2ContextMenuItem

	hr, _, _ := i.Vtbl.CreateContextMenuItem.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_Label)),
		uintptr(unsafe.Pointer(iconStream)),
		uintptr(Kind),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
