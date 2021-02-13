package bridge

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/wailsapp/wails/v2/internal/menumanager"

	"github.com/wailsapp/wails/v2/internal/messagedispatcher"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type Bridge struct {
	upgrader websocket.Upgrader
	server   *http.Server
	myLogger *logger.Logger

	bindings   string
	dispatcher *messagedispatcher.Dispatcher

	mu       sync.Mutex
	sessions map[string]*session

	ctx    context.Context
	cancel context.CancelFunc

	// Dialog client
	dialog *messagedispatcher.DispatchClient

	// Menus
	menumanager *menumanager.Manager
}

func NewBridge(myLogger *logger.Logger) *Bridge {
	result := &Bridge{
		myLogger: myLogger,
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		sessions: make(map[string]*session),
	}

	myLogger.SetLogLevel(1)

	ctx, cancel := context.WithCancel(context.Background())
	result.ctx = ctx
	result.cancel = cancel
	result.server = &http.Server{Addr: ":34115"}
	http.HandleFunc("/bridge", result.wsBridgeHandler)
	return result
}

func (b *Bridge) Run(dispatcher *messagedispatcher.Dispatcher, menumanager *menumanager.Manager, bindings string, debug bool) error {

	// Ensure we cancel the context when we shutdown
	defer b.cancel()

	b.bindings = bindings
	b.dispatcher = dispatcher
	b.menumanager = menumanager

	// Setup dialog handler
	dialogClient := NewDialogClient(b.myLogger)
	b.dialog = dispatcher.RegisterClient(dialogClient)
	dialogClient.dispatcher = b.dialog

	b.myLogger.Info("Bridge mode started.")

	err := b.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (b *Bridge) wsBridgeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := b.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	b.myLogger.Info("Connection from frontend accepted [%s].", c.RemoteAddr().String())
	b.startSession(c)

}

func (b *Bridge) startSession(conn *websocket.Conn) {

	// Create a new session for this connection
	s := newSession(conn, b.menumanager, b.bindings, b.dispatcher, b.myLogger, b.ctx)

	// Setup the close handler
	conn.SetCloseHandler(func(int, string) error {
		b.myLogger.Info("Connection dropped [%s].", s.Identifier())
		b.dispatcher.RemoveClient(s.client)
		b.mu.Lock()
		delete(b.sessions, s.Identifier())
		b.mu.Unlock()
		return nil
	})

	b.mu.Lock()
	go s.start(len(b.sessions) == 0)
	b.sessions[s.Identifier()] = s
	b.mu.Unlock()
}
