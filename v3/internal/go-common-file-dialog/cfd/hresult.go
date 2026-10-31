package cfd

import "github.com/go-ole/go-ole"

func hresultToError(hr uintptr) error {
	// HRESULT is signed 32-bit, even when a syscall returns a 64-bit uintptr.
	code := uint32(hr)
	if code == 0x800704c7 { // HRESULT_FROM_WIN32(ERROR_CANCELLED)
		return ErrorCancelled
	}
	if int32(code) < 0 {
		return ole.NewError(uintptr(code))
	}
	return nil
}
