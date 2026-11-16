//go:build linux && !gtk3 && !android && !server

package application

import (
	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// dispatchApplicationEvent queues an application event from the GTK callback.
// Both backends deliver it through this one function.
func dispatchApplicationEvent(eventID uint) {
	event := newApplicationEvent(events.ApplicationEventType(eventID))
	applicationEvents <- event
}

func (w *linuxWebviewWindow) isVisible() bool {
	// The GTK widget is created lazily in run() (windowNew). On GTK4 that only
	// happens after the application's "activate" signal fires: WebviewWindow.Run
	// sets w.impl and then blocks in waitForActivation *before* creating the
	// widget, so there is a startup window in which w.impl != nil but w.window is
	// still NULL. A window whose widget does not exist yet is, by definition, not
	// visible; without this guard a visibility poll during that gap calls
	// gtk_widget_is_visible(NULL), which trips a GTK-CRITICAL assertion and
	// returns false anyway.
	if w.window == nilPointer {
		return false
	}
	return widgetIsVisible(w.window)
}

func (a *linuxApp) setIcon(icon []byte) {
	// GTK4 removed per-window icon APIs. The application icon is determined by
	// the .desktop file's Icon= field at the desktop-integration level.
	// No programmatic equivalent exists for setting icons from bytes in GTK4.
}

func (w *linuxWebviewWindow) contextMenuSignals(menu pointer) {
	// GTK4: GtkPopoverMenu items are wired through the "app" GAction group,
	// which is attached to the window in windowNew. The popover's "closed"
	// signal (used for cleanup) is connected in show_context_menu, so there
	// is nothing to wire up here.
}

func (w *linuxWebviewWindow) ignoreMouse(ignore bool) {
	// GTK4: Input handling is different
}

func (w *linuxWebviewWindow) setIcon(icon pointer) {
	// GTK4 removed gtk_window_set_icon. Window icons are set via the
	// application's .desktop file at the desktop-integration level.
}

func (w *linuxWebviewWindow) setTransparent() {
	// GTK4: Transparency via CSS - different from GTK3
}

func (w *linuxWebviewWindow) zoomIn() {
	w.setZoom(w.getZoom() * 1.10)
}

func (w *linuxWebviewWindow) zoomOut() {
	w.setZoom(w.getZoom() / 1.10)
}

func attachMenuHandler(item *MenuItem) uint {
	gtkSignalToMenuItemLock.Lock()
	gtkSignalToMenuItem[item.id] = item
	gtkSignalToMenuItemLock.Unlock()
	return item.id
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func dialogTypeToIconName(dialogType DialogType) string {
	switch dialogType {
	case InfoDialogType:
		return "dialog-information-symbolic"
	case WarningDialogType:
		return "dialog-warning-symbolic"
	case ErrorDialogType:
		return "dialog-error-symbolic"
	case QuestionDialogType:
		return "dialog-question-symbolic"
	default:
		return ""
	}
}

func menuGetRadioGroup(item *linuxMenuItem) *GSList {
	return nil
}

func menuItemRemoveBitmap(widget pointer) {
}

func menuItemSetBitmap(widget pointer, bitmap []byte) {
}

func menuItemSetToolTip(widget pointer, tooltip string) {
}

func menuItemSignalBlock(widget pointer, handlerId uint, block bool) {
}

func nextDialogRequestID() uint {
	dialogRequestMutex.Lock()
	defer dialogRequestMutex.Unlock()
	dialogRequestCounter++
	return uint(dialogRequestCounter)
}

// menuItemActivatedByID reports a click, toggle or radio selection of the menu
// item registered under id. Both backends call it from their GTK action handler.
func menuItemActivatedByID(id uint) {
	gtkSignalToMenuItemLock.RLock()
	item, ok := gtkSignalToMenuItem[id]
	gtkSignalToMenuItemLock.RUnlock()
	if !ok {
		return
	}
	switch item.itemType {
	case text:
		menuItemClicked <- item.id
	case checkbox:
		impl := item.impl.(*linuxMenuItem)
		currentState := impl.isChecked()
		impl.setChecked(!currentState)
		menuItemClicked <- item.id
	case radio:
		// The radio action handler has already made this item the group's
		// selection, so the state cannot tell a new selection from a
		// re-selection. Report every activation, as the macOS and Windows
		// backends do.
		menuItemClicked <- item.id
	}
}

func runOpenFileDialog(dialog *OpenFileDialogStruct) (chan string, error) {
	var action int

	if dialog.canChooseDirectories {
		action = 2 // GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER
	} else {
		action = 0 // GTK_FILE_CHOOSER_ACTION_OPEN
	}

	window := nilPointer
	if dialog.window != nil {
		nativeWindow := dialog.window.NativeWindow()
		if nativeWindow != nil {
			window = nativeWindowPointer(nativeWindow)
		}
	}

	buttonText := dialog.buttonText
	if buttonText == "" {
		buttonText = "_Open"
	}

	return runChooserDialog(
		window,
		dialog.allowsMultipleSelection,
		false, // createFolders not applicable for open
		dialog.showHiddenFiles,
		dialog.directory,
		"",
		dialog.title,
		action,
		buttonText,
		dialog.filters,
	)
}
func runSaveFileDialog(dialog *SaveFileDialogStruct) (chan string, error) {
	window := nilPointer
	if dialog.window != nil {
		nativeWindow := dialog.window.NativeWindow()
		if nativeWindow != nil {
			window = nativeWindowPointer(nativeWindow)
		}
	}

	buttonText := dialog.buttonText
	if buttonText == "" {
		buttonText = "_Save"
	}

	return runChooserDialog(
		window,
		false,
		dialog.canCreateDirectories,
		dialog.showHiddenFiles,
		dialog.directory,
		dialog.filename,
		dialog.title,
		1, // GTK_FILE_CHOOSER_ACTION_SAVE
		buttonText,
		dialog.filters,
	)
}

func menuAppendItemToSection(section pointer, item *MenuItem) {
	if item.impl == nil {
		return
	}
	menuImpl := item.impl.(*linuxMenuItem)
	if menuImpl.native == nilPointer {
		return
	}

	menuImpl.parentMenu = section
	menuImpl.isHidden = item.hidden

	if !item.hidden {
		gMenuAppendItem(section, menuImpl.native)
	}
}

func (a *linuxApp) hideAllWindows() {
	for _, window := range a.getWindows() {
		widgetSetVisible(window, true)
	}
}

func (a *linuxApp) showAllWindows() {
	for _, window := range a.getWindows() {
		windowPresent(window)
	}
}

func appDestroy(application pointer) {
	webview.CloseActiveRequests()
	gApplicationQuit(application)
}

func menuSetSubmenu(item *MenuItem, menu *Menu) {
	if item.impl == nil || menu.impl == nil {
		return
	}
	itemImpl := item.impl.(*linuxMenuItem)
	menuImpl := menu.impl.(*linuxMenu)
	if itemImpl.native == nilPointer || menuImpl.native == nilPointer {
		return
	}
	gMenuItemSetSubmenu(itemImpl.native, menuImpl.native)
}

func windowDetachMenuBar(window, vbox, bar pointer, menuStyle LinuxMenuStyle) {
	if menuStyle == LinuxMenuStylePrimaryMenu {
		windowSetTitlebar(window, nilPointer)
		return
	}
	gtkWidgetUnparent(bar)
}

func windowSetGeometryHints(window pointer, minWidth, minHeight, maxWidth, maxHeight int) {
	if minWidth > 0 && minHeight > 0 {
		widgetSetSizeRequest(window, minWidth, minHeight)
	}
	if maxWidth > 0 || maxHeight > 0 {
		setWindowMaxSize(window, maxWidth, maxHeight)
	}
}

// getCurrentWindowID returns the id of the active window, or 1 when there is
// none. windowMap is read under windowMapLock: registerWindow and
// unregisterWindow write it from other goroutines (the cgo backend used to
// read it unlocked).
func (a *linuxApp) getCurrentWindowID() uint {
	window, ok := activeWindow(a.application)
	if !ok {
		return uint(1)
	}
	a.windowMapLock.Lock()
	identifier, ok := a.windowMap[window]
	a.windowMapLock.Unlock()
	if ok {
		return identifier
	}
	return uint(1)
}

func windowMenuNative(menu *Menu) pointer {
	menu.Update()
	if impl, ok := menu.impl.(*linuxMenu); ok && impl != nil {
		return impl.native
	}
	return nilPointer
}
