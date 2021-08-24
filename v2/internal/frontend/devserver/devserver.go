package devserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/assetserver"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DevServer struct {
	server           *fiber.App
	ctx              context.Context
	appoptions       *options.App
	logger           *logger.Logger
	appBindings      *binding.Bindings
	dispatcher       frontend.Dispatcher
	assetServer      *assetserver.AssetServer
	socketMutex      sync.Mutex
	websocketClients map[*websocket.Conn]struct{}
	menuManager      *menumanager.Manager
	starttime        string

	// Desktop frontend
	desktopFrontend frontend.Frontend
}

func (d *DevServer) WindowReload() {
	d.broadcast("reload")
}

func (d *DevServer) Run(ctx context.Context) error {
	d.ctx = ctx

	d.server.Get("/wails/reload", func(fctx *fiber.Ctx) error {
		d.WindowReload()
		d.desktopFrontend.WindowReload()
		return nil
	})

	d.server.Get("/wails/ipc", websocket.New(func(c *websocket.Conn) {
		d.newWebsocketSession(c)
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
			//d.logger.Info("[%p] %s", c, msg)
			if string(msg) == "drag" {
				continue
			}

			if len(msg) > 2 && strings.HasPrefix(string(msg), "EE") {
				d.notifyExcludingSender(msg, c)
				continue
			}

			result, err := d.dispatcher.ProcessMessage(string(msg))
			if err != nil {
				d.logger.Error(err.Error())
			}
			if result != "" {
				if err = c.WriteMessage(mt, []byte(result)); err != nil {
					log.Println("write:", err)
					break
				}
			}

		}
	}))

	_assetdir := ctx.Value("assetdir")
	if _assetdir == nil {
		return fmt.Errorf("no assetdir provided")
	}
	if _assetdir != nil {
		assetdir := _assetdir.(string)
		bindingsJSON, err := d.appBindings.ToJSON()
		if err != nil {
			log.Fatal(err)
		}
		d.assetServer, err = assetserver.NewAssetServer(assetdir, bindingsJSON, d.appoptions)
		if err != nil {
			log.Fatal(err)
		}
		absdir, err := filepath.Abs(assetdir)
		if err != nil {
			return err
		}
		d.LogInfo("Serving assets from: %s", absdir)
	}

	d.server.Get("*", d.loadAsset)

	// Start server
	go func(server *fiber.App, log *logger.Logger) {
		err := server.Listen(":34115")
		if err != nil {
			log.Error(err.Error())
		}
		d.LogInfo("Shutdown completed")
	}(d.server, d.logger)

	d.LogInfo("Serving application at http://localhost:34115")

	// Launch desktop app
	err := d.desktopFrontend.Run(ctx)
	d.LogInfo("Starting shutdown")
	err2 := d.server.Shutdown()
	if err2 != nil {
		d.logger.Error(err.Error())
	}

	return err
}

func (d *DevServer) Quit() {
	d.desktopFrontend.Quit()
}

func (d *DevServer) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return d.desktopFrontend.OpenFileDialog(dialogOptions)
}

func (d *DevServer) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	return d.OpenMultipleFilesDialog(dialogOptions)
}

func (d *DevServer) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return d.OpenDirectoryDialog(dialogOptions)
}

func (d *DevServer) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	return d.desktopFrontend.SaveFileDialog(dialogOptions)
}

func (d *DevServer) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	return d.desktopFrontend.MessageDialog(dialogOptions)
}

func (d *DevServer) WindowSetTitle(title string) {
	d.desktopFrontend.WindowSetTitle(title)
}

func (d *DevServer) WindowShow() {
	d.desktopFrontend.WindowShow()
}

func (d *DevServer) WindowHide() {
	d.desktopFrontend.WindowHide()
}

func (d *DevServer) WindowCenter() {
	d.desktopFrontend.WindowCenter()
}

func (d *DevServer) WindowMaximise() {
	d.desktopFrontend.WindowMaximise()
}

func (d *DevServer) WindowUnmaximise() {
	d.desktopFrontend.WindowUnmaximise()
}

func (d *DevServer) WindowMinimise() {
	d.desktopFrontend.WindowMinimise()
}

func (d *DevServer) WindowUnminimise() {
	d.desktopFrontend.WindowUnminimise()
}

func (d *DevServer) WindowSetPos(x int, y int) {
	d.desktopFrontend.WindowSetPos(x, y)
}

func (d *DevServer) WindowGetPos() (int, int) {
	return d.desktopFrontend.WindowGetPos()
}

func (d *DevServer) WindowSetSize(width int, height int) {
	d.desktopFrontend.WindowSetSize(width, height)
}

func (d *DevServer) WindowGetSize() (int, int) {
	return d.desktopFrontend.WindowGetSize()
}

func (d *DevServer) WindowSetMinSize(width int, height int) {
	d.desktopFrontend.WindowSetMinSize(width, height)
}

func (d *DevServer) WindowSetMaxSize(width int, height int) {
	d.desktopFrontend.WindowSetMaxSize(width, height)
}

func (d *DevServer) WindowFullscreen() {
	d.desktopFrontend.WindowFullscreen()
}

func (d *DevServer) WindowUnFullscreen() {
	d.desktopFrontend.WindowUnFullscreen()
}

func (d *DevServer) WindowSetColour(colour int) {
	d.desktopFrontend.WindowSetColour(colour)
}

func (d *DevServer) SetApplicationMenu(menu *menu.Menu) {
	d.desktopFrontend.SetApplicationMenu(menu)
}

func (d *DevServer) UpdateApplicationMenu() {
	d.desktopFrontend.UpdateApplicationMenu()
}

func (d *DevServer) Notify(name string, data ...interface{}) {
	d.desktopFrontend.Notify(name, data...)
	// Notify Websockets....
	d.notify(name, data...)
}

func (d *DevServer) loadAsset(ctx *fiber.Ctx) error {
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

func (d *DevServer) LogInfo(message string, args ...interface{}) {
	d.logger.Info("[DevServer] "+message, args...)
}

func (d *DevServer) newWebsocketSession(c *websocket.Conn) {
	d.socketMutex.Lock()
	defer d.socketMutex.Unlock()
	c.SetCloseHandler(func(code int, text string) error {
		d.socketMutex.Lock()
		defer d.socketMutex.Unlock()
		delete(d.websocketClients, c)
		d.LogInfo(fmt.Sprintf("Websocket client %p disconnected", c))
		return nil
	})
	d.websocketClients[c] = struct{}{}
	d.LogInfo(fmt.Sprintf("Websocket client %p connected", c))
}

type EventNotify struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

func (d *DevServer) broadcast(message string) {
	d.socketMutex.Lock()
	defer d.socketMutex.Unlock()
	for client := range d.websocketClients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			d.logger.Error(err.Error())
			return
		}
	}
}

func (d *DevServer) notify(name string, data ...interface{}) {
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

func (d *DevServer) broadcastExcludingSender(message string, sender *websocket.Conn) {
	d.socketMutex.Lock()
	defer d.socketMutex.Unlock()
	for client := range d.websocketClients {
		if client == sender {
			continue
		}
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			d.logger.Error(err.Error())
			return
		}
	}
}

func (d *DevServer) notifyExcludingSender(eventMessage []byte, sender *websocket.Conn) {
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

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher, menuManager *menumanager.Manager) *DevServer {
	result := &DevServer{
		ctx:             ctx,
		desktopFrontend: desktop.NewFrontend(ctx, appoptions, myLogger, appBindings, dispatcher),
		appoptions:      appoptions,
		logger:          myLogger,
		appBindings:     appBindings,
		dispatcher:      dispatcher,
		server: fiber.New(fiber.Config{
			ReadTimeout:           time.Second * 5,
			DisableStartupMessage: true,
		}),
		menuManager:      menuManager,
		websocketClients: make(map[*websocket.Conn]struct{}),
	}
	return result
}
