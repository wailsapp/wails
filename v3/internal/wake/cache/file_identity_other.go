//go:build !linux && !darwin && !windows

package cache

import "io/fs"

func platformFileIdentity(fs.FileInfo) (string, bool) { return "", false }

func platformIdentityTracksChanges() bool { return false }
