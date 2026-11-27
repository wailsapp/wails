//go:build windows

package w32

import (
	"github.com/wailsapp/wails/v3/internal/webview2/pkg/combridge"
)

// 386: stdcall pushes the by-value POINTL argument as two consecutive 4-byte
// stack slots. syscall.NewCallback rejects any parameter wider than a uintptr.

func _iDropTargetDragEnter(
	this uintptr,
	dataObject *IDataObject,
	grfKeyState DWORD,
	pointX, pointY int32,
	pdfEffect *DWORD,
) uintptr {
	point := POINT{X: pointX, Y: pointY}
	return combridge.Resolve[iDropTarget](this).DragEnter(dataObject, grfKeyState, point, pdfEffect)
}

func _iDropTargetDragOver(this uintptr, grfKeyState DWORD, pointX, pointY int32, pdfEffect *DWORD) uintptr {
	point := POINT{X: pointX, Y: pointY}
	return combridge.Resolve[iDropTarget](this).DragOver(grfKeyState, point, pdfEffect)
}

func _iDropTargetDrop(this uintptr, dataObject *IDataObject, grfKeyState DWORD, pointX, pointY int32, pdfEffect *DWORD) uintptr {
	point := POINT{X: pointX, Y: pointY}
	return combridge.Resolve[iDropTarget](this).Drop(dataObject, grfKeyState, point, pdfEffect)
}
