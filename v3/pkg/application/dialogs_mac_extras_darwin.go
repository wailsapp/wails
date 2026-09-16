//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13 -framework UniformTypeIdentifiers

#import "dialogs_mac_extras_darwin.h"
*/
import "C"

import (
	"strings"
	"sync"
	"unsafe"
)

// Content types map straight to allowedContentTypes here, so AddContentType
// leaves the extension filters alone.
const dialogContentTypesNative = true

// formatSeparator joins format columns for the save panel accessory. Labels
// may contain any printable character, so the unit separator is used.
const formatSeparator = "\x1f"

type promptResult struct {
	value string
	ok    bool
}

type colorResult struct {
	color   RGBA
	changed bool
}

type fontResult struct {
	font    FontDescriptor
	changed bool
}

var (
	macDialogExtrasLock sync.Mutex
	alertHelpCallbacks  = make(map[int]func())
	promptSessions      = make(map[int]chan promptResult)
	colorSessions       = make(map[int]*colorSession)
	fontSessions        = make(map[int]*fontSession)
	saveFormatDialogs   = make(map[uint]*SaveFileDialogStruct)
	macDialogExtrasNext int
)

type colorSession struct {
	done     chan colorResult
	onChange func(RGBA)
}

type fontSession struct {
	done     chan fontResult
	onChange func(FontDescriptor)
}

func nextMacDialogExtrasID() int {
	macDialogExtrasNext++
	return macDialogExtrasNext
}

// applyExtras installs the suppression checkbox, help button and accessory
// view on the freshly created NSAlert. Runs on the main thread after the
// buttons are added. It returns the help callback registration to remove
// once the alert is dismissed (-1 when there is none).
func (m *macosDialog) applyExtras() int {
	d := m.dialog
	if !d.showsSuppression && d.helpCallback == nil && d.accessoryView == nil {
		return -1
	}
	helpID := -1
	if d.helpCallback != nil {
		macDialogExtrasLock.Lock()
		helpID = nextMacDialogExtrasID()
		alertHelpCallbacks[helpID] = d.helpCallback
		macDialogExtrasLock.Unlock()
	}
	C.wailsAlertApplyExtras(m.nsDialog,
		C.bool(d.showsSuppression),
		toCString(d.suppressionLabel),
		C.bool(d.helpCallback != nil),
		C.int(helpID),
		d.accessoryView)
	return helpID
}

// finishExtras records the suppression state and releases the help
// registration. Runs on the main thread when the alert is dismissed, before
// the pressed button's callback.
func (m *macosDialog) finishExtras(helpID int) {
	if helpID >= 0 {
		macDialogExtrasLock.Lock()
		delete(alertHelpCallbacks, helpID)
		macDialogExtrasLock.Unlock()
	}
	if m.dialog.showsSuppression && m.nsDialog != nil {
		m.dialog.setSuppressed(bool(C.wailsAlertSuppressionState(m.nsDialog)))
	}
}

//export dialogHelpCallback
func dialogHelpCallback(id C.int) {
	macDialogExtrasLock.Lock()
	callback := alertHelpCallbacks[int(id)]
	macDialogExtrasLock.Unlock()
	if callback != nil {
		go func() {
			defer handlePanic()
			callback()
		}()
	}
}

// openPanelContentTypes joins the content type identifiers for
// showOpenFileDialog; nil when there are none.
func (m *macosOpenFileDialog) contentTypesCString() *C.char {
	return toCString(strings.Join(m.dialog.contentTypes, ";"))
}

// savePanelExtras builds the extra configuration handed to
// showSaveFileDialog. Strings are freed on the C side. The dialog is
// registered so format changes can be routed back to it.
func (m *macosSaveFileDialog) savePanelExtras() C.WailsSavePanelExtras {
	d := m.dialog
	var extras C.WailsSavePanelExtras
	extras.contentTypes = toCString(strings.Join(d.contentTypes, ";"))
	extras.nameFieldLabel = toCString(d.nameFieldLabel)
	extras.tags = toCString(strings.Join(d.tags, ";"))
	extras.selectedFormat = C.int(d.selectedFormat)
	if len(d.formats) > 0 {
		labels := make([]string, len(d.formats))
		extensions := make([]string, len(d.formats))
		identifiers := make([]string, len(d.formats))
		for i, format := range d.formats {
			labels[i] = strings.ReplaceAll(format.Label, formatSeparator, "")
			if labels[i] == "" {
				labels[i] = format.Extension
			}
			extensions[i] = strings.TrimPrefix(strings.TrimSpace(format.Extension), ".")
			identifiers[i] = strings.TrimSpace(format.UTI)
		}
		extras.formatLabels = C.CString(strings.Join(labels, formatSeparator))
		extras.formatExtensions = C.CString(strings.Join(extensions, formatSeparator))
		extras.formatIdentifiers = C.CString(strings.Join(identifiers, formatSeparator))
		macDialogExtrasLock.Lock()
		saveFormatDialogs[d.id] = d
		macDialogExtrasLock.Unlock()
	}
	return extras
}

// saveFormatDialogDone drops the format registration once the save panel
// has completed.
func saveFormatDialogDone(id uint) {
	macDialogExtrasLock.Lock()
	delete(saveFormatDialogs, id)
	macDialogExtrasLock.Unlock()
}

//export saveFileDialogFormatCallback
func saveFileDialogFormatCallback(cid C.uint, index C.int) {
	macDialogExtrasLock.Lock()
	dialog := saveFormatDialogs[uint(cid)]
	macDialogExtrasLock.Unlock()
	if dialog != nil {
		dialog.setSelectedFormat(int(index))
	}
}

func dialogPrompt(options PromptOptions) (string, bool, error) {
	done := make(chan promptResult, 1)
	macDialogExtrasLock.Lock()
	id := nextMacDialogExtrasID()
	promptSessions[id] = done
	macDialogExtrasLock.Unlock()

	InvokeAsync(func() {
		var parent unsafe.Pointer
		if options.Window != nil {
			parent = options.Window.NativeWindow()
		}
		C.wailsShowPromptDialog(C.int(id),
			toCString(options.Title),
			toCString(options.Message),
			toCString(options.Placeholder),
			toCString(options.DefaultValue),
			C.bool(options.Secure),
			toCString(options.OKLabel),
			toCString(options.CancelLabel),
			parent)
	})

	result := <-done
	return result.value, result.ok, nil
}

//export dialogPromptCallback
func dialogPromptCallback(id C.int, value *C.char, ok C.bool) {
	macDialogExtrasLock.Lock()
	done := promptSessions[int(id)]
	delete(promptSessions, int(id))
	macDialogExtrasLock.Unlock()
	if done == nil {
		return
	}
	result := promptResult{ok: bool(ok)}
	if value != nil {
		result.value = C.GoString(value)
	}
	done <- result
}

func dialogPickColor(options ColorPickerOptions) (RGBA, bool, error) {
	session := &colorSession{
		done:     make(chan colorResult, 1),
		onChange: options.OnChange,
	}
	macDialogExtrasLock.Lock()
	id := nextMacDialogExtrasID()
	colorSessions[id] = session
	macDialogExtrasLock.Unlock()

	// Ordering the panel front does not block, so a synchronous call is safe
	// here; the wait happens on the channel below.
	started := InvokeSyncWithResult(func() bool {
		return bool(C.wailsShowColorPanel(C.int(id),
			C.int(options.Initial.Red),
			C.int(options.Initial.Green),
			C.int(options.Initial.Blue),
			C.int(options.Initial.Alpha),
			C.bool(options.ShowsAlpha),
			toCString(options.Title),
			C.bool(options.OnChange != nil)))
	})
	if !started {
		macDialogExtrasLock.Lock()
		delete(colorSessions, id)
		macDialogExtrasLock.Unlock()
		return options.Initial, false, ErrDialogInProgress
	}

	result := <-session.done
	if !result.changed {
		return options.Initial, false, nil
	}
	return result.color, true, nil
}

//export dialogColorCallback
func dialogColorCallback(id C.int, red, green, blue, alpha C.int, closed C.bool, changed C.bool) {
	macDialogExtrasLock.Lock()
	session := colorSessions[int(id)]
	if bool(closed) {
		delete(colorSessions, int(id))
	}
	macDialogExtrasLock.Unlock()
	if session == nil {
		return
	}
	color := RGBA{
		Red:   uint8(red),
		Green: uint8(green),
		Blue:  uint8(blue),
		Alpha: uint8(alpha),
	}
	if !bool(closed) {
		if session.onChange != nil {
			go func() {
				defer handlePanic()
				session.onChange(color)
			}()
		}
		return
	}
	session.done <- colorResult{color: color, changed: bool(changed)}
}

func dialogPickFont(options FontPickerOptions) (FontDescriptor, bool, error) {
	session := &fontSession{
		done:     make(chan fontResult, 1),
		onChange: options.OnChange,
	}
	macDialogExtrasLock.Lock()
	id := nextMacDialogExtrasID()
	fontSessions[id] = session
	macDialogExtrasLock.Unlock()

	started := InvokeSyncWithResult(func() bool {
		return bool(C.wailsShowFontPanel(C.int(id),
			toCString(options.Family),
			C.double(options.Size),
			C.bool(options.OnChange != nil)))
	})
	if !started {
		macDialogExtrasLock.Lock()
		delete(fontSessions, id)
		macDialogExtrasLock.Unlock()
		return FontDescriptor{Family: options.Family, Size: options.Size}, false, ErrDialogInProgress
	}

	result := <-session.done
	return result.font, result.changed, nil
}

//export dialogFontCallback
func dialogFontCallback(id C.int, family, face, postScriptName *C.char, size C.double, closed C.bool, changed C.bool) {
	macDialogExtrasLock.Lock()
	session := fontSessions[int(id)]
	if bool(closed) {
		delete(fontSessions, int(id))
	}
	macDialogExtrasLock.Unlock()
	if session == nil {
		return
	}
	font := FontDescriptor{
		Family:         C.GoString(family),
		Face:           C.GoString(face),
		PostScriptName: C.GoString(postScriptName),
		Size:           float64(size),
	}
	if !bool(closed) {
		if session.onChange != nil {
			go func() {
				defer handlePanic()
				session.onChange(font)
			}()
		}
		return
	}
	session.done <- fontResult{font: font, changed: bool(changed)}
}
