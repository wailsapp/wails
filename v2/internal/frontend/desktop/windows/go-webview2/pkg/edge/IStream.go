package edge

import "unsafe"

type _IStreamVtbl struct {
	_IUnknownVtbl
}

type IStream struct {
	vtbl *_IStreamVtbl
}

func (i *IStream) Release() error {
	return i.vtbl.CallRelease(unsafe.Pointer(i))
}
