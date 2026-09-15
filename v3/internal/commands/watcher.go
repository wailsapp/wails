package commands

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/atterpac/refresh/engine"
	"github.com/atterpac/refresh/process"
	"gopkg.in/yaml.v3"
)

func ensureIgnored(list *[]string, pattern string) {
	for _, item := range *list {
		if item == pattern {
			return
		}
	}
	*list = append(*list, pattern)
}

type WatcherOptions struct {
	Config string `description:"The config file including path" default:"."`
}

func Watcher(options *WatcherOptions) error {
	// Parse the config file
	type devConfig struct {
		Config engine.Config `yaml:"dev_mode"`
	}

	var devconfig devConfig

	// Parse the config file
	c, err := os.ReadFile(options.Config)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(c, &devconfig)
	if err != nil {
		return err
	}

	ensureIgnored(&devconfig.Config.Ignore.File, "*_test.go")
	if err := applyFrontendReadiness(&devconfig.Config, os.Getenv("FRONTEND_DEVSERVER_URL")); err != nil {
		return err
	}

	watcherEngine, err := engine.NewEngineFromConfig(devconfig.Config)
	if err != nil {
		return err
	}

	// Refresh owns signal handling and blocks until the development session ends.
	defer watcherEngine.Stop()
	return watcherEngine.Start()
}

// Standard projects get startup ordering without rewriting their Taskfiles.
// Custom commands and explicitly configured readiness remain user-owned.
func applyFrontendReadiness(config *engine.Config, frontendURL string) error {
	if frontendURL == "" {
		return nil
	}
	for i := range config.ExecStruct {
		step := &config.ExecStruct[i]
		if step.Type != process.Background || step.Readiness != nil || len(step.Command) != 0 || strings.Join(strings.Fields(step.Cmd), " ") != "wails3 task common:dev:frontend" {
			continue
		}
		target, err := url.Parse(frontendURL)
		if err != nil || target.Hostname() == "" || (target.Scheme != "http" && target.Scheme != "https") {
			return fmt.Errorf("invalid frontend development server URL %q", frontendURL)
		}
		port := target.Port()
		if port == "" {
			port = "80"
			if target.Scheme == "https" {
				port = "443"
			}
		}
		step.Readiness = &process.Readiness{TCP: net.JoinHostPort(target.Hostname(), port), Timeout: "60s"}
	}
	return nil
}
