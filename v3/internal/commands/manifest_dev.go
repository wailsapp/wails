package commands

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type manifestProcess struct {
	cmd  *exec.Cmd
	done chan error
}

func runManifestDev(options *DevOptions) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	loaded, err := manifest.Load(root, options.Profile)
	if err != nil {
		return err
	}
	goos, goarch, err := splitTarget(options.Target)
	if err != nil {
		return err
	}
	if goos == "" {
		goos, goarch = runtime.GOOS, runtime.GOARCH
	}
	if goos != runtime.GOOS || goarch != runtime.GOARCH {
		return fmt.Errorf("dev target must be the host %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	port := options.VitePort
	if port == 0 {
		port = loaded.Config.Dev.Port
	}
	if port == 0 {
		port = defaultVitePort
	}
	host := "localhost"
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("frontend port %d is unavailable: %w", port, err)
	}
	_ = listener.Close()
	scheme := "http"
	if options.Secure {
		scheme = "https"
	}
	frontendURL := fmt.Sprintf("%s://%s:%d", scheme, host, port)
	previousURL, hadURL := os.LookupEnv("FRONTEND_DEVSERVER_URL")
	_ = os.Setenv("FRONTEND_DEVSERVER_URL", frontendURL)
	defer func() {
		if hadURL {
			_ = os.Setenv("FRONTEND_DEVSERVER_URL", previousURL)
		} else {
			_ = os.Unsetenv("FRONTEND_DEVSERVER_URL")
		}
	}()

	if err := runManifestPipeline(manifestRunOptions{Verb: "build", Profile: options.Profile, TargetOS: goos, TargetArch: goarch, Development: true, Tags: envTags()}); err != nil {
		return err
	}
	frontend, err := startFrontendDev(root, loaded.Config, port)
	if err != nil {
		return err
	}
	defer frontend.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
	app, err := startManifestApp(root, loaded.Config, goos)
	if err != nil {
		return err
	}
	defer func() {
		if app != nil {
			app.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		}
	}()

	events := make(chan notify.EventInfo, 128)
	if err := registerManifestWatches(root, loaded.Config, events); err != nil {
		return fmt.Errorf("watch project: %w", err)
	}
	defer notify.Stop(events)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)
	debounce := time.Duration(loaded.Config.Dev.DebounceMS) * time.Millisecond
	if debounce <= 0 {
		debounce = 250 * time.Millisecond
	}
	var timer *time.Timer
	var timerC <-chan time.Time
	for {
		select {
		case <-interrupt:
			return nil
		case event := <-events:
			if ignoreDevEvent(root, loaded.Config, event.Path()) {
				continue
			}
			if timer == nil {
				timer = time.NewTimer(debounce)
			} else {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(debounce)
			}
			timerC = timer.C
		case <-timerC:
			timerC = nil
			if err := runManifestPipeline(manifestRunOptions{Verb: "build", Profile: options.Profile, TargetOS: goos, TargetArch: goarch, Development: true, Tags: envTags()}); err != nil {
				fmt.Fprintf(os.Stderr, "rebuild failed; keeping the current app running: %v\n", err)
				continue
			}
			nextLoaded, err := manifest.Load(root, options.Profile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "reload config failed; keeping the current app running: %v\n", err)
				continue
			}
			next, err := startManifestApp(root, nextLoaded.Config, goos)
			if err != nil {
				fmt.Fprintf(os.Stderr, "restart failed; keeping the current app running: %v\n", err)
				continue
			}
			old := app
			app = next
			loaded = nextLoaded
			old.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		}
	}
}

func startFrontendDev(root string, config manifest.Config, port int) (*manifestProcess, error) {
	manager := config.Frontend.PackageManager
	if config.Frontend.DevCommand == "" {
		return nil, fmt.Errorf("frontend.dev_command is not set in %s", manifest.Filename)
	}
	var args []string
	switch manager {
	case "npm":
		args = []string{"run", config.Frontend.DevCommand, "--", "--port", strconv.Itoa(port), "--strictPort"}
	case "pnpm", "bun":
		args = []string{"run", config.Frontend.DevCommand, "--port", strconv.Itoa(port), "--strictPort"}
	case "yarn":
		args = []string{config.Frontend.DevCommand, "--port", strconv.Itoa(port), "--strictPort"}
	default:
		return nil, fmt.Errorf("unsupported frontend.package_manager %q", manager)
	}
	return startPackageManagerProcess(filepath.Join(root, config.Frontend.Directory), manager, args...)
}
func startPackageManagerProcess(dir, manager string, args ...string) (*manifestProcess, error) {
	name := manager
	if manager == "npm" {
		npm, err := exec.LookPath("npm")
		if err != nil {
			return nil, err
		}
		script, err := filepath.EvalSymlinks(npm)
		if err != nil {
			script = npm
		}
		node, err := exec.LookPath("node")
		if err != nil {
			return nil, err
		}
		name = node
		args = append([]string{script}, args...)
	}
	return startManifestProcess(dir, name, args...)
}
func startManifestApp(root string, config manifest.Config, goos string) (*manifestProcess, error) {
	name := config.Project.BinaryName
	if goos == "windows" {
		name += ".exe"
	}
	return startManifestProcess(root, filepath.Join(root, config.Build.OutputDirectory, name))
}
func startManifestProcess(dir, name string, args ...string) (*manifestProcess, error) {
	cmd := exec.CommandContext(context.Background(), name, args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	result := &manifestProcess{cmd: cmd, done: make(chan error, 1)}
	go func() { result.done <- cmd.Wait() }()
	return result, nil
}
func (p *manifestProcess) stop(grace time.Duration) {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return
	}
	_ = p.cmd.Process.Signal(os.Interrupt)
	if grace <= 0 {
		grace = 1500 * time.Millisecond
	}
	select {
	case <-p.done:
		return
	case <-time.After(grace):
		_ = p.cmd.Process.Kill()
		<-p.done
	}
}
func ignoreDevEvent(root string, config manifest.Config, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return true
	}
	rel = filepath.ToSlash(rel)
	if rel == "." || devPathExcluded(config, rel) {
		return true
	}
	return strings.HasSuffix(rel, "_test.go") || !matchesDevWatch(config.Dev.Watch, rel)
}

func registerManifestWatches(root string, config manifest.Config, events chan<- notify.EventInfo) error {
	return filepath.WalkDir(root, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel != "." && devPathExcluded(config, rel) {
			return filepath.SkipDir
		}
		return notify.Watch(current, events, notify.All)
	})
}

func devPathExcluded(config manifest.Config, rel string) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	for _, segment := range strings.Split(rel, "/") {
		if segment == ".git" || segment == ".wails" || segment == "node_modules" {
			return true
		}
	}
	paths := []string{config.Build.OutputDirectory, config.Frontend.Directory}
	paths = append(paths, config.Dev.Exclude...)
	for _, excluded := range paths {
		clean := strings.Trim(filepath.ToSlash(excluded), "/")
		if clean != "" && (rel == clean || strings.HasPrefix(rel, clean+"/")) {
			return true
		}
	}
	return false
}

func matchesDevWatch(patterns []string, rel string) bool {
	if len(patterns) == 0 {
		return true
	}
	rel = filepath.ToSlash(rel)
	for _, pattern := range patterns {
		pattern = filepath.ToSlash(pattern)
		if matched, _ := path.Match(pattern, rel); matched {
			return true
		}
		if strings.HasPrefix(pattern, "**/") {
			if matched, _ := path.Match(strings.TrimPrefix(pattern, "**/"), path.Base(rel)); matched {
				return true
			}
		}
	}
	return false
}
