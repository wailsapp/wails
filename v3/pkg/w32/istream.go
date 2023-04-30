//go:build windows

/*
 * Copyright (C) 2019 Tad Vizbaras. All Rights Reserved.
 * Copyright (C) 2010-2012 The W32 Authors. All Rights Reserved.
 */
package w32

import (
	"unsafe"
)

type pIStreamVtbl struct {
	pQueryInterface uintptr
	pAddRef         uintptr
	pRelease        uintptr
}

type IStream struct {
	lpVtbl *pIStreamVtbl
}

func (this *IStream) QueryInterface(id *GUID) *IDispatch {
	return ComQueryInterface((*IUnknown)(unsafe.Pointer(this)), id)
}

func (this *IStream) AddRef() int32 {
	return ComAddRef((*IUnknown)(unsafe.Pointer(this)))
}

func (this *IStream) Release() int32 {
	return ComRelease((*IUnknown)(unsafe.Pointer(this)))
}
