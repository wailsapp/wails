package bridge

import (
	"context"
	_ "embed"
	"log"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v2/internal/menumanager"

	"github.com/wailsapp/wails/v2/internal/messagedispatcher"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/internal/logger"
)

//go:embed darwin.js
var darwinRuntime string

// session represents a single websocket session
type session struct {
	bindings string
	conn     *websocket.Conn
	//eventManager interfaces.EventManager
	log *logger.Logger
	//ipc          interfaces.IPCManager

	// Mutex for writing to the socket
	shutdown  chan bool
	writeChan chan []byte

	done bool

	// context
	ctx context.Context

	// client
	client *messagedispatcher.DispatchClient

	// Menus
	menumanager *menumanager.Manager
}

func newSession(conn *websocket.Conn, menumanager *menumanager.Manager, bindings string, dispatcher *messagedispatcher.Dispatcher, logger *logger.Logger, ctx context.Context) *session {
	result := &session{
		conn:        conn,
		bindings:    bindings,
		log:         logger,
		shutdown:    make(chan bool),
		writeChan:   make(chan []byte, 100),
		ctx:         ctx,
		menumanager: menumanager,
	}

	result.client = dispatcher.RegisterClient(newBridgeClient(result))

	return result

}

// Identifier returns a string identifier for the remote connection.
// Taking the form of the client's <ip address>:<port>.
func (s *session) Identifier() string {
	if s.conn != nil {
		return s.conn.RemoteAddr().String()
	}
	return ""
}

func (s *session) sendMessage(msg string) error {
	if !s.done {
		s.writeChan <- []byte(msg)
	}
	return nil
}

func (s *session) start(firstSession bool) {
	s.log.SetLogLevel(1)
	s.log.Info("Connected to frontend.")
	go s.writePump()

	var wailsRuntime string
	switch runtime.GOOS {
	case "darwin":
		wailsRuntime = darwinRuntime
	default:
		log.Fatal("platform not supported")
	}

	bindingsMessage := "window.wailsbindings = `" + s.bindings + "`;"
	s.log.Info(bindingsMessage)
	bootstrapMessage := bindingsMessage + wailsRuntime

	s.sendMessage("b" + bootstrapMessage)

	// Send menus
	traymenus, err := s.menumanager.GetTrayMenus()
	if err != nil {
		s.log.Error(err.Error())
	}

	for _, trayMenu := range traymenus {
		s.sendMessage("TS" + trayMenu)
	}

	for {
		messageType, buffer, err := s.conn.ReadMessage()
		if messageType == -1 {
			return
		}
		if err != nil {
			s.log.Error("Error reading message: %v", err)
			err = s.conn.Close()
			return
		}

		message := string(buffer)

		s.log.Debug("Got message: %#v\n", message)

		// Dispatch message as normal
		s.client.DispatchMessage(message)

		if s.done {
			break
		}
	}
}

// Shutdown
func (s *session) Shutdown() {
	s.conn.Close()
	s.done = true
	s.log.Info("session %v exit", s.Identifier())
}

// writePump pulls messages from the writeChan and sends them to the client
// since it uses a channel to read the messages the socket is protected without locks
func (s *session) writePump() {
	s.log.Debug("Session %v - writePump start", s.Identifier())
	defer s.log.Debug("Session %v - writePump shutdown", s.Identifier())
	for {
		select {
		case <-s.ctx.Done():
			s.Shutdown()
			return
		case msg, ok := <-s.writeChan:
			s.conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
			if !ok {
				s.log.Debug("writeChan was closed!")
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := s.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				s.log.Debug(err.Error())
				return
			}
		}
	}
}
