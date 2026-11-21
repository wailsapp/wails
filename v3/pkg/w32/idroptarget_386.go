//go:build windows && 386

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

// IDropTarget::DragEnter/DragOver/Drop pass POINT by value. On 386 the
// eight-byte struct is pushed as two consecutive 32-bit stack slots (X then
// Y), while syscall.NewCallback rejects callback arguments larger than
// uintptr. The trampolines therefore declare one uintptr per slot and
// reassemble the POINT before dispatching to the Go implementation. The
// resulting argument list is six 32-bit slots, matching the stdcall stack
// layout and cleanup size of the original six-argument call.

func _iDropTargetDragEnter(
	this uintptr,
	dataObject *IDataObject,
	grfKeyState DWORD,
	pointX uintptr,
	pointY uintptr,
	pdfEffect *DWORD,
) uintptr {
	point := POINT{X: int32(pointX), Y: int32(pointY)}
	return combridge.Resolve[iDropTarget](this).DragEnter(dataObject, grfKeyState, point, pdfEffect)
}

func _iDropTargetDragOver(this uintptr, grfKeyState DWORD, pointX uintptr, pointY uintptr, pdfEffect *DWORD) uintptr {
	point := POINT{X: int32(pointX), Y: int32(pointY)}
	return combridge.Resolve[iDropTarget](this).DragOver(grfKeyState, point, pdfEffect)
}

func _iDropTargetDragLeave(this uintptr) uintptr {
	return combridge.Resolve[iDropTarget](this).DragLeave()
}

func _iDropTargetDrop(this uintptr, dataObject *IDataObject, grfKeyState DWORD, pointX uintptr, pointY uintptr, pdfEffect *DWORD) uintptr {
	point := POINT{X: int32(pointX), Y: int32(pointY)}
	return combridge.Resolve[iDropTarget](this).Drop(dataObject, grfKeyState, point, pdfEffect)
}
