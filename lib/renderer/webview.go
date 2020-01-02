package renderer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/colors"
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/messages"
	wv "github.com/wailsapp/wails/lib/renderer/webview"
)

// WebView defines the main webview application window
// Default values in []
type WebView struct {
	window        wv.WebView // The webview object
	ipc           interfaces.IPCManager
	log           *logger.CustomLogger
	config        interfaces.AppConfig
	eventManager  interfaces.EventManager
	bindingCache  []string
	enableConsole bool
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

	// Create the WebView instance
	w.window = wv.NewWebview(wv.Settings{
		Width:     config.GetWidth(),
		Height:    config.GetHeight(),
		Title:     config.GetTitle(),
		Resizable: config.GetResizable(),
		URL:       config.GetDefaultHTML(),
		Debug:     !config.GetDisableInspector(),
		ExternalInvokeCallback: func(_ wv.WebView, message string) {
			w.ipc.Dispatch(message, w.Callback)
		},
	})

	// SignalManager.OnExit(w.Exit)

	// Set colour
	err := w.SetColour(config.GetColour())
	if err != nil {
		return err
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

// EnableConsole enables the console!
func (w *WebView) EnableConsole() {
	w.enableConsole = true
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
	if w.enableConsole {
		w.log.Debug("Enabling Wails console")
		console := mewn.String("../../runtime/assets/console.js")
		w.evalJS(console)
	}

	// Runtime assets
	wailsRuntime := mewn.String("../../runtime/assets/wails.js")
	w.evalJS(wailsRuntime)

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
				w.log.Debug("Injecting Default Wails CSS")
				defaultCSS := mewn.String("../../runtime/assets/wails.css")

				w.injectCSS(defaultCSS)
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
func (w *WebView) SelectFile() string {
	var result string

	// We need to run this on the main thread, however Dispatch is
	// non-blocking so we launch this in a goroutine and wait for
	// dispatch to finish before returning the result
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w.window.Dispatch(func() {
			result = w.window.Dialog(wv.DialogTypeOpen, 0, "Select File", "")
			wg.Done()
		})
	}()
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
			result = w.window.Dialog(wv.DialogTypeOpen, wv.DialogFlagDirectory, "Select Directory", "")
			wg.Done()
		})
	}()
	wg.Wait()
	return result
}

// SelectSaveFile opens a dialog that allows the user to select a file to save
func (w *WebView) SelectSaveFile() string {
	var result string
	// We need to run this on the main thread, however Dispatch is
	// non-blocking so we launch this in a goroutine and wait for
	// dispatch to finish before returning the result
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w.window.Dispatch(func() {
			result = w.window.Dialog(wv.DialogTypeSave, 0, "Save file", "")
			wg.Done()
		})
	}()
	wg.Wait()
	return result
}

// Callback sends a callback to the frontend
func (w *WebView) Callback(data string) error {
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

	message := fmt.Sprintf("wails._.Notify('%s','%s')", event.Name, data)
	return w.evalJS(message)
}

// Fullscreen makes the main window go fullscreen
func (w *WebView) Fullscreen() {
	if w.config.GetResizable() == false {
		w.log.Warn("Cannot call Fullscreen() - App.Resizable = false")
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
