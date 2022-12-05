//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import "github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"

type ImageView struct {
	ControlBase

	bmp *Bitmap
}

func NewImageView(parent Controller) *ImageView {
	iv := new(ImageView)

	iv.InitWindow("winc_ImageView", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	RegMsgHandler(iv)

	iv.SetFont(DefaultFont)
	iv.SetText("")
	iv.SetSize(200, 65)
	return iv
}

func (iv *ImageView) DrawImageFile(filepath string) error {
	bmp, err := NewBitmapFromFile(filepath, RGB(255, 255, 0))
	if err != nil {
		return err
	}
	iv.bmp = bmp
	return nil
}

func (iv *ImageView) DrawImage(bmp *Bitmap) {
	iv.bmp = bmp
}

func (iv *ImageView) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_SIZE, w32.WM_SIZING:
		iv.Invalidate(true)

	case w32.WM_ERASEBKGND:
		return 1 // important

	case w32.WM_PAINT:
		if iv.bmp != nil {
			canvas := NewCanvasFromHwnd(iv.hwnd)
			defer canvas.Dispose()
			iv.SetSize(iv.bmp.Size())
			canvas.DrawBitmap(iv.bmp, 0, 0)
		}
	}
	return w32.DefWindowProc(iv.hwnd, msg, wparam, lparam)
}
