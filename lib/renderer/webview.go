package renderer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/runtime"

	"github.com/go-playground/colors"
	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/messages"
	wv "github.com/wailsapp/wails/lib/renderer/webview"
)

// WebView defines the main webview application window
// Default values in []

// UseFirebug indicates whether to inject the firebug console
var UseFirebug = ""

type WebView struct {
	window       wv.WebView // The webview object
	ipc          interfaces.IPCManager
	log          *logger.CustomLogger
	config       interfaces.AppConfig
	eventManager interfaces.EventManager
	bindingCache []string
	maximumSizeSet bool
}

// NewWebView returns a new WebView struct
func NewWebView() *WebView {
	return &WebView{}
}

// Initialise sets up the WebView
func (w *WebView) Initialise(config interfaces.AppConfig, ipc interfaces.IPCManager, eventManager interfaces.EventManager) error {

	// Store reference to eventManager
	w.eventManager = eventManager

	// Set up logger
	w.log = logger.NewCustomLogger("WebView")

	// Set up the dispatcher function
	w.ipc = ipc
	ipc.BindRenderer(w)

	// Save the config
	w.config = config

	width := config.GetWidth()
	height := config.GetHeight()

	// Clamp width and height
	minWidth, minHeight := config.GetMinWidth(), config.GetMinHeight()
	maxWidth, maxHeight := config.GetMaxWidth(), config.GetMaxHeight()
	setMinSize := minWidth != -1 && minHeight != -1
	setMaxSize := maxWidth != -1 && maxHeight != -1

	if setMinSize {
		if width < minWidth {
			width = minWidth
		}
		if height < minHeight {
			height = minHeight
		}
	}

	if setMaxSize {
		if width > maxWidth {
			width = maxWidth
		}
		if height > maxHeight {
			height = maxHeight
		}
	}

	// Create the WebView instance
	w.window = wv.NewWebview(wv.Settings{
		Width:     width,
		Height:    height,
		Title:     config.GetTitle(),
		Resizable: config.GetResizable(),
		URL:       config.GetHTML(),
		Debug:     !config.GetDisableInspector(),
		ExternalInvokeCallback: func(_ wv.WebView, message string) {
			w.ipc.Dispatch(message, w.callback)
		},
	})
		fmt.Println("Control")

	// Set minimum and maximum sizes
	if setMinSize {
		w.SetMinSize(minWidth, minHeight)
	}
	if setMaxSize {
		w.SetMaxSize(maxWidth, maxHeight)
		fmt.Println("Max")
	}

	// Set minimum and maximum sizes
	if setMinSize {
		w.SetMinSize(minWidth, minHeight)
	}
	if setMaxSize {
		w.SetMaxSize(maxWidth, maxHeight)
	}

	// SignalManager.OnExit(w.Exit)
	
	// Set colour
	color := config.GetColour()
	if color != "" {
		err := w.SetColour(color)
		if err != nil {
			return err
		}
	}

	w.log.Info("Initialised")

	return nil
}

// SetColour sets the window colour
func (w *WebView) SetColour(colour string) error {
	color, err := colors.Parse(colour)
	if err != nil {
		return err
	}
	rgba := color.ToRGBA()
	alpha := uint8(255 * rgba.A)
	w.window.Dispatch(func() {
		w.window.SetColor(rgba.R, rgba.G, rgba.B, alpha)
	})

	return nil
}

// evalJS evaluates the given js in the WebView
// I should rename this to evilJS lol
func (w *WebView) evalJS(js string) error {
	outputJS := fmt.Sprintf("%.45s", js)
	if len(js) > 45 {
		outputJS += "..."
	}
	w.log.DebugFields("Eval", logger.Fields{"js": outputJS})
	//
	w.window.Dispatch(func() {
		w.window.Eval(js)
	})
	return nil
}

// Escape the Javascripts!
func escapeJS(js string) (string, error) {
	result := strings.Replace(js, "\\", "\\\\", -1)
	result = strings.Replace(result, "'", "\\'", -1)
	result = strings.Replace(result, "\n", "\\n", -1)
	return result, nil
}

// evalJSSync evaluates the given js in the WebView synchronously
// Do not call this from the main thread or you'll nuke your app because
// you won't get the callback.
func (w *WebView) evalJSSync(js string) error {

	minified, err := escapeJS(js)

	if err != nil {
		return err
	}

	outputJS := fmt.Sprintf("%.45s", js)
	if len(js) > 45 {
		outputJS += "..."
	}
	w.log.DebugFields("EvalSync", logger.Fields{"js": outputJS})

	ID := fmt.Sprintf("syncjs:%d:%d", time.Now().Unix(), rand.Intn(9999))
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		exit := false
		// We are done when we receive the Callback ID
		w.log.Debug("SyncJS: sending with ID = " + ID)
		w.eventManager.On(ID, func(...interface{}) {
			w.log.Debug("SyncJS: Got callback ID = " + ID)
			wg.Done()
			exit = true
		})
		command := fmt.Sprintf("wails._.AddScript('%s', '%s')", minified, ID)
		w.window.Dispatch(func() {
			w.window.Eval(command)
		})
		for exit == false {
			time.Sleep(time.Millisecond * 1)
		}
	}()

	wg.Wait()

	return nil
}

// injectCSS adds the given CSS to the WebView
func (w *WebView) injectCSS(css string) {
	w.window.Dispatch(func() {
		w.window.InjectCSS(css)
	})
}

// Exit closes the window
func (w *WebView) Exit() {
	w.window.Exit()
}

// Run the window main loop
func (w *WebView) Run() error {

	w.log.Info("Running...")

	// Inject firebug in debug mode on Windows
	if UseFirebug != "" {
		w.log.Debug("Injecting Firebug")
		w.evalJS(`window.usefirebug=true;`)
	}

	// Runtime assets
	w.log.DebugFields("Injecting wails JS runtime", logger.Fields{"js": runtime.WailsJS})
	w.evalJS(runtime.WailsJS)

	// Ping the wait channel when the wails runtime is loaded
	w.eventManager.On("wails:loaded", func(...interface{}) {

		// Run this in a different go routine to free up the main process
		go func() {

			// Inject Bindings
			for _, binding := range w.bindingCache {
				w.evalJSSync(binding)
			}

			// Inject user CSS
			if w.config.GetCSS() != "" {
				outputCSS := fmt.Sprintf("%.45s", w.config.GetCSS())
				if len(outputCSS) > 45 {
					outputCSS += "..."
				}
				w.log.DebugFields("Inject User CSS", logger.Fields{"css": outputCSS})
				w.injectCSS(w.config.GetCSS())
			} else {
				// Use default wails css

				w.log.Debug("Injecting Default Wails CSS: " + runtime.WailsCSS)
				w.injectCSS(runtime.WailsCSS)
			}

			// Inject user JS
			if w.config.GetJS() != "" {
				outputJS := fmt.Sprintf("%.45s", w.config.GetJS())
				if len(outputJS) > 45 {
					outputJS += "..."
				}
				w.log.DebugFields("Inject User JS", logger.Fields{"js": outputJS})
				w.evalJSSync(w.config.GetJS())
			}

			// Emit that everything is loaded and ready
			w.eventManager.Emit("wails:ready")
		}()
	})

	// Kick off main window loop
	w.window.Run()

	return nil
}

// NewBinding registers a new binding with the frontend
func (w *WebView) NewBinding(methodName string) error {
	objectCode := fmt.Sprintf("window.wails._.NewBinding('%s');", methodName)
	w.bindingCache = append(w.bindingCache, objectCode)
	return nil
}

// SelectFile opens a dialog that allows the user to select a file
func (w *WebView) SelectFile(title string, filter string) string {
	var result string

	// We need to run this on the main thread, however Dispatch is
	// non-blocking so we launch this in a goroutine and wait for
	// dispatch to finish before returning the result
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w.window.Dispatch(func() {
			result = w.window.Dialog(wv.DialogTypeOpen, 0, title, "", filter)
			wg.Done()
		})
	}()

	defer w.focus() // Ensure the main window is put back into focus afterwards

	wg.Wait()
	return result
}

// SelectDirectory opens a dialog that allows the user to select a directory
func (w *WebView) SelectDirectory() string {
	var result string
	// We need to run this on the main thread, however Dispatch is
	// non-blocking so we launch this in a goroutine and wait for
	// dispatch to finish before returning the result
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w.window.Dispatch(func() {
			result = w.window.Dialog(wv.DialogTypeOpen, wv.DialogFlagDirectory, "Select Directory", "", "")
			wg.Done()
		})
	}()

	defer w.focus() // Ensure the main window is put back into focus afterwards

	wg.Wait()
	return result
}

// SelectSaveFile opens a dialog that allows the user to select a file to save
func (w *WebView) SelectSaveFile(title string, filter string) string {
	var result string
	// We need to run this on the main thread, however Dispatch is
	// non-blocking so we launch this in a goroutine and wait for
	// dispatch to finish before returning the result
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w.window.Dispatch(func() {
			result = w.window.Dialog(wv.DialogTypeSave, 0, title, "", filter)
			wg.Done()
		})
	}()

	defer w.focus() // Ensure the main window is put back into focus afterwards

	wg.Wait()
	return result
}

// focus puts the main window into focus
func (w *WebView) focus() {
	w.window.Dispatch(func() {
		w.window.Focus()
	})
}

// callback sends a callback to the frontend
func (w *WebView) callback(data string) error {
	callbackCMD := fmt.Sprintf("window.wails._.Callback('%s');", data)
	return w.evalJS(callbackCMD)
}

// NotifyEvent notifies the frontend about a backend runtime event
func (w *WebView) NotifyEvent(event *messages.EventData) error {

	// Look out! Nils about!
	var err error
	if event == nil {
		err = fmt.Errorf("Sent nil event to renderer.WebView")
		w.log.Error(err.Error())
		return err
	}

	// Default data is a blank array
	data := []byte("[]")

	// Process event data
	if event.Data != nil {
		// Marshall the data
		data, err = json.Marshal(event.Data)
		if err != nil {
			w.log.Errorf("Cannot unmarshall JSON data in event: %s ", err.Error())
			return err
		}
	}

	// Double encode data to ensure everything is escaped correctly.
	data, err = json.Marshal(string(data))
	if err != nil {
		w.log.Errorf("Cannot marshal JSON data in event: %s ", err.Error())
		return err
	}

	message := "window.wails._.Notify('" + event.Name + "'," + string(data) + ")"
	return w.evalJS(message)
}

// SetMinSize sets the minimum size of a resizable window
func (w *WebView) SetMinSize(width, height int) {
	if w.config.GetResizable() == false {
		w.log.Warn("Cannot call SetMinSize() - App.Resizable = false")
		return
	}
	w.window.Dispatch(func() {
		w.window.SetMinSize(width, height)
	})
}

// SetMaxSize sets the maximum size of a resizable window
func (w *WebView) SetMaxSize(width, height int) {
	if w.config.GetResizable() == false {
		w.log.Warn("Cannot call SetMaxSize() - App.Resizable = false")
		return
	}
	w.maximumSizeSet = true
	w.window.Dispatch(func() {
		w.window.SetMaxSize(width, height)
	})
}

// Fullscreen makes the main window go fullscreen
func (w *WebView) Fullscreen() {
	if w.config.GetResizable() == false {
		w.log.Warn("Cannot call Fullscreen() - App.Resizable = false")
		return
	} else if w.maximumSizeSet {
		w.log.Warn("Cannot call Fullscreen() - Maximum size of window set")
		return
	}
	w.window.Dispatch(func() {
		w.window.SetFullscreen(true)
	})
}

// UnFullscreen returns the window to the position prior to a fullscreen call
func (w *WebView) UnFullscreen() {
	if w.config.GetResizable() == false {
		w.log.Warn("Cannot call UnFullscreen() - App.Resizable = false")
		return
	}
	w.window.Dispatch(func() {
		w.window.SetFullscreen(false)
	})
}

// SetTitle sets the window title
func (w *WebView) SetTitle(title string) {
	w.window.Dispatch(func() {
		w.window.SetTitle(title)
	})
}

// Close closes the window
func (w *WebView) Close() {
	w.window.Dispatch(func() {
		w.window.Terminate()
	})
}
