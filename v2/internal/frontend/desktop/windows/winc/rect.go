//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type Rect struct {
	rect w32.RECT
}

func NewEmptyRect() *Rect {
	var newRect Rect
	w32.SetRectEmpty(&newRect.rect)

	return &newRect
}

func NewRect(left, top, right, bottom int) *Rect {
	var newRect Rect
	w32.SetRectEmpty(&newRect.rect)
	newRect.Set(left, top, right, bottom)

	return &newRect
}

func (re *Rect) Data() (left, top, right, bottom int32) {
	left = re.rect.Left
	top = re.rect.Top
	right = re.rect.Right
	bottom = re.rect.Bottom
	return
}

func (re *Rect) Width() int {
	return int(re.rect.Right - re.rect.Left)
}

func (re *Rect) Height() int {
	return int(re.rect.Bottom - re.rect.Top)
}

func (re *Rect) GetW32Rect() *w32.RECT {
	return &re.rect
}

func (re *Rect) Set(left, top, right, bottom int) {
	w32.SetRect(&re.rect, left, top, right, bottom)
}

func (re *Rect) IsEqual(rect *Rect) bool {
	return w32.EqualRect(&re.rect, &rect.rect)
}

func (re *Rect) Inflate(x, y int) {
	w32.InflateRect(&re.rect, x, y)
}

func (re *Rect) Intersect(src *Rect) {
	w32.IntersectRect(&re.rect, &re.rect, &src.rect)
}

func (re *Rect) IsEmpty() bool {
	return w32.IsRectEmpty(&re.rect)
}

func (re *Rect) Offset(x, y int) {
	w32.OffsetRect(&re.rect, x, y)
}

func (re *Rect) IsPointIn(x, y int) bool {
	return w32.PtInRect(&re.rect, x, y)
}

func (re *Rect) Substract(src *Rect) {
	w32.SubtractRect(&re.rect, &re.rect, &src.rect)
}

func (re *Rect) Union(src *Rect) {
	w32.UnionRect(&re.rect, &re.rect, &src.rect)
}
