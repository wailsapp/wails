package server

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
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
	config     *Config
	srv        *http.Server
	clients    map[string]client
	clientLock sync.Mutex
}

func NewServer(config *Config) *Server {
	s := &Server{
		config:  config,
		clients: map[string]client{},
	}
	return s
}

func (s *Server) Shutdown() {
	if s.srv == nil {
		return
	}
	if err := s.srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func (s *Server) handleClient(rw http.ResponseWriter, req *http.Request) {
	client := client{
		events:  make(chan message, 5),
		address: req.RemoteAddr,
	}
	clientID := req.URL.Query().Get("clientId")
	fmt.Println("client connected", "address", client.Identifier(), "clientID", clientID)
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

func (s *Server) serveRuntime(rw http.ResponseWriter, req *http.Request) {
	runtimeJS := string(bundledassets.RuntimeJS)
	rw.Header().Set("Content-Type", "application/javascript")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(runtimeJS)))
	io.WriteString(rw, runtimeJS)
}

func (s *Server) removeClient(clientID string) {
	s.clientLock.Lock()
	defer s.clientLock.Unlock()
	delete(s.clients, clientID)
	//s.api.Info("client disconnected", "address", clientID)
}

func (s *Server) sendToClient(clientID string, message message) {
	client, ok := s.clients[clientID]
	if !ok {
		fmt.Println("Failed to locate clientID", clientID)
		return
	}
	if err := client.Send(message); err != nil {
		fmt.Println("Failed to send message to client:", err)
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

func (s *Server) handleCall(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()

	params := application.QueryParams(req.URL.Query())

	argMap := map[string]any{}
	values := req.URL.Query()
	args := values.Get("args")
	if args != "" {
		json.Unmarshal([]byte(args), &argMap)
	}
	callID := argMap["call-id"].(string)
	clientID := req.Header.Get("x-wails-client-id")

	var options application.CallOptions
	err := params.ToStruct(&options)
	if err != nil {
		fmt.Println("Error parsing call options:", err)
		return
	}

	fmt.Println("Handling call", req.URL.Path, "callID", callID, "clientID", clientID)

	method := application.Get().BoundMethod(options)
	if method == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Method not found"))
		return
	}

	go func() {
		ctx, _ := context.WithCancel(context.Background())

		result, err := method.Call(ctx, options.Args)
		if err != nil {
			//s.sendToClient(clientID, message{Type: "cberror", Data: fmt.Sprintf("\"Error calling method: %s\"", err)})
			return
		}
		jsonResult, err := json.Marshal(result)
		if err == nil {
			j, err := json.Marshal(callback{
				ID:     callID,
				Result: string(jsonResult),
			})
			fmt.Println("Result:", string(j))
			if err == nil {
				s.sendToClient(clientID, message{Type: "cb", Data: string(j)})
				return
			} else {
				fmt.Println("Failed to marshal result:", err)
			}
		} else {
			s.sendToClient(clientID, message{Type: "cberror", Data: fmt.Sprintf("\"Error calling method: %s\"", err)})
		}
	}()

	rw.WriteHeader(http.StatusAccepted)
	rw.Write([]byte(""))

	fmt.Println(
		"Runtime Call:",
		"duration", time.Since(start),
	)
}

func (s *Server) handleIndex(rw http.ResponseWriter, req *http.Request) {
	html, err := fs.ReadFile(s.config.Assets, "index.html")
	if err != nil {
		http.Error(rw, "Failed to load index.html", http.StatusInternalServerError)
		return
	}
	html = bytes.Replace(html, []byte("</body>"), []byte("<script src=\"/server/ipc.js\"></script>\n</body>"), 1)

	rw.Header().Set("Content-Type", "text/html")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(html)))
	rw.WriteHeader(http.StatusOK)
	rw.Write(html)
}

func (s *Server) handleHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" || req.URL.Path == "/index.html" {
		s.handleIndex(rw, req)
		return
	}
	application.Get().AssetServerHandler()(rw, req)
}

// DispatchWailsEvent sends a WailsEvent to all connected clients
// This implements the WailsEventListener interface.
func (s *Server) DispatchWailsEvent(event *application.WailsEvent) {
	s.sendToAllClients(
		message{
			Type: "wailsevent",
			Data: event.ToJSON(),
		})
}

func (s *Server) run() {
	if s.srv != nil || s.config.Enabled == false {
		fmt.Println("Server already running or not enabled")
		return
	}
	application.Get().RegisterListener(s)
	go s.serve()
}

// ---------------- Plugin Methods ----------------
func (s *Server) serve() {
	s.srv = &http.Server{Addr: s.config.ListenAddress()}
	http.HandleFunc("/server/ipc.js", s.serveIPC)
	http.HandleFunc("/wails/runtime.js", s.serveRuntime)
	http.HandleFunc("/server/events", s.handleClient)
	http.HandleFunc("/wails/runtime", s.handleCall)
	http.HandleFunc("/", s.handleHTTP)

	//s.api.Info("Server plugin", "host", s.config.Host, "port", s.config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), nil))
}
