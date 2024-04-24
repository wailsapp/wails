package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wailsapp/wails/v3/internal/flags"
)

const defaultVitePort = 9245
const wailsVitePort = "WAILS_VITE_PORT"

type DevOptions struct {
	flags.Common

	Config   string `description:"The config file including path" default:"./build/devmode.config.toml"`
	VitePort int    `name:"port" description:"Specify the vite dev server port"`
}

func Dev(options *DevOptions) error {

	// flag takes precedence over environment variable
	var port int
	if options.VitePort != 0 {
		port = options.VitePort
	} else if p, err := strconv.Atoi(os.Getenv(wailsVitePort)); err == nil {
		port = p
	} else {
		port = defaultVitePort
	}

	// Set environment variable for the dev:frontend task
	os.Setenv(wailsVitePort, strconv.Itoa(port))

	// Set url of frontend dev server
	os.Setenv("FRONTEND_DEVSERVER_URL", fmt.Sprintf("http://%s:%d", "localhost", port))

	return Watcher(&WatcherOptions{
		Config: options.Config,
	})
}
