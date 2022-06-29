//go:build dev
// +build dev

// Package devserver provides a web-based frontend so that
// it is possible to run a Wails app in a browsers.
package devserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"golang.org/x/net/websocket"
)

type DevWebServer struct {
	server           *echo.Echo
	ctx              context.Context
	appoptions       *options.App
	logger           *logger.Logger
	appBindings      *binding.Bindings
	dispatcher       frontend.Dispatcher
	socketMutex      sync.Mutex
	websocketClients map[*websocket.Conn]*sync.Mutex
	menuManager      *menumanager.Manager
	starttime        string

	// Desktop frontend
	desktopFrontend frontend.Frontend

	devServerAddr string
}

func (d *DevWebServer) WindowSetSystemDefaultTheme() {
	d.desktopFrontend.WindowSetSystemDefaultTheme()
}

func (d *DevWebServer) WindowSetLightTheme() {
	d.desktopFrontend.WindowSetLightTheme()
}

func (d *DevWebServer) WindowSetDarkTheme() {
	d.desktopFrontend.WindowSetDarkTheme()
}

func (d *DevWebServer) Run(ctx context.Context) error {
	d.ctx = ctx

	d.server.GET("/wails/reload", d.handleReload)
	d.server.GET("/wails/ipc", d.handleIPCWebSocket)

	var assetHandler http.Handler
	_fronendDevServerURL, _ := ctx.Value("frontenddevserverurl").(string)
	if _fronendDevServerURL == "" {
		assetdir, _ := ctx.Value("assetdir").(string)
		d.server.GET("/wails/assetdir", func(c echo.Context) error {
			return c.String(http.StatusOK, assetdir)
		})

		var err error
		assetHandler, err = assetserver.NewAssetHandler(ctx, d.appoptions)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		externalURL, err := url.Parse(_fronendDevServerURL)
		if err != nil {
			return err
		}

		if externalURL.Host == "" {
			return fmt.Errorf("Invalid frontend:dev:serverUrl missing protocol scheme?")
		}

		waitCb := func() { d.LogDebug("Waiting for frontend DevServer '%s' to be ready", externalURL) }
		if !checkPortIsOpen(externalURL.Host, time.Minute, waitCb) {
			d.logger.Error("Timeout waiting for frontend DevServer")
		}

		assetHandler = newExternalDevServerAssetHandler(d.logger, externalURL, d.appoptions.AssetsHandler)
	}

	// Setup internal dev server
	bindingsJSON, err := d.appBindings.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	assetServer, err := assetserver.NewBrowserAssetServer(ctx, assetHandler, bindingsJSON)
	if err != nil {
		log.Fatal(err)
	}

	d.server.Any("/*", func(c echo.Context) error {
		assetServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	if devServerAddr := d.devServerAddr; devServerAddr != "" {
		// Start server
		go func(server *echo.Echo, log *logger.Logger) {
			err := server.Start(devServerAddr)
			if err != nil {
				log.Error(err.Error())
			}
			d.LogDebug("Shutdown completed")
		}(d.server, d.logger)

		d.LogDebug("Serving DevServer at http://%s", devServerAddr)

		defer func() {
			err := d.server.Shutdown(context.Background())
			if err != nil {
				d.logger.Error(err.Error())
			}
		}()
	}

	// Launch desktop app
	err = d.desktopFrontend.Run(ctx)
	d.LogDebug("Starting shutdown")

	return err
}

func (d *DevWebServer) Quit() {
	d.desktopFrontend.Quit()
}

func (d *DevWebServer) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return d.desktopFrontend.OpenFileDialog(dialogOptions)
}

func (d *DevWebServer) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	return d.OpenMultipleFilesDialog(dialogOptions)
}

func (d *DevWebServer) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return d.OpenDirectoryDialog(dialogOptions)
}

func (d *DevWebServer) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	return d.desktopFrontend.SaveFileDialog(dialogOptions)
}

func (d *DevWebServer) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	return d.desktopFrontend.MessageDialog(dialogOptions)
}

func (d *DevWebServer) WindowReload() {
	d.broadcast("reload")
	d.desktopFrontend.WindowReload()
}

func (d *DevWebServer) WindowReloadApp() {
	d.broadcast("reloadapp")
	d.desktopFrontend.WindowReloadApp()
}

func (d *DevWebServer) WindowSetTitle(title string) {
	d.desktopFrontend.WindowSetTitle(title)
}

func (d *DevWebServer) WindowShow() {
	d.desktopFrontend.WindowShow()
}

func (d *DevWebServer) WindowHide() {
	d.desktopFrontend.WindowHide()
}

func (d *DevWebServer) WindowCenter() {
	d.desktopFrontend.WindowCenter()
}

func (d *DevWebServer) WindowMaximise() {
	d.desktopFrontend.WindowMaximise()
}

func (d *DevWebServer) WindowToggleMaximise() {
	d.desktopFrontend.WindowToggleMaximise()
}

func (d *DevWebServer) WindowUnmaximise() {
	d.desktopFrontend.WindowUnmaximise()
}

func (d *DevWebServer) WindowMinimise() {
	d.desktopFrontend.WindowMinimise()
}

func (d *DevWebServer) WindowUnminimise() {
	d.desktopFrontend.WindowUnminimise()
}
func (d *DevWebServer) WindowSetAlwaysOnTop(b bool) {
	d.desktopFrontend.WindowSetAlwaysOnTop(b)
}

func (d *DevWebServer) WindowSetPosition(x int, y int) {
	d.desktopFrontend.WindowSetPosition(x, y)
}

func (d *DevWebServer) WindowGetPosition() (int, int) {
	return d.desktopFrontend.WindowGetPosition()
}

func (d *DevWebServer) WindowSetSize(width int, height int) {
	d.desktopFrontend.WindowSetSize(width, height)
}

func (d *DevWebServer) WindowGetSize() (int, int) {
	return d.desktopFrontend.WindowGetSize()
}

func (d *DevWebServer) WindowSetMinSize(width int, height int) {
	d.desktopFrontend.WindowSetMinSize(width, height)
}

func (d *DevWebServer) WindowSetMaxSize(width int, height int) {
	d.desktopFrontend.WindowSetMaxSize(width, height)
}

func (d *DevWebServer) WindowFullscreen() {
	d.desktopFrontend.WindowFullscreen()
}

func (d *DevWebServer) WindowUnfullscreen() {
	d.desktopFrontend.WindowUnfullscreen()
}

func (d *DevWebServer) WindowSetBackgroundColour(col *options.RGBA) {
	d.desktopFrontend.WindowSetBackgroundColour(col)
}

func (d *DevWebServer) MenuSetApplicationMenu(menu *menu.Menu) {
	d.desktopFrontend.MenuSetApplicationMenu(menu)
}

func (d *DevWebServer) MenuUpdateApplicationMenu() {
	d.desktopFrontend.MenuUpdateApplicationMenu()
}

// BrowserOpenURL uses the system default browser to open the url
func (d *DevWebServer) BrowserOpenURL(url string) {
	d.desktopFrontend.BrowserOpenURL(url)
}

func (d *DevWebServer) Notify(name string, data ...interface{}) {
	d.notify(name, data...)
}

func (d *DevWebServer) handleReload(c echo.Context) error {
	d.WindowReload()
	return c.NoContent(http.StatusNoContent)
}
func (d *DevWebServer) handleReloadApp(c echo.Context) error {
	d.WindowReloadApp()
	return c.NoContent(http.StatusNoContent)
}

func (d *DevWebServer) handleIPCWebSocket(c echo.Context) error {
	websocket.Handler(func(c *websocket.Conn) {
		d.LogDebug(fmt.Sprintf("Websocket client %p connected", c))
		d.socketMutex.Lock()
		d.websocketClients[c] = &sync.Mutex{}
		locker := d.websocketClients[c]
		d.socketMutex.Unlock()

		defer func() {
			d.socketMutex.Lock()
			delete(d.websocketClients, c)
			d.socketMutex.Unlock()
			d.LogDebug(fmt.Sprintf("Websocket client %p disconnected", c))
		}()

		var msg string
		defer c.Close()
		for {
			if err := websocket.Message.Receive(c, &msg); err != nil {
				break
			}
			// We do not support drag in browsers
			if msg == "drag" {
				continue
			}

			// Notify the other browsers of "EventEmit"
			if len(msg) > 2 && strings.HasPrefix(string(msg), "EE") {
				d.notifyExcludingSender([]byte(msg), c)
			}

			// Send the message to dispatch to the frontend
			result, err := d.dispatcher.ProcessMessage(string(msg), d)
			if err != nil {
				d.logger.Error(err.Error())
			}
			if result != "" {
				locker.Lock()
				if err = websocket.Message.Send(c, result); err != nil {
					locker.Unlock()
					break
				}
				locker.Unlock()
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (d *DevWebServer) LogDebug(message string, args ...interface{}) {
	d.logger.Debug("[DevWebServer] "+message, args...)
}

type EventNotify struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

func (d *DevWebServer) broadcast(message string) {
	d.socketMutex.Lock()
	defer d.socketMutex.Unlock()
	for client, locker := range d.websocketClients {
		go func(client *websocket.Conn, locker *sync.Mutex) {
			if client == nil {
				d.logger.Error("Lost connection to websocket server")
				return
			}
			locker.Lock()
			err := websocket.Message.Send(client, message)
			if err != nil {
				locker.Unlock()
				d.logger.Error(err.Error())
				return
			}
			locker.Unlock()
		}(client, locker)
	}
}

func (d *DevWebServer) notify(name string, data ...interface{}) {
	// Notify
	notification := EventNotify{
		Name: name,
		Data: data,
	}
	payload, err := json.Marshal(notification)
	if err != nil {
		d.logger.Error(err.Error())
		return
	}
	d.broadcast("n" + string(payload))
}

func (d *DevWebServer) broadcastExcludingSender(message string, sender *websocket.Conn) {
	d.socketMutex.Lock()
	defer d.socketMutex.Unlock()
	for client, locker := range d.websocketClients {
		go func(client *websocket.Conn, locker *sync.Mutex) {
			if client == sender {
				return
			}
			locker.Lock()
			err := websocket.Message.Send(client, message)
			if err != nil {
				locker.Unlock()
				d.logger.Error(err.Error())
				return
			}
			locker.Unlock()
		}(client, locker)
	}
}

func (d *DevWebServer) notifyExcludingSender(eventMessage []byte, sender *websocket.Conn) {
	message := "n" + string(eventMessage[2:])
	d.broadcastExcludingSender(message, sender)

	var notifyMessage EventNotify
	err := json.Unmarshal(eventMessage[2:], &notifyMessage)
	if err != nil {
		d.logger.Error(err.Error())
		return
	}
	d.desktopFrontend.Notify(notifyMessage.Name, notifyMessage.Data...)
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher, menuManager *menumanager.Manager, desktopFrontend frontend.Frontend) *DevWebServer {
	result := &DevWebServer{
		ctx:              ctx,
		desktopFrontend:  desktopFrontend,
		appoptions:       appoptions,
		logger:           myLogger,
		appBindings:      appBindings,
		dispatcher:       dispatcher,
		server:           echo.New(),
		menuManager:      menuManager,
		websocketClients: make(map[*websocket.Conn]*sync.Mutex),
	}

	result.devServerAddr, _ = ctx.Value("devserver").(string)
	result.server.HideBanner = true
	result.server.HidePort = true
	return result
}

func checkPortIsOpen(host string, timeout time.Duration, waitCB func()) (ret bool) {
	if timeout == 0 {
		timeout = time.Minute
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, _ := net.DialTimeout("tcp", host, 2*time.Second)
		if conn != nil {
			conn.Close()
			return true
		}

		waitCB()
		time.Sleep(1 * time.Second)
	}
	return false
}
