//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type Canvas struct {
	hwnd         w32.HWND
	hdc          w32.HDC
	doNotDispose bool
}

var nullBrush = NewNullBrush()

func NewCanvasFromHwnd(hwnd w32.HWND) *Canvas {
	hdc := w32.GetDC(hwnd)
	if hdc == 0 {
		panic(fmt.Sprintf("Create canvas from %v failed.", hwnd))
	}

	return &Canvas{hwnd: hwnd, hdc: hdc, doNotDispose: false}
}

func NewCanvasFromHDC(hdc w32.HDC) *Canvas {
	if hdc == 0 {
		panic("Cannot create canvas from invalid HDC.")
	}

	return &Canvas{hdc: hdc, doNotDispose: true}
}

func (ca *Canvas) Dispose() {
	if !ca.doNotDispose && ca.hdc != 0 {
		if ca.hwnd == 0 {
			w32.DeleteDC(ca.hdc)
		} else {
			w32.ReleaseDC(ca.hwnd, ca.hdc)
		}

		ca.hdc = 0
	}
}

func (ca *Canvas) DrawBitmap(bmp *Bitmap, x, y int) {
	cdc := w32.CreateCompatibleDC(0)
	defer w32.DeleteDC(cdc)

	hbmpOld := w32.SelectObject(cdc, w32.HGDIOBJ(bmp.GetHBITMAP()))
	defer w32.SelectObject(cdc, w32.HGDIOBJ(hbmpOld))

	w, h := bmp.Size()

	w32.BitBlt(ca.hdc, x, y, w, h, cdc, 0, 0, w32.SRCCOPY)
}

func (ca *Canvas) DrawStretchedBitmap(bmp *Bitmap, rect *Rect) {
	cdc := w32.CreateCompatibleDC(0)
	defer w32.DeleteDC(cdc)

	hbmpOld := w32.SelectObject(cdc, w32.HGDIOBJ(bmp.GetHBITMAP()))
	defer w32.SelectObject(cdc, w32.HGDIOBJ(hbmpOld))

	w, h := bmp.Size()

	rc := rect.GetW32Rect()
	w32.StretchBlt(ca.hdc, int(rc.Left), int(rc.Top), int(rc.Right), int(rc.Bottom), cdc, 0, 0, w, h, w32.SRCCOPY)
}

func (ca *Canvas) DrawIcon(ico *Icon, x, y int) bool {
	return w32.DrawIcon(ca.hdc, x, y, ico.Handle())
}

// DrawFillRect draw and fill rectangle with color.
func (ca *Canvas) DrawFillRect(rect *Rect, pen *Pen, brush *Brush) {
	w32Rect := rect.GetW32Rect()

	previousPen := w32.SelectObject(ca.hdc, w32.HGDIOBJ(pen.GetHPEN()))
	defer w32.SelectObject(ca.hdc, previousPen)

	previousBrush := w32.SelectObject(ca.hdc, w32.HGDIOBJ(brush.GetHBRUSH()))
	defer w32.SelectObject(ca.hdc, previousBrush)

	w32.Rectangle(ca.hdc, w32Rect.Left, w32Rect.Top, w32Rect.Right, w32Rect.Bottom)
}

func (ca *Canvas) DrawRect(rect *Rect, pen *Pen) {
	w32Rect := rect.GetW32Rect()

	previousPen := w32.SelectObject(ca.hdc, w32.HGDIOBJ(pen.GetHPEN()))
	defer w32.SelectObject(ca.hdc, previousPen)

	// nullBrush is used to make interior of the rect transparent
	previousBrush := w32.SelectObject(ca.hdc, w32.HGDIOBJ(nullBrush.GetHBRUSH()))
	defer w32.SelectObject(ca.hdc, previousBrush)

	w32.Rectangle(ca.hdc, w32Rect.Left, w32Rect.Top, w32Rect.Right, w32Rect.Bottom)
}

func (ca *Canvas) FillRect(rect *Rect, brush *Brush) {
	w32.FillRect(ca.hdc, rect.GetW32Rect(), brush.GetHBRUSH())
}

func (ca *Canvas) DrawEllipse(rect *Rect, pen *Pen) {
	w32Rect := rect.GetW32Rect()

	previousPen := w32.SelectObject(ca.hdc, w32.HGDIOBJ(pen.GetHPEN()))
	defer w32.SelectObject(ca.hdc, previousPen)

	// nullBrush is used to make interior of the rect transparent
	previousBrush := w32.SelectObject(ca.hdc, w32.HGDIOBJ(nullBrush.GetHBRUSH()))
	defer w32.SelectObject(ca.hdc, previousBrush)

	w32.Ellipse(ca.hdc, w32Rect.Left, w32Rect.Top, w32Rect.Right, w32Rect.Bottom)
}

// DrawFillEllipse draw and fill ellipse with color.
func (ca *Canvas) DrawFillEllipse(rect *Rect, pen *Pen, brush *Brush) {
	w32Rect := rect.GetW32Rect()

	previousPen := w32.SelectObject(ca.hdc, w32.HGDIOBJ(pen.GetHPEN()))
	defer w32.SelectObject(ca.hdc, previousPen)

	previousBrush := w32.SelectObject(ca.hdc, w32.HGDIOBJ(brush.GetHBRUSH()))
	defer w32.SelectObject(ca.hdc, previousBrush)

	w32.Ellipse(ca.hdc, w32Rect.Left, w32Rect.Top, w32Rect.Right, w32Rect.Bottom)
}

func (ca *Canvas) DrawLine(x, y, x2, y2 int, pen *Pen) {
	w32.MoveToEx(ca.hdc, x, y, nil)

	previousPen := w32.SelectObject(ca.hdc, w32.HGDIOBJ(pen.GetHPEN()))
	defer w32.SelectObject(ca.hdc, previousPen)

	w32.LineTo(ca.hdc, int32(x2), int32(y2))
}

// Refer win32 DrawText document for uFormat.
func (ca *Canvas) DrawText(text string, rect *Rect, format uint, font *Font, textColor Color) {
	previousFont := w32.SelectObject(ca.hdc, w32.HGDIOBJ(font.GetHFONT()))
	defer w32.SelectObject(ca.hdc, w32.HGDIOBJ(previousFont))

	previousBkMode := w32.SetBkMode(ca.hdc, w32.TRANSPARENT)
	defer w32.SetBkMode(ca.hdc, previousBkMode)

	previousTextColor := w32.SetTextColor(ca.hdc, w32.COLORREF(textColor))
	defer w32.SetTextColor(ca.hdc, previousTextColor)

	w32.DrawText(ca.hdc, text, len(text), rect.GetW32Rect(), format)
}
