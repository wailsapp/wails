//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2012 The W32 Authors. All Rights Reserved.
 */
package w32

import (
	"unsafe"
)

type pIDispatchVtbl struct {
	pQueryInterface   uintptr
	pAddRef           uintptr
	pRelease          uintptr
	pGetTypeInfoCount uintptr
	pGetTypeInfo      uintptr
	pGetIDsOfNames    uintptr
	pInvoke           uintptr
}

type IDispatch struct {
	lpVtbl *pIDispatchVtbl
}

func (this *IDispatch) QueryInterface(id *GUID) *IDispatch {
	return ComQueryInterface((*IUnknown)(unsafe.Pointer(this)), id)
}

func (this *IDispatch) AddRef() int32 {
	return ComAddRef((*IUnknown)(unsafe.Pointer(this)))
}

func (this *IDispatch) Release() int32 {
	return ComRelease((*IUnknown)(unsafe.Pointer(this)))
}

func (this *IDispatch) GetIDsOfName(names []string) []int32 {
	return ComGetIDsOfName(this, names)
}

func (this *IDispatch) Invoke(dispid int32, dispatch int16, params ...interface{}) *VARIANT {
	return ComInvoke(this, dispid, dispatch, params...)
}
