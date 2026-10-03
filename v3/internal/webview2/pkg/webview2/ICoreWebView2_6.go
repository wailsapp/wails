//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_6Vtbl struct {
	ICoreWebView2_5Vtbl
	OpenTaskManagerWindow ComProc
}

type ICoreWebView2_6 struct {
	Vtbl *ICoreWebView2_6Vtbl
}

func (i *ICoreWebView2_6) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_6) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_6 queries the object for its ICoreWebView2_6 interface. The receiver
// is the root of ICoreWebView2_6's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_6() (*ICoreWebView2_6, error) {
	var result *ICoreWebView2_6

	iidICoreWebView2_6 := NewGUID("{499aadac-d92c-4589-8a75-111bfc167795}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_6)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_6) OpenTaskManagerWindow() error {


	hr, _, _ := i.Vtbl.OpenTaskManagerWindow.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
