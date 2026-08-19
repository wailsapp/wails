//go:build ignore

// benchmark-manifest-dev measures the user-visible Dev session lifecycle with
// real filesystem notifications and child processes. It is intentionally a
// standalone acceptance tool rather than a unit benchmark.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type devSample struct {
	StartupMS        float64 `json:"startup_ms"`
	BackendRebuildMS float64 `json:"backend_rebuild_ms"`
	NoopEventMS      float64 `json:"noop_event_ms"`
	FrontendQuietMS  float64 `json:"frontend_quiet_ms"`
}

type devReport struct {
	Version                int         `json:"version"`
	GOOS                   string      `json:"goos"`
	GOARCH                 string      `json:"goarch"`
	Command                string      `json:"command"`
	Warmups                int         `json:"warmups"`
	Samples                []devSample `json:"samples"`
	MedianStartupMS        float64     `json:"median_startup_ms"`
	MedianBackendRebuildMS float64     `json:"median_backend_rebuild_ms"`
	MedianNoopEventMS      float64     `json:"median_noop_event_ms"`
	LegacyBackendRebuildMS float64     `json:"legacy_backend_rebuild_ms"`
	BackendImprovementPct  float64     `json:"backend_improvement_percent"`
	WatcherQuietWindowMS   float64     `json:"watcher_quiet_window_ms"`
}

type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (b *lockedBuffer) Write(data []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(data)
}

func (b *lockedBuffer) WriteString(value string) {
	_, _ = b.Write([]byte(value))
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}

func main() {
	var executable, output string
	var warmups, samples int
	var legacyRebuildMS, quietWindowMS float64
	flag.StringVar(&executable, "wails", "", "path to the wails3 executable")
	flag.StringVar(&output, "output", "", "optional JSON output path")
	flag.IntVar(&warmups, "warmups", 2, "unreported warmup sessions")
	flag.IntVar(&samples, "samples", 7, "measured sessions")
	flag.Float64Var(&legacyRebuildMS, "legacy-rebuild-ms", 870, "legacy live method-body rebuild baseline")
	flag.Float64Var(&quietWindowMS, "frontend-quiet-ms", 300, "window in which frontend edits must cause no Wails rebuild")
	flag.Parse()
	if runtime.GOOS == "windows" {
		fatalf("manifest Dev benchmarking is not supported on Windows")
	}
	if executable == "" || samples < 1 || warmups < 0 || quietWindowMS <= 0 {
		fatalf("-wails, -samples >= 1, -warmups >= 0, and -frontend-quiet-ms > 0 are required")
	}
	absExecutable, err := filepath.Abs(executable)
	check(err)

	root, err := os.MkdirTemp("", "wails-manifest-dev-benchmark-")
	check(err)
	defer os.RemoveAll(root)
	check(writeFixture(root))
	port := reservePort()
	backendLog := filepath.Join(root, "backend.log")
	frontendLog := filepath.Join(root, "frontend.log")
	environment := append(os.Environ(),
		"PATH="+filepath.Join(root, "tools")+string(os.PathListSeparator)+os.Getenv("PATH"),
		"WAILS_BENCH_BACKEND_LOG="+backendLog,
		"WAILS_BENCH_FRONTEND_LOG="+frontendLog,
	)

	measured := make([]devSample, 0, samples)
	for index := 0; index < warmups+samples; index++ {
		sample, runErr := measureSession(absExecutable, root, environment, port, backendLog, frontendLog, index, time.Duration(quietWindowMS*float64(time.Millisecond)))
		check(runErr)
		if index >= warmups {
			measured = append(measured, sample)
		}
	}

	report := devReport{
		Version:                1,
		GOOS:                   runtime.GOOS,
		GOARCH:                 runtime.GOARCH,
		Command:                absExecutable + " dev --port " + strconv.Itoa(port),
		Warmups:                warmups,
		Samples:                measured,
		MedianStartupMS:        median(sampleValues(measured, func(sample devSample) float64 { return sample.StartupMS })),
		MedianBackendRebuildMS: median(sampleValues(measured, func(sample devSample) float64 { return sample.BackendRebuildMS })),
		MedianNoopEventMS:      median(sampleValues(measured, func(sample devSample) float64 { return sample.NoopEventMS })),
		LegacyBackendRebuildMS: legacyRebuildMS,
		WatcherQuietWindowMS:   quietWindowMS,
	}
	if legacyRebuildMS > 0 {
		report.BackendImprovementPct = (legacyRebuildMS - report.MedianBackendRebuildMS) / legacyRebuildMS * 100
		if report.MedianBackendRebuildMS > legacyRebuildMS {
			fatalf("median backend rebuild %.3fms regressed beyond legacy %.3fms", report.MedianBackendRebuildMS, legacyRebuildMS)
		}
	}
	encoded, err := json.MarshalIndent(report, "", "  ")
	check(err)
	encoded = append(encoded, '\n')
	if output != "" {
		check(os.WriteFile(output, encoded, 0o644))
	}
	_, _ = os.Stdout.Write(encoded)
}

func measureSession(executable, root string, environment []string, port int, backendLog, frontendLog string, generation int, quietWindow time.Duration) (devSample, error) {
	beforeBackends := lineCount(backendLog)
	beforeFrontends := lineCount(frontendLog)
	command := exec.Command(executable, "dev", "--port", strconv.Itoa(port))
	command.Dir = root
	command.Env = environment
	var output lockedBuffer
	stdout, err := command.StdoutPipe()
	if err != nil {
		return devSample{}, err
	}
	command.Stderr = &output
	scanned := make(chan string)
	events := make(chan string, 4096)
	stopRelay := make(chan struct{})
	defer close(stopRelay)
	go relayOutput(scanned, events, stopRelay)
	go scanOutput(stdout, &output, scanned)
	started := time.Now()
	if err := command.Start(); err != nil {
		return devSample{}, err
	}
	stop := func() error {
		_ = command.Process.Signal(os.Interrupt)
		done := make(chan error, 1)
		go func() { done <- command.Wait() }()
		select {
		case waitErr := <-done:
			if waitErr != nil {
				return fmt.Errorf("Dev exited after interrupt: %w\n%s", waitErr, output.String())
			}
			return nil
		case <-time.After(10 * time.Second):
			_ = command.Process.Kill()
			return fmt.Errorf("Dev did not stop after interrupt\n%s", output.String())
		}
	}
	fail := func(cause error) (devSample, error) {
		_ = stop()
		return devSample{}, fmt.Errorf("%w\n%s", cause, output.String())
	}
	if err := waitForLines(backendLog, beforeBackends+1, 45*time.Second); err != nil {
		return fail(fmt.Errorf("startup: %w", err))
	}
	if err := waitForLines(frontendLog, beforeFrontends+1, 5*time.Second); err != nil {
		return fail(fmt.Errorf("frontend startup: %w", err))
	}
	if err := waitForOutput(events, "Backend built and started", 5*time.Second); err != nil {
		return fail(fmt.Errorf("controller startup: %w", err))
	}
	startup := elapsedMS(started)

	mainSource := backendSource(generation + 1)
	rebuildStarted := time.Now()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(mainSource), 0o644); err != nil {
		return fail(err)
	}
	if err := waitForLines(backendLog, beforeBackends+2, 30*time.Second); err != nil {
		return fail(fmt.Errorf("backend rebuild: %w", err))
	}
	rebuild := elapsedMS(rebuildStarted)

	noopStarted := time.Now()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(mainSource), 0o644); err != nil {
		return fail(err)
	}
	if err := waitForOutput(events, "Build is current; backend unchanged", 15*time.Second); err != nil {
		return fail(fmt.Errorf("no-op event: %w", err))
	}
	noop := elapsedMS(noopStarted)

	drain(events)
	frontendBackends := lineCount(backendLog)
	frontendStarted := time.Now()
	frontendSource := fmt.Sprintf("export const generation = %d;\n", generation)
	if err := os.WriteFile(filepath.Join(root, "frontend", "src.js"), []byte(frontendSource), 0o644); err != nil {
		return fail(err)
	}
	timer := time.NewTimer(quietWindow)
	defer timer.Stop()
	for {
		select {
		case line := <-events:
			if strings.Contains(line, "Build finished") || strings.Contains(line, "Backend rebuilt") {
				return fail(fmt.Errorf("frontend-only edit triggered Wails work: %s", line))
			}
		case <-timer.C:
			if lineCount(backendLog) != frontendBackends {
				return fail(errors.New("frontend-only edit replaced the backend"))
			}
			frontendQuiet := elapsedMS(frontendStarted)
			if err := stop(); err != nil {
				return devSample{}, err
			}
			return devSample{StartupMS: startup, BackendRebuildMS: rebuild, NoopEventMS: noop, FrontendQuietMS: frontendQuiet}, nil
		}
	}
}

func scanOutput(reader io.Reader, output *lockedBuffer, events chan<- string) {
	defer close(events)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		output.WriteString(line + "\n")
		events <- line
	}
}

func relayOutput(source <-chan string, destination chan<- string, done <-chan struct{}) {
	var pending []string
	for source != nil || len(pending) != 0 {
		var output chan<- string
		var next string
		if len(pending) != 0 {
			output = destination
			next = pending[0]
		}
		select {
		case line, open := <-source:
			if !open {
				source = nil
				continue
			}
			pending = append(pending, line)
		case output <- next:
			pending = pending[1:]
		case <-done:
			return
		}
	}
}

func waitForOutput(events <-chan string, substring string, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case line := <-events:
			if strings.Contains(line, substring) {
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("timed out waiting for output containing %q", substring)
		}
	}
}

func drain(events <-chan string) {
	for {
		select {
		case <-events:
		default:
			return
		}
	}
}

func waitForLines(path string, count int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if lineCount(path) >= count {
			return nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for %d lines in %s", count, path)
}

func lineCount(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return bytes.Count(data, []byte{'\n'})
}

func writeFixture(root string) error {
	files := map[string]struct {
		mode os.FileMode
		data string
	}{
		"go.mod":                {0o644, "module example.com/devbenchmark\n\ngo 1.24\n"},
		"main.go":               {0o644, backendSource(0)},
		"frontend/package.json": {0o644, "{}\n"},
		"frontend/src.js":       {0o644, "export const generation = 0;\n"},
		"frontend/server.mjs": {0o644, `import fs from "node:fs";
import net from "node:net";
fs.appendFileSync(process.env.WAILS_BENCH_FRONTEND_LOG, process.pid + "\n");
const server = net.createServer(socket => socket.end());
server.listen(Number(process.env.WAILS_VITE_PORT), "127.0.0.1");
for (const signal of ["SIGINT", "SIGTERM"]) process.on(signal, () => server.close(() => process.exit(0)));
`},
		"tools/npm": {0o755, `#!/bin/sh
set -eu
case "${1:-}:${2:-}" in
  install:) mkdir -p node_modules ;;
  run:build) mkdir -p dist; printf 'benchmark\n' > dist/index.html ;;
  run:dev) exec node server.mjs ;;
  *) printf 'unsupported npm invocation: %s\n' "$*" >&2; exit 2 ;;
esac
`},
		"wails.hcl": {0o644, `version = 3
project {
  name = "dev-benchmark"
  product_name = "Dev Benchmark"
  identifier = "com.example.devbenchmark"
  version = "1.0.0"
}
frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "build"]
  dev = ["npm", "run", "dev"]
  output = "frontend/dist"
}
dev {
  debounce_ms = 25
  watch = ["**/*.go", "wails.hcl"]
  exclude = [".git", ".wails", "bin", "node_modules", "frontend"]
  use_git_ignore = false
  grace_period_ms = 100
}
`},
	}
	for relative, file := range files {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(file.data), file.mode); err != nil {
			return err
		}
	}
	return nil
}

func backendSource(generation int) string {
	return fmt.Sprintf(`package main

import (
  "os"
  "os/signal"
  "strconv"
  "syscall"
)

func main() {
  _ = generation()
  file, _ := os.OpenFile(os.Getenv("WAILS_BENCH_BACKEND_LOG"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
  _, _ = file.WriteString(strconv.Itoa(os.Getpid()) + "\n")
  _ = file.Close()
  stopped := make(chan os.Signal, 1)
  signal.Notify(stopped, os.Interrupt, syscall.SIGTERM)
  <-stopped
}

func generation() int { return %d }
`, generation)
}

func reservePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	check(err)
	port := listener.Addr().(*net.TCPAddr).Port
	check(listener.Close())
	return port
}

func sampleValues(samples []devSample, value func(devSample) float64) []float64 {
	result := make([]float64, len(samples))
	for index, sample := range samples {
		result[index] = value(sample)
	}
	return result
}

func median(values []float64) float64 {
	ordered := append([]float64(nil), values...)
	for left := 1; left < len(ordered); left++ {
		for right := left; right > 0 && ordered[right] < ordered[right-1]; right-- {
			ordered[right], ordered[right-1] = ordered[right-1], ordered[right]
		}
	}
	middle := len(ordered) / 2
	if len(ordered)%2 == 1 {
		return ordered[middle]
	}
	return (ordered[middle-1] + ordered[middle]) / 2
}

func elapsedMS(start time.Time) float64 { return float64(time.Since(start).Microseconds()) / 1000 }

func check(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

func fatalf(format string, values ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", values...)
	os.Exit(1)
}
