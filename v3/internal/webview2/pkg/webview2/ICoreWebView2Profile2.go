//go:build windows

package webview2
import (
	"math"
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Profile2Vtbl struct {
	ICoreWebView2ProfileVtbl
	ClearBrowsingData ComProc
	ClearBrowsingDataInTimeRange ComProc
	ClearBrowsingDataAll ComProc
}

type ICoreWebView2Profile2 struct {
	Vtbl *ICoreWebView2Profile2Vtbl
}

func (i *ICoreWebView2Profile2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Profile2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Profile2 queries the object for its ICoreWebView2Profile2 interface. The receiver
// is the root of ICoreWebView2Profile2's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Profile) GetICoreWebView2Profile2() (*ICoreWebView2Profile2, error) {
	var result *ICoreWebView2Profile2

	iidICoreWebView2Profile2 := NewGUID("{fa740d4b-5eae-4344-a8ad-74be31925397}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Profile2) ClearBrowsingData(dataKinds COREWEBVIEW2_BROWSING_DATA_KINDS, handler *ICoreWebView2ClearBrowsingDataCompletedHandler) error {


	hr, _, _ := i.Vtbl.ClearBrowsingData.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(dataKinds),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Profile2) ClearBrowsingDataInTimeRange(dataKinds COREWEBVIEW2_BROWSING_DATA_KINDS, startTime float64, endTime float64, handler *ICoreWebView2ClearBrowsingDataCompletedHandler) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.ClearBrowsingDataInTimeRange.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(dataKinds),
			uintptr(uint32(math.Float64bits(startTime))),
			uintptr(uint32(math.Float64bits(startTime)>>32)),
			uintptr(uint32(math.Float64bits(endTime))),
			uintptr(uint32(math.Float64bits(endTime)>>32)),
			uintptr(unsafe.Pointer(handler)),
		)
	default:
		hr, _, _ = i.Vtbl.ClearBrowsingDataInTimeRange.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(dataKinds),
			uintptr(math.Float64bits(startTime)),
			uintptr(math.Float64bits(endTime)),
			uintptr(unsafe.Pointer(handler)),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Profile2) ClearBrowsingDataAll(handler *ICoreWebView2ClearBrowsingDataCompletedHandler) error {


	hr, _, _ := i.Vtbl.ClearBrowsingDataAll.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
