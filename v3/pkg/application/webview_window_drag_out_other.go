//go:build !darwin || ios || server || wails_native

package application

func (w *WebviewWindow) startDragOut(DragItems) error {
	return ErrDragOutUnsupported
}
