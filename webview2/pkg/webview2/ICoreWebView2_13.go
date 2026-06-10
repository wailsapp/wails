//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_13Vtbl struct {
	ICoreWebView2_12Vtbl
	GetProfile ComProc
}

type ICoreWebView2_13 struct {
	Vtbl *ICoreWebView2_13Vtbl
}

func (i *ICoreWebView2_13) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_13) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_13 queries the object for its ICoreWebView2_13 interface. The receiver
// is the root of ICoreWebView2_13's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_13() (*ICoreWebView2_13, error) {
	var result *ICoreWebView2_13

	iidICoreWebView2_13 := NewGUID("{f75f09a8-667e-4983-88d6-c8773f315e84}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_13)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_13) GetProfile() (*ICoreWebView2Profile, error) {

	var value *ICoreWebView2Profile

	hr, _, _ := i.Vtbl.GetProfile.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
