//go:build darwin

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// libproc is reached through purego, so the harness builds with and without
// cgo. Sizes and offsets are from <sys/proc_info.h> and <sys/resource.h>.
const (
	procPidTBSDInfo       = 3   // PROC_PIDTBSDINFO
	procBSDInfoSize       = 136 // PROC_PIDTBSDINFO_SIZE
	procBSDInfoNameAt     = 64  // offsetof(struct proc_bsdinfo, pbi_name)
	procBSDInfoNameLen    = 32
	rusageInfoCurrent     = 6   // RUSAGE_INFO_CURRENT (RUSAGE_INFO_V6)
	rusageInfoSize        = 464 // sizeof(rusage_info_current)
	rusagePhysFootprintAt = 72  // offsetof(rusage_info_current, ri_phys_footprint)
)

var (
	libprocOnce     sync.Once
	libprocErr      error
	procListAllPids func(buf unsafe.Pointer, size int32) int32
	procPidInfo     func(pid, flavor int32, arg uint64, buf unsafe.Pointer, size int32) int32
	procPidRusage   func(pid, flavor int32, buf unsafe.Pointer) int32
)

func loadLibproc() error {
	libprocOnce.Do(func() {
		lib, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			libprocErr = err
			return
		}
		purego.RegisterLibFunc(&procListAllPids, lib, "proc_listallpids")
		purego.RegisterLibFunc(&procPidInfo, lib, "proc_pidinfo")
		purego.RegisterLibFunc(&procPidRusage, lib, "proc_pid_rusage")
	})
	return libprocErr
}

// processName is pbi_name from PROC_PIDTBSDINFO.
func processName(pid int) (string, bool) {
	var info [procBSDInfoSize]byte
	if procPidInfo(int32(pid), procPidTBSDInfo, 0, unsafe.Pointer(&info[0]), procBSDInfoSize) != procBSDInfoSize {
		return "", false
	}
	name := info[procBSDInfoNameAt : procBSDInfoNameAt+procBSDInfoNameLen]
	if end := bytes.IndexByte(name, 0); end >= 0 {
		name = name[:end]
	}
	return string(name), true
}

const samplerSupported = true

const webContentProcName = "com.apple.WebKit.WebContent"

// webContentPids enumerates every WebContent process on the machine. Callers
// must diff against a pre-launch snapshot — summing all of them would fold in
// other applications' Safari tabs.
//
// proc_listallpids returns the number of PIDS written, not a byte count.
// (Verified: 907 returned against 907 pids present.) Dividing by sizeof(int)
// here would scan only a quarter of the table and could miss our own process.
func webContentPids() ([]int, error) {
	if err := loadLibproc(); err != nil {
		return nil, err
	}
	const maxPids = 32768
	buf := make([]int32, maxPids)
	n := int(procListAllPids(unsafe.Pointer(&buf[0]), int32(maxPids*unsafe.Sizeof(buf[0]))))
	if n <= 0 {
		return nil, fmt.Errorf("proc_listallpids failed (returned %d)", n)
	}
	if n > maxPids {
		n = maxPids
	}

	var out []int
	for i := 0; i < n; i++ {
		pid := int(buf[i])
		if pid <= 0 {
			continue
		}
		if name, ok := processName(pid); ok && name == webContentProcName {
			out = append(out, pid)
		}
	}
	return out, nil
}

// footprint returns ri_phys_footprint for pid. ok=false means the process is
// gone or unreadable — for a tracked WebContent pid that means it was replaced.
func footprint(pid int) (uint64, bool) {
	// ri_phys_footprint is what Activity Monitor shows and what jetsam acts
	// on. The harness's own RSS is not the signal; this is.
	if loadLibproc() != nil {
		return 0, false
	}
	var info [rusageInfoSize]byte
	if procPidRusage(int32(pid), rusageInfoCurrent, unsafe.Pointer(&info[0])) != 0 {
		return 0, false
	}
	return binary.LittleEndian.Uint64(info[rusagePhysFootprintAt:]), true
}
