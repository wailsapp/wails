//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2EnvironmentOptions8Vtbl struct {
	IUnknownVtbl
	GetScrollBarStyle ComProc
	PutScrollBarStyle ComProc
}

type ICoreWebView2EnvironmentOptions8 struct {
	Vtbl *ICoreWebView2EnvironmentOptions8Vtbl
}

func (i *ICoreWebView2EnvironmentOptions8) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2EnvironmentOptions8) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2EnvironmentOptions8) GetScrollBarStyle() (COREWEBVIEW2_SCROLLBAR_STYLE, error) {

	var value COREWEBVIEW2_SCROLLBAR_STYLE

	hr, _, err := i.Vtbl.GetScrollBarStyle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2EnvironmentOptions8) PutScrollBarStyle(value COREWEBVIEW2_SCROLLBAR_STYLE) error {


	hr, _, err := i.Vtbl.PutScrollBarStyle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
