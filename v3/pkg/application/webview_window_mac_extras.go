package application

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// This file is the cross-platform surface of the macOS window extras: the
// represented file and proxy icon, the document-edited indicator, subtitles,
// window cascading, frame autosave, critical attention requests, print
// options, and PDF or image export of the WebView. The platform functions it
// calls live in webview_window_mac_extras_darwin.go, with documented no-op
// stubs in webview_window_mac_extras_other.go, mirroring
// webview_window_tabs.go.

var (
	// ErrMacOnly is returned by features that only exist on macOS when they
	// are called on another platform.
	ErrMacOnly = errors.New("this feature is only available on macOS")
	// ErrMacWindowNotCreated means the native window does not exist yet.
	// Export and print operations need a live NSWindow, so call them after
	// the window has been shown.
	ErrMacWindowNotCreated = errors.New("the native window has not been created yet")
	// ErrMacExportOnMainThread means ExportPDF or Snapshot was called on the
	// application thread. Both wait for an asynchronous WebKit completion
	// that is delivered on that thread, so waiting there would deadlock.
	// Call them from a goroutine.
	ErrMacExportOnMainThread = errors.New("ExportPDF and Snapshot cannot be called on the application thread")
	// ErrMacExportTimeout means WebKit did not deliver the export result
	// within the configured timeout.
	ErrMacExportTimeout = errors.New("timed out waiting for the WebView export")
	// ErrMacExportUnsupported means the running macOS release lacks the
	// WebKit API required for the export (PDF export needs macOS 11).
	ErrMacExportUnsupported = errors.New("the WebView export is unsupported on this macOS release")
	// ErrMacPrinterNotFound means PrintOptions.PrinterName named a printer
	// that is not installed.
	ErrMacPrinterNotFound = errors.New("the named printer was not found")
)

// PrintOrientation selects the page orientation for PrintWithOptions.
type PrintOrientation int

const (
	// PrintOrientationAutomatic keeps the orientation of the shared print
	// settings (normally portrait, or whatever the user chose last).
	PrintOrientationAutomatic PrintOrientation = iota
	// PrintOrientationPortrait prints in portrait.
	PrintOrientationPortrait
	// PrintOrientationLandscape prints in landscape.
	PrintOrientationLandscape
)

func validPrintOrientation(orientation PrintOrientation) bool {
	switch orientation {
	case PrintOrientationAutomatic, PrintOrientationPortrait, PrintOrientationLandscape:
		return true
	}
	return false
}

// PrintMargins are page margins in points.
type PrintMargins struct {
	Top    float64
	Left   float64
	Bottom float64
	Right  float64
}

// isZero reports whether no margin was requested, in which case the print
// settings' own margins are kept.
func (m PrintMargins) isZero() bool {
	return m.Top == 0 && m.Left == 0 && m.Bottom == 0 && m.Right == 0
}

// PrintOptions configures WebviewWindow.PrintWithOptions. The zero value
// prints with the shared print settings and shows the print panel.
type PrintOptions struct {
	// Orientation is the page orientation. Automatic keeps the shared
	// print settings' orientation.
	Orientation PrintOrientation
	// Margins are the page margins in points. When every field is zero the
	// shared print settings' margins are kept. Negative values are rejected.
	Margins PrintMargins
	// Silent prints without the print panel and progress panel, using the
	// selected or default printer.
	Silent bool
	// PrinterName selects an installed printer by name. Empty keeps the
	// default printer. An unknown name fails with ErrMacPrinterNotFound.
	PrinterName string
	// Scale is the scaling factor (1 is 100%). Zero keeps the print
	// settings' scale. Negative values are rejected.
	Scale float64
	// PaperName selects a paper by its PostScript name, such as "na-letter"
	// or "iso-a4". Empty keeps the current paper.
	PaperName string
}

// validate checks the options and returns the normalised copy.
func (o PrintOptions) validate() (PrintOptions, error) {
	if !validPrintOrientation(o.Orientation) {
		return o, fmt.Errorf("unknown print orientation %d", o.Orientation)
	}
	if o.Margins.Top < 0 || o.Margins.Left < 0 || o.Margins.Bottom < 0 || o.Margins.Right < 0 {
		return o, errors.New("print margins cannot be negative")
	}
	if o.Scale < 0 {
		return o, errors.New("print scale cannot be negative")
	}
	return o, nil
}

// legacyPrintOptions are the settings WebviewWindow.Print has always used
// on macOS: landscape, 30 point margins, print and progress panels shown.
func legacyPrintOptions() PrintOptions {
	return PrintOptions{
		Orientation: PrintOrientationLandscape,
		Margins:     PrintMargins{Top: 30, Left: 30, Bottom: 30, Right: 30},
	}
}

// DefaultMacExportTimeout is how long ExportPDF and Snapshot wait for WebKit
// when the options carry no Timeout.
const DefaultMacExportTimeout = 30 * time.Second

// PDFExportOptions configures WebviewWindow.ExportPDF.
type PDFExportOptions struct {
	// Rect limits the export to a region of the WebView in points, with the
	// origin at the top left. Nil exports the whole page.
	Rect *Rect
	// Timeout bounds the wait for WebKit. Zero uses DefaultMacExportTimeout.
	Timeout time.Duration
}

// SnapshotOptions configures WebviewWindow.Snapshot.
type SnapshotOptions struct {
	// Rect limits the snapshot to a region of the WebView in points, with
	// the origin at the top left. Nil captures the visible viewport.
	Rect *Rect
	// Width scales the snapshot to this width in points, keeping the aspect
	// ratio. Zero keeps the natural size.
	Width int
	// Timeout bounds the wait for WebKit. Zero uses DefaultMacExportTimeout.
	Timeout time.Duration
}

func exportTimeout(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return DefaultMacExportTimeout
	}
	return timeout
}

// MacAttentionRequest is the handle returned by WebviewWindow.RequestAttention.
// A nil handle is safe to use.
type MacAttentionRequest struct {
	id       int
	once     sync.Once
	critical bool
}

// Critical reports whether the request was made with critical set.
func (r *MacAttentionRequest) Critical() bool {
	return r != nil && r.critical
}

// Cancel withdraws the attention request, stopping a critical Dock bounce
// that has not yet been acknowledged. It is a no-op once the request has
// completed, off macOS, and on a nil handle.
func (r *MacAttentionRequest) Cancel() {
	if r == nil {
		return
	}
	r.once.Do(func() {
		if r.id != 0 {
			macWindowExtrasCancelAttention(r.id)
		}
	})
}

// macWindowExtrasState holds the values set through the represented file,
// document edited and subtitle setters before the native window exists.
// They are applied on the application thread when the window is created.
type macWindowExtrasState struct {
	representedFile *string
	documentEdited  *bool
	subtitle        *string
}

var (
	macWindowExtrasPendingLock sync.Mutex
	macWindowExtrasPending     = map[*WebviewWindow]*macWindowExtrasState{}
)

func macWindowExtrasPendingFor(w *WebviewWindow) *macWindowExtrasState {
	state := macWindowExtrasPending[w]
	if state == nil {
		state = &macWindowExtrasState{}
		macWindowExtrasPending[w] = state
	}
	return state
}

// takeMacWindowExtrasPending removes and returns the values queued for a
// window, or nil when nothing was queued.
func takeMacWindowExtrasPending(w *WebviewWindow) *macWindowExtrasState {
	macWindowExtrasPendingLock.Lock()
	defer macWindowExtrasPendingLock.Unlock()
	state := macWindowExtrasPending[w]
	delete(macWindowExtrasPending, w)
	return state
}

// SetRepresentedFile associates the window with a file on disk: the titlebar
// shows the file's proxy icon, which can be dragged and command-clicked to
// reveal the file's path. Pass "" to remove the association. The title is
// unchanged; call SetTitle to show the file name. Before the native window
// exists the value is applied when the window is created. It is a no-op off
// macOS.
func (w *WebviewWindow) SetRepresentedFile(path string) Window {
	if w == nil {
		return w
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		macWindowExtrasPendingLock.Lock()
		macWindowExtrasPendingFor(w).representedFile = &path
		macWindowExtrasPendingLock.Unlock()
		return w
	}
	macWindowExtrasSetRepresentedFile(nsWindow, path)
	return w
}

// RepresentedFile returns the path set with SetRepresentedFile, or "" when
// the window represents no file or the native window does not exist yet.
func (w *WebviewWindow) RepresentedFile() string {
	if w == nil {
		return ""
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		macWindowExtrasPendingLock.Lock()
		defer macWindowExtrasPendingLock.Unlock()
		if state := macWindowExtrasPending[w]; state != nil && state.representedFile != nil {
			return *state.representedFile
		}
		return ""
	}
	return macWindowExtrasRepresentedFile(nsWindow)
}

// SetDocumentEdited marks the window's document as having unsaved changes:
// the close button shows a dot and the proxy icon dims. Before the native
// window exists the value is applied when the window is created. It is a
// no-op off macOS.
func (w *WebviewWindow) SetDocumentEdited(edited bool) Window {
	if w == nil {
		return w
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		macWindowExtrasPendingLock.Lock()
		macWindowExtrasPendingFor(w).documentEdited = &edited
		macWindowExtrasPendingLock.Unlock()
		return w
	}
	macWindowExtrasSetDocumentEdited(nsWindow, edited)
	return w
}

// IsDocumentEdited reports whether the window is marked as having unsaved
// changes. It is false off macOS.
func (w *WebviewWindow) IsDocumentEdited() bool {
	if w == nil {
		return false
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		macWindowExtrasPendingLock.Lock()
		defer macWindowExtrasPendingLock.Unlock()
		if state := macWindowExtrasPending[w]; state != nil && state.documentEdited != nil {
			return *state.documentEdited
		}
		return false
	}
	return macWindowExtrasIsDocumentEdited(nsWindow)
}

// SetSubtitle shows a secondary line under the title (macOS 11+). Pass "" to
// remove it. On earlier releases the call is ignored with a debug log
// entry; it is a no-op off macOS. Before the native window exists the value
// is applied when the window is created.
func (w *WebviewWindow) SetSubtitle(subtitle string) Window {
	if w == nil {
		return w
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		macWindowExtrasPendingLock.Lock()
		macWindowExtrasPendingFor(w).subtitle = &subtitle
		macWindowExtrasPendingLock.Unlock()
		return w
	}
	w.applySubtitle(nsWindow, subtitle)
	return w
}

func (w *WebviewWindow) applySubtitle(nsWindow unsafe.Pointer, subtitle string) {
	if !macWindowExtrasSetSubtitle(nsWindow, subtitle) && globalApplication != nil {
		globalApplication.debug("SetSubtitle requires macOS 11 or later; ignored", "sender", w.Name())
	}
}

// CascadeFrom positions this window so it cascades from other, offset down
// and to the right of other's top-left corner like a new document window,
// and records the new point as the application's cascade point so windows
// created with InitialPosition WindowCascade continue from it. Both native
// windows must exist. It is a no-op off macOS.
func (w *WebviewWindow) CascadeFrom(other Window) Window {
	if w == nil || other == nil {
		return w
	}
	nsWindow := w.NativeWindow()
	from := other.NativeWindow()
	if nsWindow == nil || from == nil || nsWindow == from {
		return w
	}
	macWindowExtrasCascadeFrom(nsWindow, from)
	return w
}

// SetFrameAutosaveName saves the window's frame under name in the user
// defaults whenever it moves or resizes, and restores a frame previously
// saved under that name immediately. Pass "" to stop saving. To restore a
// saved frame before the window is first shown, set
// MacWindow.FrameAutosaveName in the options instead. It requires the native
// window to exist and is a no-op off macOS.
func (w *WebviewWindow) SetFrameAutosaveName(name string) Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macWindowExtrasSetFrameAutosaveName(nsWindow, name)
	}
	return w
}

// SetWindowButtonsOffset moves the close, minimise and zoom buttons by x and
// y points from their default position, overriding
// MacTitleBar.WindowButtonsOffset. The offset is reapplied whenever AppKit
// lays the titlebar out again. It requires the native window to exist and
// is a no-op off macOS.
func (w *WebviewWindow) SetWindowButtonsOffset(x, y int) Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macWindowExtrasSetWindowButtonsOffset(nsWindow, x, y, true)
	}
	return w
}

// ResetWindowButtonsOffset restores the default window button position.
func (w *WebviewWindow) ResetWindowButtonsOffset() Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macWindowExtrasSetWindowButtonsOffset(nsWindow, 0, 0, false)
	}
	return w
}

// RequestAttention asks for the user's attention while the application is
// in the background. An informational request bounces the Dock icon once;
// a critical request keeps bouncing until the application is activated or
// the returned handle is cancelled. Flash is the cross-platform
// informational form; this one exposes the critical variant and the cancel
// handle. Off macOS it returns a handle whose Cancel is a no-op.
func (w *WebviewWindow) RequestAttention(critical bool) *MacAttentionRequest {
	request := &MacAttentionRequest{critical: critical}
	if w == nil {
		return request
	}
	request.id = macWindowExtrasRequestAttention(critical)
	return request
}

// PrintWithOptions prints the WebView's content with the given settings.
// Print keeps its historical behaviour (landscape, 30 point margins, panel
// shown) and shares this implementation. Off macOS the options are ignored
// and Print is called.
func (w *WebviewWindow) PrintWithOptions(options PrintOptions) error {
	if w == nil {
		return ErrMacWindowNotCreated
	}
	options, err := options.validate()
	if err != nil {
		return err
	}
	if !macWindowExtrasSupported() {
		return w.Print()
	}
	if w.impl == nil || w.isDestroyed() {
		return ErrMacWindowNotCreated
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return ErrMacWindowNotCreated
	}
	return macWindowExtrasPrint(nsWindow, options)
}

// ExportPDF renders the WebView's current content to a PDF document
// (macOS 11+). It waits for WebKit without blocking the application
// thread, so it must be called from a goroutine, never from the
// application thread itself. Off macOS it returns ErrMacOnly.
func (w *WebviewWindow) ExportPDF(options PDFExportOptions) ([]byte, error) {
	nsWindow, err := w.exportTarget()
	if err != nil {
		return nil, err
	}
	return macWindowExtrasExportPDF(nsWindow, options)
}

// Snapshot captures the WebView's rendered content as a PNG image. It waits
// for WebKit without blocking the application thread, so it must be called
// from a goroutine, never from the application thread itself. Off macOS it
// returns ErrMacOnly.
func (w *WebviewWindow) Snapshot(options SnapshotOptions) ([]byte, error) {
	nsWindow, err := w.exportTarget()
	if err != nil {
		return nil, err
	}
	return macWindowExtrasSnapshot(nsWindow, options)
}

func (w *WebviewWindow) exportTarget() (unsafe.Pointer, error) {
	if !macWindowExtrasSupported() {
		return nil, ErrMacOnly
	}
	if w == nil || w.impl == nil || w.isDestroyed() {
		return nil, ErrMacWindowNotCreated
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return nil, ErrMacWindowNotCreated
	}
	if globalApplication != nil && globalApplication.impl != nil && globalApplication.impl.isOnMainThread() {
		return nil, ErrMacExportOnMainThread
	}
	return nsWindow, nil
}

// normaliseInitialPosition maps placement modes a platform cannot honour to
// the closest one it can. WindowCascade becomes WindowCentered off macOS.
func normaliseInitialPosition(position WindowStartPosition) WindowStartPosition {
	if position == WindowCascade && !macWindowExtrasSupported() {
		return WindowCentered
	}
	return position
}
