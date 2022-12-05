//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

func knownFolderPath(id w32.CSIDL) (string, error) {
	var buf [w32.MAX_PATH]uint16

	if !w32.SHGetSpecialFolderPath(0, &buf[0], id, false) {
		return "", fmt.Errorf("SHGetSpecialFolderPath failed")
	}

	return syscall.UTF16ToString(buf[0:]), nil
}

func AppDataPath() (string, error) {
	return knownFolderPath(w32.CSIDL_APPDATA)
}

func CommonAppDataPath() (string, error) {
	return knownFolderPath(w32.CSIDL_COMMON_APPDATA)
}

func LocalAppDataPath() (string, error) {
	return knownFolderPath(w32.CSIDL_LOCAL_APPDATA)
}

// EnsureAppDataPath uses AppDataPath to ensure storage for local settings and databases.
func EnsureAppDataPath(company, product string) (string, error) {
	path, err := AppDataPath()
	if err != nil {
		return path, err
	}
	p := filepath.Join(path, company, product)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		// path/to/whatever does not exist
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			return p, err
		}
	}
	return p, nil
}

func DriveNames() ([]string, error) {
	bufLen := w32.GetLogicalDriveStrings(0, nil)
	if bufLen == 0 {
		return nil, fmt.Errorf("GetLogicalDriveStrings failed")
	}
	buf := make([]uint16, bufLen+1)

	bufLen = w32.GetLogicalDriveStrings(bufLen+1, &buf[0])
	if bufLen == 0 {
		return nil, fmt.Errorf("GetLogicalDriveStrings failed")
	}

	var names []string
	for i := 0; i < len(buf)-2; {
		name := syscall.UTF16ToString(buf[i:])
		names = append(names, name)
		i += len(name) + 1
	}
	return names, nil
}
