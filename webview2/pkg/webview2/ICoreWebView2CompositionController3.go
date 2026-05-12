//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2CompositionController3Vtbl struct {
	IUnknownVtbl
	DragEnter ComProc
	DragLeave ComProc
	DragOver ComProc
	Drop ComProc
}

type ICoreWebView2CompositionController3 struct {
	Vtbl *ICoreWebView2CompositionController3Vtbl
}

func (i *ICoreWebView2CompositionController3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CompositionController3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2CompositionController3() (*ICoreWebView2CompositionController3, error) {
	var result *ICoreWebView2CompositionController3

	iidICoreWebView2CompositionController3 := NewGUID("{9570570e-4d76-4361-9ee1-f04d0dbdfb1e}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2CompositionController3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2CompositionController3) DragEnter(dataObject *IDataObject, keyState uint32, point POINT) (uint32, error) {

	var effect uint32

	hr, _, err := i.Vtbl.DragEnter.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(dataObject)),
		uintptr(keyState),
		uintptr(unsafe.Pointer(&point)),
		uintptr(unsafe.Pointer(&effect)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return effect, err
}

func (i *ICoreWebView2CompositionController3) DragLeave() error {


	hr, _, err := i.Vtbl.DragLeave.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2CompositionController3) DragOver(keyState uint32, point POINT) (uint32, error) {

	var effect uint32

	hr, _, err := i.Vtbl.DragOver.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(keyState),
		uintptr(unsafe.Pointer(&point)),
		uintptr(unsafe.Pointer(&effect)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return effect, err
}

func (i *ICoreWebView2CompositionController3) Drop(dataObject *IDataObject, keyState uint32, point POINT) (uint32, error) {

	var effect uint32

	hr, _, err := i.Vtbl.Drop.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(dataObject)),
		uintptr(keyState),
		uintptr(unsafe.Pointer(&point)),
		uintptr(unsafe.Pointer(&effect)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return effect, err
}
