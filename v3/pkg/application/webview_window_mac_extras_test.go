package application

import (
	"errors"
	"testing"
	"time"
)

func TestWindowCascadeIsADistinctStartPosition(t *testing.T) {
	if WindowCascade == WindowCentered || WindowCascade == WindowXY {
		t.Fatalf("WindowCascade (%d) collides with an existing start position", WindowCascade)
	}
	if WindowCentered != 0 {
		t.Fatal("WindowCentered must stay the zero value")
	}
}

func TestNormaliseInitialPositionFollowsPlatform(t *testing.T) {
	if got := normaliseInitialPosition(WindowCentered); got != WindowCentered {
		t.Fatalf("WindowCentered normalised to %d", got)
	}
	if got := normaliseInitialPosition(WindowXY); got != WindowXY {
		t.Fatalf("WindowXY normalised to %d", got)
	}
	got := normaliseInitialPosition(WindowCascade)
	if macWindowExtrasSupported() && got != WindowCascade {
		t.Fatalf("WindowCascade normalised to %d on macOS", got)
	}
	if !macWindowExtrasSupported() && got != WindowCentered {
		t.Fatalf("WindowCascade normalised to %d off macOS, want WindowCentered", got)
	}
}

func TestNewWindowNormalisesCascadeOffMacOS(t *testing.T) {
	if macWindowExtrasSupported() {
		t.Skip("WindowCascade is honoured on macOS")
	}
	window := NewWindow(WebviewWindowOptions{InitialPosition: WindowCascade})
	if window.options.InitialPosition != WindowCentered {
		t.Fatalf("InitialPosition = %d, want WindowCentered", window.options.InitialPosition)
	}
}

func TestMacWindowExtrasOptionDefaults(t *testing.T) {
	var options WebviewWindowOptions
	if options.Mac.FrameAutosaveName != "" {
		t.Fatal("FrameAutosaveName should default to empty")
	}
	if options.Mac.TitleBar.WindowButtonsOffset != nil {
		t.Fatal("WindowButtonsOffset should default to nil")
	}
	for _, preset := range []MacTitleBar{MacTitleBarDefault, MacTitleBarHidden, MacTitleBarHiddenInset, MacTitleBarHiddenInsetUnified} {
		if preset.WindowButtonsOffset != nil {
			t.Fatal("titlebar presets must not move the window buttons")
		}
	}
}

func TestPrintOptionsValidation(t *testing.T) {
	if _, err := (PrintOptions{}).validate(); err != nil {
		t.Fatalf("zero options rejected: %v", err)
	}
	if _, err := (PrintOptions{Orientation: PrintOrientation(9)}).validate(); err == nil {
		t.Fatal("unknown orientation accepted")
	}
	if _, err := (PrintOptions{Margins: PrintMargins{Left: -1}}).validate(); err == nil {
		t.Fatal("negative margin accepted")
	}
	if _, err := (PrintOptions{Scale: -0.5}).validate(); err == nil {
		t.Fatal("negative scale accepted")
	}
	valid := PrintOptions{Orientation: PrintOrientationPortrait, Margins: PrintMargins{Top: 10, Left: 10, Bottom: 10, Right: 10}, Scale: 0.8, PaperName: "iso-a4"}
	got, err := valid.validate()
	if err != nil {
		t.Fatalf("valid options rejected: %v", err)
	}
	if got != valid {
		t.Fatalf("validate changed the options: %+v", got)
	}
}

func TestPrintOrientationValues(t *testing.T) {
	if PrintOrientationAutomatic != 0 || PrintOrientationPortrait != 1 || PrintOrientationLandscape != 2 {
		t.Fatal("print orientation values must match the native bridge (0 automatic, 1 portrait, 2 landscape)")
	}
	if validPrintOrientation(PrintOrientation(3)) {
		t.Fatal("unknown orientation accepted")
	}
}

func TestLegacyPrintOptionsMatchHistoricalPrint(t *testing.T) {
	options := legacyPrintOptions()
	if options.Orientation != PrintOrientationLandscape {
		t.Fatalf("Print has always printed landscape, got %d", options.Orientation)
	}
	if options.Margins != (PrintMargins{Top: 30, Left: 30, Bottom: 30, Right: 30}) {
		t.Fatalf("Print has always used 30 point margins, got %+v", options.Margins)
	}
	if options.Silent || options.PrinterName != "" || options.PaperName != "" || options.Scale != 0 {
		t.Fatal("Print shows the panels on the default printer at the default scale")
	}
	if !(PrintMargins{}).isZero() || options.Margins.isZero() {
		t.Fatal("PrintMargins.isZero is wrong")
	}
}

func TestExportTimeoutDefaults(t *testing.T) {
	if exportTimeout(0) != DefaultMacExportTimeout || exportTimeout(-time.Second) != DefaultMacExportTimeout {
		t.Fatal("non-positive timeouts should use the default")
	}
	if exportTimeout(time.Second) != time.Second {
		t.Fatal("explicit timeout ignored")
	}
}

func TestMacAttentionRequestHandleIsNilSafe(t *testing.T) {
	var request *MacAttentionRequest
	request.Cancel()
	if request.Critical() {
		t.Fatal("nil handle reported critical")
	}
	handle := &MacAttentionRequest{critical: true}
	if !handle.Critical() {
		t.Fatal("critical flag lost")
	}
	// An id of zero means no native request; Cancel must be a no-op twice.
	handle.Cancel()
	handle.Cancel()
}

func TestMacWindowExtrasQueueBeforeCreation(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{Name: "pending"}}
	window.SetRepresentedFile("/tmp/notes.md")
	window.SetDocumentEdited(true)
	window.SetSubtitle("Draft")
	if got := window.RepresentedFile(); got != "/tmp/notes.md" {
		t.Fatalf("RepresentedFile before creation = %q", got)
	}
	if !window.IsDocumentEdited() {
		t.Fatal("IsDocumentEdited before creation = false")
	}
	state := takeMacWindowExtrasPending(window)
	if state == nil || state.subtitle == nil || *state.subtitle != "Draft" {
		t.Fatalf("queued state = %+v", state)
	}
	if takeMacWindowExtrasPending(window) != nil {
		t.Fatal("queued state was not consumed")
	}
	if window.RepresentedFile() != "" || window.IsDocumentEdited() {
		t.Fatal("consumed state still reported")
	}
}

func TestMacWindowExtrasNilAndUncreatedWindowsAreSafe(t *testing.T) {
	var window *WebviewWindow
	window.SetRepresentedFile("x")
	window.SetDocumentEdited(true)
	window.SetSubtitle("s")
	window.SetFrameAutosaveName("n")
	window.SetWindowButtonsOffset(1, 2)
	window.ResetWindowButtonsOffset()
	window.CascadeFrom(nil)
	if window.RepresentedFile() != "" || window.IsDocumentEdited() {
		t.Fatal("nil window reported state")
	}
	if request := window.RequestAttention(true); request == nil || !request.Critical() {
		t.Fatal("RequestAttention on a nil window must still return a handle")
	}
	if _, err := window.ExportPDF(PDFExportOptions{}); err == nil {
		t.Fatal("ExportPDF on a nil window succeeded")
	}
	if _, err := window.Snapshot(SnapshotOptions{}); err == nil {
		t.Fatal("Snapshot on a nil window succeeded")
	}
	if err := window.PrintWithOptions(PrintOptions{}); !errors.Is(err, ErrMacWindowNotCreated) {
		t.Fatalf("PrintWithOptions on a nil window = %v", err)
	}

	uncreated := &WebviewWindow{options: WebviewWindowOptions{Name: "uncreated"}}
	uncreated.SetFrameAutosaveName("n")
	uncreated.SetWindowButtonsOffset(1, 2)
	uncreated.CascadeFrom(uncreated)
	if err := uncreated.PrintWithOptions(PrintOptions{Scale: -1}); err == nil {
		t.Fatal("invalid options must be rejected before the platform check")
	}
	if macWindowExtrasSupported() {
		if err := uncreated.PrintWithOptions(PrintOptions{}); !errors.Is(err, ErrMacWindowNotCreated) {
			t.Fatalf("PrintWithOptions before creation = %v", err)
		}
		if _, err := uncreated.ExportPDF(PDFExportOptions{}); !errors.Is(err, ErrMacWindowNotCreated) {
			t.Fatalf("ExportPDF before creation = %v", err)
		}
	}
}

func TestMacWindowExtrasOffMacOS(t *testing.T) {
	if macWindowExtrasSupported() {
		t.Skip("exercises the non-macOS stubs")
	}
	window := &WebviewWindow{options: WebviewWindowOptions{Name: "elsewhere"}}
	if _, err := window.ExportPDF(PDFExportOptions{}); !errors.Is(err, ErrMacOnly) {
		t.Fatalf("ExportPDF off macOS = %v, want ErrMacOnly", err)
	}
	if _, err := window.Snapshot(SnapshotOptions{}); !errors.Is(err, ErrMacOnly) {
		t.Fatalf("Snapshot off macOS = %v, want ErrMacOnly", err)
	}
	// PrintWithOptions falls back to Print, which is a no-op before the
	// native window exists.
	if err := window.PrintWithOptions(PrintOptions{Silent: true}); err != nil {
		t.Fatalf("PrintWithOptions off macOS = %v", err)
	}
	if request := window.RequestAttention(false); request == nil || request.id != 0 {
		t.Fatal("RequestAttention off macOS must return an inert handle")
	}
}
