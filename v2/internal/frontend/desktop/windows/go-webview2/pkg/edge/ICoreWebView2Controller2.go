//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2Controller2Vtbl struct {
	_IUnknownVtbl
	GetIsVisible                      ComProc
	PutIsVisible                      ComProc
	GetBounds                         ComProc
	PutBounds                         ComProc
	GetZoomFactor                     ComProc
	PutZoomFactor                     ComProc
	AddZoomFactorChanged              ComProc
	RemoveZoomFactorChanged           ComProc
	SetBoundsAndZoomFactor            ComProc
	MoveFocus                         ComProc
	AddMoveFocusRequested             ComProc
	RemoveMoveFocusRequested          ComProc
	AddGotFocus                       ComProc
	RemoveGotFocus                    ComProc
	AddLostFocus                      ComProc
	RemoveLostFocus                   ComProc
	AddAcceleratorKeyPressed          ComProc
	RemoveAcceleratorKeyPressed       ComProc
	GetParentWindow                   ComProc
	PutParentWindow                   ComProc
	NotifyParentWindowPositionChanged ComProc
	Close                             ComProc
	GetCoreWebView2                   ComProc
	GetDefaultBackgroundColor         ComProc
	PutDefaultBackgroundColor         ComProc
}

type ICoreWebView2Controller2 struct {
	vtbl *_ICoreWebView2Controller2Vtbl
}

func (i *ICoreWebView2Controller2) AddRef() uintptr {
	return i.AddRef()
}

func (i *ICoreWebView2Controller2) GetDefaultBackgroundColor() (*COREWEBVIEW2_COLOR, error) {
	var err error
	var backgroundColor *COREWEBVIEW2_COLOR
	_, _, err = i.vtbl.GetDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&backgroundColor)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return backgroundColor, nil
}

func (i *ICoreWebView2Controller2) PutDefaultBackgroundColor(backgroundColor COREWEBVIEW2_COLOR) error {
	var err error

	// Cast to a uint32 as that's what the call is expecting
	col := *(*uint32)(unsafe.Pointer(&backgroundColor))

	_, _, err = i.vtbl.PutDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(col),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
