//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Profile2Vtbl struct {
	IUnknownVtbl
	ClearBrowsingData            ComProc
	ClearBrowsingDataInTimeRange ComProc
	ClearBrowsingDataAll         ComProc
}

type ICoreWebView2Profile2 struct {
	Vtbl *ICoreWebView2Profile2Vtbl
}

func (i *ICoreWebView2Profile2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Profile2() *ICoreWebView2Profile2 {
	var result *ICoreWebView2Profile2

	iidICoreWebView2Profile2 := NewGUID("{fa740d4b-5eae-4344-a8ad-74be31925397}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile2)),
		uintptr(unsafe.Pointer(&result)))

	return result
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

	hr, _, _ := i.Vtbl.ClearBrowsingDataInTimeRange.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(dataKinds),
		uintptr(unsafe.Pointer(&startTime)),
		uintptr(unsafe.Pointer(&endTime)),
		uintptr(unsafe.Pointer(handler)),
	)
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
