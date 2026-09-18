//go:build !darwin || ios || server

package application

import "unsafe"

// Off macOS the window extras are documented no-ops; the exported methods in
// webview_window_mac_extras.go return ErrMacOnly where an error is possible.

func macWindowExtrasSupported() bool { return false }

func macWindowExtrasSetRepresentedFile(unsafe.Pointer, string)   {}
func macWindowExtrasRepresentedFile(unsafe.Pointer) string       { return "" }
func macWindowExtrasSetDocumentEdited(unsafe.Pointer, bool)      {}
func macWindowExtrasIsDocumentEdited(unsafe.Pointer) bool        { return false }
func macWindowExtrasSetSubtitle(unsafe.Pointer, string) bool     { return true }
func macWindowExtrasCascadeFrom(unsafe.Pointer, unsafe.Pointer)  {}
func macWindowExtrasSetFrameAutosaveName(unsafe.Pointer, string) {}
func macWindowExtrasSetWindowButtonsOffset(unsafe.Pointer, int, int, bool) {
}
func macWindowExtrasRequestAttention(bool) int { return 0 }
func macWindowExtrasCancelAttention(int)       {}

func macWindowExtrasPrint(unsafe.Pointer, PrintOptions) error { return ErrMacOnly }

func macWindowExtrasExportPDF(unsafe.Pointer, PDFExportOptions) ([]byte, error) {
	return nil, ErrMacOnly
}

func macWindowExtrasSnapshot(unsafe.Pointer, SnapshotOptions) ([]byte, error) {
	return nil, ErrMacOnly
}
