//go:build windows

package edge

import (
	"errors"
	"unsafe"
)

type compositionHost struct {
	device *iDCompositionDevice
	target *iDCompositionTarget
	visual *iDCompositionVisual
}

func newCompositionHost(hwnd uintptr) (*compositionHost, error) {
	device, err := dCompositionCreateDevice2()
	if err != nil {
		return nil, err
	}
	target, err := device.CreateTargetForHwnd(hwnd, true)
	if err != nil {
		device.Release()
		return nil, err
	}
	visual, err := device.CreateVisual()
	if err != nil {
		target.Release()
		device.Release()
		return nil, err
	}
	return &compositionHost{
		device: device,
		target: target,
		visual: visual,
	}, nil
}

func (h *compositionHost) attachController(controller *ICoreWebView2CompositionController) error {
	if h == nil || h.device == nil || h.target == nil || h.visual == nil || controller == nil {
		return errors.New("direct composition host is not initialized")
	}
	if err := controller.PutRootVisualTarget((*IUnknown)(unsafe.Pointer(h.visual))); err != nil {
		return err
	}
	if err := h.target.SetRoot(h.visual); err != nil {
		return err
	}
	return h.device.Commit()
}

func (h *compositionHost) release() {
	if h == nil {
		return
	}
	if h.visual != nil {
		h.visual.Release()
		h.visual = nil
	}
	if h.target != nil {
		h.target.Release()
		h.target = nil
	}
	if h.device != nil {
		h.device.Release()
		h.device = nil
	}
}
