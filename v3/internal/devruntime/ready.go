// Package devruntime implements the private CLI/application readiness protocol.
package devruntime

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const AddressEnv = "WAILS_INTERNAL_DEV_READY_URL"
const TokenEnv = "WAILS_INTERNAL_DEV_READY_TOKEN"

var once sync.Once

// Notify is called only after the webview has announced its runtime readiness.
// It never blocks the UI thread, and is inert outside a CLI-owned dev launch.
func Notify() {
	address, token := os.Getenv(AddressEnv), os.Getenv(TokenEnv)
	if address == "" || token == "" {
		return
	}
	once.Do(func() {
		go func() {
			client := http.Client{Timeout: 3 * time.Second, Transport: &http.Transport{Proxy: nil}}
			defer client.CloseIdleConnections()
			req, err := http.NewRequest(http.MethodPost, address, strings.NewReader("ready"))
			if err != nil {
				return
			}
			req.Header.Set("Authorization", "Bearer "+token)
			response, err := client.Do(req)
			if err == nil {
				response.Body.Close()
			}
		}()
	})
}
