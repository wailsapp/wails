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

func (i *IStream) AddRef() uint32 {
	// The return value is the new reference count, as documented for IUnknown::AddRef.
	// It is for diagnostics only. The Errno returned by Call is just the thread's
	// last-error and says nothing about success, so it is ignored.
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *IStream) Release() uint32 {
	// The return value is the new reference count, as documented for IUnknown::Release.
	// It is for diagnostics only. The Errno returned by Call is just the thread's
	// last-error and says nothing about success, so it is ignored.
	r, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(r)
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
