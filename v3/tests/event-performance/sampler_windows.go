//go:build windows

package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

const samplerSupported = true

// WebView2 runs page content in msedgewebview2.exe. Unlike WKWebView and
// WebKitGTK there is not a single content process: an app gets a browser
// process plus renderer, GPU and utility children, all with this name. The
// pre-launch diff in main.go handles that — every pid that appears after our
// window is created belongs to us, and they are summed.
const webContentProcName = "msedgewebview2.exe"

// PROCESS_MEMORY_COUNTERS_EX. PrivateUsage (the commit charge) is the closest
// Windows analogue to macOS phys_footprint: it counts private committed bytes
// and excludes shared pages, so it does not double count the mapped images
// every WebView2 child shares.
type processMemoryCountersEx struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
	PrivateUsage               uintptr
}

var (
	modPsapi                 = windows.NewLazySystemDLL("psapi.dll")
	procGetProcessMemoryInfo = modPsapi.NewProc("GetProcessMemoryInfo")
)

// webContentPids enumerates every WebView2 process on the machine. Callers
// diff against a pre-launch snapshot; a dev box typically has a dozen of these
// already running for other apps.
func webContentPids() ([]int, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot: %w", err)
	}
	defer windows.CloseHandle(snap)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(snap, &entry); err != nil {
		return nil, fmt.Errorf("Process32First: %w", err)
	}

	var out []int
	for {
		name := windows.UTF16ToString(entry.ExeFile[:])
		if strings.EqualFold(filepath.Base(name), webContentProcName) {
			out = append(out, int(entry.ProcessID))
		}
		if err := windows.Process32Next(snap, &entry); err != nil {
			break // ERROR_NO_MORE_FILES
		}
	}
	return out, nil
}

// footprint reports PrivateUsage for pid. ok=false means the process is gone or
// unreadable — for a tracked content process that means it was replaced.
func footprint(pid int) (uint64, bool) {
	h, err := windows.OpenProcess(
		windows.PROCESS_QUERY_LIMITED_INFORMATION|windows.PROCESS_VM_READ,
		false, uint32(pid))
	if err != nil {
		// Fall back to the narrower right; VM_READ is refused for some
		// protected processes even at the same integrity level.
		h, err = windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
		if err != nil {
			return 0, false
		}
	}
	defer windows.CloseHandle(h)

	var counters processMemoryCountersEx
	counters.CB = uint32(unsafe.Sizeof(counters))
	r, _, _ := procGetProcessMemoryInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&counters)),
		uintptr(counters.CB),
	)
	if r == 0 {
		return 0, false
	}
	return uint64(counters.PrivateUsage), true
}
