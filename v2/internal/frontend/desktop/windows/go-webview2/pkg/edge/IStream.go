//go:build windows

package edge

import (
	"io"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _IStreamVtbl struct {
	_IUnknownVtbl
	Read  ComProc
	Write ComProc
}

type IStream struct {
	vtbl *_IStreamVtbl
}

func (i *IStream) Release() error {
	return i.vtbl.CallRelease(unsafe.Pointer(i))
}

func (i *IStream) Read(p []byte) (int, error) {
	bufLen := len(p)
	if bufLen == 0 {
		return 0, nil
	}

	var n int
	res, _, err := i.vtbl.Read.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&p[0])),
		uintptr(bufLen),
		uintptr(unsafe.Pointer(&n)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0, err
	}

	switch windows.Handle(res) {
	case windows.S_OK:
		// The buffer has been completely filled
		return n, nil
	case windows.S_FALSE:
		// The buffer has been filled with less than len data and the stream is EOF
		return n, io.EOF
	default:
		return 0, syscall.Errno(res)
	}
}
