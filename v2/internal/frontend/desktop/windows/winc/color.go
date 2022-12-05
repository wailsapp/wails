//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

type Color uint32

func RGB(r, g, b byte) Color {
	return Color(uint32(r) | uint32(g)<<8 | uint32(b)<<16)
}

func (c Color) R() byte {
	return byte(c & 0xff)
}

func (c Color) G() byte {
	return byte((c >> 8) & 0xff)
}

func (c Color) B() byte {
	return byte((c >> 16) & 0xff)
}
