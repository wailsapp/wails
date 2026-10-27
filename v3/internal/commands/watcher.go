package commands

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/exec"
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

func ensurePrimaryExitPolicy(executes []process.Execute) {
	for index := range executes {
		execute := &executes[index]
		if execute.Type == process.Primary && execute.ExitPolicy == "" {
			execute.ExitPolicy = process.ExitPolicyShutdown
		}
	}
}

func isInterruptError(err error) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && isInterruptProcessState(exitErr.ProcessState)
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
	ensurePrimaryExitPolicy(devconfig.Config.ExecStruct)
	if err := applyFrontendReadiness(&devconfig.Config, os.Getenv("FRONTEND_DEVSERVER_URL")); err != nil {
		return err
	}

	watcherEngine, err := engine.NewEngineFromConfig(devconfig.Config)
	if err != nil {
		return err
	}
	if err := watcherEngine.Start(); err != nil {
		if isInterruptError(err) {
			slog.Warn("graceful exit requested", "signal", os.Interrupt)
			return nil
		}
		return err
	}
	return nil
}

// frontendTCPTarget converts the dev server URL into a host:port address.
//
// The readiness gate deliberately uses a TCP probe rather than an HTTP one:
// the HTTP probe sends no Accept header and requires exactly 200, which dev
// servers such as Vite answer with 404 (see #6209). "Port accepts connections"
// is all the gate needs to hold the primary process back until the frontend
// is listening.
func frontendTCPTarget(frontendURL string) (string, error) {
	u, err := url.Parse(frontendURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.Fragment != "" {
		return "", fmt.Errorf("frontend URL %q must be an absolute HTTP(S) URL without a fragment", frontendURL)
	}
	port := u.Port()
	if port == "" {
		port = "80"
		if u.Scheme == "https" {
			port = "443"
		}
	}
	return net.JoinHostPort(u.Hostname(), port), nil
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
		tcpTarget, err := frontendTCPTarget(frontendURL)
		if err != nil {
			return fmt.Errorf("invalid frontend readiness: %w", err)
		}
		step.Readiness = &process.Readiness{TCP: tcpTarget, Timeout: "60s"}
		if err := step.Validate(); err != nil {
			return fmt.Errorf("invalid frontend readiness: %w", err)
		}
	}
	return nil
}
