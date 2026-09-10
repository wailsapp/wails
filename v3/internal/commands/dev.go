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
	AppArgs         []string `name:"appargs" description:"Application arguments as one shell-quoted string; prefer -- followed by arguments"`
	applicationArgs []string
	Verbose         bool   `name:"verbose" description:"Stream build commands and output during development"`
	Quiet           bool   `name:"quiet" description:"Only show development failures"`
	Host            string `name:"host" description:"Frontend bind address (use this Mac’s LAN IP for a physical iOS device)"`
	Device          string `name:"device" description:"iOS device/simulator identifier or Android adb serial"`
	Emulator        string `name:"emulator" description:"Android Virtual Device to start or reuse"`
	Destination     string `name:"destination" description:"iOS destination: simulator (default) or device"`
	flags.Common

	Config   string `description:"The config file including path" default:"./build/config.yml"`
	VitePort int    `name:"port" description:"Specify the vite dev server port"`
	Secure   bool   `name:"s" description:"Enable HTTPS"`
	Tags     string `name:"tags" description:"Additional development build tags to pass to the Go compiler (comma-separated)"`
	Profile  string `name:"profile" description:"Manifest profile to apply"`
	Target   string `name:"target" description:"Target platform and architecture"`
	Plan     bool   `name:"plan" description:"Print the finite development startup plan without starting a session"`
}

func Dev(options *DevOptions, args ...string) error {
	resolved, err := resolveApplicationArguments(options.AppArgs, args)
	if err != nil {
		return err
	}
	copyOptions := *options
	options = &copyOptions
	options.applicationArgs = resolved
	active, err := activeManifestProject()
	if err != nil {
		return err
	}
	if active {
		return runManifestDev(options)
	}
	if options.applicationArgs != nil {
		return fmt.Errorf("development application arguments require an active wails.hcl and WAILS_EXP_USE_WAKE=1")
	}
	if options.Profile != "" || options.Target != "" || options.Plan || options.Host != "" || options.Device != "" || options.Emulator != "" || options.Destination != "" || options.Verbose || options.Quiet {
		return fmt.Errorf("--profile, --target, --plan, --host, --device, --emulator, --destination, --verbose and --quiet require an active %s", "wails.hcl")
	}
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

	// Environment variables such as WAILS_MCP imply extra build tags. Export them
	// via EXTRA_TAGS so the project Taskfile includes them in dev builds.
	if tags := appendUniqueStrings(splitComma(options.Tags), envTags()...); len(tags) > 0 {
		os.Setenv("EXTRA_TAGS", mergeTags(os.Getenv("EXTRA_TAGS"), tags...))
	}

	return Watcher(&WatcherOptions{
		Config: options.Config,
	})
}
