//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"errors"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type Bitmap struct {
	handle        w32.HBITMAP
	width, height int
}

func assembleBitmapFromHBITMAP(hbitmap w32.HBITMAP) (*Bitmap, error) {
	var dib w32.DIBSECTION
	if w32.GetObject(w32.HGDIOBJ(hbitmap), unsafe.Sizeof(dib), unsafe.Pointer(&dib)) == 0 {
		return nil, errors.New("GetObject for HBITMAP failed")
	}

	return &Bitmap{
		handle: hbitmap,
		width:  int(dib.DsBmih.BiWidth),
		height: int(dib.DsBmih.BiHeight),
	}, nil
}

func NewBitmapFromFile(filepath string, background Color) (*Bitmap, error) {
	var gpBitmap *uintptr
	var err error

	gpBitmap, err = w32.GdipCreateBitmapFromFile(filepath)
	if err != nil {
		return nil, err
	}
	defer w32.GdipDisposeImage(gpBitmap)

	var hbitmap w32.HBITMAP
	// Reverse RGB to BGR to satisfy gdiplus color schema.
	hbitmap, err = w32.GdipCreateHBITMAPFromBitmap(gpBitmap, uint32(RGB(background.B(), background.G(), background.R())))
	if err != nil {
		return nil, err
	}

	return assembleBitmapFromHBITMAP(hbitmap)
}

func NewBitmapFromResource(instance w32.HINSTANCE, resName *uint16, resType *uint16, background Color) (*Bitmap, error) {
	var gpBitmap *uintptr
	var err error
	var hRes w32.HRSRC

	hRes, err = w32.FindResource(w32.HMODULE(instance), resName, resType)
	if err != nil {
		return nil, err
	}
	resSize := w32.SizeofResource(w32.HMODULE(instance), hRes)
	pResData := w32.LockResource(w32.LoadResource(w32.HMODULE(instance), hRes))
	resBuffer := w32.GlobalAlloc(w32.GMEM_MOVEABLE, resSize)
	pResBuffer := w32.GlobalLock(resBuffer)
	w32.MoveMemory(pResBuffer, pResData, resSize)

	stream := w32.CreateStreamOnHGlobal(resBuffer, false)

	gpBitmap, err = w32.GdipCreateBitmapFromStream(stream)
	if err != nil {
		return nil, err
	}
	defer stream.Release()
	defer w32.GlobalUnlock(resBuffer)
	defer w32.GlobalFree(resBuffer)
	defer w32.GdipDisposeImage(gpBitmap)

	var hbitmap w32.HBITMAP
	// Reverse gform.RGB to BGR to satisfy gdiplus color schema.
	hbitmap, err = w32.GdipCreateHBITMAPFromBitmap(gpBitmap, uint32(RGB(background.B(), background.G(), background.R())))
	if err != nil {
		return nil, err
	}

	return assembleBitmapFromHBITMAP(hbitmap)
}

func (bm *Bitmap) Dispose() {
	if bm.handle != 0 {
		w32.DeleteObject(w32.HGDIOBJ(bm.handle))
		bm.handle = 0
	}
}

func (bm *Bitmap) GetHBITMAP() w32.HBITMAP {
	return bm.handle
}

func (bm *Bitmap) Size() (int, int) {
	return bm.width, bm.height
}

func (bm *Bitmap) Height() int {
	return bm.height
}

func (bm *Bitmap) Width() int {
	return bm.width
}
