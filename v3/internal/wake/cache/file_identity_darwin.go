package cache

import (
	"io/fs"
	"syscall"
)

func platformFileIdentity(info fs.FileInfo) (string, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", false
	}
	return encodeFileIdentity(uint64(stat.Dev), stat.Ino, uint64(stat.Ctimespec.Sec), uint64(stat.Ctimespec.Nsec)), true
}

func platformIdentityTracksChanges() bool { return true }
