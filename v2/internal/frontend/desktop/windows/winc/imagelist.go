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

type ImageList struct {
	handle w32.HIMAGELIST
}

func NewImageList(cx, cy int) *ImageList {
	return newImageList(cx, cy, w32.ILC_COLOR32, 0, 0)
}

func newImageList(cx, cy int, flags uint, cInitial, cGrow int) *ImageList {
	imgl := new(ImageList)
	imgl.handle = w32.ImageList_Create(cx, cy, flags, cInitial, cGrow)
	return imgl
}

func (im *ImageList) Handle() w32.HIMAGELIST {
	return im.handle
}

func (im *ImageList) Destroy() bool {
	return w32.ImageList_Destroy(im.handle)
}

func (im *ImageList) SetImageCount(uNewCount uint) bool {
	return w32.ImageList_SetImageCount(im.handle, uNewCount)
}

func (im *ImageList) ImageCount() int {
	return w32.ImageList_GetImageCount(im.handle)
}

func (im *ImageList) AddIcon(icon *Icon) int {
	return w32.ImageList_AddIcon(im.handle, icon.Handle())
}

func (im *ImageList) AddResIcon(iconID uint16) {
	if ico, err := NewIconFromResource(GetAppInstance(), iconID); err == nil {
		im.AddIcon(ico)
		return
	}
	panic(fmt.Sprintf("missing icon with icon ID: %d", iconID))
}

func (im *ImageList) RemoveAll() bool {
	return w32.ImageList_RemoveAll(im.handle)
}

func (im *ImageList) Remove(i int) bool {
	return w32.ImageList_Remove(im.handle, i)
}
