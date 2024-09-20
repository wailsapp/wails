package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "webview_window_bindings_darwin.h"

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

	"github.com/wailsapp/wails/v3/internal/runtime"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type WebviewPanel struct {
	WebviewWindow

	options WebviewPanelOptions
	impl    *macosWebviewPanel
	// keyBindings holds the keybindings for the panel
	keyBindings     map[string]func(*WebviewPanel)
}

// NewPanel creates a new panel with the given options
func NewPanel(options WebviewPanelOptions) *WebviewPanel {
	if options.Width == 0 {
		options.Width = 800
	}
	if options.Height == 0 {
		options.Height = 600
	}
	if options.URL == "" {
		options.URL = "/"
	}

	result := &WebviewPanel{
		WebviewWindow: WebviewWindow{
			id:             getWindowID(),
			options:        options.WebviewWindowOptions,
			eventListeners: make(map[uint][]*WindowEventListener),
			contextMenus:   make(map[string]*Menu),
			eventHooks:     make(map[uint][]*WindowEventListener),
			menuBindings:   make(map[string]*MenuItem),
		},
		options:        options,
	}

	result.setupEventMapping()

	// Listen for window closing events and de
	result.OnWindowEvent(events.Common.WindowClosing, func(event *WindowEvent) {
		shouldClose := true
		if result.options.ShouldClose != nil {
			shouldClose = result.options.ShouldClose(result)
		}
		if shouldClose {
			globalApplication.deleteWindowByID(result.id)
			InvokeSync(result.impl.close)
		}
	})

	// Process keybindings
	if result.options.KeyBindings != nil || result.options.WebviewWindowOptions.KeyBindings != nil {
		result.keyBindings = processKeyBindingOptionsForPanel(result.options.KeyBindings, result.options.WebviewWindowOptions.KeyBindings)
	}

	return result
}

func (p *WebviewPanel) Run() {
	if p.impl != nil {
		return
	}

	p.impl = newPanel(p)
	p.WebviewWindow.impl = &p.impl.macosWebviewWindow

	InvokeSync(p.impl.run)
}

// SetFloating makes the panel float above other application in every workspace.
func (p *WebviewPanel) SetFloating(b bool) Window {
	p.options.Floating = b
	if p.impl != nil {
		InvokeSync(func() {
			p.impl.setFloating(b)
		})
	}
	return p
}

func (p *WebviewPanel) HandleKeyEvent(acceleratorString string) {
	if p.impl == nil && !p.isDestroyed() {
		return
	}
	InvokeSync(func() {
		p.impl.handleKeyEvent(acceleratorString)
	})
}

func (p *WebviewPanel) processKeyBinding(acceleratorString string) bool {
	// Check menu bindings
	if p.menuBindings != nil {
		p.menuBindingsLock.RLock()
		defer p.menuBindingsLock.RUnlock()
		if menuItem := p.menuBindings[acceleratorString]; menuItem != nil {
			menuItem.handleClick()
			return true
		}
	}

	// Check key bindings
	if p.keyBindings != nil {
		p.keyBindingsLock.RLock()
		defer p.keyBindingsLock.RUnlock()
		if callback := p.keyBindings[acceleratorString]; callback != nil {
			// Execute callback
			go callback(p)
			return true
		}
	}

	return globalApplication.processKeyBinding(acceleratorString, &p.WebviewWindow)
}

type macosWebviewPanel struct {
	macosWebviewWindow

	nsPanel unsafe.Pointer
	parent  *WebviewPanel
}

func (p *macosWebviewPanel) setFloating(floating bool) {
	C.panelSetFloating(p.nsPanel, C.bool(floating))
}

func newPanel(parent *WebviewPanel) *macosWebviewPanel {
	result := &macosWebviewPanel{
		macosWebviewWindow: macosWebviewWindow{
			parent: &parent.WebviewWindow,
		},
		parent: parent,
	}
	result.parent.RegisterHook(events.Mac.WebViewDidFinishNavigation, func(event *WindowEvent) {
		result.execJS(runtime.Core())
	})
	return result
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
