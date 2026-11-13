//go:build windows && !386

package w32

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/webview2/pkg/combridge"
)

// Safety net compile-time assertion POINT must be exactly one word wide.
const _ = uint(unsafe.Sizeof(uintptr(0)) - unsafe.Sizeof(POINT{}))

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

func _iDropTargetDrop(this uintptr, dataObject *IDataObject, grfKeyState DWORD, point POINT, pdfEffect *DWORD) uintptr {
	return combridge.Resolve[iDropTarget](this).Drop(dataObject, grfKeyState, point, pdfEffect)
}
