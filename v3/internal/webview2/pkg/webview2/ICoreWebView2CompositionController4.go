//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2CompositionController4Vtbl struct {
	IUnknownVtbl
	GetNonClientRegionAtPoint    ComProc
	QueryNonClientRegion         ComProc
	AddNonClientRegionChanged    ComProc
	RemoveNonClientRegionChanged ComProc
}

type ICoreWebView2CompositionController4 struct {
	Vtbl *ICoreWebView2CompositionController4Vtbl
}

func (i *ICoreWebView2CompositionController4) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2CompositionController4() *ICoreWebView2CompositionController4 {
	var result *ICoreWebView2CompositionController4

	iidICoreWebView2CompositionController4 := NewGUID("{7C367B9B-3D2B-450F-9E58-D61A20F486AA}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2CompositionController4)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2CompositionController4) GetNonClientRegionAtPoint(point POINT) (COREWEBVIEW2_NON_CLIENT_REGION_KIND, error) {

	var value COREWEBVIEW2_NON_CLIENT_REGION_KIND

	hr, _, _ := i.Vtbl.GetNonClientRegionAtPoint.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&point)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2CompositionController4) QueryNonClientRegion(kind COREWEBVIEW2_NON_CLIENT_REGION_KIND) (*ICoreWebView2RegionRectCollectionView, error) {

	var rects *ICoreWebView2RegionRectCollectionView

	hr, _, _ := i.Vtbl.QueryNonClientRegion.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(kind),
		uintptr(unsafe.Pointer(&rects)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return rects, nil
}

func (i *ICoreWebView2CompositionController4) AddNonClientRegionChanged(eventHandler *ICoreWebView2NonClientRegionChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddNonClientRegionChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2CompositionController4) RemoveNonClientRegionChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveNonClientRegionChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
