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
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/assetserver"

	"github.com/wailsapp/wails/v2/internal/frontend/runtime"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/ipcauth"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/menumanager"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type Screen = frontend.Screen

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Refuse a WebSocket upgrade unless its Origin is the dev server's own. The
	// dev runtime builds its socket URL from window.location, so a genuine
	// browser client is always same-origin; there is no legitimate headerless
	// client, so a missing Origin is refused too.
	CheckOrigin: sameOrigin,
}

// sameOrigin reports whether the upgrade carries an Origin header equal to the
// request Host. A missing, malformed, non-HTTP, or foreign Origin is refused.
// Because requireAllowedHost has already validated the Host header against an
// allowlist, comparing the Origin against it is meaningful rather than a
// comparison of two attacker-controlled values.
func sameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != "" && strings.EqualFold(u.Host, r.Host)
}

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
	frontend.Frontend

	devServerAddr string
}

func (d *DevWebServer) Run(ctx context.Context) error {
	d.ctx = ctx

	d.registerRoutes()

	assetServerConfig, err := assetserver.BuildAssetServerConfig(d.appoptions)
	if err != nil {
		return err
	}

	var myLogger assetserver.Logger
	if _logger := ctx.Value("logger"); _logger != nil {
		myLogger = _logger.(*logger.Logger)
	}

	var wsHandler http.Handler

	_fronendDevServerURL, _ := ctx.Value("frontenddevserverurl").(string)
	if _fronendDevServerURL == "" {
		assetdir, _ := ctx.Value("assetdir").(string)
		d.server.GET("/wails/assetdir", func(c echo.Context) error {
			return c.String(http.StatusOK, assetdir)
		})

	} else {
		externalURL, err := url.Parse(_fronendDevServerURL)
		if err != nil {
			return err
		}

		// WebSockets aren't currently supported in prod mode, so a WebSocket connection is the result of the
		// FrontendDevServer e.g. Vite to support auto reloads.
		// Therefore we direct WebSockets directly to the FrontendDevServer instead of returning a NotImplementedStatus.
		wsHandler = httputil.NewSingleHostReverseProxy(externalURL)
	}

	assetHandler, err := assetserver.NewAssetHandler(assetServerConfig, myLogger)
	if err != nil {
		log.Fatal(err)
	}

	// Setup internal dev server
	bindingsJSON, err := d.appBindings.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	assetServer, err := assetserver.NewDevAssetServer(assetHandler, bindingsJSON, ctx.Value("assetdir") != nil, myLogger, runtime.RuntimeAssetsBundle)
	if err != nil {
		log.Fatal(err)
	}

	d.server.Any("/*", func(c echo.Context) error {
		if c.IsWebSocket() {
			wsHandler.ServeHTTP(c.Response(), c.Request())
		} else {
			assetServer.ServeHTTP(c.Response(), c.Request())
		}
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
	}

	// Launch desktop app
	err = d.Frontend.Run(ctx)

	return err
}

// registerRoutes wires the dev server's middleware and fixed routes. The host
// guard is registered first so it runs outermost (echo composes Use middleware
// outermost-first): a request whose Host is not allowed is rejected before the
// capability cookie is set or any route runs.
func (d *DevWebServer) registerRoutes() {
	d.server.Use(d.requireAllowedHost())

	// Hand the per-launch capability to every page the dev server serves, as an
	// HttpOnly, SameSite=Strict cookie — defense in depth behind the host guard
	// and origin check. A same-origin page returns it automatically on the IPC
	// WebSocket upgrade. It is not local-client authentication: a process
	// running as the same user can read it straight off a served response. Its
	// job is to raise the bar for the browser vector, not to authenticate an
	// arbitrary local caller.
	d.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !c.IsWebSocket() {
				http.SetCookie(c.Response(), &http.Cookie{
					Name:     ipcauth.CookieName,
					Value:    ipcauth.Token(),
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})
			}
			return next(c)
		}
	})

	d.server.GET("/wails/reload", d.handleReload)
	d.server.GET("/wails/ipc", d.handleIPCWebSocket)
}

func (d *DevWebServer) WindowReload() {
	d.broadcast("reload")
	d.Frontend.WindowReload()
}

func (d *DevWebServer) WindowReloadApp() {
	d.broadcast("reloadapp")
	d.Frontend.WindowReloadApp()
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
	// Require the capability cookie on the upgrade as defense in depth, behind
	// the host guard and the same-origin check. This is not local-client
	// authentication — a process running as the same user can read the cookie
	// from a served response — so it is not relied on to stop such a caller.
	if cookie, err := c.Cookie(ipcauth.CookieName); err != nil || !ipcauth.Valid(cookie.Value) {
		d.logger.Error("IPC WebSocket rejected: missing or invalid capability")
		return c.NoContent(http.StatusForbidden)
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		d.logger.Error("WebSocket upgrade failed %v", err)
		return err
	}
	d.LogDebug(fmt.Sprintf("WebSocket client %p connected", conn))

	d.socketMutex.Lock()
	d.websocketClients[conn] = &sync.Mutex{}
	locker := d.websocketClients[conn]
	d.socketMutex.Unlock()

	var wg sync.WaitGroup

	defer func() {
		wg.Wait()
		d.socketMutex.Lock()
		delete(d.websocketClients, conn)
		d.socketMutex.Unlock()
		d.LogDebug(fmt.Sprintf("WebSocket client %p disconnected", conn))
		conn.Close()
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		msg := string(msgBytes)
		wg.Add(1)

		go func(m string) {
			defer wg.Done()

			if m == "drag" {
				return
			}

			if len(m) > 2 && strings.HasPrefix(m, "EE") {
				d.notifyExcludingSender([]byte(m), conn)
			}

			result, err := d.dispatcher.ProcessMessage(m, d)
			if err != nil {
				d.logger.Error(err.Error())
			}

			if result != "" {
				locker.Lock()
				defer locker.Unlock()
				if err := conn.WriteMessage(websocket.TextMessage, []byte(result)); err != nil {
					d.logger.Error("Websocket write message failed %v", err)
				}
			}
		}(msg)
	}

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
	d.Frontend.Notify(notifyMessage.Name, notifyMessage.Data...)
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher, menuManager *menumanager.Manager, desktopFrontend frontend.Frontend) *DevWebServer {
	result := &DevWebServer{
		ctx:              ctx,
		Frontend:         desktopFrontend,
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
