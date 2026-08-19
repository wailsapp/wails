package cache

import (
	"io/fs"
	"syscall"
)

func platformFileIdentity(info fs.FileInfo) (string, bool) {
	stat, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return "", false
	}
	created := uint64(stat.CreationTime.HighDateTime)<<32 | uint64(stat.CreationTime.LowDateTime)
	return encodeFileIdentity(created, uint64(stat.FileAttributes), uint64(stat.FileSizeHigh), uint64(stat.FileSizeLow)), true
}

func platformIdentityTracksChanges() bool { return false }
