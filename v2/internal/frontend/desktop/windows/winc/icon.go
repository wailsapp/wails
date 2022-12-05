//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type Icon struct {
	handle w32.HICON
}

func NewIconFromFile(path string) (*Icon, error) {
	ico := new(Icon)
	var err error
	if ico.handle = w32.LoadIcon(0, syscall.StringToUTF16Ptr(path)); ico.handle == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from %s", path))
	}
	return ico, err
}

func NewIconFromResource(instance w32.HINSTANCE, resId uint16) (*Icon, error) {
	ico := new(Icon)
	var err error
	if ico.handle = w32.LoadIconWithResourceID(instance, resId); ico.handle == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from resource with id %v", resId))
	}
	return ico, err
}

func ExtractIcon(fileName string, index int) (*Icon, error) {
	ico := new(Icon)
	var err error
	if ico.handle = w32.ExtractIcon(fileName, index); ico.handle == 0 || ico.handle == 1 {
		err = errors.New(fmt.Sprintf("Cannot extract icon from %s at index %v", fileName, index))
	}
	return ico, err
}

func (ic *Icon) Destroy() bool {
	return w32.DestroyIcon(ic.handle)
}

func (ic *Icon) Handle() w32.HICON {
	return ic.handle
}
