package application

import "C"
import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/assetserver"
	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
	assetserveroptions "github.com/wailsapp/wails/v2/pkg/options/assetserver"

	wailsruntime "github.com/wailsapp/wails/v3/internal/runtime"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/logger"
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
		pid:                       os.Getpid(),
	}
	globalApplication = result

	if !appOptions.Logger.Silent {
		result.log.AddOutput(&logger.Console{})
	}

	result.Events = NewWailsEventProcessor(result.dispatchEventToWindows)

	opts := assetserveroptions.Options{
		Assets:     appOptions.Assets.FS,
		Handler:    appOptions.Assets.Handler,
		Middleware: assetserveroptions.Middleware(appOptions.Assets.Middleware),
	}

	// TODO ServingFrom disk?
	srv, err := assetserver.NewAssetServer("", opts, false, nil, wailsruntime.RuntimeAssetsBundle)
	if err != nil {
		result.fatal(err.Error())
	}

	srv.UseRuntimeHandler(NewMessageProcessor())
	result.assets = srv

	result.bindings, err = NewBindings(appOptions.Bind)
	if err != nil {
		println("Fatal error in application initialisation: ", err.Error())
		os.Exit(1)
	}

	result.plugins = NewPluginManager(appOptions.Plugins, srv)
	err = result.plugins.Init()
	if err != nil {
		result.Quit()
		os.Exit(1)
	}

	err = result.bindings.AddPlugins(appOptions.Plugins)
	if err != nil {
		println("Fatal error in application initialisation: ", err.Error())
		os.Exit(1)
	}

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

var _ webview.Request = &webViewAssetRequest{}

const webViewRequestHeaderWindowId = "x-wails-window-id"
const webViewRequestHeaderWindowName = "x-wails-window-name"

type webViewAssetRequest struct {
	webview.Request
	windowId   uint
	windowName string
}

func (r *webViewAssetRequest) Header() (http.Header, error) {
	h, err := r.Request.Header()
	if err != nil {
		return nil, err
	}

	hh := h.Clone()
	hh.Set(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(r.windowId), 10))
	return hh, nil
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
	plugins  *PluginManager

	// platform app
	impl platformApp

	// The main application menu
	ApplicationMenu *Menu

	clipboard *Clipboard
	Events    *EventProcessor
	log       *logger.Logger

	contextMenus     map[string]*Menu
	contextMenusLock sync.Mutex

	assets *assetserver.AssetServer

	// Hooks
	windowCreatedCallbacks []func(window *WebviewWindow)
	pid                    int
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

func (a *App) deleteWindowByID(id uint) {
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	delete(a.windows, id)
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
	return a.NewWebviewWindowWithOptions(&WebviewWindowOptions{})
}

func (a *App) GetPID() int {
	return a.pid
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

	// Call hooks
	for _, hook := range a.windowCreatedCallbacks {
		hook(newWindow)
	}

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
			request := <-webviewRequests
			a.handleWebViewRequest(request)
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

	err := a.impl.run()
	if err != nil {
		return err
	}

	a.plugins.Shutdown()

	return nil
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

func (a *App) handleWebViewRequest(request *webViewAssetRequest) {
	// Get window from window map
	url, _ := request.URL()
	a.info("Window: '%s', Request: %s", request.windowName, url)
	a.assets.ServeWebViewRequest(request)
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
	if a.impl != nil {
		a.impl.destroy()
	}
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

func (a *App) dispatchEventToWindows(event *WailsEvent) {
	for _, window := range a.windows {
		window.dispatchWailsEvent(event)
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

func (a *App) OnWindowCreation(callback func(window *WebviewWindow)) {
	a.windowCreatedCallbacks = append(a.windowCreatedCallbacks, callback)
}

func (a *App) GetWindowByName(name string) *WebviewWindow {
	a.windowsLock.Lock()
	defer a.windowsLock.Unlock()
	for _, window := range a.windows {
		if window.Name() == name {
			return window
		}
	}
	return nil
}
