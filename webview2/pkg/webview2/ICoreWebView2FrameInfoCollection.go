//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2FrameInfoCollectionVtbl struct {
	IUnknownVtbl
	GetIterator ComProc
}

type ICoreWebView2FrameInfoCollection struct {
	Vtbl *ICoreWebView2FrameInfoCollectionVtbl
}

func (i *ICoreWebView2FrameInfoCollection) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameInfoCollection) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2FrameInfoCollection) GetIterator() (*ICoreWebView2FrameInfoCollectionIterator, error) {

	var value *ICoreWebView2FrameInfoCollectionIterator

	hr, _, err := i.Vtbl.GetIterator.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
