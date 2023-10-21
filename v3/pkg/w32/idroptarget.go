//go:build windows

package w32

import (
	"github.com/wailsapp/go-webview2/pkg/combridge"
	"golang.org/x/sys/windows"
)

const (
	DROPEFFECT_NONE = 0
	DROPEFFECT_COPY = 1
	DROPEFFECT_MOVE = 2
	DROPEFFECT_LINK = 4
)

func _NOP(_ uintptr) uintptr {
	return uintptr(windows.S_FALSE)
}

func init() {
	combridge.RegisterVTable[combridge.IUnknown, iDropTarget](
		"{00000122-0000-0000-C000-000000000046}",
		_iDropTargetDragEnter,
		_iDropTargetDragOver,
		_iDropTargetDragLeave,
		_iDropTargetDrop,
		_NOP,
	)
}

func _iDropTargetDragEnter(this uintptr, dataObject *IDataObject, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr {
	return combridge.Resolve[iDropTarget](this).DragEnter(dataObject, grfKeyState, point, pdfEffect)
}

func _iDropTargetDragOver(this uintptr, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr {
	return combridge.Resolve[iDropTarget](this).DragOver(grfKeyState, point, pdfEffect)
}

func _iDropTargetDragLeave(this uintptr) uintptr {
	return combridge.Resolve[iDropTarget](this).DragLeave()
}

func _iDropTargetDrop(this uintptr, dataObject *IDataObject, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr {
	return combridge.Resolve[iDropTarget](this).Drop(dataObject, grfKeyState, point, pdfEffect)
}

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
	OnDrop        func(filenames []string)
}

func NewDropTarget() *DropTarget {
	return &DropTarget{
		OnEnterEffect: DROPEFFECT_COPY,
		OnOverEffect:  DROPEFFECT_COPY,
	}
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

	d.OnDrop(filenames)

	return uintptr(windows.S_OK)
}
