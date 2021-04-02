package renderer

import (
	"time"

	"github.com/wailsapp/wails/runtime"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
)

// TODO Move this back into bridge.go

// session represents a single websocket session
type session struct {
	bindingCache []string
	conn         *websocket.Conn
	eventManager interfaces.EventManager
	log          *logger.CustomLogger
	ipc          interfaces.IPCManager

	// Mutex for writing to the socket
	shutdown  chan bool
	writeChan chan []byte

	done bool
}

func newSession(conn *websocket.Conn, bindingCache []string, ipc interfaces.IPCManager, logger *logger.CustomLogger, eventMgr interfaces.EventManager) *session {
	return &session{
		conn:         conn,
		bindingCache: bindingCache,
		ipc:          ipc,
		log:          logger,
		eventManager: eventMgr,
		shutdown:     make(chan bool),
		writeChan:    make(chan []byte, 100),
	}
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
	s.log.Infof("Connected to frontend.")
	go s.writePump()

	s.evalJS(runtime.WailsJS, wailsRuntimeMessage)

	// Inject bindings
	for _, binding := range s.bindingCache {
		s.evalJS(binding, bindingMessage)
	}
	s.eventManager.Emit("wails:bridge:session:started", s.Identifier())

	// Emit that everything is loaded and ready
	if firstSession {
		s.eventManager.Emit("wails:ready")
	}

	for {
		messageType, buffer, err := s.conn.ReadMessage()
		if messageType == -1 {
			return
		}
		if err != nil {
			s.log.Errorf("Error reading message: %v", err)
			continue
		}

		s.log.Debugf("Got message: %#v\n", string(buffer))

		s.ipc.Dispatch(string(buffer), s.Callback)

		if s.done {
			break
		}
	}
}

// Callback sends a callback to the frontend
func (s *session) Callback(data string) error {
	return s.evalJS(data, callbackMessage)
}

func (s *session) evalJS(js string, mtype messageType) error {
	// Prepend message type to message
	return s.sendMessage(mtype.toString() + js)
}

// Shutdown
func (s *session) Shutdown() {
	s.done = true
	s.shutdown <- true
	s.log.Debugf("session %v exit", s.Identifier())
}

// writePump pulls messages from the writeChan and sends them to the client
// since it uses a channel to read the messages the socket is protected without locks
func (s *session) writePump() {
	s.log.Debugf("Session %v - writePump start", s.Identifier())
	for {
		select {
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
		case <-s.shutdown:
			break
		}
	}
	s.log.Debug("writePump exiting...")
}
