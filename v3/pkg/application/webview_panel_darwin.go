package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "webview_window_bindings_darwin.h"

void *panelNew(unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences) {
	return createWindow(WindowTypePanel, id, width, height, fraudulentWebsiteWarningEnabled, frameless, enableDragAndDrop, preferences);
}

// Set NSPanel floating
void panelSetFloating(void* nsPanel, bool floating) {
	// Set panel floating on main thread
	NSWindow *window = ((WebviewWindow*)nsPanel).w;
	NSPanel *panel = (NSPanel *) window;

	[panel setLevel:floating ? NSFloatingWindowLevel : NSNormalWindowLevel];
	[panel setFloatingPanel:floating ? YES : NO];
	[panel setStyleMask:floating ? panel.styleMask | NSWindowStyleMaskNonactivatingPanel : panel.styleMask & ~NSWindowStyleMaskNonactivatingPanel];
	NSWindowCollectionBehavior panelCB = NSWindowCollectionBehaviorCanJoinAllSpaces | NSWindowCollectionBehaviorFullScreenAuxiliary;
	[panel setCollectionBehavior:floating ? panel.collectionBehavior | panelCB : panel.collectionBehavior & ~panelCB];
}
*/
import "C"
import (
	"unsafe"
)

type macosWebviewPanel struct {
	*macosWebviewWindow

	nsPanel unsafe.Pointer
	parent  *WebviewPanel
}

func newPanelImpl(parent *WebviewPanel) *macosWebviewPanel {
	result := &macosWebviewPanel{
		macosWebviewWindow: newWindowImpl(parent.WebviewWindow),
		parent:             parent,
	}
	return result
}

func (p *macosWebviewPanel) getWebviewWindowImpl() webviewWindowImpl {
	return p.macosWebviewWindow
}

func (p *macosWebviewPanel) run() {
	for eventId := range p.parent.eventListeners {
		p.on(eventId)
	}
	globalApplication.dispatchOnMainThread(func() {
		options := p.parent.options
		macOptions := options.Mac

		p.nsPanel = C.panelNew(C.uint(p.parent.id),
			C.int(options.Width),
			C.int(options.Height),
			C.bool(macOptions.EnableFraudulentWebsiteWarnings),
			C.bool(options.Frameless),
			C.bool(options.EnableDragAndDrop),
			p.getWebviewPreferences(),
		)
		p.macosWebviewWindow.nsWindow = p.nsPanel

		p.setup(&options.WebviewWindowOptions, &macOptions)
		p.setFloating(options.Floating)
	})
}

func (p *macosWebviewPanel) handleKeyEvent(acceleratorString string) {
	// Parse acceleratorString
	accelerator, err := parseAccelerator(acceleratorString)
	if err != nil {
		globalApplication.error("unable to parse accelerator: %s", err.Error())
		return
	}
	p.parent.processKeyBinding(accelerator.String())
}

func (p *macosWebviewPanel) setFloating(floating bool) {
	C.panelSetFloating(p.nsPanel, C.bool(floating))
}
