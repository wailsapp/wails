//go:build darwin && purego && !ios && !server

// Package application - CGO-free macOS dialogs.
//
// This is the purego counterpart to dialogs_darwin.go. Instead of compiling the
// AppKit calls through cgo it drives NSAlert / NSOpenPanel / NSSavePanel via the
// Objective-C runtime helpers defined in darwin_purego_cocoa.go.
//
// The Go-visible surface (struct names, constructors and methods) is kept
// identical to the cgo backend so the higher-level dialog code links unchanged.
// Because only one of the two files is ever compiled (mutually exclusive build
// tags) the shared package-level symbols are redefined here rather than shared.
package application

import (
	"strings"
	"sync"

	"github.com/ebitengine/purego/objc"
)

// ---------------------------------------------------------------------------
// NSAlert style constants + mapping
// ---------------------------------------------------------------------------

const (
	nsAlertStyleWarning       = 0
	nsAlertStyleInformational = 1
	nsAlertStyleCritical      = 2
)

// NSModalResponse values we care about.
const (
	nsModalResponseOK        = 1 // NSModalResponseOK
	nsAlertFirstButtonReturn = 1000
)

var alertTypeMap = map[DialogType]int{
	WarningDialogType:  nsAlertStyleWarning,
	InfoDialogType:     nsAlertStyleInformational,
	ErrorDialogType:    nsAlertStyleCritical,
	QuestionDialogType: nsAlertStyleInformational,
}

// ---------------------------------------------------------------------------
// Pure-Go dialog callback registry (message dialogs)
//
// Mirrors the cgo backend's callback map. In the cgo version the completion is
// delivered from an Objective-C block via the exported dialogCallback; here the
// modal is run synchronously on the main thread so we can invoke the callback
// inline, but we keep the registry for structural parity.
// ---------------------------------------------------------------------------

type dialogResultCallback func(int)

var (
	callbacks = make(map[int]dialogResultCallback)
	mutex     = &sync.Mutex{}
)

func addDialogCallback(callback dialogResultCallback) int {
	mutex.Lock()
	defer mutex.Unlock()

	// Find the first free integer key.
	var id int
	for {
		if _, exists := callbacks[id]; !exists {
			break
		}
		id++
	}
	callbacks[id] = callback
	return id
}

func removeDialogCallback(id int) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(callbacks, id)
}

func dialogCallback(id int, buttonPressed int) {
	mutex.Lock()
	callback, exists := callbacks[id]
	mutex.Unlock()
	if !exists {
		return
	}
	callback(buttonPressed)
}

// ---------------------------------------------------------------------------
// Small helpers
// ---------------------------------------------------------------------------

// getButtonNumber maps an NSModalResponse from an NSAlert to a zero-based button
// index, matching the cgo helper of the same name (clamped to 3).
func getButtonNumber(response int) int {
	switch response {
	case nsAlertFirstButtonReturn:
		return 0
	case nsAlertFirstButtonReturn + 1:
		return 1
	case nsAlertFirstButtonReturn + 2:
		return 2
	default:
		return 3
	}
}

// nsImageFromBytes builds an autoreleased-ish NSImage from raw image bytes.
func nsImageFromBytes(b []byte) id {
	return class("NSImage").send("alloc").send("initWithData:", nsData(b))
}

// ---------------------------------------------------------------------------
// About box
// ---------------------------------------------------------------------------

func (m *macosApp) showAboutDialog(title string, message string, icon []byte) {
	InvokeAsync(func() {
		alert := class("NSAlert").send("alloc").send("init")
		if title != "" {
			alert.send("setMessageText:", nsString(title))
		}
		if message != "" {
			alert.send("setInformativeText:", nsString(message))
		}
		if len(icon) > 0 {
			alert.send("setIcon:", nsImageFromBytes(icon))
		}
		alert.send("setAlertStyle:", nsAlertStyleInformational)
		alert.send("runModal")
		alert.send("release")
	})
}

// ---------------------------------------------------------------------------
// Message dialog (NSAlert)
// ---------------------------------------------------------------------------

type macosDialog struct {
	dialog *MessageDialog

	nsDialog id
}

func newDialogImpl(d *MessageDialog) *macosDialog {
	return &macosDialog{
		dialog: d,
	}
}

func (m *macosDialog) show() {
	InvokeAsync(func() {
		// Mac can only have 4 buttons on a dialog.
		if len(m.dialog.Buttons) > 4 {
			m.dialog.Buttons = m.dialog.Buttons[:4]
		}

		if !m.nsDialog.isNil() {
			m.nsDialog.send("release")
			m.nsDialog = 0
		}

		alertType, ok := alertTypeMap[m.dialog.DialogType]
		if !ok {
			alertType = nsAlertStyleInformational
		}

		alert := class("NSAlert").send("alloc").send("init")
		alert.send("setAlertStyle:", alertType)
		if m.dialog.Title != "" {
			alert.send("setMessageText:", nsString(m.dialog.Title))
		}
		if m.dialog.Message != "" {
			alert.send("setInformativeText:", nsString(m.dialog.Message))
		}

		// Icon: explicit icon bytes, else the app icon for errors, else the
		// default system caution/info image (matching the cgo createAlert).
		var iconBytes []byte
		if len(m.dialog.Icon) > 0 {
			iconBytes = m.dialog.Icon
		} else if m.dialog.DialogType == ErrorDialogType && len(globalApplication.options.Icon) > 0 {
			iconBytes = globalApplication.options.Icon
		}
		if len(iconBytes) > 0 {
			alert.send("setIcon:", nsImageFromBytes(iconBytes))
		} else {
			var imageName string
			if alertType == nsAlertStyleCritical || alertType == nsAlertStyleWarning {
				imageName = "NSCaution"
			} else {
				imageName = "NSInfo"
			}
			img := class("NSImage").send("imageNamed:", nsString(imageName))
			if !img.isNil() {
				alert.send("setIcon:", img)
			}
		}

		// Reverse the buttons so that the default ends up on the right.
		reversedButtons := make([]*Button, len(m.dialog.Buttons))
		count := 0
		for i := len(m.dialog.Buttons) - 1; i >= 0; i-- {
			button := m.dialog.Buttons[i]
			nsButton := alert.send("addButtonWithTitle:", nsString(button.Label))
			switch {
			case button.IsDefault:
				nsButton.send("setKeyEquivalent:", nsString("\r"))
			case button.IsCancel:
				nsButton.send("setKeyEquivalent:", nsString("\033"))
			default:
				nsButton.send("setKeyEquivalent:", nsString(""))
			}
			reversedButtons[count] = button
			count++
		}

		m.nsDialog = alert

		var callBackID int
		callBackID = addDialogCallback(func(buttonPressed int) {
			if len(m.dialog.Buttons) > buttonPressed {
				button := reversedButtons[buttonPressed]
				if button.Callback != nil {
					button.Callback()
				}
			}
			removeDialogCallback(callBackID)
		})

		// Run the alert. When a parent window is supplied the cgo backend uses a
		// sheet; here we run it application-modal on the main thread which is the
		// simplest block-free equivalent and delivers the same button result.
		response := get[int](alert, "runModal")
		dialogCallback(callBackID, getButtonNumber(response))
	})
}

// ---------------------------------------------------------------------------
// Open panel filter delegate (NSOpenPanelDelegate)
//
// Recreates OpenPanelDelegate.panel:shouldEnableURL: from the cgo delegate. The
// per-delegate allowed extension list is stored in a Go map keyed by the
// delegate instance pointer (dialogs are shown one at a time, but the map keeps
// it correct even if not).
// ---------------------------------------------------------------------------

var (
	openPanelDelegateClass id
	openPanelDelegateOnce  sync.Once

	openPanelExtMu   sync.Mutex
	openPanelExtByID = map[uintptr][]string{}
)

func openPanelShouldEnableURL(self objc.ID, _ objc.SEL, _ objc.ID, url objc.ID) bool {
	u := id(url)
	if u.isNil() {
		return false
	}
	// Always allow directories so the user can navigate into them.
	if get[bool](u, "hasDirectoryPath") {
		return true
	}

	openPanelExtMu.Lock()
	exts := openPanelExtByID[uintptr(self)]
	openPanelExtMu.Unlock()

	if len(exts) == 0 {
		return true
	}

	ext := u.send("pathExtension").send("lowercaseString").string()
	if ext == "" {
		return false
	}
	for _, allowed := range exts {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}
	return false
}

func registerOpenPanelDelegateClass() {
	openPanelDelegateOnce.Do(func() {
		openPanelDelegateClass = registerDelegateClass(
			"WailsOpenPanelDelegatePurego",
			"NSObject",
			nil,
			[]objc.MethodDef{
				{
					Cmd: sel_("panel:shouldEnableURL:"),
					Fn:  openPanelShouldEnableURL,
				},
			},
		)
	})
}

// newOpenPanelDelegate creates a delegate instance whose shouldEnableURL uses
// the supplied extension list.
func newOpenPanelDelegate(extensions []string) id {
	registerOpenPanelDelegateClass()
	delegate := openPanelDelegateClass.send("alloc").send("init")
	openPanelExtMu.Lock()
	openPanelExtByID[delegate.ptr()] = extensions
	openPanelExtMu.Unlock()
	return delegate
}

func releaseOpenPanelDelegate(delegate id) {
	if delegate.isNil() {
		return
	}
	openPanelExtMu.Lock()
	delete(openPanelExtByID, delegate.ptr())
	openPanelExtMu.Unlock()
	delegate.send("release")
}

// ---------------------------------------------------------------------------
// Open file dialog (NSOpenPanel)
// ---------------------------------------------------------------------------

type macosOpenFileDialog struct {
	dialog *OpenFileDialogStruct
}

func newOpenFileDialogImpl(d *OpenFileDialogStruct) *macosOpenFileDialog {
	return &macosOpenFileDialog{
		dialog: d,
	}
}

// openFileDialogCallback pushes a single selected path onto the response channel.
func openFileDialogCallback(cid uint, path string) {
	channel, ok := openFileResponses[cid]
	if ok {
		channel <- path
	} else {
		panic("No channel found for open file dialog")
	}
}

// openFileDialogCallbackEnd closes the response channel and frees the dialog id.
func openFileDialogCallbackEnd(cid uint) {
	channel, ok := openFileResponses[cid]
	if ok {
		close(channel)
		delete(openFileResponses, cid)
		freeDialogID(cid)
	} else {
		panic("No channel found for open file dialog")
	}
}

func (m *macosOpenFileDialog) show() (chan string, error) {
	openFileResponses[m.dialog.id] = make(chan string)

	// Massage filter patterns into macOS format: a single ";"-joined list of
	// bare extensions (e.g. "png;jpg;gif").
	var filterPatterns string
	if len(m.dialog.filters) > 0 {
		var allPatterns []string
		for _, filter := range m.dialog.filters {
			patternComponents := strings.Split(filter.Pattern, ";")
			for i, component := range patternComponents {
				fp := strings.TrimSpace(component)
				fp = strings.TrimPrefix(fp, "*.")
				patternComponents[i] = fp
			}
			allPatterns = append(allPatterns, strings.Join(patternComponents, ";"))
		}
		filterPatterns = strings.Join(allPatterns, ";")
	}

	panel := class("NSOpenPanel").send("openPanel")

	var delegate id
	if filterPatterns != "" {
		extensions := strings.Split(filterPatterns, ";")
		delegate = newOpenPanelDelegate(extensions)
		panel.send("setDelegate:", delegate)

		// UTType-based content type filtering (macOS 11+). Any extension that
		// does not resolve to a UTType is skipped.
		filterTypes := class("NSMutableArray").send("array")
		utTypeClass := class("UTType")
		for _, ext := range extensions {
			if ext == "" {
				continue
			}
			ut := utTypeClass.send("typeWithFilenameExtension:", nsString(ext))
			if !ut.isNil() {
				filterTypes.send("addObject:", ut)
			}
		}
		if get[int](filterTypes, "count") > 0 {
			panel.send("setAllowedContentTypes:", filterTypes)
		}
	}

	if m.dialog.message != "" {
		panel.send("setMessage:", nsString(m.dialog.message))
	}
	if m.dialog.directory != "" {
		panel.send("setDirectoryURL:", fileURL(m.dialog.directory))
	}
	if m.dialog.buttonText != "" {
		panel.send("setPrompt:", nsString(m.dialog.buttonText))
	}

	panel.send("setCanChooseFiles:", m.dialog.canChooseFiles)
	panel.send("setCanChooseDirectories:", m.dialog.canChooseDirectories)
	panel.send("setCanCreateDirectories:", m.dialog.canCreateDirectories)
	panel.send("setShowsHiddenFiles:", m.dialog.showHiddenFiles)
	panel.send("setAllowsMultipleSelection:", m.dialog.allowsMultipleSelection)
	panel.send("setResolvesAliases:", m.dialog.resolvesAliases)
	panel.send("setExtensionHidden:", m.dialog.hideExtension)
	panel.send("setTreatsFilePackagesAsDirectories:", m.dialog.treatsFilePackagesAsDirectories)
	panel.send("setAllowsOtherFileTypes:", m.dialog.allowsOtherFileTypes)

	// Run the panel modally on the main thread and collect the results. We then
	// feed the channel from a goroutine so this (main-thread) call can return
	// the channel to the caller, which reads it on its own goroutine.
	dialogID := m.dialog.id
	var paths []string
	if get[int](panel, "runModal") == nsModalResponseOK {
		urls := panel.send("URLs")
		count := get[int](urls, "count")
		if count > 0 {
			for i := 0; i < count; i++ {
				url := urls.send("objectAtIndex:", uint(i))
				paths = append(paths, url.send("path").string())
			}
		} else {
			url := panel.send("URL")
			paths = append(paths, url.send("path").string())
		}
	}

	if !delegate.isNil() {
		releaseOpenPanelDelegate(delegate)
	}

	go func() {
		for _, p := range paths {
			openFileDialogCallback(dialogID, p)
		}
		openFileDialogCallbackEnd(dialogID)
	}()

	return openFileResponses[dialogID], nil
}

// ---------------------------------------------------------------------------
// Save file dialog (NSSavePanel)
// ---------------------------------------------------------------------------

type macosSaveFileDialog struct {
	dialog *SaveFileDialogStruct
}

func newSaveFileDialogImpl(d *SaveFileDialogStruct) *macosSaveFileDialog {
	return &macosSaveFileDialog{
		dialog: d,
	}
}

// saveFileDialogCallback delivers the chosen path (possibly empty on cancel) and
// tears down the response channel.
func saveFileDialogCallback(cid uint, path string) {
	channel, ok := saveFileResponses[cid]
	if ok {
		channel <- path
		close(channel)
		delete(saveFileResponses, cid)
		freeDialogID(cid)
	} else {
		panic("No channel found for save file dialog")
	}
}

func (m *macosSaveFileDialog) show() (chan string, error) {
	saveFileResponses[m.dialog.id] = make(chan string)

	panel := class("NSSavePanel").send("savePanel")

	if m.dialog.message != "" {
		panel.send("setMessage:", nsString(m.dialog.message))
	}
	if m.dialog.directory != "" {
		panel.send("setDirectoryURL:", fileURL(m.dialog.directory))
	}
	if m.dialog.filename != "" {
		panel.send("setNameFieldStringValue:", nsString(m.dialog.filename))
	}
	if m.dialog.buttonText != "" {
		panel.send("setPrompt:", nsString(m.dialog.buttonText))
	}

	panel.send("setCanCreateDirectories:", m.dialog.canCreateDirectories)
	panel.send("setShowsHiddenFiles:", m.dialog.showHiddenFiles)
	panel.send("setCanSelectHiddenExtension:", m.dialog.canSelectHiddenExtension)
	panel.send("setExtensionHidden:", m.dialog.hideExtension)
	panel.send("setTreatsFilePackagesAsDirectories:", m.dialog.treatsFilePackagesAsDirectories)
	panel.send("setAllowsOtherFileTypes:", m.dialog.allowOtherFileTypes)

	dialogID := m.dialog.id
	var path string
	if get[int](panel, "runModal") == nsModalResponseOK {
		url := panel.send("URL")
		path = url.send("path").string()
	}

	go func() {
		saveFileDialogCallback(dialogID, path)
	}()

	return saveFileResponses[dialogID], nil
}
