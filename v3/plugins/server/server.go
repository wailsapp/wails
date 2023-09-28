package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type message struct {
	Type string
	Data string
}

type callback struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

type client struct {
	address string
	events  chan message
}

func (c client) close() {
	if _, ok := (<-c.events); ok {
		close(c.events)
	}
}

func (c client) Identifier() string {
	return c.address
}

func (c *client) Send(msg message) error {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("connection lost")
		}
	}()
	c.events <- msg
	return err
}

type Server struct {
	id         uint // to allow for registration as a window
	app        *application.App
	config     *Config
	srv        *http.Server
	window     Window
	clients    map[string]client
	clientLock sync.Mutex
}

func NewServer(config *Config) *Server {
	s := &Server{
		config:  config,
		clients: map[string]client{},
	}
	s.window.server = s
	return s
}

func (s *Server) Info(msg string) {
	// s.app.Log(&logger.Message{
	// 	Level:   "INFO",
	// 	Message: fmt.Sprintf("[plugin/server]: %v", msg),
	// })
}

func (s *Server) Shutdown() {
	if err := s.srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func (s *Server) handleClient(rw http.ResponseWriter, req *http.Request) {
	client := client{
		events:  make(chan message, 5),
		address: req.RemoteAddr,
	}
	s.Info(fmt.Sprintf("client %v connected", client.Identifier()))
	clientID := req.URL.Query().Get("clientId")
	if clientID != "" {
		// we only save if we have an identifier
		s.clientLock.Lock()
		s.clients[clientID] = client
		s.clientLock.Unlock()
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	for header, value := range s.config.Headers {
		rw.Header().Set(header, value)
	}
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Connection does not support streaming", http.StatusBadRequest)
		return
	}

	for {
		timeout := time.After(500 * time.Millisecond)
		select {
		case <-req.Context().Done():
			client.close()
			s.removeClient(client.Identifier())
			return
		case msg := <-client.events:
			fmt.Fprintf(rw, "event: %s\n", msg.Type)
			fmt.Fprintf(rw, "data: %v\n\n", msg.Data)
		case <-timeout:
			continue
		}
		flusher.Flush()
	}

}

func (s *Server) serveIPC(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/javascript")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(clientJS)))
	io.WriteString(rw, clientJS)
}

func (s *Server) removeClient(clientID string) {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	delete(s.clients, clientID)
	s.Info(fmt.Sprintf("client %v disconnected", clientID))
}

func (s *Server) sendToClient(requestID string, message message) {
	client, ok := s.clients[requestID]
	if !ok {
		return
	}
	if err := client.Send(message); err != nil {
		s.removeClient(client.Identifier())
	}
}

func (s *Server) sendToAllClients(msg message) {
	if len(s.clients) == 0 {
		return
	}
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	dead := []client{}
	for _, client := range s.clients {
		if err := client.Send(msg); err != nil {
			dead = append(dead, client)
		}
	}
	for _, d := range dead {
		s.removeClient(d.Identifier())
	}
}

func updateCallID(windowID uint, req *http.Request) *http.Request {
	argMap := map[string]any{}
	values := req.URL.Query()
	args := values.Get("args")
	if args != "" {
		json.Unmarshal([]byte(args), &argMap)
	}
	callID := argMap["call-id"]
	clientID := req.Header.Get("x-wails-client-id")
	if clientID != "" {
		argMap["call-id"] = fmt.Sprintf("%s|%s", clientID, callID)
	}
	newArgs, _ := json.Marshal(argMap)
	values.Set("args", string(newArgs))
	req.Header.Add("x-wails-window-id", fmt.Sprintf("%d", windowID))
	req.URL.RawQuery = values.Encode()
	return req
}

func (s *Server) handleHTTP(rw http.ResponseWriter, req *http.Request) {
	req = updateCallID(s.window.id, req)
	s.app.AssetServerHandler()(
		rw,
		req,
	)
}

func (s *Server) run() {
	if s.srv != nil || s.config.Enabled == false {
		return
	}
	address := s.config.ListenAddress()
	s.srv = &http.Server{Addr: address}
	http.HandleFunc("/wails/ipc.js", s.serveIPC)
	http.HandleFunc("/server/events", s.handleClient)
	http.HandleFunc("/", s.handleHTTP)

	s.window.id = s.app.RegisterWindow(s.window)
	go s.serve()
}

// ---------------- Plugin Methods ----------------
func (s *Server) serve() {
	s.Info(fmt.Sprintf("listening %s", s.config.ListenAddress()))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), nil))
}
