//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2CompositionController3Vtbl struct {
	ICoreWebView2CompositionController2Vtbl
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


// GetICoreWebView2CompositionController3 queries the object for its ICoreWebView2CompositionController3 interface. The receiver
// is the root of ICoreWebView2CompositionController3's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2CompositionController) GetICoreWebView2CompositionController3() (*ICoreWebView2CompositionController3, error) {
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
	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.DragEnter.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(unsafe.Pointer(dataObject)),
			uintptr(keyState),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[1]),
			uintptr(unsafe.Pointer(&effect)),
		)
	default:
		hr, _, _ = i.Vtbl.DragEnter.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(unsafe.Pointer(dataObject)),
			uintptr(keyState),
			uintptr(*(*uint64)(unsafe.Pointer(&point))),
			uintptr(unsafe.Pointer(&effect)),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return effect, nil
}

func (i *ICoreWebView2CompositionController3) DragLeave() error {


	hr, _, _ := i.Vtbl.DragLeave.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CompositionController3) DragOver(keyState uint32, point POINT) (uint32, error) {

	var effect uint32
	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.DragOver.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(keyState),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[1]),
			uintptr(unsafe.Pointer(&effect)),
		)
	default:
		hr, _, _ = i.Vtbl.DragOver.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(keyState),
			uintptr(*(*uint64)(unsafe.Pointer(&point))),
			uintptr(unsafe.Pointer(&effect)),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return effect, nil
}

func (i *ICoreWebView2CompositionController3) Drop(dataObject *IDataObject, keyState uint32, point POINT) (uint32, error) {

	var effect uint32
	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.Drop.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(unsafe.Pointer(dataObject)),
			uintptr(keyState),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&point)))[1]),
			uintptr(unsafe.Pointer(&effect)),
		)
	default:
		hr, _, _ = i.Vtbl.Drop.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(unsafe.Pointer(dataObject)),
			uintptr(keyState),
			uintptr(*(*uint64)(unsafe.Pointer(&point))),
			uintptr(unsafe.Pointer(&effect)),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return effect, nil
}
