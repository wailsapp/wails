package bridge

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
)

type Bridge struct {
	upgrader websocket.Upgrader
	server   *http.Server
	myLogger *logger.Logger
}

func NewBridge(myLogger *logger.Logger) *Bridge {
	result := &Bridge{
		myLogger: myLogger,
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}

	result.server = &http.Server{Addr: ":34115"}
	http.HandleFunc("/bridge", result.wsBridgeHandler)
	return result
}

func (b *Bridge) Run(dispatcher *messagedispatcher.Dispatcher, bindingDump string, debug bool) error {

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
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
