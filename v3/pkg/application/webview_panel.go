package application

import "github.com/wailsapp/wails/v3/pkg/events"

type webviewPanelImpl interface {
	webviewWindowImpl
	getWebviewWindowImpl() webviewWindowImpl
	setFloating(floating bool)
}

type WebviewPanel struct {
	WebviewWindow

	options WebviewPanelOptions
	impl    webviewPanelImpl
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

	p.impl = newPanelImpl(p)
	p.WebviewWindow.impl = p.impl.getWebviewWindowImpl()

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
