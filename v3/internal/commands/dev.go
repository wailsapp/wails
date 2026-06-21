package commands

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/wailsapp/wails/v3/internal/flags"
)

const defaultVitePort = 9245
const wailsVitePort = "WAILS_VITE_PORT"

type DevOptions struct {
	flags.Common

	Config   string `description:"The config file including path" default:"./build/config.yml"`
	VitePort int    `name:"port" description:"Specify the vite dev server port"`
	Secure   bool   `name:"s" description:"Enable HTTPS"`
}

func Dev(options *DevOptions) error {
	host := "localhost"

	// flag takes precedence over environment variable
	var port int
	if options.VitePort != 0 {
		port = options.VitePort
	} else if p, err := strconv.Atoi(os.Getenv(wailsVitePort)); err == nil {
		port = p
	} else {
		port = defaultVitePort
	}

	// check if port is already in use
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	if err = l.Close(); err != nil {
		return err
	}

	// Set environment variable for the dev:frontend task
	os.Setenv(wailsVitePort, strconv.Itoa(port))

	// Set url of frontend dev server
	if options.Secure {
		os.Setenv("FRONTEND_DEVSERVER_URL", fmt.Sprintf("https://%s:%d", host, port))
	} else {
		os.Setenv("FRONTEND_DEVSERVER_URL", fmt.Sprintf("http://%s:%d", host, port))
	}

	return Watcher(&WatcherOptions{
		Config: options.Config,
	})
}
