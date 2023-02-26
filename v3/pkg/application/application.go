package application

import "C"
import (
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/logger"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var globalApplication *App

func init() {
	runtime.LockOSThread()
}

func New(appOptions Options) *App {
	if globalApplication != nil {
		return globalApplication
	}

	mergeApplicationDefaults(&appOptions)

	result := &App{
		options:                   appOptions,
		applicationEventListeners: make(map[uint][]func()),
		systemTrays:               make(map[uint]*SystemTray),
		log:                       logger.New(appOptions.Logger.CustomLoggers...),
		contextMenus:              make(map[string]*Menu),
	}

	if !appOptions.Logger.Silent {
		result.log.AddOutput(&logger.Console{})
	}

	result.Events = NewCustomEventProcessor(result.dispatchEventToWindows)
	globalApplication = result
	return result
}

func mergeApplicationDefaults(o *Options) {
	if o.Name == "" {
		o.Name = "My Wails Application"
	}
	if o.Description == "" {
		o.Description = "An application written using Wails"
	}
	if o.Icon == nil {
		o.Icon = DefaultApplicationIcon
	}

}

type platformApp interface {
	run() error
	destroy()
	setApplicationMenu(menu *Menu)
	name() string
	getCurrentWindowID() uint
	showAboutDialog(name string, description string, icon []byte)
	setIcon(icon []byte)
	on(id uint)
	dispatchOnMainThread(id uint)
	hide()
	show()
}

// Messages sent from javascript get routed here
type windowMessage struct {
	windowId uint
	message  string
}

var windowMessageBuffer = make(chan *windowMessage)

type dragAndDropMessage struct {
	windowId  uint
	filenames []string
}

var windowDragAndDropBuffer = make(chan *dragAndDropMessage)

type webViewAssetRequest struct {
	windowId uint
	request  webview.Request
}

var webviewRequests = make(chan *webViewAssetRequest)

type App struct {
	options                       Options
	applicationEventListeners     map[uint][]func()
	applicationEventListenersLock sync.RWMutex

	// Windows
	windows     map[uint]*WebviewWindow
	windowsLock sync.Mutex

	// System Trays
	systemTrays      map[uint]*SystemTray
	systemTraysLock  sync.Mutex
	systemTrayID     uint
	systemTrayIDLock sync.RWMutex

	// MenuItems
	menuItems     map[uint]*MenuItem
	menuItemsLock sync.Mutex

	// Running
	running  bool
	bindings *Bindings

	// platform app
	impl platformApp

	// The main application menu
	ApplicationMenu *Menu

	clipboard *Clipboard
	Events    *EventProcessor
	log       *logger.Logger

	contextMenus     map[string]*Menu
	contextMenusLock sync.Mutex
}

func (a *App) getSystemTrayID() uint {
	a.systemTrayIDLock.Lock()
	defer a.systemTrayIDLock.Unlock()
	a.systemTrayID++
	return a.systemTrayID
}

func (a *App) getWindowForID(id uint) *WebviewWindow {
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	return a.windows[id]
}

func (a *App) On(eventType events.ApplicationEventType, callback func()) {
	eventID := uint(eventType)
	a.applicationEventListenersLock.Lock()
	defer a.applicationEventListenersLock.Unlock()
	a.applicationEventListeners[eventID] = append(a.applicationEventListeners[eventID], callback)
	if a.impl != nil {
		go a.impl.on(eventID)
	}
}
func (a *App) NewWebviewWindow() *WebviewWindow {
	return a.NewWebviewWindowWithOptions(nil)
}

func (a *App) info(message string, args ...any) {
	a.Log(&logger.Message{
		Level:   "INFO",
		Message: message,
		Data:    args,
		Sender:  "Wails",
	})
}

func (a *App) fatal(message string, args ...any) {
	msg := "************** FATAL **************\n"
	msg += message
	msg += "***********************************\n"

	a.Log(&logger.Message{
		Level:   "FATAL",
		Message: msg,
		Data:    args,
		Sender:  "Wails",
	})

	a.log.Flush()
	os.Exit(1)
}

func (a *App) error(message string, args ...any) {
	a.Log(&logger.Message{
		Level:   "ERROR",
		Message: message,
		Data:    args,
		Sender:  "Wails",
	})
}

func (a *App) NewWebviewWindowWithOptions(windowOptions *WebviewWindowOptions) *WebviewWindow {
	// Ensure we have sane defaults
	if windowOptions == nil {
		windowOptions = WebviewWindowDefaults
	}

	newWindow := NewWindow(windowOptions)
	id := newWindow.id
	if a.windows == nil {
		a.windows = make(map[uint]*WebviewWindow)
	}
	a.windowsLock.Lock()
	a.windows[id] = newWindow
	a.windowsLock.Unlock()

	if a.running {
		newWindow.run()
	}

	return newWindow
}

func (a *App) NewSystemTray() *SystemTray {
	id := a.getSystemTrayID()
	newSystemTray := NewSystemTray(id)
	a.systemTraysLock.Lock()
	a.systemTrays[id] = newSystemTray
	a.systemTraysLock.Unlock()

	if a.running {
		newSystemTray.Run()
	}
	return newSystemTray
}

func (a *App) Run() error {
	a.info("Starting application")
	a.impl = newPlatformApp(a)

	a.running = true
	go func() {
		for {
			event := <-applicationEvents
			a.handleApplicationEvent(event)
		}
	}()
	go func() {
		for {
			event := <-windowEvents
			a.handleWindowEvent(event)
		}
	}()
	go func() {
		for {
			event := <-webviewRequests
			a.handleWebViewRequest(event)
			err := event.request.Release()
			if err != nil {
				a.error("Failed to release webview request: %s", err.Error())
			}
		}
	}()
	go func() {
		for {
			event := <-windowMessageBuffer
			a.handleWindowMessage(event)
		}
	}()
	go func() {
		for {
			dragAndDropMessage := <-windowDragAndDropBuffer
			a.handleDragAndDropMessage(dragAndDropMessage)
		}
	}()

	go func() {
		for {
			menuItemID := <-menuItemClicked
			a.handleMenuItemClicked(menuItemID)
		}
	}()

	var err error
	a.bindings, err = NewBindings(a.options.Bind)
	if err != nil {
		return err
	}

	// run windows
	for _, window := range a.windows {
		go window.run()
	}

	// run system trays
	for _, systray := range a.systemTrays {
		go systray.Run()
	}

	// set the application menu
	a.impl.setApplicationMenu(a.ApplicationMenu)

	// set the application Icon
	a.impl.setIcon(a.options.Icon)

	return a.impl.run()
}

func (a *App) handleApplicationEvent(event uint) {
	a.applicationEventListenersLock.RLock()
	listeners, ok := a.applicationEventListeners[event]
	a.applicationEventListenersLock.RUnlock()
	if !ok {
		return
	}
	for _, listener := range listeners {
		go listener()
	}
}

func (a *App) handleDragAndDropMessage(event *dragAndDropMessage) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.handleDragAndDropMessage(event)
}

func (a *App) handleWindowMessage(event *windowMessage) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.handleMessage(event.message)
}

func (a *App) handleWebViewRequest(event *webViewAssetRequest) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.windowId]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.windowId)
		return
	}
	// Get callback from window
	window.handleWebViewRequest(event.request)
}

func (a *App) handleWindowEvent(event *WindowEvent) {
	// Get window from window map
	a.windowsLock.Lock()
	window, ok := a.windows[event.WindowID]
	a.windowsLock.Unlock()
	if !ok {
		log.Printf("WebviewWindow #%d not found", event.WindowID)
		return
	}
	window.handleWindowEvent(event.EventID)
}

func (a *App) handleMenuItemClicked(menuItemID uint) {
	menuItem := getMenuItemByID(menuItemID)
	if menuItem == nil {
		log.Printf("MenuItem #%d not found", menuItemID)
		return
	}
	menuItem.handleClick()
}

func (a *App) CurrentWindow() *WebviewWindow {
	if a.impl == nil {
		return nil
	}
	id := a.impl.getCurrentWindowID()
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	return a.windows[id]
}

func (a *App) Quit() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		a.windowsLock.Lock()
		for _, window := range a.windows {
			window.Destroy()
		}
		a.windowsLock.Unlock()
		wg.Done()
	}()
	go func() {
		a.systemTraysLock.Lock()
		for _, systray := range a.systemTrays {
			systray.Destroy()
		}
		a.systemTraysLock.Unlock()
		wg.Done()
	}()
	wg.Wait()
	a.impl.destroy()
}

func (a *App) SetMenu(menu *Menu) {
	a.ApplicationMenu = menu
	if a.impl != nil {
		a.impl.setApplicationMenu(menu)
	}
}
func (a *App) ShowAboutDialog() {
	if a.impl != nil {
		a.impl.showAboutDialog(a.options.Name, a.options.Description, a.options.Icon)
	}
}

func (a *App) InfoDialog() *MessageDialog {
	return newMessageDialog(InfoDialog)
}

func (a *App) QuestionDialog() *MessageDialog {
	return newMessageDialog(QuestionDialog)
}

func (a *App) WarningDialog() *MessageDialog {
	return newMessageDialog(WarningDialog)
}

func (a *App) ErrorDialog() *MessageDialog {
	return newMessageDialog(ErrorDialog)
}

func (a *App) OpenDirectoryDialog() *MessageDialog {
	return newMessageDialog(OpenDirectoryDialog)
}

func (a *App) OpenFileDialog() *OpenFileDialog {
	return newOpenFileDialog()
}

func (a *App) SaveFileDialog() *SaveFileDialog {
	return newSaveFileDialog()
}

func (a *App) GetPrimaryScreen() (*Screen, error) {
	return getPrimaryScreen()
}

func (a *App) GetScreens() ([]*Screen, error) {
	return getScreens()
}

func (a *App) Clipboard() *Clipboard {
	if a.clipboard == nil {
		a.clipboard = newClipboard()
	}
	return a.clipboard
}

func (a *App) dispatchOnMainThread(fn func()) {
	mainThreadFunctionStoreLock.Lock()
	id := generateFunctionStoreID()
	mainThreadFunctionStore[id] = fn
	mainThreadFunctionStoreLock.Unlock()
	// Call platform specific dispatch function
	a.impl.dispatchOnMainThread(id)
}

func (a *App) OpenFileDialogWithOptions(options *OpenFileDialogOptions) *OpenFileDialog {
	result := a.OpenFileDialog()
	result.SetOptions(options)
	return result
}

func (a *App) SaveFileDialogWithOptions(s *SaveFileDialogOptions) *SaveFileDialog {
	result := a.SaveFileDialog()
	result.SetOptions(s)
	return result
}

func (a *App) dispatchEventToWindows(event *CustomEvent) {
	for _, window := range a.windows {
		window.dispatchCustomEvent(event)
	}
}

func (a *App) Hide() {
	if a.impl != nil {
		a.impl.hide()
	}
}

func (a *App) Show() {
	if a.impl != nil {
		a.impl.show()
	}
}

func (a *App) Log(message *logger.Message) {
	a.log.Log(message)
}

func (a *App) RegisterContextMenu(name string, menu *Menu) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	a.contextMenus[name] = menu
}

func (a *App) getContextMenu(name string) (*Menu, bool) {
	a.contextMenusLock.Lock()
	defer a.contextMenusLock.Unlock()
	menu, ok := a.contextMenus[name]
	return menu, ok

}
