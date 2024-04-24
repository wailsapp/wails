package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wailsapp/wails/v3/internal/flags"
)

type DevOptions struct {
	flags.Common

	Config   string `description:"The config file including path" default:"./build/devmode.config.toml"`
	ViteHost string `name:"vhost" description:"The vite dev server host" default:"localhost"`
	VitePort int    `name:"vport" description:"The vite dev server port" default:"5173"`
}

func Dev(options *DevOptions) error {

	// Set variables for the dev:frontend task
	os.Setenv("VITE_HOST", options.ViteHost)
	os.Setenv("VITE_PORT", strconv.Itoa(options.VitePort))

	// Set url of frontend dev server
	os.Setenv("FRONTEND_DEVSERVER_URL", fmt.Sprintf("http://%s:%d", options.ViteHost, options.VitePort))

	return Watcher(&WatcherOptions{
		Config: options.Config,
	})
}
