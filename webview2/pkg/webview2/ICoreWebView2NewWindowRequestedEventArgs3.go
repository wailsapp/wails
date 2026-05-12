//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NewWindowRequestedEventArgs3Vtbl struct {
	IUnknownVtbl
	GetOriginalSourceFrameInfo ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs3 struct {
	Vtbl *ICoreWebView2NewWindowRequestedEventArgs3Vtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NewWindowRequestedEventArgs3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2NewWindowRequestedEventArgs3() (*ICoreWebView2NewWindowRequestedEventArgs3, error) {
	var result *ICoreWebView2NewWindowRequestedEventArgs3

	iidICoreWebView2NewWindowRequestedEventArgs3 := NewGUID("{842bed3c-6ad6-4dd9-b938-28c96667ad66}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NewWindowRequestedEventArgs3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2NewWindowRequestedEventArgs3) GetOriginalSourceFrameInfo() (*ICoreWebView2FrameInfo, error) {

	var value *ICoreWebView2FrameInfo

	hr, _, err := i.Vtbl.GetOriginalSourceFrameInfo.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
