//go:build windows

package w32

import (
	"github.com/wailsapp/wails/v3/internal/webview2/pkg/combridge"
	"golang.org/x/sys/windows"
)

var (
	DROPEFFECT_NONE DWORD = 0
	DROPEFFECT_COPY DWORD = 1
	DROPEFFECT_MOVE DWORD = 2
	DROPEFFECT_LINK DWORD = 4
)

const (
	DRAGDROP_E_ALREADYREGISTERED = 0x80040101
	DRAGDROP_E_INVALIDHWND       = 0x80040102
)

func _NOP(_ uintptr) uintptr {
	return uintptr(windows.S_FALSE)
}

// The IDropTarget vtable trampolines live in architecture-specific files:
// syscall.NewCallback rejects callback arguments larger than uintptr, and
// IDropTarget::DragEnter/DragOver/Drop pass POINT by value, which is eight
// bytes. On 64-bit architectures that matches uintptr, but on 386 the
// eight-byte struct is pushed as two 32-bit stack slots and must be
// reassembled from two uintptr arguments instead.
//
//	idroptarget_64bit.go: windows && (amd64 || arm64)
//	idroptarget_386.go:   windows && 386

type iDropTarget interface {
	combridge.IUnknown

	DragEnter(dataObject *IDataObject, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr
	DragOver(grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr
	DragLeave() uintptr
	Drop(dataObject *IDataObject, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr
}

var _ iDropTarget = &DropTarget{}

type DropTarget struct {
	combridge.IUnknownImpl
	OnEnterEffect DWORD
	OnOverEffect  DWORD
	OnEnter       func()
	OnLeave       func()
	OnOver        func()
	OnDrop        func(filenames []string, x int, y int)
}

func NewDropTarget() *DropTarget {
	result := &DropTarget{
		OnEnterEffect: DROPEFFECT_COPY,
		OnOverEffect:  DROPEFFECT_COPY,
	}
	return result
}

func (d *DropTarget) DragEnter(dataObject *IDataObject, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr {
	*pdfEffect = d.OnEnterEffect
	if d.OnEnter != nil {
		d.OnEnter()
	}
	return uintptr(windows.S_OK)
}

func (d *DropTarget) DragOver(grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr {
	*pdfEffect = d.OnOverEffect
	if d.OnOver != nil {
		d.OnOver()
	}
	return uintptr(windows.S_OK)
}

func (d *DropTarget) DragLeave() uintptr {
	if d.OnLeave != nil {
		d.OnLeave()
	}
	return uintptr(windows.S_OK)
}

func (d *DropTarget) Drop(dataObject *IDataObject, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr {

	if d.OnDrop == nil {
		return uintptr(windows.S_OK)
	}

	// Extract filenames from dataObject
	var filenames []string
	var formatETC = FORMATETC{
		CfFormat: CF_HDROP,
		Tymed:    TYMED_HGLOBAL,
	}

	var stgMedium STGMEDIUM

	err := dataObject.GetData(&formatETC, &stgMedium)
	if err != nil && err != windows.ERROR_SUCCESS {
		return uintptr(windows.S_FALSE)
	}
	defer stgMedium.Release()
	hDrop := stgMedium.Union
	_, numFiles := DragQueryFile(hDrop, 0xFFFFFFFF)
	for i := uint(0); i < numFiles; i++ {
		filename, _ := DragQueryFile(hDrop, i)
		filenames = append(filenames, filename)
	}

	d.OnDrop(filenames, int(point.X), int(point.Y))

	return uintptr(windows.S_OK)
}
