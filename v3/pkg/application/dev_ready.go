//go:build !production

package application

import (
	"net/url"
	"os"

	"github.com/wailsapp/wails/v3/internal/devruntime"
)

func notifyDevRuntimeReady() { devruntime.Notify() }

// Called by the Android debug launcher before the application main function.
func configureDevRuntime(frontend, ready, token string) {
	if frontend == "" || ready == "" || token == "" {
		return
	}
	// adb reverse exposes only loopback endpoints inside the Android sandbox.
	for _, address := range []string{frontend, ready} {
		u, e := url.Parse(address)
		if e != nil || u.Hostname() != "127.0.0.1" || (u.Scheme != "http" && u.Scheme != "https") {
			return
		}
	}
	os.Setenv("FRONTEND_DEVSERVER_URL", frontend)
	os.Setenv(devruntime.AddressEnv, ready)
	os.Setenv(devruntime.TokenEnv, token)
	configureDevOutput()
}
