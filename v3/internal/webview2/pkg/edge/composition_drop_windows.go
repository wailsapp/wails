//go:build windows

package edge

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/webview2/pkg/combridge"
	"golang.org/x/sys/windows"
)

var dropOle32 = windows.NewLazySystemDLL("ole32.dll")
var dropRegister = dropOle32.NewProc("RegisterDragDrop")
var dropRevoke = dropOle32.NewProc("RevokeDragDrop")
var dropOleInitialize = dropOle32.NewProc("OleInitialize")
var dropOleUninitialize = dropOle32.NewProc("OleUninitialize")
var dropScreenToClient = windows.NewLazySystemDLL("user32.dll").NewProc("ScreenToClient")

type compositionController3Vtbl struct {
	ICoreWebView2CompositionControllerVtbl
	GetUIAProvider ComProc
	DragEnter      ComProc
	DragLeave      ComProc
	DragOver       ComProc
	Drop           ComProc
}

type compositionController3 struct{ Vtbl *compositionController3Vtbl }

// POINTL is passed by value as one register on Windows amd64 and arm64,
// just like POINT in the composition controller's SendMouseInput method.
type compositionDropInterface interface {
	combridge.IUnknown
	enter(uintptr, uint32, uintptr, *uint32) uintptr
	over(uint32, uintptr, *uint32) uintptr
	leave() uintptr
	drop(uintptr, uint32, uintptr, *uint32) uintptr
}

func init() {
	combridge.RegisterVTable[combridge.IUnknown, compositionDropInterface](
		"{00000122-0000-0000-C000-000000000046}",
		func(this, data uintptr, keys uint32, point uintptr, effect *uint32) uintptr {
			return combridge.Resolve[compositionDropInterface](this).enter(data, keys, point, effect)
		},
		func(this uintptr, keys uint32, point uintptr, effect *uint32) uintptr {
			return combridge.Resolve[compositionDropInterface](this).over(keys, point, effect)
		},
		func(this uintptr) uintptr { return combridge.Resolve[compositionDropInterface](this).leave() },
		func(this, data uintptr, keys uint32, point uintptr, effect *uint32) uintptr {
			return combridge.Resolve[compositionDropInterface](this).drop(data, keys, point, effect)
		},
	)
}

type compositionDropTarget struct {
	combridge.IUnknownImpl
	hwnd       uintptr
	controller *compositionController3
	object     *combridge.ComObject[compositionDropInterface]
}

func (e *Chromium) registerCompositionDropTarget() error {
	var controller *compositionController3
	iid := NewGUID("{9570570e-4d76-4361-9ee1-f04d0dbdfb1e}")
	hr, _, _ := e.compositionController.Vtbl.QueryInterface.Call(uintptr(unsafe.Pointer(e.compositionController)), uintptr(unsafe.Pointer(iid)), uintptr(unsafe.Pointer(&controller)))
	if int32(hr) < 0 || controller == nil {
		return fmt.Errorf("query composition drag controller: %08x", hr)
	}
	hr, _, _ = dropOleInitialize.Call(0)
	if int32(hr) < 0 {
		controller.Vtbl.Release.Call(uintptr(unsafe.Pointer(controller)))
		return fmt.Errorf("initialize OLE drag/drop: %08x", hr)
	}
	target := &compositionDropTarget{hwnd: e.hwnd, controller: controller}
	target.object = combridge.New[compositionDropInterface](target)
	hr, _, _ = dropRegister.Call(e.hwnd, target.object.Ref())
	if int32(hr) < 0 {
		target.object.Close()
		controller.Vtbl.Release.Call(uintptr(unsafe.Pointer(controller)))
		dropOleUninitialize.Call()
		return fmt.Errorf("register composition drop target: %w", syscall.Errno(hr))
	}
	e.compositionDropTarget = target
	return nil
}

func (e *Chromium) releaseCompositionDropTarget() {
	d := e.compositionDropTarget
	if d == nil {
		return
	}
	// Detach first: native teardown can reenter shutdown, and Close also calls
	// this after ShuttingDown. Each registration owns one OLE initialization.
	e.compositionDropTarget = nil
	dropRevoke.Call(d.hwnd)
	d.object.Close()
	d.controller.Vtbl.Release.Call(uintptr(unsafe.Pointer(d.controller)))
	dropOleUninitialize.Call()
}

func (d *compositionDropTarget) clientPoint(packed uintptr) (uintptr, bool) {
	p := POINT{X: int32(uint32(packed)), Y: int32(uint64(packed) >> 32)}
	ok, _, _ := dropScreenToClient.Call(d.hwnd, uintptr(unsafe.Pointer(&p)))
	return p.uintptr(), ok != 0
}

func (d *compositionDropTarget) enter(data uintptr, keys uint32, point uintptr, effect *uint32) uintptr {
	p, ok := d.clientPoint(point)
	if !ok {
		*effect = 0
		return 0x80004005
	}
	hr, _, _ := d.controller.Vtbl.DragEnter.Call(uintptr(unsafe.Pointer(d.controller)), data, uintptr(keys), p, uintptr(unsafe.Pointer(effect)))
	return hr
}

func (d *compositionDropTarget) over(keys uint32, point uintptr, effect *uint32) uintptr {
	p, ok := d.clientPoint(point)
	if !ok {
		*effect = 0
		return 0x80004005
	}
	hr, _, _ := d.controller.Vtbl.DragOver.Call(uintptr(unsafe.Pointer(d.controller)), uintptr(keys), p, uintptr(unsafe.Pointer(effect)))
	return hr
}

func (d *compositionDropTarget) leave() uintptr {
	hr, _, _ := d.controller.Vtbl.DragLeave.Call(uintptr(unsafe.Pointer(d.controller)))
	return hr
}

func (d *compositionDropTarget) drop(data uintptr, keys uint32, point uintptr, effect *uint32) uintptr {
	p, ok := d.clientPoint(point)
	if !ok {
		*effect = 0
		return 0x80004005
	}
	hr, _, _ := d.controller.Vtbl.Drop.Call(uintptr(unsafe.Pointer(d.controller)), data, uintptr(keys), p, uintptr(unsafe.Pointer(effect)))
	return hr
}
