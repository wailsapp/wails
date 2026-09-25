//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Frame5Vtbl struct {
	ICoreWebView2Frame4Vtbl
	GetFrameId ComProc
}

type ICoreWebView2Frame5 struct {
	Vtbl *ICoreWebView2Frame5Vtbl
}

func (i *ICoreWebView2Frame5) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Frame5) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Frame5 queries the object for its ICoreWebView2Frame5 interface. The receiver
// is the root of ICoreWebView2Frame5's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Frame) GetICoreWebView2Frame5() (*ICoreWebView2Frame5, error) {
	var result *ICoreWebView2Frame5

	iidICoreWebView2Frame5 := NewGUID("{99d199c4-7305-11ee-b962-0242ac120002}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame5)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Frame5) GetFrameId() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetFrameId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
