//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PointerInfoVtbl struct {
	IUnknownVtbl
	GetPointerKind         ComProc
	PutPointerKind         ComProc
	GetPointerId           ComProc
	PutPointerId           ComProc
	GetFrameId             ComProc
	PutFrameId             ComProc
	GetPointerFlags        ComProc
	PutPointerFlags        ComProc
	GetPointerDeviceRect   ComProc
	PutPointerDeviceRect   ComProc
	GetDisplayRect         ComProc
	PutDisplayRect         ComProc
	GetPixelLocation       ComProc
	PutPixelLocation       ComProc
	GetHimetricLocation    ComProc
	PutHimetricLocation    ComProc
	GetPixelLocationRaw    ComProc
	PutPixelLocationRaw    ComProc
	GetHimetricLocationRaw ComProc
	PutHimetricLocationRaw ComProc
	GetTime                ComProc
	PutTime                ComProc
	GetHistoryCount        ComProc
	PutHistoryCount        ComProc
	GetInputData           ComProc
	PutInputData           ComProc
	GetKeyStates           ComProc
	PutKeyStates           ComProc
	GetPerformanceCount    ComProc
	PutPerformanceCount    ComProc
	GetButtonChangeKind    ComProc
	PutButtonChangeKind    ComProc
	GetPenFlags            ComProc
	PutPenFlags            ComProc
	GetPenMask             ComProc
	PutPenMask             ComProc
	GetPenPressure         ComProc
	PutPenPressure         ComProc
	GetPenRotation         ComProc
	PutPenRotation         ComProc
	GetPenTiltX            ComProc
	PutPenTiltX            ComProc
	GetPenTiltY            ComProc
	PutPenTiltY            ComProc
	GetTouchFlags          ComProc
	PutTouchFlags          ComProc
	GetTouchMask           ComProc
	PutTouchMask           ComProc
	GetTouchContact        ComProc
	PutTouchContact        ComProc
	GetTouchContactRaw     ComProc
	PutTouchContactRaw     ComProc
	GetTouchOrientation    ComProc
	PutTouchOrientation    ComProc
	GetTouchPressure       ComProc
	PutTouchPressure       ComProc
}

type ICoreWebView2PointerInfo struct {
	Vtbl *ICoreWebView2PointerInfoVtbl
}

func (i *ICoreWebView2PointerInfo) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2PointerInfo) GetPointerKind() (uint32, error) {

	var pointerKind uint32

	hr, _, _ := i.Vtbl.GetPointerKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerKind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return pointerKind, nil
}

func (i *ICoreWebView2PointerInfo) PutPointerKind(pointerKind uint32) error {

	hr, _, _ := i.Vtbl.PutPointerKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerKind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPointerId() (uint32, error) {

	var pointerId uint32

	hr, _, _ := i.Vtbl.GetPointerId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return pointerId, nil
}

func (i *ICoreWebView2PointerInfo) PutPointerId(pointerId uint32) error {

	hr, _, _ := i.Vtbl.PutPointerId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetFrameId() (uint32, error) {

	var frameId uint32

	hr, _, _ := i.Vtbl.GetFrameId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&frameId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return frameId, nil
}

func (i *ICoreWebView2PointerInfo) PutFrameId(frameId uint32) error {

	hr, _, _ := i.Vtbl.PutFrameId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&frameId)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPointerFlags() (uint32, error) {

	var pointerFlags uint32

	hr, _, _ := i.Vtbl.GetPointerFlags.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerFlags)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return pointerFlags, nil
}

func (i *ICoreWebView2PointerInfo) PutPointerFlags(pointerFlags uint32) error {

	hr, _, _ := i.Vtbl.PutPointerFlags.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerFlags)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPointerDeviceRect() (RECT, error) {

	var pointerDeviceRect RECT

	hr, _, _ := i.Vtbl.GetPointerDeviceRect.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerDeviceRect)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return RECT{}, syscall.Errno(hr)
	}
	return pointerDeviceRect, nil
}

func (i *ICoreWebView2PointerInfo) PutPointerDeviceRect(pointerDeviceRect RECT) error {

	hr, _, _ := i.Vtbl.PutPointerDeviceRect.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pointerDeviceRect)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetDisplayRect() (RECT, error) {

	var displayRect RECT

	hr, _, _ := i.Vtbl.GetDisplayRect.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&displayRect)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return RECT{}, syscall.Errno(hr)
	}
	return displayRect, nil
}

func (i *ICoreWebView2PointerInfo) PutDisplayRect(displayRect RECT) error {

	hr, _, _ := i.Vtbl.PutDisplayRect.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&displayRect)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPixelLocation() (POINT, error) {

	var pixelLocation POINT

	hr, _, _ := i.Vtbl.GetPixelLocation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pixelLocation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return pixelLocation, nil
}

func (i *ICoreWebView2PointerInfo) PutPixelLocation(pixelLocation POINT) error {

	hr, _, _ := i.Vtbl.PutPixelLocation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pixelLocation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetHimetricLocation() (POINT, error) {

	var himetricLocation POINT

	hr, _, _ := i.Vtbl.GetHimetricLocation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&himetricLocation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return himetricLocation, nil
}

func (i *ICoreWebView2PointerInfo) PutHimetricLocation(himetricLocation POINT) error {

	hr, _, _ := i.Vtbl.PutHimetricLocation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&himetricLocation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPixelLocationRaw() (POINT, error) {

	var pixelLocationRaw POINT

	hr, _, _ := i.Vtbl.GetPixelLocationRaw.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pixelLocationRaw)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return pixelLocationRaw, nil
}

func (i *ICoreWebView2PointerInfo) PutPixelLocationRaw(pixelLocationRaw POINT) error {

	hr, _, _ := i.Vtbl.PutPixelLocationRaw.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&pixelLocationRaw)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetHimetricLocationRaw() (POINT, error) {

	var himetricLocationRaw POINT

	hr, _, _ := i.Vtbl.GetHimetricLocationRaw.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&himetricLocationRaw)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return himetricLocationRaw, nil
}

func (i *ICoreWebView2PointerInfo) PutHimetricLocationRaw(himetricLocationRaw POINT) error {

	hr, _, _ := i.Vtbl.PutHimetricLocationRaw.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&himetricLocationRaw)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetTime() (uint32, error) {

	var time uint32

	hr, _, _ := i.Vtbl.GetTime.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&time)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return time, nil
}

func (i *ICoreWebView2PointerInfo) PutTime(time uint32) error {

	hr, _, _ := i.Vtbl.PutTime.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&time)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetHistoryCount() (uint32, error) {

	var historyCount uint32

	hr, _, _ := i.Vtbl.GetHistoryCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&historyCount)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return historyCount, nil
}

func (i *ICoreWebView2PointerInfo) PutHistoryCount(historyCount uint32) error {

	hr, _, _ := i.Vtbl.PutHistoryCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&historyCount)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetInputData() (int32, error) {

	var inputData int32

	hr, _, _ := i.Vtbl.GetInputData.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&inputData)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return inputData, nil
}

func (i *ICoreWebView2PointerInfo) PutInputData(inputData int32) error {

	hr, _, _ := i.Vtbl.PutInputData.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&inputData)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetKeyStates() (uint32, error) {

	var keyStates uint32

	hr, _, _ := i.Vtbl.GetKeyStates.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&keyStates)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return keyStates, nil
}

func (i *ICoreWebView2PointerInfo) PutKeyStates(keyStates uint32) error {

	hr, _, _ := i.Vtbl.PutKeyStates.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&keyStates)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPerformanceCount() (uint64, error) {

	var performanceCount uint64

	hr, _, _ := i.Vtbl.GetPerformanceCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&performanceCount)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return performanceCount, nil
}

func (i *ICoreWebView2PointerInfo) PutPerformanceCount(performanceCount uint64) error {

	hr, _, _ := i.Vtbl.PutPerformanceCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&performanceCount)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetButtonChangeKind() (int32, error) {

	var buttonChangeKind int32

	hr, _, _ := i.Vtbl.GetButtonChangeKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&buttonChangeKind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return buttonChangeKind, nil
}

func (i *ICoreWebView2PointerInfo) PutButtonChangeKind(buttonChangeKind int32) error {

	hr, _, _ := i.Vtbl.PutButtonChangeKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&buttonChangeKind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPenFlags() (uint32, error) {

	var penFLags uint32

	hr, _, _ := i.Vtbl.GetPenFlags.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penFLags)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return penFLags, nil
}

func (i *ICoreWebView2PointerInfo) PutPenFlags(penFLags uint32) error {

	hr, _, _ := i.Vtbl.PutPenFlags.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penFLags)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPenMask() (uint32, error) {

	var penMask uint32

	hr, _, _ := i.Vtbl.GetPenMask.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penMask)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return penMask, nil
}

func (i *ICoreWebView2PointerInfo) PutPenMask(penMask uint32) error {

	hr, _, _ := i.Vtbl.PutPenMask.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penMask)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPenPressure() (uint32, error) {

	var penPressure uint32

	hr, _, _ := i.Vtbl.GetPenPressure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penPressure)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return penPressure, nil
}

func (i *ICoreWebView2PointerInfo) PutPenPressure(penPressure uint32) error {

	hr, _, _ := i.Vtbl.PutPenPressure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penPressure)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPenRotation() (uint32, error) {

	var penRotation uint32

	hr, _, _ := i.Vtbl.GetPenRotation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penRotation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return penRotation, nil
}

func (i *ICoreWebView2PointerInfo) PutPenRotation(penRotation uint32) error {

	hr, _, _ := i.Vtbl.PutPenRotation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penRotation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPenTiltX() (int32, error) {

	var penTiltX int32

	hr, _, _ := i.Vtbl.GetPenTiltX.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penTiltX)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return penTiltX, nil
}

func (i *ICoreWebView2PointerInfo) PutPenTiltX(penTiltX int32) error {

	hr, _, _ := i.Vtbl.PutPenTiltX.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penTiltX)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetPenTiltY() (int32, error) {

	var penTiltY int32

	hr, _, _ := i.Vtbl.GetPenTiltY.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penTiltY)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return penTiltY, nil
}

func (i *ICoreWebView2PointerInfo) PutPenTiltY(penTiltY int32) error {

	hr, _, _ := i.Vtbl.PutPenTiltY.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&penTiltY)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetTouchFlags() (uint32, error) {

	var touchFlags uint32

	hr, _, _ := i.Vtbl.GetTouchFlags.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchFlags)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return touchFlags, nil
}

func (i *ICoreWebView2PointerInfo) PutTouchFlags(touchFlags uint32) error {

	hr, _, _ := i.Vtbl.PutTouchFlags.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchFlags)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetTouchMask() (uint32, error) {

	var touchMask uint32

	hr, _, _ := i.Vtbl.GetTouchMask.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchMask)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return touchMask, nil
}

func (i *ICoreWebView2PointerInfo) PutTouchMask(touchMask uint32) error {

	hr, _, _ := i.Vtbl.PutTouchMask.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchMask)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetTouchContact() (RECT, error) {

	var touchContact RECT

	hr, _, _ := i.Vtbl.GetTouchContact.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchContact)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return RECT{}, syscall.Errno(hr)
	}
	return touchContact, nil
}

func (i *ICoreWebView2PointerInfo) PutTouchContact(touchContact RECT) error {

	hr, _, _ := i.Vtbl.PutTouchContact.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchContact)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetTouchContactRaw() (RECT, error) {

	var touchContactRaw RECT

	hr, _, _ := i.Vtbl.GetTouchContactRaw.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchContactRaw)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return RECT{}, syscall.Errno(hr)
	}
	return touchContactRaw, nil
}

func (i *ICoreWebView2PointerInfo) PutTouchContactRaw(touchContactRaw RECT) error {

	hr, _, _ := i.Vtbl.PutTouchContactRaw.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchContactRaw)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetTouchOrientation() (uint32, error) {

	var touchOrientation uint32

	hr, _, _ := i.Vtbl.GetTouchOrientation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchOrientation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return touchOrientation, nil
}

func (i *ICoreWebView2PointerInfo) PutTouchOrientation(touchOrientation uint32) error {

	hr, _, _ := i.Vtbl.PutTouchOrientation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchOrientation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PointerInfo) GetTouchPressure() (uint32, error) {

	var touchPressure uint32

	hr, _, _ := i.Vtbl.GetTouchPressure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchPressure)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return touchPressure, nil
}

func (i *ICoreWebView2PointerInfo) PutTouchPressure(touchPressure uint32) error {

	hr, _, _ := i.Vtbl.PutTouchPressure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&touchPressure)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
