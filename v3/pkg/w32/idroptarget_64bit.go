//go:build windows && (amd64 || arm64)

package w32

import (
	"github.com/wailsapp/wails/v3/internal/webview2/pkg/combridge"
)

func init() {
	combridge.RegisterVTable[combridge.IUnknown, iDropTarget](
		"{00000122-0000-0000-C000-000000000046}",
		_iDropTargetDragEnter,
		_iDropTargetDragOver,
		_iDropTargetDragLeave,
		_iDropTargetDrop,
	)
}

// On 64-bit architectures POINT is eight bytes, matching uintptr, so the
// by-value POINT arguments of IDropTarget::DragEnter/DragOver/Drop can be
// declared as plain struct values in the trampolines.

func _iDropTargetDragEnter(
	this uintptr,
	dataObject *IDataObject,
	grfKeyState DWORD,
	point POINT,
	pdfEffect *DWORD,
) uintptr {
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
