//go:build dev
// +build dev

// Package devserver provides a web-based frontend so that
// it is possible to run a Wails app in a browsers.
package devserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type DevWebServer struct {
	server           *fiber.App
	ctx              context.Context
	appoptions       *options.App
	logger           *logger.Logger
	appBindings      *binding.Bindings
	dispatcher       frontend.Dispatcher
	assetServer      *assetserver.BrowserAssetServer
	socketMutex      sync.Mutex
	websocketClients map[*websocket.Conn]*sync.Mutex
	menuManager      *menumanager.Manager
	starttime        string

	// Desktop frontend
	desktopFrontend frontend.Frontend
}

func (d *DevWebServer) WindowReload() {
	d.broadcast("reload")
}

func (d *DevWebServer) Run(ctx context.Context) error {
	d.ctx = ctx

	assetdir, _ := ctx.Value("assetdir").(string)
	d.server.Get("/wails/assetdir", func(fctx *fiber.Ctx) error {
		return fctx.SendString(assetdir)
	})

	d.server.Get("/wails/reload", func(fctx *fiber.Ctx) error {
		d.WindowReload()
		d.desktopFrontend.WindowReload()
		return nil
	})

	d.server.Get("/wails/ipc", websocket.New(func(c *websocket.Conn) {
		d.newWebsocketSession(c)
		locker := d.websocketClients[c]
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break
			}
			// We do not support drag in browsers
			if string(msg) == "drag" {
				continue
			}

			// Notify the other browsers of "EventEmit"
			if len(msg) > 2 && strings.HasPrefix(string(msg), "EE") {
				d.notifyExcludingSender(msg, c)
			}

			// Send the message to dispatch to the frontend
			result, err := d.dispatcher.ProcessMessage(string(msg), d)
			if err != nil {
				d.logger.Error(err.Error())
			}
			if result != "" {
				locker.Lock()
				if err = c.WriteMessage(mt, []byte(result)); err != nil {
					locker.Unlock()
					break
				}
				locker.Unlock()
			}

		}
	}))

	_devServerURL := ctx.Value("devserverurl")
	if _devServerURL == "http://localhost:34115" {
		// Setup internal dev server
		bindingsJSON, err := d.appBindings.ToJSON()
		if err != nil {
			log.Fatal(err)
		}

		d.assetServer, err = assetserver.NewBrowserAssetServer(ctx, d.appoptions.Assets, bindingsJSON)
		if err != nil {
			log.Fatal(err)
		}

		d.server.Get("*", d.loadAsset)

		// Start server
		go func(server *fiber.App, log *logger.Logger) {
			err := server.Listen("localhost:34115")
			if err != nil {
				log.Error(err.Error())
			}
			d.LogDebug("Shutdown completed")
		}(d.server, d.logger)

		d.LogDebug("Serving application at http://localhost:34115")

		defer func() {
			err := d.server.Shutdown()
			if err != nil {
				d.logger.Error(err.Error())
			}
		}()
	}

	// Launch desktop app
	err := d.desktopFrontend.Run(ctx)
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

func (d *DevWebServer) WindowSetRGBA(col *options.RGBA) {
	d.desktopFrontend.WindowSetRGBA(col)
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

func (d *DevWebServer) loadAsset(ctx *fiber.Ctx) error {
	data, mimetype, err := d.assetServer.Load(ctx.Path())
	if err != nil {
		_, ok := err.(*fs.PathError)
		if !ok {
			return err
		}
		err := ctx.SendStatus(404)
		if err != nil {
			return err
		}
		return nil
	}
	err = ctx.SendStatus(200)
	if err != nil {
		return err
	}
	ctx.Set("Content-Type", mimetype)
	err = ctx.Send(data)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevWebServer) LogDebug(message string, args ...interface{}) {
	d.logger.Debug("[DevWebServer] "+message, args...)
}

func (d *DevWebServer) newWebsocketSession(c *websocket.Conn) {
	d.socketMutex.Lock()
	defer d.socketMutex.Unlock()
	c.SetCloseHandler(func(code int, text string) error {
		d.socketMutex.Lock()
		defer d.socketMutex.Unlock()
		delete(d.websocketClients, c)
		d.LogDebug(fmt.Sprintf("Websocket client %p disconnected", c))
		return nil
	})
	d.websocketClients[c] = &sync.Mutex{}
	d.LogDebug(fmt.Sprintf("Websocket client %p connected", c))
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
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
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
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
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
		ctx:             ctx,
		desktopFrontend: desktopFrontend,
		appoptions:      appoptions,
		logger:          myLogger,
		appBindings:     appBindings,
		dispatcher:      dispatcher,
		server: fiber.New(fiber.Config{

			ReadTimeout:           time.Second * 5,
			DisableStartupMessage: true,
		}),
		menuManager:      menuManager,
		websocketClients: make(map[*websocket.Conn]*sync.Mutex),
	}
	return result
}
