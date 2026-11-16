package dev

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"strings"
)

const devControlsHelp = "Controls (press Enter): r rebuild/restart, b bindings, f frontend restart, c clear, q quit, h help"

type forceRebuildKey struct{}

func WithForcedRebuild(ctx context.Context) context.Context {
	return context.WithValue(ctx, forceRebuildKey{}, true)
}
func IsForcedRebuild(ctx context.Context) bool {
	value, _ := ctx.Value(forceRebuildKey{}).(bool)
	return value
}

type sessionControls struct {
	Output          *Output
	Rebuild         func(force, bindings bool)
	RestartFrontend func() error
}

func (control sessionControls) handle(command string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(command)) {
	case "q", "quit", "exit":
		return true, nil
	case "r", "restart":
		control.Rebuild(true, false)
	case "b", "bindings":
		control.Rebuild(false, true)
	case "f", "frontend":
		return false, control.RestartFrontend()
	case "c", "clear":
		control.Output.Clear()
	case "h", "help", "?":
		control.Output.Status(0, "session", "controls", devControlsHelp)
	}
	return false, nil
}

type frontendEndpoint struct {
	root, host, url string
	port            int
}

func restartSessionFrontend(ctx context.Context, previous Process, config manifest.Config, endpoint frontendEndpoint, ops Operations) (Process, error) {
	if config.Frontend.Disabled {
		return previous, fmt.Errorf("frontend is disabled")
	}
	previous.Stop(time.Duration(config.Dev.GracePeriodMS) * time.Millisecond)
	next, startErr := ops.StartFrontend(endpoint.root, config, endpoint.host, endpoint.port, endpoint.url)
	if startErr == nil {
		startErr = ops.WaitFrontendReady(ctx, next, net.JoinHostPort(endpoint.host, strconv.Itoa(endpoint.port)), 30*time.Second)
	}
	if startErr == nil {
		return next, nil
	}
	if next != nil {
		next.Stop(time.Duration(config.Dev.GracePeriodMS) * time.Millisecond)
	}
	restored, err := ops.RestoreFrontend(ctx, endpoint.root, config, endpoint.port, endpoint.url, endpoint.host)
	if err != nil {
		return nil, fmt.Errorf("frontend restart failed (%v), restore failed: %w", startErr, err)
	}
	return restored, startErr
}

type rebuildTask struct {
	generation        uint64
	force, bindings   bool
	root, frontendURL string
	target            pipeline.Target
	port              int
	args              []string
	operations        Operations
}

func (task rebuildTask) run(ctx context.Context) rebuild {
	result := rebuild{generation: task.generation, force: task.force}
	result.loaded, result.err = task.operations.Load(task.root, "")
	if result.err != nil {
		return result
	}
	result.args, result.err = applicationArguments(result.loaded.Config, task.target.OS, task.target.Arch, task.args)
	if result.err != nil {
		return result
	}
	if task.bindings {
		if task.operations.GenerateBindings == nil {
			result.err = fmt.Errorf("bindings regeneration is unavailable")
			return result
		}
		result.err = task.operations.GenerateBindings(ctx, result.loaded)
		if result.err != nil {
			return result
		}
	}
	result.run, result.err = task.operations.Build(WithGeneration(ctx, task.generation), result.loaded, task.target.OS, task.target.Arch, task.frontendURL, task.port)
	return result
}
