//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment12Vtbl struct {
	IUnknownVtbl
	CreateSharedBuffer ComProc
}

type ICoreWebView2Environment12 struct {
	Vtbl *ICoreWebView2Environment12Vtbl
}

func (i *ICoreWebView2Environment12) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment12) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Environment12() (*ICoreWebView2Environment12, error) {
	var result *ICoreWebView2Environment12

	iidICoreWebView2Environment12 := NewGUID("{f503db9b-739f-48dd-b151-fdfcf253f54e}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment12)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment12) CreateSharedBuffer(Size uint64) (*ICoreWebView2SharedBuffer, error) {

	var value *ICoreWebView2SharedBuffer

	hr, _, _ := i.Vtbl.CreateSharedBuffer.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(Size),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
