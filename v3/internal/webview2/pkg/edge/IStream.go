//go:build windows

package edge

import (
	"errors"
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

func (i *IStream) AddRef() error {
	_, _, err := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	if err != nil && !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}

	return nil
}

func (i *IStream) Release() error {
	_, _, err := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	if err != nil && !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}

	return nil
}

func (i *IStream) Read(p []byte) (int, error) {
	bufLen := len(p)
	if bufLen == 0 {
		return 0, nil
	}

	var n int
	hr, _, _ := i.vtbl.Read.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&p[0])),
		uintptr(bufLen),
		uintptr(unsafe.Pointer(&n)),
	)

	switch windows.Handle(hr) {
	case windows.S_OK:
		// The buffer has been completely filled
		return n, nil
	case windows.S_FALSE:
		// The buffer has been filled with less than len data and the stream is EOF
		return n, io.EOF
	default:
		return 0, syscall.Errno(hr)
	}
}
