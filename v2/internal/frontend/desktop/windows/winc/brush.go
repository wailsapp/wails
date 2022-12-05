//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

var DefaultBackgroundBrush = NewSystemColorBrush(w32.COLOR_BTNFACE)

type Brush struct {
	hBrush   w32.HBRUSH
	logBrush w32.LOGBRUSH
}

func NewSolidColorBrush(color Color) *Brush {
	lb := w32.LOGBRUSH{LbStyle: w32.BS_SOLID, LbColor: w32.COLORREF(color)}
	hBrush := w32.CreateBrushIndirect(&lb)
	if hBrush == 0 {
		panic("Faild to create solid color brush")
	}

	return &Brush{hBrush, lb}
}

func NewSystemColorBrush(colorIndex int) *Brush {
	//lb := w32.LOGBRUSH{LbStyle: w32.BS_SOLID, LbColor: w32.COLORREF(colorIndex)}
	lb := w32.LOGBRUSH{LbStyle: w32.BS_NULL}
	hBrush := w32.GetSysColorBrush(colorIndex)
	if hBrush == 0 {
		panic("GetSysColorBrush failed")
	}
	return &Brush{hBrush, lb}
}

func NewHatchedColorBrush(color Color) *Brush {
	lb := w32.LOGBRUSH{LbStyle: w32.BS_HATCHED, LbColor: w32.COLORREF(color)}
	hBrush := w32.CreateBrushIndirect(&lb)
	if hBrush == 0 {
		panic("Faild to create solid color brush")
	}

	return &Brush{hBrush, lb}
}

func NewNullBrush() *Brush {
	lb := w32.LOGBRUSH{LbStyle: w32.BS_NULL}
	hBrush := w32.CreateBrushIndirect(&lb)
	if hBrush == 0 {
		panic("Failed to create null brush")
	}

	return &Brush{hBrush, lb}
}

func (br *Brush) GetHBRUSH() w32.HBRUSH {
	return br.hBrush
}

func (br *Brush) GetLOGBRUSH() *w32.LOGBRUSH {
	return &br.logBrush
}

func (br *Brush) Dispose() {
	if br.hBrush != 0 {
		w32.DeleteObject(w32.HGDIOBJ(br.hBrush))
		br.hBrush = 0
	}
}
