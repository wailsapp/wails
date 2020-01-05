package renderer

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/leaanthony/mewn"
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
	lock sync.Mutex
}

// Identifier returns a string identifier for the remote connection.
// Taking the form of the client's <ip address>:<port>.
func (s *session) Identifier() string {
	return s.conn.RemoteAddr().String()
}

func (s *session) sendMessage(msg string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s *session) start(firstSession bool) {
	s.log.Infof("Connected to frontend.")

	wailsRuntime := mewn.String("../../runtime/assets/wails.js")
	s.evalJS(wailsRuntime, wailsRuntimeMessage)

	// Inject bindings
	for _, binding := range s.bindingCache {
		s.evalJS(binding, bindingMessage)
	}

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
