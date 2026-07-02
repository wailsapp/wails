//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// COREWEBVIEW2_PROCESS_FAILED_REASON describes why a WebView2 process failed
// (ICoreWebView2ProcessFailedEventArgs2.GetReason).
type COREWEBVIEW2_PROCESS_FAILED_REASON uint32

const (
	// An unexpected process failure occurred.
	COREWEBVIEW2_PROCESS_FAILED_REASON_UNEXPECTED = 0
	// The process became unresponsive.
	COREWEBVIEW2_PROCESS_FAILED_REASON_UNRESPONSIVE = 1
	// The process was terminated (e.g. from Task Manager or by another process).
	COREWEBVIEW2_PROCESS_FAILED_REASON_TERMINATED = 2
	// The process crashed.
	COREWEBVIEW2_PROCESS_FAILED_REASON_CRASHED = 3
	// The process failed to launch.
	COREWEBVIEW2_PROCESS_FAILED_REASON_LAUNCH_FAILED = 4
	// The process terminated due to running out of memory.
	COREWEBVIEW2_PROCESS_FAILED_REASON_OUT_OF_MEMORY = 5
	// The process's profile was deleted.
	COREWEBVIEW2_PROCESS_FAILED_REASON_PROFILE_DELETED = 6
)

// ICoreWebView2ProcessFailedEventArgs2 extends the ProcessFailed event args
// with the failure reason, exit code and process description — the diagnostics
// needed to tell a crash from an external kill, a launch failure or an
// out-of-memory condition. COM vtable order (after IUnknown and the base
// interface's GetProcessFailedKind): get_Reason, get_ExitCode,
// get_ProcessDescription, get_FrameInfosForFailedProcess.
type _ICoreWebView2ProcessFailedEventArgs2Vtbl struct {
	_IUnknownVtbl
	GetProcessFailedKind          ComProc
	GetReason                     ComProc
	GetExitCode                   ComProc
	GetProcessDescription         ComProc
	GetFrameInfosForFailedProcess ComProc
}

type ICoreWebView2ProcessFailedEventArgs2 struct {
	vtbl *_ICoreWebView2ProcessFailedEventArgs2Vtbl
}

// GetArgs2 queries the event args for ICoreWebView2ProcessFailedEventArgs2.
// Returns nil when the interface is unavailable (WebView2 runtime predating
// SDK 1.0.1054.31). The caller must Release() a non-nil result: QueryInterface
// adds a reference beyond the event args' own lifetime.
func (i *ICoreWebView2ProcessFailedEventArgs) GetArgs2() *ICoreWebView2ProcessFailedEventArgs2 {
	var result *ICoreWebView2ProcessFailedEventArgs2

	iid := NewGUID("{4dab9422-46fa-4c3e-a5d2-41d2071d3680}")
	i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2ProcessFailedEventArgs2) Release() {
	i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
}

// GetReason reports why the process failed. Browser-process exits always
// report UNEXPECTED; render-process-unresponsive always reports UNRESPONSIVE;
// other kinds may report any reason.
func (i *ICoreWebView2ProcessFailedEventArgs2) GetReason() (COREWEBVIEW2_PROCESS_FAILED_REASON, error) {
	var reason COREWEBVIEW2_PROCESS_FAILED_REASON
	hr, _, _ := i.vtbl.GetReason.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&reason)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return reason, nil
}

// GetExitCode reports the exit code of the failing process. It is always
// STILL_ACTIVE (259) when the kind is RENDER_PROCESS_UNRESPONSIVE.
func (i *ICoreWebView2ProcessFailedEventArgs2) GetExitCode() (int32, error) {
	var exitCode int32
	hr, _, _ := i.vtbl.GetExitCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&exitCode)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return exitCode, nil
}

// GetProcessDescription reports the WebView2 Runtime's technical description
// of the failed process (e.g. its utility subtype); empty for renderers.
func (i *ICoreWebView2ProcessFailedEventArgs2) GetProcessDescription() (string, error) {
	var _desc *uint16
	hr, _, _ := i.vtbl.GetProcessDescription.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_desc)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	desc := windows.UTF16PtrToString(_desc)
	windows.CoTaskMemFree(unsafe.Pointer(_desc))
	return desc, nil
}
