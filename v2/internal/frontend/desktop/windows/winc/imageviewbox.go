//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import (
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type direction int

const (
	DirNone direction = iota
	DirX
	DirY
	DirX2
	DirY2
)

var ImageBoxPen = NewPen(w32.PS_GEOMETRIC, 2, NewSolidColorBrush(RGB(140, 140, 220)))
var ImageBoxHiPen = NewPen(w32.PS_GEOMETRIC, 2, NewSolidColorBrush(RGB(220, 140, 140)))
var ImageBoxMarkBrush = NewSolidColorBrush(RGB(40, 40, 40))
var ImageBoxMarkPen = NewPen(w32.PS_GEOMETRIC, 2, ImageBoxMarkBrush)

type ImageBox struct {
	Name         string
	Type         int
	X, Y, X2, Y2 int

	underMouse bool // dynamic value
}

func (b *ImageBox) Rect() *Rect {
	return NewRect(b.X, b.Y, b.X2, b.Y2)
}

// ImageViewBox is image view with boxes.
type ImageViewBox struct {
	ControlBase

	bmp       *Bitmap
	mouseLeft bool
	modified  bool // used by GUI to see if any image box modified

	add bool

	Boxes   []*ImageBox // might be persisted to file
	dragBox *ImageBox
	selBox  *ImageBox

	dragStartX, dragStartY int
	resize                 direction

	onSelectedChange EventManager
	onAdd            EventManager
	onModify         EventManager
}

func NewImageViewBox(parent Controller) *ImageViewBox {
	iv := new(ImageViewBox)

	iv.InitWindow("winc_ImageViewBox", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	RegMsgHandler(iv)

	iv.SetFont(DefaultFont)
	iv.SetText("")
	iv.SetSize(200, 65)

	return iv
}

func (iv *ImageViewBox) OnSelectedChange() *EventManager {
	return &iv.onSelectedChange
}

func (iv *ImageViewBox) OnAdd() *EventManager {
	return &iv.onAdd
}

func (iv *ImageViewBox) OnModify() *EventManager {
	return &iv.onModify
}

func (iv *ImageViewBox) IsModified() bool          { return iv.modified }
func (iv *ImageViewBox) SetModified(modified bool) { iv.modified = modified }
func (iv *ImageViewBox) IsLoaded() bool            { return iv.bmp != nil }
func (iv *ImageViewBox) AddMode() bool             { return iv.add }
func (iv *ImageViewBox) SetAddMode(add bool)       { iv.add = add }
func (iv *ImageViewBox) HasSelected() bool         { return iv.selBox != nil && iv.bmp != nil }

func (iv *ImageViewBox) wasModified() {
	iv.modified = true
	iv.onModify.Fire(NewEvent(iv, nil))
}

func (iv *ImageViewBox) DeleteSelected() {
	if iv.selBox != nil {
		for i, b := range iv.Boxes {
			if b == iv.selBox {
				iv.Boxes = append(iv.Boxes[:i], iv.Boxes[i+1:]...)
				iv.selBox = nil
				iv.Invalidate(true)
				iv.wasModified()
				iv.onSelectedChange.Fire(NewEvent(iv, nil))
				return
			}
		}
	}
}

func (iv *ImageViewBox) NameSelected() string {
	if iv.selBox != nil {
		return iv.selBox.Name
	}
	return ""
}

func (iv *ImageViewBox) SetNameSelected(name string) {
	if iv.selBox != nil {
		iv.selBox.Name = name
		iv.wasModified()
	}
}

func (iv *ImageViewBox) TypeSelected() int {
	if iv.selBox != nil {
		return iv.selBox.Type
	}
	return 0
}

func (iv *ImageViewBox) SetTypeSelected(typ int) {
	if iv.selBox != nil {
		iv.selBox.Type = typ
		iv.wasModified()
	}
}

func (ib *ImageViewBox) updateHighlight(x, y int) bool {
	var changed bool
	for _, b := range ib.Boxes {
		under := x >= b.X && y >= b.Y && x <= b.X2 && y <= b.Y2
		if b.underMouse != under {
			changed = true
		}
		b.underMouse = under
		/*if sel {
			break // allow only one to be underMouse
		}*/
	}
	return changed
}

func (ib *ImageViewBox) isUnderMouse(x, y int) *ImageBox {
	for _, b := range ib.Boxes {
		if x >= b.X && y >= b.Y && x <= b.X2 && y <= b.Y2 {
			return b
		}
	}
	return nil
}

func (ib *ImageViewBox) getCursor(x, y int) uint16 {
	for _, b := range ib.Boxes {
		switch d := ib.resizingDirection(b, x, y); d {
		case DirY, DirY2:
			return w32.IDC_SIZENS
		case DirX, DirX2:
			return w32.IDC_SIZEWE
		}
		// w32.IDC_SIZEALL or w32.IDC_SIZE for resize
	}
	return w32.IDC_ARROW
}

func (ib *ImageViewBox) resizingDirection(b *ImageBox, x, y int) direction {
	if b == nil {
		return DirNone
	}
	switch {
	case b.X == x || b.X == x-1 || b.X == x+1:
		return DirX
	case b.X2 == x || b.X2 == x-1 || b.X2 == x+1:
		return DirX2
	case b.Y == y || b.Y == y-1 || b.Y == y+1:
		return DirY
	case b.Y2 == y || b.Y2 == y-1 || b.Y2 == y+1:
		return DirY2
	}
	return DirNone
}

func (ib *ImageViewBox) resizeToDirection(b *ImageBox, x, y int) {
	switch ib.resize {
	case DirX:
		b.X = x
	case DirY:
		b.Y = y
	case DirX2:
		b.X2 = x
	case DirY2:
		b.Y2 = y
	}
}

func (ib *ImageViewBox) drag(b *ImageBox, x, y int) {
	w, h := b.X2-b.X, b.Y2-b.Y

	nx := ib.dragStartX - b.X
	ny := ib.dragStartY - b.Y

	b.X = x - nx
	b.Y = y - ny
	b.X2 = b.X + w
	b.Y2 = b.Y + h

	ib.dragStartX, ib.dragStartY = x, y
}

func (iv *ImageViewBox) DrawImageFile(filepath string) (err error) {
	iv.bmp, err = NewBitmapFromFile(filepath, RGB(255, 255, 0))
	iv.selBox = nil
	iv.modified = false
	iv.onSelectedChange.Fire(NewEvent(iv, nil))
	iv.onModify.Fire(NewEvent(iv, nil))
	return
}

func (iv *ImageViewBox) DrawImage(bmp *Bitmap) {
	iv.bmp = bmp
	iv.selBox = nil
	iv.modified = false
	iv.onSelectedChange.Fire(NewEvent(iv, nil))
	iv.onModify.Fire(NewEvent(iv, nil))
}

func (iv *ImageViewBox) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_SIZE, w32.WM_SIZING:
		iv.Invalidate(true)

	case w32.WM_ERASEBKGND:
		return 1 // important

	case w32.WM_CREATE:
		internalTrackMouseEvent(iv.hwnd)

	case w32.WM_PAINT:
		if iv.bmp != nil {
			canvas := NewCanvasFromHwnd(iv.hwnd)
			defer canvas.Dispose()
			iv.SetSize(iv.bmp.Size())
			canvas.DrawBitmap(iv.bmp, 0, 0)

			for _, b := range iv.Boxes {
				// old code used NewSystemColorBrush(w32.COLOR_BTNFACE) w32.COLOR_WINDOW
				pen := ImageBoxPen
				if b.underMouse {
					pen = ImageBoxHiPen
				}
				canvas.DrawRect(b.Rect(), pen)

				if b == iv.selBox {
					x1 := []int{b.X, b.X2, b.X2, b.X}
					y1 := []int{b.Y, b.Y, b.Y2, b.Y2}

					for i := 0; i < len(x1); i++ {
						r := NewRect(x1[i]-2, y1[i]-2, x1[i]+2, y1[i]+2)
						canvas.DrawFillRect(r, ImageBoxMarkPen, ImageBoxMarkBrush)
					}

				}
			}
		}

	case w32.WM_MOUSEMOVE:
		x, y := genPoint(lparam)

		if iv.dragBox != nil {
			if iv.resize == DirNone {
				iv.drag(iv.dragBox, x, y)
				iv.wasModified()
			} else {
				iv.resizeToDirection(iv.dragBox, x, y)
				iv.wasModified()
			}
			iv.Invalidate(true)

		} else {
			if !iv.add {
				w32.SetCursor(w32.LoadCursorWithResourceID(0, iv.getCursor(x, y)))
			}
			//  do not call repaint if underMouse item did not change.
			if iv.updateHighlight(x, y) {
				iv.Invalidate(true)
			}
		}

		if iv.mouseLeft {
			internalTrackMouseEvent(iv.hwnd)
			iv.mouseLeft = false
		}

	case w32.WM_MOUSELEAVE:
		iv.dragBox = nil
		iv.mouseLeft = true
		iv.updateHighlight(-1, -1)
		iv.Invalidate(true)

	case w32.WM_LBUTTONUP:
		iv.dragBox = nil

	case w32.WM_LBUTTONDOWN:
		x, y := genPoint(lparam)
		if iv.add {
			now := time.Now()
			s := fmt.Sprintf("field%s", now.Format("020405"))
			b := &ImageBox{Name: s, underMouse: true, X: x, Y: y, X2: x + 150, Y2: y + 30}
			iv.Boxes = append(iv.Boxes, b)
			iv.selBox = b
			iv.wasModified()
			iv.onAdd.Fire(NewEvent(iv, nil))
		} else {
			iv.dragBox = iv.isUnderMouse(x, y)
			iv.selBox = iv.dragBox
			iv.dragStartX, iv.dragStartY = x, y
			iv.resize = iv.resizingDirection(iv.dragBox, x, y)
		}
		iv.Invalidate(true)
		iv.onSelectedChange.Fire(NewEvent(iv, nil))

	case w32.WM_RBUTTONDOWN:

	}
	return w32.DefWindowProc(iv.hwnd, msg, wparam, lparam)
}
