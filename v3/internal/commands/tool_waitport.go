package commands

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const portWaitInterval = 100 * time.Millisecond

type ToolWaitPortOptions struct {
	Host    string `name:"h" description:"Host to check" default:"localhost"`
	Port    int    `name:"p" description:"Port to check; defaults to WAILS_VITE_PORT when set"`
	Timeout int    `name:"timeout" description:"Maximum number of seconds to wait for the port to open" default:"60"`
}

func waitForPort(check func() bool, timeout time.Duration) bool {
	if check() {
		return true
	}
	if timeout <= 0 {
		return false
	}

	ticker := time.NewTicker(portWaitInterval)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			if check() {
				return true
			}
		case <-timer.C:
			return false
		}
	}
}

func ToolWaitPort(options *ToolWaitPortOptions) error {
	DisableFooter = true

	if options.Port == 0 {
		port := os.Getenv(wailsVitePort)
		if port == "" {
			return fmt.Errorf("please use the -p flag to specify a port or set %s", wailsVitePort)
		}
		var err error
		options.Port, err = strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid %s value %q: %w", wailsVitePort, port, err)
		}
	}

	timeout := time.Duration(options.Timeout) * time.Second
	if !waitForPort(func() bool { return isPortOpen(options.Host, options.Port) }, timeout) {
		return fmt.Errorf("timed out after %s waiting for port %d to open on %s", timeout, options.Port, options.Host)
	}
	return nil
}
