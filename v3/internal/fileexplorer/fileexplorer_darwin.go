//go:build darwin

package fileexplorer

import "syscall"

func explorerBinArgs(path string, selectFile bool) (string, []string, error) {
	args := []string{}
	if selectFile {
		args = append(args, "-R")
	}

	args = append(args, path)
	return "open", args, nil
}

func sysProcAttr(path string, selectFile bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
