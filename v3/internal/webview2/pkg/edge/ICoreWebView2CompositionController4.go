//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICoreWebView2CompositionController4Vtbl struct {
	_IUnknownVtbl
	GetRootVisualTarget          ComProc
	PutRootVisualTarget          ComProc
	SendMouseInput               ComProc
	SendPointerInput             ComProc
	GetCursor                    ComProc
	GetSystemCursorId            ComProc
	AddCursorChanged             ComProc
	RemoveCursorChanged          ComProc
	GetUIAProvider               ComProc
	DragEnter                    ComProc
	DragLeave                    ComProc
	DragOver                     ComProc
	Drop                         ComProc
	GetNonClientRegionAtPoint    ComProc
	QueryNonClientRegion         ComProc
	AddNonClientRegionChanged    ComProc
	RemoveNonClientRegionChanged ComProc
}

type ICoreWebView2CompositionController4 struct {
	Vtbl *ICoreWebView2CompositionController4Vtbl
}

func (i *ICoreWebView2CompositionController4) AddRef() uintptr {
	ret, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2CompositionController4) Release() uintptr {
	ret, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2CompositionController) GetICoreWebView2CompositionController4() *ICoreWebView2CompositionController4 {
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
		point.uintptr(),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
