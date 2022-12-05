//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type Pen struct {
	hPen  w32.HPEN
	style uint
	brush *Brush
}

func NewPen(style uint, width uint, brush *Brush) *Pen {
	if brush == nil {
		panic("Brush cannot be nil")
	}

	hPen := w32.ExtCreatePen(style, width, brush.GetLOGBRUSH(), 0, nil)
	if hPen == 0 {
		panic("Failed to create pen")
	}

	return &Pen{hPen, style, brush}
}

func NewNullPen() *Pen {
	lb := w32.LOGBRUSH{LbStyle: w32.BS_NULL}

	hPen := w32.ExtCreatePen(w32.PS_COSMETIC|w32.PS_NULL, 1, &lb, 0, nil)
	if hPen == 0 {
		panic("failed to create null brush")
	}

	return &Pen{hPen: hPen}
}

func (pen *Pen) Style() uint {
	return pen.style
}

func (pen *Pen) Brush() *Brush {
	return pen.brush
}

func (pen *Pen) GetHPEN() w32.HPEN {
	return pen.hPen
}

func (pen *Pen) Dispose() {
	if pen.hPen != 0 {
		w32.DeleteObject(w32.HGDIOBJ(pen.hPen))
		pen.hPen = 0
	}
}
