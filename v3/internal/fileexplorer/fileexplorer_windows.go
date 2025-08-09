//go:build windows

package fileexplorer

import (
	"fmt"
	"syscall"
)

func explorerBinArgs(path string, selectFile bool) (string, []string, error) {
	return "explorer.exe", []string{}, nil
}

func sysProcAttr(path string, selectFile bool) *syscall.SysProcAttr {
	if selectFile {
		return &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf("explorer.exe /select,\"%s\"", path),
		}
	} else {
		return &syscall.SysProcAttr{
			CmdLine: fmt.Sprintf("explorer.exe \"%s\"", path),
		}
	}
}
