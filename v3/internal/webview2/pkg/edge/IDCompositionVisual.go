//go:build windows

package edge

import "unsafe"

type iDCompositionVisualVtbl struct {
	_IUnknownVtbl
}

type iDCompositionVisual struct {
	vtbl *iDCompositionVisualVtbl
}

func (v *iDCompositionVisual) AddRef() uintptr {
	ret, _, _ := v.vtbl.AddRef.Call(uintptr(unsafe.Pointer(v)))

	return ret
}

func (v *iDCompositionVisual) Release() uintptr {
	ret, _, _ := v.vtbl.Release.Call(uintptr(unsafe.Pointer(v)))

	return ret
}
